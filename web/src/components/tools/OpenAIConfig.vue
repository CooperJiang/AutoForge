<template>
  <div class="space-y-4">
    <div class="bg-primary-light border-l-4 border-primary p-3 rounded">
      <p class="text-sm text-text-primary">💡 API 凭证已在系统配置中统一管理，无需每次填写</p>
    </div>

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

    <div class="bg-blue-500/10 border border-blue-500/20 rounded-lg p-3">
      <p class="text-sm text-blue-800 dark:text-blue-300 font-medium mb-2">💡 配置方式选择</p>
      <ul class="text-xs text-blue-700 dark:text-blue-400 space-y-1">
        <li>• <strong>简单场景</strong>：填写"提示词"和"系统消息"</li>
        <li>• <strong>多轮对话</strong>：使用下方"消息数组 (JSON)"，配合上下文管理器</li>
        <li>• <strong>优先级</strong>：消息数组 > 提示词+系统消息</li>
      </ul>
    </div>

    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2"> 提示词 (可选) </label>
      <textarea
        v-model="localConfig.prompt"
        rows="4"
        placeholder="请输入要发送给 ChatGPT 的问题或指令&#10;&#10;示例：分析以下数据"
        class="w-full px-3 py-2 border-2 border-border-primary rounded-lg focus:outline-none focus:border-primary bg-bg-elevated text-text-primary font-mono text-sm"
      />
      <p class="mt-1 text-xs text-text-tertiary">
        支持使用变量，格式：<code class="text-primary font-mono"
          >&#123;&#123;nodes.xxx.yyy&#125;&#125;</code
        >
        或 <code class="text-primary font-mono">&#123;&#123;external.zzz&#125;&#125;</code>
      </p>
    </div>

    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2"> 系统消息 (可选) </label>
      <textarea
        v-model="localConfig.system_message"
        rows="3"
        placeholder="设定 AI 的角色和行为，例如：你是一个专业的数据分析师"
        class="w-full px-3 py-2 border-2 border-border-primary rounded-lg focus:outline-none focus:border-primary bg-bg-elevated text-text-primary font-mono text-sm"
      />
      <p class="mt-1 text-xs text-text-tertiary">定义 AI 的角色和行为方式</p>
    </div>

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

    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2"> 超时时间 (秒) </label>
      <BaseInput v-model.number="localConfig.timeout" type="number" placeholder="300" />
      <p class="mt-1 text-xs text-text-tertiary">默认 300 秒，对于图片生成等耗时操作可适当增加</p>
    </div>

    <div class="border-t border-border-primary pt-4">
      <label class="block text-sm font-medium text-text-secondary mb-2">
        消息数组 (JSON，可选，存在时优先使用)
      </label>
      <textarea
        v-model="localConfig.messages_json"
        rows="6"
        placeholder='[ {"role":"system","content":"你是一个有帮助的助手"}, {"role":"user","content":"你好"} ]'
        class="w-full px-3 py-2 border-2 border-border-primary rounded-lg focus:outline-none focus:border-primary bg-bg-elevated text-text-primary font-mono text-sm"
      />
      <p class="mt-1 text-xs text-text-tertiary">
        填写符合 OpenAI Chat Completions 的 messages 数组 JSON 字符串；若填写，将忽略上面的
        system/prompt 构造。
      </p>
    </div>

    <div class="border-t border-border-primary pt-4 space-y-3">
      <div class="bg-blue-500/10 border border-blue-500/20 rounded-lg p-3 mb-3">
        <p class="text-sm text-blue-800 dark:text-blue-300">
          💬 <strong>对话记忆</strong>：启用后，AI 会记住本次会话的历史对话内容，实现多轮对话
        </p>
      </div>

      <div class="flex items-center justify-between">
        <div>
          <div class="text-sm font-medium text-text-secondary">启用对话记忆</div>
          <div class="text-xs text-text-tertiary">适用于聊天机器人、客服助手等多轮对话场景</div>
        </div>
        <label class="inline-flex items-center cursor-pointer">
          <input
            type="checkbox"
            v-model="localConfig.context_config.enabled"
            class="w-4 h-4 text-primary bg-bg-elevated border-border-primary rounded focus:ring-primary"
          />
        </label>
      </div>

      <div v-if="localConfig.context_config.enabled" class="space-y-3 pt-2">
        <div>
          <label class="block text-sm font-medium text-text-secondary mb-1">
            会话ID <span class="text-red-500">*</span>
          </label>
          <BaseInput
            v-model="localConfig.context_config.session_key"
            :placeholder="placeholderSessionId"
          />
          <p class="mt-1 text-xs text-text-tertiary">
            用于区分不同用户的对话。支持变量：
            <code class="text-primary font-mono">&#123;&#123;params.user_id&#125;&#125;</code>、
            <code class="text-primary font-mono">&#123;&#123;params.session_id&#125;&#125;</code> 等
          </p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
          <div>
            <label class="block text-sm font-medium text-text-secondary mb-1">窗口条数</label>
            <BaseInput
              v-model.number="localConfig.context_config.window_size"
              type="number"
              placeholder="10"
            />
            <p class="mt-1 text-xs text-text-tertiary">保留最近 N 条消息（一问一答算2条）</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-text-secondary mb-1">过期时间（秒）</label>
            <BaseInput
              v-model.number="localConfig.context_config.ttl_seconds"
              type="number"
              placeholder="604800"
            />
            <p class="mt-1 text-xs text-text-tertiary">默认 7 天（604800 秒）</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import BaseInput from '@/components/BaseInput'

interface Props {
  config: Record<string, any>
  previousNodes?: Array<{ id: string; name: string; type: string; toolCode?: string }>
  envVars?: Array<{ key: string; value: string; description?: string }>
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:config': [config: Record<string, any>]
}>()

// Placeholder 常量（避免 Vue 模板解析错误）
const placeholderSessionId = '{{params.session_id}}'

const localConfig = ref({
  model: props.config.model || 'gpt-3.5-turbo',
  prompt: props.config.prompt || '',
  system_message: props.config.system_message || '',
  temperature: props.config.temperature ?? 0.7,
  max_tokens: props.config.max_tokens || '',
  timeout: props.config.timeout || 300,
  messages_json: props.config.messages_json || '',
  context_config: {
    enabled: props.config.context_config?.enabled ?? false,
    session_key: props.config.context_config?.session_key || '{{params.session_id}}',
    window_size: props.config.context_config?.window_size ?? 10,
    ttl_seconds: props.config.context_config?.ttl_seconds ?? 604800,
  },
})

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
        messages_json: newConfig.messages_json || '',
        context_config: {
          enabled: newConfig.context_config?.enabled ?? false,
          session_key: newConfig.context_config?.session_key || '{{params.session_id}}',
          window_size: newConfig.context_config?.window_size ?? 10,
          ttl_seconds: newConfig.context_config?.ttl_seconds ?? 604800,
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
