<template>
  <div class="space-y-4">
    <!-- 提示 -->
    <div class="bg-primary-light border border-primary rounded-lg p-3 text-xs text-primary">
      💡 小提示：按
      <kbd class="px-1.5 py-0.5 bg-bg-elevated border border-primary rounded">{{
        isMac ? 'Cmd' : 'Ctrl'
      }}</kbd>
      + <kbd class="px-1.5 py-0.5 bg-bg-elevated border border-primary rounded">V</kbd> 可直接粘贴
      cURL 命令自动解析
    </div>

    <!-- 请求方式 -->
    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2">
        请求方式 <span class="text-red-500">*</span>
      </label>
      <BaseSelect v-model="config.method" :options="methodOptions" required />
    </div>

    <!-- 接口地址 -->
    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2">
        接口地址 <span class="text-red-500">*</span>
      </label>
      <BaseInput
        v-model="config.url"
        placeholder="https://api.example.com/checkin"
        required
      />
    </div>

    <!-- 请求头 -->
    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2"> 请求头（可选） </label>
      <div class="space-y-2">
        <ParamInput
          v-for="(header, index) in config.headers"
          :key="index"
          :param="header"
          key-placeholder="Header名称"
          value-placeholder="Header值"
          @update:param="$emit('update-header', index, $event)"
          @remove="$emit('remove-header', index)"
        />
        <button
          type="button"
          @click="$emit('add-header')"
          class="w-full py-2 text-sm text-text-secondary border-2 border-dashed border-slate-300 rounded-lg hover:border-slate-400 hover:text-text-secondary transition-colors"
        >
          + 添加请求头
        </button>
      </div>
    </div>

    <!-- 请求参数 -->
    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2"> 请求参数（可选） </label>
      <div class="space-y-2">
        <ParamInput
          v-for="(param, index) in config.params"
          :key="index"
          :param="param"
          key-placeholder="参数名"
          value-placeholder="参数值"
          @update:param="$emit('update-param', index, $event)"
          @remove="$emit('remove-param', index)"
        />
        <button
          type="button"
          @click="$emit('add-param')"
          class="w-full py-2 text-sm text-text-secondary border-2 border-dashed border-slate-300 rounded-lg hover:border-slate-400 hover:text-text-secondary transition-colors"
        >
          + 添加参数
        </button>
      </div>
    </div>

    <!-- 请求体 -->
    <div>
      <button
        type="button"
        @click="$emit('toggle-body')"
        class="flex items-center justify-between w-full mb-2 text-left"
      >
        <label class="block text-sm font-medium text-text-secondary cursor-pointer">
          {{ bodyExpanded ? '▼' : '▶' }} 请求体 (Body)
          <span class="text-xs text-text-tertiary">(POST/PUT/PATCH)</span>
        </label>
      </button>
      <div v-show="bodyExpanded" class="space-y-1">
        <textarea
          v-model="config.body"
          class="w-full px-3 py-2 border-2 border-border-primary rounded-lg focus:outline-none focus:border-primary bg-bg-elevated text-text-primary font-mono text-sm"
          rows="8"
          placeholder='{"key": "value"}'
        />
        <div class="text-xs text-text-tertiary">支持 JSON、文本等格式</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import BaseInput from '@/components/BaseInput'
import BaseSelect from '@/components/BaseSelect'
import ParamInput from '@/pages/Workflows/components/ParamInput.vue'
import type { HttpRequestConfig, SelectOption } from './types'

interface Props {
  config: HttpRequestConfig
  methodOptions: SelectOption[]
  bodyExpanded: boolean
  isMac: boolean
}

defineProps<Props>()

defineEmits<{
  'add-header': []
  'update-header': [index: number, param: { key: string; value: string }]
  'remove-header': [index: number]
  'add-param': []
  'update-param': [index: number, param: { key: string; value: string }]
  'remove-param': [index: number]
  'toggle-body': []
}>()
</script>

