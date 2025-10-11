<template>
  <div class="flex gap-2 items-start">
    <div class="flex-1">
      <input
        :value="param.key"
        @input="updateKey(($event.target as HTMLInputElement).value)"
        type="text"
        :placeholder="keyPlaceholder"
        class="w-full px-2 py-1 text-xs text-slate-900 bg-white border-2 border-slate-200 rounded transition-all duration-200 focus:border-blue-400 focus:ring-2 focus:ring-blue-50 focus:outline-none hover:border-slate-300"
      />
    </div>
    <div class="flex-1">
      <input
        :value="param.value"
        @input="updateValue(($event.target as HTMLInputElement).value)"
        type="text"
        :placeholder="valuePlaceholder"
        class="w-full px-2 py-1 text-xs text-slate-900 bg-white border-2 border-slate-200 rounded transition-all duration-200 focus:border-blue-400 focus:ring-2 focus:ring-blue-50 focus:outline-none hover:border-slate-300"
      />
    </div>
    <button
      type="button"
      @click="$emit('remove')"
      class="p-1 text-rose-500 hover:bg-rose-50 rounded transition-colors duration-200"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
      </svg>
    </button>
  </div>
</template>

<script setup lang="ts">
interface Param {
  key: string
  value: string
}

interface Props {
  param: Param
  keyPlaceholder?: string
  valuePlaceholder?: string
}

const props = withDefaults(defineProps<Props>(), {
  keyPlaceholder: '键',
  valuePlaceholder: '值'
})

const emit = defineEmits<{
  'update:param': [param: Param]
  'remove': []
}>()

const updateKey = (key: string) => {
  emit('update:param', { ...props.param, key })
}

const updateValue = (value: string) => {
  emit('update:param', { ...props.param, value })
}
</script>
