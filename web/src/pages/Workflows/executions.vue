<template>
  <div class="space-y-6">
    
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-3">
        <div>
          <h2 class="text-lg font-semibold text-text-primary">执行历史</h2>
          <p class="text-sm text-text-tertiary">{{ workflowName }}</p>
        </div>
      </div>
      <div class="flex items-center gap-2">
        <BaseButton size="sm" @click="handleRefresh">
          <RefreshCw :class="['w-4 h-4 mr-1', refreshing && 'animate-spin']" />
          刷新
        </BaseButton>
        <ExecuteWorkflowButton
          :workflow-id="workflowId"
          :workflow="workflowInfo"
          @executed="onExecuted"
        />
      </div>
    </div>

    
    <div
      v-if="
        workflowInfo &&
        workflowInfo.enabled &&
        workflowInfo.schedule_type &&
        workflowInfo.next_run_time
      "
      class="bg-primary-light border border-primary rounded-lg p-4"
    >
      <div class="flex items-center gap-3">
        <Clock class="w-5 h-5 text-primary" />
        <div class="flex-1">
          <div class="text-sm font-semibold text-primary">定时任务已启用</div>
          <div class="text-xs text-primary mt-1">
            下次执行时间：{{ formatNextRunTime(workflowInfo.next_run_time) }}
            <span class="ml-2 px-2 py-0.5 bg-primary text-white rounded font-mono">
              <CountdownDisplay
                :timestamp="workflowInfo.next_run_time"
                @finish="handleCountdownFinish"
              />
            </span>
          </div>
          <div class="text-xs text-text-secondary mt-1">
            {{ getScheduleDescription(workflowInfo.schedule_type, workflowInfo.schedule_value) }}
          </div>
        </div>
      </div>
    </div>

    
    <div class="bg-bg-elevated rounded-lg border border-border-primary p-4">
      <div class="flex items-center gap-3">
        <BaseSelect v-model="statusFilter" :options="statusOptions" class="w-40" />
        <BaseInput v-model="searchQuery" placeholder="搜索执行ID..." class="w-64" />
      </div>
    </div>

    
    <div v-if="loading" class="flex justify-center items-center py-20">
      <div class="text-text-tertiary">加载中...</div>
    </div>

    
    <div v-else class="space-y-3">
      <div
        v-for="execution in filteredExecutions"
        :key="execution.id"
        class="bg-bg-elevated rounded-lg border-2 border-border-primary p-4 hover:border-green-400 transition-all"
      >
        <div class="flex items-start justify-between">
          <div class="flex-1 cursor-pointer" @click="handleViewExecution(execution)">
            <div class="flex items-center gap-3 mb-2">
              
              <span
                :class="[
                  'px-2.5 py-1 rounded-full text-xs font-semibold',
                  getStatusClass(execution.status),
                ]"
              >
                <span
                  :class="[
                    'inline-block w-1.5 h-1.5 rounded-full mr-1.5',
                    getStatusDotClass(execution.status),
                  ]"
                ></span>
                {{ getStatusText(execution.status) }}
              </span>

              
              <span class="text-xs font-mono text-text-tertiary">{{ execution.id }}</span>

              
              <span class="text-xs text-text-tertiary flex items-center gap-1">
                <component :is="getTriggerIcon(execution.trigger_type)" class="w-3.5 h-3.5" />
                {{ getTriggerText(execution.trigger_type) }}
              </span>
            </div>

            
            <div class="flex items-center gap-4 text-sm text-text-secondary">
              <div class="flex items-center gap-1">
                <Clock class="w-4 h-4" />
                开始：{{ formatTime(execution.start_time) }}
              </div>
              <div v-if="execution.end_time" class="flex items-center gap-1">
                <Timer class="w-4 h-4" />
                耗时：{{ formatDurationMs(execution.duration_ms) }}
              </div>
              <div v-else class="text-primary font-medium">执行中...</div>
            </div>

            
            <div class="mt-3 flex items-center gap-2">
              <div class="flex-1 bg-bg-tertiary rounded-full h-2 overflow-hidden">
                <div
                  :style="{ width: getProgress(execution) + '%' }"
                  :class="[
                    'h-full transition-all',
                    execution.status === 'success'
                      ? 'bg-green-500'
                      : execution.status === 'failed'
                        ? 'bg-red-500'
                        : execution.status === 'running'
                          ? 'bg-[var(--color-primary)]'
                          : 'bg-slate-400',
                  ]"
                ></div>
              </div>
              <span class="text-xs text-text-tertiary">
                {{ execution.success_nodes + execution.failed_nodes + execution.skipped_nodes }}/{{
                  execution.total_nodes
                }}
              </span>
            </div>

            
            <div
              v-if="execution.error"
              class="mt-2 text-sm text-red-600 bg-red-50 rounded px-2 py-1"
            >
              {{ execution.error }}
            </div>
          </div>

          
          <div class="flex items-center gap-2 flex-shrink-0 ml-4">
            <button
              @click.stop="handleDeleteExecution(execution)"
              class="w-8 h-8 flex items-center justify-center rounded-lg text-text-tertiary hover:text-error hover:bg-error-light transition-all"
              title="删除记录"
            >
              <Trash2 class="w-4 h-4" />
            </button>
            <ChevronRight
              class="w-5 h-5 text-text-placeholder cursor-pointer"
              @click="handleViewExecution(execution)"
            />
          </div>
        </div>
      </div>

      
      <div v-if="filteredExecutions.length === 0" class="text-center py-16">
        <Activity class="w-16 h-16 mx-auto mb-4 text-slate-300" />
        <p class="text-text-tertiary">暂无执行记录</p>
      </div>
    </div>

    
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
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  RefreshCw,
  Clock,
  Timer,
  ChevronRight,
  Activity,
  Play,
  Webhook,
  MousePointerClick,
  Trash2,
} from 'lucide-vue-next'
import BaseButton from '@/components/BaseButton'
import BaseSelect from '@/components/BaseSelect'
import BaseInput from '@/components/BaseInput'
import Dialog from '@/components/Dialog'
import CountdownDisplay from '@/components/CountdownDisplay'
import ExecuteWorkflowButton from '@/components/ExecuteWorkflowButton.vue'
import type { WorkflowExecution } from '@/types/workflow'
import { workflowApi } from '@/api/workflow'
import { message } from '@/utils/message'
import { formatTimestamp } from '@/composables/useCountdown'

