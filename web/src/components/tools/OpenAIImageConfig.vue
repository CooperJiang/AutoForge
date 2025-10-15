<template>
  <div class="space-y-4">
    <!-- 提示信息 -->
    <div class="bg-primary-light border-l-4 border-primary p-3 rounded">
      <p class="text-sm text-text-primary">💡 API 凭证已在系统配置中统一管理，无需每次填写</p>
    </div>

    <!-- 模型选择 -->
    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2">
        模型 <span class="text-red-500">*</span>
      </label>
      <BaseInput v-model="localConfig.model" placeholder="例如：dall-e-3, dall-e-2, gpt-image-1" />
      <p class="mt-1 text-xs text-text-tertiary">填写要使用的图片生成模型名称</p>
    </div>

    <!-- 提示词 -->
    <div>
      <label
        class="block text-sm font-medium text-text-secondary mb-2 flex items-center justify-between"
      >
        <span>提示词 <span class="text-red-500">*</span></span>
        <button
          type="button"
          @click="showPromptHelper = !showPromptHelper"
          class="text-xs text-primary hover:text-primary"
        >
          {{ showPromptHelper ? '隐藏' : '显示' }}变量助手
        </button>
      </label>

      <!-- 变量助手 -->
      <VariableHelper
        :show="showPromptHelper"
        :previous-nodes="previousNodes"
        :env-vars="formattedEnvVars"
        @insert-field="
          (nodeId, fieldName) => insertFieldVariable(nodeId, fieldName, promptTextareaRef)
        "
        @insert-node="(nodeId) => insertNodeVariable(nodeId, promptTextareaRef)"
        @insert-env="(key) => insertEnvVariable(key, promptTextareaRef)"
      />

      <textarea
        ref="promptTextareaRef"
        v-model="localConfig.prompt"
        rows="4"
        placeholder="描述你想要生成的图片内容&#10;&#10;示例：一只可爱的橘猫坐在窗台上，阳光透过玻璃洒在它身上"
        class="w-full px-3 py-2 border-2 border-border-primary rounded-lg focus:outline-none focus:border-primary bg-bg-elevated text-text-primary font-mono text-sm"
      />
      <p class="mt-1 text-xs text-text-tertiary">
        支持使用变量，点击上方"显示变量助手"按钮选择前置节点的输出字段
      </p>
    </div>

    <!-- 图片尺寸 -->
    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2"> 图片尺寸 </label>
      <BaseSelect v-model="localConfig.size" :options="sizeOptions" />
      <p class="mt-1 text-xs text-text-tertiary">dall-e-3 支持 1024x1024、1792x1024、1024x1792</p>
    </div>

    <!-- 图片质量（可选） -->
    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2"> 图片质量（可选） </label>
      <BaseSelect v-model="localConfig.quality" :options="qualityOptions" />
      <p class="mt-1 text-xs text-text-tertiary">仅部分模型支持，dall-e-3 支持 hd 高清模式</p>
    </div>

    <!-- 生成数量 -->
    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2"> 生成数量（可选） </label>
      <BaseInput v-model.number="localConfig.n" type="number" min="1" max="10" placeholder="1" />
      <p class="mt-1 text-xs text-text-tertiary">
        部分模型支持一次生成多张图片（dall-e-2 支持 1-10 张，dall-e-3 仅支持 1 张）
      </p>
    </div>

    <!-- 返回格式 -->
    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2"> 返回格式 </label>
      <BaseSelect v-model="localConfig.response_format" :options="responseFormatOptions" />
      <p class="mt-1 text-xs text-text-tertiary">
        url 返回图片链接，b64_json 返回 base64 编码的图片数据
      </p>
    </div>

    <!-- 超时时间 -->
    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2"> 超时时间 (秒) </label>
      <BaseInput v-model.number="localConfig.timeout" type="number" placeholder="300" />
      <p class="mt-1 text-xs text-text-tertiary">默认 300 秒，图片生成通常需要 30-60 秒</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import BaseInput from '@/components/BaseInput'
import BaseSelect from '@/components/BaseSelect'
import VariableHelper from '@/components/VariableHelper'

