<template>
  <label :for="computedId" class="inline-flex items-center gap-2 cursor-pointer select-none">
    <input
      :id="computedId"
      type="checkbox"
      class="rounded border-slate-300 text-green-600 focus:ring-green-500"
      :checked="modelValue"
      :disabled="disabled"
      @change="onChange"
    />
    <span v-if="label" class="text-sm text-text-secondary">{{ label }}</span>
    <slot />
  </label>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  modelValue: boolean
  label?: string
  id?: string
  disabled?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: false,
  label: undefined,
  id: undefined,
  disabled: false,
})

const emit = defineEmits<{ 'update:modelValue': [value: boolean] }>()

const computedId = computed(() => props.id || `chk-${Math.random().toString(36).slice(2)}`)

const onChange = (e: Event) => {
  const target = e.target as HTMLInputElement
  emit('update:modelValue', !!target.checked)
}
</script>

