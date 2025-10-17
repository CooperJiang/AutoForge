/**
 * FeishuBotConfig 类型定义
 */

export interface FeishuConfig {
  webhook_url: string
  sign_secret?: string
  msg_type: 'text' | 'post' | 'image' | 'interactive'
  content?: string
  title?: string
  post_content?: string
  image_url?: string
  card_template?: 'notification' | 'update' | 'alert' | 'custom'
  card_content?: string
  card_status?: 'info' | 'success' | 'warning' | 'error'
  card_fields?: string
  card_buttons?: string
  card_custom_json?: string
}

export interface FeishuBotConfigProps {
  config: FeishuConfig
  previousNodes?: any[]
  envVars?: any[]
}

export interface SelectOption {
  label: string
  value: string
}

