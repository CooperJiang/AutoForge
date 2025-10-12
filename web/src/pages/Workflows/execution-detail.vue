<template>
  <div class="space-y-6">
    <!-- 顶部信息 -->
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-3">
        <BaseButton size="sm" variant="ghost" @click="handleBack">
          <ArrowLeft class="w-4 h-4 mr-1" />
          返回列表
        </BaseButton>
        <div class="h-6 w-px bg-slate-200"></div>
        <div>
          <h2 class="text-lg font-semibold text-slate-900">执行详情</h2>
          <p class="text-sm text-slate-500 font-mono">{{ executionId }}</p>
        </div>
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
            <component :is="getTriggerIcon(execution?.trigger.type || '')" class="w-4 h-4" />
            {{ getTriggerText(execution?.trigger.type || '') }}
          </div>
        </div>
        <div>
          <div class="text-xs text-slate-500 mb-1">开始时间</div>
          <div class="text-sm font-medium text-slate-900">{{ formatTime(execution?.startTime || '') }}</div>
        </div>
        <div>
          <div class="text-xs text-slate-500 mb-1">结束时间</div>
          <div class="text-sm font-medium text-slate-900">
            {{ execution?.endTime ? formatTime(execution.endTime) : '-' }}
          </div>
        </div>
        <div>
          <div class="text-xs text-slate-500 mb-1">执行耗时</div>
          <div class="text-sm font-medium text-slate-900">
            {{ execution?.endTime ? formatDuration(execution.startTime, execution.endTime) : '-' }}
          </div>
        </div>
      </div>

      <!-- 触发器数据 -->
      <div v-if="execution?.trigger.data" class="mt-4 pt-4 border-t border-slate-200">
        <div class="text-xs font-semibold text-slate-700 mb-2">触发器数据</div>
        <pre class="bg-slate-900 text-slate-100 rounded-lg p-3 text-xs overflow-x-auto">{{ JSON.stringify(execution.trigger.data, null, 2) }}</pre>
      </div>
    </div>

    <!-- 节点执行详情 -->
    <div class="bg-white rounded-lg border border-slate-200">
      <div class="px-6 py-4 border-b border-slate-200">
        <h3 class="text-sm font-semibold text-slate-900">节点执行详情</h3>
      </div>

      <div class="divide-y divide-slate-200">
        <div
          v-for="(nodeExec, index) in execution?.nodeExecutions"
          :key="nodeExec.nodeId"
          class="p-6"
        >
          <div class="flex items-start gap-4">
            <!-- 序号和连接线 -->
            <div class="flex flex-col items-center">
              <div
                :class="[
                  'w-8 h-8 rounded-full flex items-center justify-center text-sm font-semibold',
                  getNodeStatusBgClass(nodeExec.status)
                ]"
              >
                {{ index + 1 }}
              </div>
              <div
                v-if="index < (execution?.nodeExecutions.length || 0) - 1"
                class="flex-1 w-0.5 bg-slate-200 my-2 min-h-[20px]"
              ></div>
            </div>

            <!-- 节点信息 -->
            <div class="flex-1">
              <div class="flex items-center justify-between mb-2">
                <div class="flex items-center gap-2">
                  <h4 class="text-sm font-semibold text-slate-900">{{ nodeExec.nodeName }}</h4>
                  <span class="text-xs text-slate-500 font-mono">{{ nodeExec.nodeId }}</span>
                </div>
                <span
                  :class="[
                    'px-2 py-0.5 rounded-full text-xs font-medium',
                    getNodeStatusClass(nodeExec.status)
                  ]"
                >
                  {{ getNodeStatusText(nodeExec.status) }}
                </span>
              </div>

              <!-- 时间信息 -->
              <div v-if="nodeExec.startTime" class="flex items-center gap-4 text-xs text-slate-600 mb-3">
                <span>开始：{{ formatTime(nodeExec.startTime) }}</span>
                <span v-if="nodeExec.endTime">
                  耗时：{{ formatDuration(nodeExec.startTime, nodeExec.endTime) }}
                </span>
              </div>

              <!-- 输入数据 -->
              <div v-if="nodeExec.input" class="mb-3">
                <button
                  type="button"
                  @click="toggleSection(nodeExec.nodeId, 'input')"
                  class="flex items-center gap-1 text-xs font-medium text-slate-700 hover:text-slate-900 mb-1"
                >
                  <ChevronDown :class="['w-3 h-3 transition-transform', isSectionOpen(nodeExec.nodeId, 'input') && 'rotate-180']" />
                  输入数据
                </button>
                <pre
                  v-show="isSectionOpen(nodeExec.nodeId, 'input')"
                  class="bg-slate-50 border border-slate-200 rounded p-3 text-xs overflow-x-auto max-h-60"
                >{{ JSON.stringify(nodeExec.input, null, 2) }}</pre>
              </div>

              <!-- 输出数据 -->
              <div v-if="nodeExec.output" class="mb-3">
                <button
                  type="button"
                  @click="toggleSection(nodeExec.nodeId, 'output')"
                  class="flex items-center gap-1 text-xs font-medium text-slate-700 hover:text-slate-900 mb-1"
                >
                  <ChevronDown :class="['w-3 h-3 transition-transform', isSectionOpen(nodeExec.nodeId, 'output') && 'rotate-180']" />
                  输出数据
                </button>
                <pre
                  v-show="isSectionOpen(nodeExec.nodeId, 'output')"
                  class="bg-slate-50 border border-slate-200 rounded p-3 text-xs overflow-x-auto max-h-60"
                >{{ JSON.stringify(nodeExec.output, null, 2) }}</pre>
              </div>

              <!-- 错误信息 -->
              <div v-if="nodeExec.error" class="bg-red-50 border border-red-200 rounded-lg p-3">
                <div class="flex items-start gap-2">
                  <AlertCircle class="w-4 h-4 text-red-600 flex-shrink-0 mt-0.5" />
                  <div>
                    <div class="text-xs font-semibold text-red-900 mb-1">错误信息</div>
                    <div class="text-xs text-red-700">{{ nodeExec.error }}</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  ArrowLeft,
  Clock,
  Webhook,
  MousePointerClick,
  Play,
  ChevronDown,
  AlertCircle
} from 'lucide-vue-next'
import BaseButton from '@/components/BaseButton'
import type { WorkflowExecution } from '@/types/workflow'

