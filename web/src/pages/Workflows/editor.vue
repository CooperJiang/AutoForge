<template>
  <div class="fixed inset-0 bg-bg-hover flex flex-col">
    <!-- 顶部工具栏 -->
    <div
      class="bg-bg-elevated border-b border-border-primary px-6 py-3 flex items-center justify-between flex-shrink-0"
    >
      <div class="flex items-center gap-4">
        <BaseButton size="sm" variant="ghost" @click="handleBack" class="shrink-0">
          <ArrowLeft class="w-4 h-4" />
        </BaseButton>
        <div
          class="input-wrapper flex items-center gap-2 px-2.5 py-1 rounded-md bg-bg-hover border border-border-primary hover:border-slate-300 transition-all duration-200"
        >
          <Workflow class="w-3.5 h-3.5 text-text-placeholder shrink-0" />
          <input
            v-model="workflow.name"
            type="text"
            placeholder="工作流名称"
            class="w-32 bg-transparent text-xs font-medium text-text-primary placeholder:text-text-placeholder"
            style="border: none; outline: none; box-shadow: none"
            @focus="$event.target.parentElement.classList.add('input-focused')"
            @blur="$event.target.parentElement.classList.remove('input-focused')"
          />
        </div>
      </div>
      <div class="flex items-center gap-3">
        <!-- 状态指示 -->
        <div
          :class="[
            'px-3 py-1.5 rounded-full text-xs font-medium flex items-center gap-1.5 transition-colors border',
            !workflow.id
              ? 'bg-slate-500/10 text-slate-600 dark:text-slate-400 border-slate-500/20 opacity-50 cursor-not-allowed'
              : workflow.enabled
                ? 'bg-green-500/10 text-green-600 dark:text-green-400 border-green-500/20 hover:bg-green-500/20 cursor-pointer'
                : 'bg-slate-500/10 text-slate-600 dark:text-slate-400 border-slate-500/20 hover:bg-slate-500/20 cursor-pointer',
          ]"
          @click="workflow.id && handleToggleEnabled()"
          :title="!workflow.id ? '请先保存工作流后才能启用/禁用' : ''"
        >
          <Power :class="['w-3.5 h-3.5', workflow.enabled && 'animate-pulse']" />
          {{ workflow.enabled ? '已启用' : '已禁用' }}
        </div>

        <div class="h-6 w-px bg-bg-tertiary"></div>

        <!-- 次要操作 -->
        <Tooltip text="API 设置" position="bottom">
          <BaseButton size="sm" variant="ghost" @click="showAPISettings = true">
            <Globe class="w-4 h-4" />
          </BaseButton>
        </Tooltip>
        <Tooltip text="环境变量配置" position="bottom">
          <BaseButton size="sm" variant="ghost" @click="showEnvVarManager = true">
            <Settings class="w-4 h-4" />
          </BaseButton>
        </Tooltip>
        <Tooltip text="导入工作流" position="bottom">
          <BaseButton size="sm" variant="ghost" @click="showImportDialog = true">
            <Upload class="w-4 h-4" />
          </BaseButton>
        </Tooltip>
        <Tooltip text="导出工作流" position="bottom">
          <BaseButton size="sm" variant="ghost" @click="showExportDialog = true">
            <Download class="w-4 h-4" />
          </BaseButton>
        </Tooltip>
        <Tooltip text="清除本地草稿" position="bottom">
          <BaseButton size="sm" variant="ghost" @click="handleClearDraft">
            <Trash2 class="w-4 h-4" />
          </BaseButton>
        </Tooltip>

        <div class="h-6 w-px bg-bg-tertiary"></div>

        <!-- 主要操作 -->
        <BaseButton
          size="sm"
          variant="secondary"
          @click="handleExecute"
          :disabled="!workflow.enabled || nodes.length === 0"
        >
          <Play class="w-4 h-4 mr-1.5" />
          执行
        </BaseButton>
        <BaseButton size="sm" @click="handleSave">
          <Save class="w-4 h-4 mr-1.5" />
          保存
        </BaseButton>
      </div>
    </div>

    <!-- 主要内容区 -->
    <div class="flex-1 flex overflow-hidden">
      <!-- 左侧工具面板 -->
      <ToolPanel @add-node="handleAddNode" />

      <!-- 画布区域 -->
      <div class="flex-1 relative" @drop="handleDrop" @dragover.prevent @dragenter.prevent>
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
          <Background variant="dots" pattern-color="#94a3b8" :gap="16" :size="1" />
          <Controls />

          <template #node-tool="{ data }">
            <ToolNode :data="data" @delete="handleNodeDeleteFromCanvas" />
          </template>

          <template #node-external_trigger="{ data }">
            <ExternalTriggerNode :data="data" @delete="handleNodeDeleteFromCanvas" />
          </template>

          <template #node-trigger="{ data }">
            <TriggerNode :data="data" @delete="handleNodeDeleteFromCanvas" />
          </template>

          <template #node-condition="{ data }">
            <ConditionNode :data="data" @delete="handleNodeDeleteFromCanvas" />
          </template>

          <template #node-delay="{ data }">
            <DelayNode :data="data" @delete="handleNodeDeleteFromCanvas" />
          </template>

          <template #node-switch="{ data }">
            <SwitchNode :data="data" @delete="handleNodeDeleteFromCanvas" />
          </template>
        </VueFlow>

        <!-- 空状态提示 -->
        <div
          v-if="nodes.length === 0"
          class="absolute inset-0 flex items-center justify-center pointer-events-none"
        >
          <div class="text-center text-text-placeholder">
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

    <!-- API 设置 -->
    <WorkflowAPISettings v-model="showAPISettings" :workflow="workflow" @refresh="loadWorkflow" />

    <!-- 环境变量管理 -->
    <EnvVarManager
      v-model="showEnvVarManager"
      :env-vars="envVars"
      @update:env-vars="handleUpdateEnvVars"
    />

    <!-- 导入对话框 -->
    <ImportExportDialog v-model="showImportDialog" mode="import" @import="handleImportData" />

    <!-- 导出对话框 -->
    <ImportExportDialog
      v-model="showExportDialog"
      mode="export"
      :workflow-data="exportWorkflowData"
    />

    <!-- 执行参数对话框 -->
    <ExecuteWithParamsDialog
      :visible="showExecuteDialog"
      :workflow="workflow"
      @close="showExecuteDialog = false"
      @execute="executeWorkflow"
    />

    <!-- 确认对话框 -->
    <ConfirmDialog
      v-model="confirmDialog.show"
      :title="confirmDialog.title"
      :message="confirmDialog.message"
      :confirm-text="confirmDialog.confirmText"
      :cancel-text="confirmDialog.cancelText"
      :variant="confirmDialog.variant"
      @confirm="confirmDialog.resolve?.(true)"
      @cancel="confirmDialog.resolve?.(false)"
    />

    <!-- 画布删除节点确认对话框 -->
    <Dialog
      v-model="showDeleteConfirm"
      title="删除节点"
      max-width="max-w-md"
      @confirm="confirmDeleteFromCanvas"
      @cancel="showDeleteConfirm = false"
    >
      <div class="py-2">
        <p class="text-text-primary">
          确定要删除节点 <span class="font-semibold text-primary">{{ deleteNodeName }}</span> 吗？
        </p>
        <p class="text-text-secondary text-sm mt-1">此操作无法撤销。</p>
      </div>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { useRouter, useRoute, onBeforeRouteLeave } from 'vue-router'
