<template>
  <Teleport to="body">
    <Transition name="dialog">
      <div
        v-if="modelValue"
        class="fixed inset-0 z-50 flex items-center justify-center p-4"
        @click.self="onClose"
      >
        <div class="absolute inset-0 bg-black bg-opacity-50 transition-opacity"></div>

        <div
          class="relative bg-bg-elevated rounded-lg shadow-xl border-2 border-border-primary w-full max-w-2xl transform transition-all max-h-[90vh] flex flex-col"
        >
          <div
            class="px-5 py-4 border-b-2 border-border-primary flex items-center justify-between flex-shrink-0"
          >
            <div class="flex items-center gap-3">
              <div
                v-if="result?.success"
                class="w-10 h-10 rounded-full bg-emerald-100 flex items-center justify-center"
              >
                <CheckCircle :size="24" class="text-emerald-600" />
              </div>
              <div
                v-else
                class="w-10 h-10 rounded-full bg-rose-100 flex items-center justify-center"
              >
                <XCircle :size="24" class="text-error" />
              </div>
              <div>
                <h3
                  class="text-base font-semibold"
                  :class="result?.success ? 'text-emerald-700' : 'text-rose-700'"
                >
                  {{ result?.success ? '测试成功' : '测试失败' }}
                </h3>
                <p class="text-xs text-text-tertiary mt-0.5">接口配置测试结果</p>
              </div>
            </div>
            <button
              @click="onClose"
              class="text-text-tertiary hover:text-text-secondary transition-colors"
            >
              <X :size="20" />
            </button>
          </div>

          <div class="px-5 py-4 overflow-y-auto flex-1">
            <div class="space-y-4">
              <div class="grid grid-cols-2 gap-3">
                <div class="bg-bg-hover border-2 border-border-primary rounded-lg p-3">
                  <div class="text-xs text-text-tertiary mb-1">状态码</div>
                  <div
                    class="text-lg font-semibold"
                    :class="
                      result?.status_code && result.status_code >= 200 && result.status_code < 300
                        ? 'text-emerald-600'
                        : 'text-error'
                    "
                  >
                    {{ result?.status_code || '-' }}
                  </div>
                </div>
                <div class="bg-bg-hover border-2 border-border-primary rounded-lg p-3">
                  <div class="text-xs text-text-tertiary mb-1">响应时间</div>
                  <div class="text-lg font-semibold text-primary">{{ result?.duration_ms }}ms</div>
                </div>
              </div>

              <div
                v-if="result?.error_message"
                class="bg-rose-50 border-2 border-rose-200 rounded-lg p-4"
              >
                <div class="flex items-center gap-2 mb-2">
                  <AlertCircle :size="16" class="text-error" />
                  <div class="font-semibold text-sm text-rose-700">错误信息</div>
                </div>
                <div class="text-sm text-error font-mono whitespace-pre-wrap break-words">
                  {{ result.error_message }}
                </div>
              </div>

              <div v-if="result?.response_body">
                <div class="flex items-center justify-between mb-2">
                  <div class="font-semibold text-sm text-text-secondary">响应内容</div>
                  <div class="text-xs text-text-tertiary">
                    {{ formatSize(result.response_body.length) }}
                  </div>
                </div>
                <JsonViewer :content="result.response_body" />
              </div>
            </div>
          </div>

          <div class="px-5 py-4 border-t-2 border-border-primary flex justify-end flex-shrink-0">
            <BaseButton variant="primary" @click="onClose"> 关闭 </BaseButton>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { CheckCircle, XCircle, X, AlertCircle } from 'lucide-vue-next'
import BaseButton from '../BaseButton/index.vue'
import JsonViewer from '../JsonViewer/index.vue'

interface TestResult {
  success: boolean
  status_code: number
  response_body: string
  duration_ms: number
  error_message?: string
}

interface Props {
  modelValue: boolean
  result: TestResult | null
}

defineProps<Props>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
}>()

const onClose = () => {
  emit('update:modelValue', false)
}

const formatSize = (bytes: number) => {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}
</script>

<style scoped>
.dialog-enter-active,
.dialog-leave-active {
  transition: opacity 0.2s ease;
}

.dialog-enter-from,
.dialog-leave-to {
  opacity: 0;
}

.dialog-enter-active .relative,
.dialog-leave-active .relative {
  transition:
    transform 0.2s ease,
    opacity 0.2s ease;
}

.dialog-enter-from .relative,
.dialog-leave-to .relative {
  transform: scale(0.95);
  opacity: 0;
}
</style>
