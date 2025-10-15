<template>
  <Dialog
    :model-value="visible"
    title="执行工作流"
    max-width="max-w-2xl"
    @update:model-value="(value) => !value && $emit('close')"
  >
    <div class="space-y-4">
      
      <div class="bg-bg-tertiary rounded-lg p-3">
        <div class="text-sm font-medium text-text-primary mb-1">
          {{ workflow?.name }}
        </div>
        <div class="text-xs text-text-secondary">
          {{ workflow?.description || '暂无描述' }}
        </div>
      </div>

      
      <div v-if="hasParams" class="text-sm text-text-secondary">
        <div class="flex items-center gap-2 mb-2">
          <Info class="w-4 h-4 text-primary" />
          <span>该工作流需要以下参数：</span>
        </div>
      </div>

      
      <div v-if="hasParams" class="space-y-3">
        <div v-for="param in externalParams" :key="param.key" class="space-y-1.5">
          <label class="block text-sm font-medium text-text-primary">
            {{ param.key }}
            <span v-if="param.required" class="text-red-500 ml-0.5">*</span>
          </label>
          <div v-if="param.description" class="text-xs text-text-tertiary">
            {{ param.description }}
          </div>

          
          <BaseInput
            v-if="param.type === 'string'"
            v-model="paramValues[param.key]"
            :placeholder="param.example ? `例如：${param.example}` : `请输入${param.key}`"
            class="w-full"
          />
          <BaseInput
            v-else-if="param.type === 'number'"
            v-model="paramValues[param.key]"
            type="number"
            :placeholder="param.example ? `例如：${param.example}` : `请输入${param.key}`"
            class="w-full"
          />
          <textarea
            v-else
            v-model="paramValues[param.key]"
            :placeholder="param.example ? `例如：${param.example}` : `请输入${param.key}`"
            rows="3"
            class="w-full px-3 py-2 bg-bg-primary border border-border-primary rounded-lg text-sm text-text-primary placeholder-text-tertiary focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent resize-none"
          />
        </div>
      </div>

      
      <div v-else class="text-center py-8 text-text-tertiary">
        <Play class="w-12 h-12 mx-auto mb-2 opacity-50" />
        <p>该工作流无需参数，可直接执行</p>
      </div>
    </div>

    
    <div class="space-y-1.5">
      <label class="block text-sm font-medium text-text-primary">会话ID（可选）</label>
      <BaseInput v-model="sessionId" placeholder="例如：user-123，或留空使用默认" class="w-full" />
      <div class="text-xs text-text-tertiary">
        用于对话记忆的分组键；留空时由后端根据用户或使用 global
      </div>
    </div>

    <template #footer>
      <div class="flex justify-end gap-3">
        <BaseButton variant="secondary" @click="$emit('close')"> 取消 </BaseButton>
        <BaseButton variant="primary" @click="handleExecute" :disabled="!isValid">
          <Play class="w-4 h-4 mr-1" />
          执行
        </BaseButton>
      </div>
    </template>
  </Dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { Info, Play } from 'lucide-vue-next'
import Dialog from '@/components/Dialog'
import BaseButton from '@/components/BaseButton'
import BaseInput from '@/components/BaseInput'
import type { Workflow, WorkflowNode } from '@/types/workflow'

interface Props {
  visible: boolean
  workflow: Workflow | null
}

const props = defineProps<Props>()

const emit = defineEmits<{
  close: []
  execute: [params: Record<string, any>]
}>()

// 参数值
const paramValues = ref<Record<string, any>>({})
const sessionId = ref('')

// 提取外部触发器参数
const externalParams = computed(() => {
  if (!props.workflow) return []

  // 找到第一个节点
  const firstNode = props.workflow.nodes.find((node: WorkflowNode) => {
    // 检查是否有入边
    const hasIncomingEdge = props.workflow!.edges.some((edge) => edge.target === node.id)
    return !hasIncomingEdge
  })

  // 如果第一个节点是 external_trigger，提取参数
  if (firstNode?.type === 'external_trigger' && firstNode.config?.params) {
    return firstNode.config.params as Array<{
      key: string
      type: string
      required: boolean
      description?: string
      example?: any
      defaultValue?: any
    }>
  }

  return []
})

// 是否有参数
const hasParams = computed(() => externalParams.value.length > 0)

// 验证参数是否有效
const isValid = computed(() => {
  if (!hasParams.value) return true

  // 检查必填参数是否都已填写
  return externalParams.value.every((param) => {
    if (!param.required) return true
    const value = paramValues.value[param.key]
    return value !== undefined && value !== null && value !== ''
  })
})

// 初始化参数值
watch(
  () => props.visible,
  (visible) => {
    if (visible) {
      // 重置参数值，使用默认值
      paramValues.value = {}
      externalParams.value.forEach((param) => {
        if (param.defaultValue !== undefined && param.defaultValue !== null) {
          paramValues.value[param.key] = param.defaultValue
        } else {
          paramValues.value[param.key] = ''
        }
      })
    }
  },
  { immediate: true }
)

const handleExecute = () => {
  if (!isValid.value) return

  // 转换参数值类型
  const params: Record<string, any> = {}
  externalParams.value.forEach((param) => {
    const value = paramValues.value[param.key]

    // 跳过空值（非必填）
    if (!param.required && (value === '' || value === null || value === undefined)) {
      return
    }

    // 类型转换
    if (param.type === 'number') {
      params[param.key] = Number(value)
    } else if (param.type === 'boolean') {
      params[param.key] = value === 'true' || value === true
    } else {
      params[param.key] = value
    }
  })

  if (sessionId.value) {
    params.session_id = sessionId.value
  }
  emit('execute', params)
}
</script>
