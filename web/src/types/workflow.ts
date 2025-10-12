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
  id?: string
  name: string
  description: string
  trigger: WorkflowTrigger
  nodes: WorkflowNode[]
  edges: WorkflowEdge[]
  envVars?: WorkflowEnvVar[] // 环境变量
  enabled?: boolean
  created_at?: string
  updated_at?: string
}

export interface WorkflowExecution {
  id: string
  workflowId: string
  status: 'running' | 'success' | 'failed' | 'cancelled'
  startTime: string
  endTime?: string
  trigger: {
    type: string
    data?: any
  }
  nodeExecutions: NodeExecution[]
  error?: string
}

export interface NodeExecution {
  nodeId: string
  nodeName: string
  status: 'pending' | 'running' | 'success' | 'failed' | 'skipped'
  startTime?: string
  endTime?: string
  input?: any
  output?: any
  error?: string
}
