<template>
  <Teleport to="body">
    <Transition name="dialog">
      <div
        v-if="modelValue"
        class="fixed inset-0 z-50 flex items-center justify-center"
        @click.self="handleCancel"
      >
        <!-- 遮罩层 -->
        <div class="absolute inset-0 bg-black/50 backdrop-blur-sm"></div>

        <!-- 对话框 -->
        <div class="relative bg-white rounded-lg shadow-xl max-w-md w-full mx-4 overflow-hidden">
          <!-- 头部 -->
          <div class="px-6 py-4 border-b border-slate-200">
            <div class="flex items-center gap-3">
              <div
                :class="[
                  'w-10 h-10 rounded-full flex items-center justify-center',
                  variantClasses[variant].iconBg
                ]"
              >
                <component :is="variantIcons[variant]" :class="['w-5 h-5', variantClasses[variant].icon]" />
              </div>
              <h3 class="text-lg font-semibold text-slate-900">
                {{ title }}
              </h3>
            </div>
          </div>

          <!-- 内容 -->
          <div class="px-6 py-4">
            <p class="text-slate-600 leading-relaxed">{{ message }}</p>
          </div>

          <!-- 底部按钮 -->
          <div class="px-6 py-4 bg-slate-50 border-t border-slate-200 flex items-center justify-end gap-3">
            <BaseButton
              size="sm"
              variant="ghost"
              @click="handleCancel"
            >
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
import BaseButton from './BaseButton'

interface Props {
  modelValue: boolean
  title?: string
  message: string
  confirmText?: string
  cancelText?: string
  variant?: 'info' | 'warning' | 'danger' | 'question'
}

const props = withDefaults(defineProps<Props>(), {
  title: '确认',
  confirmText: '确定',
  cancelText: '取消',
  variant: 'question'
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
  question: HelpCircle
}

const variantClasses = {
  info: {
    iconBg: 'bg-blue-100',
    icon: 'text-blue-600'
  },
  warning: {
    iconBg: 'bg-amber-100',
    icon: 'text-amber-600'
  },
  danger: {
    iconBg: 'bg-red-100',
    icon: 'text-red-600'
  },
  question: {
    iconBg: 'bg-slate-100',
    icon: 'text-slate-600'
  }
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
