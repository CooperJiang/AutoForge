<template>
  <Dialog
    :model-value="modelValue"
    :title="mode === 'import' ? '导入工作流' : '导出工作流'"
    :width="700"
    @update:model-value="$emit('update:modelValue', $event)"
    @confirm="handleConfirm"
    @cancel="handleCancel"
  >
    <!-- 导入模式 -->
    <div v-if="mode === 'import'" class="space-y-4">
      <!-- JSON编辑区 -->
      <div class="space-y-3">
        <div class="flex items-center justify-between">
          <label class="text-sm font-medium text-text-secondary">工作流JSON：</label>
          <div class="flex items-center gap-2">
            <BaseButton size="sm" variant="ghost" @click="triggerFileInput">
              <Upload class="w-3.5 h-3.5 mr-1.5" />
              上传文件
            </BaseButton>
            <input
              ref="fileInputRef"
              type="file"
              accept="application/json"
              class="hidden"
              @change="handleFileSelect"
            />
          </div>
        </div>

        <textarea
          v-model="jsonText"
          placeholder="粘贴或上传工作流JSON..."
          class="w-full h-96 px-3 py-2 border-2 border-border-primary rounded-lg focus:outline-none focus:border-primary font-mono text-xs resize-none"
        />

        <div v-if="jsonError" class="p-3 bg-red-50 border border-red-200 rounded-lg">
          <p class="text-xs text-red-700">{{ jsonError }}</p>
        </div>
      </div>

      <!-- 导入提示 -->
      <div class="bg-amber-50 border border-amber-200 rounded-lg p-3">
        <p class="text-xs text-amber-800">
          ⚠️ 导入将覆盖当前工作流内容，请确保已保存重要数据
        </p>
      </div>
    </div>

    <!-- 导出模式 -->
    <div v-if="mode === 'export'" class="space-y-4">
      <!-- JSON编辑区 -->
      <div class="space-y-3">
        <div class="flex items-center justify-between">
          <label class="text-sm font-medium text-text-secondary">编辑JSON内容：</label>
          <div class="flex items-center gap-2">
            <BaseButton size="sm" variant="ghost" @click="copyToClipboard">
              <Copy class="w-3.5 h-3.5 mr-1.5" />
              复制
            </BaseButton>
            <BaseButton size="sm" variant="ghost" @click="handleDownload">
              <Download class="w-3.5 h-3.5 mr-1.5" />
              下载
            </BaseButton>
          </div>
        </div>

        <textarea
          v-model="editableJson"
          class="w-full h-96 px-3 py-2 border-2 border-border-primary rounded-lg focus:outline-none focus:border-primary font-mono text-xs resize-none"
        />

        <div v-if="exportError" class="p-3 bg-red-50 border border-red-200 rounded-lg">
          <p class="text-xs text-red-700">{{ exportError }}</p>
        </div>
      </div>

      <!-- 导出提示 -->
      <div class="bg-primary-light border border-primary rounded-lg p-3">
        <p class="text-xs text-primary mb-2">
          💡 <strong>提示：</strong>
        </p>
        <ul class="text-xs text-primary space-y-1 ml-4">
          <li>• 你可以直接编辑JSON内容，删除敏感信息（如token、邮箱等）</li>
          <li>• 编辑后点击"复制"或"下载"按钮即可分享</li>
          <li>• 建议删除包含敏感信息的字段值，保留字段名方便他人填写</li>
        </ul>
      </div>
    </div>

    <template #footer>
      <div class="flex justify-end gap-3">
        <BaseButton variant="ghost" @click="handleCancel">
          {{ mode === 'export' ? '关闭' : '取消' }}
        </BaseButton>
        <BaseButton v-if="mode === 'import'" @click="handleConfirm">
          导入
        </BaseButton>
      </div>
    </template>
  </Dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { Upload, Copy, Download } from 'lucide-vue-next'
import Dialog from '@/components/Dialog'
import BaseButton from '@/components/BaseButton'
import { message } from '@/utils/message'

interface Props {
  modelValue: boolean
  mode: 'import' | 'export'
  workflowData?: any
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'import': [data: any]
}>()

// 导入相关
const jsonText = ref('')
const jsonError = ref('')
const fileInputRef = ref<HTMLInputElement>()

// 导出相关
const editableJson = ref('')
const exportError = ref('')

// 触发文件选择
const triggerFileInput = () => {
  fileInputRef.value?.click()
}

// 文件选择
const handleFileSelect = async (e: Event) => {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (!file) return

  try {
    const text = await file.text()
    // 格式化JSON
    const data = JSON.parse(text)
    jsonText.value = JSON.stringify(data, null, 2)
    jsonError.value = ''
  } catch (error: any) {
    jsonError.value = '文件格式错误: ' + (error.message || '无法解析JSON')
    message.error('文件读取失败')
  }
}

// 复制到剪贴板
const copyToClipboard = async () => {
  try {
    // 先验证JSON格式
    JSON.parse(editableJson.value)
    await navigator.clipboard.writeText(editableJson.value)
    message.success('已复制到剪贴板')
    exportError.value = ''
  } catch (error: any) {
    if (error instanceof SyntaxError) {
      exportError.value = 'JSON格式错误: ' + error.message
      message.error('JSON格式不正确，请检查')
    } else {
      message.error('复制失败')
    }
  }
}

// 下载文件
const handleDownload = () => {
  try {
    // 验证JSON格式
    const data = JSON.parse(editableJson.value)
    const blob = new Blob([editableJson.value], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `${data.name || 'workflow'}.json`
    a.click()
    URL.revokeObjectURL(url)
    message.success('文件已下载')
    exportError.value = ''
  } catch (error: any) {
    exportError.value = 'JSON格式错误: ' + (error.message || '无法解析')
    message.error('JSON格式不正确，请检查')
  }
}

// 验证JSON
const validateJson = (text: string): any => {
  try {
    const data = JSON.parse(text)
    if (!data.nodes || !data.edges) {
      throw new Error('JSON格式不正确，缺少nodes或edges字段')
    }
    return data
  } catch (error: any) {
    throw new Error(error.message || 'JSON格式错误')
  }
}

// 确认导入
const handleConfirm = () => {
  try {
    if (!jsonText.value.trim()) {
      message.error('请粘贴或上传JSON内容')
      return
    }

    const data = validateJson(jsonText.value)
    jsonError.value = ''
    emit('import', data)
    emit('update:modelValue', false)
  } catch (error: any) {
    jsonError.value = error.message
    message.error('导入失败: ' + error.message)
  }
}

// 取消操作
const handleCancel = () => {
  emit('update:modelValue', false)
}

// 监听modelValue变化，初始化数据
watch(() => props.modelValue, (newVal) => {
  if (newVal) {
    if (props.mode === 'import') {
      jsonText.value = ''
      jsonError.value = ''
    } else {
      // 导出模式：初始化可编辑的JSON
      editableJson.value = props.workflowData
        ? JSON.stringify(props.workflowData, null, 2)
        : ''
      exportError.value = ''
    }
  }
})
</script>