const router = useRouter()
const route = useRoute()

const workflowId = route.params.id as string
const executionId = route.params.executionId as string

// 记录展开/收起状态
const openSections = ref<Record<string, Set<string>>>({})

// Mock 数据
const execution = ref<WorkflowExecution>({
  id: executionId,
  workflowId: workflowId,
  status: 'success',
  startTime: new Date(Date.now() - 3600000).toISOString(),
  endTime: new Date(Date.now() - 3500000).toISOString(),
  trigger: {
    type: 'webhook',
    data: {
      user: 'admin',
      action: 'test',
      timestamp: Date.now()
    }
  },
  nodeExecutions: [
    {
      nodeId: 'node_1',
      nodeName: 'HTTP请求 - 获取数据',
      status: 'success',
      startTime: new Date(Date.now() - 3600000).toISOString(),
      endTime: new Date(Date.now() - 3550000).toISOString(),
      input: {
        url: 'https://api.example.com/data',
        method: 'GET',
        headers: {
          'Authorization': 'Bearer token123'
        }
      },
      output: {
        status: 200,
        data: {
          result: 'success',
          items: [1, 2, 3]
        },
        headers: {
          'content-type': 'application/json'
        }
      }
    },
    {
      nodeId: 'node_2',
      nodeName: '条件判断',
      status: 'success',
      startTime: new Date(Date.now() - 3550000).toISOString(),
      endTime: new Date(Date.now() - 3540000).toISOString(),
      input: {
        field: 'node_1.response.data.result',
        value: 'success'
      },
      output: {
        result: true,
        branch: 'true'
      }
    },
    {
      nodeId: 'node_3',
      nodeName: '发送邮件通知',
      status: 'success',
      startTime: new Date(Date.now() - 3540000).toISOString(),
      endTime: new Date(Date.now() - 3500000).toISOString(),
      input: {
        to: 'admin@example.com',
        subject: '工作流执行成功',
        body: '数据已成功获取'
      },
      output: {
        success: true,
        messageId: 'msg_1234567890'
      }
    }
  ]
})

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

// 时间格式化
const formatTime = (time: string) => {
  if (!time) return '-'
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

  if (minutes > 0) {
    return `${minutes}分${seconds % 60}秒`
  } else {
    return `${seconds}秒`
  }
}

onMounted(() => {
  // TODO: 加载执行详情
})
</script>
