<template>
  <div
    class="flex-shrink-0 h-16 flex items-center justify-between px-6 bg-black/50 z-20"
    @click.stop
  >
    <div class="flex items-center gap-4">
      <div class="flex items-center gap-2 bg-white/10 rounded-lg px-3 py-2">
        <button
          class="p-1 hover:bg-white/20 rounded transition-colors"
          @click.stop="$emit('zoom-out')"
          title="缩小 (滚轮向下)"
        >
          <ZoomOut class="w-5 h-5 text-white" />
        </button>
        <span class="text-white text-sm min-w-[60px] text-center">
          {{ Math.round(scale * 100) }}%
        </span>
        <button
          class="p-1 hover:bg-white/20 rounded transition-colors"
          @click.stop="$emit('zoom-in')"
          title="放大 (滚轮向上)"
        >
          <ZoomIn class="w-5 h-5 text-white" />
        </button>
      </div>

      <button
        class="p-2 bg-white/10 hover:bg-white/20 rounded-lg transition-colors"
        @click.stop="$emit('rotate-left')"
        title="逆时针旋转"
      >
        <RotateCcw class="w-5 h-5 text-white" />
      </button>
      <button
        class="p-2 bg-white/10 hover:bg-white/20 rounded-lg transition-colors"
        @click.stop="$emit('rotate-right')"
        title="顺时针旋转"
      >
        <RotateCw class="w-5 h-5 text-white" />
      </button>

      <button
        class="p-2 bg-white/10 hover:bg-white/20 rounded-lg transition-colors"
        @click.stop="$emit('reset')"
        title="重置 (双击图片)"
      >
        <RefreshCw class="w-5 h-5 text-white" />
      </button>
    </div>

    <div class="flex items-center gap-2">
      <button
        class="px-4 py-2 bg-white/10 hover:bg-white/20 rounded-lg transition-colors text-white text-sm"
        @click.stop="$emit('open-new-tab')"
        title="在新标签页打开"
      >
        新窗口打开
      </button>

      <button
        class="p-2 hover:bg-white/20 rounded-lg transition-colors"
        @click.stop="$emit('close')"
        title="关闭 (ESC)"
      >
        <X class="w-6 h-6 text-white" />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ZoomIn, ZoomOut, RotateCcw, RotateCw, RefreshCw, X } from 'lucide-vue-next'
import type { ImageToolbarProps, ImageToolbarEmits } from '@/types/image-viewer'

defineProps<ImageToolbarProps>()

defineEmits<ImageToolbarEmits>()
</script>
