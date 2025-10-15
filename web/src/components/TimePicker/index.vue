<template>
  <div class="space-y-2" @click.stop @mousedown.stop @mouseup.stop>
    <label v-if="label" class="block text-sm font-medium text-text-secondary">
      {{ label }}
      <span v-if="required" class="text-red-500">*</span>
    </label>

    <div class="relative">
      <input
        type="text"
        :value="displayValue"
        @click.stop="showPicker = true"
        @focus.stop="showPicker = true"
        @mousedown.stop
        readonly
        placeholder="选择时间"
        class="w-full px-3 py-2 border-2 border-border-primary rounded-lg focus:outline-none focus:border-primary bg-bg-elevated text-text-primary cursor-pointer"
      />

      
      <div
        v-if="showPicker"
        class="absolute z-[9999] mt-1 bg-bg-elevated border-2 border-border-primary rounded-lg shadow-lg p-4"
        @click.stop
        @mousedown.stop
        @mouseup.stop
      >
        <div class="flex gap-2 items-center">
          
          <div class="flex flex-col">
            <label class="text-xs text-text-secondary mb-1">时</label>
            <input
              v-model.number="hours"
              type="number"
              min="0"
              max="23"
              class="w-16 px-2 py-1 border border-border-primary rounded text-center focus:outline-none focus:border-primary"
              @input="handleHourInput"
              @click="handleInputClick('hour', $event)"
              @mousedown="handleInputMouseDown('hour', $event)"
              @mouseup.stop
              @focus.stop
            />
          </div>

          <span class="text-text-primary mt-5">:</span>

          
          <div class="flex flex-col">
            <label class="text-xs text-text-secondary mb-1">分</label>
            <input
              v-model.number="minutes"
              type="number"
              min="0"
              max="59"
              class="w-16 px-2 py-1 border border-border-primary rounded text-center focus:outline-none focus:border-primary"
              @input="updateTime"
              @click.stop
              @mousedown.stop
              @mouseup.stop
              @focus.stop
            />
          </div>

          <span class="text-text-primary mt-5">:</span>

          
          <div class="flex flex-col">
            <label class="text-xs text-text-secondary mb-1">秒</label>
            <input
              v-model.number="seconds"
              type="number"
              min="0"
              max="59"
              class="w-16 px-2 py-1 border border-border-primary rounded text-center focus:outline-none focus:border-primary"
              @input="updateTime"
              @click.stop
              @mousedown.stop
              @mouseup.stop
              @focus.stop
            />
          </div>
        </div>

        <div class="flex gap-2 mt-3">
          <button
            @click.stop="confirmTime"
            @mousedown.stop
            class="flex-1 px-3 py-1.5 bg-primary text-white rounded hover:opacity-90 text-sm"
          >
            确定
          </button>
          <button
            @click.stop="showPicker = false"
            @mousedown.stop
            class="flex-1 px-3 py-1.5 bg-bg-tertiary text-text-secondary rounded hover:bg-bg-hover text-sm"
          >
            取消
          </button>
        </div>
      </div>
    </div>

    <p v-if="hint" class="text-xs text-text-tertiary">{{ hint }}</p>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'

interface Props {
  modelValue?: string
  label?: string
  hint?: string
  required?: boolean
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const showPicker = ref(false)
const hours = ref(9)
const minutes = ref(0)
const seconds = ref(0)

// 显示值 HH:MM:SS
const displayValue = computed(() => {
  if (!props.modelValue) return ''
  return props.modelValue
})

// 初始化时间
const initTime = () => {
  if (props.modelValue) {
    const parts = props.modelValue.split(':')
    if (parts.length >= 2) {
      hours.value = parseInt(parts[0]) || 0
      minutes.value = parseInt(parts[1]) || 0
      seconds.value = parseInt(parts[2]) || 0
    }
  }
}

// 监听 modelValue 变化
watch(
  () => props.modelValue,
  () => {
    initTime()
  },
  { immediate: true }
)

// 实时更新时间（输入时）
const updateTime = () => {
  // 确保数值在有效范围内
  if (hours.value < 0) hours.value = 0
  if (hours.value > 23) hours.value = 23
  if (minutes.value < 0) minutes.value = 0
  if (minutes.value > 59) minutes.value = 59
  if (seconds.value < 0) seconds.value = 0
  if (seconds.value > 59) seconds.value = 59
}

const handleHourInput = () => {
  updateTime()
}

const handleInputClick = (_type: string, e: Event) => {
  e.stopPropagation()
  e.preventDefault()
}

const handleInputMouseDown = (_type: string, e: Event) => {
  e.stopPropagation()
}

// 确认时间
const confirmTime = () => {
  const h = String(hours.value).padStart(2, '0')
  const m = String(minutes.value).padStart(2, '0')
  const s = String(seconds.value).padStart(2, '0')
  const timeString = `${h}:${m}:${s}`
  emit('update:modelValue', timeString)
  showPicker.value = false
}

// 点击外部关闭选择器
const handleClickOutside = (e: MouseEvent) => {
  const target = e.target as HTMLElement
  if (showPicker.value && !target.closest('.relative')) {
    showPicker.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>
