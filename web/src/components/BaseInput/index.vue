<template>
  <div class="w-full">
    <label v-if="label" class="block text-sm font-medium text-text-primary mb-2">
      {{ label }}
      <span v-if="required" class="text-error ml-1">*</span>
    </label>
    <input
      :value="modelValue"
      @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
      :type="type"
      :placeholder="placeholder"
      :required="required"
      :min="min"
      :max="max"
      class="w-full px-3 py-1.5 text-sm text-text-primary bg-bg-primary border-2 border-border-primary rounded-md transition-all duration-200 focus:border-border-focus focus:ring-2 focus:ring-primary-light focus:outline-none hover:border-border-secondary placeholder:text-text-placeholder"
      :class="inputClass"
    />
    <p v-if="hint" class="mt-1.5 text-xs text-text-tertiary">{{ hint }}</p>
  </div>
</template>

<script setup lang="ts">
interface Props {
  modelValue: string | number
  label?: string
  type?: string
  placeholder?: string
  required?: boolean
  hint?: string
  inputClass?: string
  min?: string | number
  max?: string | number
}

withDefaults(defineProps<Props>(), {
  type: 'text',
  placeholder: '',
  required: false,
  hint: '',
  inputClass: '',
})

defineEmits<{
  'update:modelValue': [value: string]
}>()
</script>

<style scoped>
/* 隐藏 number 类型的上下箭头 */
input[type='number']::-webkit-inner-spin-button,
input[type='number']::-webkit-outer-spin-button {
  -webkit-appearance: none;
  margin: 0;
}

input[type='number'] {
  -moz-appearance: textfield;
}
</style>
