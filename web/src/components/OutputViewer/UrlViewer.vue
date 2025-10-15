<template>
  <div class="url-viewer-container">
    <!-- 工具栏 -->
    <div class="flex items-center justify-between mb-2 gap-2">
      <div class="text-xs text-text-secondary">URL 预览</div>
      <div class="flex items-center gap-2">
        <button
          @click="copyUrl"
          class="px-3 py-1.5 text-xs bg-bg-elevated hover:bg-bg-hover border border-border-primary rounded flex items-center gap-1.5 transition-colors"
          :class="copied ? 'text-green-600' : 'text-text-secondary hover:text-text-primary'"
        >
          <component :is="copied ? Check : Copy" class="w-3.5 h-3.5" />
          {{ copied ? '已复制' : '复制分享链接' }}
        </button>
        <button
          @click="openInNewWindow"
          class="px-3 py-1.5 text-xs bg-bg-elevated hover:bg-bg-hover border border-border-primary rounded flex items-center gap-1.5 text-text-secondary hover:text-text-primary transition-colors"
        >
          <ExternalLink class="w-3.5 h-3.5" />
          新窗口预览
        </button>
      </div>
    </div>

    <!-- iframe 预览 -->
    <iframe
      ref="iframeRef"
      :src="fullUrl"
      class="w-full border border-border-primary rounded-lg bg-white"
      style="min-height: 400px; max-height: 800px"
      sandbox="allow-same-origin allow-scripts allow-forms allow-popups"
      @load="adjustIframeHeight"
    />

    <!-- 分享链接显示 -->
    <div
      class="mt-2 flex items-center gap-2 p-2 bg-bg-hover rounded border border-border-primary overflow-hidden"
    >
      <span class="text-xs text-text-tertiary flex-shrink-0">分享链接:</span>
      <code class="flex-1 text-xs font-mono text-text-secondary truncate min-w-0">{{
        fullUrl
      }}</code>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { Copy, Check, ExternalLink } from 'lucide-vue-next'
import { message } from '@/utils/message'

interface Props {
  src: string
}

const props = defineProps<Props>()
const iframeRef = ref<HTMLIFrameElement>()
const copied = ref(false)

const fullUrl = computed(() => {
  const url = props.src

  if (url.startsWith('/')) {
    if (import.meta.env.DEV) {
      return `http://localhost:7777${url}`
    }
    return `${window.location.origin}${url}`
  }

  return url
})

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

const openInNewWindow = () => {
  window.open(fullUrl.value, '_blank')
}

const copyUrl = async () => {
  try {
    await navigator.clipboard.writeText(fullUrl.value)
    copied.value = true
    message.success('Link copied to clipboard')
    setTimeout(() => {
      copied.value = false
    }, 2000)
  } catch {
    message.error('Failed to copy')
  }
}
</script>

<style scoped>
.url-viewer-container {
  max-width: 100%;
  overflow: hidden;
}
</style>
