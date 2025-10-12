<template>
  <div
    class="switch-node bg-white rounded-lg shadow-lg border-2 border-indigo-400 min-w-[200px] hover:shadow-xl transition-shadow"
    :class="{ 'ring-2 ring-indigo-500': data.selected }"
  >
    <!-- 顶部输入点 -->
    <Handle
      type="target"
      :position="Position.Top"
      class="w-3 h-3 !bg-indigo-500 !border-2 !border-white"
    />

    <!-- 节点内容 -->
    <div class="px-4 py-3">
      <!-- 图标和标题 -->
      <div class="flex items-center gap-2 mb-2">
        <div class="flex-shrink-0 w-8 h-8 rounded-lg bg-gradient-to-br from-indigo-400 to-indigo-600 flex items-center justify-center text-white shadow-sm">
          <Split class="w-4 h-4" />
        </div>
        <div class="flex-1 min-w-0">
          <div class="text-sm font-semibold text-slate-900 truncate">
            {{ data.name || '开关分支' }}
          </div>
          <div class="text-xs text-slate-500">
            多路分支
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
          <span>{{ branchCount }} 个分支</span>
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

    <!-- 底部输出点 - 多个分支 -->
    <Handle
      v-for="(branch, index) in branches"
      :key="branch.value"
      :id="branch.value"
      type="source"
      :position="Position.Bottom"
      :style="{ left: `${(index + 1) * (100 / (branches.length + 1))}%` }"
      class="w-3 h-3 !bg-indigo-500 !border-2 !border-white"
    >
      <div
        class="absolute -bottom-5 left-1/2 -translate-x-1/2 text-xs font-medium whitespace-nowrap"
        :style="{ color: getBranchColor(index) }"
      >
        {{ branch.label }}
      </div>
    </Handle>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Handle, Position } from '@vue-flow/core'
import { Split, CheckCircle2, AlertCircle } from 'lucide-vue-next'
import type { WorkflowNode } from '@/types/workflow'

interface Props {
  data: WorkflowNode
}

const props = defineProps<Props>()

const hasConfig = computed(() => {
  return !!(props.data.config?.field && props.data.config?.cases?.length > 0)
})

const branchCount = computed(() => {
  const cases = props.data.config?.cases || []
  return cases.length + 1 // +1 for default branch
})

const branches = computed(() => {
  const cases = props.data.config?.cases || []
  const result = cases.map((c: any, index: number) => ({
    value: `case_${index}`,
    label: c.label || `Case ${index + 1}`
  }))
  result.push({
    value: 'default',
    label: 'Default'
  })
  return result
})

const getBranchColor = (index: number) => {
  const colors = ['#3b82f6', '#8b5cf6', '#ec4899', '#f59e0b', '#10b981']
  return colors[index % colors.length]
}
</script>

<style scoped>
.switch-node {
  position: relative;
}

:deep(.vue-flow__handle) {
  width: 12px;
  height: 12px;
}
</style>
