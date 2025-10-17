<template>
  <div class="h-full flex flex-col">
    <!-- 顶部操作栏 -->
    <div class="flex-shrink-0 mb-4">
      <div class="bg-bg-elevated rounded-xl border border-border-primary p-4 shadow-sm">
        <div class="flex items-center justify-between gap-4">
          <!-- 左侧筛选区域 -->
          <div class="flex items-center gap-4">
            <!-- 搜索框 -->
            <BaseInput
              v-model="searchKeyword"
              placeholder="搜索工作流..."
              style="width: 300px"
              @keyup.enter="handleSearch"
            >
              <template #prefix>
                <Search class="w-4 h-4" />
              </template>
            </BaseInput>

            <!-- 状态筛选 -->
            <div class="flex items-center gap-2">
              <span class="text-xs" style="color: var(--color-text-secondary)">状态:</span>
              <RadioGroup
                v-model="filterStatus"
                :options="statusOptions"
                size="sm"
                @change="handleFilterChange"
              />
            </div>
          </div>

          <!-- 右侧按钮 -->
          <BaseButton size="md" @click="router.push('/workflows/create')">
            <Plus class="w-4 h-4 mr-1" />
            创建工作流
          </BaseButton>
        </div>
      </div>
    </div>

    <!-- 内容区域 - 带背景 -->
    <div
      class="flex-1 flex flex-col bg-bg-elevated rounded-xl border border-border-primary overflow-hidden shadow-sm"
    >
      <!-- 可滚动的工作流网格 -->
      <div class="flex-1 overflow-y-auto p-6">
        <!-- 工作流卡片网格 -->
        <div
          v-if="!loading && workflows.length > 0"
          class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5 gap-4"
        >
          <WorkflowCard
            v-for="workflow in workflows"
            :key="workflow.id"
            :workflow="workflow"
            @edit="handleEdit"
            @executions="handleViewExecutions"
            @execute="handleExecute"
            @delete="handleDelete"
            @toggle="handleToggle"
            @refresh="handleRefresh"
          />
        </div>

        <!-- 空状态 -->
        <div
          v-else-if="!loading && workflows.length === 0"
          class="flex items-center justify-center h-full"
        >
          <div class="text-text-placeholder text-center">
            <Workflow class="w-16 h-16 mx-auto mb-4 opacity-50" />
            <p class="text-lg font-medium">暂无工作流</p>
            <p class="text-sm">点击上方按钮创建第一个工作流</p>
          </div>
        </div>

        <!-- 加载中 -->
        <div v-else class="flex items-center justify-center h-full">
          <div class="text-text-tertiary">加载中...</div>
        </div>
      </div>

      <!-- 分页 - 固定在底部 -->
      <div
        v-if="!loading && totalWorkflows > pageSize"
        class="border-t border-border-primary flex-shrink-0"
      >
        <Pagination
          :current="currentPage"
          :page-size="pageSize"
          :total="totalWorkflows"
          :show-range="true"
          :show-jumper="totalWorkflows > 100"
          @change="handlePageChange"
        />
      </div>
    </div>

    <!-- 对话框 -->
    <ExecuteWithParamsDialog
      :visible="executeDialogVisible"
      :workflow="selectedWorkflow"
      @close="executeDialogVisible = false"
      @execute="handleExecuteWithParams"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Plus, Workflow, Search } from 'lucide-vue-next'
import BaseButton from '@/components/BaseButton'
import BaseInput from '@/components/BaseInput'
import RadioGroup from '@/components/RadioGroup/index.vue'
import Pagination from '@/components/Pagination'
import WorkflowCard from './components/WorkflowCard.vue'
import ExecuteWithParamsDialog from './components/ExecuteWithParamsDialog.vue'
import { workflowApi } from '@/api/workflow'
import { message } from '@/utils/message'
import type { Workflow as WorkflowType } from '@/types/workflow'

const router = useRouter()
const loading = ref(false)
const workflows = ref<WorkflowType[]>([])
const executeDialogVisible = ref(false)
const selectedWorkflow = ref<WorkflowType | null>(null)
const searchKeyword = ref('')
const filterStatus = ref<string>('all')

// 状态筛选选项
const statusOptions = [
  { label: '全部', value: 'all' },
  { label: '启用', value: 'enabled' },
  { label: '禁用', value: 'disabled' },
]

// 分页相关
const currentPage = ref(1)
const pageSize = ref(20)
const totalWorkflows = ref(0)

