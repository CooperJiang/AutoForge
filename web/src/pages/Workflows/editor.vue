<template>
  <div class="fixed inset-0 bg-slate-50 flex flex-col">
    <!-- 顶部工具栏 -->
    <div class="bg-white border-b border-slate-200 px-4 py-3 flex items-center justify-between flex-shrink-0">
      <div class="flex items-center gap-3">
        <BaseButton size="sm" variant="ghost" @click="handleBack">
          <ArrowLeft class="w-4 h-4 mr-1" />
          返回
        </BaseButton>
        <div class="h-6 w-px bg-slate-200"></div>
        <BaseInput
          v-model="workflow.name"
          placeholder="工作流名称"
          class="w-64"
        />
      </div>
      <div class="flex items-center gap-2">
        <BaseButton
          size="sm"
          :variant="workflow.enabled ? 'success' : 'ghost'"
          @click="handleToggleEnabled"
        >
          <Power class="w-4 h-4 mr-1" />
          {{ workflow.enabled ? '已启用' : '已禁用' }}
        </BaseButton>
        <div class="h-6 w-px bg-slate-200"></div>
        <BaseButton size="sm" variant="ghost" @click="showEnvVarManager = true">
          <Settings class="w-4 h-4 mr-1" />
          环境变量
        </BaseButton>
        <BaseButton size="sm" variant="ghost" @click="handleExport">
          <FileJson class="w-4 h-4 mr-1" />
          导出JSON
        </BaseButton>
        <BaseButton
          size="sm"
          variant="secondary"
          @click="handleExecute"
          :disabled="!workflow.enabled || nodes.length === 0"
        >
          <Play class="w-4 h-4 mr-1" />
          执行
        </BaseButton>
        <BaseButton size="sm" @click="handleSave">
          <Save class="w-4 h-4 mr-1" />
          保存
        </BaseButton>
      </div>
    </div>

    <!-- 主要内容区 -->
    <div class="flex-1 flex overflow-hidden">
      <!-- 左侧工具面板 -->
      <ToolPanel @add-node="handleAddNode" />

      <!-- 画布区域 -->
      <div
        class="flex-1 relative"
        @drop="handleDrop"
        @dragover.prevent
        @dragenter.prevent
      >
        <VueFlow
          v-model:nodes="vueFlowNodes"
          v-model:edges="vueFlowEdges"
          :default-zoom="1"
          :min-zoom="0.2"
          :max-zoom="4"
          @node-click="handleNodeClick"
          @edge-click="handleEdgeClick"
          @connect="handleConnect"
        >
          <Background pattern-color="#e2e8f0" :gap="16" />
          <Controls />

          <template #node-tool="{ data }">
            <ToolNode :data="data" />
          </template>

          <template #node-trigger="{ data }">
            <TriggerNode :data="data" />
          </template>

          <template #node-condition="{ data }">
            <ConditionNode :data="data" />
          </template>

          <template #node-delay="{ data }">
            <DelayNode :data="data" />
          </template>

          <template #node-switch="{ data }">
            <SwitchNode :data="data" />
          </template>
        </VueFlow>

        <!-- 空状态提示 -->
        <div
          v-if="nodes.length === 0"
          class="absolute inset-0 flex items-center justify-center pointer-events-none"
        >
          <div class="text-center text-slate-400">
            <Workflow class="w-16 h-16 mx-auto mb-4 opacity-50" />
            <p class="text-lg font-medium mb-2">从左侧拖拽工具开始构建工作流</p>
            <p class="text-sm">或点击工具添加到画布</p>
          </div>
        </div>
      </div>
    </div>

    <!-- 节点配置抽屉 -->
    <NodeConfigDrawer
      v-model="showConfigDrawer"
      :node="selectedNode"
      :previous-nodes="selectedNode ? getPreviousNodes(selectedNode.id) : []"
      :env-vars="envVars"
      @update="handleUpdateNode"
      @delete="handleDeleteNode"
    />

    <!-- 环境变量管理 -->
    <EnvVarManager
      v-model="showEnvVarManager"
      :env-vars="envVars"
      @update:env-vars="handleUpdateEnvVars"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { VueFlow, useVueFlow } from '@vue-flow/core'
