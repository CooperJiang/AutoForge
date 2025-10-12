<template>
  <div class="tool-node bg-white rounded-lg shadow-lg border-2 border-slate-200 hover:border-blue-400 transition-all">
    <!-- 节点头部 -->
    <div :class="['px-3 py-2 rounded-t-lg flex items-center gap-2', getToolBgClass(data.toolCode)]">
      <component :is="getToolIcon(data.toolCode)" class="w-4 h-4 text-white flex-shrink-0" />
      <span class="text-sm font-medium text-white truncate">{{ data.name }}</span>
    </div>

    <!-- 节点内容 -->
    <div class="px-3 py-2 text-xs text-slate-600">
      <div v-if="hasConfig" class="flex items-center gap-1">
        <CheckCircle2 class="w-3 h-3 text-green-600" />
        <span>已配置</span>
      </div>
      <div v-else class="flex items-center gap-1">
        <AlertCircle class="w-3 h-3 text-orange-500" />
        <span>未配置</span>
      </div>
    </div>

    <!-- 连接点 -->
    <Handle type="target" :position="Position.Top" class="handle-top" />
    <Handle type="source" :position="Position.Bottom" class="handle-bottom" />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Handle, Position } from '@vue-flow/core'
import { Globe, Mail, HeartPulse, CheckCircle2, AlertCircle } from 'lucide-vue-next'
import type { WorkflowNode } from '@/types/workflow'

interface Props {
  data: WorkflowNode
}

const props = defineProps<Props>()

// 是否已配置
const hasConfig = computed(() => {
  return Object.keys(props.data.config).length > 0
})

// 获取工具图标
const getToolIcon = (code?: string) => {
  const iconMap: Record<string, any> = {
    'http_request': Globe,
    'email_sender': Mail,
    'health_checker': HeartPulse
  }
  return iconMap[code || ''] || Globe
}

// 获取工具背景色
const getToolBgClass = (code?: string) => {
  const colorMap: Record<string, string> = {
    'http_request': 'bg-gradient-to-r from-blue-500 to-purple-600',
    'email_sender': 'bg-gradient-to-r from-purple-500 to-pink-600',
    'health_checker': 'bg-gradient-to-r from-indigo-500 to-blue-600'
  }
  return colorMap[code || ''] || 'bg-gradient-to-r from-blue-500 to-purple-600'
}
</script>

<style scoped>
.tool-node {
  width: 200px;
}

:deep(.handle-top),
:deep(.handle-bottom) {
  width: 10px;
  height: 10px;
  background: #3b82f6;
  border: 2px solid white;
}

:deep(.handle-top) {
  top: -5px;
}

:deep(.handle-bottom) {
  bottom: -5px;
}
</style>
