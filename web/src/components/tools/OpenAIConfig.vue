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
      <BaseInput
        v-model="localConfig.model"
        placeholder="例如：gpt-3.5-turbo, gpt-4, gpt-4o, dall-e-3"
      />
      <p class="mt-1 text-xs text-text-tertiary">填写要使用的模型名称</p>
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
        placeholder="请输入要发送给 ChatGPT 的问题或指令&#10;&#10;示例：分析以下数据：{{nodes.xxx.data}}"
        class="w-full px-3 py-2 border-2 border-border-primary rounded-lg focus:outline-none focus:border-primary bg-bg-elevated text-text-primary font-mono text-sm"
      />
      <p class="mt-1 text-xs text-text-tertiary">
        支持使用变量，点击上方"显示变量助手"按钮选择前置节点的输出字段
      </p>
    </div>

    <!-- 系统消息 -->
    <div>
      <label
        class="block text-sm font-medium text-text-secondary mb-2 flex items-center justify-between"
      >
        <span>系统消息 (可选)</span>
        <button
          type="button"
          @click="showSystemHelper = !showSystemHelper"
          class="text-xs text-primary hover:text-primary"
        >
          {{ showSystemHelper ? '隐藏' : '显示' }}变量助手
        </button>
      </label>

      <!-- 变量助手 -->
      <VariableHelper
        :show="showSystemHelper"
        :previous-nodes="previousNodes"
        :env-vars="formattedEnvVars"
        @insert-field="
          (nodeId, fieldName) => insertFieldVariable(nodeId, fieldName, systemTextareaRef)
        "
        @insert-node="(nodeId) => insertNodeVariable(nodeId, systemTextareaRef)"
        @insert-env="(key) => insertEnvVariable(key, systemTextareaRef)"
      />

      <textarea
        ref="systemTextareaRef"
        v-model="localConfig.system_message"
        rows="3"
        placeholder="设定 AI 的角色和行为，例如：你是一个专业的数据分析师"
        class="w-full px-3 py-2 border-2 border-border-primary rounded-lg focus:outline-none focus:border-primary bg-bg-elevated text-text-primary font-mono text-sm"
      />
    </div>

    <!-- 温度 -->
    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2">
        温度: {{ localConfig.temperature }}
      </label>
      <input
        v-model.number="localConfig.temperature"
        type="range"
        min="0"
        max="2"
        step="0.1"
        class="w-full h-2 bg-bg-hover rounded-lg appearance-none cursor-pointer"
      />
      <div class="flex justify-between text-xs text-text-tertiary mt-1">
        <span>精确 (0)</span>
        <span>平衡 (1)</span>
        <span>创造性 (2)</span>
      </div>
      <p class="mt-1 text-xs text-text-tertiary">控制回复的随机性，值越高越有创意但可能不够准确</p>
    </div>

    <!-- 最大 Token 数 -->
    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2">
        最大 Token 数 (可选)
      </label>
      <BaseInput
        v-model.number="localConfig.max_tokens"
        type="number"
        placeholder="留空使用模型默认值"
      />
      <p class="mt-1 text-xs text-text-tertiary">限制生成回复的长度，留空则使用模型默认值</p>
    </div>

    <!-- 超时时间 -->
    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2"> 超时时间 (秒) </label>
      <BaseInput v-model.number="localConfig.timeout" type="number" placeholder="300" />
      <p class="mt-1 text-xs text-text-tertiary">默认 300 秒，对于图片生成等耗时操作可适当增加</p>
    </div>

    <!-- 对话记忆 -->
    <div class="border-t border-border-primary pt-4 space-y-3">
      <div class="flex items-center justify-between">
        <div>
          <div class="text-sm font-medium text-text-secondary">对话记忆</div>
          <div class="text-xs text-text-tertiary">启用后会按会话ID保留上下文</div>
        </div>
        <label class="inline-flex items-center cursor-pointer">
          <input type="checkbox" v-model="localConfig.memory.enabled" class="mr-2" />
          <span class="text-sm text-text-secondary">启用</span>
        </label>
      </div>

      <div v-if="localConfig.memory.enabled" class="space-y-3">
        <div>
          <label class="block text-sm font-medium text-text-secondary mb-1">会话ID</label>
          <BaseInput
            v-model="localConfig.memory.sessionKey"
            placeholder="例如：{{params.session_id}} 或 {{user.id}}"
          />
          <p class="mt-1 text-xs text-text-tertiary">
            用于区分不同对话（为空时回退到 {{ user.id }} 或 global）
          </p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
          <div>
            <label class="block text-sm font-medium text-text-secondary mb-1">窗口条数</label>
            <BaseInput
              v-model.number="localConfig.memory.windowSize"
              type="number"
              placeholder="10"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-text-secondary mb-1"
              >最大上下文 Tokens</label
            >
            <BaseInput
              v-model.number="localConfig.memory.maxTokens"
              type="number"
              placeholder="2000"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-text-secondary mb-1">TTL（秒）</label>
            <BaseInput
              v-model.number="localConfig.memory.ttlSeconds"
              type="number"
              placeholder="604800"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import BaseInput from '@/components/BaseInput'
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
const showSystemHelper = ref(false)

