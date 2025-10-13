<template>
  <div class="w-full relative" ref="pickerRef">
    <label v-if="label" class="block text-sm font-medium text-text-primary mb-2">
      {{ label }}
      <span v-if="required" class="text-error ml-1">*</span>
    </label>
    <button
      type="button"
      @click="togglePicker"
      class="w-full px-3 py-1.5 text-sm text-left text-text-primary bg-bg-primary border-2 border-border-primary rounded-md transition-all duration-200 hover:border-border-secondary focus:border-border-focus focus:ring-2 focus:ring-primary-light focus:outline-none"
      :class="{ 'border-border-focus ring-2 ring-primary-light': isOpen }"
    >
      <div class="flex items-center justify-between">
        <span class="font-mono">{{ displayTime }}</span>
        <svg class="w-4 h-4 text-text-tertiary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      </div>
    </button>

    <!-- Time Picker Dropdown using Teleport to body -->
    <Teleport to="body">
      <Transition
        enter-active-class="transition ease-out duration-200"
        enter-from-class="transform opacity-0 scale-95 translate-y-2"
        enter-to-class="transform opacity-100 scale-100 translate-y-0"
        leave-active-class="transition ease-in duration-150"
        leave-from-class="transform opacity-100 scale-100 translate-y-0"
        leave-to-class="transform opacity-0 scale-95 translate-y-2"
      >
        <div
          v-show="isOpen"
          class="fixed z-50 bg-bg-elevated border border-border-secondary rounded-xl shadow-2xl overflow-hidden"
          :style="dropdownStyle"
        >
          <!-- 时间轮选择器 -->
          <div class="flex items-center justify-center gap-1 p-4 bg-gradient-to-b from-bg-secondary to-bg-elevated">
            <!-- Hours -->
            <div class="flex-1 relative">
              <div class="text-xs text-text-tertiary text-center mb-2 font-medium">时</div>
              <div class="h-40 overflow-y-auto scroll-smooth scrollbar-thin" ref="hoursRef">
                <div class="py-14">
                  <div
                    v-for="h in hours"
                    :key="h"
                    @click="selectHour(h)"
                    class="h-8 flex items-center justify-center cursor-pointer text-center transition-all duration-200 font-mono text-sm relative"
                    :class="getTimeItemClass(h, selectedHour)"
                  >
                    {{ h.toString().padStart(2, '0') }}
                  </div>
                </div>
              </div>
              <!-- 选中指示器 -->
              <div class="absolute top-[52px] left-0 right-0 h-8 bg-primary/10 border-y-2 border-primary pointer-events-none rounded-md"></div>
            </div>

            <div class="text-xl text-text-primary font-bold py-1">:</div>

            <!-- Minutes -->
            <div class="flex-1 relative">
              <div class="text-xs text-text-tertiary text-center mb-2 font-medium">分</div>
              <div class="h-40 overflow-y-auto scroll-smooth scrollbar-thin" ref="minutesRef">
                <div class="py-14">
                  <div
                    v-for="m in minutes"
                    :key="m"
                    @click="selectMinute(m)"
                    class="h-8 flex items-center justify-center cursor-pointer text-center transition-all duration-200 font-mono text-sm relative"
                    :class="getTimeItemClass(m, selectedMinute)"
                  >
                    {{ m.toString().padStart(2, '0') }}
                  </div>
                </div>
              </div>
              <!-- 选中指示器 -->
              <div class="absolute top-[52px] left-0 right-0 h-8 bg-primary/10 border-y-2 border-primary pointer-events-none rounded-md"></div>
            </div>

            <div class="text-xl text-text-primary font-bold py-1">:</div>

            <!-- Seconds -->
            <div class="flex-1 relative">
              <div class="text-xs text-text-tertiary text-center mb-2 font-medium">秒</div>
              <div class="h-40 overflow-y-auto scroll-smooth scrollbar-thin" ref="secondsRef">
                <div class="py-14">
                  <div
                    v-for="s in seconds"
                    :key="s"
                    @click="selectSecond(s)"
                    class="h-8 flex items-center justify-center cursor-pointer text-center transition-all duration-200 font-mono text-sm relative"
                    :class="getTimeItemClass(s, selectedSecond)"
                  >
                    {{ s.toString().padStart(2, '0') }}
                  </div>
                </div>
              </div>
              <!-- 选中指示器 -->
              <div class="absolute top-[52px] left-0 right-0 h-8 bg-primary/10 border-y-2 border-primary pointer-events-none rounded-md"></div>
            </div>
          </div>

          <!-- 底部按钮 -->
          <div class="px-4 py-3 bg-bg-hover border-t-2 border-border-secondary flex justify-end gap-2">
            <button
              type="button"
              @click="isOpen = false"
              class="px-3 py-1.5 text-xs bg-bg-elevated text-text-secondary rounded-lg hover:bg-bg-hover transition-all duration-200 font-medium border border-border-primary"
            >
              取消
            </button>
            <button
              type="button"
              @click="confirmTime"
              class="px-3 py-1.5 text-xs bg-[var(--color-primary)] text-white rounded-lg hover:bg-primary-hover transition-all duration-200 font-medium shadow-md hover:shadow-lg"
            >
              确定
            </button>
          </div>
        </div>
      </Transition>
    </Teleport>

    <p v-if="hint" class="mt-1.5 text-xs text-text-tertiary">{{ hint }}</p>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'

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
const hoursRef = ref<HTMLElement>()
const minutesRef = ref<HTMLElement>()
const secondsRef = ref<HTMLElement>()
const dropdownPosition = ref({ top: 0, left: 0, width: 0 })

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