const router = useRouter()
const route = useRoute()

const workflowId = route.params.id as string
const workflowName = ref('')
const workflowInfo = ref<any>(null)
const statusFilter = ref('all')
const searchQuery = ref('')
const refreshing = ref(false)
const loading = ref(false)

const executions = ref<WorkflowExecution[]>([])

// 删除确认对话框
const showDeleteDialog = ref(false)
const executionToDelete = ref<WorkflowExecution | null>(null)

const statusOptions = [
  { label: '全部状态', value: 'all' },
  { label: '运行中', value: 'running' },
  { label: '成功', value: 'success' },
  { label: '失败', value: 'failed' },
  { label: '已取消', value: 'cancelled' },
]

// 过滤执行记录
const filteredExecutions = computed(() => {
  let result = executions.value

  // 状态筛选
  if (statusFilter.value !== 'all') {
    result = result.filter((e) => e.status === statusFilter.value)
  }

  // 搜索筛选
  if (searchQuery.value) {
    result = result.filter((e) => e.id.toLowerCase().includes(searchQuery.value.toLowerCase()))
  }

  return result
})

// 加载工作流详情
const loadWorkflow = async () => {
  try {
    const data = await workflowApi.getById(workflowId)
    workflowName.value = data.name
    workflowInfo.value = data
  } catch (error) {
    console.error('Load workflow failed:', error)
  }
}

// 加载执行历史
const loadExecutions = async () => {
  loading.value = true
  try {
    const data = await workflowApi.getExecutions(workflowId, {
      page: 1,
      page_size: 100,
      status: statusFilter.value === 'all' ? undefined : statusFilter.value,
    })
    executions.value = data.items || []
  } catch (error: any) {
    console.error('Load executions failed:', error)
    message.error(error.response?.data?.message || '加载执行历史失败')
  } finally {
    loading.value = false
  }
}

const handleRefresh = async () => {
  refreshing.value = true
  try {
    await loadExecutions()
    message.success('刷新成功')
  } catch {
    message.error('刷新失败')
  } finally {
    refreshing.value = false
  }
}

