<template>
  <div class="image-viewer">
    
    <div class="relative inline-block">
      <img
        ref="thumbnailRef"
        :src="src"
        :alt="alt"
        class="rounded-lg border border-border-primary shadow-sm cursor-pointer hover:opacity-90 transition-opacity"
        style="max-width: 200px; max-height: 200px; object-fit: contain"
        @error="handleError"
        @click="openPreview"
      />
      <div
        class="absolute inset-0 flex items-center justify-center opacity-0 hover:opacity-100 transition-opacity bg-black/50 rounded-lg cursor-pointer"
        @click="openPreview"
      >
        <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0zM10 7v6m3-3H7"
          />
        </svg>
      </div>
    </div>

    <div v-if="description" class="mt-2 text-sm text-text-secondary">
      {{ description }}
    </div>

    
    <Teleport to="body">
      <Transition name="fade-bg">
        <div v-if="showPreview" class="fixed inset-0 z-50 bg-black/90 backdrop-blur-sm">
          
          <div
            class="absolute top-0 left-0 right-0 h-16 flex items-center justify-between px-6 bg-black/50 z-20"
            @click.stop
          >
            <div class="flex items-center gap-4">
              
              <div class="flex items-center gap-2 bg-white/10 rounded-lg px-3 py-2">
                <button
                  class="p-1 hover:bg-white/20 rounded transition-colors"
                  @click.stop="zoomOut"
                  title="缩小 (滚轮向下)"
                >
                  <svg
                    class="w-5 h-5 text-white"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0zM13 10H7"
                    />
                  </svg>
                </button>
                <span class="text-white text-sm min-w-[60px] text-center"
                  >{{ Math.round(scale * 100) }}%</span
                >
                <button
                  class="p-1 hover:bg-white/20 rounded transition-colors"
                  @click.stop="zoomIn"
                  title="放大 (滚轮向上)"
                >
                  <svg
                    class="w-5 h-5 text-white"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0zM10 7v6m3-3H7"
                    />
                  </svg>
                </button>
              </div>

              
              <button
                class="p-2 bg-white/10 hover:bg-white/20 rounded-lg transition-colors"
                @click.stop="rotateLeft"
                title="逆时针旋转"
              >
                <svg
                  class="w-5 h-5 text-white"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6"
                  />
                </svg>
              </button>
              <button
                class="p-2 bg-white/10 hover:bg-white/20 rounded-lg transition-colors"
                @click.stop="rotateRight"
                title="顺时针旋转"
              >
                <svg
                  class="w-5 h-5 text-white"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M21 10H11a8 8 0 00-8 8v2m18-10l-6 6m6-6l-6-6"
                  />
                </svg>
              </button>

              
              <button
                class="p-2 bg-white/10 hover:bg-white/20 rounded-lg transition-colors"
                @click.stop="reset"
                title="重置 (双击图片)"
              >
                <svg
                  class="w-5 h-5 text-white"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
                  />
                </svg>
              </button>
            </div>

            <div class="flex items-center gap-2">
              
              <button
                class="px-4 py-2 bg-white/10 hover:bg-white/20 rounded-lg transition-colors text-white text-sm"
                @click.stop="openInNewTab"
                title="在新标签页打开"
              >
                新窗口打开
              </button>

              
              <button
                class="p-2 hover:bg-white/20 rounded-lg transition-colors"
                @click.stop="closePreview"
                title="关闭 (ESC)"
              >
                <svg
                  class="w-6 h-6 text-white"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M6 18L18 6M6 6l12 12"
                  />
                </svg>
              </button>
            </div>
          </div>

          
          <div
            ref="containerRef"
            class="absolute inset-0 overflow-hidden z-10"
            style="padding-top: 64px"
            @mousedown="startDrag"
            @mousemove="onDrag"
            @mouseup="stopDrag"
            @mouseleave="stopDrag"
            @wheel.prevent="onWheel"
            @click="closePreview"
          >
            <Transition name="zoom" @enter="onEnter" @leave="onLeave">
              <img
                v-if="showPreview"
                ref="imageRef"
                :src="src"
                :alt="alt"
                class="absolute select-none"
                :style="imageStyle"
                @click.stop
                @dblclick.stop="reset"
                @dragstart.prevent
              />
            </Transition>
          </div>

          
          <div
            class="absolute bottom-6 left-1/2 -translate-x-1/2 bg-black/50 text-white text-xs px-4 py-2 rounded-lg pointer-events-none z-20"
          >
            拖动图片移动 | 滚轮缩放 | 双击重置
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'

interface Props {
  src: string
  alt?: string
  description?: string
}

const props = withDefaults(defineProps<Props>(), {
  alt: '图片',
})

const emit = defineEmits<{
  error: [event: Event]
}>()

