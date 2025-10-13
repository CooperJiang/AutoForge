<template>
  <Drawer
    :model-value="modelValue"
    title="配置工具参数"
    confirm-text="保存配置"
    cancel-text="取消"
    size="xl"
    @update:model-value="$emit('update:modelValue', $event)"
    @confirm="handleSave"
    @cancel="$emit('update:modelValue', false)"
  >
    <div v-if="toolCode === 'http_request'" class="space-y-4">
      <!-- Curl 粘贴提示 -->
      <div class="bg-primary-light border border-primary rounded-lg p-3 text-xs text-primary">
        💡 小提示：按 <kbd class="px-1.5 py-0.5 bg-bg-elevated border border-primary rounded">{{ isMac ? 'Cmd' : 'Ctrl' }}</kbd> + <kbd class="px-1.5 py-0.5 bg-bg-elevated border border-primary rounded">V</kbd> 可直接粘贴 cURL 命令自动解析
      </div>

      <div>
        <label class="block text-sm font-medium text-text-secondary mb-2">
          请求方式 <span class="text-red-500">*</span>
        </label>
        <BaseSelect
          v-model="localConfig.method"
          :options="methodOptions"
          required
        />
      </div>

      <div>
        <label class="block text-sm font-medium text-text-secondary mb-2">
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
        <label class="block text-sm font-medium text-text-secondary mb-2">
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
            class="w-full py-2 text-sm text-text-secondary border-2 border-dashed border-slate-300 rounded-lg hover:border-slate-400 hover:text-text-secondary transition-colors"
          >
            + 添加请求头
          </button>
        </div>
      </div>

      <!-- Params -->
      <div>
        <label class="block text-sm font-medium text-text-secondary mb-2">
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
            class="w-full py-2 text-sm text-text-secondary border-2 border-dashed border-slate-300 rounded-lg hover:border-slate-400 hover:text-text-secondary transition-colors"
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
          <label class="block text-sm font-medium text-text-secondary cursor-pointer">
            {{ bodyExpanded ? '▼' : '▶' }} 请求体 (Body) <span class="text-xs text-text-tertiary">(POST/PUT/PATCH)</span>
          </label>
        </button>
        <div v-show="bodyExpanded" class="space-y-1">
          <textarea
            v-model="localConfig.body"
            class="w-full px-3 py-2 border-2 border-border-primary rounded-lg focus:outline-none focus:border-primary bg-bg-elevated text-text-primary font-mono text-sm"
            rows="8"
            placeholder='{"key": "value"}'
          />
          <div class="text-xs text-text-tertiary">支持 JSON、文本等格式</div>
        </div>
      </div>
    </div>

    <!-- 邮件发送工具配置 -->
    <div v-else-if="toolCode === 'email_sender'" class="space-y-4">
      <div class="bg-primary-light border-l-4 border-primary p-3 mb-4">
        <p class="text-sm text-primary">
          <svg class="inline-block w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd"/>
          </svg>
          邮件发送使用系统配置，只需填写收件人和邮件内容
        </p>
      </div>

      <div>
        <label class="block text-sm font-medium text-text-secondary mb-2">
          收件人 <span class="text-red-500">*</span>
        </label>
        <BaseInput
          v-model="emailConfig.to"
          placeholder="recipient@example.com, another@example.com"
          required
        />
        <p class="text-xs text-text-tertiary mt-1">多个收件人用逗号分隔</p>
      </div>

      <div>
        <label class="block text-sm font-medium text-text-secondary mb-2">
          抄送人
        </label>
        <BaseInput
          v-model="emailConfig.cc"
          placeholder="cc@example.com"
        />
      </div>

      <div>
        <label class="block text-sm font-medium text-text-secondary mb-2">
          邮件主题 <span class="text-red-500">*</span>
        </label>
        <BaseInput
          v-model="emailConfig.subject"
          placeholder="定时任务执行通知"
          required
        />
      </div>

      <div>
        <label class="block text-sm font-medium text-text-secondary mb-2">
          邮件正文 <span class="text-red-500">*</span>
        </label>
        <textarea
          v-model="emailConfig.body"
          class="w-full px-3 py-2 border-2 border-border-primary rounded-lg focus:outline-none focus:border-primary bg-bg-elevated text-text-primary font-mono text-sm"
          rows="8"
          placeholder="尊敬的用户，您好！&#10;&#10;您正在使用【自动任务系统】进行身份验证，您的验证码为：&#10;&#10;      123456&#10;&#10;验证码有效期为 10 分钟，请勿泄露给他人。&#10;如非本人操作，请忽略此邮件。&#10;&#10;感谢您的使用！&#10;&#10;---&#10;【自动任务系统】&#10;support@yourdomain.com"
          required
        />
        <div class="space-y-1 mt-2">
          <p class="text-xs text-amber-600">
            💡 <strong>避免被拦截的建议：</strong>
          </p>
          <ul class="text-xs text-text-secondary ml-4 space-y-0.5">
            <li>• 使用完整的邮件格式（称呼、正文、签名）</li>
            <li>• 说明邮件来源和目的</li>
            <li>• 验证码邮件需包含有效期、安全提示</li>
            <li>• 避免纯数字或过于简短的内容</li>
          </ul>
        </div>
      </div>

      <div>
        <label class="block text-sm font-medium text-text-secondary mb-2">
          内容类型
        </label>
        <BaseSelect
          v-model="emailConfig.content_type"
          :options="contentTypeOptions"
        />
      </div>
    </div>

    <!-- 健康检查工具配置 -->
    <div v-else-if="toolCode === 'health_checker'" class="space-y-4">
      <div class="bg-primary-light border-l-4 border-primary p-3 mb-4">
        <p class="text-sm text-primary">
          <svg class="inline-block w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd"/>
          </svg>
          支持粘贴 cURL 命令自动填充配置（{{ isMac ? 'Cmd+V' : 'Ctrl+V' }}）
        </p>
      </div>

      <div>
        <label class="block text-sm font-medium text-text-secondary mb-2">
          检查 URL <span class="text-red-500">*</span>
        </label>
        <BaseInput
          v-model="healthConfig.url"
          placeholder="https://api.example.com/health"
          required
        />
      </div>

      <div>
        <label class="block text-sm font-medium text-text-secondary mb-2">
          请求方法
        </label>
        <BaseSelect
          v-model="healthConfig.method"
          :options="healthMethodOptions"
        />
      </div>

      <!-- Headers -->
      <div>
        <div class="flex items-center justify-between mb-2">
          <label class="block text-sm font-medium text-text-secondary">
            请求头 (Headers)
          </label>
          <button
            type="button"
            @click="addHealthHeader"
            class="text-xs text-emerald-600 hover:text-emerald-700"
          >
            + 添加
          </button>
        </div>
        <div class="space-y-2">
          <ParamInput
            v-for="(header, index) in healthHeaders"
            :key="index"
            :param="header"
            @update="(p) => updateHealthHeader(index, p)"
            @remove="removeHealthHeader(index)"
          />
        </div>
      </div>

      <!-- Body -->
      <div>
        <label class="block text-sm font-medium text-text-secondary mb-2">
          请求体 (Body)
        </label>
        <textarea
          v-model="healthBody"
          class="w-full px-3 py-2 border-2 border-border-primary rounded-lg focus:outline-none focus:border-primary bg-bg-elevated text-text-primary font-mono text-sm"
          :rows="bodyExpanded ? 12 : 4"
          placeholder='{"key": "value"}'
        />
        <div class="flex items-center justify-between mt-1">
          <p class="text-xs text-text-tertiary">支持 JSON 或纯文本</p>
          <button
            type="button"
            @click="bodyExpanded = !bodyExpanded"
            class="text-xs text-primary hover:text-primary"
          >
            {{ bodyExpanded ? '收起' : '展开' }}
          </button>
        </div>
      </div>

      <div>
        <label class="block text-sm font-medium text-text-secondary mb-2">
          超时时间（秒）
        </label>
        <BaseInput
          v-model.number="healthConfig.timeout"
          type="number"
          placeholder="10"
        />
      </div>

      <div>
        <label class="block text-sm font-medium text-text-secondary mb-2">
          期望状态码
        </label>
        <BaseInput
          v-model.number="healthConfig.expected_status"
          type="number"
          placeholder="200"
        />
        <p class="text-xs text-text-tertiary mt-1">设置为 0 表示任意 2xx 状态码</p>
      </div>

      <div>
        <label class="block text-sm font-medium text-text-secondary mb-2">
          期望内容
        </label>
        <BaseInput
          v-model="healthConfig.expected_content"
          placeholder="success"
        />
        <p class="text-xs text-text-tertiary mt-1">响应体中应包含的内容</p>
      </div>

      <div class="flex items-center gap-2">
        <input
          type="checkbox"
          v-model="healthConfig.use_regex"
          id="use_regex"
          class="rounded border-slate-300"
        />
        <label for="use_regex" class="text-sm text-text-secondary cursor-pointer">
          使用正则表达式匹配
        </label>
      </div>

      <div class="flex items-center gap-2">
        <input
          type="checkbox"
          v-model="healthConfig.check_ssl"
          id="check_ssl"
          class="rounded border-slate-300"
        />
        <label for="check_ssl" class="text-sm text-text-secondary cursor-pointer">
          检查 SSL 证书有效期
        </label>
      </div>

      <div v-if="healthConfig.check_ssl">
        <label class="block text-sm font-medium text-text-secondary mb-2">
          SSL 到期告警天数
        </label>
        <BaseInput
          v-model.number="healthConfig.ssl_warning_days"
          type="number"
          placeholder="30"
        />
      </div>

      <div class="flex items-center gap-2">
        <input
          type="checkbox"
          v-model="healthConfig.follow_redirects"
          id="follow_redirects"
          class="rounded border-slate-300"
        />
        <label for="follow_redirects" class="text-sm text-text-secondary cursor-pointer">
          跟随重定向
        </label>
      </div>

      <div class="flex items-center gap-2">
        <input
          type="checkbox"
          v-model="healthConfig.verify_ssl"
          id="verify_ssl"
          class="rounded border-slate-300"
        />
        <label for="verify_ssl" class="text-sm text-text-secondary cursor-pointer">
          验证 SSL 证书有效性
        </label>
      </div>
    </div>

    <div v-else class="text-center py-8 text-text-tertiary">
      该工具暂无需配置参数
    </div>
  </Drawer>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue'
