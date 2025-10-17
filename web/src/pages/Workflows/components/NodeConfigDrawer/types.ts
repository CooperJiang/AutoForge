import type { WorkflowNode, WorkflowEnvVar } from '@/types/workflow'

/**
 * NodeConfigDrawer 类型定义
 */

export interface NodeConfigDrawerProps {
  modelValue: boolean
  node: WorkflowNode | null
  previousNodes?: WorkflowNode[]
  envVars?: WorkflowEnvVar[]
}

export interface NodeConfigDrawerEmits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'update', node: WorkflowNode): void
  (e: 'delete', nodeId: string): void
}

export interface Param {
  key: string
  value: string
}

export interface HttpRequestConfig {
  method: string
  url: string
  headers: Param[]
  params: Param[]
  body: string
}

export interface SelectOption {
  label: string
  value: string
}

export interface FormattedEnvVar {
  name: string
  value: string
}

