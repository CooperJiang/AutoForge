/**
 * 图片交互逻辑 Hook
 * 负责：拖拽、滚轮缩放、键盘快捷键
 */

import { ref, onMounted, onUnmounted, type Ref } from 'vue'

interface UseImageInteractionOptions {
  translateX: Ref<number>
  translateY: Ref<number>
  onZoomIn: () => void
  onZoomOut: () => void
  onRotateLeft: () => void
  onRotateRight: () => void
  onReset: () => void
  onClose: () => void
}

export function useImageInteraction(options: UseImageInteractionOptions) {
  const { translateX, translateY, onZoomIn, onZoomOut, onRotateLeft, onRotateRight, onReset, onClose } =
    options

  // 拖拽状态
  const isDragging = ref(false)
  const startX = ref(0)
  const startY = ref(0)
  const initialTranslateX = ref(0)
  const initialTranslateY = ref(0)

  // 开始拖动
  const startDrag = (e: MouseEvent) => {
    isDragging.value = true
    startX.value = e.clientX
    startY.value = e.clientY
    initialTranslateX.value = translateX.value
    initialTranslateY.value = translateY.value
  }

  // 拖动中
  const onDrag = (e: MouseEvent) => {
    if (!isDragging.value) return

    const deltaX = e.clientX - startX.value
    const deltaY = e.clientY - startY.value

    translateX.value = initialTranslateX.value + deltaX
    translateY.value = initialTranslateY.value + deltaY
  }

  // 停止拖动
  const stopDrag = () => {
    isDragging.value = false
  }

  // 滚轮缩放
  const onWheel = (e: WheelEvent) => {
    e.preventDefault()

    if (e.deltaY < 0) {
      // 向上滚动 - 放大
      onZoomIn()
    } else {
      // 向下滚动 - 缩小
      onZoomOut()
    }
  }

  // 键盘快捷键
  const handleKeyDown = (e: KeyboardEvent) => {
    switch (e.key) {
      case 'Escape':
        onClose()
        break
      case '+':
      case '=':
        onZoomIn()
        break
      case '-':
      case '_':
        onZoomOut()
        break
      case 'ArrowLeft':
        onRotateLeft()
        break
      case 'ArrowRight':
        onRotateRight()
        break
      case 'r':
      case 'R':
        onReset()
        break
    }
  }

  // 注册和注销事件监听
  const register = () => {
    document.addEventListener('keydown', handleKeyDown)
  }

  const unregister = () => {
    document.removeEventListener('keydown', handleKeyDown)
  }

  // 生命周期
  onMounted(() => {
    register()
  })

  onUnmounted(() => {
    unregister()
  })

  return {
    // 状态
    isDragging,

    // 方法
    startDrag,
    onDrag,
    stopDrag,
    onWheel,
    register,
    unregister,
  }
}

