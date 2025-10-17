import { ref, computed } from 'vue'
import type { WorkflowExecution, SelectOption } from '../types'

/**
 * 执行历史状态管理
 */
export function useExecutionsState() {
  const workflowName = ref('')
  const workflowInfo = ref<any>(null)
  const statusFilter = ref('all')
  const searchQuery = ref('')
  const refreshing = ref(false)
  const loading = ref(false)
  const executions = ref<WorkflowExecution[]>([])
  const showDeleteDialog = ref(false)
  const executionToDelete = ref<WorkflowExecution | null>(null)

  // 状态选项
  const statusOptions: SelectOption[] = [
    { label: '全部状态', value: 'all' },
    { label: '运行中', value: 'running' },
    { label: '成功', value: 'success' },
    { label: '失败', value: 'failed' },
    { label: '已取消', value: 'cancelled' },
  ]

  // 过滤执行记录
  const filteredExecutions = computed(() => {
    let result = executions.value

    // 状态筛选
    if (statusFilter.value !== 'all') {
      result = result.filter((e) => e.status === statusFilter.value)
    }

    // 搜索筛选
    if (searchQuery.value) {
      result = result.filter((e) => e.id.toLowerCase().includes(searchQuery.value.toLowerCase()))
    }

    return result
  })

  return {
    workflowName,
    workflowInfo,
    statusFilter,
    searchQuery,
    refreshing,
    loading,
    executions,
    showDeleteDialog,
    executionToDelete,
    statusOptions,
    filteredExecutions,
  }
}

