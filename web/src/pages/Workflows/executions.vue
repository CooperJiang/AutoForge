<template>
  <div class="space-y-6">
    <!-- 顶部面包屑 -->
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-3">
        <BaseButton size="sm" variant="ghost" @click="handleBack">
          <ArrowLeft class="w-4 h-4 mr-1" />
          返回工作流
        </BaseButton>
        <div class="h-6 w-px bg-slate-200"></div>
        <div>
          <h2 class="text-lg font-semibold text-slate-900">执行历史</h2>
          <p class="text-sm text-slate-500">{{ workflowName }}</p>
        </div>
      </div>
      <BaseButton size="sm" @click="handleRefresh">
        <RefreshCw :class="['w-4 h-4 mr-1', refreshing && 'animate-spin']" />
        刷新
      </BaseButton>
    </div>

    <!-- 筛选器 -->
    <div class="bg-white rounded-lg border border-slate-200 p-4">
      <div class="flex items-center gap-3">
        <BaseSelect
          v-model="statusFilter"
          :options="statusOptions"
          class="w-40"
        />
        <BaseInput
          v-model="searchQuery"
          placeholder="搜索执行ID..."
          class="w-64"
        />
      </div>
    </div>

    <!-- 执行列表 -->
    <div class="space-y-3">
      <div
        v-for="execution in filteredExecutions"
        :key="execution.id"
        @click="handleViewExecution(execution)"
        class="bg-white rounded-lg border-2 border-slate-200 p-4 hover:border-green-400 transition-all cursor-pointer"
      >
        <div class="flex items-start justify-between">
          <div class="flex-1">
            <div class="flex items-center gap-3 mb-2">
              <!-- 状态标签 -->
              <span
                :class="[
                  'px-2.5 py-1 rounded-full text-xs font-semibold',
                  getStatusClass(execution.status)
                ]"
              >
                <span :class="['inline-block w-1.5 h-1.5 rounded-full mr-1.5', getStatusDotClass(execution.status)]"></span>
                {{ getStatusText(execution.status) }}
              </span>

              <!-- 执行ID -->
              <span class="text-xs font-mono text-slate-500">{{ execution.id }}</span>

              <!-- 触发方式 -->
              <span class="text-xs text-slate-500 flex items-center gap-1">
                <component :is="getTriggerIcon(execution.trigger.type)" class="w-3.5 h-3.5" />
                {{ getTriggerText(execution.trigger.type) }}
              </span>
            </div>

            <!-- 时间信息 -->
            <div class="flex items-center gap-4 text-sm text-slate-600">
              <div class="flex items-center gap-1">
                <Clock class="w-4 h-4" />
                开始：{{ formatTime(execution.startTime) }}
              </div>
              <div v-if="execution.endTime" class="flex items-center gap-1">
                <Timer class="w-4 h-4" />
                耗时：{{ formatDuration(execution.startTime, execution.endTime) }}
              </div>
              <div v-else class="text-blue-600 font-medium">
                执行中...
              </div>
            </div>

            <!-- 节点执行进度 -->
            <div class="mt-3 flex items-center gap-2">
              <div class="flex-1 bg-slate-100 rounded-full h-2 overflow-hidden">
                <div
                  :style="{ width: getProgress(execution) + '%' }"
                  :class="[
                    'h-full transition-all',
                    execution.status === 'success' ? 'bg-green-500' :
                    execution.status === 'failed' ? 'bg-red-500' :
                    execution.status === 'running' ? 'bg-blue-500' :
                    'bg-slate-400'
                  ]"
                ></div>
              </div>
              <span class="text-xs text-slate-500">
                {{ getCompletedNodeCount(execution) }}/{{ execution.nodeExecutions.length }}
              </span>
            </div>

            <!-- 错误信息 -->
            <div v-if="execution.error" class="mt-2 text-sm text-red-600 bg-red-50 rounded px-2 py-1">
              {{ execution.error }}
            </div>
          </div>

          <!-- 查看详情按钮 -->
          <ChevronRight class="w-5 h-5 text-slate-400 flex-shrink-0 ml-4" />
        </div>
      </div>

      <!-- 空状态 -->
      <div v-if="filteredExecutions.length === 0" class="text-center py-16">
        <Activity class="w-16 h-16 mx-auto mb-4 text-slate-300" />
        <p class="text-slate-500">暂无执行记录</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  ArrowLeft,
  RefreshCw,
  Clock,
  Timer,
  ChevronRight,
  Activity,
  Play,
  Webhook,
  MousePointerClick
} from 'lucide-vue-next'
import BaseButton from '@/components/BaseButton'
import BaseSelect from '@/components/BaseSelect'
import BaseInput from '@/components/BaseInput'
import type { WorkflowExecution } from '@/types/workflow'
import { message } from '@/utils/message'

const router = useRouter()
const route = useRoute()

const workflowId = route.params.id as string
const workflowName = ref('我的工作流') // TODO: 从API获取
const statusFilter = ref('all')
const searchQuery = ref('')
const refreshing = ref(false)

