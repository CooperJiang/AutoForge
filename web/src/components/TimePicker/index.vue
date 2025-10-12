<template>
  <div class="w-full relative" ref="pickerRef">
    <label v-if="label" class="block text-sm font-medium text-slate-700 mb-2">
      {{ label }}
      <span v-if="required" class="text-rose-500 ml-1">*</span>
    </label>
    <button
      type="button"
      @click="togglePicker"
      class="w-full px-3 py-1.5 text-sm text-left bg-white border-2 border-slate-200 rounded-md transition-all duration-200 hover:border-slate-300 focus:border-blue-400 focus:ring-2 focus:ring-blue-50 focus:outline-none"
      :class="{ 'border-blue-400 ring-2 ring-blue-50': isOpen }"
    >
      <div class="flex items-center justify-between">
        <span class="text-slate-900 font-mono">{{ displayTime }}</span>
        <svg class="w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      </div>
    </button>

    <!-- Time Picker Dropdown -->
    <div
      v-show="isOpen"
      class="absolute z-50 mt-1 bg-white border-2 border-slate-200 rounded-md shadow-xl p-3"
    >
      <div class="flex gap-2">
        <!-- Hours -->
        <div class="flex-1">
          <div class="text-xs text-slate-500 text-center mb-1">时</div>
          <div class="h-32 w-14 overflow-y-auto border-2 border-slate-200 rounded">
            <div
              v-for="h in hours"
              :key="h"
              @click="selectHour(h)"
              class="px-2 py-1 text-xs cursor-pointer text-center transition-colors duration-150 hover:bg-blue-50"
              :class="{ 'bg-blue-500 text-white font-medium': h === selectedHour }"
            >
              {{ h.toString().padStart(2, '0') }}
            </div>
          </div>
        </div>

        <!-- Minutes -->
        <div class="flex-1">
          <div class="text-xs text-slate-500 text-center mb-1">分</div>
          <div class="h-32 w-14 overflow-y-auto border-2 border-slate-200 rounded">
            <div
              v-for="m in minutes"
              :key="m"
              @click="selectMinute(m)"
              class="px-2 py-1 text-xs cursor-pointer text-center transition-colors duration-150 hover:bg-blue-50"
              :class="{ 'bg-blue-500 text-white font-medium': m === selectedMinute }"
            >
              {{ m.toString().padStart(2, '0') }}
            </div>
          </div>
        </div>

        <!-- Seconds -->
        <div class="flex-1">
          <div class="text-xs text-slate-500 text-center mb-1">秒</div>
          <div class="h-32 w-14 overflow-y-auto border-2 border-slate-200 rounded">
            <div
              v-for="s in seconds"
              :key="s"
              @click="selectSecond(s)"
              class="px-2 py-1 text-xs cursor-pointer text-center transition-colors duration-150 hover:bg-blue-50"
              :class="{ 'bg-blue-500 text-white font-medium': s === selectedSecond }"
            >
              {{ s.toString().padStart(2, '0') }}
            </div>
          </div>
        </div>
      </div>

      <div class="mt-2 pt-2 border-t-2 border-slate-200 flex gap-2">
        <button
          type="button"
          @click="confirmTime"
          class="flex-1 px-2 py-1 text-xs bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors duration-200 font-medium"
        >
          确定
        </button>
        <button
          type="button"
          @click="isOpen = false"
          class="flex-1 px-2 py-1 text-xs bg-slate-100 text-slate-700 rounded hover:bg-slate-200 transition-colors duration-200 font-medium"
        >
          取消
        </button>
      </div>
    </div>

    <p v-if="hint" class="mt-1.5 text-xs text-slate-500">{{ hint }}</p>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'

interface Props {
  modelValue: string
  label?: string
  required?: boolean
  hint?: string
}

const props = withDefaults(defineProps<Props>(), {
  required: false,
  hint: ''
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const isOpen = ref(false)
const pickerRef = ref<HTMLElement>()

// Generate hours (0-23), minutes (0-59), seconds (0-59)
const hours = Array.from({ length: 24 }, (_, i) => i)
const minutes = Array.from({ length: 60 }, (_, i) => i)
const seconds = Array.from({ length: 60 }, (_, i) => i)

const selectedHour = ref(9)
const selectedMinute = ref(0)
const selectedSecond = ref(0)

// Parse initial value
const parseTime = () => {
  if (props.modelValue) {
    const parts = props.modelValue.split(':')
    selectedHour.value = parseInt(parts[0]) || 0
    selectedMinute.value = parseInt(parts[1]) || 0
    selectedSecond.value = parseInt(parts[2]) || 0
  }
}

const displayTime = computed(() => {
  const h = selectedHour.value.toString().padStart(2, '0')
  const m = selectedMinute.value.toString().padStart(2, '0')
  const s = selectedSecond.value.toString().padStart(2, '0')
  return `${h}:${m}:${s}`
})

const togglePicker = () => {
  isOpen.value = !isOpen.value
}

const selectHour = (h: number) => {
  selectedHour.value = h
}

const selectMinute = (m: number) => {
  selectedMinute.value = m
}

const selectSecond = (s: number) => {
  selectedSecond.value = s
}

const confirmTime = () => {
  emit('update:modelValue', displayTime.value)
  isOpen.value = false
}

const handleClickOutside = (event: MouseEvent) => {
  if (pickerRef.value && !pickerRef.value.contains(event.target as Node)) {
    isOpen.value = false
  }
}

watch(() => props.modelValue, () => {
  parseTime()
})

onMounted(() => {
  parseTime()
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>
