<template>
  <div class="relative inline-block" @mouseenter="show = true" @mouseleave="show = false">
    <slot />
    <Transition name="tooltip">
      <div
        v-if="show && text"
        :class="[
          'absolute z-50 px-2 py-1 text-xs font-medium text-white bg-slate-900 rounded shadow-lg whitespace-nowrap pointer-events-none',
          positionClasses[position]
        ]"
      >
        {{ text }}
        <div :class="['absolute w-2 h-2 bg-slate-900 transform rotate-45', arrowClasses[position]]"></div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

interface Props {
  text: string
  position?: 'top' | 'bottom' | 'left' | 'right'
}

const props = withDefaults(defineProps<Props>(), {
  position: 'bottom'
})

const show = ref(false)

const positionClasses = {
  top: 'bottom-full left-1/2 -translate-x-1/2 mb-2',
  bottom: 'top-full left-1/2 -translate-x-1/2 mt-2',
  left: 'right-full top-1/2 -translate-y-1/2 mr-2',
  right: 'left-full top-1/2 -translate-y-1/2 ml-2'
}

const arrowClasses = {
  top: 'top-full left-1/2 -translate-x-1/2 -mt-1',
  bottom: 'bottom-full left-1/2 -translate-x-1/2 -mb-1',
  left: 'left-full top-1/2 -translate-y-1/2 -ml-1',
  right: 'right-full top-1/2 -translate-y-1/2 -mr-1'
}
</script>

<style scoped>
.tooltip-enter-active,
.tooltip-leave-active {
  transition: opacity 0.15s ease;
}

.tooltip-enter-from,
.tooltip-leave-to {
  opacity: 0;
}
</style>
