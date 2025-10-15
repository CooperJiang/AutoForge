<template>
  <div class="inline-flex items-center">
    <BaseButton :size="size" :variant="variant" :disabled="busy" @click="handleClick">
      <Play class="w-4 h-4 mr-1" />
      {{ label }}
    </BaseButton>

    
    <ExecuteWithParamsDialog
      v-if="showParamsDialog"
      :visible="showParamsDialog"
      :workflow="workflowData"
      @close="showParamsDialog = false"
      @execute="handleExecuteWithParams"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Play } from 'lucide-vue-next'
import BaseButton from '@/components/BaseButton'
import ExecuteWithParamsDialog from '@/pages/Workflows/components/ExecuteWithParamsDialog.vue'
import { workflowApi } from '@/api/workflow'
import { message } from '@/utils/message'
import { useRouter } from 'vue-router'

interface Props {
  workflowId: string
  workflow?: any
  label?: string
  size?: 'sm' | 'md' | 'lg'
  variant?: 'primary' | 'secondary' | 'ghost'
  navigateToDetail?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  label: '执行',
  size: 'sm',
  variant: 'secondary',
  navigateToDetail: true,
})

const emit = defineEmits<{ (e: 'executed', executionId: string): void }>()

const router = useRouter()
const busy = ref(false)
const showParamsDialog = ref(false)
const workflowData = ref<any | null>(props.workflow || null)

const loadWorkflowIfNeeded = async () => {
  if (!workflowData.value) {
    workflowData.value = await workflowApi.getById(props.workflowId)
  }
}

const needsExternalParams = (wf: any): boolean => {
  if (!wf || !Array.isArray(wf.nodes) || !Array.isArray(wf.edges)) return false
  // 找到起始节点（没有入边的节点）
  const targets = new Set(wf.edges.map((e: any) => e.target))
  const startNodes = wf.nodes.filter((n: any) => !targets.has(n.id))
  const first = startNodes[0]
  if (first && first.type === 'external_trigger') {
    const params = first.config?.params
    return Array.isArray(params) && params.length > 0
  }
  return false
}

const handleExecute = async (params?: Record<string, any>) => {
  busy.value = true
  try {
    const data = await workflowApi.execute(props.workflowId, params ? { params } : undefined)
    message.success('工作流已开始执行')
    emit('executed', data.execution_id)
    if (props.navigateToDetail) {
      router.push(`/workflows/${props.workflowId}/executions/${data.execution_id}`)
    }
  } catch (error: any) {
    console.error('Execute workflow failed:', error)
    message.error(error.response?.data?.message || '执行失败')
  } finally {
    busy.value = false
    showParamsDialog.value = false
  }
}

const handleClick = async () => {
  try {
    await loadWorkflowIfNeeded()
    if (!workflowData.value?.enabled) {
      message.warning('工作流已禁用，请先启用')
      return
    }
    if (needsExternalParams(workflowData.value)) {
      showParamsDialog.value = true
      return
    }
    await handleExecute()
  } catch (error: any) {
    console.error(error)
    message.error('加载工作流失败')
  }
}

const handleExecuteWithParams = async (params: Record<string, any>) => {
  await handleExecute(params)
}
</script>
