/**
 * execution-detail 类型定义
 */

export interface WorkflowExecution {
  id: string
  workflow_id: string
  status: 'running' | 'success' | 'failed' | 'cancelled'
  trigger_type: string
  start_time: string
  end_time?: string
  duration_ms?: number
  node_results: NodeResult[]
  error?: string
}

export interface NodeResult {
  node_id: string
  node_name: string
  node_type: string
  status: 'pending' | 'running' | 'success' | 'failed' | 'skipped'
  start_time?: string
  end_time?: string
  duration_ms?: number
  input?: any
  output?: any
  error?: string
}

export interface ExecutionDetailProps {
  workflowId: string
  executionId: string
}

export type StatusType = 'running' | 'success' | 'failed' | 'cancelled'
export type NodeStatusType = 'pending' | 'running' | 'success' | 'failed' | 'skipped'

