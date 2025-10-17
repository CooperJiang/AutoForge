import { onMounted, onUnmounted, type Ref } from 'vue'
import { useRouter } from 'vue-router'
import { workflowApi } from '@/api/workflow'
import { message } from '@/utils/message'
import type { WorkflowExecution } from '../types'

/**
 * 执行详情操作逻辑
 */
export function useExecutionActions(options: {
  workflowId: string
  executionId: string
  execution: Ref<WorkflowExecution | null>
  loading: Ref<boolean>
  polling: Ref<boolean>
}) {
  const { workflowId, executionId, execution, loading, polling } = options
  const router = useRouter()

  let pollingTimer: NodeJS.Timeout | null = null

  // ============ 加载执行详情 ============

  const loadExecutionDetail = async () => {
    try {
      loading.value = true
      const data = await workflowApi.getExecutionDetail(workflowId, executionId)
      execution.value = data as WorkflowExecution
    } catch (error: any) {
      message.error(error.response?.data?.message || '加载执行详情失败')
    } finally {
      loading.value = false
    }
  }

  // ============ 轮询更新 ============

  const startPolling = () => {
    if (polling.value) return

    polling.value = true
    pollingTimer = setInterval(async () => {
      if (!execution.value) return

      // 如果已经完成，停止轮询
      if (execution.value.status !== 'running') {
        stopPolling()
        return
      }

      try {
        const data = await workflowApi.getExecutionDetail(workflowId, executionId)
        execution.value = data as WorkflowExecution
      } catch (error) {
        stopPolling()
      }
    }, 2000) // 每2秒轮询一次
  }

  const stopPolling = () => {
    if (pollingTimer) {
      clearInterval(pollingTimer)
      pollingTimer = null
    }
    polling.value = false
  }

  // ============ 导航操作 ============

  const handleBack = () => {
    router.push(`/workflows/${workflowId}/executions`)
  }

  const handleExecuted = (newExecutionId: string) => {
    // 直接跳转到新的执行详情页
    router.push(`/workflows/${workflowId}/executions/${newExecutionId}`)
  }

  // ============ 格式化函数 ============

  const formatTime = (time?: string) => {
    if (!time) return '-'
    return new Date(time).toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
    })
  }

  const formatDurationMs = (ms: number) => {
    if (ms < 1000) return `${ms}ms`
    if (ms < 60000) return `${(ms / 1000).toFixed(2)}s`
    const minutes = Math.floor(ms / 60000)
    const seconds = ((ms % 60000) / 1000).toFixed(0)
    return `${minutes}m ${seconds}s`
  }

  const formatInputData = (input: any): string => {
    if (!input) return ''

    try {
      const formatted: any = {}

      if (input.external_params && Object.keys(input.external_params).length > 0) {
        formatted['外部参数 (External Params)'] = input.external_params
      }

      if (input.config) {
        formatted['节点配置 (Config)'] = input.config
      }

      if (input.resolved_config) {
        formatted['解析后配置 (Resolved Config)'] = input.resolved_config
      }

      return JSON.stringify(formatted, null, 2)
    } catch {
      return typeof input === 'string' ? input : JSON.stringify(input, null, 2)
    }
  }

  // ============ 生命周期 ============

  onMounted(async () => {
    await loadExecutionDetail()

    // 如果状态是运行中，开始轮询
    if (execution.value?.status === 'running') {
      startPolling()
    }
  })

  onUnmounted(() => {
    stopPolling()
  })

  return {
    loadExecutionDetail,
    startPolling,
    stopPolling,
    handleBack,
    handleExecuted,
    formatTime,
    formatDurationMs,
    formatInputData,
  }
}

