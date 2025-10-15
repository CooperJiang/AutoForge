<template>
  <div class="space-y-6">
    <!-- 加载状态 -->
    <div v-if="loading" class="flex justify-center items-center py-20">
      <div class="text-text-tertiary">加载中...</div>
    </div>

    <!-- 数据为空 -->
    <div v-else-if="!execution" class="flex flex-col justify-center items-center py-20">
      <div class="text-text-tertiary mb-4">执行记录不存在</div>
      <BaseButton size="sm" @click="handleBack">返回列表</BaseButton>
    </div>

    <!-- 执行详情 -->
    <template v-else>
      <!-- 运行中提示 -->
      <div
        v-if="execution?.status === 'running'"
        class="relative bg-primary-light border-l-4 border-primary p-4 rounded-lg mb-4 overflow-hidden"
      >
        <div class="flex items-center gap-3 relative z-10">
          <div class="flex items-center gap-1.5">
            <div
              class="w-1.5 h-1.5 bg-primary rounded-full animate-bounce"
              style="animation-delay: 0ms"
            ></div>
            <div
              class="w-1.5 h-1.5 bg-primary rounded-full animate-bounce"
              style="animation-delay: 150ms"
            ></div>
            <div
              class="w-1.5 h-1.5 bg-primary rounded-full animate-bounce"
              style="animation-delay: 300ms"
            ></div>
          </div>
          <div>
            <div class="text-sm font-semibold text-primary">工作流正在执行中</div>
            <div class="text-xs text-text-secondary mt-0.5">正在实时更新执行进度...</div>
          </div>
        </div>
        <!-- 进度条动画 -->
        <div class="absolute bottom-0 left-0 right-0 h-0.5 bg-primary/20">
          <div class="h-full bg-primary animate-progress"></div>
        </div>
      </div>

      <!-- 顶部信息 -->
      <div class="flex items-center justify-between">
        <div>
          <p class="text-sm text-text-tertiary font-mono">{{ executionId }}</p>
        </div>

        <div class="flex items-center gap-2">
          <span
            :class="[
              'px-3 py-1.5 rounded-full text-sm font-semibold',
              getStatusClass(execution?.status || ''),
            ]"
          >
            <span
              :class="[
                'inline-block w-2 h-2 rounded-full mr-2',
                getStatusDotClass(execution?.status || ''),
              ]"
            ></span>
            {{ getStatusText(execution?.status || '') }}
          </span>
        </div>
      </div>

      <!-- 执行概览 -->
      <div class="bg-bg-elevated rounded-lg border border-border-primary p-6">
        <h3 class="text-sm font-semibold text-text-primary mb-4">执行概览</h3>
        <div class="grid grid-cols-4 gap-6">
          <div>
            <div class="text-xs text-text-tertiary mb-1">触发方式</div>
            <div class="flex items-center gap-2 text-sm font-medium text-text-primary">
              <component :is="getTriggerIcon(execution?.trigger_type || '')" class="w-4 h-4" />
              {{ getTriggerText(execution?.trigger_type || '') }}
            </div>
          </div>
          <div>
            <div class="text-xs text-text-tertiary mb-1">开始时间</div>
            <div class="text-sm font-medium text-text-primary">
              {{ formatTime(execution?.start_time) }}
            </div>
          </div>
          <div>
            <div class="text-xs text-text-tertiary mb-1">结束时间</div>
            <div class="text-sm font-medium text-text-primary">
              {{ execution?.end_time ? formatTime(execution.end_time) : '-' }}
            </div>
          </div>
          <div>
            <div class="text-xs text-text-tertiary mb-1">执行耗时</div>
            <div class="text-sm font-medium text-text-primary">
              {{ execution?.duration_ms != null ? formatDurationMs(execution.duration_ms) : '-' }}
            </div>
          </div>
        </div>
      </div>

      <!-- 节点执行详情 -->
      <div class="bg-bg-elevated rounded-lg border border-border-primary">
        <div class="px-6 py-4 border-b border-border-primary">
          <div class="flex items-center justify-between">
            <h3 class="text-sm font-semibold text-text-primary">节点执行详情</h3>
            <div class="text-xs text-text-secondary">
              进度：{{ getCompletedNodesCount() }} / {{ execution?.total_nodes || 0 }} 个节点
            </div>
          </div>
        </div>

        <div class="divide-y divide-slate-200">
          <div
            v-for="(nodeLog, index) in execution?.node_logs"
            :key="nodeLog.node_id"
            :ref="(el) => setNodeRef(nodeLog.node_id, el)"
            class="p-6"
          >
            <div class="flex items-start gap-4">
              <!-- 序号和连接线 -->
              <div class="flex flex-col items-center">
                <div
                  :class="[
                    'w-8 h-8 rounded-full flex items-center justify-center text-sm font-semibold',
                    getNodeStatusBgClass(nodeLog.status),
                  ]"
                >
                  {{ index + 1 }}
                </div>
                <div
                  v-if="index < (execution?.node_logs.length || 0) - 1"
                  class="flex-1 w-0.5 bg-bg-tertiary my-2 min-h-[20px]"
                ></div>
              </div>

              <!-- 节点信息 -->
              <div class="flex-1">
                <div class="flex items-center justify-between mb-2">
                  <div class="flex items-center gap-2">
                    <h4 class="text-sm font-semibold text-text-primary">{{ nodeLog.node_name }}</h4>
                    <span
                      v-if="nodeLog.tool_code"
                      class="px-2 py-0.5 bg-primary-light text-primary text-xs font-mono rounded"
                      >{{ nodeLog.tool_code }}</span
                    >
                    <span class="text-xs text-text-placeholder font-mono">{{
                      nodeLog.node_id
                    }}</span>
                    <span
                      v-if="nodeLog.retry_count > 0"
                      class="px-2 py-0.5 bg-orange-100 text-orange-700 text-xs rounded"
                    >
                      重试 {{ nodeLog.retry_count }} 次
                    </span>
                  </div>
                  <div class="flex items-center gap-2">
                    <!-- Running 状态显示加载动画 -->
                    <div v-if="nodeLog.status === 'running'" class="flex items-center gap-1.5">
                      <div
                        class="w-1.5 h-1.5 bg-blue-500 rounded-full animate-bounce"
                        style="animation-delay: 0ms"
                      ></div>
                      <div
                        class="w-1.5 h-1.5 bg-blue-500 rounded-full animate-bounce"
                        style="animation-delay: 150ms"
                      ></div>
                      <div
                        class="w-1.5 h-1.5 bg-blue-500 rounded-full animate-bounce"
                        style="animation-delay: 300ms"
                      ></div>
                    </div>
                    <span
                      :class="[
                        'px-2 py-0.5 rounded-full text-xs font-medium',
                        getNodeStatusClass(nodeLog.status),
                      ]"
                    >
                      {{ getNodeStatusText(nodeLog.status) }}
                    </span>
                  </div>
                </div>

                <!-- 时间信息 -->
                <div
                  v-if="nodeLog.start_time"
                  class="flex items-center gap-4 text-xs text-text-secondary mb-3"
                >
                  <span
                    >开始：{{ formatTime(nodeLog.start_time) }} ({{
                      getRelativeTime(nodeLog.start_time)
                    }})</span
                  >
                  <span v-if="nodeLog.end_time"> 结束：{{ formatTime(nodeLog.end_time) }} </span>
                  <span v-if="nodeLog.duration_ms != null">
                    耗时：{{ formatDurationMs(nodeLog.duration_ms) }}
                  </span>
                  <span v-else-if="nodeLog.status === 'running'">
                    执行中：{{ getCurrentDuration(nodeLog.start_time) }}
                  </span>
                </div>

                <!-- 输入数据 -->
                <div v-if="nodeLog.input && Object.keys(nodeLog.input).length > 0" class="mb-3">
                  <button
                    type="button"
                    @click="toggleSection(nodeLog.node_id, 'input')"
                    class="flex items-center gap-1 text-xs font-medium text-text-secondary hover:text-text-primary mb-2"
                  >
                    <ChevronDown
                      :class="[
                        'w-3 h-3 transition-transform',
                        (isSectionOpen(nodeLog.node_id, 'input') || nodeLog.status === 'running') &&
                          'rotate-180',
                      ]"
                    />
                    输入数据
                    <span v-if="nodeLog.status === 'running'" class="text-blue-500">(执行中)</span>
                  </button>
                  <div
                    v-show="isSectionOpen(nodeLog.node_id, 'input') || nodeLog.status === 'running'"
                  >
                    <JsonViewer :content="formatInputData(nodeLog.input)" />
                  </div>
                </div>

                <!-- 输出数据 -->
                <div v-if="nodeLog.output || nodeLog.status === 'running'" class="mb-3">
                  <button
                    type="button"
                    @click="toggleSection(nodeLog.node_id, 'output')"
                    class="flex items-center gap-1 text-xs font-medium text-text-secondary hover:text-text-primary mb-2"
                  >
                    <ChevronDown
                      :class="[
                        'w-3 h-3 transition-transform',
                        (isSectionOpen(nodeLog.node_id, 'output') ||
                          nodeLog.status === 'running') &&
                          'rotate-180',
                      ]"
                    />
                    输出数据
                    <span v-if="nodeLog.status === 'running'" class="text-blue-500">(执行中)</span>
                  </button>
                  <div
                    v-show="
                      isSectionOpen(nodeLog.node_id, 'output') || nodeLog.status === 'running'
                    "
                  >
                    <div
                      v-if="nodeLog.status === 'running' && !nodeLog.output"
                      class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4"
                    >
                      <div class="flex items-center gap-2 text-blue-600 dark:text-blue-400">
                        <div class="flex items-center gap-1">
                          <div class="w-1.5 h-1.5 bg-blue-500 rounded-full animate-pulse"></div>
                          <div
                            class="w-1.5 h-1.5 bg-blue-500 rounded-full animate-pulse"
                            style="animation-delay: 200ms"
                          ></div>
                          <div
                            class="w-1.5 h-1.5 bg-blue-500 rounded-full animate-pulse"
                            style="animation-delay: 400ms"
                          ></div>
                        </div>
                        <span class="text-sm">正在执行中，等待执行结束...</span>
                      </div>
                    </div>
                    <OutputViewer
                      v-else-if="nodeLog.output"
                      :output="nodeLog.output"
                      :output-render="nodeLog.output_render"
                    />
                  </div>
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
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Clock, Webhook, MousePointerClick, Play, ChevronDown, AlertCircle } from 'lucide-vue-next'
import BaseButton from '@/components/BaseButton'
import JsonViewer from '@/components/JsonViewer'
import OutputViewer from '@/components/OutputViewer'
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

