<template>
  <div class="space-y-4">
    <div class="bg-info-light border-l-4 border-info rounded-lg p-3">
      <p class="text-sm text-info-text">
        <svg class="inline-block w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
          <path
            fill-rule="evenodd"
            d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z"
            clip-rule="evenodd"
          />
        </svg>
        在继续执行下一个节点之前等待指定的时间
      </p>
    </div>

    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2">
        等待时长 <span class="text-red-500">*</span>
      </label>
      <div class="grid grid-cols-2 gap-2">
        <BaseInput
          v-model.number="localConfig.duration"
          type="number"
          :min="1"
          placeholder="1"
          @update:model-value="emitUpdate"
        />
        <BaseSelect
          v-model="localConfig.unit"
          :options="unitOptions"
          @update:model-value="emitUpdate"
        />
      </div>
    </div>

    <!-- 示例说明 -->
    <div class="bg-bg-hover rounded-lg p-3">
      <div class="text-xs font-semibold text-text-secondary mb-2">使用场景：</div>
      <div class="text-xs text-text-secondary space-y-1">
        <div>• 等待API限流冷却时间</div>
        <div>• 在发送通知前等待系统稳定</div>
        <div>• 分批处理任务之间的间隔</div>
        <div>• 等待外部系统处理完成</div>
      </div>
    </div>

    <div class="bg-warning-light border border-warning rounded-lg p-3">
      <div class="text-xs font-semibold text-warning-text mb-1">⚠️ 注意事项</div>
      <p class="text-xs text-warning-text">
        延迟期间工作流将暂停执行，请合理设置等待时长，避免影响整体执行效率
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import BaseInput from '@/components/BaseInput'
import BaseSelect from '@/components/BaseSelect'

interface Props {
  config: Record<string, any>
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:config': [config: Record<string, any>]
}>()

const localConfig = ref({
  duration: 5,
  unit: 'seconds',
  ...props.config,
})

const unitOptions = [
  { label: '秒', value: 'seconds' },
  { label: '分钟', value: 'minutes' },
  { label: '小时', value: 'hours' },
]

watch(
  () => props.config,
  (newVal) => {
    localConfig.value = { ...localConfig.value, ...newVal }
  },
  { deep: true }
)

const emitUpdate = () => {
  emit('update:config', localConfig.value)
}
</script>
