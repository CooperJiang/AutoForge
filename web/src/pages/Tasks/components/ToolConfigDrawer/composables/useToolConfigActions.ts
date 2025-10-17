import { watch, onMounted, onUnmounted, type Ref } from 'vue'
import { message } from '@/utils/message'
import { parseCurl } from '@/utils/curlParser'
import type {
  HttpRequestConfig,
  EmailSenderConfig,
  HealthCheckerConfig,
  Param,
} from '../types'

/**
 * 工具配置操作逻辑
 */
export function useToolConfigActions(options: {
  toolCode: Ref<string>
  httpConfig: Ref<HttpRequestConfig>
  emailConfig: Ref<EmailSenderConfig>
  healthConfig: Ref<HealthCheckerConfig>
  healthHeaders: Ref<Param[]>
  healthBody: Ref<string>
  feishuConfig: Ref<Record<string, any>>
  bodyExpanded: Ref<boolean>
  modelValue: Ref<boolean>
  emit: (event: string, ...args: any[]) => void
}) {
  const {
    toolCode,
    httpConfig,
    emailConfig,
    healthConfig,
    healthHeaders,
    healthBody,
    feishuConfig,
    bodyExpanded,
    modelValue,
    emit,
  } = options

  // ============ HTTP Request 操作 ============

  const addHeader = () => {
    httpConfig.value.headers.push({ key: '', value: '' })
  }

  const updateHeader = (index: number, param: Param) => {
    httpConfig.value.headers[index] = param
  }

  const removeHeader = (index: number) => {
    httpConfig.value.headers.splice(index, 1)
  }

  const addParam = () => {
    httpConfig.value.params.push({ key: '', value: '' })
  }

  const updateParam = (index: number, param: Param) => {
    httpConfig.value.params[index] = param
  }

  const removeParam = (index: number) => {
    httpConfig.value.params.splice(index, 1)
  }

  // ============ Health Checker 操作 ============

  const addHealthHeader = () => {
    healthHeaders.value.push({ key: '', value: '' })
  }

  const updateHealthHeader = (index: number, param: Param) => {
    healthHeaders.value[index] = param
  }

  const removeHealthHeader = (index: number) => {
    healthHeaders.value.splice(index, 1)
  }

  const syncHealthHeaders = () => {
    const headersObj: Record<string, string> = {}
    healthHeaders.value.forEach((h) => {
      if (h.key) headersObj[h.key] = h.value
    })
    healthConfig.value.headers = JSON.stringify(headersObj)
  }

  const syncHealthBody = () => {
    healthConfig.value.body = healthBody.value
  }

  // 监听健康检查的 headers 和 body 变化
  watch(healthHeaders, syncHealthHeaders, { deep: true })
  watch(healthBody, syncHealthBody)

  // ============ cURL 粘贴解析 ============

  const handlePaste = (e: ClipboardEvent) => {
    if (!modelValue.value) return
    if (toolCode.value !== 'http_request' && toolCode.value !== 'health_checker') return

    const text = e.clipboardData?.getData('text')
    if (!text || !text.trim().startsWith('curl')) return

    e.preventDefault()

    const parsed = parseCurl(text)
    if (parsed) {
      let formattedBody = parsed.body || ''
      if (formattedBody) {
        try {
          const bodyObj = JSON.parse(formattedBody)
          formattedBody = JSON.stringify(bodyObj, null, 2)
        } catch {
          // 如果不是 JSON，保持原样
        }
      }

      if (toolCode.value === 'http_request') {
        httpConfig.value = {
          url: parsed.url,
          method: parsed.method,
          headers: parsed.headers,
          params: parsed.params,
          body: formattedBody,
        }

        if (formattedBody) {
          bodyExpanded.value = true
        }
      } else if (toolCode.value === 'health_checker') {
        healthConfig.value.url = parsed.url
        healthConfig.value.method = parsed.method
        healthHeaders.value = parsed.headers
        healthBody.value = formattedBody

        if (formattedBody) {
          bodyExpanded.value = true
        }
      }

      message.success('cURL 命令解析成功')
    } else {
      message.error('cURL 命令解析失败')
    }
  }

  // ============ 配置同步 ============

  const syncHttpConfig = () => {
    const headersObj: Record<string, string> = {}
    httpConfig.value.headers.forEach((h) => {
      if (h.key) headersObj[h.key] = h.value
    })

    let bodyObj: any = {}
    if (httpConfig.value.body) {
      try {
        bodyObj = JSON.parse(httpConfig.value.body)
      } catch {
        bodyObj = httpConfig.value.body
      }
    }

    emit('update:config', {
      url: httpConfig.value.url,
      method: httpConfig.value.method,
      headers: JSON.stringify(headersObj),
      body: typeof bodyObj === 'string' ? bodyObj : JSON.stringify(bodyObj),
    })
  }

  // ============ 保存配置 ============

  const validateHttpConfig = (): boolean => {
    if (!httpConfig.value.url) {
      message.error('请输入请求URL')
      return false
    }
    return true
  }

  const validateEmailConfig = (): boolean => {
    if (!emailConfig.value.to) {
      message.error('请输入收件人')
      return false
    }
    if (!emailConfig.value.subject) {
      message.error('请输入邮件主题')
      return false
    }
    if (!emailConfig.value.body) {
      message.error('请输入邮件正文')
      return false
    }
    return true
  }

  const validateHealthConfig = (): boolean => {
    if (!healthConfig.value.url) {
      message.error('请输入检查URL')
      return false
    }
    return true
  }

  const validateFeishuConfig = (): boolean => {
    if (!feishuConfig.value.webhook_url) {
      message.error('请输入 Webhook URL')
      return false
    }

    const msgType = feishuConfig.value.msg_type
    if (msgType === 'text' && !feishuConfig.value.content) {
      message.error('请输入文本消息内容')
      return false
    }
    if (msgType === 'post' && !feishuConfig.value.post_content) {
      message.error('请输入富文本内容')
      return false
    }
    if (msgType === 'image' && !feishuConfig.value.image_url) {
      message.error('请输入图片 URL')
      return false
    }
    if (msgType === 'interactive') {
      if (feishuConfig.value.card_template === 'custom' && !feishuConfig.value.card_custom_json) {
        message.error('请输入自定义卡片 JSON')
        return false
      }
      if (feishuConfig.value.card_template !== 'custom' && !feishuConfig.value.title) {
        message.error('请输入卡片标题')
        return false
      }
    }
    return true
  }

  const handleSave = () => {
    let isValid = false

    switch (toolCode.value) {
      case 'http_request':
        isValid = validateHttpConfig()
        if (isValid) syncHttpConfig()
        break

      case 'email_sender':
        isValid = validateEmailConfig()
        if (isValid) emit('update:config', emailConfig.value as any)
        break

      case 'health_checker':
        isValid = validateHealthConfig()
        if (isValid) emit('update:config', healthConfig.value as any)
        break

      case 'feishu_bot':
        isValid = validateFeishuConfig()
        if (isValid) emit('update:config', feishuConfig.value as any)
        break

      default:
        isValid = true
    }

    if (isValid) {
      emit('save')
    }
  }

  // ============ 生命周期 ============

  onMounted(() => {
    window.addEventListener('paste', handlePaste)
  })

  onUnmounted(() => {
    window.removeEventListener('paste', handlePaste)
  })

  return {
    // HTTP Request
    addHeader,
    updateHeader,
    removeHeader,
    addParam,
    updateParam,
    removeParam,
    syncHttpConfig,

    // Health Checker
    addHealthHeader,
    updateHealthHeader,
    removeHealthHeader,
    syncHealthHeaders,
    syncHealthBody,

    // cURL 解析
    handlePaste,

    // 保存
    handleSave,
  }
}

