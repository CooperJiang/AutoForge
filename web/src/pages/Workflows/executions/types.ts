/**
 * Executions 类型定义
 */

export interface WorkflowExecution {
  id: string
  workflow_id: string
  status: 'running' | 'success' | 'failed' | 'cancelled'
  trigger_type: 'schedule' | 'scheduled' | 'webhook' | 'manual'
  start_time: number
  end_time?: number
  duration_ms: number
  total_nodes: number
  success_nodes: number
  failed_nodes: number
  skipped_nodes: number
  error?: string
}

export interface SelectOption {
  label: string
  value: string
}

