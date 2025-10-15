<template>
  <div
    class="markdown-viewer prose dark:prose-invert max-w-none bg-bg-hover border border-border-primary p-4 rounded-lg overflow-auto"
    style="max-height: 600px"
    v-html="renderedHtml"
  ></div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  content: string
}

const props = defineProps<Props>()

const renderedHtml = computed(() => {
  let html = props.content

  // 标题
  html = html.replace(/^### (.*$)/gim, '<h3 class="text-lg font-semibold mb-2">$1</h3>')
  html = html.replace(/^## (.*$)/gim, '<h2 class="text-xl font-bold mb-3">$1</h2>')
  html = html.replace(/^# (.*$)/gim, '<h1 class="text-2xl font-bold mb-4">$1</h1>')

  // 粗体和斜体
  html = html.replace(/\*\*\*(.+?)\*\*\*/g, '<strong><em>$1</em></strong>')
  html = html.replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>')
  html = html.replace(/\*(.+?)\*/g, '<em>$1</em>')

  // 代码块
  html = html.replace(
    /```(\w+)?\n([\s\S]+?)```/g,
    '<pre class="bg-bg-primary p-3 rounded border border-border-primary overflow-x-auto"><code>$2</code></pre>'
  )
  html = html.replace(
    /`([^`]+)`/g,
    '<code class="bg-bg-primary px-1.5 py-0.5 rounded text-sm">$1</code>'
  )

  // 链接
  html = html.replace(
    /\[([^\]]+)\]\(([^\)]+)\)/g,
    '<a href="$2" class="text-primary hover:underline" target="_blank">$1</a>'
  )

  // 列表
  html = html.replace(/^\* (.+)$/gim, '<li class="ml-4">$1</li>')
  html = html.replace(/(<li.*<\/li>\n?)+/g, '<ul class="list-disc space-y-1 my-2">$&</ul>')

  // 换行
  html = html.replace(/\n\n/g, '</p><p class="mb-2">')
  html = html.replace(/\n/g, '<br>')

  return `<div>${html}</div>`
})
</script>

<style scoped>
.prose :deep(h1),
.prose :deep(h2),
.prose :deep(h3) {
  color: var(--color-text-primary);
}

.prose :deep(code) {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}
</style>
