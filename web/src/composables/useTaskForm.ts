import { ref, computed } from 'vue'
import type { Task } from '@/api/task'
import * as toolApi from '@/api/tool'

export interface TaskFormData {
  name: string
  scheduleType: string
  scheduleValue: string
  tool_code: string
}

export function useTaskForm() {
  const taskForm = ref<TaskFormData>({
    name: '',
    scheduleType: 'daily',
    scheduleValue: '09:00:00',
    tool_code: '',
  })

  const toolConfig = ref<Record<string, any>>({
    url: '',
    method: 'GET',
    headers: '{}',
    body: '{}',
  })

  const isConfigured = ref(false)
  const editingTask = ref<Task | null>(null)
  const tools = ref<toolApi.Tool[]>([])

  const toolOptions = computed(() =>
    tools.value.map((tool) => ({ label: tool.name, value: tool.code }))
  )

  const resetForm = () => {
    editingTask.value = null
    taskForm.value = {
      name: '',
      scheduleType: 'daily',
      scheduleValue: '09:00:00',
      tool_code: '',
    }
    toolConfig.value = {
      url: '',
      method: 'GET',
      headers: '{}',
      body: '{}',
    }
    isConfigured.value = false
  }

  const handleToolChange = () => {
    toolConfig.value = {
      url: '',
      method: 'GET',
      headers: '{}',
      body: '{}',
    }
    isConfigured.value = false
  }

  const loadTaskForEdit = (task: Task) => {
    editingTask.value = task
    taskForm.value = {
      name: task.name,
      scheduleType: task.schedule_type,
      scheduleValue: task.schedule_value,
      tool_code: task.tool_code,
    }

    // Parse existing config
    if (task.config) {
      try {
        const config = JSON.parse(task.config)
        toolConfig.value = {
          url: config.url || '',
          method: config.method || 'GET',
          headers: JSON.stringify(config.headers || {}, null, 2),
          body: JSON.stringify(config.body || {}, null, 2),
        }
        isConfigured.value = true
      } catch (error) {
        console.error('Failed to parse task config:', error)
      }
    }
  }

  const loadTaskForClone = (task: Task) => {
    editingTask.value = null
    taskForm.value = {
      name: task.name + ' (副本)',
      scheduleType: task.schedule_type,
      scheduleValue: task.schedule_value,
      tool_code: task.tool_code,
    }

    // Parse existing config
    if (task.config) {
      try {
        const config = JSON.parse(task.config)
        toolConfig.value = {
          url: config.url || '',
          method: config.method || 'GET',
          headers: JSON.stringify(config.headers || {}, null, 2),
          body: JSON.stringify(config.body || {}, null, 2),
        }
        isConfigured.value = true
      } catch (error) {
        console.error('Failed to parse task config:', error)
      }
    }
  }

  return {
    taskForm,
    toolConfig,
    isConfigured,
    editingTask,
    tools,
    toolOptions,
    resetForm,
    handleToolChange,
    loadTaskForEdit,
    loadTaskForClone,
  }
}
