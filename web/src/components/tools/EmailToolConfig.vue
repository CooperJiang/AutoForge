<template>
  <div class="space-y-4">
    <div class="bg-primary-light border-l-4 border-border-focus p-3">
      <p class="text-sm text-primary">
        <svg class="inline-block w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
          <path
            fill-rule="evenodd"
            d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z"
            clip-rule="evenodd"
          />
        </svg>
        邮件发送使用系统配置，只需填写收件人和邮件内容。支持使用变量引用前置节点的输出数据
      </p>
    </div>

    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2">
        收件人 <span class="text-red-500">*</span>
      </label>
      <BaseInput
        v-model="localConfig.to"
        placeholder="recipient@example.com, another@example.com"
        @update:model-value="emitUpdate"
      />
      <p class="text-xs text-text-tertiary mt-1">多个收件人用逗号分隔</p>
    </div>

    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2"> 抄送人 </label>
      <BaseInput
        v-model="localConfig.cc"
        placeholder="cc@example.com"
        @update:model-value="emitUpdate"
      />
    </div>

    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2">
        邮件主题 <span class="text-red-500">*</span>
      </label>
      <VariableSelector
        v-model="localConfig.subject"
        placeholder="定时任务执行通知"
        :previous-nodes="previousNodes"
        :env-vars="formattedEnvVars"
        @update:model-value="emitUpdate"
      />
    </div>

    <div>
      <label
        class="block text-sm font-medium text-text-secondary mb-2 flex items-center justify-between"
      >
        <span>邮件正文 <span class="text-red-500">*</span></span>
        <button
          type="button"
          @click="showVariableHelper = !showVariableHelper"
          class="text-xs text-primary hover:text-primary"
        >
          {{ showVariableHelper ? '隐藏' : '显示' }}变量助手
        </button>
      </label>

      <div
        v-if="showVariableHelper"
        class="mb-2 p-3 bg-bg-hover rounded-lg border border-border-primary"
      >
        <div class="text-xs font-semibold text-text-secondary mb-2">可用变量：</div>
        <div class="space-y-2">
          <div v-if="previousNodes && previousNodes.length > 0">
            <div class="text-xs text-text-secondary mb-1">前置节点输出：</div>
            <div class="space-y-2">
              <div
                v-for="node in previousNodes"
                :key="node.id"
                class="text-xs border border-border-primary rounded p-2"
              >
                <div class="font-semibold text-text-primary mb-1">{{ node.name }}</div>
                <div class="text-text-tertiary text-[10px] mb-2">ID: {{ node.id }}</div>

                <div class="text-[10px] text-text-secondary mb-1">常见字段（点击插入）：</div>
                <div class="flex flex-wrap gap-1">
                  <button
                    v-for="field in getCommonFields(node.type, node.toolCode)"
                    :key="field.name"
                    type="button"
                    @click="insertFieldVariable(node.id, field.name)"
                    class="px-1.5 py-0.5 bg-primary-light hover:bg-primary-light text-primary rounded text-[10px] font-mono"
                    :title="field.description"
                  >
                    {{ field.name }}
                  </button>
                </div>

                <div class="mt-2 text-[10px] text-text-secondary">
                  也可手动输入：
                  <button
                    type="button"
                    @click="insertNodeVariable(node.id)"
                    class="font-mono text-primary hover:underline"
                  >
                    <span v-text="getNodeVariableText(node.id)"></span>
                  </button>
                </div>
              </div>
            </div>
          </div>

          <div v-if="formattedEnvVars && formattedEnvVars.length > 0">
            <div class="text-xs text-text-secondary mb-1">环境变量：</div>
            <div class="space-y-1">
              <div v-for="envVar in formattedEnvVars" :key="envVar.key" class="text-xs">
                <button
                  type="button"
                  @click="insertEnvVariable(envVar.key)"
                  class="font-mono text-primary hover:text-primary hover:underline"
                >
                  <span v-text="getEnvVariableText(envVar.key)"></span>
                </button>
                <span class="text-text-tertiary ml-1">- {{ envVar.description }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <textarea
        ref="bodyTextareaRef"
        v-model="localConfig.body"
        @input="emitUpdate"
        class="w-full px-3 py-1.5 text-sm text-text-primary bg-bg-primary border-2 border-border-primary rounded-md transition-all duration-200 focus:border-border-focus focus:ring-2 focus:ring-primary-light focus:outline-none hover:border-border-secondary placeholder:text-text-placeholder font-mono"
        rows="8"
        placeholder="尊敬的用户，您好！&#10;&#10;您的邮件内容...&#10;&#10;示例：执行结果：{{nodes.node_xxx.message}}"
      />
      <div class="space-y-1 mt-2">
        <p class="text-xs text-amber-600">💡 <strong>避免被拦截的建议：</strong></p>
        <ul class="text-xs text-text-secondary ml-4 space-y-0.5">
          <li>• 使用完整的邮件格式（称呼、正文、签名）</li>
          <li>• 说明邮件来源和目的</li>
          <li>• 验证码邮件需包含有效期、安全提示</li>
          <li>• 避免纯数字或过于简短的内容</li>
        </ul>
      </div>
    </div>

    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2"> 内容类型 </label>
      <BaseSelect
        v-model="localConfig.content_type"
        :options="contentTypeOptions"
        @update:model-value="emitUpdate"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import BaseInput from '@/components/BaseInput'
import BaseSelect from '@/components/BaseSelect'
import VariableSelector from '@/components/VariableSelector'

interface Props {
  config: Record<string, any>
  previousNodes?: Array<{ id: string; name: string; type: string }>
  envVars?: Array<{ key: string; value: string; description?: string }>
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:config': [config: Record<string, any>]
}>()

const localConfig = ref({
  to: '',
  cc: '',
  subject: '',
  body: '',
  content_type: 'html',
  ...props.config,
})

const contentTypeOptions = [
  { label: 'HTML', value: 'html' },
  { label: '纯文本', value: 'plain' },
]

const showVariableHelper = ref(false)
const bodyTextareaRef = ref<HTMLTextAreaElement>()

// 格式化环境变量
const formattedEnvVars = computed(() => {
  return props.envVars || []
})

// 获取节点变量文本
const getNodeVariableText = (nodeId: string) => {
  return `{{nodes.${nodeId}.*}}`
}

// 获取环境变量文本
const getEnvVariableText = (key: string) => {
  return `{{env.${key}}}`
}

// 获取节点的常见字段
const getCommonFields = (nodeType: string, toolCode?: string) => {
  // HTTP 请求的常见返回字段
  if (toolCode === 'http_request') {
    return [
      { name: 'status', description: 'HTTP 状态码' },
      { name: 'message', description: '返回消息' },
      { name: 'success', description: '成功标识' },
      { name: 'data', description: '返回数据' },
      { name: 'code', description: '业务代码' },
      { name: 'error', description: '错误信息' },
    ]
  }

  // 邮件发送的常见返回字段
  if (toolCode === 'email_sender') {
    return [
      { name: 'success', description: '发送成功' },
      { name: 'message_id', description: '消息ID' },
      { name: 'error', description: '错误信息' },
    ]
  }

  // 健康检查的常见返回字段
  if (toolCode === 'health_checker') {
    return [
      { name: 'message', description: '状态消息' },
      { name: 'url', description: '检测URL' },
      { name: 'status_code', description: 'HTTP状态码' },
      { name: 'response_time', description: '响应时间(毫秒)' },
      { name: 'response_body', description: '接口原始返回内容' },
      { name: 'body_size', description: '返回体大小' },
      { name: 'headers', description: '响应头' },
      { name: 'ssl', description: 'SSL证书信息' },
      { name: 'issues', description: '检测到的问题' },
      { name: 'warnings', description: '警告信息' },
    ]
  }

  // 条件节点的返回字段
  if (nodeType === 'condition') {
    return [
      { name: 'result', description: '条件判断结果 (true/false)' },
      { name: 'message', description: '判断说明' },
    ]
  }

  // 通用字段
  return [
    { name: 'success', description: '成功标识' },
    { name: 'message', description: '消息' },
    { name: 'data', description: '数据' },
    { name: 'error', description: '错误' },
  ]
}

// 插入字段变量
const insertFieldVariable = (nodeId: string, fieldName: string) => {
  insertToBody(`{{nodes.${nodeId}.${fieldName}}}`)
}

// 插入节点变量
const insertNodeVariable = (nodeId: string) => {
  insertToBody(`{{nodes.${nodeId}.`)
}

// 插入环境变量
const insertEnvVariable = (key: string) => {
  insertToBody(`{{env.${key}}}`)
}

// 插入变量到邮件正文
const insertToBody = (text: string) => {
  const textarea = bodyTextareaRef.value
  if (!textarea) return

  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const currentValue = localConfig.value.body || ''

  localConfig.value.body = currentValue.substring(0, start) + text + currentValue.substring(end)
  emitUpdate()

  // 恢复光标位置
  setTimeout(() => {
    textarea.focus()
    const newPos = start + text.length
    textarea.setSelectionRange(newPos, newPos)
  }, 0)
}

watch(
  () => props.config,
  (newVal) => {
    localConfig.value = { ...localConfig.value, ...newVal }
  },
  { deep: true }
)

const emitUpdate = () => {
  emit('update:config', localConfig.value)
}
</script>
