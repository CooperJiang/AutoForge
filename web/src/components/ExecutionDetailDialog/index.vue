<template>
  <div
    v-if="modelValue"
    class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
  >
    <div
      class="bg-bg-elevated rounded-xl shadow-2xl max-w-4xl w-full max-h-[90vh] overflow-hidden flex flex-col"
    >
      <div
        class="px-6 py-4 border-b-2 border-border-primary flex justify-between items-center flex-shrink-0"
      >
        <h3 class="text-xl font-bold text-text-primary">执行详情</h3>
        <button
          @click="$emit('update:modelValue', false)"
          class="text-text-tertiary hover:text-text-secondary transition-colors"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M6 18L18 6M6 6l12 12"
            />
          </svg>
        </button>
      </div>

      <div class="flex-1 overflow-y-auto p-6">
        <div v-if="execution" class="space-y-6">
          <div>
            <h4 class="text-sm font-semibold text-text-secondary mb-3 flex items-center gap-2">
              <span class="w-1 h-4 bg-[var(--color-primary)] rounded"></span>
              基本信息
            </h4>
            <div class="grid grid-cols-3 gap-4 bg-bg-hover rounded-lg p-4">
              <div>
                <p class="text-xs text-text-tertiary mb-1">执行ID</p>
                <p class="text-sm text-text-primary font-mono">{{ execution.id }}</p>
              </div>
              <div>
                <p class="text-xs text-text-tertiary mb-1">任务ID</p>
                <p class="text-sm text-text-primary font-mono">{{ execution.task_id }}</p>
              </div>
              <div>
                <p class="text-xs text-text-tertiary mb-1">用户ID</p>
                <p class="text-sm text-text-primary font-mono">{{ execution.user_id }}</p>
              </div>
            </div>
          </div>

          <div>
            <h4 class="text-sm font-semibold text-text-secondary mb-3 flex items-center gap-2">
              <span class="w-1 h-4 bg-success-light0 rounded"></span>
              执行状态
            </h4>
            <div class="grid grid-cols-3 gap-4 bg-bg-hover rounded-lg p-4">
              <div>
                <p class="text-xs text-text-tertiary mb-1">状态</p>
                <span
                  :class="[
                    'inline-flex px-2 py-1 text-xs font-medium rounded-full',
                    execution.status === 'success'
                      ? 'bg-green-100 text-success'
                      : 'bg-red-100 text-red-700',
                  ]"
                >
                  {{ execution.status === 'success' ? '成功' : '失败' }}
                </span>
              </div>
              <div>
                <p class="text-xs text-text-tertiary mb-1">HTTP状态码</p>
                <span
                  :class="[
                    'inline-flex px-2 py-1 text-xs font-semibold rounded',
                    execution.response_status >= 200 && execution.response_status < 300
                      ? 'bg-green-100 text-success'
                      : execution.response_status >= 400
                        ? 'bg-red-100 text-red-700'
                        : 'bg-yellow-100 text-yellow-700',
                  ]"
                >
                  {{ execution.response_status || 'N/A' }}
                </span>
              </div>
              <div>
                <p class="text-xs text-text-tertiary mb-1">执行时长</p>
                <p class="text-sm text-text-primary">{{ execution.duration_ms }}ms</p>
              </div>
            </div>
          </div>

          <div>
            <h4 class="text-sm font-semibold text-text-secondary mb-3 flex items-center gap-2">
              <span class="w-1 h-4 bg-purple-500 rounded"></span>
              响应数据
            </h4>
            <div v-if="execution.response_body">
              <JsonViewer :content="execution.response_body" />
            </div>
            <div v-else class="bg-bg-hover rounded-lg p-4">
              <p class="text-sm text-text-tertiary">无响应数据</p>
            </div>
          </div>

          <div v-if="execution.error_message">
            <h4 class="text-sm font-semibold text-text-secondary mb-3 flex items-center gap-2">
              <span class="w-1 h-4 bg-red-500 rounded"></span>
              错误信息
            </h4>
            <div class="bg-red-50 border border-red-200 rounded-lg p-4">
              <p class="text-sm text-red-800">{{ execution.error_message }}</p>
            </div>
          </div>

          <div>
            <h4 class="text-sm font-semibold text-text-secondary mb-3 flex items-center gap-2">
              <span class="w-1 h-4 bg-orange-500 rounded"></span>
              时间信息
            </h4>
            <div class="grid grid-cols-2 gap-4 bg-bg-hover rounded-lg p-4">
              <div>
                <p class="text-xs text-text-tertiary mb-1">开始时间</p>
                <p class="text-sm text-text-primary">{{ formatTime(execution.started_at) }}</p>
              </div>
              <div>
                <p class="text-xs text-text-tertiary mb-1">结束时间</p>
                <p class="text-sm text-text-primary">{{ formatTime(execution.finished_at) }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="px-6 py-4 border-t-2 border-border-primary flex justify-end gap-3 flex-shrink-0">
        <button
          @click="$emit('update:modelValue', false)"
          class="px-4 py-2 bg-bg-tertiary hover:bg-bg-tertiary text-text-secondary font-medium rounded-lg transition-colors"
        >
          关闭
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps, defineEmits } from 'vue'
import JsonViewer from '../JsonViewer/index.vue'

defineProps<{
  modelValue: boolean
  execution: any | null
}>()

defineEmits(['update:modelValue'])

const formatTime = (time: any) => {
  if (!time) return '--'
  const timestamp = typeof time === 'number' ? time * 1000 : time
  return new Date(timestamp).toLocaleString('zh-CN')
}
</script>
