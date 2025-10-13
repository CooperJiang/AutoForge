<template>
  <!-- Loading -->
  <div v-if="loading" class="flex justify-center items-center py-20">
    <div class="text-text-tertiary">加载中...</div>
  </div>

  <!-- Main Content -->
  <main v-else>
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Task Form Panel -->
      <div class="lg:col-span-1">
        <TaskFormPanel
          :task-form="taskForm"
          :tool-options="toolOptions"
          :is-configured="isConfigured"
          :testing="testing"
          :editing-task="editingTask"
          @submit="handleSubmitTask"
          @cancel="cancelEdit"
          @tool-change="handleToolChange"
          @config-click="showConfigDialog = true"
          @test-click="handleTestConfig"
        />
      </div>

      <!-- Task List Panel -->
      <div class="lg:col-span-2">
        <TaskListPanel
          :tasks="tasks"
          :filtered-tasks="filteredTasks"
          :executions="executions"
          :selected-task-id="selectedTaskId"
          v-model:search-keyword="searchKeyword"
          v-model:status-filter="statusFilter"
          @test-task="handleTestTask"
          @edit-task="editTask"
          @clone-task="cloneTask"
          @toggle-task="toggleTask"
          @delete-task="confirmDelete"
          @refresh-executions="loadExecutions"
          @delete-all-executions="deleteAllTaskExecutions"
          @delete-execution="deleteExecution"
          @toggle-expanded="toggleExpanded"
        />
      </div>
    </div>
  </main>

  <!-- Dialogs -->
  <Dialog
    v-model="dialogVisible"
    :title="dialogConfig.title"
    :message="dialogConfig.message"
    :confirm-text="dialogConfig.confirmText"
    :cancel-text="dialogConfig.cancelText"
    :confirm-variant="dialogConfig.confirmVariant"
    @confirm="dialogConfig.onConfirm"
  />

  <ToolConfigDrawer
    v-model="showConfigDialog"
    :tool-code="taskForm.tool_code"
    :config="toolConfig"
    @update:config="toolConfig = $event"
    @save="saveConfig"
  />

  <TestResultDialog
    v-model="showTestDialog"
    :result="testResult"
  />
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import Dialog from '@/components/Dialog'
import ToolConfigDrawer from './components/ToolConfigDrawer.vue'
import TestResultDialog from './components/TestResultDialog.vue'
import TaskFormPanel from './components/TaskFormPanel.vue'
import TaskListPanel from './components/TaskListPanel.vue'
import { message } from '@/utils/message'
import * as taskApi from '@/api/task'
import * as toolApi from '@/api/tool'
import { useTaskForm } from '@/composables/useTaskForm'
import { useTaskList } from '@/composables/useTaskList'
import { getDefaultScheduleValue } from '@/utils/taskHelpers'
import type { Task } from '@/api/task'

const route = useRoute()

// Use composables
const {
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
  selectTool,
} = useTaskForm()

const {
  tasks,
  executions,
  selectedTaskId,
  loading,
  searchKeyword,
  statusFilter,
  filteredTasks,
  loadTasks,
  loadExecutions,
  toggleTask,
  deleteExecution,
  deleteAllTaskExecutions,
} = useTaskList()

// Dialog state
const dialogVisible = ref(false)
const dialogConfig = ref({
  title: '确认',
  message: '确定要执行此操作吗？',
  confirmText: '确定',
  cancelText: '取消',
  confirmVariant: 'primary' as 'primary' | 'danger',
  onConfirm: () => {},
})

const showConfigDialog = ref(false)

// Test state
const testing = ref(false)
const testResult = ref<any>(null)
const showTestDialog = ref(false)

// Load tools
const loadTools = async () => {
  try {
    tools.value = await toolApi.getToolList()
  } catch (error: any) {
    message.error('加载工具列表失败: ' + (error.message || '未知错误'))
  }
}

// Handle task submission
const handleSubmitTask = async () => {
  try {
    if (!isConfigured.value) {
      message.error('请先配置工具参数')
      return
    }

    let config: Record<string, any> = {}
    if (taskForm.value.tool_code === 'http_request') {
      try {
        config = {
          url: toolConfig.value.url,
          method: toolConfig.value.method,
          headers: JSON.parse(toolConfig.value.headers || '{}'),
          body: JSON.parse(toolConfig.value.body || '{}'),
        }
      } catch {
        message.error('Headers 或 Body 不是有效的 JSON 格式')
        return
      }
    }

    if (editingTask.value) {
      await taskApi.updateTask(editingTask.value.id, {
        name: taskForm.value.name,
        tool_code: taskForm.value.tool_code,
        config,
        schedule_type: taskForm.value.scheduleType,
        schedule_value: taskForm.value.scheduleValue,
      })
      message.success('任务更新成功')
    } else {
      await taskApi.createTask({
        name: taskForm.value.name,
        tool_code: taskForm.value.tool_code,
        config,
        schedule_type: taskForm.value.scheduleType,
        schedule_value: taskForm.value.scheduleValue,
      })
      message.success('任务创建成功')
    }
    resetForm()
    await loadTasks()
  } catch (error: any) {
    message.error('操作失败: ' + (error.response?.data?.message || error.message))
  }
}

