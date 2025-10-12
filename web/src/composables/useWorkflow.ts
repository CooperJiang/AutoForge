import { ref, computed } from 'vue'
import type { Workflow, WorkflowNode, WorkflowEdge, WorkflowEnvVar } from '@/types/workflow'

export function useWorkflow() {
  const workflow = ref<Workflow>({
    name: '',
    description: '',
    trigger: {
      type: 'manual'
    },
    nodes: [],
    edges: [],
    envVars: []
  })

  const nodes = computed(() => workflow.value.nodes)
  const edges = computed(() => workflow.value.edges)
  const envVars = computed(() => workflow.value.envVars || [])

  // 添加节点
  const addNode = (node: WorkflowNode) => {
    workflow.value.nodes.push(node)
  }

  // 更新节点
  const updateNode = (nodeId: string, updates: Partial<WorkflowNode>) => {
    const index = workflow.value.nodes.findIndex(n => n.id === nodeId)
    if (index !== -1) {
      workflow.value.nodes[index] = {
        ...workflow.value.nodes[index],
        ...updates
      }
    }
  }

  // 删除节点
  const deleteNode = (nodeId: string) => {
    workflow.value.nodes = workflow.value.nodes.filter(n => n.id !== nodeId)
    workflow.value.edges = workflow.value.edges.filter(
      e => e.source !== nodeId && e.target !== nodeId
    )
  }

  // 添加连线
  const addEdge = (edge: WorkflowEdge) => {
    workflow.value.edges.push(edge)
  }

  // 删除连线
  const deleteEdge = (edgeId: string) => {
    workflow.value.edges = workflow.value.edges.filter(e => e.id !== edgeId)
  }

  // 添加环境变量
  const addEnvVar = (envVar: WorkflowEnvVar) => {
    if (!workflow.value.envVars) {
      workflow.value.envVars = []
    }
    workflow.value.envVars.push(envVar)
  }

  // 更新环境变量
  const updateEnvVar = (key: string, updates: Partial<WorkflowEnvVar>) => {
    if (!workflow.value.envVars) return
    const index = workflow.value.envVars.findIndex(v => v.key === key)
    if (index !== -1) {
      workflow.value.envVars[index] = {
        ...workflow.value.envVars[index],
        ...updates
      }
    }
  }

  // 删除环境变量
  const deleteEnvVar = (key: string) => {
    if (!workflow.value.envVars) return
    workflow.value.envVars = workflow.value.envVars.filter(v => v.key !== key)
  }

  // 获取前置节点（用于变量引用）
  const getPreviousNodes = (currentNodeId: string): WorkflowNode[] => {
    const previousNodeIds = new Set<string>()

    // 找到所有指向当前节点的边
    const incomingEdges = workflow.value.edges.filter(e => e.target === currentNodeId)

    // 递归找到所有前置节点
    const findPrevious = (nodeId: string) => {
      const edges = workflow.value.edges.filter(e => e.target === nodeId)
      edges.forEach(edge => {
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

    return workflow.value.nodes.filter(n => previousNodeIds.has(n.id))
  }

  // 重置工作流
  const resetWorkflow = () => {
    workflow.value = {
      name: '',
      description: '',
      trigger: {
        type: 'manual'
      },
      nodes: [],
      edges: [],
      envVars: []
    }
  }

  // 加载工作流
  const loadWorkflow = (data: Workflow) => {
    workflow.value = data
  }

  // 导出工作流JSON
  const exportWorkflow = () => {
    return JSON.stringify(workflow.value, null, 2)
  }

  // 切换启用/禁用状态
  const toggleEnabled = () => {
    workflow.value.enabled = !workflow.value.enabled
  }

  // 验证工作流
  const validateWorkflow = () => {
    if (!workflow.value.name.trim()) {
      return { valid: false, message: '请输入工作流名称' }
    }
    if (workflow.value.nodes.length === 0) {
      return { valid: false, message: '工作流至少需要一个节点' }
    }
    return { valid: true, message: '' }
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
    validateWorkflow
  }
}
