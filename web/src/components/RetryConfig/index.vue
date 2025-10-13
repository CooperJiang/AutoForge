<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <label class="text-sm font-medium text-text-secondary">
        错误重试
      </label>
      <div class="flex items-center gap-2">
        <input
          v-model="localConfig.enabled"
          type="checkbox"
          :id="`retry-enabled-${uniqueId}`"
          class="rounded border-slate-300 text-green-600 focus:ring-green-500"
          @change="emitUpdate"
        />
        <label :for="`retry-enabled-${uniqueId}`" class="text-sm text-text-secondary">
          启用自动重试
        </label>
      </div>
    </div>

    <div v-if="localConfig.enabled" class="space-y-3 pl-4 border-l-2 border-green-200">
      <div class="bg-primary-light border-l-4 border-border-focus p-3 text-xs">
        <p class="text-primary">
          <svg class="inline-block w-3.5 h-3.5 mr-1" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd"/>
          </svg>
          当节点执行失败时，系统将自动重试指定次数
        </p>
      </div>

      <div>
        <label class="block text-sm font-medium text-text-secondary mb-1">
          最大重试次数 <span class="text-red-500">*</span>
        </label>
        <BaseInput
          v-model.number="localConfig.maxRetries"
          type="number"
          min="1"
          max="10"
          placeholder="3"
          @update:model-value="emitUpdate"
        />
        <p class="text-xs text-text-tertiary mt-1">建议设置 1-5 次，最多 10 次</p>
      </div>

      <div>
        <label class="block text-sm font-medium text-text-secondary mb-1">
          重试间隔（秒） <span class="text-red-500">*</span>
        </label>
        <BaseInput
          v-model.number="localConfig.retryInterval"
          type="number"
          min="1"
          max="300"
          placeholder="5"
          @update:model-value="emitUpdate"
        />
        <p class="text-xs text-text-tertiary mt-1">每次重试之间的等待时间</p>
      </div>

      <div class="flex items-center gap-2">
        <input
          v-model="localConfig.exponentialBackoff"
          type="checkbox"
          :id="`exponential-${uniqueId}`"
          class="rounded border-slate-300 text-green-600 focus:ring-green-500"
          @change="emitUpdate"
        />
        <label :for="`exponential-${uniqueId}`" class="text-sm text-text-secondary">
          使用指数退避策略
        </label>
      </div>

      <div v-if="localConfig.exponentialBackoff" class="bg-bg-hover rounded-lg p-3 text-xs">
        <div class="font-semibold text-text-secondary mb-1">重试时间预览：</div>
        <div class="text-text-secondary space-y-0.5">
          <div v-for="i in Math.min(localConfig.maxRetries, 5)" :key="i">
            第 {{ i }} 次重试：等待 {{ calculateBackoff(i) }} 秒
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import BaseInput from '@/components/BaseInput'
import type { NodeRetryConfig } from '@/types/workflow'

interface Props {
  config: NodeRetryConfig
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:config': [config: NodeRetryConfig]
}>()

const uniqueId = Math.random().toString(36).substring(7)

const localConfig = ref<NodeRetryConfig>({
  enabled: false,
  maxRetries: 3,
  retryInterval: 5,
  exponentialBackoff: false,
  ...props.config
})

watch(() => props.config, (newVal) => {
  localConfig.value = { ...localConfig.value, ...newVal }
}, { deep: true })

const calculateBackoff = (attempt: number): number => {
  if (!localConfig.value.exponentialBackoff) {
    return localConfig.value.retryInterval
  }
  return localConfig.value.retryInterval * Math.pow(2, attempt - 1)
}

const emitUpdate = () => {
  emit('update:config', localConfig.value)
}
</script>