// textarea refs
const promptTextareaRef = ref<HTMLTextAreaElement | null>(null)
const systemTextareaRef = ref<HTMLTextAreaElement | null>(null)

const localConfig = ref({
  model: props.config.model || 'gpt-3.5-turbo',
  prompt: props.config.prompt || '',
  system_message: props.config.system_message || '',
  temperature: props.config.temperature ?? 0.7,
  max_tokens: props.config.max_tokens || '',
  timeout: props.config.timeout || 300,
  memory: {
    enabled: props.config.memory?.enabled ?? false,
    sessionKey: props.config.memory?.sessionKey || '',
    windowSize: props.config.memory?.windowSize ?? 10,
    maxTokens: props.config.memory?.maxTokens ?? 2000,
    ttlSeconds: props.config.memory?.ttlSeconds ?? 604800,
  },
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

  // 更新对应的 localConfig
  if (textarea === promptTextareaRef.value) {
    localConfig.value.prompt = before + variable + after
  } else if (textarea === systemTextareaRef.value) {
    localConfig.value.system_message = before + variable + after
  }

  // 设置光标位置
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

  // 更新对应的 localConfig
  if (textarea === promptTextareaRef.value) {
    localConfig.value.prompt = before + variable + after
  } else if (textarea === systemTextareaRef.value) {
    localConfig.value.system_message = before + variable + after
  }

  // 设置光标位置
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

  // 更新对应的 localConfig
  if (textarea === promptTextareaRef.value) {
    localConfig.value.prompt = before + variable + after
  } else if (textarea === systemTextareaRef.value) {
    localConfig.value.system_message = before + variable + after
  }

  // 设置光标位置
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
        model: newConfig.model || 'gpt-3.5-turbo',
        prompt: newConfig.prompt || '',
        system_message: newConfig.system_message || '',
        temperature: newConfig.temperature ?? 0.7,
        max_tokens: newConfig.max_tokens || '',
        timeout: newConfig.timeout || 300,
        memory: {
          enabled: newConfig.memory?.enabled ?? false,
          sessionKey: newConfig.memory?.sessionKey || '',
          windowSize: newConfig.memory?.windowSize ?? 10,
          maxTokens: newConfig.memory?.maxTokens ?? 2000,
          ttlSeconds: newConfig.memory?.ttlSeconds ?? 604800,
        },
      }
    }
  },
  { immediate: true }
)
</script>

<style scoped>
input[type='range']::-webkit-slider-thumb {
  appearance: none;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: var(--color-primary);
  cursor: pointer;
}

input[type='range']::-moz-range-thumb {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: var(--color-primary);
  cursor: pointer;
  border: none;
}
</style>
