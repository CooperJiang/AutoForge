<template>
  <div class="space-y-6">
    <!-- 加载状态 -->
    <div v-if="loading" class="flex justify-center items-center py-20">
      <div class="text-slate-500">加载中...</div>
    </div>

    <!-- 数据为空 -->
    <div v-else-if="!execution" class="flex flex-col justify-center items-center py-20">
      <div class="text-slate-500 mb-4">执行记录不存在</div>
      <BaseButton size="sm" @click="handleBack">返回列表</BaseButton>
    </div>

    <!-- 执行详情 -->
    <template v-else>
    <!-- 顶部信息 -->
    <div class="flex items-center justify-between">
      <div>
        <p class="text-sm text-slate-500 font-mono">{{ executionId }}</p>
      </div>

      <div class="flex items-center gap-2">
        <span
          :class="[
            'px-3 py-1.5 rounded-full text-sm font-semibold',
            getStatusClass(execution?.status || '')
          ]"
        >
          <span :class="['inline-block w-2 h-2 rounded-full mr-2', getStatusDotClass(execution?.status || '')]"></span>
          {{ getStatusText(execution?.status || '') }}
        </span>
      </div>
    </div>

    <!-- 执行概览 -->
    <div class="bg-white rounded-lg border border-slate-200 p-6">
      <h3 class="text-sm font-semibold text-slate-900 mb-4">执行概览</h3>
      <div class="grid grid-cols-4 gap-6">
        <div>
          <div class="text-xs text-slate-500 mb-1">触发方式</div>
          <div class="flex items-center gap-2 text-sm font-medium text-slate-900">
            <component :is="getTriggerIcon(execution?.trigger_type || '')" class="w-4 h-4" />
            {{ getTriggerText(execution?.trigger_type || '') }}
          </div>
        </div>
        <div>
          <div class="text-xs text-slate-500 mb-1">开始时间</div>
          <div class="text-sm font-medium text-slate-900">{{ formatTime(execution?.start_time) }}</div>
        </div>
        <div>
          <div class="text-xs text-slate-500 mb-1">结束时间</div>
          <div class="text-sm font-medium text-slate-900">
            {{ execution?.end_time ? formatTime(execution.end_time) : '-' }}
          </div>
        </div>
        <div>
          <div class="text-xs text-slate-500 mb-1">执行耗时</div>
          <div class="text-sm font-medium text-slate-900">
            {{ execution?.duration_ms != null ? formatDurationMs(execution.duration_ms) : '-' }}
          </div>
        </div>
      </div>
    </div>

    <!-- 节点执行详情 -->
    <div class="bg-white rounded-lg border border-slate-200">
      <div class="px-6 py-4 border-b border-slate-200">
        <h3 class="text-sm font-semibold text-slate-900">节点执行详情</h3>
      </div>

      <div class="divide-y divide-slate-200">
        <div
          v-for="(nodeLog, index) in execution?.node_logs"
          :key="nodeLog.node_id"
          class="p-6"
        >
          <div class="flex items-start gap-4">
            <!-- 序号和连接线 -->
            <div class="flex flex-col items-center">
              <div
                :class="[
                  'w-8 h-8 rounded-full flex items-center justify-center text-sm font-semibold',
                  getNodeStatusBgClass(nodeLog.status)
                ]"
              >
                {{ index + 1 }}
              </div>
              <div
                v-if="index < (execution?.node_logs.length || 0) - 1"
                class="flex-1 w-0.5 bg-slate-200 my-2 min-h-[20px]"
              ></div>
            </div>

            <!-- 节点信息 -->
            <div class="flex-1">
              <div class="flex items-center justify-between mb-2">
                <div class="flex items-center gap-2">
                  <h4 class="text-sm font-semibold text-slate-900">{{ nodeLog.node_name }}</h4>
                  <span v-if="nodeLog.tool_code" class="px-2 py-0.5 bg-blue-50 text-blue-600 text-xs font-mono rounded">{{ nodeLog.tool_code }}</span>
                  <span class="text-xs text-slate-400 font-mono">{{ nodeLog.node_id }}</span>
                </div>
                <span
                  :class="[
                    'px-2 py-0.5 rounded-full text-xs font-medium',
                    getNodeStatusClass(nodeLog.status)
                  ]"
                >
                  {{ getNodeStatusText(nodeLog.status) }}
                </span>
              </div>

              <!-- 时间信息 -->
              <div v-if="nodeLog.start_time" class="flex items-center gap-4 text-xs text-slate-600 mb-3">
                <span>开始：{{ formatTime(nodeLog.start_time) }}</span>
                <span v-if="nodeLog.duration_ms != null">
                  耗时：{{ formatDurationMs(nodeLog.duration_ms) }}
                </span>
              </div>

              <!-- 输出数据 -->
              <div v-if="nodeLog.output" class="mb-3">
                <button
                  type="button"
                  @click="toggleSection(nodeLog.node_id, 'output')"
                  class="flex items-center gap-1 text-xs font-medium text-slate-700 hover:text-slate-900 mb-1"
                >
                  <ChevronDown :class="['w-3 h-3 transition-transform', isSectionOpen(nodeLog.node_id, 'output') && 'rotate-180']" />
                  输出数据
                </button>
                <pre
                  v-show="isSectionOpen(nodeLog.node_id, 'output')"
                  class="bg-slate-50 border border-slate-200 rounded p-3 text-xs overflow-x-auto max-h-60"
                >{{ JSON.stringify(nodeLog.output, null, 2) }}</pre>
              </div>

              <!-- 错误信息 -->
              <div v-if="nodeLog.error" class="bg-red-50 border border-red-200 rounded-lg p-3">
                <div class="flex items-start gap-2">
                  <AlertCircle class="w-4 h-4 text-red-600 flex-shrink-0 mt-0.5" />
                  <div>
                    <div class="text-xs font-semibold text-red-900 mb-1">错误信息</div>
                    <div class="text-xs text-red-700">{{ nodeLog.error }}</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  Clock,
  Webhook,
  MousePointerClick,
  Play,
  ChevronDown,
  AlertCircle
} from 'lucide-vue-next'
import BaseButton from '@/components/BaseButton'
import type { WorkflowExecution } from '@/types/workflow'
import { message } from '@/utils/message'

