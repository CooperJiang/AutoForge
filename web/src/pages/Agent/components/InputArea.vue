<template>
  <div class="border-t border-border-primary bg-bg-secondary p-4">
    <div class="max-w-4xl mx-auto">
      <!-- 配置栏 -->
      <div class="flex items-center gap-4 mb-3 text-sm">
        <!-- 模式选择 -->
        <div class="flex items-center gap-2">
          <label class="text-text-secondary">模式:</label>
          <select
            v-model="config.mode"
            class="px-3 py-1.5 bg-bg-primary border border-border-primary rounded-lg text-text-primary text-sm focus:outline-none focus:ring-2 focus:ring-primary"
          >
            <option value="direct">Direct (边思考边执行)</option>
            <option value="plan">Plan (先规划再执行)</option>
          </select>
        </div>

        <!-- 模型选择 -->
        <div class="flex items-center gap-2">
          <label class="text-text-secondary">模型:</label>
          <select
            v-model="config.model"
            class="px-3 py-1.5 bg-bg-primary border border-border-primary rounded-lg text-text-primary text-sm focus:outline-none focus:ring-2 focus:ring-primary"
          >
            <option value="gpt-4o-mini">GPT-4o Mini</option>
            <option value="gpt-4o">GPT-4o</option>
            <option value="gemini-1.5-flash">Gemini 1.5 Flash</option>
            <option value="gemini-1.5-pro">Gemini 1.5 Pro</option>
          </select>
        </div>
      </div>

      <!-- 输入框 -->
      <div class="flex items-end gap-3">
        <div class="flex-1 relative">
          <textarea
            ref="textareaRef"
            v-model="message"
            @keydown.enter="handleKeydown"
            :disabled="disabled"
            placeholder="输入你的问题或任务..."
            rows="1"
            class="w-full px-4 py-3 pr-12 bg-bg-primary border border-border-primary rounded-xl text-text-primary placeholder-text-tertiary resize-none focus:outline-none focus:ring-2 focus:ring-primary disabled:opacity-50 disabled:cursor-not-allowed"
            style="max-height: 200px"
          />

          <!-- 文件上传按钮 -->
          <button
            class="absolute right-3 bottom-3 p-1.5 rounded-lg hover:bg-bg-tertiary transition-colors"
            @click="handleFileClick"
            :disabled="disabled"
          >
            <Paperclip class="w-5 h-5 text-text-secondary" />
          </button>

          <input
            ref="fileInputRef"
            type="file"
            multiple
            accept="image/*,application/pdf"
            class="hidden"
            @change="handleFileChange"
          />
        </div>

        <!-- 发送按钮 -->
        <button
          @click="handleSend"
          :disabled="disabled || !message.trim()"
          class="flex-shrink-0 p-3 bg-primary text-white rounded-xl hover:bg-primary-dark transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <Send class="w-5 h-5" />
        </button>
      </div>

      <!-- 文件预览 -->
      <div v-if="files.length > 0" class="mt-3 flex flex-wrap gap-2">
        <div
          v-for="(file, index) in files"
          :key="index"
          class="flex items-center gap-2 px-3 py-2 bg-bg-tertiary rounded-lg text-sm"
        >
          <FileIcon class="w-4 h-4 text-text-secondary" />
          <span class="text-text-primary">{{ file.name }}</span>
          <button @click="removeFile(index)" class="text-text-tertiary hover:text-error">
            <X class="w-4 h-4" />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'
import { Send, Paperclip, FileIcon, X } from 'lucide-vue-next'
import type { SendMessageRequest } from '@/api/agent'

interface Props {
  disabled?: boolean
}

defineProps<Props>()

const emit = defineEmits<{
  send: [data: SendMessageRequest]
}>()

const message = ref('')
const files = ref<File[]>([])
const fileInputRef = ref<HTMLInputElement>()
const textareaRef = ref<HTMLTextAreaElement>()

const config = ref({
  mode: 'plan' as 'direct' | 'plan',
  model: 'gpt-4o-mini',
  max_steps: 10,
  temperature: 0.7,
})

// 自动调整 textarea 高度
watch(message, () => {
  nextTick(() => {
    if (textareaRef.value) {
      textareaRef.value.style.height = 'auto'
      textareaRef.value.style.height = textareaRef.value.scrollHeight + 'px'
    }
  })
})

// 处理键盘事件
const handleKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    handleSend()
  }
}

// 发送消息
const handleSend = () => {
  if (!message.value.trim()) return

  emit('send', {
    message: message.value.trim(),
    config: config.value,
  })

  message.value = ''
  files.value = []
}

// 处理文件点击
const handleFileClick = () => {
  fileInputRef.value?.click()
}

// 处理文件选择
const handleFileChange = (e: Event) => {
  const target = e.target as HTMLInputElement
  if (target.files) {
    files.value.push(...Array.from(target.files))
  }
  target.value = ''
}

// 移除文件
const removeFile = (index: number) => {
  files.value.splice(index, 1)
}
</script>



