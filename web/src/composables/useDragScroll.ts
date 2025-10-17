import { ref, watch, onUnmounted, type Ref } from 'vue'

/**
 * 为容器添加鼠标拖拽滚动功能
 * @param containerRef 容器元素的 ref
 * @returns 拖拽状态
 */
export function useDragScroll(containerRef: Ref<HTMLElement | null>) {
  const isDragging = ref(false)
  const startX = ref(0)
  const scrollLeft = ref(0)
  let currentContainer: HTMLElement | null = null

  const handleMouseDown = (e: MouseEvent) => {
    if (!containerRef.value) return

    isDragging.value = true
    startX.value = e.clientX
    scrollLeft.value = containerRef.value.scrollLeft

    containerRef.value.style.cursor = 'grabbing'
    containerRef.value.style.userSelect = 'none'
  }

  const handleMouseMove = (e: MouseEvent) => {
    if (!isDragging.value || !containerRef.value) return

    e.preventDefault()

    const deltaX = e.clientX - startX.value
    containerRef.value.scrollLeft = scrollLeft.value - deltaX
  }

  const handleMouseUp = () => {
    if (!containerRef.value) return

    isDragging.value = false
    containerRef.value.style.cursor = 'grab'
    containerRef.value.style.userSelect = ''
  }

  const handleMouseLeave = () => {
    if (!containerRef.value) return

    isDragging.value = false
    containerRef.value.style.cursor = 'grab'
    containerRef.value.style.userSelect = ''
  }

  const bindEvents = (container: HTMLElement) => {
    container.style.cursor = 'grab'
    container.addEventListener('mousedown', handleMouseDown)
    container.addEventListener('mousemove', handleMouseMove)
    container.addEventListener('mouseup', handleMouseUp)
    container.addEventListener('mouseleave', handleMouseLeave)
  }

  const unbindEvents = (container: HTMLElement) => {
    container.removeEventListener('mousedown', handleMouseDown)
    container.removeEventListener('mousemove', handleMouseMove)
    container.removeEventListener('mouseup', handleMouseUp)
    container.removeEventListener('mouseleave', handleMouseLeave)
  }

  watch(
    containerRef,
    (newContainer, oldContainer) => {
      if (oldContainer) {
        unbindEvents(oldContainer)
      }

      if (newContainer) {
        bindEvents(newContainer)
        currentContainer = newContainer
      }
    },
    { immediate: true }
  )

  onUnmounted(() => {
    if (currentContainer) {
      unbindEvents(currentContainer)
    }
  })

  return {
    isDragging,
  }
}
