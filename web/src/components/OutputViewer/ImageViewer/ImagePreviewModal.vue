<template>
  <Teleport to="body">
    <Transition name="fade-bg">
      <div v-if="visible" class="fixed inset-0 z-50 bg-black/90 backdrop-blur-sm flex flex-col">
        <!-- 工具栏 -->
        <ImageToolbar
          :scale="scale"
          @zoom-in="$emit('zoom-in')"
          @zoom-out="$emit('zoom-out')"
          @rotate-left="$emit('rotate-left')"
          @rotate-right="$emit('rotate-right')"
          @reset="$emit('reset')"
          @open-new-tab="$emit('open-new-tab')"
          @close="$emit('close')"
        />

        <!-- 画布 -->
        <ImageCanvas
          ref="canvasRef"
          :src="src"
          :alt="alt"
          :visible="visible"
          :image-style="imageStyle"
          @mousedown="$emit('mousedown', $event)"
          @mousemove="$emit('mousemove', $event)"
          @mouseup="$emit('mouseup')"
          @mouseleave="$emit('mouseleave')"
          @wheel="$emit('wheel', $event)"
          @click="$emit('canvas-click')"
          @dblclick="$emit('dblclick')"
        />

        <!-- 底部提示 -->
        <div
          class="absolute bottom-6 left-1/2 -translate-x-1/2 bg-black/50 text-white text-xs px-4 py-2 rounded-lg pointer-events-none z-20"
        >
          拖动图片移动 | 滚轮缩放 | 双击重置
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, type CSSProperties } from 'vue'
import ImageToolbar from './ImageToolbar.vue'
import ImageCanvas from './ImageCanvas.vue'

interface Props {
  visible: boolean
  src: string
  alt?: string
  scale: number
  imageStyle: CSSProperties
}

defineProps<Props>()

defineEmits<{
  'zoom-in': []
  'zoom-out': []
  'rotate-left': []
  'rotate-right': []
  reset: []
  'open-new-tab': []
  close: []
  mousedown: [e: MouseEvent]
  mousemove: [e: MouseEvent]
  mouseup: []
  mouseleave: []
  wheel: [e: WheelEvent]
  'canvas-click': []
  dblclick: []
}>()

const canvasRef = ref<InstanceType<typeof ImageCanvas> | null>(null)

// 暴露画布的 refs
defineExpose({
  get containerRef() {
    return canvasRef.value?.containerRef || null
  },
  get imageRef() {
    return canvasRef.value?.imageRef || null
  },
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