// 节点 DOM 引用
const nodeRefs = ref<Record<string, HTMLElement>>({})

// 轮询相关
let pollingTimer: NodeJS.Timeout | null = null
const POLLING_INTERVAL = 2000 // 2秒轮询一次

// 上一次的节点数量，用于检测是否有新节点
let previousNodeCount = 0

// 加载执行详情
const loadExecution = async (showLoading = true) => {
  if (showLoading) {
    loading.value = true
  }
  try {
    const { workflowApi } = await import('@/api/workflow')
    const data = await workflowApi.getExecutionDetail(workflowId, executionId)
    execution.value = data

    // 如果状态是运行中，启动轮询
    if (data.status === 'running') {
      startPolling()
    } else {
      stopPolling()
    }
  } catch (error: any) {
    console.error('Load execution failed:', error)
    if (showLoading) {
      message.error('加载执行详情失败')
    }
  } finally {
    if (showLoading) {
      loading.value = false
    }
  }
}

// 启动轮询
const startPolling = () => {
  if (pollingTimer) return

  pollingTimer = setInterval(() => {
    loadExecution(false) // 轮询时不显示 loading
  }, POLLING_INTERVAL)
}

// 停止轮询
const stopPolling = () => {
  if (pollingTimer) {
    clearInterval(pollingTimer)
    pollingTimer = null
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
    success: '执行成功',
    failed: '执行失败',
    cancelled: '已取消',
  }
  return texts[status as keyof typeof texts] || '未知'
}

