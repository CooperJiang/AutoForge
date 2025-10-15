<template>
  <div
    class="tool-node bg-bg-elevated rounded-lg shadow-lg border-2 border-border-primary hover:border-primary transition-all group relative"
  >
    <!-- 删除按钮 (hover 显示) -->
    <button
      class="absolute -top-2 -right-2 w-6 h-6 bg-red-500 hover:bg-red-600 text-white rounded-full opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center z-10 shadow-lg"
      @click.stop="handleDelete"
      title="删除节点"
    >
      <X class="w-4 h-4" />
    </button>

    <!-- 节点头部 -->
    <div :class="['px-3 py-2 rounded-t-lg flex items-center gap-2', getToolBgClass(data.toolCode)]">
      <component :is="getToolIcon(data.toolCode)" class="w-4 h-4 text-white flex-shrink-0" />
      <span class="text-sm font-medium text-white truncate">{{ data.name }}</span>
    </div>

    <!-- 节点内容 -->
    <div class="px-3 py-2 text-xs text-text-secondary">
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
import {
  Globe,
  Mail,
  HeartPulse,
  Braces,
  Database,
  CheckCircle2,
  AlertCircle,
  X,
} from 'lucide-vue-next'
import type { WorkflowNode } from '@/types/workflow'

interface Props {
  data: WorkflowNode
}

const props = defineProps<Props>()

const emit = defineEmits<{
  delete: [id: string]
}>()

const handleDelete = (e: Event) => {
  e.stopPropagation()
  emit('delete', props.data.id)
}

// 是否已配置
const hasConfig = computed(() => {
  return Object.keys(props.data.config).length > 0
})

// 获取工具图标
const getToolIcon = (code?: string) => {
  const iconMap: Record<string, any> = {
    http_request: Globe,
    email_sender: Mail,
    health_checker: HeartPulse,
    json_transform: Braces,
    redis_context: Database,
  }
  return iconMap[code || ''] || Globe
}

// 获取工具背景色
const getToolBgClass = (code?: string) => {
  const colorMap: Record<string, string> = {
    http_request: 'bg-gradient-to-r from-primary to-accent',
    email_sender: 'bg-gradient-to-r from-purple-500 to-pink-600',
    health_checker: 'bg-gradient-to-r from-primary to-accent',
    json_transform: 'bg-gradient-to-r from-emerald-500 to-teal-600',
    redis_context: 'bg-gradient-to-r from-slate-500 to-slate-700',
  }
  return colorMap[code || ''] || 'bg-gradient-to-r from-primary to-accent'
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
