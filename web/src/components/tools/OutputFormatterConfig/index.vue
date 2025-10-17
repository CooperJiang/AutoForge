<template>
  <div class="space-y-4">
    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2">
        输出类型 <span class="text-red-500">*</span>
      </label>
      <BaseSelect v-model="localConfig.output_type" :options="outputTypeOptions" />
      <p class="text-xs text-text-tertiary mt-1">选择最终输出的展示类型</p>
    </div>

    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2">
        主要内容 <span class="text-red-500">*</span>
      </label>
      <textarea
        v-model="localConfig.content"
        rows="3"
        :placeholder="getContentPlaceholder()"
        class="w-full px-3 py-2 border-2 border-border-primary rounded-lg focus:outline-none focus:border-primary bg-bg-elevated text-text-primary font-mono text-sm"
      />
      <p class="text-xs text-text-tertiary mt-1">
        {{ getContentDescription() }}
      </p>
    </div>

    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2"> 标题（可选） </label>
      <input
        v-model="localConfig.title"
        type="text"
        placeholder="如：生成的图片"
        class="w-full px-3 py-2 border-2 border-border-primary rounded-lg focus:outline-none focus:border-primary bg-bg-elevated text-text-primary font-mono text-sm"
      />
    </div>

    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2"> 描述（可选） </label>
      <textarea
        v-model="localConfig.description"
        rows="2"
        placeholder="如：{{nodes.xxx.response.data[0].revised_prompt}}"
        class="w-full px-3 py-2 border-2 border-border-primary rounded-lg focus:outline-none focus:border-primary bg-bg-elevated text-text-primary font-mono text-sm"
      />
    </div>

    <div v-if="localConfig.output_type === 'image' || localConfig.output_type === 'video'">
      <label class="block text-sm font-medium text-text-secondary mb-2"> 替代文本（可选） </label>
      <input
        v-model="localConfig.alt_text"
        type="text"
        placeholder="图片的替代描述"
        class="w-full px-3 py-2 border-2 border-border-primary rounded-lg focus:outline-none focus:border-primary bg-bg-elevated text-text-primary font-mono text-sm"
      />
    </div>

    <div v-if="localConfig.output_type === 'video'">
      <label class="block text-sm font-medium text-text-secondary mb-2"> 缩略图 URL（可选） </label>
      <input
        v-model="localConfig.thumbnail"
        type="text"
        placeholder="视频封面图片 URL"
        class="w-full px-3 py-2 border-2 border-border-primary rounded-lg focus:outline-none focus:border-primary bg-bg-elevated text-text-primary font-mono text-sm"
      />
    </div>

    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2"> 元数据（可选） </label>
      <textarea
        v-model="localConfig.metadata"
        rows="2"
        placeholder='JSON 格式的元数据，如 {"model": "dall-e-3"}'
        class="w-full px-3 py-2 border-2 border-border-primary rounded-lg focus:outline-none focus:border-primary bg-bg-elevated text-text-primary font-mono text-sm"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import BaseSelect from '@/components/BaseSelect'
import type { WorkflowNode, WorkflowEnvVar } from '@/types/workflow'

interface Props {
  config: Record<string, any>
  previousNodes?: WorkflowNode[]
  envVars?: WorkflowEnvVar[]
}

const props = withDefaults(defineProps<Props>(), {
  previousNodes: () => [],
  envVars: () => [],
})

const emit = defineEmits<{
  'update:config': [config: Record<string, any>]
}>()

// 本地配置
const localConfig = ref({
  output_type: props.config.output_type || 'json',
  content: props.config.content || '',
  title: props.config.title || '',
  description: props.config.description || '',
  alt_text: props.config.alt_text || '',
  thumbnail: props.config.thumbnail || '',
  metadata: props.config.metadata || '',
})

// 输出类型选项
const outputTypeOptions = [
  { label: 'JSON 数据', value: 'json' },
  { label: '图片', value: 'image' },
  { label: '视频', value: 'video' },
  { label: 'HTML', value: 'html' },
  { label: 'HTML URL 预览', value: 'html-url' },
  { label: 'Markdown', value: 'markdown' },
  { label: '纯文本', value: 'text' },
  { label: '图片画廊', value: 'gallery' },
]

// 根据类型返回内容占位符
const getContentPlaceholder = () => {
  const type = localConfig.value.output_type
  const placeholders: Record<string, string> = {
    image: '图片 URL，如 {{nodes.xxx.response.data[0].url}}',
    video: '视频 URL',
    html: 'HTML 字符串',
    'html-url': 'HTML 页面的 URL 地址，如 {{nodes.xxx.url}}',
    markdown: 'Markdown 文本',
    text: '纯文本内容',
    gallery: '图片 URL 数组（JSON 字符串）',
    json: 'JSON 数据',
  }
  return placeholders[type] || '请输入内容'
}

// 根据类型返回内容说明
const getContentDescription = () => {
  const type = localConfig.value.output_type
  const descriptions: Record<string, string> = {
    image: '填写图片的 URL 地址，支持从上一个节点引用',
    video: '填写视频的 URL 地址',
    html: '填写要渲染的 HTML 内容（会自动移除 script 标签）',
    'html-url': '填写 HTML 页面的 URL 地址，将使用 iframe 预览并提供分享链接',
    markdown: '填写 Markdown 格式的文本',
    text: '填写纯文本内容',
    gallery: '填写图片 URL 的数组，可以是 JSON 字符串或直接引用数组变量',
    json: '填写 JSON 格式的数据',
  }
  return descriptions[type] || ''
}

// 标记是否正在从外部更新
const updatingFromProps = ref(false)

// 监听配置变化
watch(
  localConfig,
  (newConfig) => {
    if (!updatingFromProps.value) {
      emit('update:config', { ...newConfig })
    }
  },
  { deep: true }
)

// 监听外部配置变化
watch(
  () => props.config,
  (newConfig) => {
    updatingFromProps.value = true
    localConfig.value = {
      output_type: newConfig.output_type || 'json',
      content: newConfig.content || '',
      title: newConfig.title || '',
      description: newConfig.description || '',
      alt_text: newConfig.alt_text || '',
      thumbnail: newConfig.thumbnail || '',
      metadata: newConfig.metadata || '',
    }
    setTimeout(() => {
      updatingFromProps.value = false
    }, 0)
  },
  { immediate: true, deep: true }
)
</script>
