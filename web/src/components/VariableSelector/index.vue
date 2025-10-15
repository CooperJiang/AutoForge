<template>
  <div class="relative">
    <BaseInput
      :model-value="modelValue"
      @update:model-value="handleInput"
      :placeholder="placeholder"
    />
  </div>
</template>

<script setup lang="ts">
import BaseInput from '@/components/BaseInput'

interface Props {
  modelValue: string
  placeholder?: string
  previousNodes?: Array<{ id: string; name: string; type: string }>
  envVars?: Array<{ key: string; value: string; description?: string }>
}

defineProps<Props>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const handleInput = (value: string) => {
  emit('update:modelValue', value)
}
</script>
