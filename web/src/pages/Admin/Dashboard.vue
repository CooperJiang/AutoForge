<template>
  <div class="px-6 py-6">
    <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-6 gap-4 mb-6">
      <div
        class="bg-bg-elevated border-2 border-border-primary rounded-xl p-5 shadow-sm hover:shadow-lg transition-all group"
      >
        <div class="flex items-start justify-between mb-3">
          <div class="p-3 bg-blue-500/10 rounded-lg">
            <Users class="w-6 h-6 text-blue-600 dark:text-blue-400" />
          </div>
          <div
            class="text-xs font-semibold px-2 py-1 bg-blue-500/10 text-blue-600 dark:text-blue-400 rounded-full"
          >
            用户
          </div>
        </div>
        <div>
          <p class="text-text-tertiary text-xs font-medium mb-1">总用户数</p>
          <p class="text-3xl font-bold text-text-primary">{{ stats.total_users }}</p>
        </div>
      </div>

      <div
        class="bg-bg-elevated border-2 border-border-primary rounded-xl p-5 shadow-sm hover:shadow-lg transition-all group"
      >
        <div class="flex items-start justify-between mb-3">
          <div class="p-3 bg-purple-500/10 rounded-lg">
            <ListTodo class="w-6 h-6 text-purple-600 dark:text-purple-400" />
          </div>
          <div
            class="text-xs font-semibold px-2 py-1 bg-purple-500/10 text-purple-600 dark:text-purple-400 rounded-full"
          >
            任务
          </div>
        </div>
        <div>
          <p class="text-text-tertiary text-xs font-medium mb-1">总任务数</p>
          <p class="text-3xl font-bold text-text-primary">{{ stats.total_tasks }}</p>
        </div>
      </div>

      <div
        class="bg-bg-elevated border-2 border-border-primary rounded-xl p-5 shadow-sm hover:shadow-lg transition-all group"
      >
        <div class="flex items-start justify-between mb-3">
          <div class="p-3 bg-green-500/10 rounded-lg">
            <Activity class="w-6 h-6 text-green-600 dark:text-green-400" />
          </div>
          <div
            class="text-xs font-semibold px-2 py-1 bg-green-500/10 text-green-600 dark:text-green-400 rounded-full"
          >
            执行
          </div>
        </div>
        <div>
          <p class="text-text-tertiary text-xs font-medium mb-1">今日执行</p>
          <p class="text-3xl font-bold text-text-primary">{{ stats.today_executions }}</p>
        </div>
      </div>

      <div
        class="bg-bg-elevated border-2 border-border-primary rounded-xl p-5 shadow-sm hover:shadow-lg transition-all group"
      >
        <div class="flex items-start justify-between mb-3">
          <div class="p-3 bg-orange-500/10 rounded-lg">
            <TrendingUp class="w-6 h-6 text-orange-600 dark:text-orange-400" />
          </div>
          <div
            class="text-xs font-semibold px-2 py-1 bg-orange-500/10 text-orange-600 dark:text-orange-400 rounded-full"
          >
            成功率
          </div>
        </div>
        <div>
          <p class="text-text-tertiary text-xs font-medium mb-1">执行成功率</p>
          <p class="text-3xl font-bold text-text-primary">
            {{ stats.success_rate.toFixed(1) }}<span class="text-xl">%</span>
          </p>
        </div>
      </div>

      <div
        class="bg-bg-elevated border-2 border-border-primary rounded-xl p-5 shadow-sm hover:shadow-lg transition-all group"
      >
        <div class="flex items-start justify-between mb-3">
          <div class="p-3 bg-indigo-500/10 rounded-lg">
            <Workflow class="w-6 h-6 text-indigo-600 dark:text-indigo-400" />
          </div>
          <div
            class="text-xs font-semibold px-2 py-1 bg-indigo-500/10 text-indigo-600 dark:text-indigo-400 rounded-full"
          >
            工作流
          </div>
        </div>
        <div>
          <p class="text-text-tertiary text-xs font-medium mb-1">总工作流数</p>
          <p class="text-3xl font-bold text-text-primary">{{ stats.total_workflows }}</p>
        </div>
      </div>

      <div
        class="bg-bg-elevated border-2 border-border-primary rounded-xl p-5 shadow-sm hover:shadow-lg transition-all group"
      >
        <div class="flex items-start justify-between mb-3">
          <div class="p-3 bg-cyan-500/10 rounded-lg">
            <Package class="w-6 h-6 text-cyan-600 dark:text-cyan-400" />
          </div>
          <div
            class="text-xs font-semibold px-2 py-1 bg-cyan-500/10 text-cyan-600 dark:text-cyan-400 rounded-full"
          >
            工作流
          </div>
        </div>
        <div>
          <p class="text-text-tertiary text-xs font-medium mb-1">工作流市场</p>
          <p class="text-3xl font-bold text-text-primary">{{ stats.total_templates }}</p>
        </div>
      </div>
    </div>

    <!-- Tab Content -->
    <div class="bg-bg-elevated rounded-xl shadow-lg border-2 border-border-primary p-6">
      <div v-show="activeTab === 'executions'">
        <ExecutionsTab />
      </div>

      <div v-show="activeTab === 'tasks'">
        <TasksTab
          :tasks="tasks"
          :total="total"
          :current-page="currentPage"
          :total-pages="totalPages"
          @search="handleTaskSearch"
          @copy-id="handleCopy"
          @view="viewTask"
          @toggle-status="toggleTaskStatus"
          @execute="executeTaskNow"
          @delete="deleteTaskConfirm"
          @page-change="handleTaskPageChange"
        />
      </div>

      <div v-show="activeTab === 'users'">
        <UsersTab />
      </div>

      <div v-show="activeTab === 'workflows'">
        <TemplateManagement />
      </div>

      <div v-show="activeTab === 'categories'">
        <CategoryManagement />
      </div>

      <div v-show="activeTab === 'tools'">
        <ToolManagement />
      </div>
    </div>

    <Dialog
      v-model="showDeleteDialog"
      title="确认删除"
      :message="`确定要删除任务 &quot;${taskToDelete?.name}&quot; 吗？此操作不可恢复！`"
      confirm-text="删除"
      cancel-text="取消"
      confirm-variant="danger"
      @confirm="deleteTask"
    />

    <TestResultDialog v-model="showTestResult" :result="testResult" />

    <TaskDetailDialog v-model="showTaskDetail" :task="selectedTask" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { Users, ListTodo, Activity, TrendingUp, Workflow, Package } from 'lucide-vue-next'