const showPreview = ref(false)
const scale = ref(1)
const rotation = ref(0)
const translateX = ref(0)
const translateY = ref(0)
const isDragging = ref(false)
const dragStartX = ref(0)
const dragStartY = ref(0)
const containerRef = ref<HTMLElement>()
const imageRef = ref<HTMLImageElement>()
const thumbnailRef = ref<HTMLImageElement>()
const isAnimating = ref(false)

// 图片样式
const imageStyle = computed(() => ({
  transform: `translate(${translateX.value}px, ${translateY.value}px) scale(${scale.value}) rotate(${rotation.value}deg)`,
  transition: isDragging.value || isAnimating.value ? 'none' : 'transform 0.3s ease',
  cursor: isDragging.value ? 'grabbing' : 'grab',
  transformOrigin: 'center center',
}))

// 打开预览
const openPreview = () => {
  if (!thumbnailRef.value) return

  // 获取缩略图位置
  const thumbRect = thumbnailRef.value.getBoundingClientRect()

  // 设置初始位置和缩放（从缩略图位置开始）
  translateX.value = thumbRect.left
  translateY.value = thumbRect.top
  scale.value = thumbRect.width / thumbnailRef.value.naturalWidth
  rotation.value = 0

  showPreview.value = true
  isAnimating.value = true

  // 等待 DOM 更新后执行动画
  setTimeout(() => {
    centerImage()
    scale.value = 1

    // 动画结束后重置状态
    setTimeout(() => {
      isAnimating.value = false
    }, 300)
  }, 50)
}

// 关闭预览
const closePreview = () => {
  showPreview.value = false
}

// 放大
const zoomIn = () => {
  scale.value = Math.min(scale.value + 0.25, 5)
}

// 缩小
const zoomOut = () => {
  scale.value = Math.max(scale.value - 0.25, 0.1)
}

// 左旋转
const rotateLeft = () => {
  rotation.value -= 90
}

// 右旋转
const rotateRight = () => {
  rotation.value += 90
}

// 重置
const reset = () => {
  scale.value = 1
  rotation.value = 0
  translateX.value = 0
  translateY.value = 0
  centerImage()
}

// 居中图片
const centerImage = () => {
  if (!imageRef.value || !containerRef.value) {
    // 如果元素还未准备好，延迟重试
    setTimeout(() => {
      if (showPreview.value) {
        centerImage()
      }
    }, 100)
    return
  }

  const img = imageRef.value

  // 如果图片还没加载，等待加载完成
  if (img.naturalWidth === 0) {
    img.onload = () => {
      centerImage()
    }
    return
  }

  const container = containerRef.value.getBoundingClientRect()

  // 计算居中位置
  translateX.value = (container.width - img.naturalWidth) / 2
  translateY.value = (container.height - img.naturalHeight) / 2
}

// 开始拖动
const startDrag = (e: MouseEvent) => {
  isDragging.value = true
  dragStartX.value = e.clientX - translateX.value
  dragStartY.value = e.clientY - translateY.value
}

// 拖动中
const onDrag = (e: MouseEvent) => {
  if (!isDragging.value) return
  translateX.value = e.clientX - dragStartX.value
  translateY.value = e.clientY - dragStartY.value
}

// 停止拖动
const stopDrag = () => {
  isDragging.value = false
}

// 滚轮缩放
const onWheel = (e: WheelEvent) => {
  e.preventDefault()
  const delta = e.deltaY > 0 ? -0.1 : 0.1
  scale.value = Math.max(0.1, Math.min(5, scale.value + delta))
}

// 键盘事件
const handleKeydown = (e: KeyboardEvent) => {
  if (!showPreview.value) return

  switch (e.key) {
    case 'Escape':
      closePreview()
      break
    case '+':
    case '=':
      zoomIn()
      break
    case '-':
    case '_':
      zoomOut()
      break
    case '0':
      reset()
      break
  }
}

// 错误处理
const handleError = (event: Event) => {
  emit('error', event)
}

// 新窗口打开
const openInNewTab = () => {
  window.open(props.src, '_blank')
}

// 进入动画钩子
const onEnter = (el: Element) => {
  const element = el as HTMLElement
  element.style.transition = 'transform 0.3s cubic-bezier(0.4, 0, 0.2, 1), opacity 0.3s'
}

// 离开动画钩子
const onLeave = (el: Element) => {
  const element = el as HTMLElement
  element.style.transition = 'transform 0.3s cubic-bezier(0.4, 0, 0.2, 1), opacity 0.3s'
}

onMounted(() => {
  document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
})
</script>

<style scoped>
.fade-bg-enter-active,
.fade-bg-leave-active {
  transition: opacity 0.3s ease;
}

.fade-bg-enter-from,
.fade-bg-leave-to {
  opacity: 0;
}
</style>