interface Props {
  config: Record<string, any>
  previousNodes?: Array<{ id: string; name: string; type: string; toolCode?: string }>
  envVars?: Array<{ key: string; value: string; description?: string }>
}

const props = defineProps<Props>()

// 格式化环境变量
const formattedEnvVars = computed(() => {
  return props.envVars || []
})

const emit = defineEmits<{
  'update:config': [config: Record<string, any>]
}>()

// 变量助手显示状态
const showPromptHelper = ref(false)

// textarea ref
const promptTextareaRef = ref<HTMLTextAreaElement | null>(null)

// 尺寸选项
const sizeOptions = [
  { label: '1024x1024（正方形）', value: '1024x1024' },
  { label: '1792x1024（横向）', value: '1792x1024' },
  { label: '1024x1792（纵向）', value: '1024x1792' },
  { label: '512x512（DALL-E 2）', value: '512x512' },
  { label: '256x256（DALL-E 2）', value: '256x256' },
]

// 质量选项
const qualityOptions = [
  { label: '不使用（使用模型默认）', value: '' },
  { label: 'Standard（标准）', value: 'standard' },
  { label: 'HD（高清）', value: 'hd' },
]

// 返回格式选项
const responseFormatOptions = [
  { label: 'URL（图片链接）', value: 'url' },
  { label: 'Base64（编码数据）', value: 'b64_json' },
]

const localConfig = ref({
  model: props.config.model || '',
  prompt: props.config.prompt || '',
  n: props.config.n || 1,
  size: props.config.size || '1024x1024',
  quality: props.config.quality || '',
  response_format: props.config.response_format || 'url',
  timeout: props.config.timeout || 300,
})

// 插入字段变量
const insertFieldVariable = (
  nodeId: string,
  fieldName: string,
  textarea: HTMLTextAreaElement | null
) => {
  if (!textarea) return
  const variable = `{{nodes.${nodeId}.${fieldName}}}`
  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const text = textarea.value
  const before = text.substring(0, start)
  const after = text.substring(end, text.length)

  localConfig.value.prompt = before + variable + after

  setTimeout(() => {
    textarea.focus()
    textarea.selectionStart = textarea.selectionEnd = start + variable.length
  }, 0)
}

// 插入节点所有输出
const insertNodeVariable = (nodeId: string, textarea: HTMLTextAreaElement | null) => {
  if (!textarea) return
  const variable = `{{nodes.${nodeId}}}`
  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const text = textarea.value
  const before = text.substring(0, start)
  const after = text.substring(end, text.length)

  localConfig.value.prompt = before + variable + after

  setTimeout(() => {
    textarea.focus()
    textarea.selectionStart = textarea.selectionEnd = start + variable.length
  }, 0)
}

// 插入环境变量
const insertEnvVariable = (key: string, textarea: HTMLTextAreaElement | null) => {
  if (!textarea) return
  const variable = `{{env.${key}}}`
  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const text = textarea.value
  const before = text.substring(0, start)
  const after = text.substring(end, text.length)

  localConfig.value.prompt = before + variable + after

  setTimeout(() => {
    textarea.focus()
    textarea.selectionStart = textarea.selectionEnd = start + variable.length
  }, 0)
}

// 监听本地配置变化，同步到父组件
watch(
  localConfig,
  (newConfig) => {
    emit('update:config', { ...newConfig })
  },
  { deep: true }
)

// 监听外部配置变化
watch(
  () => props.config,
  (newConfig) => {
    const hasChanged = Object.keys(localConfig.value).some((key) => {
      return localConfig.value[key] !== (newConfig[key] ?? localConfig.value[key])
    })
    if (hasChanged) {
      localConfig.value = {
        model: newConfig.model || '',
        prompt: newConfig.prompt || '',
        n: newConfig.n || 1,
        size: newConfig.size || '1024x1024',
        quality: newConfig.quality || '',
        response_format: newConfig.response_format || 'url',
        timeout: newConfig.timeout || 300,
      }
    }
  },
  { immediate: true }
)
</script>
