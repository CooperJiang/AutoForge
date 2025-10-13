<template>
  <div class="space-y-4">
    <div class="bg-primary-light border border-primary rounded-lg p-3 text-xs text-primary">
      💡 小提示：按 <kbd class="px-1.5 py-0.5 bg-bg-elevated border border-primary rounded">{{ isMac ? 'Cmd' : 'Ctrl' }}</kbd> + <kbd class="px-1.5 py-0.5 bg-bg-elevated border border-primary rounded">V</kbd> 可直接粘贴 cURL 命令自动解析
    </div>

    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2">
        检查 URL <span class="text-red-500">*</span>
      </label>
      <BaseInput
        v-model="localConfig.url"
        placeholder="https://api.example.com/health"
        @update:model-value="emitUpdate"
      />
    </div>

    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2">
        请求方法
      </label>
      <BaseSelect
        v-model="localConfig.method"
        :options="methodOptions"
        @update:model-value="emitUpdate"
      />
    </div>

    <div>
      <div class="flex items-center justify-between mb-2">
        <label class="block text-sm font-medium text-text-secondary">
          请求头 (Headers)
        </label>
        <button
          type="button"
          @click="addHeader"
          class="text-xs text-emerald-600 hover:text-emerald-700"
        >
          + 添加
        </button>
      </div>
      <div class="space-y-2">
        <ParamInput
          v-for="(header, index) in localConfig.headers"
          :key="index"
          :param="header"
          @update="(p) => updateHeader(index, p)"
          @remove="removeHeader(index)"
        />
      </div>
    </div>

    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2">
        请求体 (Body)
      </label>
      <textarea
        v-model="localConfig.body"
        @input="emitUpdate"
        class="w-full px-3 py-1.5 text-sm text-text-primary bg-bg-primary border-2 border-border-primary rounded-md transition-all duration-200 focus:border-border-focus focus:ring-2 focus:ring-primary-light focus:outline-none hover:border-border-secondary placeholder:text-text-placeholder font-mono"
        rows="4"
        placeholder='{"key": "value"}'
      />
      <p class="text-xs text-text-tertiary mt-1">支持 JSON 或纯文本</p>
    </div>

    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2">
        超时时间（秒）
      </label>
      <BaseInput
        v-model.number="localConfig.timeout"
        type="number"
        placeholder="10"
        @update:model-value="emitUpdate"
      />
    </div>

    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2">
        期望状态码
      </label>
      <BaseInput
        v-model.number="localConfig.expected_status"
        type="number"
        placeholder="200"
        @update:model-value="emitUpdate"
      />
    </div>

    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2">
        响应内容匹配（正则表达式，可选）
      </label>
      <BaseInput
        v-model="localConfig.response_pattern"
        placeholder="^success$"
        @update:model-value="emitUpdate"
      />
      <p class="text-xs text-text-tertiary mt-1">留空则不检查响应内容</p>
    </div>

    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2">
        SSL 证书检查（天数警告）
      </label>
      <BaseInput
        v-model.number="localConfig.ssl_expiry_days"
        type="number"
        placeholder="30"
        @update:model-value="emitUpdate"
      />
      <p class="text-xs text-text-tertiary mt-1">SSL证书剩余天数少于此值时发出警告</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue'
import BaseInput from '@/components/BaseInput'
import BaseSelect from '@/components/BaseSelect'
import ParamInput from '@/components/ParamInput'
import { message } from '@/utils/message'
import { parseCurl } from '@/utils/curlParser'

interface Props {
  config: Record<string, any>
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:config': [config: Record<string, any>]
}>()

const isMac = ref(/Mac/.test(navigator.userAgent))

const localConfig = ref({
  url: '',
  method: 'GET',
  headers: [],
  body: '',
  timeout: 10,
  expected_status: 200,
  response_pattern: '',
  ssl_expiry_days: 30,
  ...props.config
})

const methodOptions = [
  { label: 'GET', value: 'GET' },
  { label: 'POST', value: 'POST' },
  { label: 'PUT', value: 'PUT' },
  { label: 'DELETE', value: 'DELETE' }
]

watch(() => props.config, (newVal) => {
  localConfig.value = { ...localConfig.value, ...newVal }
}, { deep: true })

const addHeader = () => {
  if (!Array.isArray(localConfig.value.headers)) {
    localConfig.value.headers = []
  }
  localConfig.value.headers.push({ key: '', value: '', enabled: true })
  emitUpdate()
}

const updateHeader = (index: number, param: any) => {
  localConfig.value.headers[index] = param
  emitUpdate()
}

const removeHeader = (index: number) => {
  localConfig.value.headers.splice(index, 1)
  emitUpdate()
}

const emitUpdate = () => {
  emit('update:config', localConfig.value)
}

// cURL 粘贴解析
const handlePaste = (e: ClipboardEvent) => {
  const text = e.clipboardData?.getData('text')
  if (!text || !text.trim().startsWith('curl')) return

  e.preventDefault()

  const parsed = parseCurl(text)
  if (parsed) {
    // 格式化 body
    let formattedBody = parsed.body || ''
    if (formattedBody) {
      try {
        const bodyObj = JSON.parse(formattedBody)
        formattedBody = JSON.stringify(bodyObj, null, 2)
      } catch {
        // 如果不是 JSON，保持原样
      }
    }

    localConfig.value.url = parsed.url
    localConfig.value.method = parsed.method
    localConfig.value.headers = parsed.headers
    localConfig.value.body = formattedBody

    emitUpdate()
    message.success('cURL 命令解析成功')
  } else {
    message.error('cURL 命令解析失败')
  }
}

// 监听粘贴事件
onMounted(() => {
  window.addEventListener('paste', handlePaste)
})

onUnmounted(() => {
  window.removeEventListener('paste', handlePaste)
})
</script>
