<template>
  <div v-if="show" class="mb-2 p-3 bg-bg-hover rounded-lg border border-border-primary">
    <div class="text-xs font-semibold text-text-secondary mb-2">可用变量：</div>
    <div class="space-y-3">
      
      <div v-if="previousNodes && previousNodes.length > 0">
        <div class="text-xs text-text-secondary mb-2">前置节点输出：</div>
        <VariableTreeSelector :previous-nodes="previousNodes" />
      </div>

      
      <div v-if="envVars && envVars.length > 0">
        <div class="text-xs text-text-secondary mb-1">环境变量：</div>
        <div class="space-y-1">
          <div v-for="envVar in envVars" :key="envVar.key" class="text-xs">
            <button
              type="button"
              @click="copyToClipboard(`{{env.${envVar.key}}}`, envVar.key)"
              class="font-mono text-primary hover:text-primary hover:underline"
            >
              {{ getEnvVariableText(envVar.key) }}
            </button>
            <span class="text-text-tertiary ml-1">- {{ envVar.description }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import VariableTreeSelector from '@/components/VariableTreeSelector'
import { message } from '@/utils/message'

interface Props {
  show: boolean
  previousNodes?: Array<{ id: string; name: string; type: string; toolCode?: string }>
  envVars?: Array<{ key: string; value: string; description?: string }>
}

defineProps<Props>()

// 复制到剪贴板
const copyToClipboard = async (text: string, _label: string) => {
  try {
    await navigator.clipboard.writeText(text)
    message.success(`已复制: ${text}`)
  } catch (err) {
    console.error('复制失败:', err)
    message.error('复制失败,请手动复制')
  }
}

// 获取环境变量文本
const getEnvVariableText = (key: string) => {
  return `{{env.${key}}}`
}
</script>
