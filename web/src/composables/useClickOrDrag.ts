import { ref } from 'vue'
import { DRAG_THRESHOLD } from '@/constants/interaction'

/**
 * 区分点击和拖拽行为
 * 当鼠标移动距离超过阈值时，判定为拖拽而非点击
 *
 * @param threshold 拖拽判定阈值（像素），默认为 DRAG_THRESHOLD
 * @returns 鼠标事件处理器和状态
 *
 * @example
 * ```vue
 * <script setup>
 * const { handleMouseDown, handleMouseMove, handleMouseUp, onClick } = useClickOrDrag()
 *
 * onClick(() => {
 *   // 处理点击事件
 * })
 * </script>
 *
 * <template>
 *   <div
 *     @mousedown="handleMouseDown"
 *     @mousemove="handleMouseMove"
 *     @mouseup="handleMouseUp"
 *   >
 *     点击或拖拽我
 *   </div>
 * </template>
 * ```
 */
export function useClickOrDrag(threshold = DRAG_THRESHOLD) {
  const mouseDownPos = ref<{ x: number; y: number } | null>(null)
  const hasMoved = ref(false)
  let clickCallback: (() => void) | null = null

  const handleMouseDown = (e: MouseEvent) => {
    mouseDownPos.value = { x: e.clientX, y: e.clientY }
    hasMoved.value = false
  }

  const handleMouseMove = (e: MouseEvent) => {
    if (!mouseDownPos.value) return

    const deltaX = Math.abs(e.clientX - mouseDownPos.value.x)
    const deltaY = Math.abs(e.clientY - mouseDownPos.value.y)

    if (deltaX > threshold || deltaY > threshold) {
      hasMoved.value = true
    }
  }

  const handleMouseUp = () => {
    if (!hasMoved.value && mouseDownPos.value && clickCallback) {
      clickCallback()
    }

    mouseDownPos.value = null
    hasMoved.value = false
  }

  const onClick = (callback: () => void) => {
    clickCallback = callback
  }

  return {
    handleMouseDown,
    handleMouseMove,
    handleMouseUp,
    onClick,
    hasMoved,
  }
}

