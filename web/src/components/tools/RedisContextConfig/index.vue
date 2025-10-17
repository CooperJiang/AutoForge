<template>
  <div class="space-y-4">
    <div class="bg-primary-light border-l-4 border-primary rounded p-3 text-xs text-text-secondary">
      使用缓存(优先 Redis, 无则内存)存取上下文键值；支持 get / set / delete 与可选 TTL
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
      <div>
        <label class="block text-sm font-medium text-text-secondary mb-1">操作类型</label>
        <BaseSelect v-model="localConfig.action" :options="actionOptions" />
      </div>
      <div>
        <label class="block text-sm font-medium text-text-secondary mb-1">TTL（秒，可选）</label>
        <BaseInput v-model.number="localConfig.ttl_seconds" type="number" placeholder="0" />
        <p class="mt-1 text-xs text-text-tertiary">
          仅在 set 时生效；0 或留空表示永不过期。例如: 3600 (1小时)
        </p>
      </div>
    </div>

    <div>
      <label class="block text-sm font-medium text-text-secondary mb-1">键名</label>
      <BaseInput
        v-model="localConfig.key"
        placeholder="例如：session:user_123 或 {{params.session_id}}"
      />
    </div>

    <div v-if="localConfig.action === 'set'">
      <label class="block text-sm font-medium text-text-secondary mb-1">值（JSON 或文本）</label>
      <textarea
        v-model="localConfig.value"
        rows="5"
        class="w-full px-3 py-2 border-2 border-border-primary rounded-lg focus:outline-none focus:border-primary bg-bg-elevated text-text-primary font-mono text-sm"
        placeholder='{"token":"xxx","exp":1700000000} 或 任意文本'
      />
      <p class="mt-1 text-xs text-text-tertiary">可直接填写 JSON 字符串或普通文本</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import BaseInput from '@/components/BaseInput'
import BaseSelect from '@/components/BaseSelect'
import type { WorkflowNode, WorkflowEnvVar } from '@/types/workflow'

interface Props {
  config: Record<string, any>
  previousNodes?: WorkflowNode[]
  envVars?: WorkflowEnvVar[]
}

const props = withDefaults(defineProps<Props>(), {
  previousNodes: () => [],
  envVars: () => [],
})

const emit = defineEmits<{ (e: 'update:config', config: Record<string, any>): void }>()

const actionOptions = [
  { label: '读取 (get)', value: 'get' },
  { label: '写入 (set)', value: 'set' },
  { label: '删除 (delete)', value: 'delete' },
]

const localConfig = ref({
  action: props.config.action || 'get',
  key: props.config.key || '',
  value: props.config.value || '',
  ttl_seconds: props.config.ttl_seconds ?? 0,
})

watch(
  () => props.config,
  (val) => {
    localConfig.value = {
      action: val.action || 'get',
      key: val.key || '',
      value: val.value || '',
      ttl_seconds: val.ttl_seconds ?? 0,
    }
  },
  { deep: true }
)

watch(
  localConfig,
  (val) => {
    emit('update:config', { ...val })
  },
  { deep: true }
)
</script>
