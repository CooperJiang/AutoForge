<template>
  <div class="space-y-3">
    <div class="flex items-center gap-2 mb-2">
      <Zap class="w-4 h-4 text-primary" />
      <h4 class="text-sm font-semibold text-text-primary">执行过程</h4>
    </div>

    <div
      v-for="(step, index) in steps"
      :key="index"
      class="border border-border-secondary rounded-lg p-3 bg-bg-tertiary/20 hover:bg-bg-tertiary/30 transition-colors"
    >
      <div class="flex items-start gap-3">
        <!-- 步骤编号 -->
        <div
          class="flex-shrink-0 w-6 h-6 rounded-full bg-primary/10 text-primary text-xs font-semibold flex items-center justify-center"
        >
          {{ step.step }}
        </div>

        <!-- 步骤内容 -->
        <div class="flex-1 min-w-0 space-y-2">
          <!-- 工具信息 -->
          <div v-if="step.action?.tool" class="flex items-center gap-2">
            <span class="text-sm font-medium text-text-primary">
              🔧 {{ step.action.tool }}
            </span>
            <span class="text-xs text-text-tertiary">{{ formatDuration(step.elapsed_ms) }}</span>
          </div>

          <!-- 工具参数 -->
          <details v-if="step.action?.args && Object.keys(step.action.args).length > 0" class="group">
            <summary class="text-xs text-text-tertiary cursor-pointer hover:text-text-secondary flex items-center gap-1">
              <ChevronRight class="w-3 h-3 transition-transform group-open:rotate-90" />
              查看参数
            </summary>
            <div class="mt-2 text-xs bg-bg-primary border border-border-primary rounded p-2 overflow-x-auto">
              <pre class="text-text-secondary">{{ JSON.stringify(step.action.args, null, 2) }}</pre>
            </div>
          </details>

          <!-- 观察结果 -->
          <div v-if="step.observation" class="text-sm text-text-secondary bg-bg-primary/50 rounded p-2">
            <div class="font-medium text-text-primary mb-1 text-xs">结果：</div>
            {{ step.observation }}
          </div>

          <!-- 工具输出 -->
          <details v-if="step.tool_output && Object.keys(step.tool_output).length > 0" class="group">
            <summary class="text-xs text-text-tertiary cursor-pointer hover:text-text-secondary flex items-center gap-1">
              <ChevronRight class="w-3 h-3 transition-transform group-open:rotate-90" />
              查看详细输出
            </summary>
            <div class="mt-2 text-xs bg-bg-primary border border-border-primary rounded p-2 overflow-x-auto">
              <pre class="text-text-secondary">{{ JSON.stringify(step.tool_output, null, 2) }}</pre>
            </div>
          </details>

          <!-- 错误信息 -->
          <div v-if="step.error" class="text-sm text-error bg-error/10 border border-error/20 rounded p-2">
            <div class="font-medium mb-1 text-xs">❌ 错误：</div>
            {{ step.error }}
          </div>
        </div>
      </div>
    </div>

    <!-- 加载中的步骤 -->
    <div v-if="isStreaming" class="border border-border-secondary rounded-lg p-3 bg-bg-tertiary/20">
      <div class="flex items-center gap-3">
        <div class="w-6 h-6 rounded-full bg-primary/20 animate-pulse" />
        <div class="flex items-center gap-2 text-text-tertiary">
          <div class="flex gap-1">
            <span class="w-1.5 h-1.5 bg-current rounded-full animate-bounce" style="animation-delay: 0s"></span>
            <span class="w-1.5 h-1.5 bg-current rounded-full animate-bounce" style="animation-delay: 0.15s"></span>
            <span class="w-1.5 h-1.5 bg-current rounded-full animate-bounce" style="animation-delay: 0.3s"></span>
          </div>
          <span class="text-sm">执行中...</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Zap, ChevronRight } from 'lucide-vue-next'
import type { AgentStep } from '@/api/agent'

interface Props {
  steps: AgentStep[]
  isStreaming?: boolean
}

defineProps<Props>()

// 格式化时长
const formatDuration = (ms: number) => {
  if (ms < 1000) return `${ms}ms`
  return `${(ms / 1000).toFixed(2)}s`
}

// 截断文本
const truncate = (text: string, maxLength: number) => {
  if (text.length <= maxLength) return text
  return text.substring(0, maxLength) + '...'
}
</script>