// Mock 数据
const executions = ref<WorkflowExecution[]>([
  {
    id: 'exec_1234567890',
    workflowId: workflowId,
    status: 'success',
    startTime: new Date(Date.now() - 3600000).toISOString(),
    endTime: new Date(Date.now() - 3500000).toISOString(),
    trigger: {
      type: 'schedule'
    },
    nodeExecutions: [
      {
        nodeId: 'node_1',
        nodeName: 'HTTP请求',
        status: 'success',
        startTime: new Date(Date.now() - 3600000).toISOString(),
        endTime: new Date(Date.now() - 3550000).toISOString(),
        output: { status: 200, data: { result: 'ok' } }
      },
      {
        nodeId: 'node_2',
        nodeName: '发送邮件',
        status: 'success',
        startTime: new Date(Date.now() - 3550000).toISOString(),
        endTime: new Date(Date.now() - 3500000).toISOString(),
        output: { success: true, messageId: 'msg_123' }
      }
    ]
  },
  {
    id: 'exec_0987654321',
    workflowId: workflowId,
    status: 'running',
    startTime: new Date(Date.now() - 60000).toISOString(),
    trigger: {
      type: 'webhook',
      data: { user: 'test' }
    },
    nodeExecutions: [
      {
        nodeId: 'node_1',
        nodeName: 'HTTP请求',
        status: 'success',
        startTime: new Date(Date.now() - 60000).toISOString(),
        endTime: new Date(Date.now() - 50000).toISOString(),
        output: { status: 200 }
      },
      {
        nodeId: 'node_2',
        nodeName: '发送邮件',
        status: 'running',
        startTime: new Date(Date.now() - 50000).toISOString()
      }
    ]
  },
  {
    id: 'exec_1111222233',
    workflowId: workflowId,
    status: 'failed',
    startTime: new Date(Date.now() - 7200000).toISOString(),
    endTime: new Date(Date.now() - 7150000).toISOString(),
    trigger: {
      type: 'manual'
    },
    error: '节点执行失败：HTTP请求超时',
    nodeExecutions: [
      {
        nodeId: 'node_1',
        nodeName: 'HTTP请求',
        status: 'failed',
        startTime: new Date(Date.now() - 7200000).toISOString(),
        endTime: new Date(Date.now() - 7150000).toISOString(),
        error: '请求超时'
      }
    ]
  }
])

const statusOptions = [
  { label: '全部状态', value: 'all' },
  { label: '运行中', value: 'running' },
  { label: '成功', value: 'success' },
  { label: '失败', value: 'failed' },
  { label: '已取消', value: 'cancelled' }
]

// 过滤执行记录
const filteredExecutions = computed(() => {
  let result = executions.value

  // 状态筛选
  if (statusFilter.value !== 'all') {
    result = result.filter(e => e.status === statusFilter.value)
  }

  // 搜索筛选
  if (searchQuery.value) {
    result = result.filter(e =>
      e.id.toLowerCase().includes(searchQuery.value.toLowerCase())
    )
  }

  return result
})

const handleBack = () => {
  router.push('/workflows')
}

const handleRefresh = async () => {
  refreshing.value = true
  // TODO: 调用API刷新数据
  await new Promise(resolve => setTimeout(resolve, 1000))
  refreshing.value = false
  message.success('刷新成功')
}

const handleViewExecution = (execution: WorkflowExecution) => {
  router.push(`/workflows/${workflowId}/executions/${execution.id}`)
}

// 状态相关
const getStatusClass = (status: string) => {
  const classes = {
    running: 'bg-blue-100 text-blue-700',
    success: 'bg-green-100 text-green-700',
    failed: 'bg-red-100 text-red-700',
    cancelled: 'bg-slate-100 text-slate-700'
  }
  return classes[status as keyof typeof classes] || classes.cancelled
}

const getStatusDotClass = (status: string) => {
  const classes = {
    running: 'bg-blue-500 animate-pulse',
    success: 'bg-green-500',
    failed: 'bg-red-500',
    cancelled: 'bg-slate-500'
  }
  return classes[status as keyof typeof classes] || classes.cancelled
}

const getStatusText = (status: string) => {
  const texts = {
    running: '运行中',
    success: '成功',
    failed: '失败',
    cancelled: '已取消'
  }
  return texts[status as keyof typeof texts] || '未知'
}

// 触发器相关
const getTriggerIcon = (type: string) => {
  const icons = {
    schedule: Clock,
    webhook: Webhook,
    manual: MousePointerClick
  }
  return icons[type as keyof typeof icons] || Play
}

const getTriggerText = (type: string) => {
  const texts = {
    schedule: '定时触发',
    webhook: 'Webhook',
    manual: '手动触发'
  }
  return texts[type as keyof typeof texts] || '未知'
}

// 进度计算
const getProgress = (execution: WorkflowExecution) => {
  const completed = getCompletedNodeCount(execution)
  const total = execution.nodeExecutions.length
  return total > 0 ? Math.round((completed / total) * 100) : 0
}

const getCompletedNodeCount = (execution: WorkflowExecution) => {
  return execution.nodeExecutions.filter(n =>
    n.status === 'success' || n.status === 'failed' || n.status === 'skipped'
  ).length
}

// 时间格式化
const formatTime = (time: string) => {
  const date = new Date(time)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

const formatDuration = (start: string, end: string) => {
  const duration = new Date(end).getTime() - new Date(start).getTime()
  const seconds = Math.floor(duration / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)

  if (hours > 0) {
    return `${hours}小时${minutes % 60}分钟`
  } else if (minutes > 0) {
    return `${minutes}分钟${seconds % 60}秒`
  } else {
    return `${seconds}秒`
  }
}

onMounted(() => {
  // TODO: 加载执行历史
})
</script>
