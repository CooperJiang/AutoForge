<template>
  <div>
      <!-- 页面标题 -->
      <div class="flex items-center justify-between mb-6">
        <div>
          <h1 class="text-2xl font-bold text-text-primary mb-1">
            工作流管理
          </h1>
          <p class="text-sm text-text-secondary">
            创建和管理自动化工作流程
          </p>
        </div>
        <BaseButton size="md" @click="router.push('/workflows/create')">
          <Plus class="w-4 h-4 mr-1" />
          创建工作流
        </BaseButton>
      </div>

      <!-- 工作流列表 -->
      <div v-if="!loading && workflows.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5 gap-4">
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
        <div class="text-text-placeholder mb-4">
          <Workflow class="w-16 h-16 mx-auto mb-4" />
          <p class="text-lg">暂无工作流</p>
          <p class="text-sm">点击上方按钮创建第一个工作流</p>
        </div>
      </div>

      <!-- 加载状态 -->
      <div v-else class="flex justify-center items-center py-20">
        <div class="text-text-tertiary">加载中...</div>
      </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Plus, Workflow } from 'lucide-vue-next'
import BaseButton from '@/components/BaseButton'
import WorkflowCard from './components/WorkflowCard.vue'
import { workflowApi } from '@/api/workflow'
import { message } from '@/utils/message'
import type { Workflow as WorkflowType } from '@/types/workflow'

const router = useRouter()
const loading = ref(false)
const workflows = ref<WorkflowType[]>([])

// 加载工作流列表
const loadWorkflows = async () => {
  loading.value = true
  try {
    const data = await workflowApi.list()
    workflows.value = data.items || []
  } catch (error) {
    console.error('Failed to load workflows:', error)
    message.error('加载工作流列表失败')
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
  try {
    await workflowApi.delete(workflow.id)
    message.success('删除成功')
    await loadWorkflows()
  } catch (error) {
    console.error('Delete workflow failed:', error)
    message.error('删除失败')
  }
}

// 执行工作流
const handleExecute = async (workflow: WorkflowType) => {
  if (!workflow.enabled) {
    message.warning('工作流未启用')
    return
  }

  try {
    const data = await workflowApi.execute(workflow.id)
    message.success('工作流已开始执行')
    router.push(`/workflows/${workflow.id}/executions/${data.execution_id}`)
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

onMounted(() => {
  loadWorkflows()
})
</script>
