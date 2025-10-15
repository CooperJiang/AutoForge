<template>
  <div class="slider-container">
    <div class="flex items-center gap-3">
      <input
        :value="modelValue"
        @input="handleInput"
        type="range"
        :min="min"
        :max="max"
        :step="step"
        :disabled="disabled"
        class="slider"
        :class="{ 'slider-disabled': disabled }"
      />
      <span v-if="showValue" class="slider-value">
        {{ displayValue }}
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  modelValue: number
  min?: number
  max?: number
  step?: number
  disabled?: boolean
  showValue?: boolean
  valueFormatter?: (value: number) => string
}

const props = withDefaults(defineProps<Props>(), {
  min: 0,
  max: 100,
  step: 1,
  disabled: false,
  showValue: true,
})

const emit = defineEmits<{
  'update:modelValue': [value: number]
}>()

const handleInput = (event: Event) => {
  const target = event.target as HTMLInputElement
  emit('update:modelValue', Number(target.value))
}

const displayValue = computed(() => {
  if (props.valueFormatter) {
    return props.valueFormatter(props.modelValue)
  }
  return String(props.modelValue)
})
</script>

<style scoped>
.slider-container {
  @apply w-full;
}

.slider {
  @apply flex-1 h-2 rounded-lg appearance-none cursor-pointer;
  @apply bg-bg-tertiary;
  @apply transition-all;
}

.slider::-webkit-slider-thumb {
  @apply appearance-none w-5 h-5 rounded-full cursor-pointer;
  @apply bg-primary;
  @apply transition-all;
  @apply shadow-md;
}

.slider::-webkit-slider-thumb:hover {
  @apply scale-110;
  @apply shadow-lg;
}

.slider::-moz-range-thumb {
  @apply w-5 h-5 rounded-full cursor-pointer;
  @apply bg-primary;
  @apply border-0;
  @apply transition-all;
  @apply shadow-md;
}

.slider::-moz-range-thumb:hover {
  @apply scale-110;
  @apply shadow-lg;
}

.slider::-webkit-slider-track {
  @apply h-2 rounded-lg;
  @apply bg-bg-tertiary;
}

.slider::-moz-range-track {
  @apply h-2 rounded-lg;
  @apply bg-bg-tertiary;
}

.slider:focus {
  @apply outline-none ring-2 ring-primary ring-opacity-50;
}

.slider-disabled {
  @apply opacity-50 cursor-not-allowed;
}

.slider-disabled::-webkit-slider-thumb {
  @apply cursor-not-allowed;
}

.slider-disabled::-webkit-slider-thumb:hover {
  @apply scale-100;
}

.slider-disabled::-moz-range-thumb {
  @apply cursor-not-allowed;
}

.slider-disabled::-moz-range-thumb:hover {
  @apply scale-100;
}

.slider-value {
  @apply text-sm text-text-secondary w-20 text-right;
}
</style>
