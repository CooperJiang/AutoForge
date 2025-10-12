<template>
  <div
    class="condition-node bg-white rounded-lg shadow-lg border-2 border-amber-400 min-w-[200px] hover:shadow-xl transition-shadow"
    :class="{ 'ring-2 ring-amber-500': data.selected }"
  >
    <!-- 顶部输入点 -->
    <Handle
      type="target"
      :position="Position.Top"
      class="w-3 h-3 !bg-amber-500 !border-2 !border-white"
    />

    <!-- 节点内容 -->
    <div class="px-4 py-3">
      <!-- 图标和标题 -->
      <div class="flex items-center gap-2 mb-2">
        <div class="flex-shrink-0 w-8 h-8 rounded-lg bg-gradient-to-br from-amber-400 to-orange-500 flex items-center justify-center text-white shadow-sm">
          <GitBranch class="w-4 h-4" />
        </div>
        <div class="flex-1 min-w-0">
          <div class="text-sm font-semibold text-slate-900 truncate">
            {{ data.name || '条件判断' }}
          </div>
          <div class="text-xs text-slate-500">
            {{ conditionTypeLabel }}
          </div>
        </div>
      </div>

      <!-- 配置状态 -->
      <div class="flex items-center gap-1 text-xs">
        <div
          v-if="hasConfig"
          class="flex items-center gap-1 text-emerald-600"
        >
          <CheckCircle2 class="w-3 h-3" />
          <span>已配置</span>
        </div>
        <div
          v-else
          class="flex items-center gap-1 text-amber-600"
        >
          <AlertCircle class="w-3 h-3" />
          <span>待配置</span>
        </div>
      </div>
    </div>

    <!-- 底部输出点 - True分支 -->
    <Handle
      id="true"
      type="source"
      :position="Position.Bottom"
      :style="{ left: '35%' }"
      class="w-3 h-3 !bg-emerald-500 !border-2 !border-white"
    >
      <div class="absolute -bottom-5 left-1/2 -translate-x-1/2 text-xs font-medium text-emerald-600 whitespace-nowrap">
        True
      </div>
    </Handle>

    <!-- 底部输出点 - False分支 -->
    <Handle
      id="false"
      type="source"
      :position="Position.Bottom"
      :style="{ left: '65%' }"
      class="w-3 h-3 !bg-rose-500 !border-2 !border-white"
    >
      <div class="absolute -bottom-5 left-1/2 -translate-x-1/2 text-xs font-medium text-rose-600 whitespace-nowrap">
        False
      </div>
    </Handle>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Handle, Position } from '@vue-flow/core'
import { GitBranch, CheckCircle2, AlertCircle } from 'lucide-vue-next'
import type { WorkflowNode } from '@/types/workflow'

interface Props {
  data: WorkflowNode
}

const props = defineProps<Props>()

const hasConfig = computed(() => {
  return !!(props.data.config?.field && props.data.config?.operator)
})

const conditionTypeLabel = computed(() => {
  const config = props.data.config
  if (!config?.conditionType) return '简单条件'

  const typeMap: Record<string, string> = {
    simple: '简单条件',
    expression: '表达式',
    script: '脚本'
  }
  return typeMap[config.conditionType] || '简单条件'
})
</script>

<style scoped>
.condition-node {
  position: relative;
}

:deep(.vue-flow__handle) {
  width: 12px;
  height: 12px;
}
</style>
