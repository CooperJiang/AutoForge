<template>
  <Teleport to="body">
    <Transition name="dialog">
      <div
        v-if="modelValue"
        class="fixed inset-0 z-50 flex items-center justify-center"
        @click.self="handleCancel"
      >
        
        <div class="absolute inset-0 bg-black/50 backdrop-blur-sm"></div>

        
        <div
          class="relative bg-bg-elevated rounded-lg shadow-xl max-w-md w-full mx-4 overflow-hidden"
        >
          
          <div class="px-6 py-4 border-b border-border-primary">
            <div class="flex items-center gap-3">
              <div
                :class="[
                  'w-10 h-10 rounded-full flex items-center justify-center',
                  variantClasses[variant].iconBg,
                ]"
              >
                <component
                  :is="variantIcons[variant]"
                  :class="['w-5 h-5', variantClasses[variant].icon]"
                />
              </div>
              <h3 class="text-lg font-semibold text-text-primary">
                {{ title }}
              </h3>
            </div>
          </div>

          
          <div class="px-6 py-4">
            <p class="text-text-secondary leading-relaxed">{{ message }}</p>
          </div>

          
          <div
            class="px-6 py-4 bg-bg-hover border-t border-border-primary flex items-center justify-end gap-3"
          >
            <BaseButton size="sm" variant="ghost" @click="handleCancel">
              {{ cancelText }}
            </BaseButton>
            <BaseButton
              size="sm"
              :variant="variant === 'danger' ? 'danger' : 'primary'"
              @click="handleConfirm"
            >
              {{ confirmText }}
            </BaseButton>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { AlertCircle, AlertTriangle, Info, HelpCircle } from 'lucide-vue-next'
import BaseButton from '@/components/BaseButton'

interface Props {
  modelValue: boolean
  title?: string
  message: string
  confirmText?: string
  cancelText?: string
  variant?: 'info' | 'warning' | 'danger' | 'question'
}

withDefaults(defineProps<Props>(), {
  title: '确认',
  confirmText: '确定',
  cancelText: '取消',
  variant: 'question',
})

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  confirm: []
  cancel: []
}>()

const variantIcons = {
  info: Info,
  warning: AlertTriangle,
  danger: AlertCircle,
  question: HelpCircle,
}

const variantClasses = {
  info: {
    iconBg: 'bg-primary-light',
    icon: 'text-primary',
  },
  warning: {
    iconBg: 'bg-amber-100',
    icon: 'text-amber-600',
  },
  danger: {
    iconBg: 'bg-red-100',
    icon: 'text-red-600',
  },
  question: {
    iconBg: 'bg-bg-tertiary',
    icon: 'text-text-secondary',
  },
}

const handleConfirm = () => {
  emit('confirm')
  emit('update:modelValue', false)
}

const handleCancel = () => {
  emit('cancel')
  emit('update:modelValue', false)
}
</script>

<style scoped>
.dialog-enter-active,
.dialog-leave-active {
  transition: opacity 0.2s ease;
}

.dialog-enter-from,
.dialog-leave-to {
  opacity: 0;
}

.dialog-enter-active > div:last-child,
.dialog-leave-active > div:last-child {
  transition: transform 0.2s ease;
}

.dialog-enter-from > div:last-child,
.dialog-leave-to > div:last-child {
  transform: scale(0.95);
}
</style>
