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

          <!-- 文件上传 -->
          <div v-if="param.type === 'file'" class="space-y-2">
            <input
              :ref="(el) => (fileInputRefs[param.key] = el)"
              type="file"
              :accept="param.accept || '*'"
              @change="handleFileChange(param.key, $event)"
              class="hidden"
            />
            <div
              @click="triggerFileInput(param.key)"
              class="w-full px-4 py-3 border-2 border-dashed border-border-primary rounded-lg cursor-pointer hover:border-primary hover:bg-primary-light transition-colors"
            >
              <div v-if="!fileNames[param.key]" class="flex items-center justify-center gap-2 text-text-tertiary">
                <Upload class="w-5 h-5" />
                <span class="text-sm">点击选择文件</span>
              </div>
              <div v-else class="flex items-center justify-between">
                <div class="flex items-center gap-2 text-text-primary">
                  <FileIcon class="w-5 h-5 text-primary" />
                  <span class="text-sm">{{ fileNames[param.key] }}</span>
                </div>
                <button
                  @click.stop="clearFile(param.key)"
                  class="text-text-tertiary hover:text-error transition-colors"
                >
                  <X class="w-4 h-4" />
                </button>
              </div>
            </div>
            <div v-if="param.accept || param.maxSize" class="text-xs text-text-tertiary">
              <span v-if="param.accept">支持格式：{{ param.accept }}</span>
              <span v-if="param.accept && param.maxSize"> | </span>
              <span v-if="param.maxSize">最大 {{ param.maxSize }}MB</span>
            </div>
          </div>

          <!-- 文本输入 -->
          <BaseInput
            v-else-if="param.type === 'string'"
            v-model="paramValues[param.key]"
            :placeholder="param.example ? `例如：${param.example}` : `请输入${param.key}`"
            class="w-full"
          />

          <!-- 数字输入 -->
          <BaseInput
            v-else-if="param.type === 'number'"
            v-model="paramValues[param.key]"
            type="number"
            :placeholder="param.example ? `例如：${param.example}` : `请输入${param.key}`"
            class="w-full"
          />

          <!-- 其他类型使用 textarea -->
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
import { Info, Play, Upload, FileIcon, X } from 'lucide-vue-next'
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

// 文件相关
const fileInputRefs = ref<Record<string, any>>({})
const fileNames = ref<Record<string, string>>({})
const selectedFiles = ref<Record<string, File>>({})

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

    // 文件类型检查
    if (param.type === 'file') {
      return selectedFiles.value[param.key] !== undefined
    }

    // 其他类型检查
    const value = paramValues.value[param.key]
    return value !== undefined && value !== null && value !== ''
  })
})

// 触发文件选择
const triggerFileInput = (key: string) => {
  const input = fileInputRefs.value[key]
  if (input) {
    input.click()
  }
}

// 处理文件选择
const handleFileChange = (key: string, event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]

  if (file) {
    selectedFiles.value[key] = file
    fileNames.value[key] = file.name
    paramValues.value[key] = file
  }
}

// 清除文件
const clearFile = (key: string) => {
  delete selectedFiles.value[key]
  delete fileNames.value[key]
  delete paramValues.value[key]

  // 重置 input
  const input = fileInputRefs.value[key]
  if (input) {
    input.value = ''
  }
}

// 初始化参数值
watch(
  () => props.visible,
  (visible) => {
    if (visible) {
      // 重置参数值，使用默认值
      paramValues.value = {}
      selectedFiles.value = {}
      fileNames.value = {}

      externalParams.value.forEach((param) => {
        // 文件类型不设置默认值
        if (param.type === 'file') {
          return
        }

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

  // 检查是否包含文件类型参数
  const hasFileParam = externalParams.value.some((param) => param.type === 'file')

  if (hasFileParam) {
    // 使用 FormData 处理包含文件的请求
    const formData = new FormData()

    // 收集非文件参数
    const jsonParams: Record<string, any> = {}
    externalParams.value.forEach((param) => {
      if (param.type === 'file') {
        const file = selectedFiles.value[param.key]
        if (file) {
          formData.append(param.key, file)
        }
      } else {
        const value = paramValues.value[param.key]

        // 跳过空值（非必填）
        if (!param.required && (value === '' || value === null || value === undefined)) {
          return
        }

        // 类型转换后添加到 jsonParams
        if (param.type === 'number') {
          jsonParams[param.key] = Number(value)
        } else if (param.type === 'boolean') {
          jsonParams[param.key] = value === 'true' || value === true
        } else {
          jsonParams[param.key] = value
        }
      }
    })

    if (Object.keys(jsonParams).length > 0) {
      formData.append('params', JSON.stringify(jsonParams))
    }

    if (sessionId.value) {
      formData.append('webhook_url', '')
    }

    emit('execute', formData)
  } else {
    // 不包含文件，使用普通 JSON 对象
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
}
</script>
