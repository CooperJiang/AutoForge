<template>
  <div class="output-viewer max-w-full overflow-auto">
    <!-- 有渲染配置时使用专门的子组件 -->
    <component
      v-if="hasRenderConfig"
      :is="viewerComponent"
      v-bind="viewerProps"
      @error="handleError"
    />

    <!-- 无渲染配置时使用默认 JSON 显示 -->
    <JsonViewer v-else :content="jsonContent" />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import JsonViewer from '@/components/JsonViewer'
import ImageViewer from './ImageViewer.vue'
import VideoViewer from './VideoViewer.vue'
import HtmlViewer from './HtmlViewer.vue'
import MarkdownViewer from './MarkdownViewer.vue'
import TextViewer from './TextViewer.vue'
import GalleryViewer from './GalleryViewer.vue'
import UrlViewer from './UrlViewer.vue'
import type { OutputRenderConfig } from '@/types/workflow'

interface Props {
  output: any
  outputRender?: OutputRenderConfig
}

const props = defineProps<Props>()

// 是否有渲染配置
const hasRenderConfig = computed(() => {
  return props.outputRender && props.outputRender.type !== 'json'
})

// 渲染类型
const renderType = computed(() => {
  return props.outputRender?.type || 'json'
})

// 主要内容
const primaryContent = computed(() => {
  if (!props.outputRender) return ''
  const primaryField = props.outputRender.primary
  return getNestedValue(props.output, primaryField) || ''
})

// 描述信息
const description = computed(() => props.output?.description || '')

// 替代文本
const altText = computed(() => props.output?.alt_text || props.output?.title || '')

// 缩略图
const thumbnail = computed(() => props.output?.thumbnail || '')

// 画廊项目
const galleryItems = computed(() => {
  const content = primaryContent.value
  if (typeof content === 'string') {
    try {
      const parsed = JSON.parse(content)
      return Array.isArray(parsed) ? parsed : []
    } catch {
      return []
    }
  }
  return Array.isArray(content) ? content : []
})

// JSON 内容
const jsonContent = computed(() => {
  return typeof props.output === 'string' ? props.output : JSON.stringify(props.output, null, 2)
})

// 动态选择渲染组件
const viewerComponent = computed(() => {
  const componentMap: Record<string, any> = {
    image: ImageViewer,
    video: VideoViewer,
    html: HtmlViewer,
    markdown: MarkdownViewer,
    text: TextViewer,
    gallery: GalleryViewer,
    url: UrlViewer,
  }
  return componentMap[renderType.value] || JsonViewer
})

// 动态生成组件 props
const viewerProps = computed(() => {
  const type = renderType.value

  switch (type) {
    case 'image':
      return {
        src: primaryContent.value,
        alt: altText.value,
        description: description.value,
      }

    case 'video':
      return {
        src: primaryContent.value,
        poster: thumbnail.value,
        description: description.value,
      }

    case 'html':
      return {
        content: primaryContent.value,
      }

    case 'markdown':
      return {
        content: primaryContent.value,
      }

    case 'text':
      return {
        content: primaryContent.value,
      }

    case 'gallery':
      return {
        items: galleryItems.value,
      }

    case 'url':
      return {
        src: primaryContent.value,
        alt: altText.value,
        description: description.value,
      }

    default:
      return {}
  }
})

// 获取嵌套对象的值
const getNestedValue = (obj: any, path: string): any => {
  if (!path) return obj

  const keys = path.split('.')
  let value = obj

  for (const key of keys) {
    // 处理数组索引，如 data[0]
    const arrayMatch = key.match(/^(.+)\[(\d+)\]$/)
    if (arrayMatch) {
      const [, arrayKey, index] = arrayMatch
      value = value?.[arrayKey]?.[parseInt(index)]
    } else {
      value = value?.[key]
    }

    if (value === undefined) break
  }

  return value
}

// 事件处理
const handleError = (event: Event) => {
  console.error('内容加载失败:', event)
}
</script>
