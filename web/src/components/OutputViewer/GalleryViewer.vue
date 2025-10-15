<template>
  <div class="gallery-viewer">
    <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
      <div v-for="(item, index) in items" :key="index" class="gallery-item group relative">
        <img
          :src="item"
          :alt="`图片 ${index + 1}`"
          class="w-full h-48 object-cover rounded-lg border border-border-primary cursor-pointer transition-transform group-hover:scale-105"
          @click="handleClick(item, index)"
          @error="handleError(index, $event)"
        />
        <div
          class="absolute top-2 right-2 bg-bg-elevated/90 backdrop-blur px-2 py-1 rounded text-xs text-text-secondary"
        >
          {{ index + 1 }} / {{ items.length }}
        </div>
      </div>
    </div>

    <div v-if="items.length === 0" class="text-center py-8 text-text-tertiary">暂无图片</div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  items: string[]
}

defineProps<Props>()

const emit = defineEmits<{
  click: [url: string, index: number]
  error: [index: number, event: Event]
}>()

const handleClick = (url: string, index: number) => {
  emit('click', url, index)
}

const handleError = (index: number, event: Event) => {
  emit('error', index, event)
}
</script>
