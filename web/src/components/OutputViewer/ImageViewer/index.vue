<template>
  <div class="image-viewer">
    <!-- 缩略图 -->
    <ImageThumbnail
      :src="src"
      :alt="alt"
      :description="description"
      :show-centered-description="showCenteredDescription"
      @click="openPreview"
      @error="handleError"
    />

    <!-- 预览弹窗 -->
    <ImagePreviewModal
      ref="previewModalRef"
      :visible="showPreview"
      :src="src"
      :alt="alt"
      :scale="scale"
      :image-style="imageStyle"
      @zoom-in="zoomIn"
      @zoom-out="zoomOut"
      @rotate-left="rotateLeft"
      @rotate-right="rotateRight"
      @reset="reset"
      @open-new-tab="openInNewTab"
      @close="closePreview"
      @mousedown="startDrag"
      @mousemove="onDrag"
      @mouseup="stopDrag"
      @mouseleave="stopDrag"
      @wheel="onWheel"
      @canvas-click="closePreview"
      @dblclick="reset"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, nextTick } from 'vue'
import ImageThumbnail from './ImageThumbnail.vue'
import ImagePreviewModal from './ImagePreviewModal.vue'
import { useImageTransform } from './useImageTransform'
import { useImageInteraction } from './useImageInteraction'
import type { ImageViewerProps } from './types'

const props = withDefaults(defineProps<ImageViewerProps>(), {
  alt: '',
  description: '',
  showCenteredDescription: false,
})

// Refs
const previewModalRef = ref<InstanceType<typeof ImagePreviewModal> | null>(null)
const showPreview = ref(false)

// 动态获取 refs
const containerRef = computed(() => previewModalRef.value?.containerRef || null)
const imageRef = computed(() => previewModalRef.value?.imageRef || null)

// 图片变换逻辑
const transform = useImageTransform({
  containerRef: containerRef as any,
  imageRef: imageRef as any,
})

const {
  scale,
  rotation,
  translateX,
  translateY,
  imageStyle,
  zoomIn,
  zoomOut,
  rotateLeft,
  rotateRight,
  reset: resetTransform,
  centerImage,
  initialize,
} = transform

// 关闭预览
const closePreview = () => {
  showPreview.value = false
}

// 在新标签页打开
const openInNewTab = () => {
  window.open(props.src, '_blank')
}

// 错误处理
const handleError = (event: Event) => {
  console.error('图片加载失败:', event)
}

// 重置（双击）
const reset = () => {
  resetTransform()
}

// 图片交互逻辑（需要在 closePreview 定义后）
const { startDrag, onDrag, stopDrag, onWheel } = useImageInteraction({
  translateX,
  translateY,
  onZoomIn: zoomIn,
  onZoomOut: zoomOut,
  onRotateLeft: rotateLeft,
  onRotateRight: rotateRight,
  onReset: resetTransform,
  onClose: closePreview,
})

// 打开预览
const openPreview = async () => {
  showPreview.value = true

  // 等待 DOM 更新后初始化图片
  await nextTick()
  await nextTick() // 双重 nextTick 确保 refs 已经更新

  if (containerRef.value && imageRef.value) {
    initialize()
  }
}
</script>

<style scoped>
.image-viewer {
  display: inline-block;
}
</style>