import { VueFlow, useVueFlow } from '@vue-flow/core'
import { Background } from '@vue-flow/background'
import { Controls } from '@vue-flow/controls'
import {
  ArrowLeft,
  Save,
  Play,
  Upload,
  Download,
  Workflow,
  Settings,
  Power,
  Trash2,
  Globe,
} from 'lucide-vue-next'
import BaseButton from '@/components/BaseButton'
import ConfirmDialog from '@/components/ConfirmDialog'
import Dialog from '@/components/Dialog'
import Tooltip from '@/components/Tooltip'
import ToolPanel from './components/ToolPanel.vue'
import ToolNode from './components/ToolNode.vue'
import ExternalTriggerNode from './components/ExternalTriggerNode.vue'
import TriggerNode from './components/TriggerNode.vue'
import ConditionNode from './components/ConditionNode.vue'
import DelayNode from './components/DelayNode.vue'
import SwitchNode from './components/SwitchNode.vue'
import NodeConfigDrawer from './components/NodeConfigDrawer.vue'
import WorkflowAPISettings from './components/WorkflowAPISettings.vue'
import EnvVarManager from './components/EnvVarManager.vue'
import ImportExportDialog from './components/ImportExportDialog.vue'
import ExecuteWithParamsDialog from './components/ExecuteWithParamsDialog.vue'
import { useWorkflow } from '@/composables/useWorkflow'
import { message } from '@/utils/message'
import type { WorkflowNode, WorkflowEnvVar } from '@/types/workflow'

