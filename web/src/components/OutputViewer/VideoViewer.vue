<template>
  <div class="video-viewer">
    <video
      :src="src"
      :poster="poster"
      controls
      class="max-w-full rounded-lg border border-border-primary shadow-sm"
      @error="handleError"
    ></video>
    <div v-if="description" class="mt-2 text-sm text-text-secondary">
      {{ description }}
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  src: string
  poster?: string
  description?: string
}

defineProps<Props>()

const emit = defineEmits<{
  error: [event: Event]
}>()

const handleError = (event: Event) => {
  emit('error', event)
}
</script>
