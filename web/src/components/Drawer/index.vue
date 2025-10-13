<template>
  <Teleport to="body">
    <Transition name="drawer">
      <div v-if="modelValue" class="fixed inset-0 z-50 flex" @click.self="onClose">
        <!-- 背景遮罩 -->
        <div class="absolute inset-0 bg-black bg-opacity-50 transition-opacity"></div>

        <!-- 抽屉内容 -->
        <div
          :class="[
            'relative bg-bg-elevated shadow-2xl flex flex-col transition-all duration-300',
            positionClass,
            sizeClass
          ]"
        >
          <!-- 标题栏 -->
          <div class="px-6 py-4 border-b-2 border-border-primary flex items-center justify-between flex-shrink-0">
            <h3 class="text-lg font-semibold text-text-primary">{{ title }}</h3>
            <button
              @click="onClose"
              class="p-1 text-text-tertiary hover:text-text-secondary hover:bg-bg-tertiary rounded transition-colors"
            >
              <X :size="20" />
            </button>
          </div>

          <!-- 内容区 - 可滚动 -->
          <div class="flex-1 overflow-y-auto px-6 py-4">
            <slot></slot>
          </div>

          <!-- 底部按钮栏 -->
          <div v-if="showFooter" class="px-6 py-4 border-t-2 border-border-primary flex gap-3 justify-end flex-shrink-0">
            <BaseButton v-if="cancelText" variant="secondary" @click="onCancel">
              {{ cancelText }}
            </BaseButton>
            <BaseButton v-if="confirmText" :variant="confirmVariant" @click="onConfirm">
              {{ confirmText }}
            </BaseButton>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { computed, watch, onMounted, onUnmounted } from 'vue'
import { X } from 'lucide-vue-next'
import BaseButton from '../BaseButton/index.vue'

interface Props {
  modelValue: boolean
  title?: string
  position?: 'left' | 'right'
  size?: 'sm' | 'md' | 'lg' | 'xl'
  confirmText?: string
  cancelText?: string
  confirmVariant?: 'primary' | 'secondary' | 'danger' | 'success' | 'ghost'
  showFooter?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  title: '',
  position: 'right',
  size: 'lg',
  confirmText: '',
  cancelText: '',
  confirmVariant: 'primary',
  showFooter: true
})

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  confirm: []
  cancel: []
  close: []
}>()

const positionClass = computed(() => {
  return props.position === 'right' ? 'ml-auto' : 'mr-auto'
})

const sizeClass = computed(() => {
  const sizes = {
    sm: 'w-80',
    md: 'w-96',
    lg: 'w-[32rem]',
    xl: 'w-[40rem]'
  }
  return sizes[props.size]
})

const onConfirm = () => {
  emit('confirm')
}

const onCancel = () => {
  emit('cancel')
  emit('update:modelValue', false)
  emit('close')
}

const onClose = () => {
  emit('update:modelValue', false)
  emit('close')
}

// ESC 键关闭抽屉
const handleEscapeKey = (event: KeyboardEvent) => {
  if (event.key === 'Escape' && props.modelValue) {
    onClose()
  }
}

// 监听 modelValue 变化，动态添加/移除事件监听器
watch(() => props.modelValue, (isOpen) => {
  if (isOpen) {
    document.addEventListener('keydown', handleEscapeKey)
  } else {
    document.removeEventListener('keydown', handleEscapeKey)
  }
})

// 组件挂载时，如果抽屉已打开，添加监听器
onMounted(() => {
  if (props.modelValue) {
    document.addEventListener('keydown', handleEscapeKey)
  }
})

// 组件卸载时移除监听器
onUnmounted(() => {
  document.removeEventListener('keydown', handleEscapeKey)
})
</script>

<style scoped>
/* 抽屉进入/离开动画 */
.drawer-enter-active,
.drawer-leave-active {
  transition: opacity 0.3s ease;
}

.drawer-enter-from,
.drawer-leave-to {
  opacity: 0;
}

/* 抽屉面板滑动动画 */
.drawer-enter-active > div:not(.absolute),
.drawer-leave-active > div:not(.absolute) {
  transition: transform 0.3s ease;
}

.drawer-enter-from > div:not(.absolute) {
  transform: translateX(100%);
}

.drawer-leave-to > div:not(.absolute) {
  transform: translateX(100%);
}

/* 左侧抽屉动画 */
.drawer-enter-from > div.mr-auto,
.drawer-leave-to > div.mr-auto {
  transform: translateX(-100%);
}
</style>