// 确认对话框状态
interface ConfirmDialogState {
  show: boolean
  title?: string
  message: string
  confirmText?: string
  cancelText?: string
  variant?: 'info' | 'warning' | 'danger' | 'question'
  resolve?: (value: boolean) => void
}

const confirmDialog = ref<ConfirmDialogState>({
  show: false,
  message: '',
})

// 通用确认函数
const confirm = (options: Omit<ConfirmDialogState, 'show' | 'resolve'>): Promise<boolean> => {
  return new Promise((resolve) => {
    confirmDialog.value = {
      ...options,
      show: true,
      resolve,
    }
  })
}

// Import styles
import '@vue-flow/core/dist/style.css'
import '@vue-flow/core/dist/theme-default.css'
import '@vue-flow/controls/dist/style.css'

const router = useRouter()
const route = useRoute()
const workflowId = computed(() => route.params.id as string)
const {
  workflow,
  nodes,
  edges,
  envVars,
  addNode,
  updateNode,
  deleteNode,
  addEdge,
  deleteEdge,
  getPreviousNodes,
  toggleEnabled,
  validateWorkflow,
  loadWorkflow: loadWorkflowData,
  clearDraft,
} = useWorkflow(workflowId.value)
const { project } = useVueFlow()

// Vue Flow state
const vueFlowNodes = ref<any[]>([])
const vueFlowEdges = ref<any[]>([])
const showConfigDrawer = ref(false)
const showEnvVarManager = ref(false)
const showAPISettings = ref(false)
const showImportDialog = ref(false)
const showExportDialog = ref(false)
const showExecuteDialog = ref(false)
const selectedNode = ref<WorkflowNode | null>(null)

// 追踪是否有未保存的更改
const hasUnsavedChanges = ref(false)

// 监听数据变化，标记为有未保存的更改
watch(
  [nodes, edges, workflow, envVars],
  () => {
    hasUnsavedChanges.value = true
  },
  { deep: true }
)

// 将工作流节点转换为VueFlow节点
const syncToVueFlow = () => {
  vueFlowNodes.value = nodes.value.map((node) => ({
    id: node.id,
    type: node.type,
    position: node.position,
    data: {
      ...node,
      config: node.config || {},
      retry: node.retry || {
        enabled: false,
        maxRetries: 3,
        retryInterval: 1000,
        exponentialBackoff: false,
      },
    },
  }))

  vueFlowEdges.value = edges.value.map((edge) => {
    // 查找源节点
    const sourceNode = nodes.value.find((n) => n.id === edge.source)
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
      labelBgStyle: { fill: 'white', fillOpacity: 0.9 },
    }
  })
}

