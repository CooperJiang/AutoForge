<template>
  <div class="trigger-node bg-white rounded-lg shadow-lg border-2 border-blue-300 hover:border-blue-500 transition-all">
    <!-- 节点头部 -->
    <div class="px-3 py-2 bg-gradient-to-r from-blue-500 to-purple-600 rounded-t-lg flex items-center gap-2">
      <Clock class="w-4 h-4 text-white flex-shrink-0" />
      <span class="text-sm font-medium text-white truncate">{{ data.name }}</span>
    </div>

    <!-- 节点内容 -->
    <div class="px-3 py-2 text-xs text-slate-600">
      <div class="flex items-center gap-1">
        <Zap class="w-3 h-3 text-blue-600" />
        <span>工作流入口</span>
      </div>
    </div>

    <!-- 连接点 -->
    <Handle type="source" :position="Position.Bottom" class="handle-bottom" />
  </div>
</template>

<script setup lang="ts">
import { Handle, Position } from '@vue-flow/core'
import { Clock, Zap } from 'lucide-vue-next'
import type { WorkflowNode } from '@/types/workflow'

interface Props {
  data: WorkflowNode
}

defineProps<Props>()
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