const handleViewExecution = (execution: WorkflowExecution) => {
  router.push(`/workflows/${workflowId}/executions/${execution.id}`)
}

const handleDeleteExecution = (execution: WorkflowExecution) => {
  executionToDelete.value = execution
  showDeleteDialog.value = true
}

const confirmDelete = async () => {
  if (!executionToDelete.value) return

  try {
    await workflowApi.deleteExecution(workflowId, executionToDelete.value.id)
    message.success('删除成功')
    // 从列表中移除
    executions.value = executions.value.filter((e) => e.id !== executionToDelete.value!.id)
    executionToDelete.value = null
  } catch (error: any) {
    message.error(error.message || '删除失败')
  }
}

// 状态相关
const getStatusClass = (status: string) => {
  const classes = {
    running: 'bg-primary-light text-primary',
    success: 'bg-green-100 text-green-700',
    failed: 'bg-red-100 text-red-700',
    cancelled: 'bg-bg-tertiary text-text-secondary',
  }
  return classes[status as keyof typeof classes] || classes.cancelled
}

const getStatusDotClass = (status: string) => {
  const classes = {
    running: 'bg-[var(--color-primary)] animate-pulse',
    success: 'bg-green-500',
    failed: 'bg-red-500',
    cancelled: 'bg-bg-hover0',
  }
  return classes[status as keyof typeof classes] || classes.cancelled
}

const getStatusText = (status: string) => {
  const texts = {
    running: '运行中',
    success: '成功',
    failed: '失败',
    cancelled: '已取消',
  }
  return texts[status as keyof typeof texts] || '未知'
}

// 触发器相关
const getTriggerIcon = (type: string) => {
  const icons = {
    schedule: Clock,
    scheduled: Clock,
    webhook: Webhook,
    manual: MousePointerClick,
  }
  return icons[type as keyof typeof icons] || Play
}

const getTriggerText = (type: string) => {
  const texts = {
    schedule: '定时触发',
    scheduled: '定时触发',
    webhook: 'Webhook',
    manual: '手动触发',
  }
  return texts[type as keyof typeof texts] || '未知'
}

// 进度计算
const getProgress = (execution: WorkflowExecution) => {
  const completed = execution.success_nodes + execution.failed_nodes + execution.skipped_nodes
  const total = execution.total_nodes
  return total > 0 ? Math.round((completed / total) * 100) : 0
}

// 时间格式化
const formatTime = (timestamp?: number) => {
  if (!timestamp) return '-'
  const date = new Date(timestamp * 1000)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  })
}

// 倒计时结束回调 - 自动刷新执行历史
const handleCountdownFinish = async () => {
  setTimeout(async () => {
    try {
      await Promise.all([loadWorkflow(), loadExecutions()])
      message.success('工作流已执行，列表已自动刷新')
    } catch (error) {
      console.error('Auto refresh failed:', error)
      message.error('自动刷新失败，请手动刷新')
    }
  }, 1000)
}

// 格式化下次执行时间
const formatNextRunTime = formatTimestamp

// 获取调度描述
const getScheduleDescription = (type: string, value: string) => {
  if (!type || !value) return ''

  switch (type) {
    case 'daily':
      return `每天 ${value}`
    case 'weekly':
      const [days, time] = value.split(':').slice(0, 2).join(':').split(':')
      return `每周 ${days} ${time || ''}`
    case 'monthly':
      const [day, ...timeParts] = value.split(':')
      return `每月 ${day}号 ${timeParts.join(':')}`
    case 'hourly':
      return `每小时 ${value}`
    case 'interval':
      return `每 ${value} 秒`
    case 'cron':
      return `Cron: ${value}`
    default:
      return type
  }
}

const formatDurationMs = (ms: number) => {
  const seconds = Math.floor(ms / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)

  if (hours > 0) {
    return `${hours}小时${minutes % 60}分钟${seconds % 60}秒`
  } else if (minutes > 0) {
    return `${minutes}分钟${seconds % 60}秒`
  } else if (seconds > 0) {
    return `${seconds}秒`
  } else {
    return `${ms}毫秒`
  }
}

onMounted(() => {
  loadWorkflow()
  loadExecutions()
})

const onExecuted = async () => {
  // 执行后刷新列表
  await loadExecutions()
}
</script>
