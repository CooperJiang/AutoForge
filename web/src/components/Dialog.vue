<template>
  <Teleport to="body">
    <Transition name="dialog">
      <div v-if="modelValue" class="fixed inset-0 z-50 flex items-center justify-center p-4" @click.self="onCancel">
        <!-- 背景遮罩 -->
        <div class="absolute inset-0 bg-black bg-opacity-50 transition-opacity"></div>

        <!-- 对话框内容 -->
        <div class="relative bg-white rounded-lg shadow-xl border-2 border-slate-200 w-full max-w-2xl max-h-[90vh] flex flex-col transform transition-all">
          <!-- 标题 -->
          <div class="px-5 py-4 border-b-2 border-slate-100 flex items-center justify-between flex-shrink-0">
            <h3 class="text-base font-semibold text-slate-900">{{ title }}</h3>
            <button @click="onCancel" class="text-slate-400 hover:text-slate-600 transition-colors">
              <X :size="20" />
            </button>
          </div>

          <!-- 内容 - 可滚动 -->
          <div class="px-5 py-4 overflow-y-auto flex-1">
            <slot>
              <p class="text-sm text-slate-600">{{ message }}</p>
            </slot>
          </div>

          <!-- 按钮 -->
          <div class="px-5 py-4 border-t-2 border-slate-100 flex gap-2 justify-end flex-shrink-0">
            <BaseButton v-if="cancelText" variant="secondary" @click="onCancel">
              {{ cancelText }}
            </BaseButton>
            <BaseButton :variant="confirmVariant" @click="onConfirm">
              {{ confirmText }}
            </BaseButton>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { X } from 'lucide-vue-next'
import BaseButton from './BaseButton.vue'

interface Props {
  modelValue: boolean
  title?: string
  message?: string
  confirmText?: string
  cancelText?: string
  confirmVariant?: 'primary' | 'secondary' | 'danger' | 'success' | 'ghost'
}

const props = withDefaults(defineProps<Props>(), {
  title: '确认',
  message: '确定要执行此操作吗？',
  confirmText: '确定',
  cancelText: '取消',
  confirmVariant: 'primary'
})

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  confirm: []
  cancel: []
}>()

const onConfirm = () => {
  emit('confirm')
  emit('update:modelValue', false)
}

const onCancel = () => {
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

.dialog-enter-active .relative,
.dialog-leave-active .relative {
  transition: transform 0.2s ease, opacity 0.2s ease;
}

.dialog-enter-from .relative,
.dialog-leave-to .relative {
  transform: scale(0.95);
  opacity: 0;
}
</style>
