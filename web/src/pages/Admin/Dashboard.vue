<template>
  <div class="max-w-7xl mx-auto px-4 py-6">
    <!-- 统计卡片 -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
      <div class="bg-gradient-to-br from-blue-500 to-blue-600 rounded-xl p-6 text-white shadow-lg">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-blue-100 text-sm font-medium">总用户数</p>
            <p class="text-3xl font-bold mt-2">{{ stats.total_users }}</p>
          </div>
          <div class="text-4xl opacity-50">👥</div>
        </div>
      </div>

      <div class="bg-gradient-to-br from-green-500 to-green-600 rounded-xl p-6 text-white shadow-lg">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-green-100 text-sm font-medium">总任务数</p>
            <p class="text-3xl font-bold mt-2">{{ stats.total_tasks }}</p>
          </div>
          <div class="text-4xl opacity-50">📋</div>
        </div>
      </div>

      <div class="bg-gradient-to-br from-purple-500 to-purple-600 rounded-xl p-6 text-white shadow-lg">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-purple-100 text-sm font-medium">今日执行</p>
            <p class="text-3xl font-bold mt-2">{{ stats.today_executions }}</p>
          </div>
          <div class="text-4xl opacity-50">✅</div>
        </div>
      </div>

      <div class="bg-gradient-to-br from-orange-500 to-orange-600 rounded-xl p-6 text-white shadow-lg">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-orange-100 text-sm font-medium">成功率</p>
            <p class="text-3xl font-bold mt-2">{{ stats.success_rate.toFixed(1) }}%</p>
          </div>
          <div class="text-4xl opacity-50">📊</div>
        </div>
      </div>
    </div>

    <!-- Tab 导航 -->
    <div class="bg-white rounded-t-xl shadow-lg border-2 border-b-0 border-slate-200">
      <div class="flex border-b border-slate-200">
        <button
          @click="activeTab = 'executions'"
          :class="[
            'px-6 py-3 font-semibold transition-colors relative',
            activeTab === 'executions'
              ? 'text-green-600 bg-slate-50'
              : 'text-slate-600 hover:text-slate-900 hover:bg-slate-50'
          ]"
        >
          执行记录
          <span v-if="activeTab === 'executions'" class="absolute bottom-0 left-0 right-0 h-0.5 bg-green-600"></span>
        </button>
        <button
          @click="activeTab = 'tasks'"
          :class="[
            'px-6 py-3 font-semibold transition-colors relative',
            activeTab === 'tasks'
              ? 'text-green-600 bg-slate-50'
              : 'text-slate-600 hover:text-slate-900 hover:bg-slate-50'
          ]"
        >
          任务管理
          <span v-if="activeTab === 'tasks'" class="absolute bottom-0 left-0 right-0 h-0.5 bg-green-600"></span>
        </button>
        <button
          @click="activeTab = 'users'"
          :class="[
            'px-6 py-3 font-semibold transition-colors relative',
            activeTab === 'users'
              ? 'text-green-600 bg-slate-50'
              : 'text-slate-600 hover:text-slate-900 hover:bg-slate-50'
          ]"
        >
          用户管理
          <span v-if="activeTab === 'users'" class="absolute bottom-0 left-0 right-0 h-0.5 bg-green-600"></span>
        </button>
      </div>
    </div>

    <!-- Tab 内容 -->
    <!-- 执行记录 Tab -->
    <div v-show="activeTab === 'executions'" class="bg-white rounded-b-xl shadow-lg border-2 border-t-0 border-slate-200 p-6">
      <ExecutionsTab />
    </div>

    <!-- 任务列表 Tab -->
    <div v-show="activeTab === 'tasks'" class="bg-white rounded-b-xl shadow-lg border-2 border-t-0 border-slate-200 p-6">
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
        @prev-page="prevPage"
        @next-page="nextPage"
      />
    </div>

    <!-- 用户管理 Tab -->
    <div v-show="activeTab === 'users'" class="bg-white rounded-b-xl shadow-lg border-2 border-t-0 border-slate-200 p-6">
      <UsersTab />
    </div>

    <!-- 删除确认对话框 -->
    <Dialog
      v-model="showDeleteDialog"
      title="确认删除"
      :message="`确定要删除任务 &quot;${taskToDelete?.name}&quot; 吗？此操作不可恢复！`"
      confirm-text="删除"
      cancel-text="取消"
      confirm-variant="danger"
      @confirm="deleteTask"
    />

    <!-- 测试结果对话框 -->
    <TestResultDialog
      v-model="showTestResult"
      :result="testResult"
    />

    <!-- 任务详情对话框 -->
    <TaskDetailDialog
      v-model="showTaskDetail"
      :task="selectedTask"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
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

// Tab 状态（默认显示执行记录）
const activeTab = ref('executions')

// 统计数据
const stats = ref<StatsResponse>({
  total_users: 0,
  total_tasks: 0,
  active_tasks: 0,
  today_executions: 0,
  success_rate: 0,
  recent_users: []
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
  keyword: ''
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
      ...filters.value
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

// 上一页
const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
    loadTasks()
  }
}

// 下一页
const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
    loadTasks()
  }
}

onMounted(() => {
  loadStats()
  loadTasks()
})
</script>
