import { onMounted, type Ref } from 'vue'
import { useRouter } from 'vue-router'
import { workflowApi } from '@/api/workflow'
import { message } from '@/utils/message'
import type { WorkflowExecution } from '../types'

/**
 * 执行历史操作逻辑
 */
export function useExecutionsActions(options: {
  workflowId: string
  workflowName: Ref<string>
  workflowInfo: Ref<any>
  statusFilter: Ref<string>
  loading: Ref<boolean>
  refreshing: Ref<boolean>
  executions: Ref<WorkflowExecution[]>
  executionToDelete: Ref<WorkflowExecution | null>
  showDeleteDialog: Ref<boolean>
}) {
  const {
    workflowId,
    workflowName,
    workflowInfo,
    statusFilter,
    loading,
    refreshing,
    executions,
    executionToDelete,
    showDeleteDialog,
  } = options

  const router = useRouter()

  // 加载工作流详情
  const loadWorkflow = async () => {
    try {
      const data = await workflowApi.getById(workflowId)
      workflowName.value = data.name
      workflowInfo.value = data
    } catch (error) {
      console.error('Load workflow failed:', error)
    }
  }

  // 加载执行历史
  const loadExecutions = async () => {
    loading.value = true
    try {
      const data = await workflowApi.getExecutions(workflowId, {
        page: 1,
        page_size: 100,
        status: statusFilter.value === 'all' ? undefined : statusFilter.value,
      })
      executions.value = data.items || []
    } catch (error: any) {
      console.error('Load executions failed:', error)
      message.error(error.response?.data?.message || '加载执行历史失败')
    } finally {
      loading.value = false
    }
  }

  // 刷新
  const handleRefresh = async () => {
    refreshing.value = true
    try {
      await loadExecutions()
      message.success('刷新成功')
    } catch {
      message.error('刷新失败')
    } finally {
      refreshing.value = false
    }
  }

  // 查看执行详情
  const handleViewExecution = (execution: WorkflowExecution) => {
    router.push(`/workflows/${workflowId}/executions/${execution.id}`)
  }

  // 删除执行记录
  const handleDeleteExecution = (execution: WorkflowExecution) => {
    executionToDelete.value = execution
    showDeleteDialog.value = true
  }

  // 确认删除
  const confirmDelete = async () => {
    if (!executionToDelete.value) return

    try {
      await workflowApi.deleteExecution(workflowId, executionToDelete.value.id)
      message.success('删除成功')
      executions.value = executions.value.filter((e) => e.id !== executionToDelete.value!.id)
      executionToDelete.value = null
    } catch (error: any) {
      message.error(error.message || '删除失败')
    }
  }

  // 倒计时结束回调 - 自动刷新执行历史
  const handleCountdownFinish = async () => {
    setTimeout(async () => {
      try {
        await Promise.all([loadWorkflow(), loadExecutions()])
        message.success('工作流已执行，列表已自动刷新')
      } catch (error) {
        console.error('Auto refresh failed:', error)
        message.error('自动刷新失败，请手动刷新')
      }
    }, 1000)
  }

  // 执行后回调
  const onExecuted = async () => {
    await loadExecutions()
  }

  // 初始化
  onMounted(() => {
    loadWorkflow()
    loadExecutions()
  })

  return {
    loadWorkflow,
    loadExecutions,
    handleRefresh,
    handleViewExecution,
    handleDeleteExecution,
    confirmDelete,
    handleCountdownFinish,
    onExecuted,
  }
}

