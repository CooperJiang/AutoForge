<template>
  <div>
    <!-- 搜索筛选栏 -->
    <div class="flex gap-3 mb-6 items-center">
      <div class="flex-shrink-0" style="width: 200px;">
        <BaseInput
          v-model="filters.user_id"
          placeholder="用户ID"
        />
      </div>
      <div class="flex-shrink-0" style="width: 200px;">
        <BaseInput
          v-model="filters.task_id"
          placeholder="任务ID"
        />
      </div>
      <div class="flex-shrink-0" style="width: 150px;">
        <BaseSelect
          v-model="filters.status"
          :options="[
            { label: '全部状态', value: '' },
            { label: '成功', value: 'success' },
            { label: '失败', value: 'failed' }
          ]"
          placeholder="全部状态"
        />
      </div>
      <BaseButton
        @click="loadExecutions"
        variant="primary"
        class="flex-shrink-0"
      >
        搜索
      </BaseButton>
    </div>

    <!-- 表格 -->
    <div class="overflow-x-auto">
      <table class="w-full">
        <thead>
          <tr class="border-b-2 border-slate-200 text-left">
            <th class="pb-3 text-sm font-semibold text-slate-700">执行ID</th>
            <th class="pb-3 text-sm font-semibold text-slate-700">任务名称</th>
            <th class="pb-3 text-sm font-semibold text-slate-700">用户ID</th>
            <th class="pb-3 text-sm font-semibold text-slate-700">状态</th>
            <th class="pb-3 text-sm font-semibold text-slate-700">HTTP码</th>
            <th class="pb-3 text-sm font-semibold text-slate-700">耗时</th>
            <th class="pb-3 text-sm font-semibold text-slate-700">开始时间</th>
            <th class="pb-3 text-sm font-semibold text-slate-700">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="exec in executions"
            :key="exec.id"
            class="border-b border-slate-100 hover:bg-slate-50 transition-colors"
          >
            <td class="py-3 text-sm text-slate-900 font-mono text-xs">{{ truncateId(exec.id, 8) }}</td>
            <td class="py-3 text-sm text-slate-900 font-medium max-w-[200px] truncate" :title="exec.task?.name || exec.task_id">
              {{ exec.task?.name || '未知任务' }}
            </td>
            <td class="py-3 text-sm text-slate-900 font-mono">{{ maskUserId(exec.user_id) }}</td>
            <td class="py-3">
              <span
                :class="[
                  'px-2 py-1 text-xs font-medium rounded-full',
                  exec.status === 'success'
                    ? 'bg-green-100 text-green-700'
                    : 'bg-red-100 text-red-700'
                ]"
              >
                {{ exec.status === 'success' ? '成功' : '失败' }}
              </span>
            </td>
            <td class="py-3">
              <span
                :class="[
                  'px-2 py-1 text-xs font-semibold rounded',
                  exec.response_status >= 200 && exec.response_status < 300
                    ? 'bg-green-100 text-green-700'
                    : exec.response_status >= 400
                    ? 'bg-red-100 text-red-700'
                    : 'bg-yellow-100 text-yellow-700'
                ]"
              >
                {{ exec.response_status || 'N/A' }}
              </span>
            </td>
            <td class="py-3 text-sm text-slate-600">{{ exec.duration_ms }}ms</td>
            <td class="py-3 text-sm text-slate-600">{{ formatTime(exec.started_at) }}</td>
            <td class="py-3">
              <button
                @click="viewExecution(exec)"
                class="p-1.5 bg-blue-100 hover:bg-blue-200 text-blue-700 rounded-lg transition-colors"
                title="查看详情"
              >
                <Eye :size="16" />
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 分页 -->
    <div class="flex justify-between items-center mt-6">
      <div class="text-sm text-slate-600">
        共 {{ total }} 条记录
      </div>
      <div class="flex gap-2">
        <button
          @click="prevPage"
          :disabled="currentPage === 1"
          class="px-4 py-2 bg-slate-100 text-slate-700 text-sm font-medium rounded-lg hover:bg-slate-200 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          上一页
        </button>
        <span class="px-4 py-2 text-sm text-slate-700">
          {{ currentPage }} / {{ totalPages }}
        </span>
        <button
          @click="nextPage"
          :disabled="currentPage >= totalPages"
          class="px-4 py-2 bg-slate-100 text-slate-700 text-sm font-medium rounded-lg hover:bg-slate-200 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          下一页
        </button>
      </div>
    </div>

    <!-- 执行详情对话框 -->
    <ExecutionDetailDialog
      v-model="showExecutionDetail"
      :execution="selectedExecution"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { Eye } from 'lucide-vue-next'
import * as adminApi from '@/api/admin'
import { message } from '@/utils/message'
import { maskUserId, formatTime, truncateId } from '@/utils/format'
import BaseInput from '@/components/BaseInput.vue'
import BaseSelect from '@/components/BaseSelect.vue'
import BaseButton from '@/components/BaseButton.vue'
import ExecutionDetailDialog from '@/components/ExecutionDetailDialog.vue'

// 执行记录列表
const executions = ref<any[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

// 筛选条件
const filters = ref({
  user_id: '',
  task_id: '',
  status: ''
})

// 详情对话框
const showExecutionDetail = ref(false)
const selectedExecution = ref<any>(null)

// 计算总页数
const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

// 加载执行记录
const loadExecutions = async () => {
  try {
    const res = await adminApi.getExecutions({
      page: currentPage.value,
      page_size: pageSize.value,
      ...filters.value
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