const dropdownStyle = computed(() => ({
  top: `${dropdownPosition.value.top}px`,
  left: `${dropdownPosition.value.left}px`,
  minWidth: `${dropdownPosition.value.width}px`,
  zIndex: 9999
}))

const calculatePosition = () => {
  if (!pickerRef.value) return

  const button = pickerRef.value.querySelector('button')
  const buttonRect = button?.getBoundingClientRect()

  if (buttonRect) {
    dropdownPosition.value = {
      top: buttonRect.bottom + 4,
      left: buttonRect.left,
      width: buttonRect.width
    }
  }
}

const scrollToSelected = () => {
  nextTick(() => {
    if (hoursRef.value) {
      const hourElement = hoursRef.value.querySelector(`[data-hour="${selectedHour.value}"]`) as HTMLElement
      if (hourElement) {
        hoursRef.value.scrollTop = hourElement.offsetTop - hoursRef.value.offsetHeight / 2 + hourElement.offsetHeight / 2
      } else {
        hoursRef.value.scrollTop = selectedHour.value * 32 - 80
      }
    }

    if (minutesRef.value) {
      minutesRef.value.scrollTop = selectedMinute.value * 32 - 80
    }

    if (secondsRef.value) {
      secondsRef.value.scrollTop = selectedSecond.value * 32 - 80
    }
  })
}

const togglePicker = () => {
  if (!isOpen.value) {
    calculatePosition()
    scrollToSelected()
  }
  isOpen.value = !isOpen.value
}

const getTimeItemClass = (value: number, selected: number) => {
  if (value === selected) {
    return 'text-primary font-bold scale-110'
  }
  return 'text-text-secondary hover:text-primary hover:scale-105'
}

const selectHour = (h: number) => {
  selectedHour.value = h
  if (hoursRef.value) {
    hoursRef.value.scrollTop = h * 32 - 80
  }
}

const selectMinute = (m: number) => {
  selectedMinute.value = m
  if (minutesRef.value) {
    minutesRef.value.scrollTop = m * 32 - 80
  }
}

const selectSecond = (s: number) => {
  selectedSecond.value = s
  if (secondsRef.value) {
    secondsRef.value.scrollTop = s * 32 - 80
  }
}

const confirmTime = () => {
  emit('update:modelValue', displayTime.value)
  isOpen.value = false
}

const handleClickOutside = (event: MouseEvent) => {
  if (pickerRef.value && !pickerRef.value.contains(event.target as Node)) {
    const dropdown = document.querySelector('[data-timepicker-dropdown]')
    if (dropdown && dropdown.contains(event.target as Node)) {
      return
    }
    isOpen.value = false
  }
}

const handleScroll = () => {
  if (isOpen.value) {
    calculatePosition()
  }
}

watch(() => props.modelValue, () => {
  parseTime()
})

watch(isOpen, (newVal) => {
  if (newVal) {
    scrollToSelected()
  }
})

onMounted(() => {
  parseTime()
  document.addEventListener('click', handleClickOutside)
  window.addEventListener('scroll', handleScroll, true)
  window.addEventListener('resize', handleScroll)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
  window.removeEventListener('scroll', handleScroll, true)
  window.removeEventListener('resize', handleScroll)
})
</script>

<style scoped>
.scrollbar-thin::-webkit-scrollbar {
  width: 4px;
}

.scrollbar-thin::-webkit-scrollbar-track {
  background: transparent;
}

.scrollbar-thin::-webkit-scrollbar-thumb {
  background: var(--color-border-secondary);
  border-radius: 2px;
}

.scrollbar-thin::-webkit-scrollbar-thumb:hover {
  background: var(--color-primary);
}
</style>
