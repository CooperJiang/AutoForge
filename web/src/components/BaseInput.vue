<template>
  <div class="w-full">
    <label v-if="label" class="block text-sm font-medium text-slate-700 mb-2">
      {{ label }}
      <span v-if="required" class="text-rose-500 ml-1">*</span>
    </label>
    <input
      :value="modelValue"
      @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
      :type="type"
      :placeholder="placeholder"
      :required="required"
      :min="min"
      :max="max"
      class="w-full px-3 py-1.5 text-sm text-slate-900 bg-white border-2 border-slate-200 rounded-md transition-all duration-200 focus:border-blue-400 focus:ring-2 focus:ring-blue-50 focus:outline-none hover:border-slate-300"
      :class="inputClass"
    />
    <p v-if="hint" class="mt-1.5 text-xs text-slate-500">{{ hint }}</p>
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
  inputClass: ''
})

defineEmits<{
  'update:modelValue': [value: string]
}>()
</script>
