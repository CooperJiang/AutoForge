<template>
  <div
    class="trigger-node bg-bg-elevated rounded-lg shadow-lg border-2 border-primary hover:border-primary transition-all group relative"
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
    <div
      class="px-3 py-2 bg-gradient-to-r from-primary to-accent rounded-t-lg flex items-center gap-2"
    >
      <Clock class="w-4 h-4 text-white flex-shrink-0" />
      <span class="text-sm font-medium text-white truncate">{{ data.name }}</span>
    </div>

    <!-- 节点内容 -->
    <div class="px-3 py-2 text-xs text-text-secondary">
      <div class="flex items-center gap-1">
        <Zap class="w-3 h-3 text-primary" />
        <span>工作流入口</span>
      </div>
    </div>

    <!-- 连接点 -->
    <Handle type="source" :position="Position.Bottom" class="handle-bottom" />
  </div>
</template>

<script setup lang="ts">
import { Handle, Position } from '@vue-flow/core'
import { Clock, Zap, X } from 'lucide-vue-next'
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
</script>

<style scoped>
.trigger-node {
  width: 200px;
}

:deep(.handle-bottom) {
  width: 10px;
  height: 10px;
  background: #3b82f6;
  border: 2px solid white;
  bottom: -5px;
}
</style>
