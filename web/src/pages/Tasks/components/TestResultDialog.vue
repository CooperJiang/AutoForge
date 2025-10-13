<template>
  <Dialog
    :model-value="modelValue"
    title="测试结果"
    confirm-text="关闭"
    cancel-text=""
    @update:model-value="$emit('update:modelValue', $event)"
    @confirm="$emit('update:modelValue', false)"
  >
    <div v-if="result" class="space-y-3">
      <!-- 状态 -->
      <div class="flex items-center gap-2">
        <span class="text-sm font-medium text-text-secondary">状态:</span>
        <span
          :class="[
            'px-2 py-1 rounded text-xs font-medium',
            result.success
              ? 'bg-green-100 text-green-800'
              : 'bg-red-100 text-red-800'
          ]"
        >
          {{ result.success ? '成功' : '失败' }}
        </span>
      </div>

      <!-- HTTP 状态码 -->
      <div v-if="result.status_code" class="flex items-center gap-2">
        <span class="text-sm font-medium text-text-secondary">HTTP 状态码:</span>
        <span class="text-sm text-text-primary">{{ result.status_code }}</span>
      </div>

      <!-- 响应时间 -->
      <div v-if="result.duration_ms !== undefined" class="flex items-center gap-2">
        <span class="text-sm font-medium text-text-secondary">响应时间:</span>
        <span class="text-sm text-text-primary">{{ result.duration_ms }} ms</span>
      </div>

      <!-- 错误信息 -->
      <div v-if="result.error_message" class="space-y-1">
        <span class="text-sm font-medium text-text-secondary">错误信息:</span>
        <div class="bg-red-50 border border-red-200 rounded-lg p-3">
          <pre class="text-xs text-red-800 whitespace-pre-wrap break-words">{{ result.error_message }}</pre>
        </div>
      </div>

      <!-- 响应内容 -->
      <div v-if="result.response_body" class="space-y-1">
        <span class="text-sm font-medium text-text-secondary">响应内容:</span>
        <JsonViewer :content="result.response_body" />
      </div>
    </div>
  </Dialog>
</template>

<script setup lang="ts">
import Dialog from '@/components/Dialog'
import JsonViewer from '@/components/JsonViewer'

interface TestResult {
  success: boolean
  status_code?: number
  response_body?: string
  duration_ms?: number
  error_message?: string
}

defineProps<{
  modelValue: boolean
  result: TestResult | null
}>()

defineEmits<{
  'update:modelValue': [value: boolean]
}>()
</script>
