/**
 * WorkflowAPISettings 类型定义
 */

export interface WorkflowAPISettingsProps {
  modelValue: boolean
  workflow: any
}

export interface WorkflowAPISettingsEmits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'refresh'): void
}

export interface Tab {
  key: string
  label: string
}

export interface CodeExample {
  label: string
  value: string
  code: string
}

