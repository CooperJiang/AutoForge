<template>
  <div
    class="relative overflow-hidden rounded-lg"
    :style="{ width, height, maxWidth, maxHeight }"
  >
    <!-- Loading State -->
    <div
      v-if="loading"
      class="absolute inset-0 flex items-center justify-center bg-surface-tertiary"
    >
      <div class="flex flex-col items-center gap-2">
        <div class="w-8 h-8 border-3 border-primary border-t-transparent rounded-full animate-spin"></div>
        <span class="text-xs text-text-tertiary">加载中...</span>
      </div>
    </div>

    <!-- Error State -->
    <div
      v-else-if="error"
      class="absolute inset-0 flex items-center justify-center bg-surface-tertiary"
    >
      <div class="flex flex-col items-center gap-2 text-text-tertiary">
        <svg
          class="w-8 h-8"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
          />
        </svg>
        <span class="text-xs">加载失败</span>
      </div>
    </div>

    <!-- Image -->
    <img
      v-show="!loading && !error"
      :src="src"
      :alt="alt"
      :class="[
        'w-full h-full transition-opacity duration-300',
        objectFit === 'cover' ? 'object-cover' : 'object-contain',
        clickable ? 'cursor-pointer hover:opacity-90' : 'pointer-events-none',
      ]"
      @load="handleLoad"
      @error="handleError"
      @click="handleClick"
    />

    <!-- Hover Overlay (optional) -->
    <div
      v-if="clickable && !loading && !error"
      class="absolute inset-0 bg-black/50 opacity-0 hover:opacity-100 transition-opacity flex items-center justify-center"
    >
      <slot name="overlay">
        <span class="text-white text-sm font-medium">点击查看大图</span>
      </slot>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

interface Props {
  src: string
  alt?: string
  width?: string
  height?: string
  maxWidth?: string
  maxHeight?: string
  objectFit?: 'cover' | 'contain'
  clickable?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  alt: '',
  width: '100%',
  height: '100%',
  maxWidth: undefined,
  maxHeight: undefined,
  objectFit: 'cover',
  clickable: false,
})

const emit = defineEmits<{
  click: [event: MouseEvent]
  load: []
  error: []
}>()

const loading = ref(true)
const error = ref(false)

const handleLoad = () => {
  loading.value = false
  error.value = false
  emit('load')
}

const handleError = () => {
  loading.value = false
  error.value = true
  emit('error')
}

const handleClick = (event: MouseEvent) => {
  if (props.clickable) {
    emit('click', event)
  }
}
</script>
