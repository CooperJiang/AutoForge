<template>
  <Dialog
    :model-value="modelValue"
    title="配置工具参数"
    confirm-text="保存配置"
    @update:model-value="$emit('update:modelValue', $event)"
    @confirm="handleSave"
  >
    <div v-if="toolCode === 'http_request'" class="space-y-4">
      <!-- Curl 粘贴提示 -->
      <div class="bg-blue-50 border border-blue-200 rounded-lg p-3 text-xs text-blue-800">
        💡 小提示：按 <kbd class="px-1.5 py-0.5 bg-white border border-blue-300 rounded">{{ isMac ? 'Cmd' : 'Ctrl' }}</kbd> + <kbd class="px-1.5 py-0.5 bg-white border border-blue-300 rounded">V</kbd> 可直接粘贴 cURL 命令自动解析
      </div>

      <div>
        <label class="block text-sm font-medium text-slate-700 mb-2">
          请求方式 <span class="text-red-500">*</span>
        </label>
        <BaseSelect
          v-model="localConfig.method"
          :options="methodOptions"
          required
        />
      </div>

      <div>
        <label class="block text-sm font-medium text-slate-700 mb-2">
          接口地址 <span class="text-red-500">*</span>
        </label>
        <BaseInput
          v-model="localConfig.url"
          placeholder="https://api.example.com/checkin"
          required
        />
      </div>

      <!-- Headers -->
      <div>
        <label class="block text-sm font-medium text-slate-700 mb-2">
          请求头（可选）
        </label>
        <div class="space-y-2">
          <ParamInput
            v-for="(header, index) in localConfig.headers"
            :key="index"
            :param="header"
            key-placeholder="Header名称"
            value-placeholder="Header值"
            @update:param="updateHeader(index, $event)"
            @remove="removeHeader(index)"
          />
          <button
            type="button"
            @click="addHeader"
            class="w-full py-2 text-sm text-slate-600 border-2 border-dashed border-slate-300 rounded-lg hover:border-slate-400 hover:text-slate-700 transition-colors"
          >
            + 添加请求头
          </button>
        </div>
      </div>

      <!-- Params -->
      <div>
        <label class="block text-sm font-medium text-slate-700 mb-2">
          请求参数（可选）
        </label>
        <div class="space-y-2">
          <ParamInput
            v-for="(param, index) in localConfig.params"
            :key="index"
            :param="param"
            key-placeholder="参数名"
            value-placeholder="参数值"
            @update:param="updateParam(index, $event)"
            @remove="removeParam(index)"
          />
          <button
            type="button"
            @click="addParam"
            class="w-full py-2 text-sm text-slate-600 border-2 border-dashed border-slate-300 rounded-lg hover:border-slate-400 hover:text-slate-700 transition-colors"
          >
            + 添加参数
          </button>
        </div>
      </div>

      <!-- Body -->
      <div>
        <button
          type="button"
          @click="bodyExpanded = !bodyExpanded"
          class="flex items-center justify-between w-full mb-2 text-left"
        >
          <label class="block text-sm font-medium text-slate-700 cursor-pointer">
            {{ bodyExpanded ? '▼' : '▶' }} 请求体 (Body) <span class="text-xs text-slate-500">(POST/PUT/PATCH)</span>
          </label>
        </button>
        <div v-show="bodyExpanded" class="space-y-1">
          <textarea
            v-model="localConfig.body"
            class="w-full px-3 py-2 border-2 border-slate-200 rounded-lg focus:outline-none focus:border-green-500 font-mono text-sm"
            rows="6"
            placeholder='{"key": "value"}'
          />
          <div class="text-xs text-slate-500">支持 JSON、文本等格式</div>
        </div>
      </div>
    </div>

    <div v-else class="text-center py-8 text-slate-500">
      该工具暂无需配置参数
    </div>
  </Dialog>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue'
import Dialog from '@/components/Dialog.vue'
import BaseInput from '@/components/BaseInput.vue'
import BaseSelect from '@/components/BaseSelect.vue'
import ParamInput from '@/components/ParamInput.vue'
import { message } from '@/utils/message'
import { parseCurl } from '@/utils/curlParser'

