<template>
  <div class="space-y-4">
    <h3 class="text-sm font-semibold text-text-primary mb-3">二维码生成配置</h3>

    <!-- 二维码内容 -->
    <div>
      <label class="block text-xs font-medium text-text-secondary mb-1.5">
        二维码内容 <span class="text-error">*</span>
      </label>
      <textarea
        v-model="localConfig.content"
        placeholder="输入要编码的内容（URL、文本等）&#10;支持变量：{{nodes.xxx.yyy}} 或 {{external.zzz}}"
        rows="3"
        class="w-full px-3 py-2 text-sm bg-bg-primary text-text-primary border border-border-primary rounded-md focus:ring-2 focus:ring-primary focus:border-primary resize-none"
      />
      <p class="mt-1 text-xs text-text-tertiary">支持 URL、文本、vCard 等任意内容</p>
    </div>

    <!-- 图片尺寸 -->
    <div>
      <label class="block text-xs font-medium text-text-secondary mb-1.5"> 图片尺寸（像素） </label>
      <BaseInput v-model="localConfig.size" type="number" min="64" max="2048" placeholder="256" />
      <p class="mt-1 text-xs text-text-tertiary">范围：64-2048，默认 256</p>
    </div>

    <!-- 错误纠正级别 -->
    <div>
      <label class="block text-xs font-medium text-text-secondary mb-1.5"> 错误纠正级别 </label>
      <BaseSelect v-model="localConfig.level" :options="levelOptions" />
      <p class="mt-1 text-xs text-text-tertiary">容错能力越高，二维码尺寸越大，但损坏后仍可识别</p>
    </div>

    <!-- 输出格式 -->
    <div>
      <label class="block text-xs font-medium text-text-secondary mb-1.5"> 输出格式 </label>
      <BaseSelect v-model="localConfig.output_format" :options="formatOptions" />
      <p class="mt-1 text-xs text-text-tertiary">
        Base64：可用于直接显示；File 对象：可传递给上传工具
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import BaseInput from '@/components/BaseInput'
import BaseSelect from '@/components/BaseSelect'

interface Props {
  config: Record<string, any>
}

const props = defineProps<Props>()
const emit = defineEmits<{
  (e: 'update:config', value: Record<string, any>): void
}>()

const levelOptions = [
  { label: 'Low（低 - 7% 容错）', value: 'Low' },
  { label: 'Medium（中 - 15% 容错）', value: 'Medium' },
  { label: 'High（高 - 25% 容错）', value: 'High' },
  { label: 'Highest（最高 - 30% 容错）', value: 'Highest' },
]

const formatOptions = [
  { label: 'Base64 编码', value: 'base64' },
  { label: 'File 对象', value: 'file' },
]

const localConfig = ref({
  content: props.config.content || '',
  size: props.config.size || 256,
  level: props.config.level || 'Medium',
  output_format: props.config.output_format || 'base64',
})

watch(
  localConfig,
  (newConfig) => {
    emit('update:config', newConfig)
  },
  { deep: true }
)
</script>
