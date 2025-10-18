<template>
  <div class="flex" :class="message.role === 'user' ? 'justify-end' : 'justify-start'">
    <div
      class="max-w-[85%]"
      :class="
        message.role === 'user'
          ? 'bg-primary text-white rounded-2xl rounded-tr-sm px-4 py-3'
          : 'bg-bg-secondary border border-border-primary rounded-2xl rounded-tl-sm p-4'
      "
    >
      <!-- 用户消息 -->
      <div v-if="message.role === 'user'" class="text-sm whitespace-pre-wrap break-words">
        {{ message.content }}
      </div>

      <!-- Agent 消息 -->
      <div v-else class="space-y-4">
        <!-- Plan 显示 -->
        <PlanView v-if="currentPlan || message.plan" :plan="currentPlan || message.plan!" />

        <!-- 步骤列表 -->
        <StepsList
          v-if="(currentSteps && currentSteps.length > 0) || (message.trace && message.trace.steps.length > 0)"
          :steps="currentSteps.length > 0 ? currentSteps : message.trace?.steps || []"
          :is-streaming="isStreaming"
        />

        <!-- 最终答案 -->
        <div
          v-if="message.content"
          class="prose prose-sm max-w-none text-text-primary"
          v-html="renderMarkdown(message.content)"
        />

        <!-- 加载中 -->
        <div v-if="isStreaming && !message.content" class="flex items-center gap-2 text-text-tertiary">
          <div class="flex gap-1">
            <span class="w-2 h-2 bg-current rounded-full animate-bounce" style="animation-delay: 0s"></span>
            <span class="w-2 h-2 bg-current rounded-full animate-bounce" style="animation-delay: 0.15s"></span>
            <span class="w-2 h-2 bg-current rounded-full animate-bounce" style="animation-delay: 0.3s"></span>
          </div>
          <span class="text-sm">思考中...</span>
        </div>

        <!-- Token 使用情况 -->
        <div v-if="message.token_usage" class="text-xs text-text-tertiary border-t border-border-primary pt-2 mt-2">
          Token: {{ message.token_usage.total_tokens }}
          ({{ message.token_usage.prompt_tokens }} + {{ message.token_usage.completion_tokens }})
        </div>

        <!-- 错误信息 -->
        <div v-if="message.error" class="text-sm text-error bg-error/10 border-l-4 border-error rounded-lg p-4 mt-3">
          <div class="flex items-start gap-2">
            <span class="text-xl">⚠️</span>
            <div class="flex-1">
              <div class="font-semibold mb-1">执行失败</div>
              <div class="text-error/90">{{ message.error }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { marked } from 'marked'
import PlanView from './PlanView.vue'
import StepsList from './StepsList.vue'
import type { AgentMessage, AgentPlan, AgentStep } from '@/api/agent'

interface Props {
  message: AgentMessage
  isStreaming?: boolean
  currentPlan?: AgentPlan | null
  currentSteps?: AgentStep[]
}

defineProps<Props>()

// 渲染 Markdown
const renderMarkdown = (text: string) => {
  return marked(text, {
    breaks: true,
    gfm: true,
  })
}
</script>