interface Param {
  key: string
  value: string
}

interface LocalConfig {
  url: string
  method: string
  headers: Param[]
  params: Param[]
  body: string
}

interface ToolConfig {
  url: string
  method: string
  headers: string
  body: string
}

const props = defineProps<{
  modelValue: boolean
  toolCode: string
  config: ToolConfig
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'update:config': [config: ToolConfig]
  save: []
}>()

const isMac = ref(/Mac/.test(navigator.userAgent))
const bodyExpanded = ref(false)

const localConfig = ref<LocalConfig>({
  url: '',
  method: 'GET',
  headers: [],
  params: [],
  body: ''
})

const methodOptions = [
  { label: 'GET', value: 'GET' },
  { label: 'POST', value: 'POST' },
  { label: 'PUT', value: 'PUT' },
  { label: 'DELETE', value: 'DELETE' },
  { label: 'PATCH', value: 'PATCH' }
]

// 监听 props.config 变化，转换为本地格式
watch(() => props.config, (newConfig) => {
  if (newConfig) {
    try {
      const headers = JSON.parse(newConfig.headers || '{}')
      const body = JSON.parse(newConfig.body || '{}')

      localConfig.value = {
        url: newConfig.url || '',
        method: newConfig.method || 'GET',
        headers: Object.entries(headers).map(([key, value]) => ({
          key,
          value: String(value)
        })),
        params: [],
        body: typeof body === 'object' && Object.keys(body).length > 0
          ? JSON.stringify(body, null, 2)
          : ''
      }
    } catch {
      localConfig.value = {
        url: newConfig.url || '',
        method: newConfig.method || 'GET',
        headers: [],
        params: [],
        body: ''
      }
    }
  }
}, { immediate: true })

// 同步本地配置回父组件
const syncConfig = () => {
  const headersObj: Record<string, string> = {}
  localConfig.value.headers.forEach(h => {
    if (h.key) headersObj[h.key] = h.value
  })

  let bodyObj: any = {}
  if (localConfig.value.body) {
    try {
      bodyObj = JSON.parse(localConfig.value.body)
    } catch {
      // 如果不是JSON，保持原样
      bodyObj = localConfig.value.body
    }
  }

  emit('update:config', {
    url: localConfig.value.url,
    method: localConfig.value.method,
    headers: JSON.stringify(headersObj),
    body: typeof bodyObj === 'string' ? bodyObj : JSON.stringify(bodyObj)
  })
}

// Headers 操作
const addHeader = () => {
  localConfig.value.headers.push({ key: '', value: '' })
}

const updateHeader = (index: number, param: Param) => {
  localConfig.value.headers[index] = param
}

const removeHeader = (index: number) => {
  localConfig.value.headers.splice(index, 1)
}

// Params 操作
const addParam = () => {
  localConfig.value.params.push({ key: '', value: '' })
}

const updateParam = (index: number, param: Param) => {
  localConfig.value.params[index] = param
}

const removeParam = (index: number) => {
  localConfig.value.params.splice(index, 1)
}

// cURL 粘贴解析
const handlePaste = (e: ClipboardEvent) => {
  if (!props.modelValue || props.toolCode !== 'http_request') return

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

    localConfig.value = {
      url: parsed.url,
      method: parsed.method,
      headers: parsed.headers,
      params: parsed.params,
      body: formattedBody
    }
    message.success('cURL 命令解析成功')
  } else {
    message.error('cURL 命令解析失败')
  }
}

// 保存配置
const handleSave = () => {
  if (props.toolCode === 'http_request') {
    if (!localConfig.value.url) {
      message.error('请输入请求URL')
      return
    }
  }
  syncConfig()
  emit('save')
}

// 监听粘贴事件
onMounted(() => {
  window.addEventListener('paste', handlePaste)
})

onUnmounted(() => {
  window.removeEventListener('paste', handlePaste)
})
</script>