// 节点状态样式
const getNodeStatusClass = (status: string) => {
  const classes = {
    pending: 'bg-bg-tertiary text-text-secondary',
    running: 'bg-primary-light text-primary',
    success: 'bg-green-100 text-green-700',
    failed: 'bg-red-100 text-red-700',
    skipped: 'bg-amber-100 text-amber-700',
  }
  return classes[status as keyof typeof classes] || classes.pending
}

const getNodeStatusBgClass = (status: string) => {
  const classes = {
    pending: 'bg-bg-tertiary text-text-secondary',
    running: 'bg-[var(--color-primary)] text-white',
    success: 'bg-green-500 text-white',
    failed: 'bg-red-500 text-white',
    skipped: 'bg-amber-500 text-white',
  }
  return classes[status as keyof typeof classes] || classes.pending
}

const getNodeStatusText = (status: string) => {
  const texts = {
    pending: '等待中',
    running: '运行中',
    success: '成功',
    failed: '失败',
    skipped: '已跳过',
  }
  return texts[status as keyof typeof texts] || '未知'
}

// 格式化输入数据
const formatInputData = (input: any): string => {
  if (!input) return ''

  try {
    const formatted: any = {}

    // 外部参数
    if (input.external_params && Object.keys(input.external_params).length > 0) {
      formatted['外部参数 (External Params)'] = input.external_params
    }

    // 原始配置
    if (input.config) {
      formatted['原始配置 (Original Config)'] = input.config
    }

    // 解析后的配置（变量替换后）
    if (input.resolved_config) {
      formatted['解析配置 (Resolved Config)'] = input.resolved_config
    }

    return JSON.stringify(formatted, null, 2)
  } catch (error) {
    console.error('Failed to format input data:', error)
    return '[数据格式化失败]'
  }
}

