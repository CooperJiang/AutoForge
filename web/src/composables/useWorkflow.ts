import { ref, computed, watch } from 'vue'
import type { Workflow, WorkflowNode, WorkflowEdge, WorkflowEnvVar } from '@/types/workflow'
import SecureStorage from '@/utils/storage'

const WORKFLOW_DRAFT_KEY = 'workflow_draft'

export function useWorkflow(workflowId?: string) {
  const storageKey = workflowId ? `${WORKFLOW_DRAFT_KEY}_${workflowId}` : WORKFLOW_DRAFT_KEY

  // 尝试从本地存储恢复草稿
  const loadDraft = () => {
    const draft = SecureStorage.getItem<{
      workflow: Partial<Workflow>
      nodes: WorkflowNode[]
      edges: WorkflowEdge[]
      envVars: WorkflowEnvVar[]
    }>(storageKey, null, true)
    return draft
  }

  const draft = loadDraft()

  const workflow = ref<Partial<Workflow>>(draft?.workflow || {
    name: '',
    description: '',
    nodes: [],
    edges: [],
    env_vars: [],
    enabled: false
  })

  const nodes = ref<WorkflowNode[]>(draft?.nodes || [])
  const edges = ref<WorkflowEdge[]>(draft?.edges || [])
  const envVars = ref<WorkflowEnvVar[]>(draft?.envVars || [])

  // 自动保存到本地存储
  const saveDraft = () => {
    SecureStorage.setItem(storageKey, {
      workflow: workflow.value,
      nodes: nodes.value,
      edges: edges.value,
      envVars: envVars.value
    }, { encrypt: true })
  }

  // 监听变化自动保存
  watch([workflow, nodes, edges, envVars], () => {
    saveDraft()
  }, { deep: true })

  // 添加节点
  const addNode = (node: WorkflowNode) => {
    nodes.value.push(node)
  }

  // 更新节点
  const updateNode = (nodeId: string, updates: Partial<WorkflowNode>) => {
    const index = nodes.value.findIndex(n => n.id === nodeId)
    if (index !== -1) {
      nodes.value[index] = {
        ...nodes.value[index],
        ...updates
      }
    }
  }

  // 删除节点
  const deleteNode = (nodeId: string) => {
    nodes.value = nodes.value.filter(n => n.id !== nodeId)
    edges.value = edges.value.filter(
      e => e.source !== nodeId && e.target !== nodeId
    )
  }

  // 添加连线
  const addEdge = (edge: WorkflowEdge) => {
    edges.value.push(edge)
  }

  // 删除连线
  const deleteEdge = (edgeId: string) => {
    edges.value = edges.value.filter(e => e.id !== edgeId)
  }

  // 添加环境变量
  const addEnvVar = (envVar: WorkflowEnvVar) => {
    envVars.value.push(envVar)
  }

  // 更新环境变量
  const updateEnvVar = (key: string, updates: Partial<WorkflowEnvVar>) => {
    const index = envVars.value.findIndex(v => v.key === key)
    if (index !== -1) {
      envVars.value[index] = {
        ...envVars.value[index],
        ...updates
      }
    }
  }

  // 删除环境变量
  const deleteEnvVar = (key: string) => {
    envVars.value = envVars.value.filter(v => v.key !== key)
  }

  // 获取前置节点（用于变量引用）
  const getPreviousNodes = (currentNodeId: string): WorkflowNode[] => {
    const previousNodeIds = new Set<string>()

    // 找到所有指向当前节点的边
    const incomingEdges = edges.value.filter(e => e.target === currentNodeId)

    // 递归找到所有前置节点
    const findPrevious = (nodeId: string) => {
      const edgeList = edges.value.filter(e => e.target === nodeId)
      edgeList.forEach(edge => {
        if (!previousNodeIds.has(edge.source)) {
          previousNodeIds.add(edge.source)
          findPrevious(edge.source)
        }
      })
    }

    incomingEdges.forEach(edge => {
      previousNodeIds.add(edge.source)
      findPrevious(edge.source)
    })

    return nodes.value.filter(n => previousNodeIds.has(n.id))
  }

  // 重置工作流
  const resetWorkflow = () => {
    workflow.value = {
      name: '',
      description: '',
      nodes: [],
      edges: [],
      env_vars: [],
      enabled: false
    }
    nodes.value = []
    edges.value = []
    envVars.value = []
  }

  // 加载工作流
  const loadWorkflow = (data: Workflow) => {
    workflow.value = data
    nodes.value = data.nodes || []
    edges.value = data.edges || []
    envVars.value = data.env_vars || []
  }

  // 导出工作流JSON
  const exportWorkflow = () => {
    return JSON.stringify({
      ...workflow.value,
      nodes: nodes.value,
      edges: edges.value,
      env_vars: envVars.value
    }, null, 2)
  }

  // 切换启用/禁用状态
  const toggleEnabled = () => {
    workflow.value.enabled = !workflow.value.enabled
  }

  // 验证工作流
  const validateWorkflow = () => {
    if (!workflow.value.name?.trim()) {
      return { valid: false, message: '请输入工作流名称' }
    }
    if (nodes.value.length === 0) {
      return { valid: false, message: '工作流至少需要一个节点' }
    }
    return { valid: true, message: '' }
  }

  // 清除草稿
  const clearDraft = () => {
    SecureStorage.removeItem(storageKey)
  }

  return {
    workflow,
    nodes,
    edges,
    envVars,
    addNode,
    updateNode,
    deleteNode,
    addEdge,
    deleteEdge,
    addEnvVar,
    updateEnvVar,
    deleteEnvVar,
    getPreviousNodes,
    toggleEnabled,
    resetWorkflow,
    loadWorkflow,
    exportWorkflow,
    validateWorkflow,
    saveDraft,
    clearDraft
  }
}