import { Background } from '@vue-flow/background'
import { Controls } from '@vue-flow/controls'
import { ArrowLeft, Save, Play, FileJson, Workflow, Settings, Power } from 'lucide-vue-next'
import BaseButton from '@/components/BaseButton'
import BaseInput from '@/components/BaseInput'
import ToolPanel from './components/ToolPanel.vue'
import ToolNode from './components/ToolNode.vue'
import TriggerNode from './components/TriggerNode.vue'
import ConditionNode from './components/ConditionNode.vue'
import DelayNode from './components/DelayNode.vue'
import SwitchNode from './components/SwitchNode.vue'
import NodeConfigDrawer from './components/NodeConfigDrawer.vue'
import EnvVarManager from './components/EnvVarManager.vue'
import { useWorkflow } from '@/composables/useWorkflow'
import { message } from '@/utils/message'
import type { WorkflowNode, WorkflowEnvVar } from '@/types/workflow'

// Import styles
import '@vue-flow/core/dist/style.css'
import '@vue-flow/core/dist/theme-default.css'
import '@vue-flow/controls/dist/style.css'

const router = useRouter()
const route = useRoute()
const { workflow, nodes, edges, envVars, addNode, updateNode, deleteNode, addEdge, deleteEdge, getPreviousNodes, toggleEnabled, validateWorkflow, exportWorkflow } = useWorkflow()
const { project } = useVueFlow()

// Vue Flow state
const vueFlowNodes = ref<any[]>([])
const vueFlowEdges = ref<any[]>([])
const showConfigDrawer = ref(false)
const showEnvVarManager = ref(false)
const selectedNode = ref<WorkflowNode | null>(null)
const dropZone = ref<HTMLElement | null>(null)

// 将工作流节点转换为VueFlow节点
const syncToVueFlow = () => {
  vueFlowNodes.value = nodes.value.map(node => ({
    id: node.id,
    type: node.type,
    position: node.position,
    data: node
  }))

  vueFlowEdges.value = edges.value.map(edge => {
    // 查找源节点
    const sourceNode = nodes.value.find(n => n.id === edge.source)
    const isConditionNode = sourceNode?.type === 'condition'

    // 根据sourceHandle确定标签和颜色
    let label = ''
    let style: any = {}

    if (isConditionNode && edge.sourceHandle) {
      if (edge.sourceHandle === 'true') {
        label = '✓ True'
        style = { stroke: '#10b981', strokeWidth: 2 }
      } else if (edge.sourceHandle === 'false') {
        label = '✗ False'
        style = { stroke: '#f43f5e', strokeWidth: 2 }
      }
    }

    return {
      id: edge.id,
      source: edge.source,
      target: edge.target,
      sourceHandle: edge.sourceHandle,
      targetHandle: edge.targetHandle,
      type: 'smoothstep',
      animated: true,
      label,
      style,
      labelStyle: { fill: style.stroke, fontWeight: 600, fontSize: 12 },
      labelBgStyle: { fill: 'white', fillOpacity: 0.9 }
    }
  })
}

// 将VueFlow节点转换回工作流节点
watch([vueFlowNodes, vueFlowEdges], () => {
  workflow.value.nodes = vueFlowNodes.value.map(node => ({
    ...node.data,
    position: node.position
  }))

  workflow.value.edges = vueFlowEdges.value.map(edge => ({
    id: edge.id,
    source: edge.source,
    target: edge.target
  }))
}, { deep: true })

// 添加节点
const handleAddNode = (toolCode: string, toolName: string, nodeType?: string) => {
  let type: 'trigger' | 'tool' | 'condition' | 'delay' | 'switch' = 'tool'

  if (nodeType) {
    type = nodeType as 'trigger' | 'tool' | 'condition' | 'delay' | 'switch'
  } else {
    // 向后兼容
    if (toolCode === 'trigger') type = 'trigger'
    else if (toolCode === 'condition') type = 'condition'
    else if (toolCode === 'delay') type = 'delay'
    else if (toolCode === 'switch') type = 'switch'
  }

  const newNode: WorkflowNode = {
    id: `node_${Date.now()}`,
    type,
    toolCode: type === 'tool' ? toolCode : undefined,
    name: toolName,
    config: {},
    position: { x: 250, y: 100 + nodes.value.length * 100 }
  }
  addNode(newNode)
  syncToVueFlow()
  message.success(`已添加 ${toolName} 节点`)
}

