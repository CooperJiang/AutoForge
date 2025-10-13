<template>
  <div
    class="delay-node bg-bg-elevated rounded-lg shadow-lg border-2 border-purple-400 min-w-[180px] hover:shadow-xl transition-shadow"
    :class="{ 'ring-2 ring-purple-500': data.selected }"
  >
    <!-- 顶部输入点 -->
    <Handle
      type="target"
      :position="Position.Top"
      class="w-3 h-3 !bg-purple-500 !border-2 !border-bg-elevated"
    />

    <!-- 节点内容 -->
    <div class="px-4 py-3">
      <!-- 图标和标题 -->
      <div class="flex items-center gap-2 mb-2">
        <div class="flex-shrink-0 w-8 h-8 rounded-lg bg-gradient-to-br from-accent to-accent-hover flex items-center justify-center text-white shadow-sm">
          <Timer class="w-4 h-4" />
        </div>
        <div class="flex-1 min-w-0">
          <div class="text-sm font-semibold text-text-primary truncate">
            {{ data.name || '延迟' }}
          </div>
          <div class="text-xs text-text-tertiary">
            {{ delayInfo }}
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

    <!-- 底部输出点 -->
    <Handle
      type="source"
      :position="Position.Bottom"
      class="w-3 h-3 !bg-purple-500 !border-2 !border-bg-elevated"
    />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Handle, Position } from '@vue-flow/core'
import { Timer, CheckCircle2, AlertCircle } from 'lucide-vue-next'
import type { WorkflowNode } from '@/types/workflow'

interface Props {
  data: WorkflowNode
}

const props = defineProps<Props>()

const hasConfig = computed(() => {
  return !!(props.data.config?.duration && props.data.config?.unit)
})

const delayInfo = computed(() => {
  const config = props.data.config
  if (!config?.duration) return '等待指定时间'

  const unitMap: Record<string, string> = {
    seconds: '秒',
    minutes: '分钟',
    hours: '小时'
  }

  return `等待 ${config.duration} ${unitMap[config.unit] || '秒'}`
})
</script>

<style scoped>
.delay-node {
  position: relative;
}

:deep(.vue-flow__handle) {
  width: 12px;
  height: 12px;
}
</style>
