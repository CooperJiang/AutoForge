import { ref, watch } from 'vue'
import type { Workflow, WorkflowNode, WorkflowEdge, WorkflowEnvVar } from '@/types/workflow'
import SecureStorage from '@/utils/storage'

const WORKFLOW_DRAFT_KEY = 'workflow_draft'

export function useWorkflow(workflowId?: string) {
  const storageKey = workflowId ? `${WORKFLOW_DRAFT_KEY}_${workflowId}` : WORKFLOW_DRAFT_KEY


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

  const workflow = ref<Partial<Workflow>>(
    draft?.workflow || {
      name: '',
      description: '',
      nodes: [],
      edges: [],
      env_vars: [],
      enabled: false,
    }
  )

  const nodes = ref<WorkflowNode[]>(draft?.nodes || [])
  const edges = ref<WorkflowEdge[]>(draft?.edges || [])
  const envVars = ref<WorkflowEnvVar[]>(draft?.envVars || [])


  const saveDraft = () => {
    SecureStorage.setItem(
      storageKey,
      {
        workflow: workflow.value,
        nodes: nodes.value,
        edges: edges.value,
        envVars: envVars.value,
      },
      { encrypt: true }
    )
  }


  watch(
    [workflow, nodes, edges, envVars],
    () => {
      saveDraft()
    },
    { deep: true }
  )


  const addNode = (node: WorkflowNode) => {
    nodes.value.push(node)
  }


  const updateNode = (nodeId: string, updates: Partial<WorkflowNode>) => {
    const index = nodes.value.findIndex((n) => n.id === nodeId)
    if (index !== -1) {
      nodes.value[index] = {
        ...nodes.value[index],
        ...updates,
      }
    }
  }


  const deleteNode = (nodeId: string) => {
    nodes.value = nodes.value.filter((n) => n.id !== nodeId)
    edges.value = edges.value.filter((e) => e.source !== nodeId && e.target !== nodeId)
  }


  const addEdge = (edge: WorkflowEdge) => {
    edges.value.push(edge)
  }


  const deleteEdge = (edgeId: string) => {
    edges.value = edges.value.filter((e) => e.id !== edgeId)
  }


  const addEnvVar = (envVar: WorkflowEnvVar) => {
    envVars.value.push(envVar)
  }


  const updateEnvVar = (key: string, updates: Partial<WorkflowEnvVar>) => {
    const index = envVars.value.findIndex((v) => v.key === key)
    if (index !== -1) {
      envVars.value[index] = {
        ...envVars.value[index],
        ...updates,
      }
    }
  }


  const deleteEnvVar = (key: string) => {
    envVars.value = envVars.value.filter((v) => v.key !== key)
  }


  const getPreviousNodes = (currentNodeId: string): WorkflowNode[] => {
    const previousNodeIds = new Set<string>()


    const incomingEdges = edges.value.filter((e) => e.target === currentNodeId)


    const findPrevious = (nodeId: string) => {
      const edgeList = edges.value.filter((e) => e.target === nodeId)
      edgeList.forEach((edge) => {
        if (!previousNodeIds.has(edge.source)) {
          previousNodeIds.add(edge.source)
          findPrevious(edge.source)
        }
      })
    }

    incomingEdges.forEach((edge) => {
      previousNodeIds.add(edge.source)
      findPrevious(edge.source)
    })

    return nodes.value.filter((n) => previousNodeIds.has(n.id))
  }


  const resetWorkflow = () => {
    workflow.value = {
      name: '',
      description: '',
      nodes: [],
      edges: [],
      env_vars: [],
      enabled: false,
    }
    nodes.value = []
    edges.value = []
    envVars.value = []
  }


  const loadWorkflow = (data: Workflow) => {
    workflow.value = data
    nodes.value = data.nodes || []
    edges.value = data.edges || []
    envVars.value = data.env_vars || []
  }


  const exportWorkflow = () => {
    return JSON.stringify(
      {
        ...workflow.value,
        nodes: nodes.value,
        edges: edges.value,
        env_vars: envVars.value,
      },
      null,
      2
    )
  }


  const toggleEnabled = () => {
    workflow.value.enabled = !workflow.value.enabled
  }


  const validateWorkflow = () => {
    if (!workflow.value.name?.trim()) {
      return { valid: false, message: '请输入工作流名称' }
    }
    if (nodes.value.length === 0) {
      return { valid: false, message: '工作流至少需要一个节点' }
    }


    const externalTriggerNode = nodes.value.find((n) => n.type === 'external_trigger')
    if (externalTriggerNode) {

      const targetNodes = new Set(edges.value.map((e) => e.target))
      const startNodes = nodes.value.filter((n) => !targetNodes.has(n.id))

      if (startNodes.length === 0) {
        return { valid: false, message: '工作流必须有起始节点' }
      }


      if (!startNodes.some((n) => n.id === externalTriggerNode.id)) {
        return { valid: false, message: '外部 API 触发节点必须是工作流的起始节点' }
      }


      if (externalTriggerNode.config?.params && externalTriggerNode.config.params.length > 0) {

        for (const param of externalTriggerNode.config.params) {
          if (!param.key?.trim()) {
            return { valid: false, message: '外部 API 触发节点的参数名称不能为空' }
          }
        }
      }
    }

    return { valid: true, message: '' }
  }


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
    clearDraft,
  }
}
