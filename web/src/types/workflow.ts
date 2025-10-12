export type NodeType = 'trigger' | 'tool' | 'condition' | 'delay' | 'switch' | 'end'

export interface NodeRetryConfig {
  enabled: boolean
  maxRetries: number
  retryInterval: number
  exponentialBackoff: boolean
}

export interface WorkflowNode {
  id: string
  type: NodeType
  toolCode?: string
  name: string
  config: Record<string, any>
  retry?: NodeRetryConfig
  position: { x: number; y: number }
}

export interface WorkflowEdge {
  id: string
  source: string
  target: string
  sourceHandle?: string
  targetHandle?: string
  condition?: string
}

export interface WorkflowTrigger {
  type: 'schedule' | 'manual' | 'webhook'
  scheduleType?: string
  scheduleValue?: string
  webhookPath?: string
  webhookMethod?: string
}

export interface WorkflowEnvVar {
  key: string
  value: string
  description?: string
  encrypted?: boolean
}

export interface Workflow {
  id: string
  user_id?: string
  name: string
  description: string
  nodes: WorkflowNode[]
  edges: WorkflowEdge[]
  env_vars?: WorkflowEnvVar[]
  schedule_type?: string
  schedule_value?: string
  enabled: boolean
  next_run_time?: number
  total_executions?: number
  success_count?: number
  failed_count?: number
  last_executed_at?: number
  created_at: number
  updated_at: number
}

export interface WorkflowExecution {
  id: string
  workflow_id: string
  user_id: string
  status: 'pending' | 'running' | 'success' | 'failed' | 'cancelled'
  trigger_type: string
  start_time?: number
  end_time?: number
  duration_ms: number
  total_nodes: number
  success_nodes: number
  failed_nodes: number
  skipped_nodes: number
  node_logs: NodeExecutionLog[]
  error?: string
  created_at: number
  updated_at: number
}

export interface NodeExecutionLog {
  node_id: string
  node_type: string
  node_name: string
  status: 'pending' | 'running' | 'success' | 'failed' | 'skipped'
  start_time?: number
  end_time?: number
  duration_ms: number
  retry_count: number
  output?: Record<string, any>
  error?: string
}

// DTO types for API requests
export interface CreateWorkflowDto {
  name: string
  description: string
  nodes: WorkflowNode[]
  edges: WorkflowEdge[]
  env_vars?: WorkflowEnvVar[]
  schedule_type?: string
  schedule_value?: string
  enabled?: boolean
}

export interface UpdateWorkflowDto {
  name?: string
  description?: string
  nodes?: WorkflowNode[]
  edges?: WorkflowEdge[]
  env_vars?: WorkflowEnvVar[]
  schedule_type?: string
  schedule_value?: string
  enabled?: boolean
}

export interface ExecuteWorkflowDto {
  env_vars?: Record<string, string>
}

// Same as WorkflowExecution for now
export type WorkflowExecutionDetail = WorkflowExecution