// 处理拖拽放置
const handleDrop = (event: DragEvent) => {
  event.preventDefault()

  const data = event.dataTransfer?.getData('application/vueflow')
  if (!data) return

  try {
    const { toolCode, toolName, nodeType } = JSON.parse(data)

    // 获取鼠标在画布上的位置
    const position = project({ x: event.clientX - 100, y: event.clientY - 50 })

    let type: 'trigger' | 'tool' | 'condition' | 'delay' | 'switch' = 'tool'
    if (nodeType) {
      type = nodeType as 'trigger' | 'tool' | 'condition' | 'delay' | 'switch'
    } else {
      // 向后兼容
      if (toolCode === 'trigger') type = 'trigger'
      else if (toolCode === 'condition') type = 'condition'
      else if (toolCode === 'delay') type = 'delay'
      else if (toolCode === 'switch') type = 'switch'
    }

    const newNode: WorkflowNode = {
      id: `node_${Date.now()}`,
      type,
      toolCode: type === 'tool' ? toolCode : undefined,
      name: toolName,
      config: {},
      position
    }

    addNode(newNode)
    syncToVueFlow()
    message.success(`已添加 ${toolName} 节点`)
  } catch (error) {
    console.error('Failed to parse drag data:', error)
  }
}

// 点击节点
const handleNodeClick = (event: any) => {
  selectedNode.value = event.node.data
  showConfigDrawer.value = true
}

// 点击连线
const handleEdgeClick = (event: any) => {
  if (confirm('是否删除此连接？')) {
    deleteEdge(event.edge.id)
    syncToVueFlow()
  }
}

// 连接节点
const handleConnect = (params: any) => {
  // 检查源节点是否是条件节点
  const sourceNode = nodes.value.find(n => n.id === params.source)
  const isConditionNode = sourceNode?.type === 'condition'

  // 如果是条件节点，需要记录是从哪个分支出来的（true或false）
  const edgeLabel = isConditionNode && params.sourceHandle ?
    (params.sourceHandle === 'true' ? 'True' : 'False') :
    undefined

  addEdge({
    id: `edge_${Date.now()}`,
    source: params.source,
    target: params.target,
    sourceHandle: params.sourceHandle,
    targetHandle: params.targetHandle
  })
  syncToVueFlow()
  message.success(`节点已连接${edgeLabel ? ` (${edgeLabel} 分支)` : ''}`)
}

// 更新节点配置
const handleUpdateNode = (nodeId: string, updates: Partial<WorkflowNode>) => {
  updateNode(nodeId, updates)
  syncToVueFlow()
  message.success('节点配置已更新')
}

// 删除节点
const handleDeleteNode = (nodeId: string) => {
  deleteNode(nodeId)
  syncToVueFlow()
  showConfigDrawer.value = false
  message.success('节点已删除')
}

// 更新环境变量
const handleUpdateEnvVars = (newEnvVars: WorkflowEnvVar[]) => {
  workflow.value.envVars = newEnvVars
}

// 返回
const handleBack = () => {
  if (nodes.value.length > 0) {
    if (confirm('有未保存的更改，确定要离开吗？')) {
      router.push('/workflows')
    }
  } else {
    router.push('/workflows')
  }
}

// 保存
const handleSave = () => {
  const validation = validateWorkflow()
  if (!validation.valid) {
    message.error(validation.message)
    return
  }

  // TODO: 调用API保存工作流
  console.log('Save workflow:', workflow.value)
  message.success('工作流已保存')
}

// 切换启用/禁用
const handleToggleEnabled = () => {
  toggleEnabled()
  const status = workflow.value.enabled ? '已启用' : '已禁用'
  message.success(`工作流${status}`)
}

// 执行工作流
const handleExecute = () => {
  if (!workflow.value.enabled) {
    message.warning('工作流已禁用，请先启用工作流')
    return
  }

  const validation = validateWorkflow()
  if (!validation.valid) {
    message.error(validation.message)
    return
  }

  // TODO: 调用API执行工作流
  message.info('正在执行工作流...')
  console.log('Execute workflow:', workflow.value)

  // 执行成功后跳转到执行历史
  setTimeout(() => {
    message.success('工作流执行成功')
    router.push(`/workflows/${route.params.id}/executions`)
  }, 1000)
}

// 导出JSON
const handleExport = () => {
  const json = exportWorkflow()
  const blob = new Blob([json], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${workflow.value.name || 'workflow'}.json`
  a.click()
  URL.revokeObjectURL(url)
  message.success('工作流JSON已导出')
}

// 初始化
syncToVueFlow()
</script>

<style scoped>
:deep(.vue-flow__node) {
  cursor: pointer;
}

:deep(.vue-flow__edge-path) {
  stroke-width: 2;
}
</style>
