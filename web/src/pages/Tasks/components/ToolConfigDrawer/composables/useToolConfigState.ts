import { ref, computed, watch, type Ref } from 'vue'
import type {
  HttpRequestConfig,
  EmailSenderConfig,
  HealthCheckerConfig,
  Param,
  SelectOption,
} from '../types'

/**
 * 工具配置状态管理
 * @param toolCode 工具代码
 * @param initialConfig 初始配置
 */
export function useToolConfigState(toolCode: Ref<string>, initialConfig: Ref<any>) {
  // UI 状态
  const bodyExpanded = ref(false)
  const isMac = computed(() => /Mac/.test(navigator.userAgent))

  // HTTP 请求配置
  const httpConfig = ref<HttpRequestConfig>({
    url: '',
    method: 'GET',
    headers: [],
    params: [],
    body: '',
  })

  // 邮件发送配置
  const emailConfig = ref<EmailSenderConfig>({
    to: '',
    cc: '',
    subject: '',
    body: '',
    content_type: 'text/plain',
  })

  // 健康检查配置
  const healthConfig = ref<HealthCheckerConfig>({
    url: '',
    method: 'GET',
    headers: '{}',
    body: '',
    timeout: 10,
    expected_status: 200,
    expected_text: '',
  })

  const healthHeaders = ref<Param[]>([])
  const healthBody = ref('')

  // 飞书机器人配置
  const feishuConfig = ref<Record<string, any>>({
    webhook_url: '',
    sign_secret: '',
    msg_type: 'text',
    content: '',
    title: '',
    post_content: '',
    image_url: '',
    card_template: 'notification',
    card_content: '',
    card_status: 'info',
    card_fields: '',
    card_buttons: '',
    card_custom_json: '',
  })

  // 选项列表
  const methodOptions: SelectOption[] = [
    { label: 'GET', value: 'GET' },
    { label: 'POST', value: 'POST' },
    { label: 'PUT', value: 'PUT' },
    { label: 'DELETE', value: 'DELETE' },
    { label: 'PATCH', value: 'PATCH' },
  ]

  const healthMethodOptions: SelectOption[] = [
    { label: 'GET', value: 'GET' },
    { label: 'POST', value: 'POST' },
    { label: 'HEAD', value: 'HEAD' },
    { label: 'PUT', value: 'PUT' },
    { label: 'DELETE', value: 'DELETE' },
    { label: 'PATCH', value: 'PATCH' },
  ]

  const contentTypeOptions: SelectOption[] = [
    { label: '纯文本', value: 'text/plain' },
    { label: 'HTML', value: 'text/html' },
  ]

  /**
   * 初始化 HTTP 配置
   */
  const initHttpConfig = (config: any) => {
    if (!config) {
      httpConfig.value = {
        url: '',
        method: 'GET',
        headers: [],
        params: [],
        body: '',
      }
      return
    }

    try {
      const headers = JSON.parse(config.headers || '{}')
      const body = config.body ? JSON.parse(config.body) : {}

      httpConfig.value = {
        url: config.url || '',
        method: config.method || 'GET',
        headers: Object.entries(headers).map(([key, value]) => ({
          key,
          value: String(value),
        })),
        params: [],
        body:
          typeof body === 'object' && Object.keys(body).length > 0
            ? JSON.stringify(body, null, 2)
            : '',
      }
    } catch {
      httpConfig.value = {
        url: config.url || '',
        method: config.method || 'GET',
        headers: [],
        params: [],
        body: '',
      }
    }
  }

  /**
   * 初始化邮件配置
   */
  const initEmailConfig = (config: any) => {
    emailConfig.value = {
      to: config?.to || '',
      cc: config?.cc || '',
      subject: config?.subject || '',
      body: config?.body || '',
      content_type: config?.content_type || 'text/plain',
    }
  }

  /**
   * 初始化健康检查配置
   */
  const initHealthConfig = (config: any) => {
    healthConfig.value = {
      url: config?.url || '',
      method: config?.method || 'GET',
      headers: config?.headers || '{}',
      body: config?.body || '',
      timeout: config?.timeout || 10,
      expected_status: config?.expected_status || 200,
      expected_text: config?.expected_text || '',
    }

    // 初始化 headers
    try {
      const headers = JSON.parse(healthConfig.value.headers)
      healthHeaders.value = Object.entries(headers).map(([key, value]) => ({
        key,
        value: String(value),
      }))
    } catch {
      healthHeaders.value = []
    }

    healthBody.value = healthConfig.value.body
  }

  /**
   * 初始化飞书配置
   */
  const initFeishuConfig = (config: any) => {
    feishuConfig.value = {
      webhook_url: config?.webhook_url || '',
      sign_secret: config?.sign_secret || '',
      msg_type: config?.msg_type || 'text',
      content: config?.content || '',
      title: config?.title || '',
      post_content: config?.post_content || '',
      image_url: config?.image_url || '',
      card_template: config?.card_template || 'notification',
      card_content: config?.card_content || '',
      card_status: config?.card_status || 'info',
      card_fields: config?.card_fields || '',
      card_buttons: config?.card_buttons || '',
      card_custom_json: config?.card_custom_json || '',
    }
  }

  /**
   * 根据工具类型初始化配置
   */
  const initConfig = () => {
    const config = initialConfig.value

    switch (toolCode.value) {
      case 'http_request':
        initHttpConfig(config)
        break
      case 'email_sender':
        initEmailConfig(config)
        break
      case 'health_checker':
        initHealthConfig(config)
        break
      case 'feishu_bot':
        initFeishuConfig(config)
        break
    }
  }

  // 监听配置变化
  watch(
    initialConfig,
    () => {
      initConfig()
    },
    { immediate: true, deep: true }
  )

  return {
    // UI 状态
    bodyExpanded,
    isMac,

    // 配置状态
    httpConfig,
    emailConfig,
    healthConfig,
    healthHeaders,
    healthBody,
    feishuConfig,

    // 选项列表
    methodOptions,
    healthMethodOptions,
    contentTypeOptions,

    // 方法
    initConfig,
  }
}

