<template>
  <div class="space-y-4">
    
    <div class="bg-primary-light border-l-4 border-primary p-3 rounded">
      <p class="text-sm text-text-primary">💡 API 凭证已在系统配置中统一管理，无需每次填写</p>
    </div>

    
    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2">
        模型 <span class="text-red-500">*</span>
      </label>
      <BaseInput v-model="localConfig.model" placeholder="dall-e-3" />
      <p class="mt-1 text-xs text-text-tertiary">
        常用模型：dall-e-3、dall-e-2、gpt-image-1，支持变量
      </p>
    </div>

    
    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2">
        提示词 <span class="text-red-500">*</span>
      </label>
      <textarea
        v-model="localConfig.prompt"
        rows="4"
        placeholder="描述你想要生成的图片内容&#10;&#10;示例：一只可爱的橘猫坐在窗台上，阳光透过玻璃洒在它身上"
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
      <label class="block text-sm font-medium text-text-secondary mb-2"> 图片尺寸 </label>
      <BaseInput v-model="localConfig.size" placeholder="1024x1024" />
      <p class="mt-1 text-xs text-text-tertiary">
        dall-e-2: 256x256, 512x512, 1024x1024 / dall-e-3: 1024x1024, 1792x1024, 1024x1792，支持变量
      </p>
    </div>

    
    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2"> 图片质量（可选） </label>
      <BaseInput v-model="localConfig.quality" placeholder="standard" />
      <p class="mt-1 text-xs text-text-tertiary">
        standard（标准）或 hd（高清），留空使用模型默认值，支持变量
      </p>
    </div>

    
    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2"> 生成数量（可选） </label>
      <BaseInput v-model="localConfig.n" placeholder="1" />
      <p class="mt-1 text-xs text-text-tertiary">
        dall-e-2 支持 1-10 张，dall-e-3 仅支持 1 张，支持变量
      </p>
    </div>

    
    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2"> 返回格式 </label>
      <BaseInput v-model="localConfig.response_format" placeholder="url" />
      <p class="mt-1 text-xs text-text-tertiary">
        url（图片链接）或 b64_json（base64 编码），支持变量
      </p>
    </div>

    
    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2"> 超时时间 (秒) </label>
      <BaseInput v-model="localConfig.timeout" placeholder="300" />
      <p class="mt-1 text-xs text-text-tertiary">
        默认 300 秒，图片生成通常需要 30-60 秒，支持变量
      </p>
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

const localConfig = ref({
  model: props.config.model || '',
  prompt: props.config.prompt || '',
  n: props.config.n || '',
  size: props.config.size || '',
  quality: props.config.quality || '',
  response_format: props.config.response_format || '',
  timeout: props.config.timeout || '',
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
        model: newConfig.model || '',
        prompt: newConfig.prompt || '',
        n: newConfig.n || '',
        size: newConfig.size || '',
        quality: newConfig.quality || '',
        response_format: newConfig.response_format || '',
        timeout: newConfig.timeout || '',
      }
    }
  },
  { immediate: true }
)
</script>