// 将VueFlow节点转换回工作流节点
watch(
  [vueFlowNodes, vueFlowEdges],
  () => {
    nodes.value = vueFlowNodes.value.map((node) => ({
      ...node.data,
      position: node.position,
    }))

    edges.value = vueFlowEdges.value.map((edge) => ({
      id: edge.id,
      source: edge.source,
      target: edge.target,
      sourceHandle: edge.sourceHandle,
      targetHandle: edge.targetHandle,
    }))
  },
  { deep: true }
)

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
    position: { x: 250, y: 100 + nodes.value.length * 100 },
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
      position,
    }

    addNode(newNode)
    syncToVueFlow()
    message.success(`已添加 ${toolName} 节点`)
  } catch (error: any) {
    console.error('Failed to parse drag data:', error)
  }
}

// 点击节点
const handleNodeClick = (event: any) => {
  selectedNode.value = event.node.data
  showConfigDrawer.value = true
}

// 点击连线
const handleEdgeClick = async (event: any) => {
  const confirmed = await confirm({
    title: '删除连接',
    message: '确定要删除这条连接吗？',
    variant: 'warning',
  })

  if (confirmed) {
    deleteEdge(event.edge.id)
    syncToVueFlow()
  }
}

// 连接节点
const handleConnect = (params: any) => {
  // 检查源节点是否是条件节点
  const sourceNode = nodes.value.find((n) => n.id === params.source)
  const isConditionNode = sourceNode?.type === 'condition'

  // 如果是条件节点，需要记录是从哪个分支出来的（true或false）
  const edgeLabel =
    isConditionNode && params.sourceHandle
      ? params.sourceHandle === 'true'
        ? 'True'
        : 'False'
      : undefined

  addEdge({
    id: `edge_${Date.now()}`,
    source: params.source,
    target: params.target,
    sourceHandle: params.sourceHandle,
    targetHandle: params.targetHandle,
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

// 从画布删除节点（带确认对话框）
const deleteNodeId = ref<string | null>(null)
const deleteNodeName = ref<string>('')
const showDeleteConfirm = ref(false)

const handleNodeDeleteFromCanvas = (nodeId: string) => {
  const node = nodes.value.find((n) => n.id === nodeId)
  deleteNodeId.value = nodeId
  deleteNodeName.value = node?.name || '未命名节点'
  showDeleteConfirm.value = true
}

const confirmDeleteFromCanvas = () => {
  if (deleteNodeId.value) {
    deleteNode(deleteNodeId.value)
    syncToVueFlow()
    message.success('节点已删除')
    deleteNodeId.value = null
    deleteNodeName.value = ''
    showDeleteConfirm.value = false
  }
}

// 更新环境变量
const handleUpdateEnvVars = (newEnvVars: WorkflowEnvVar[]) => {
  envVars.value = newEnvVars
}

// 返回
const handleBack = () => {
  router.push('/workflows')
}

// 保存
const handleSave = async () => {
  const validation = validateWorkflow()
  if (!validation.valid) {
    message.error(validation.message)
    return
  }

  try {
    const currentWorkflowId = workflowId.value

    // 从触发器节点中提取调度配置
    const triggerNode = nodes.value.find((n) => n.type === 'trigger')
    let scheduleType = ''
    let scheduleValue = ''

    if (triggerNode && triggerNode.config) {
      const config = triggerNode.config

      // 根据触发器配置构建调度信息
      if (config.scheduleType === 'interval' && config.scheduleValue) {
        scheduleType = 'interval'
        scheduleValue = String(config.scheduleValue) // 秒数
      } else if (config.scheduleType === 'daily' && config.scheduleValue) {
        scheduleType = 'daily'
        scheduleValue = config.scheduleValue // 直接使用 scheduleValue (HH:MM:SS)
      } else if (config.scheduleType === 'weekly' && config.scheduleValue) {
        scheduleType = 'weekly'
        scheduleValue = config.scheduleValue // 格式：day1,day2:HH:MM:SS
      } else if (config.scheduleType === 'monthly' && config.scheduleValue) {
        scheduleType = 'monthly'
        scheduleValue = config.scheduleValue // 格式：day:HH:MM:SS
      } else if (config.scheduleType === 'cron' && config.scheduleValue) {
        scheduleType = 'cron'
        scheduleValue = config.scheduleValue // cron 表达式
      }
    }

    // 去重 edges：根据 source 和 target 组合去重
    const uniqueEdges = edges.value.filter(
      (edge, index, self) =>
        index === self.findIndex((e) => e.source === edge.source && e.target === edge.target)
    )

    const workflowData = {
      name: workflow.value.name,
      description: workflow.value.description,
      nodes: nodes.value,
      edges: uniqueEdges,
      env_vars: envVars.value,
      schedule_type: scheduleType,
      schedule_value: scheduleValue,
      enabled: workflow.value.enabled,
    }

    if (currentWorkflowId && currentWorkflowId !== 'create') {
      // 更新已有工作流
      const { workflowApi } = await import('@/api/workflow')
      await workflowApi.update(currentWorkflowId, workflowData)
      message.success('工作流已更新')
      // 保存成功后清除草稿和未保存标记
      clearDraft()
      hasUnsavedChanges.value = false
    } else {
      // 创建新工作流
      const { workflowApi } = await import('@/api/workflow')
      const data = await workflowApi.create(workflowData)
      message.success('工作流已创建')
      // 清除创建页面的草稿和未保存标记
      clearDraft()
      hasUnsavedChanges.value = false
      // 跳转到编辑页面
      router.replace(`/workflows/${data.id}/edit`)
    }
  } catch (error: any) {
    console.error('Save workflow failed:', error)
    message.error(error.response?.data?.message || '保存失败')
  }
}

// 切换启用/禁用
const handleToggleEnabled = async () => {
  if (!workflow.value.id) {
    message.warning('请先保存工作流后才能启用/禁用')
    return
  }

  try {
    const newStatus = !workflow.value.enabled
    await workflowApi.toggleEnabled(workflow.value.id, newStatus)
    toggleEnabled()
    message.success(`工作流已${newStatus ? '启用' : '禁用'}`)
  } catch (error: any) {
    message.error(error.response?.data?.message || '操作失败')
  }
}

// 检查工作流是否需要外部参数
const needsExternalParams = () => {
  // 找到第一个节点（没有入边的节点）
  const firstNode = nodes.value.find((node) => {
    const hasIncomingEdge = edges.value.some((edge) => edge.target === node.id)
    return !hasIncomingEdge
  })

  // 如果第一个节点是 external_trigger 且有参数配置
  return firstNode?.type === 'external_trigger' && firstNode.config?.params?.length > 0
}

// 执行工作流
const handleExecute = async () => {
  if (!workflow.value.enabled) {
    message.warning('工作流已禁用，请先启用工作流')
    return
  }

  const validation = validateWorkflow()
  if (!validation.valid) {
    message.error(validation.message)
    return
  }

  const workflowId = route.params.id as string
  if (!workflowId || workflowId === 'create') {
    message.warning('请先保存工作流')
    return
  }

  // 如果需要参数，显示参数输入对话框
  if (needsExternalParams()) {
    showExecuteDialog.value = true
    return
  }

  // 直接执行
  await executeWorkflow({})
}

// 执行工作流（带参数）
const executeWorkflow = async (params: Record<string, any>) => {
  const workflowId = route.params.id as string

  try {
    message.info('正在执行工作流...')
    const { workflowApi } = await import('@/api/workflow')
    const data = await workflowApi.execute(workflowId, { params })
    message.success('工作流已开始执行')
    showExecuteDialog.value = false
    // 跳转到执行详情
    router.push(`/workflows/${workflowId}/executions/${data.execution_id}`)
  } catch (error: any) {
    console.error('Execute workflow failed:', error)
    message.error(error.message || '执行失败')
  }
}

// 清除草稿
const handleClearDraft = async () => {
  const confirmed = await confirm({
    title: '清除草稿',
    message: '确定要清除本地草稿吗？此操作不可恢复',
    variant: 'danger',
    confirmText: '清除',
    cancelText: '取消',
  })

  if (confirmed) {
    clearDraft()
    // 重新加载服务器数据
    if (workflowId.value && workflowId.value !== 'create') {
      loadWorkflow()
    } else {
      // 创建页面则重置为空
      nodes.value = []
      edges.value = []
      envVars.value = []
      workflow.value = {
        name: '',
        description: '',
        enabled: false,
      }
      syncToVueFlow()
    }
    message.success('草稿已清除')
  }
}

// 导出工作流数据（用于传递给导出对话框）
const exportWorkflowData = computed(() => {
  return {
    name: workflow.value.name,
    description: workflow.value.description,
    enabled: workflow.value.enabled,
    nodes: nodes.value,
    edges: edges.value,
    envVars: envVars.value,
  }
})

// 处理导入数据
const handleImportData = (data: any) => {
  // 导入数据
  workflow.value = {
    name: data.name || '导入的工作流',
    description: data.description || '',
    enabled: data.enabled !== undefined ? data.enabled : false,
  }
  nodes.value = data.nodes || []
  edges.value = data.edges || []
  envVars.value = data.env_vars || data.envVars || []

  // 同步到画布
  syncToVueFlow()

  // 标记为有未保存的更改
  hasUnsavedChanges.value = true

  message.success('工作流已导入')
}

// 加载工作流数据
const loadWorkflow = async () => {
  const currentWorkflowId = workflowId.value
  if (!currentWorkflowId || currentWorkflowId === 'create') {
    return
  }

  try {
    const { workflowApi } = await import('@/api/workflow')
    const data = await workflowApi.getById(currentWorkflowId)

    // 编辑模式：直接加载服务器数据，不提示
    // 清除草稿
    clearDraft()

    // 加载服务器数据
    loadWorkflowData(data)

    // 使用 nextTick 确保数据更新后再同步到 Vue Flow
    await nextTick()
    syncToVueFlow()

    // 加载完成后重置未保存标记
    hasUnsavedChanges.value = false
  } catch (error: any) {
    console.error('Load workflow failed:', error)
    message.error('加载工作流失败')
  }
}

// 路由守卫 - 阻止用户在有未保存更改时离开
onBeforeRouteLeave(async (to, from, next) => {
  if (hasUnsavedChanges.value) {
    const confirmed = await confirm({
      title: '离开编辑器',
      message: '有未保存的更改，确定要离开吗？（草稿将保留）',
      variant: 'warning',
      confirmText: '离开',
      cancelText: '继续编辑',
    })

    next(confirmed)
  } else {
    next()
  }
})

// 浏览器刷新/关闭提示
const handleBeforeUnload = (e: BeforeUnloadEvent) => {
  if (hasUnsavedChanges.value) {
    e.preventDefault()
    e.returnValue = ''
  }
}

// 初始化
onMounted(() => {
  loadWorkflow()
  window.addEventListener('beforeunload', handleBeforeUnload)
})

// 清理
onUnmounted(() => {
  window.removeEventListener('beforeunload', handleBeforeUnload)
})

syncToVueFlow()
</script>

<style scoped>
/* 输入框焦点样式 */
.input-wrapper input {
  border: none !important;
  outline: none !important;
  box-shadow: none !important;
}

.input-wrapper input:focus {
  border: none !important;
  outline: none !important;
  box-shadow: none !important;
}

.input-wrapper.input-focused {
  @apply border-slate-400 bg-bg-elevated shadow-sm;
}

:deep(.vue-flow__node) {
  cursor: pointer;
}

:deep(.vue-flow__edge-path) {
  stroke-width: 2;
}

/* VueFlow Controls 暗色模式适配 */
:deep(.vue-flow__controls) {
  @apply bg-bg-elevated border-2 border-border-primary rounded-lg shadow-lg;
}

:deep(.vue-flow__controls-button) {
  @apply bg-bg-elevated border border-border-primary;
  @apply hover:bg-bg-hover transition-colors duration-200;
}

:deep(.vue-flow__controls-button):hover {
  @apply bg-bg-hover;
}

:deep(.vue-flow__controls-button) svg {
  @apply fill-text-primary stroke-text-primary;
  fill: currentColor !important;
  stroke: currentColor !important;
  color: var(--color-text-primary) !important;
}
</style>
