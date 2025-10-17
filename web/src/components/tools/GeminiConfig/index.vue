<template>
  <div class="space-y-4">
    <h3 class="text-sm font-semibold text-text-primary mb-3">Gemini AI 配置</h3>

    <!-- 模型名称 -->
    <div>
      <label class="block text-xs font-medium text-text-secondary mb-1.5"> 模型名称 </label>
      <BaseInput v-model="localConfig.model" placeholder="gemini-pro 或 {{nodes.xxx.model}}" />
      <p class="mt-1 text-xs text-text-tertiary">
        常用模型：gemini-pro、gemini-1.5-pro、gemini-1.5-flash、gemini-pro-vision
      </p>
    </div>

    <!-- 提示词 -->
    <div>
      <label class="block text-xs font-medium text-text-secondary mb-1.5">
        提示词 <span class="text-error">*</span>
      </label>
      <textarea
        v-model="localConfig.prompt"
        placeholder="输入你想让 AI 回答或处理的内容&#10;支持变量：{{nodes.xxx.yyy}} 或 {{external.zzz}}"
        rows="4"
        class="w-full px-3 py-2 text-sm bg-bg-primary text-text-primary border border-border-primary rounded-md focus:ring-2 focus:ring-primary focus:border-primary resize-none"
      />
      <p class="mt-1 text-xs text-text-tertiary">支持多行文本，可以引用前置节点的输出</p>
    </div>

    <!-- 图片输入 -->
    <div>
      <label class="block text-xs font-medium text-text-secondary mb-1.5"> 图片输入（可选） </label>
      <BaseInput
        v-model="localConfig.image"
        placeholder="{{nodes.xxx.file}} 或 {{nodes.xxx.data}}"
      />
      <p class="mt-1 text-xs text-text-tertiary">
        仅 vision 模型支持，可传入：① 文件对象 ② Base64 字符串 ③ Data URI
      </p>
    </div>

    <!-- 系统指令 -->
    <div>
      <label class="block text-xs font-medium text-text-secondary mb-1.5"> 系统指令（可选） </label>
      <textarea
        v-model="localConfig.system_instruction"
        placeholder="设定 AI 的角色和行为，例如：&#10;你是一个专业的文案写作助手，请用简洁、有吸引力的语言回答。"
        rows="3"
        class="w-full px-3 py-2 text-sm bg-bg-primary text-text-primary border border-border-primary rounded-md focus:ring-2 focus:ring-primary focus:border-primary resize-none"
      />
      <p class="mt-1 text-xs text-text-tertiary">定义 AI 的性格、专业领域或回答风格</p>
    </div>

    <!-- 高级参数（可折叠） -->
    <details class="border border-border-primary rounded-md p-3">
      <summary class="cursor-pointer text-sm font-medium text-text-primary">
        高级参数（可选）
      </summary>
      <div class="mt-3 space-y-4">
        <!-- 温度 -->
        <div>
          <label class="block text-xs font-medium text-text-secondary mb-1.5">
            温度 (Temperature)
          </label>
          <BaseInput
            v-model.number="localConfig.temperature"
            type="number"
            min="0"
            max="2"
            step="0.1"
            placeholder="0.7"
          />
          <p class="mt-1 text-xs text-text-tertiary">
            控制输出的随机性，0-2 之间。越高越随机，越低越确定
          </p>
        </div>

        <!-- 最大 Token 数 -->
        <div>
          <label class="block text-xs font-medium text-text-secondary mb-1.5">
            最大 Token 数
          </label>
          <BaseInput
            v-model.number="localConfig.max_tokens"
            type="number"
            min="1"
            placeholder="2048"
          />
          <p class="mt-1 text-xs text-text-tertiary">生成内容的最大长度，不同模型限制不同</p>
        </div>

        <!-- Top P -->
        <div>
          <label class="block text-xs font-medium text-text-secondary mb-1.5"> Top P </label>
          <BaseInput
            v-model.number="localConfig.top_p"
            type="number"
            min="0"
            max="1"
            step="0.05"
            placeholder="0.95"
          />
          <p class="mt-1 text-xs text-text-tertiary">核采样参数，0-1 之间，通常保持默认</p>
        </div>

        <!-- Top K -->
        <div>
          <label class="block text-xs font-medium text-text-secondary mb-1.5"> Top K </label>
          <BaseInput v-model.number="localConfig.top_k" type="number" min="1" placeholder="40" />
          <p class="mt-1 text-xs text-text-tertiary">采样时考虑的 token 数量</p>
        </div>
      </div>
    </details>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import BaseInput from '@/components/BaseInput'

interface Props {
  config: Record<string, any>
}

const props = defineProps<Props>()
const emit = defineEmits<{
  (e: 'update:config', value: Record<string, any>): void
}>()

const localConfig = ref({
  model: props.config.model || 'gemini-pro',
  prompt: props.config.prompt || '',
  image: props.config.image || '',
  system_instruction: props.config.system_instruction || '',
  temperature: props.config.temperature !== undefined ? props.config.temperature : 0.7,
  max_tokens: props.config.max_tokens || 2048,
  top_p: props.config.top_p !== undefined ? props.config.top_p : 0.95,
  top_k: props.config.top_k || 40,
})

watch(
  localConfig,
  (newConfig) => {
    emit('update:config', newConfig)
  },
  { deep: true }
)
</script>
