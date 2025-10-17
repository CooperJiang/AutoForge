<template>
  <div class="space-y-4">
    <!-- 提示 -->
    <div class="bg-primary-light border-l-4 border-primary p-3 mb-4">
      <p class="text-sm text-primary">
        <svg class="inline-block w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
          <path
            fill-rule="evenodd"
            d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z"
            clip-rule="evenodd"
          />
        </svg>
        支持粘贴 cURL 命令自动填充配置（{{ isMac ? 'Cmd+V' : 'Ctrl+V' }}）
      </p>
    </div>

    <!-- 检查 URL -->
    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2">
        检查 URL <span class="text-red-500">*</span>
      </label>
      <BaseInput v-model="config.url" placeholder="https://api.example.com/health" required />
    </div>

    <!-- 请求方法 -->
    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2"> 请求方法 </label>
      <BaseSelect v-model="config.method" :options="methodOptions" />
    </div>

    <!-- 请求头 -->
    <div>
      <div class="flex items-center justify-between mb-2">
        <label class="block text-sm font-medium text-text-secondary"> 请求头 (Headers) </label>
        <button
          type="button"
          @click="$emit('add-header')"
          class="text-xs text-emerald-600 hover:text-emerald-700"
        >
          + 添加
        </button>
      </div>
      <div class="space-y-2">
        <ParamInput
          v-for="(header, index) in headers"
          :key="index"
          :param="header"
          @update:param="$emit('update-header', index, $event)"
          @remove="$emit('remove-header', index)"
        />
      </div>
    </div>

    <!-- 请求体 -->
    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2"> 请求体 (Body) </label>
      <textarea
        :value="body"
        @input="$emit('update:body', ($event.target as HTMLTextAreaElement).value)"
        class="w-full px-3 py-2 border-2 border-border-primary rounded-lg focus:outline-none focus:border-primary bg-bg-elevated text-text-primary font-mono text-sm"
        :rows="bodyExpanded ? 12 : 4"
        placeholder='{"key": "value"}'
      />
      <div class="flex items-center justify-between mt-1">
        <p class="text-xs text-text-tertiary">支持 JSON 或纯文本</p>
        <button
          type="button"
          @click="$emit('toggle-body')"
          class="text-xs text-primary hover:text-primary"
        >
          {{ bodyExpanded ? '收起' : '展开' }}
        </button>
      </div>
    </div>

    <!-- 超时时间 -->
    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2"> 超时时间（秒） </label>
      <BaseInput v-model.number="config.timeout" type="number" placeholder="10" />
    </div>

    <!-- 期望状态码 -->
    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2"> 期望状态码 </label>
      <BaseInput v-model.number="config.expected_status" type="number" placeholder="200" />
      <p class="text-xs text-text-tertiary mt-1">设置为 0 表示任意 2xx 状态码</p>
    </div>

    <!-- 期望内容 -->
    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2"> 期望内容 </label>
      <BaseInput v-model="config.expected_text" placeholder="success" />
      <p class="text-xs text-text-tertiary mt-1">响应体中应包含的内容</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import BaseInput from '@/components/BaseInput'
import BaseSelect from '@/components/BaseSelect'
import ParamInput from '@/pages/Workflows/components/ParamInput.vue'
import type { HealthCheckerConfig, Param, SelectOption } from './types'

interface Props {
  config: HealthCheckerConfig
  headers: Param[]
  body: string
  methodOptions: SelectOption[]
  bodyExpanded: boolean
  isMac: boolean
}

defineProps<Props>()

defineEmits<{
  'add-header': []
  'update-header': [index: number, param: Param]
  'remove-header': [index: number]
  'update:body': [value: string]
  'toggle-body': []
}>()
</script>

