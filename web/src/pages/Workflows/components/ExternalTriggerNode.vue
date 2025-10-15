<template>
  <div
    class="external-trigger-node bg-bg-elevated rounded-lg shadow-lg border-2 border-blue-500 hover:border-blue-600 transition-all group relative"
  >
    
    <button
      class="absolute -top-2 -right-2 w-6 h-6 bg-red-500 hover:bg-red-600 text-white rounded-full opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center z-10 shadow-lg"
      @click.stop="handleDelete"
      title="删除节点"
    >
      <X class="w-4 h-4" />
    </button>

    
    <div
      class="px-3 py-2 bg-gradient-to-r from-blue-500 to-blue-600 rounded-t-lg flex items-center gap-2"
    >
      <Globe class="w-4 h-4 text-white flex-shrink-0" />
      <span class="text-sm font-medium text-white truncate">{{ data.name }}</span>
    </div>

    
    <div class="px-3 py-2.5 space-y-2">
      <div class="flex items-center gap-1 text-xs text-text-secondary">
        <Zap class="w-3 h-3 text-blue-500" />
        <span>外部 API 触发</span>
      </div>

      
      <div
        v-if="paramCount > 0"
        class="flex items-center gap-1 text-xs text-blue-600 bg-blue-50 px-2 py-1 rounded"
      >
        <Settings class="w-3 h-3" />
        <span>{{ paramCount }} 个参数</span>
      </div>

      <div v-else class="text-xs text-text-tertiary">未配置参数</div>
    </div>

    
    <Handle type="source" :position="Position.Bottom" class="handle-bottom" />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Handle, Position } from '@vue-flow/core'
import { Globe, Zap, X, Settings } from 'lucide-vue-next'
import type { WorkflowNode } from '@/types/workflow'

interface Props {
  data: WorkflowNode
}

const props = defineProps<Props>()

const emit = defineEmits<{
  delete: [id: string]
}>()

const paramCount = computed(() => {
  return props.data.config?.params?.length || 0
})

const handleDelete = (e: Event) => {
  e.stopPropagation()
  emit('delete', props.data.id)
}
</script>

<style scoped>
.external-trigger-node {
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