const router = useRouter()
const route = useRoute()

const workflowId = route.params.id as string
const executionId = route.params.executionId as string

// 记录展开/收起状态
const openSections = ref<Record<string, Set<string>>>({})

// 执行详情数据
const execution = ref<WorkflowExecution | null>(null)
const loading = ref(true)

// 加载执行详情
const loadExecution = async () => {
  loading.value = true
  try {
    const { workflowApi } = await import('@/api/workflow')
    const data = await workflowApi.getExecutionDetail(workflowId, executionId)
    execution.value = data
  } catch (error: any) {
    console.error('Load execution failed:', error)
    message.error('加载执行详情失败')
  } finally {
    loading.value = false
  }
}

const handleBack = () => {
  router.push(`/workflows/${workflowId}/executions`)
}

// 切换展开/收起
const toggleSection = (nodeId: string, section: 'input' | 'output') => {
  if (!openSections.value[nodeId]) {
    openSections.value[nodeId] = new Set()
  }
  if (openSections.value[nodeId].has(section)) {
    openSections.value[nodeId].delete(section)
  } else {
    openSections.value[nodeId].add(section)
  }
}

const isSectionOpen = (nodeId: string, section: 'input' | 'output') => {
  return openSections.value[nodeId]?.has(section) || false
}

// 状态样式
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
    success: '执行成功',
    failed: '执行失败',
    cancelled: '已取消'
  }
  return texts[status as keyof typeof texts] || '未知'
}

// 节点状态样式
const getNodeStatusClass = (status: string) => {
  const classes = {
    pending: 'bg-slate-100 text-slate-600',
    running: 'bg-blue-100 text-blue-700',
    success: 'bg-green-100 text-green-700',
    failed: 'bg-red-100 text-red-700',
    skipped: 'bg-amber-100 text-amber-700'
  }
  return classes[status as keyof typeof classes] || classes.pending
}

const getNodeStatusBgClass = (status: string) => {
  const classes = {
    pending: 'bg-slate-200 text-slate-600',
    running: 'bg-blue-500 text-white',
    success: 'bg-green-500 text-white',
    failed: 'bg-red-500 text-white',
    skipped: 'bg-amber-500 text-white'
  }
  return classes[status as keyof typeof classes] || classes.pending
}

const getNodeStatusText = (status: string) => {
  const texts = {
    pending: '等待中',
    running: '运行中',
    success: '成功',
    failed: '失败',
    skipped: '已跳过'
  }
  return texts[status as keyof typeof texts] || '未知'
}

// 触发器
const getTriggerIcon = (type: string) => {
  const icons = {
    scheduled: Clock,
    schedule: Clock,
    webhook: Webhook,
    manual: MousePointerClick
  }
  return icons[type as keyof typeof icons] || Play
}

const getTriggerText = (type: string) => {
  const texts = {
    scheduled: '定时触发',
    schedule: '定时触发',
    webhook: 'Webhook',
    manual: '手动触发'
  }
  return texts[type as keyof typeof texts] || '未知'
}

// 时间格式化
const formatTime = (timestamp?: number) => {
  if (!timestamp) return '-'
  const date = new Date(timestamp * 1000) // 后端返回秒级时间戳，转为毫秒
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

  if (minutes > 0) {
    return `${minutes}分${seconds % 60}秒`
  } else {
    return `${seconds}秒`
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
  loadExecution()
})
</script>
