<template>
  <div class="space-y-4">
    <h3 class="text-sm font-semibold text-text-primary mb-3">腾讯云 COS 上传配置</h3>

    <!-- 文件参数 -->
    <div>
      <label class="block text-xs font-medium text-text-secondary mb-1.5">
        文件 <span class="text-error">*</span>
      </label>
      <input
        v-model="localConfig.file"
        type="text"
        placeholder="文件路径或使用变量 &#123;&#123;nodes.xxx.file&#125;&#125;"
        class="w-full px-3 py-2 text-sm font-mono border border-border-primary rounded-md focus:ring-2 focus:ring-primary focus:border-primary"
        style="color: var(--color-text-primary); background-color: var(--color-bg-elevated);"
      />
      <p class="text-xs text-text-tertiary mt-1">
        支持变量：&#123;&#123;nodes.xxx.file&#125;&#125;、&#123;&#123;external.file&#125;&#125;
      </p>
    </div>

    <!-- COS 存储路径（可选） -->
    <div>
      <label class="block text-xs font-medium text-text-secondary mb-1.5">
        存储路径（可选）
      </label>
      <input
        v-model="localConfig.path"
        type="text"
        placeholder="默认使用文件名，如: images/2024/01/15/photo.jpg"
        class="w-full px-3 py-2 text-sm border border-border-primary rounded-md focus:ring-2 focus:ring-primary focus:border-primary"
        style="color: var(--color-text-primary); background-color: var(--color-bg-elevated);"
      />
      <p class="text-xs text-text-tertiary mt-1">
        不填则使用文件原名，支持路径层级
      </p>
    </div>

    <!-- 说明 -->
    <div class="bg-bg-elevated rounded-lg p-3 border border-border-primary">
      <p class="text-xs text-text-secondary leading-relaxed">
        💡 <span class="font-medium">配置说明：</span><br />
        COS 配置（SecretId、SecretKey、Bucket等）在后端配置文件中统一管理，无需在此填写。
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

interface Props {
  config: Record<string, any>
}

interface Emits {
  (e: 'update:config', value: Record<string, any>): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 本地配置状态
const localConfig = ref({
  file: props.config.file || '',
  path: props.config.path || '',
})

// 监听配置变化并向父组件发送更新
watch(
  localConfig,
  (newConfig) => {
    emit('update:config', newConfig)
  },
  { deep: true }
)
</script>

