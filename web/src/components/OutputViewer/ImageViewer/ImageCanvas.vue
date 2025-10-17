<template>
  <div
    ref="containerRef"
    class="flex-1 relative overflow-hidden z-10"
    @mousedown="$emit('mousedown', $event)"
    @mousemove="$emit('mousemove', $event)"
    @mouseup="$emit('mouseup')"
    @mouseleave="$emit('mouseleave')"
    @wheel.prevent="$emit('wheel', $event)"
    @click="$emit('click')"
  >
    <Transition name="zoom" @enter="onEnter" @leave="onLeave">
      <img
        v-if="visible"
        ref="imageRef"
        :src="src"
        :alt="alt"
        class="absolute left-1/2 top-1/2 select-none"
        :style="imageStyle"
        @click.stop
        @dblclick.stop="$emit('dblclick')"
        @dragstart.prevent
      />
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, type CSSProperties } from 'vue'

interface Props {
  src: string
  alt?: string
  visible: boolean
  imageStyle: CSSProperties
}

defineProps<Props>()

defineEmits<{
  mousedown: [e: MouseEvent]
  mousemove: [e: MouseEvent]
  mouseup: []
  mouseleave: []
  wheel: [e: WheelEvent]
  click: []
  dblclick: []
  ready: [container: HTMLElement, image: HTMLImageElement]
}>()

const containerRef = ref<HTMLElement | null>(null)
const imageRef = ref<HTMLImageElement | null>(null)

// 动画钩子
const onEnter = (el: Element) => {
  const element = el as HTMLElement
  element.style.opacity = '0'
  element.style.transform = 'scale(0.8)'

  requestAnimationFrame(() => {
    element.style.transition = 'opacity 0.3s ease, transform 0.3s ease'
    element.style.opacity = '1'
    element.style.transform = 'scale(1)'
  })
}

const onLeave = (el: Element) => {
  const element = el as HTMLElement
  element.style.transition = 'opacity 0.2s ease, transform 0.2s ease'
  element.style.opacity = '0'
  element.style.transform = 'scale(0.9)'
}

// 暴露 refs
defineExpose({
  containerRef,
  imageRef,
})
</script>

<style scoped>
.zoom-enter-active,
.zoom-leave-active {
  transition:
    opacity 0.3s ease,
    transform 0.3s ease;
}

.zoom-enter-from {
  opacity: 0;
  transform: scale(0.8);
}

.zoom-leave-to {
  opacity: 0;
  transform: scale(0.9);
}
</style>
