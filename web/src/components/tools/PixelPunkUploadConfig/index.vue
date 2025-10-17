<template>
  <div class="space-y-4">
    <h3 class="text-sm font-semibold text-text-primary mb-3">PixelPunk 图床上传配置</h3>

    <!-- 文件参数说明 -->
    <div class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-3">
      <p class="text-xs text-blue-800 dark:text-blue-300">
        <strong>📌 文件来源：</strong>从 <code class="px-1 bg-blue-100 dark:bg-blue-800 rounded">external_trigger</code>
        节点接收文件对象，使用 <code class="px-1 bg-blue-100 dark:bg-blue-800 rounded">&#123;&#123;external.image&#125;&#125;</code> 引用。
      </p>
    </div>

    <!-- 文件对象 -->
    <div>
      <label class="block text-xs font-medium text-text-secondary mb-1.5">
        文件对象 <span class="text-error">*</span>
      </label>
      <input
        v-model="localConfig.file"
        type="text"
        placeholder="{{external.image}}"
        class="w-full px-3 py-2 text-sm border border-border-primary rounded-md focus:ring-2 focus:ring-primary focus:border-primary font-mono"
        style="color: var(--color-text-primary); background-color: var(--color-bg-elevated);"
      />
      <p class="text-xs text-text-tertiary mt-1">
        通常使用 &#123;&#123;external.xxx&#125;&#125; 引用从 API 接收的文件
      </p>
    </div>

    <!-- 访问级别 -->
    <div>
      <label class="block text-xs font-medium text-text-secondary mb-1.5">
        访问级别
      </label>
      <BaseSelect
        v-model="localConfig.access_level"
        :options="accessLevelOptions"
      />
      <p class="text-xs text-text-tertiary mt-1">
        控制图片的访问权限
      </p>
    </div>

    <!-- 优化图片 -->
    <div class="flex items-center justify-between">
      <div>
        <label class="block text-xs font-medium text-text-secondary mb-1">
          优化图片
        </label>
        <p class="text-xs text-text-tertiary">
          自动压缩图片以减少文件大小
        </p>
      </div>
      <label class="relative inline-flex items-center cursor-pointer">
        <input
          v-model="localConfig.optimize"
          type="checkbox"
          class="sr-only peer"
        />
        <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary/20 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-primary"></div>
      </label>
    </div>

    <!-- 虚拟路径（可选） -->
    <div>
      <label class="block text-xs font-medium text-text-secondary mb-1.5">
        虚拟路径（可选）
      </label>
      <input
        v-model="localConfig.file_path"
        type="text"
        placeholder="例如：projects/website"
        class="w-full px-3 py-2 text-sm border border-border-primary rounded-md focus:ring-2 focus:ring-primary focus:border-primary"
        style="color: var(--color-text-primary); background-color: var(--color-bg-elevated);"
      />
      <p class="text-xs text-text-tertiary mt-1">
        在 PixelPunk 中的存储路径（用于分类管理）
      </p>
    </div>

    <!-- 文件夹ID（可选） -->
    <div>
      <label class="block text-xs font-medium text-text-secondary mb-1.5">
        文件夹 ID（可选）
      </label>
      <input
        v-model="localConfig.folder_id"
        type="text"
        placeholder="输入文件夹 ID"
        class="w-full px-3 py-2 text-sm border border-border-primary rounded-md focus:ring-2 focus:ring-primary focus:border-primary"
        style="color: var(--color-text-primary); background-color: var(--color-bg-elevated);"
      />
      <p class="text-xs text-text-tertiary mt-1">
        指定目标文件夹（优先级高于虚拟路径）
      </p>
    </div>

    <!-- 输出说明 -->
    <div class="bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg p-3">
      <p class="text-xs font-medium text-green-800 dark:text-green-300 mb-2">
        ✅ 输出字段
      </p>
      <ul class="text-xs text-green-700 dark:text-green-400 space-y-1">
        <li><code class="px-1 bg-green-100 dark:bg-green-800 rounded">url</code> - 图片 CDN 地址</li>
        <li><code class="px-1 bg-green-100 dark:bg-green-800 rounded">thumb_url</code> - 缩略图地址</li>
        <li><code class="px-1 bg-green-100 dark:bg-green-800 rounded">id</code> - 图片唯一ID</li>
        <li><code class="px-1 bg-green-100 dark:bg-green-800 rounded">width / height</code> - 图片尺寸</li>
        <li><code class="px-1 bg-green-100 dark:bg-green-800 rounded">format</code> - 图片格式</li>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import BaseSelect from '@/components/BaseSelect'

interface Props {
  config: Record<string, any>
}

interface Emits {
  (e: 'update:config', value: Record<string, any>): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 下拉框选项
const accessLevelOptions = [
  { label: 'Public（公开）', value: 'public' },
  { label: 'Private（私有）', value: 'private' },
  { label: 'Protected（受保护）', value: 'protected' },
]

// 本地配置状态
const localConfig = ref({
  file: props.config.file || '{{external.image}}',
  access_level: props.config.access_level || 'public',
  optimize: props.config.optimize !== undefined ? props.config.optimize : true,
  file_path: props.config.file_path || '',
  folder_id: props.config.folder_id || '',
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

<style scoped>
code {
  font-family: 'Monaco', 'Courier New', monospace;
  font-size: 0.85em;
}
</style>

