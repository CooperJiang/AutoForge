/**
 * ToolConfigDrawer 类型定义
 */

export interface Param {
  key: string
  value: string
}

export interface HttpRequestConfig {
  url: string
  method: string
  headers: Param[]
  params: Param[]
  body: string
}

export interface EmailSenderConfig {
  to: string
  cc: string
  subject: string
  body: string
  content_type: string
}

export interface HealthCheckerConfig {
  url: string
  method: string
  headers: string
  body: string
  timeout: number
  expected_status: number
  expected_text: string
}

export interface FeishuBotConfig {
  webhook_url: string
  msg_type: string
  content: string
}

export type ToolConfig = HttpRequestConfig | EmailSenderConfig | HealthCheckerConfig | FeishuBotConfig

export interface ToolConfigDrawerProps {
  modelValue: boolean
  toolCode: string
  currentConfig: any
}

export interface ToolConfigDrawerEmits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'save', config: any): void
}

export interface ParsedCurlResult {
  url: string
  method: string
  headers: Param[]
  params: Param[]
  body: string
}

export interface SelectOption {
  label: string
  value: string
}

