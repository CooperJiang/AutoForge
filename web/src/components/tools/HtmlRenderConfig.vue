<template>
  <div class="space-y-4">
    
    <div class="bg-bg-hover rounded-lg p-3 border border-border-primary">
      <VariableHelper
        v-if="showHelper"
        :show="true"
        :previous-nodes="props.previousNodes"
        :env-vars="props.envVars"
        @insert-field="handleInsertField"
        @insert-node="handleInsertNode"
        @insert-env="handleInsertEnv"
      />
    </div>

    
    <div>
      <label class="block text-sm font-medium text-text-primary mb-2">
        HTML 内容 <span class="text-error">*</span>
      </label>
      <textarea
        ref="contentInputRef"
        v-model="localConfig.content"
        @focus="activeInputRef = $event.target"
        rows="8"
        placeholder="输入 HTML 内容，支持变量如 {{nodes.xxx.content}}"
        class="w-full px-3 py-2 bg-bg-elevated border border-border-primary rounded-lg text-text-primary font-mono text-sm focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent resize-y"
      ></textarea>
      <p class="mt-1 text-xs text-text-tertiary">要保存的 HTML 内容，将自动生成可访问的预览 URL</p>
    </div>

    
    <div>
      <label class="block text-sm font-medium text-text-primary mb-2"> 页面标题 </label>
      <input
        ref="titleInputRef"
        v-model="localConfig.title"
        @focus="activeInputRef = $event.target"
        type="text"
        placeholder="可选，网页标题"
        class="w-full px-3 py-2 bg-bg-elevated border border-border-primary rounded-lg text-text-primary focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
      />
    </div>

    
    <div>
      <label class="block text-sm font-medium text-text-primary mb-2"> 过期时间（小时） </label>
      <input
        v-model.number="localConfig.expires_hours"
        type="number"
        min="0"
        placeholder="0 表示永不过期"
        class="w-full px-3 py-2 bg-bg-elevated border border-border-primary rounded-lg text-text-primary focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
      />
      <p class="mt-1 text-xs text-text-tertiary">生成的预览链接过期时间，0 表示永不过期</p>
    </div>

    
    <div class="bg-primary-light border border-primary rounded-lg p-3">
      <div class="text-xs text-primary space-y-1">
        <div class="font-medium">功能说明：</div>
        <div>• 将 HTML 内容保存为静态网页文件</div>
        <div>• 生成可访问的 URL 地址，支持分享和预览</div>
        <div>• 输出的 URL 可以被后续节点使用</div>
        <div>• 支持设置过期时间自动清理文件</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import VariableHelper from '@/components/VariableHelper'
import type { WorkflowNode, WorkflowEnvVar } from '@/types/workflow'
import { message } from '@/utils/message'

interface Props {
  config: Record<string, any>
  previousNodes?: WorkflowNode[]
  envVars?: WorkflowEnvVar[]
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:config': [config: Record<string, any>]
}>()

const showHelper = ref(true)
const activeInputRef = ref<HTMLTextAreaElement | HTMLInputElement | null>(null)
const contentInputRef = ref<HTMLTextAreaElement>()
const titleInputRef = ref<HTMLInputElement>()

// 本地配置
const localConfig = ref({
  content: props.config?.content || '',
  title: props.config?.title || '',
  expires_hours: props.config?.expires_hours || 0,
})

// 插入变量到聚焦的输入框
const insertVariable = (variable: string) => {
  const element = activeInputRef.value
  if (!element) {
    navigator.clipboard.writeText(variable).then(() => {
      message.success('变量已复制到剪贴板')
    })
    return
  }

  const start = element.selectionStart || 0
  const end = element.selectionEnd || 0
  const text = element.value || ''

  const newText = text.substring(0, start) + variable + text.substring(end)

  // 更新对应的字段
  if (element === contentInputRef.value) {
    localConfig.value.content = newText
  } else if (element === titleInputRef.value) {
    localConfig.value.title = newText
  }

  // 恢复光标位置
  setTimeout(() => {
    element.focus()
    element.setSelectionRange(start + variable.length, start + variable.length)
  }, 0)
}

// 处理字段变量插入
const handleInsertField = (field: string) => {
  insertVariable(`{{${field}}}`)
}

// 处理节点变量插入
const handleInsertNode = (path: string) => {
  insertVariable(`{{${path}}}`)
}

// 处理环境变量插入
const handleInsertEnv = (envVar: string) => {
  insertVariable(`{{env.${envVar}}}`)
}

// 防止循环更新
const updatingFromProps = ref(false)

// 监听本地配置变化
watch(
  localConfig,
  (newConfig) => {
    if (!updatingFromProps.value) {
      emit('update:config', { ...newConfig })
    }
  },
  { deep: true }
)

// 监听 props 配置变化
watch(
  () => props.config,
  (newConfig) => {
    if (newConfig) {
      updatingFromProps.value = true
      localConfig.value = {
        content: newConfig.content || '',
        title: newConfig.title || '',
        expires_hours: newConfig.expires_hours || 0,
      }
      setTimeout(() => {
        updatingFromProps.value = false
      }, 0)
    }
  },
  { deep: true }
)
</script>
