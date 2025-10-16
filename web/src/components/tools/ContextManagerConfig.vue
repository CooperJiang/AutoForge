<template>
  <div class="space-y-4">
    <div class="bg-indigo-500/10 border-l-4 border-indigo-500 rounded p-3">
      <p class="text-sm text-indigo-800 dark:text-indigo-300">
        💬 <strong>对话上下文管理器</strong>：管理多轮对话的历史记录，支持两种模式
      </p>
      <ul class="mt-2 text-xs text-indigo-700 dark:text-indigo-400 space-y-1 ml-4">
        <li>• <strong>Prepare</strong>：读取历史 + 拼接当前消息 → 传给 LLM 工具</li>
        <li>• <strong>Persist</strong>：将 LLM 回复追加到历史 → 保存到 Redis</li>
      </ul>
    </div>

    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2">
        工作模式 <span class="text-red-500">*</span>
      </label>
      <BaseSelect
        v-model="localConfig.mode"
        :options="modeOptions"
        @update:model-value="onModeChange"
      />
      <p class="mt-1 text-xs text-text-tertiary">Prepare = 准备消息，Persist = 保存消息</p>
    </div>

    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2">
        会话标识 <span class="text-red-500">*</span>
      </label>
      <BaseInput
        v-model="localConfig.session_key"
        placeholder="例如：{{params.session_id}} 或 {{params.user_id}}"
      />
      <p class="mt-1 text-xs text-text-tertiary">用于区分不同用户的对话。支持变量引用。</p>
    </div>

    <div
      v-if="localConfig.mode === 'prepare'"
      class="border-t border-border-primary pt-4 space-y-3"
    >
      <h4 class="text-sm font-medium text-text-primary">Prepare 模式配置</h4>

      <div>
        <label class="block text-sm font-medium text-text-secondary mb-1">
          用户输入 <span class="text-red-500">*</span>
        </label>
        <BaseInput v-model="localConfig.user_input" placeholder="{{params.user_message}}" />
        <p class="mt-1 text-xs text-text-tertiary">当前用户发送的消息，通常引用外部参数</p>
      </div>

      <div>
        <label class="block text-sm font-medium text-text-secondary mb-1"> 系统消息（可选） </label>
        <textarea
          v-model="localConfig.system_message"
          rows="3"
          class="w-full px-3 py-2 border-2 border-border-primary rounded-lg focus:outline-none focus:border-primary bg-bg-elevated text-text-primary font-mono text-sm"
          placeholder="你是一个友好的 AI 助手..."
        />
        <p class="mt-1 text-xs text-text-tertiary">设定 AI 的角色和行为</p>
      </div>
    </div>

    <div
      v-if="localConfig.mode === 'persist'"
      class="border-t border-border-primary pt-4 space-y-3"
    >
      <h4 class="text-sm font-medium text-text-primary">Persist 模式配置</h4>

      <div>
        <label class="block text-sm font-medium text-text-secondary mb-1"> 用户输入（可选） </label>
        <BaseInput v-model="localConfig.user_input" placeholder="{{params.user_message}}" />
        <p class="mt-1 text-xs text-text-tertiary">
          如果需要保存用户消息（通常在 Prepare 中已保存）
        </p>
      </div>

      <div>
        <label class="block text-sm font-medium text-text-secondary mb-1">
          AI 回复 <span class="text-red-500">*</span>
        </label>
        <BaseInput
          v-model="localConfig.assistant_output"
          placeholder="{{nodes.openai_chat.content}}"
        />
        <p class="mt-1 text-xs text-text-tertiary">
          AI 助手的回复内容，通常引用上游 LLM 节点的输出
        </p>
      </div>
    </div>

    <div class="border-t border-border-primary pt-4 space-y-3">
      <h4 class="text-sm font-medium text-text-primary">通用配置</h4>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <div>
          <label class="block text-sm font-medium text-text-secondary mb-1">窗口大小</label>
          <BaseInput v-model.number="localConfig.window_size" type="number" placeholder="10" />
          <p class="mt-1 text-xs text-text-tertiary">保留最近 N 条消息</p>
        </div>

        <div>
          <label class="block text-sm font-medium text-text-secondary mb-1">过期时间（秒）</label>
          <BaseInput v-model.number="localConfig.ttl_seconds" type="number" placeholder="604800" />
          <p class="mt-1 text-xs text-text-tertiary">默认 7 天（604800 秒）</p>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <input
          type="checkbox"
          v-model="localConfig.clear_history"
          class="w-4 h-4 text-primary bg-bg-elevated border-border-primary rounded focus:ring-primary"
        />
        <label class="text-sm text-text-secondary"> 执行前清空历史记录（用于重新开始对话） </label>
      </div>
    </div>

    <div class="bg-bg-hover rounded-lg p-3">
      <div class="text-xs font-semibold text-text-secondary mb-2">典型工作流配置：</div>
      <div class="text-xs text-text-secondary space-y-1 font-mono">
        <div>1. 上下文管理器 (Prepare) → 读取历史 + 拼接消息</div>
        <div>2. OpenAI Chat → 使用上游的 messages_json</div>
        <div>3. 上下文管理器 (Persist) → 保存 AI 回复</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'
import BaseInput from '@/components/BaseInput'
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

const modeOptions = [
  { label: 'Prepare - 准备消息', value: 'prepare' },
  { label: 'Persist - 保存消息', value: 'persist' },
]

// 初始化本地配置
const localConfig = ref({
  mode: props.config.mode || 'prepare',
  session_key: props.config.session_key || '{{external.session_id}}',
  user_input: props.config.user_input || '{{external.user_input}}',
  assistant_output: props.config.assistant_output || '',
  system_message: props.config.system_message || '',
  window_size: props.config.window_size ?? 10,
  ttl_seconds: props.config.ttl_seconds ?? 604800,
  clear_history: props.config.clear_history ?? false,
})

// 阻止循环更新的标志
let isUpdatingFromProps = false

const onModeChange = () => {
  // 根据模式设置默认值
  if (localConfig.value.mode === 'prepare') {
    if (!localConfig.value.user_input) {
      localConfig.value.user_input = '{{external.user_input}}'
    }
  } else if (localConfig.value.mode === 'persist') {
    if (!localConfig.value.assistant_output) {
      localConfig.value.assistant_output = '{{nodes.openai_chat.content}}'
    }
  }
}

// 监听 localConfig 变化，触发更新事件
watch(
  localConfig,
  (val) => {
    // 如果是从 props 同步过来的更新，不再触发 emit
    if (isUpdatingFromProps) {
      return
    }
    emit('update:config', { ...val })
  },
  { deep: true }
)

// 监听外部 config 变化，同步到 localConfig
watch(
  () => props.config,
  (newVal) => {
    if (!newVal) return

    // 设置标志，防止触发 emit
    isUpdatingFromProps = true

    localConfig.value = {
      mode: newVal.mode || 'prepare',
      session_key: newVal.session_key || '{{params.session_id}}',
      user_input: newVal.user_input || '',
      assistant_output: newVal.assistant_output || '',
      system_message: newVal.system_message || '',
      window_size: newVal.window_size ?? 10,
      ttl_seconds: newVal.ttl_seconds ?? 604800,
      clear_history: newVal.clear_history ?? false,
    }

    // 下一个 tick 后重置标志
    nextTick(() => {
      isUpdatingFromProps = false
    })
  },
  { immediate: true, deep: true }
)
</script>
