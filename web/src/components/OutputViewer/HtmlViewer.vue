<template>
  <div class="html-viewer-container">
    <div class="flex items-center justify-between mb-2 gap-2">
      <div class="text-xs text-text-secondary">HTML 预览</div>
      <button
        @click="openInNewWindow"
        class="px-3 py-1.5 text-xs bg-bg-elevated hover:bg-bg-hover border border-border-primary rounded flex items-center gap-1.5 text-text-secondary hover:text-text-primary transition-colors"
      >
        <Maximize2 class="w-3.5 h-3.5" />
        全屏预览
      </button>
    </div>

    <iframe
      ref="iframeRef"
      :srcdoc="sanitizedHtml"
      class="w-full border border-border-primary rounded-lg bg-white"
      style="min-height: 400px; max-height: 800px"
      sandbox="allow-same-origin"
      @load="adjustIframeHeight"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { Maximize2 } from 'lucide-vue-next'

interface Props {
  content: string
}

const props = defineProps<Props>()
const iframeRef = ref<HTMLIFrameElement>()

const sanitizedHtml = computed(() => {
  let html = props.content

  // 移除 markdown 代码块标记（```html 或 ``` 开头/结尾）
  html = html.replace(/^```html\s*\n/i, '').replace(/^```\s*\n/i, '')
  html = html.replace(/\n```\s*$/i, '')

  // 移除潜在危险的标签和属性
  html = html.replace(/<script\b[^<]*(?:(?!<\/script>)<[^<]*)*<\/script>/gi, '')
  html = html.replace(/on\w+\s*=\s*["'][^"']*["']/gi, '') // 移除内联事件处理器

  return html.trim()
})

// 调整 iframe 高度以适应内容
const adjustIframeHeight = () => {
  if (!iframeRef.value) return

  try {
    const iframe = iframeRef.value
    const iframeDoc = iframe.contentDocument || iframe.contentWindow?.document

    if (iframeDoc) {
      setTimeout(() => {
        const height = iframeDoc.documentElement.scrollHeight
        iframe.style.height = `${Math.min(Math.max(height, 400), 800)}px`
      }, 100)
    }
  } catch {
    // Cross-origin iframe cannot adjust height
  }
}

// 全屏预览 - 新窗口打开
const openInNewWindow = () => {
  const newWindow = window.open('', '_blank', 'width=1200,height=800')
  if (newWindow) {
    newWindow.document.write(sanitizedHtml.value)
    newWindow.document.close()
  }
}
</script>

<style scoped>
.html-viewer-container {
  max-width: 100%;
  overflow: hidden;
}
</style>