// 加载工作流列表
const loadWorkflows = async () => {
  loading.value = true
  try {
    const params: any = {
      page: currentPage.value,
      page_size: pageSize.value,
    }

    // 添加搜索关键词
    if (searchKeyword.value.trim()) {
      params.keyword = searchKeyword.value.trim()
    }

    // 添加状态筛选
    if (filterStatus.value === 'enabled') {
      params.enabled = true
    } else if (filterStatus.value === 'disabled') {
      params.enabled = false
    }
    // 'all' 时不传 enabled 参数

    const data = await workflowApi.list(params)
    workflows.value = data.items || []
    totalWorkflows.value = data.total || 0
  } catch (error) {
    console.error('Failed to load workflows:', error)
    message.error('加载工作流列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索处理
const handleSearch = () => {
  currentPage.value = 1 // 搜索时重置到第一页
  loadWorkflows()
}

// 筛选变化处理
const handleFilterChange = () => {
  currentPage.value = 1 // 筛选时重置到第一页
  loadWorkflows()
}

// 页码变化
const handlePageChange = (page: number) => {
  currentPage.value = page
  loadWorkflows()
}

// 编辑工作流
const handleEdit = (workflow: WorkflowType) => {
  router.push(`/workflows/${workflow.id}/edit`)
}

// 查看执行历史
const handleViewExecutions = (workflow: WorkflowType) => {
  router.push(`/workflows/${workflow.id}/executions`)
}

// 删除工作流
const handleDelete = async (workflow: WorkflowType) => {
  try {
    await workflowApi.delete(workflow.id)
    message.success('删除成功')
    await loadWorkflows()
  } catch (error) {
    console.error('Delete workflow failed:', error)
    message.error('删除失败')
  }
}

// 检查是否需要外部参数
const needsExternalParams = (workflow: WorkflowType): boolean => {
  // 找到第一个节点（没有入边的节点）
  const firstNode = workflow.nodes.find((node) => {
    const hasIncomingEdge = workflow.edges.some((edge) => edge.target === node.id)
    return !hasIncomingEdge
  })

  // 如果第一个节点是 external_trigger 且有参数配置，则需要参数
  return (
    firstNode?.type === 'external_trigger' &&
    firstNode.config?.params &&
    Array.isArray(firstNode.config.params) &&
    firstNode.config.params.length > 0
  )
}

// 执行工作流
const handleExecute = async (workflow: WorkflowType) => {
  if (!workflow.enabled) {
    message.warning('工作流未启用')
    return
  }

  // 检查是否需要参数
  if (needsExternalParams(workflow)) {
    selectedWorkflow.value = workflow
    executeDialogVisible.value = true
    return
  }

  // 无需参数，直接执行
  try {
    const data = await workflowApi.execute(workflow.id)
    message.success('工作流已开始执行')
    router.push(`/workflows/${workflow.id}/executions/${data.execution_id}`)
  } catch (error) {
    console.error('Execute workflow failed:', error)
    message.error('执行失败')
  }
}

// 带参数执行工作流
const handleExecuteWithParams = async (params: Record<string, any> | FormData) => {
  if (!selectedWorkflow.value) return

  const workflowId = selectedWorkflow.value.id

  try {
    // 如果是 FormData，直接传递；否则包装成 {params}
    const data = await workflowApi.execute(
      workflowId,
      params instanceof FormData ? params : { params }
    )
    message.success('工作流已开始执行')
    executeDialogVisible.value = false
    selectedWorkflow.value = null
    router.push(`/workflows/${workflowId}/executions/${data.execution_id}`)
  } catch (error) {
    console.error('Execute workflow failed:', error)
    message.error('执行失败')
  }
}

// 切换工作流状态
const handleToggle = async (workflow: WorkflowType) => {
  try {
    const newEnabled = !workflow.enabled
    await workflowApi.toggleEnabled(workflow.id, newEnabled)
    workflow.enabled = newEnabled
    message.success(newEnabled ? '已启用' : '已禁用')
  } catch (error) {
    console.error('Toggle workflow failed:', error)
    message.error('操作失败')
  }
}

// 工作流到达执行时间时刷新列表
const handleRefresh = async () => {
  try {
    await loadWorkflows()
  } catch (error) {
    console.error('Refresh workflows failed:', error)
  }
}

onMounted(() => {
  loadWorkflows()
})
</script>
