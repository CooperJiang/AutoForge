<template>
  <div class="image-thumbnail">
    <div
      class="relative inline-block cursor-pointer"
      @mousedown="handleMouseDown"
      @mouseup="handleMouseUp"
      @mousemove="handleMouseMove"
    >
      <ImageLoader
        ref="thumbnailRef"
        :src="src"
        :alt="alt"
        width="100%"
        :height="maxHeight"
        object-fit="contain"
        :clickable="false"
        class="border border-border-primary shadow-sm rounded-lg"
        @error="$emit('error', $event)"
      >
        <template #overlay>
          <ZoomIn class="w-8 h-8 text-white" />
        </template>
      </ImageLoader>

      <div
        v-if="showCenteredDescription && alt"
        class="absolute inset-0 flex items-center justify-center pointer-events-none"
      >
        <div class="bg-black/60 backdrop-blur-sm px-4 py-2 rounded-lg">
          <span class="text-white text-sm font-medium">{{ alt }}</span>
        </div>
      </div>
    </div>

    <div v-if="description && !showCenteredDescription" class="mt-2 text-sm text-text-secondary">
      {{ description }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ZoomIn } from 'lucide-vue-next'
import ImageLoader from '@/components/ImageLoader'
import { useClickOrDrag } from '@/composables/useClickOrDrag'
import { THUMBNAIL_MAX_HEIGHT } from '@/constants/image-viewer'
import type { ImageThumbnailProps, ImageThumbnailEmits } from '@/types/image-viewer'

defineProps<ImageThumbnailProps>()

const emit = defineEmits<ImageThumbnailEmits>()

const { handleMouseDown, handleMouseMove, handleMouseUp, onClick } = useClickOrDrag()

onClick(() => {
  emit('click')
})

const maxHeight = `${THUMBNAIL_MAX_HEIGHT}px`
</script>
