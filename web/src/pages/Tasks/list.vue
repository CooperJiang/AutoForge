<template>
  <!-- Loading -->
  <div v-if="loading" class="flex justify-center items-center py-20">
    <div class="text-text-tertiary">加载中...</div>
  </div>

  <!-- Main Content -->
  <main v-else class="max-w-7xl mx-auto px-4 py-4">
    <!-- 顶部操作栏 -->
    <div class="mb-4 flex items-center justify-between">
      <h1 class="text-2xl font-bold text-text-primary">定时任务管理</h1>
      <BaseButton variant="primary" @click="goToCreateTask">
        ➕ 创建任务
      </BaseButton>
    </div>

    <!-- 任务列表 -->
    <div class="bg-bg-elevated border-2 border-border-primary rounded-lg shadow-sm">
      <Table :data="tasks" :loading="loading">
        <template #header>
          <th class="w-[15%]">任务名称</th>
          <th class="w-[20%]">描述</th>
          <th class="w-[12%]">工具</th>
          <th class="w-[8%]">调度类型</th>
          <th class="w-[10%]">下次执行</th>
          <th class="w-[8%]">状态</th>
          <th class="w-[27%] text-center">操作</th>
        </template>
        <template #body>
          <tr v-for="task in tasks" :key="task.id" class="hover:bg-bg-hover">
            <td class="font-medium w-[15%]" :title="task.name">{{ task.name }}</td>
            <td class="text-text-secondary w-[20%]" :title="task.description || '-'">{{ task.description || '-' }}</td>
            <td class="w-[12%]">
              <span class="inline-block px-2 py-1 text-xs rounded-md bg-primary-light text-primary border border-primary max-w-full overflow-hidden text-ellipsis" :title="task.tool_code">
                {{ task.tool_code }}
              </span>
            </td>
            <td class="w-[8%]">{{ getScheduleTypeName(task.schedule_type) }}</td>
            <td class="text-text-secondary w-[10%]">
              <span v-if="task.next_run_time">{{ formatTime(task.next_run_time) }}</span>
              <span v-else class="text-text-placeholder">-</span>
            </td>
            <td class="w-[8%]">
              <span
                class="inline-block px-2 py-1 text-xs rounded-md whitespace-nowrap"
                :class="task.enabled
                  ? 'bg-emerald-50 text-emerald-700 border border-emerald-200'
                  : 'bg-bg-hover text-text-secondary border border-border-primary'"
              >
                {{ task.enabled ? '✓ 启用' : '✗ 禁用' }}
              </span>
            </td>
            <td class="w-[27%]">
              <div class="flex gap-1 justify-center">
                <BaseButton size="sm" variant="ghost" @click="handleTrigger(task)">
                  执行
                </BaseButton>
                <BaseButton size="sm" variant="ghost" @click="handleViewExecutions(task)">
                  记录
                </BaseButton>
                <BaseButton size="sm" variant="ghost" @click="handleEdit(task)">
                  编辑
                </BaseButton>
                <BaseButton
                  size="sm"
                  :variant="task.enabled ? 'secondary' : 'success'"
                  @click="handleToggleStatus(task)"
                >
                  {{ task.enabled ? '禁用' : '启用' }}
                </BaseButton>
                <BaseButton size="sm" variant="danger" @click="handleDelete(task)">
                  删除
                </BaseButton>
              </div>
            </td>
          </tr>
        </template>
      </Table>

      <!-- 分页 -->
      <Pagination
        v-model:current-page="currentPage"
        :page-size="pageSize"
        :total="total"
        @change="loadTasks"
      />
    </div>

    <!-- 执行记录对话框 -->
    <TaskDetailDialog
      v-if="currentTaskId"
      v-model="showTaskDetailDialog"
      :task-id="currentTaskId"
      @close="showTaskDetailDialog = false"
    />

    <!-- 删除确认对话框 -->
    <Dialog
      v-model="showDeleteDialog"
      title="删除任务"
      :message="`确定要删除任务【${deletingTask?.name}】吗？删除后将无法恢复。`"
      confirm-text="删除"
      cancel-text="取消"
      confirm-variant="danger"
      @confirm="confirmDelete"
    />
  </main>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { message } from '@/utils/message'
import type { Task } from '@/api/task'
import {
  getTaskList,
  deleteTask,
  enableTask,
  disableTask,
  triggerTask,
} from '@/api/task'
import BaseButton from '@/components/BaseButton'
import Table from '@/components/Table'
import Pagination from '@/components/Pagination'
import TaskDetailDialog from '@/components/TaskDetailDialog'
import Dialog from '@/components/Dialog'

const router = useRouter()

const tasks = ref<Task[]>([])
const loading = ref(true)
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

const showTaskDetailDialog = ref(false)
const currentTaskId = ref('')

const showDeleteDialog = ref(false)
const deletingTask = ref<Task | null>(null)

const scheduleTypeMap: Record<string, string> = {
  daily: '每天',
  weekly: '每周',
  monthly: '每月',
  hourly: '每小时',
  interval: '间隔',
  cron: 'Cron',
}

onMounted(async () => {
  await loadTasks()
})

const loadTasks = async () => {
  loading.value = true
  try {
    const result = await getTaskList(currentPage.value, pageSize.value)
    tasks.value = result.items
    total.value = result.total
  } catch (error: any) {
    message.error(error.message || '加载任务列表失败')
  } finally {
    loading.value = false
  }
}


const getScheduleTypeName = (type: string) => {
  return scheduleTypeMap[type] || type
}

const formatTime = (timestamp: number) => {
  const date = new Date(timestamp * 1000)
  return date.toLocaleString('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const goToCreateTask = () => {
  router.push('/tasks/create')
}

const handleEdit = (task: Task) => {
  router.push(`/tasks/${task.id}/edit`)
}

const handleToggleStatus = async (task: Task) => {
  try {
    if (task.enabled) {
      await disableTask(task.id)
      message.success('已禁用')
    } else {
      await enableTask(task.id)
      message.success('已启用')
    }
    await loadTasks()
  } catch (error: any) {
    message.error(error.message || '操作失败')
  }
}

const handleDelete = (task: Task) => {
  deletingTask.value = task
  showDeleteDialog.value = true
}

const confirmDelete = async () => {
  if (!deletingTask.value) return

  try {
    await deleteTask(deletingTask.value.id)
    message.success('删除成功')
    await loadTasks()
  } catch (error: any) {
    message.error(error.message || '删除失败')
  } finally {
    deletingTask.value = null
  }
}

const handleTrigger = async (task: Task) => {
  try {
    await triggerTask(task.id)
    message.success('任务已触发执行')
  } catch (error: any) {
    message.error(error.message || '触发失败')
  }
}

const handleViewExecutions = (task: Task) => {
  currentTaskId.value = task.id
  showTaskDetailDialog.value = true
}
</script>