// Edit task
const editTask = (task: Task) => {
  loadTaskForEdit(task)
}

// Cancel edit
const cancelEdit = () => {
  resetForm()
}

// Clone task
const cloneTask = (task: Task) => {
  loadTaskForClone(task)
  message.success('已克隆任务配置，请修改后保存')
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

// Confirm delete
const confirmDelete = (task: Task) => {
  dialogConfig.value = {
    title: '确认删除',
    message: `确定要删除任务 "${task.name}" 吗？`,
    confirmText: '删除',
    cancelText: '取消',
    confirmVariant: 'danger',
    onConfirm: async () => {
      try {
        await taskApi.deleteTask(task.id)
        message.success('删除成功')
        await loadTasks()
      } catch (error: any) {
        message.error('删除失败: ' + (error.response?.data?.message || error.message))
      }
    },
  }
  dialogVisible.value = true
}

// Save config
const saveConfig = () => {
  if (taskForm.value.tool_code === 'http_request') {
    if (!toolConfig.value.url) {
      message.error('请输入请求URL')
      return
    }
  }
  isConfigured.value = true
  showConfigDialog.value = false
}

// Test config
const handleTestConfig = async () => {
  if (!isConfigured.value) {
    message.error('请先配置工具参数')
    return
  }

  testing.value = true
  try {
    let testConfig: any = {
      tool_code: taskForm.value.tool_code
    }

    // 根据不同工具类型构建测试配置
    if (taskForm.value.tool_code === 'http_request') {
      let headersObj = {}
      let bodyObj = {}

      try {
        headersObj = JSON.parse(toolConfig.value.headers || '{}')
      } catch {
        message.error('Headers 不是有效的 JSON 格式')
        testing.value = false
        return
      }

      try {
        bodyObj = JSON.parse(toolConfig.value.body || '{}')
      } catch {
        message.error('Body 不是有效的 JSON 格式')
        testing.value = false
        return
      }

      const headers = Object.entries(headersObj).map(([key, value]) => ({
        key,
        value: String(value),
      }))

      testConfig = {
        tool_code: taskForm.value.tool_code,
        url: toolConfig.value.url,
        method: toolConfig.value.method,
        headers,
        params: [],
        body: JSON.stringify(bodyObj),
      }
    } else {
      // 其他工具（邮件、健康检查等）使用 config 字段
      testConfig = {
        tool_code: taskForm.value.tool_code,
        config: toolConfig.value
      }
    }

    const response = await taskApi.testTask(testConfig)

    testResult.value = response
    showTestDialog.value = true

    if (response.success) {
      message.success('测试成功')
    } else {
      // 检查是否是反垃圾邮件拦截
      const errorMsg = response.error_message || ''
      const isSpamBlocked = errorMsg.toLowerCase().includes('spam') ||
                           errorMsg.includes('垃圾邮件') ||
                           errorMsg.includes('ANTISPAM') ||
                           errorMsg.includes('Reject by content')

      if (isSpamBlocked) {
        message.warning('邮件发送成功，但被反垃圾邮件系统拦截。建议使用更完整的邮件内容')
      } else if (errorMsg) {
        message.warning('测试失败: ' + errorMsg)
      } else {
        message.warning('测试完成，但请求未成功')
      }
    }
  } catch (error: any) {
    message.error('测试失败: ' + (error.response?.data?.message || error.message))
  } finally {
    testing.value = false
  }
}

// Test task from card
const handleTestTask = async (task: Task) => {
  try {
    const config = JSON.parse(task.config)

    const headersObj = config.headers || {}
    const bodyObj = config.body || {}

    const headers = Object.entries(headersObj).map(([key, value]) => ({
      key,
      value: String(value),
    }))

    const response = await taskApi.testTask({
      url: config.url,
      method: config.method,
      headers,
      params: [],
      body: JSON.stringify(bodyObj),
    })

    testResult.value = response
    showTestDialog.value = true

    if (response.success) {
      message.success('测试成功')
    } else {
      message.warning('测试完成，但请求未成功')
    }
  } catch (error: any) {
    message.error('测试失败: ' + (error.response?.data?.message || error.message))
  }
}

// Toggle expanded
const toggleExpanded = (taskId: string) => {
  if (selectedTaskId.value === taskId) {
    selectedTaskId.value = null
  } else {
    loadExecutions(taskId, true)
  }
}

// Watch schedule type changes
watch(
  () => taskForm.value.scheduleType,
  (newType, oldType) => {
    if (editingTask.value) return
    if (newType && newType !== oldType) {
      taskForm.value.scheduleValue = getDefaultScheduleValue(newType)
    }
  }
)

// Initialize
onMounted(async () => {
  await loadTools()
  await loadTasks()

  // Handle tool selection from query parameter
  if (route.query.tool) {
    selectTool(route.query.tool as string)
  }
})
</script>
