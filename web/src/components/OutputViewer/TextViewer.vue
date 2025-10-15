<template>
  <pre
    class="text-viewer bg-bg-hover border border-border-primary p-4 rounded-lg font-mono text-sm overflow-auto whitespace-pre-wrap"
    style="max-height: 600px"
    >{{ content }}</pre
  >
</template>

<script setup lang="ts">
interface Props {
  content: string
}

defineProps<Props>()
</script>
