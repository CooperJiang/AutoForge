<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50 to-purple-50">
    <div class="max-w-7xl mx-auto px-4 py-8">
      <!-- 页面标题 -->
      <div class="flex items-center justify-between mb-6">
        <div>
          <h1 class="text-2xl font-bold text-slate-900 mb-1">
            工作流管理
          </h1>
          <p class="text-sm text-slate-600">
            创建和管理自动化工作流程
          </p>
        </div>
        <BaseButton size="md" @click="router.push('/workflows/create')">
          <Plus class="w-4 h-4 mr-1" />
          创建工作流
        </BaseButton>
      </div>

      <!-- 工作流列表 -->
      <div v-if="!loading && workflows.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <WorkflowCard
          v-for="workflow in workflows"
          :key="workflow.id"
          :workflow="workflow"
          @edit="handleEdit"
          @executions="handleViewExecutions"
          @execute="handleExecute"
          @delete="handleDelete"
          @toggle="handleToggle"
        />
      </div>

      <!-- 空状态 -->
      <div v-else-if="!loading && workflows.length === 0" class="text-center py-20">
        <div class="text-slate-400 mb-4">
          <Workflow class="w-16 h-16 mx-auto mb-4" />
          <p class="text-lg">暂无工作流</p>
          <p class="text-sm">点击上方按钮创建第一个工作流</p>
        </div>
      </div>

      <!-- 加载状态 -->
      <div v-else class="flex justify-center items-center py-20">
        <div class="text-slate-500">加载中...</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Plus, Workflow } from 'lucide-vue-next'
import BaseButton from '@/components/BaseButton'
import WorkflowCard from './components/WorkflowCard.vue'
import type { Workflow as WorkflowType } from '@/types/workflow'

const router = useRouter()
const loading = ref(false)
const workflows = ref<WorkflowType[]>([])

// 加载工作流列表
const loadWorkflows = async () => {
  loading.value = true
  try {
    // TODO: 调用API加载工作流列表
    // const response = await workflowApi.getWorkflowList()
    // workflows.value = response.data

    // 临时Mock数据
    workflows.value = []
  } catch (error) {
    console.error('Failed to load workflows:', error)
  } finally {
    loading.value = false
  }
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
  // TODO: 实现删除逻辑
  console.log('Delete workflow:', workflow.id)
  await loadWorkflows()
}

// 执行工作流
const handleExecute = async (workflow: WorkflowType) => {
  if (!workflow.enabled) {
    console.warn('工作流已禁用')
    return
  }

  try {
    console.log('执行工作流:', workflow.id)
    // TODO: 调用API执行工作流
    // const response = await workflowApi.executeWorkflow(workflow.id)
    // router.push(`/workflows/${workflow.id}/executions/${response.data.executionId}`)
    router.push(`/workflows/${workflow.id}/executions`)
  } catch (error) {
    console.error('工作流执行失败:', error)
  }
}

// 切换工作流状态
const handleToggle = async (workflow: WorkflowType) => {
  // TODO: 实现切换状态逻辑
  console.log('Toggle workflow:', workflow.id)
}

onMounted(() => {
  loadWorkflows()
})
</script>
