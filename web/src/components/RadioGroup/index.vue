<template>
  <div :class="containerClasses">
    <label
      v-for="option in options"
      :key="option.value"
      :class="getOptionClasses(option.value)"
      @click="handleSelect(option.value)"
    >
      <input
        type="radio"
        :value="option.value"
        :checked="modelValue === option.value"
        :disabled="disabled || option.disabled"
        class="sr-only"
      />
      <span class="radio-dot"></span>
      <span class="radio-label">{{ option.label }}</span>
    </label>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface RadioOption {
  label: string
  value: string | number | boolean
  disabled?: boolean
}

interface Props {
  modelValue: string | number | boolean
  options: RadioOption[]
  disabled?: boolean
  direction?: 'horizontal' | 'vertical'
  size?: 'sm' | 'md' | 'lg'
}

const props = withDefaults(defineProps<Props>(), {
  disabled: false,
  direction: 'horizontal',
  size: 'md',
})

const emit = defineEmits<{
  'update:modelValue': [value: string | number | boolean]
  change: [value: string | number | boolean]
}>()

const containerClasses = computed(() => {
  const base = 'flex gap-3'
  const direction = props.direction === 'vertical' ? 'flex-col' : 'flex-row flex-wrap'
  return `${base} ${direction}`
})

const getOptionClasses = (value: string | number | boolean) => {
  const base =
    'inline-flex items-center gap-2 cursor-pointer transition-all duration-200 select-none'
  const isSelected = props.modelValue === value
  const disabled = props.disabled || props.options.find((o) => o.value === value)?.disabled

  const sizes = {
    sm: 'text-xs',
    md: 'text-sm',
    lg: 'text-base',
  }

  const state = disabled
    ? 'opacity-50 cursor-not-allowed'
    : isSelected
      ? 'text-text-primary font-medium'
      : 'text-text-secondary hover:text-text-primary'

  return `${base} ${sizes[props.size]} ${state}`
}

const handleSelect = (value: string | number | boolean) => {
  if (props.disabled) return
  const option = props.options.find((o) => o.value === value)
  if (option?.disabled) return

  emit('update:modelValue', value)
  emit('change', value)
}
</script>

<style scoped>
.radio-dot {
  position: relative;
  width: 16px;
  height: 16px;
  border: 2px solid var(--color-border-primary);
  border-radius: 50%;
  flex-shrink: 0;
  transition: all 0.2s;
  background-color: var(--color-bg-primary);
}

.radio-dot::after {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%) scale(0);
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: var(--color-primary);
  transition: transform 0.2s;
}

input:checked + .radio-dot {
  border-color: var(--color-primary);
}

input:checked + .radio-dot::after {
  transform: translate(-50%, -50%) scale(1);
}

input:disabled + .radio-dot {
  opacity: 0.5;
  cursor: not-allowed;
}

label:hover:not(:has(input:disabled)) .radio-dot {
  border-color: var(--color-primary);
}
</style>