import Drawer from '@/components/Drawer'
import BaseInput from '@/components/BaseInput'
import BaseSelect from '@/components/BaseSelect'
import ParamInput from '@/components/ParamInput'
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
  url?: string
  method?: string
  headers?: string
  body?: string
  [key: string]: any
}

interface EmailConfig {
  to: string
  cc: string
  subject: string
  body: string
  content_type: string
}

interface HealthConfig {
  url: string
  method: string
  headers: string
  body: string
  timeout: number
  expected_status: number
  expected_content: string
  use_regex: boolean
  check_ssl: boolean
  ssl_warning_days: number
  follow_redirects: boolean
  verify_ssl: boolean
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

const emailConfig = ref<EmailConfig>({
  to: '',
  cc: '',
  subject: '',
  body: '',
  content_type: 'text/plain'
})

const healthConfig = ref<HealthConfig>({
  url: '',
  method: 'GET',
  headers: '{}',
  body: '',
  timeout: 10,
  expected_status: 200,
  expected_content: '',
  use_regex: false,
  check_ssl: true,
  ssl_warning_days: 30,
  follow_redirects: true,
  verify_ssl: true
})

const healthHeaders = ref<Param[]>([])
const healthBody = ref('')

const methodOptions = [
  { label: 'GET', value: 'GET' },
  { label: 'POST', value: 'POST' },
  { label: 'PUT', value: 'PUT' },
  { label: 'DELETE', value: 'DELETE' },
  { label: 'PATCH', value: 'PATCH' }
]

const healthMethodOptions = [
  { label: 'GET', value: 'GET' },
  { label: 'POST', value: 'POST' },
  { label: 'HEAD', value: 'HEAD' },
  { label: 'PUT', value: 'PUT' },
  { label: 'DELETE', value: 'DELETE' },
  { label: 'PATCH', value: 'PATCH' }
]

const contentTypeOptions = [
  { label: '纯文本', value: 'text/plain' },
  { label: 'HTML', value: 'text/html' }
]

// 监听 props.config 变化，转换为本地格式
watch(() => props.config, (newConfig) => {
  if (newConfig) {
    if (props.toolCode === 'http_request') {
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
    } else if (props.toolCode === 'email_sender') {
      emailConfig.value = {
        to: newConfig.to || '',
        cc: newConfig.cc || '',
        subject: newConfig.subject || '',
        body: newConfig.body || '',
        content_type: newConfig.content_type || 'text/plain'
      }
    } else if (props.toolCode === 'health_checker') {
      healthConfig.value = {
        url: newConfig.url || '',
        method: newConfig.method || 'GET',
        headers: newConfig.headers || '{}',
        body: newConfig.body || '',
        timeout: newConfig.timeout || 10,
        expected_status: newConfig.expected_status || 200,
        expected_content: newConfig.expected_content || '',
        use_regex: newConfig.use_regex || false,
        check_ssl: newConfig.check_ssl !== undefined ? newConfig.check_ssl : true,
        ssl_warning_days: newConfig.ssl_warning_days || 30,
        follow_redirects: newConfig.follow_redirects !== undefined ? newConfig.follow_redirects : true,
        verify_ssl: newConfig.verify_ssl !== undefined ? newConfig.verify_ssl : true
      }

      // 解析 headers
      try {
        const headers = JSON.parse(newConfig.headers || '{}')
        healthHeaders.value = Object.entries(headers).map(([key, value]) => ({
          key,
          value: String(value)
        }))
      } catch {
        healthHeaders.value = []
      }

      // 设置 body
      healthBody.value = newConfig.body || ''
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

// Health Headers 操作
const addHealthHeader = () => {
  healthHeaders.value.push({ key: '', value: '' })
}

const updateHealthHeader = (index: number, param: Param) => {
  healthHeaders.value[index] = param
}

const removeHealthHeader = (index: number) => {
  healthHeaders.value.splice(index, 1)
}

// 同步 health headers 到 config
const syncHealthHeaders = () => {
  const headersObj: Record<string, string> = {}
  healthHeaders.value.forEach(h => {
    if (h.key) headersObj[h.key] = h.value
  })
  healthConfig.value.headers = JSON.stringify(headersObj)
}

// 同步 health body 到 config
const syncHealthBody = () => {
  healthConfig.value.body = healthBody.value
}

// 监听 health headers 变化
watch(healthHeaders, syncHealthHeaders, { deep: true })
watch(healthBody, syncHealthBody)

// cURL 粘贴解析
const handlePaste = (e: ClipboardEvent) => {
  if (!props.modelValue) return
  if (props.toolCode !== 'http_request' && props.toolCode !== 'health_checker') return

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

    if (props.toolCode === 'http_request') {
      localConfig.value = {
        url: parsed.url,
        method: parsed.method,
        headers: parsed.headers,
        params: parsed.params,
        body: formattedBody
      }

      // 如果有 body，自动展开
      if (formattedBody) {
        bodyExpanded.value = true
      }
    } else if (props.toolCode === 'health_checker') {
      // 更新 health checker 配置
      healthConfig.value.url = parsed.url
      healthConfig.value.method = parsed.method

      // 设置 headers
      healthHeaders.value = parsed.headers

      // 设置 body
      healthBody.value = formattedBody

      // 如果有 body，自动展开
      if (formattedBody) {
        bodyExpanded.value = true
      }
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
    syncConfig()
  } else if (props.toolCode === 'email_sender') {
    // 验证必填字段
    if (!emailConfig.value.to) {
      message.error('请输入收件人')
      return
    }
    if (!emailConfig.value.subject) {
      message.error('请输入邮件主题')
      return
    }
    if (!emailConfig.value.body) {
      message.error('请输入邮件正文')
      return
    }

    // 直接发送 emailConfig
    emit('update:config', emailConfig.value as any)
  } else if (props.toolCode === 'health_checker') {
    if (!healthConfig.value.url) {
      message.error('请输入检查URL')
      return
    }

    // 直接发送 healthConfig
    emit('update:config', healthConfig.value as any)
  }

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
