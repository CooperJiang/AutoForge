<template>
  <div>
    
    <div class="flex gap-3 mb-6 items-center">
      <div class="flex-shrink-0" style="width: 200px">
        <BaseInput v-model="filters.user_id" placeholder="用户ID" />
      </div>
      <div class="flex-shrink-0" style="width: 200px">
        <BaseInput v-model="filters.task_id" placeholder="任务ID" />
      </div>
      <div class="flex-shrink-0" style="width: 150px">
        <BaseSelect
          v-model="filters.status"
          :options="[
            { label: '全部状态', value: '' },
            { label: '成功', value: 'success' },
            { label: '失败', value: 'failed' },
          ]"
          placeholder="全部状态"
        />
      </div>
      <BaseButton @click="loadExecutions" variant="primary" class="flex-shrink-0">
        搜索
      </BaseButton>
    </div>

    
    <div class="overflow-x-auto">
      <table class="w-full">
        <thead>
          <tr class="border-b-2 border-border-primary text-left">
            <th class="pb-3 text-sm font-semibold text-text-secondary">执行ID</th>
            <th class="pb-3 text-sm font-semibold text-text-secondary">任务名称</th>
            <th class="pb-3 text-sm font-semibold text-text-secondary">用户ID</th>
            <th class="pb-3 text-sm font-semibold text-text-secondary">状态</th>
            <th class="pb-3 text-sm font-semibold text-text-secondary">HTTP码</th>
            <th class="pb-3 text-sm font-semibold text-text-secondary">耗时</th>
            <th class="pb-3 text-sm font-semibold text-text-secondary">开始时间</th>
            <th class="pb-3 text-sm font-semibold text-text-secondary">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="exec in executions"
            :key="exec.id"
            class="border-b border-border-primary hover:bg-bg-hover transition-colors"
          >
            <td class="py-3 text-sm text-text-primary font-mono text-xs">
              {{ truncateId(exec.id, 8) }}
            </td>
            <td
              class="py-3 text-sm text-text-primary font-medium max-w-[200px] truncate"
              :title="exec.task?.name || exec.task_id"
            >
              {{ exec.task?.name || '未知任务' }}
            </td>
            <td class="py-3 text-sm text-text-primary font-mono">{{ maskUserId(exec.user_id) }}</td>
            <td class="py-3">
              <span
                :class="[
                  'px-2 py-1 text-xs font-medium rounded-full border',
                  exec.status === 'success'
                    ? 'bg-green-500/10 text-green-600 dark:text-green-400 border-green-500/20'
                    : 'bg-red-500/10 text-red-600 dark:text-red-400 border-red-500/20',
                ]"
              >
                {{ exec.status === 'success' ? '成功' : '失败' }}
              </span>
            </td>
            <td class="py-3">
              <span
                :class="[
                  'px-2 py-1 text-xs font-semibold rounded border',
                  exec.response_status >= 200 && exec.response_status < 300
                    ? 'bg-green-500/10 text-green-600 dark:text-green-400 border-green-500/20'
                    : exec.response_status >= 400
                      ? 'bg-red-500/10 text-red-600 dark:text-red-400 border-red-500/20'
                      : 'bg-yellow-500/10 text-yellow-600 dark:text-yellow-400 border-yellow-500/20',
                ]"
              >
                {{ exec.response_status || 'N/A' }}
              </span>
            </td>
            <td class="py-3 text-sm text-text-secondary">{{ exec.duration_ms }}ms</td>
            <td class="py-3 text-sm text-text-secondary">{{ formatTime(exec.started_at) }}</td>
            <td class="py-3">
              <div class="flex items-center gap-2">
                <button
                  @click="viewExecution(exec)"
                  class="p-1.5 bg-primary-light hover:bg-primary-active text-primary rounded-lg transition-colors"
                  title="查看详情"
                >
                  <Eye :size="16" />
                </button>
                <button
                  @click="deleteExecution(exec)"
                  class="p-1.5 bg-error-light hover:bg-error-active text-error rounded-lg transition-colors"
                  title="删除记录"
                >
                  <Trash2 :size="16" />
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    
    <div class="flex justify-between items-center mt-6">
      <div class="text-sm text-text-secondary">共 {{ total }} 条记录</div>
      <div class="flex gap-2">
        <button
          @click="prevPage"
          :disabled="currentPage === 1"
          class="px-4 py-2 bg-bg-tertiary text-text-secondary text-sm font-medium rounded-lg hover:bg-bg-tertiary disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          上一页
        </button>
        <span class="px-4 py-2 text-sm text-text-secondary">
          {{ currentPage }} / {{ totalPages }}
        </span>
        <button
          @click="nextPage"
          :disabled="currentPage >= totalPages"
          class="px-4 py-2 bg-bg-tertiary text-text-secondary text-sm font-medium rounded-lg hover:bg-bg-tertiary disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          下一页
        </button>
      </div>
    </div>

    
    <ExecutionDetailDialog v-model="showExecutionDetail" :execution="selectedExecution" />

    
    <Dialog
      v-model="showDeleteDialog"
      title="确认删除"
      :message="`确定要删除执行记录吗？此操作不可恢复！`"
      confirm-text="删除"
      cancel-text="取消"
      confirm-variant="danger"
      @confirm="confirmDelete"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { Eye, Trash2 } from 'lucide-vue-next'
import * as adminApi from '@/api/admin'
import { message } from '@/utils/message'
import { maskUserId, formatTime, truncateId } from '@/utils/format'
import BaseInput from '@/components/BaseInput'
import BaseSelect from '@/components/BaseSelect'
import BaseButton from '@/components/BaseButton'
import Dialog from '@/components/Dialog'
import ExecutionDetailDialog from '@/components/ExecutionDetailDialog'

// 执行记录列表
const executions = ref<any[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

// 筛选条件
const filters = ref({
  user_id: '',
  task_id: '',
  status: '',
})

// 详情对话框
const showExecutionDetail = ref(false)
const selectedExecution = ref<any>(null)

// 删除确认对话框
const showDeleteDialog = ref(false)
const executionToDelete = ref<any>(null)

// 计算总页数
const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

// 加载执行记录
const loadExecutions = async () => {
  try {
    const res = await adminApi.getExecutions({
      page: currentPage.value,
      page_size: pageSize.value,
      ...filters.value,
    })
    executions.value = res.data.executions || []
    total.value = res.data.total
  } catch (error: any) {
    message.error(error.response?.data?.message || '加载失败')
  }
}

/**
 * 查看执行详情
 */
const viewExecution = (execution: any) => {
  selectedExecution.value = execution
  showExecutionDetail.value = true
}

/**
 * 删除执行记录
 */
const deleteExecution = (execution: any) => {
  executionToDelete.value = execution
  showDeleteDialog.value = true
}

const confirmDelete = async () => {
  if (!executionToDelete.value) return

  try {
    await adminApi.deleteExecution(executionToDelete.value.id)
    message.success('删除成功')
    executionToDelete.value = null
    // 重新加载当前页数据
    loadExecutions()
  } catch (error: any) {
    message.error(error.message || '删除失败')
  }
}

// 上一页
const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
    loadExecutions()
  }
}

// 下一页
const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
    loadExecutions()
  }
}

onMounted(() => {
  loadExecutions()
})
</script>
