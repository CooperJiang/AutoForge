<template>
  <div class="space-y-4">
    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2">
        下载链接 <span class="text-red-500">*</span>
      </label>
      <VariableSelector
        v-model="localConfig.url"
        placeholder="https://example.com/file.png 或 {{变量}}"
        :previous-nodes="previousNodes"
        :env-vars="envVars"
      />
      <p class="text-xs text-text-tertiary mt-1">
        支持变量：
        <code v-pre>{{external.xxx}}</code>
        /
        <code v-pre>{{nodes.xxx.yyy}}</code>
      </p>
    </div>

    <div>
      <div class="flex items-center justify-between mb-2">
        <label class="block text-sm font-medium text-text-secondary"> 请求头 (Headers) </label>
        <button type="button" @click="addHeader" class="text-xs text-emerald-600 hover:text-emerald-700">
          + 添加
        </button>
      </div>
      <div class="space-y-2">
        <ParamInput
          v-for="(header, index) in localHeaders"
          :key="index"
          :param="header"
          @update:param="updateHeader(index, $event)"
          @remove="removeHeader(index)"
        />
      </div>
    </div>

    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2"> 自定义文件名（可选） </label>
      <BaseInput v-model="localConfig.filename" placeholder="如：my-image.png 或 {{变量}}" />
      <p class="text-xs text-text-tertiary mt-1">留空则从响应头或 URL 推断</p>
    </div>

    <div class="grid grid-cols-2 gap-3">
      <div>
        <label class="block text-sm font-medium text-text-secondary mb-2"> 超时时间（秒） </label>
        <BaseInput v-model.number="localConfig.timeout" type="number" min="1" placeholder="60" />
      </div>
      <div class="flex items-center gap-2 mt-6">
        <BaseCheckbox v-model="localConfig.verify_ssl" label="验证 SSL 证书" />
      </div>
    </div>

    <div class="flex items-center gap-2">
      <BaseCheckbox v-model="localConfig.follow_redirects" label="跟随重定向" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import BaseInput from '@/components/BaseInput'
import ParamInput from '@/components/ParamInput'
import VariableSelector from '@/components/VariableSelector'
import BaseCheckbox from '@/components/BaseCheckbox/index.vue'
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

const previousNodes = props.previousNodes
const envVars = props.envVars

const localConfig = ref({
  url: props.config.url || '',
  headers: props.config.headers || [],
  filename: props.config.filename || '',
  timeout: props.config.timeout ?? 60,
  verify_ssl: props.config.verify_ssl ?? true,
  follow_redirects: props.config.follow_redirects ?? true,
})

// UI 使用数组表单展示 headers（key/value）
const localHeaders = ref<any[]>(Array.isArray(localConfig.value.headers) ? localConfig.value.headers : [])

// 防止 props->localConfig 回填触发递归更新
const updatingFromProps = ref(false)

// 唯一 ID 用于 checkbox 绑定 label
const uniqueId = Math.random().toString(36).substring(7)

const addHeader = () => {
  localHeaders.value.push({ key: '', value: '' })
  syncHeaders()
}
const updateHeader = (index: number, param: any) => {
  localHeaders.value[index] = param
  syncHeaders()
}
const removeHeader = (index: number) => {
  localHeaders.value.splice(index, 1)
  syncHeaders()
}

const syncHeaders = () => {
  localConfig.value.headers = localHeaders.value
  if (!updatingFromProps.value) {
    emit('update:config', { ...localConfig.value })
  }
}

watch(
  localConfig,
  (v) => {
    if (!updatingFromProps.value) {
      emit('update:config', { ...v, headers: localHeaders.value })
    }
  },
  { deep: true }
)

watch(
  () => props.config,
  (cfg) => {
    updatingFromProps.value = true
    localConfig.value = {
      url: cfg?.url || '',
      headers: cfg?.headers || [],
      filename: cfg?.filename || '',
      timeout: cfg?.timeout ?? 60,
      verify_ssl: cfg?.verify_ssl ?? true,
      follow_redirects: cfg?.follow_redirects ?? true,
    }
    localHeaders.value = Array.isArray(cfg?.headers) ? cfg.headers as any[] : []
    setTimeout(() => {
      updatingFromProps.value = false
    }, 0)
  },
  { deep: true }
)
</script>