const props = defineProps<{
  activeTab?: string
}>()

const emit = defineEmits<{
  'update:active-tab': [value: string]
}>()

// Use props or default to 'executions'
const internalActiveTab = ref(props.activeTab || 'executions')

// Computed to sync with props
const activeTab = computed({
  get: () => props.activeTab || internalActiveTab.value,
  set: (value) => {
    internalActiveTab.value = value
    emit('update:active-tab', value)
  },
})

import * as adminApi from '@/api/admin'
import type { Task } from '@/api/task'
import type { StatsResponse } from '@/api/admin'
import { message } from '@/utils/message'
import Dialog from '@/components/Dialog'
import TestResultDialog from '@/components/TestResultDialog'
import TaskDetailDialog from '@/components/TaskDetailDialog'
import { copyToClipboard } from '@/utils/clipboard'
import ExecutionsTab from './ExecutionsTab.vue'
import TasksTab from './components/TasksTab.vue'
import UsersTab from './UsersTab.vue'
import TemplateManagement from './TemplateManagement.vue'
import CategoryManagement from './CategoryManagement.vue'
import ToolManagement from './ToolManagement.vue'

// 统计数据
const stats = ref<StatsResponse>({
  total_users: 0,
  total_tasks: 0,
  active_tasks: 0,
  today_executions: 0,
  success_rate: 0,
  total_workflows: 0,
  total_templates: 0,
  recent_users: [],
})

// 任务列表
const tasks = ref<Task[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

// 筛选条件
const filters = ref({
  user_id: '',
  status: '',
  keyword: '',
})

// Dialog控制
const showDeleteDialog = ref(false)
const taskToDelete = ref<Task | null>(null)

// 测试结果Dialog控制
const showTestResult = ref(false)
const testResult = ref<any>(null)

// 任务详情Dialog控制
const showTaskDetail = ref(false)
const selectedTask = ref<Task | null>(null)

// 计算总页数
const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

// 加载统计数据
const loadStats = async () => {
  try {
    const res = await adminApi.getStats()
    stats.value = res.data
  } catch (error: any) {
    message.error(error.response?.data?.error || '加载统计数据失败')
  }
}

// 加载任务列表
const loadTasks = async () => {
  try {
    const res = await adminApi.getTasks({
      page: currentPage.value,
      page_size: pageSize.value,
      ...filters.value,
    })
    tasks.value = res.data.tasks || []
    total.value = res.data.total
  } catch (error: any) {
    message.error(error.response?.data?.error || '加载任务列表失败')
  }
}

// 处理任务搜索
const handleTaskSearch = (searchFilters: { user_id: string; status: string; keyword: string }) => {
  filters.value = searchFilters
  currentPage.value = 1
  loadTasks()
}

// 切换任务状态
const toggleTaskStatus = async (task: Task) => {
  try {
    await adminApi.updateTaskStatus(task.id, !task.enabled)
    message.success(task.enabled ? '已禁用任务' : '已启用任务')
    loadTasks()
    loadStats()
  } catch (error: any) {
    message.error(error.response?.data?.error || '操作失败')
  }
}

// 立即执行任务
const executeTaskNow = async (task: Task) => {
  try {
    await adminApi.executeTask(task.id)
    message.success('任务已提交执行，请稍后查看执行记录')
    // 延迟刷新任务列表和统计信息
    setTimeout(() => {
      loadTasks()
      loadStats()
    }, 1000)
  } catch (error: any) {
    message.error(error.response?.data?.message || error.message || '执行失败')
  }
}

// 查看任务详情
const viewTask = (task: Task) => {
  selectedTask.value = task
  showTaskDetail.value = true
}

// 删除任务确认
const deleteTaskConfirm = (task: Task) => {
  taskToDelete.value = task
  showDeleteDialog.value = true
}

// 删除任务
const deleteTask = async () => {
  if (!taskToDelete.value) return

  try {
    await adminApi.deleteTask(taskToDelete.value.id)
    message.success('任务已删除')
    loadTasks()
    loadStats()
  } catch (error: any) {
    message.error(error.response?.data?.message || '删除失败')
  }
}

// 复制ID到剪贴板
const handleCopy = (id: string) => {
  copyToClipboard(id)
}

// 任务页码变化
const handleTaskPageChange = (page: number) => {
  currentPage.value = page
  loadTasks()
}

onMounted(() => {
  loadStats()
  loadTasks()
})
</script>