// 触发器
const getTriggerIcon = (type: string) => {
  const icons = {
    scheduled: Clock,
    schedule: Clock,
    webhook: Webhook,
    manual: MousePointerClick,
  }
  return icons[type as keyof typeof icons] || Play
}

const getTriggerText = (type: string) => {
  const texts = {
    scheduled: '定时触发',
    schedule: '定时触发',
    webhook: 'Webhook',
    manual: '手动触发',
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
    second: '2-digit',
  })
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

// 获取相对时间（例如："3秒前"）
const getRelativeTime = (timestamp: number) => {
  const now = Math.floor(Date.now() / 1000)
  const diff = now - timestamp

  if (diff < 60) {
    return `${diff}秒前`
  } else if (diff < 3600) {
    return `${Math.floor(diff / 60)}分钟前`
  } else if (diff < 86400) {
    return `${Math.floor(diff / 3600)}小时前`
  } else {
    return `${Math.floor(diff / 86400)}天前`
  }
}

// 获取当前执行时长（用于 running 状态）
const getCurrentDuration = (startTime: number) => {
  const now = Math.floor(Date.now() / 1000)
  const durationSeconds = now - startTime

  if (durationSeconds < 60) {
    return `${durationSeconds}秒`
  } else if (durationSeconds < 3600) {
    const minutes = Math.floor(durationSeconds / 60)
    const seconds = durationSeconds % 60
    return `${minutes}分${seconds}秒`
  } else {
    const hours = Math.floor(durationSeconds / 3600)
    const minutes = Math.floor((durationSeconds % 3600) / 60)
    return `${hours}小时${minutes}分钟`
  }
}

// 获取已完成节点数量
const getCompletedNodesCount = () => {
  if (!execution.value) return 0
  return execution.value.node_logs.filter(
    (log) => log.status === 'success' || log.status === 'failed' || log.status === 'skipped'
  ).length
}

// 设置节点引用
const setNodeRef = (nodeId: string, el: any) => {
  if (el) {
    nodeRefs.value[nodeId] = el
  }
}

// 滚动到最新的 running 节点
const scrollToRunningNode = async () => {
  await nextTick()

  if (!execution.value) return

  // 找到 running 状态的节点
  const runningNode = execution.value.node_logs.find((log) => log.status === 'running')

  if (runningNode && nodeRefs.value[runningNode.node_id]) {
    const element = nodeRefs.value[runningNode.node_id]
    element.scrollIntoView({
      behavior: 'smooth',
      block: 'center',
    })
  }
}

// 监听执行数据变化，自动滚动到新节点
watch(
  () => execution.value?.node_logs.length,
  (newCount) => {
    if (newCount && newCount > previousNodeCount) {
      // 有新节点添加，滚动到最新的 running 节点
      scrollToRunningNode()
      previousNodeCount = newCount
    }
  }
)

onMounted(() => {
  loadExecution()
})

onUnmounted(() => {
  stopPolling()
})
</script>

<style scoped>
@keyframes progress {
  0% {
    width: 0%;
    opacity: 1;
  }
  50% {
    opacity: 1;
  }
  100% {
    width: 100%;
    opacity: 0;
  }
}

.animate-progress {
  animation: progress 2s ease-in-out infinite;
}
</style>
