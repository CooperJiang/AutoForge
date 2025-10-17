<template>
  <div ref="containerRef" class="relative group">
    <!-- Images Container -->
    <div class="overflow-hidden rounded-lg">
      <div
        class="flex transition-transform duration-300 ease-in-out"
        :style="{ transform: `translateX(-${currentIndex * (imageWidth + gap)}px)` }"
      >
        <div
          v-for="(imageUrl, index) in images"
          :key="index"
          class="flex-shrink-0"
          :style="{ width: `${imageWidth}px`, marginRight: index < images.length - 1 ? `${gap}px` : '0' }"
        >
          <slot name="image" :src="imageUrl" :index="index">
            <ImageLoader
              :src="imageUrl"
              :alt="`案例 ${index + 1}`"
              :height="height"
              :clickable="clickable"
              @click="handleImageClick(imageUrl, index)"
            />
          </slot>
        </div>
      </div>
    </div>

    <!-- Left Arrow -->
    <button
      v-if="showControls && canGoPrev"
      class="absolute left-2 top-1/2 -translate-y-1/2 w-8 h-8 rounded-full bg-black/50 hover:bg-black/70 text-white flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity z-10"
      @click="prev"
    >
      <ChevronLeft class="w-5 h-5" />
    </button>

    <!-- Right Arrow -->
    <button
      v-if="showControls && canGoNext"
      class="absolute right-2 top-1/2 -translate-y-1/2 w-8 h-8 rounded-full bg-black/50 hover:bg-black/70 text-white flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity z-10"
      @click="next"
    >
      <ChevronRight class="w-5 h-5" />
    </button>

    <!-- Indicators -->
    <div
      v-if="showIndicators && images.length > visibleCount"
      class="absolute bottom-3 left-1/2 -translate-x-1/2 flex gap-1.5 z-10"
    >
      <button
        v-for="(_, index) in pageCount"
        :key="index"
        class="w-1.5 h-1.5 rounded-full transition-all"
        :class="currentPage === index ? 'bg-white w-4' : 'bg-white/50 hover:bg-white/75'"
        @click="goToPage(index)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { ChevronLeft, ChevronRight } from 'lucide-vue-next'
import ImageLoader from '@/components/ImageLoader'

interface Props {
  images: string[]
  imageWidth?: number
  height?: string
  gap?: number
  clickable?: boolean
  showIndicators?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  imageWidth: 320,
  height: '200px',
  gap: 16,
  clickable: true,
  showIndicators: true,
})

const emit = defineEmits<{
  imageClick: [url: string, index: number]
}>()

const containerRef = ref<HTMLElement>()
const currentIndex = ref(0)
const containerWidth = ref(0)

// 根据容器宽度计算可见数量
const visibleCount = computed(() => {
  if (containerWidth.value === 0) return 3 // 默认值
  const availableWidth = containerWidth.value
  const itemWidth = props.imageWidth + props.gap
  return Math.max(1, Math.floor(availableWidth / itemWidth))
})

// 计算总页数
const pageCount = computed(() => {
  return Math.ceil(props.images.length / visibleCount.value)
})

// 当前页码
const currentPage = computed(() => {
  return Math.floor(currentIndex.value / visibleCount.value)
})

// 计算所有图片的总宽度
const totalImagesWidth = computed(() => {
  return props.images.length * props.imageWidth + (props.images.length - 1) * props.gap
})

// 是否显示控制按钮（图片总宽度超过容器宽度时显示）
const showControls = computed(() => {
  return totalImagesWidth.value > containerWidth.value
})

// 能否向前
const canGoPrev = computed(() => {
  return currentIndex.value > 0
})

// 能否向后
const canGoNext = computed(() => {
  return currentIndex.value < props.images.length - visibleCount.value
})

const prev = () => {
  if (canGoPrev.value) {
    currentIndex.value = Math.max(0, currentIndex.value - visibleCount.value)
  }
}

const next = () => {
  if (canGoNext.value) {
    currentIndex.value = Math.min(
      props.images.length - visibleCount.value,
      currentIndex.value + visibleCount.value
    )
  }
}

const goToPage = (pageIndex: number) => {
  currentIndex.value = pageIndex * visibleCount.value
}

const handleImageClick = (url: string, index: number) => {
  emit('imageClick', url, index)
}

// 更新容器宽度
const updateContainerWidth = () => {
  if (containerRef.value) {
    containerWidth.value = containerRef.value.clientWidth
  }
}

// 监听窗口大小变化
let resizeObserver: ResizeObserver | null = null

onMounted(() => {
  updateContainerWidth()

  // 使用 ResizeObserver 监听容器大小变化
  if (containerRef.value) {
    resizeObserver = new ResizeObserver(() => {
      updateContainerWidth()
    })
    resizeObserver.observe(containerRef.value)
  }

  // 兼容旧浏览器，同时监听 window resize
  window.addEventListener('resize', updateContainerWidth)
})

onUnmounted(() => {
  if (resizeObserver && containerRef.value) {
    resizeObserver.unobserve(containerRef.value)
  }
  window.removeEventListener('resize', updateContainerWidth)
})

// 重置索引当图片列表变化时
watch(() => props.images, () => {
  currentIndex.value = 0
})
</script>
