import { ref, computed } from 'vue'
import type { Task, TaskExecution } from '@/api/task'
import * as taskApi from '@/api/task'
import { message } from '@/utils/message'

export function useTaskList() {
  const tasks = ref<Task[]>([])
  const executions = ref<TaskExecution[]>([])
  const selectedTaskId = ref<string | null>(null)
  const loading = ref(false)
  const refreshing = ref(false)


  const searchKeyword = ref('')
  const statusFilter = ref<'all' | 'enabled' | 'disabled'>('all')


  const filteredTasks = computed(() => {
    if (!tasks.value || !Array.isArray(tasks.value)) {
      return []
    }

    return tasks.value.filter((task) => {
      const keyword = searchKeyword.value.toLowerCase()
      const matchSearch =
        keyword === '' ||
        task.name.toLowerCase().includes(keyword) ||
        task.tool_code.toLowerCase().includes(keyword)

      const matchStatus =
        statusFilter.value === 'all' ||
        (statusFilter.value === 'enabled' && task.enabled) ||
        (statusFilter.value === 'disabled' && !task.enabled)

      return matchSearch && matchStatus
    })
  })

  const loadTasks = async () => {
    loading.value = true
    try {
      const response = await taskApi.getTaskList()
      if (response && response.items) {
        tasks.value = response.items
      } else {
        tasks.value = []
      }
    } catch (error: any) {
      tasks.value = []
      message.error('加载任务失败: ' + (error.message || '未知错误'))
    } finally {
      loading.value = false
    }
  }

  const loadExecutions = async (taskId: string, showLoading = false) => {
    selectedTaskId.value = taskId
    if (showLoading) refreshing.value = true
    try {
      const response = await taskApi.getTaskExecutions(taskId)
      executions.value = response.items
    } catch (error: any) {
      message.error('加载执行记录失败: ' + (error.message || '未知错误'))
    } finally {
      if (showLoading) refreshing.value = false
    }
  }

  const deleteTask = async (taskId: string) => {
    try {
      await taskApi.deleteTask(taskId)
      message.success('删除成功')
      await loadTasks()
    } catch (error: any) {
      message.error('删除失败: ' + (error.response?.data?.message || error.message))
    }
  }

  const toggleTask = async (task: Task) => {
    try {
      if (task.enabled) {
        await taskApi.disableTask(task.id)
      } else {
        await taskApi.enableTask(task.id)
      }
      await loadTasks()
    } catch (error: any) {
      message.error('操作失败: ' + (error.response?.data?.message || error.message))
    }
  }

  const deleteExecution = async (executionId: string) => {
    try {
      await taskApi.deleteExecution(executionId)
      message.success('已删除执行记录')
      if (selectedTaskId.value) {
        await loadExecutions(selectedTaskId.value, true)
      }
    } catch (error: any) {
      message.error('删除失败: ' + (error.response?.data?.message || error.message))
    }
  }

  const deleteAllTaskExecutions = async (taskId: string) => {
    try {
      await taskApi.deleteAllExecutions(taskId)
      message.success('已删除所有执行记录')
      await loadExecutions(taskId, true)
    } catch (error: any) {
      message.error('删除失败: ' + (error.response?.data?.message || error.message))
    }
  }

  return {
    tasks,
    executions,
    selectedTaskId,
    loading,
    refreshing,
    searchKeyword,
    statusFilter,
    filteredTasks,
    loadTasks,
    loadExecutions,
    deleteTask,
    toggleTask,
    deleteExecution,
    deleteAllTaskExecutions,
  }
}
