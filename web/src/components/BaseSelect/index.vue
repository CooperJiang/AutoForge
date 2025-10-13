<template>
  <div class="w-full relative" ref="selectRef">
    <label v-if="label" class="block text-sm font-medium text-text-primary mb-2">
      {{ label }}
      <span v-if="required" class="text-error ml-1">*</span>
    </label>
    <button
      type="button"
      @click="toggleDropdown"
      class="w-full px-3 py-1.5 text-sm text-left bg-bg-primary border-2 border-border-primary rounded-md transition-all duration-200 hover:border-border-secondary focus:border-border-focus focus:ring-2 focus:ring-primary-light focus:outline-none"
      :class="{ 'border-border-focus ring-2 ring-primary-light': isOpen }"
    >
      <div class="flex items-center justify-between">
        <span class="text-text-primary">{{ selectedLabel || placeholder }}</span>
        <svg
          class="w-5 h-5 text-text-tertiary transition-transform duration-200"
          :class="{ 'rotate-180': isOpen }"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
        </svg>
      </div>
    </button>

    <!-- Dropdown using Teleport to body -->
    <Teleport to="body">
      <Transition
        enter-active-class="transition ease-out duration-100"
        enter-from-class="transform opacity-0 scale-95"
        enter-to-class="transform opacity-100 scale-100"
        leave-active-class="transition ease-in duration-75"
        leave-from-class="transform opacity-100 scale-100"
        leave-to-class="transform opacity-0 scale-95"
      >
        <div
          v-show="isOpen"
          class="fixed bg-bg-elevated border-2 border-border-primary rounded-md shadow-lg max-h-48 overflow-y-auto"
          :style="dropdownStyle"
        >
          <div
            v-for="option in options"
            :key="option.value"
            @click="selectOption(option)"
            class="px-3 py-1.5 text-sm cursor-pointer transition-colors duration-150 hover:bg-bg-hover text-text-primary"
            :class="{ 'bg-primary-light font-medium': option.value === modelValue }"
          >
            {{ option.label }}
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'

interface Option {
  label: string
  value: string
}

interface Props {
  modelValue: string
  options: Option[]
  label?: string
  placeholder?: string
  required?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: '请选择',
  required: false
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const isOpen = ref(false)
const selectRef = ref<HTMLElement>()
const dropdownPosition = ref({ top: 0, left: 0, width: 0 })

const selectedLabel = computed(() => {
  const option = props.options.find(opt => opt.value === props.modelValue)
  return option?.label || ''
})

const dropdownStyle = computed(() => ({
  top: `${dropdownPosition.value.top}px`,
  left: `${dropdownPosition.value.left}px`,
  width: `${dropdownPosition.value.width}px`,
  zIndex: 9999
}))

const calculatePosition = () => {
  if (!selectRef.value) return

  const rect = selectRef.value.getBoundingClientRect()
  const button = selectRef.value.querySelector('button')
  const buttonRect = button?.getBoundingClientRect()

  if (buttonRect) {
    dropdownPosition.value = {
      top: buttonRect.bottom + 4,
      left: buttonRect.left,
      width: buttonRect.width
    }
  }
}

const toggleDropdown = () => {
  if (!isOpen.value) {
    calculatePosition()
  }
  isOpen.value = !isOpen.value
}

const selectOption = (option: Option) => {
  emit('update:modelValue', option.value)
  isOpen.value = false
}

const handleClickOutside = (event: MouseEvent) => {
  if (selectRef.value && !selectRef.value.contains(event.target as Node)) {
    isOpen.value = false
  }
}

const handleScroll = () => {
  if (isOpen.value) {
    calculatePosition()
  }
}

onMounted(() => {
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
