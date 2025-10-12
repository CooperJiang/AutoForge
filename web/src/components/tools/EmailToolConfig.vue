<template>
  <div class="space-y-4">
    <div class="bg-blue-50 border-l-4 border-blue-400 p-3">
      <p class="text-sm text-blue-700">
        <svg class="inline-block w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd"/>
        </svg>
        邮件发送使用系统配置，只需填写收件人和邮件内容
      </p>
    </div>

    <div>
      <label class="block text-sm font-medium text-slate-700 mb-2">
        收件人 <span class="text-red-500">*</span>
      </label>
      <BaseInput
        v-model="localConfig.to"
        placeholder="recipient@example.com, another@example.com"
        @update:model-value="emitUpdate"
      />
      <p class="text-xs text-slate-500 mt-1">多个收件人用逗号分隔</p>
    </div>

    <div>
      <label class="block text-sm font-medium text-slate-700 mb-2">
        抄送人
      </label>
      <BaseInput
        v-model="localConfig.cc"
        placeholder="cc@example.com"
        @update:model-value="emitUpdate"
      />
    </div>

    <div>
      <label class="block text-sm font-medium text-slate-700 mb-2">
        邮件主题 <span class="text-red-500">*</span>
      </label>
      <BaseInput
        v-model="localConfig.subject"
        placeholder="定时任务执行通知"
        @update:model-value="emitUpdate"
      />
    </div>

    <div>
      <label class="block text-sm font-medium text-slate-700 mb-2">
        邮件正文 <span class="text-red-500">*</span>
      </label>
      <textarea
        v-model="localConfig.body"
        @input="emitUpdate"
        class="w-full px-3 py-2 border-2 border-slate-200 rounded-lg focus:outline-none focus:border-emerald-500 font-mono text-sm"
        rows="8"
        placeholder="尊敬的用户，您好！&#10;&#10;您的邮件内容..."
      />
      <div class="space-y-1 mt-2">
        <p class="text-xs text-amber-600">
          💡 <strong>避免被拦截的建议：</strong>
        </p>
        <ul class="text-xs text-slate-600 ml-4 space-y-0.5">
          <li>• 使用完整的邮件格式（称呼、正文、签名）</li>
          <li>• 说明邮件来源和目的</li>
          <li>• 验证码邮件需包含有效期、安全提示</li>
          <li>• 避免纯数字或过于简短的内容</li>
        </ul>
      </div>
    </div>

    <div>
      <label class="block text-sm font-medium text-slate-700 mb-2">
        内容类型
      </label>
      <BaseSelect
        v-model="localConfig.content_type"
        :options="contentTypeOptions"
        @update:model-value="emitUpdate"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import BaseInput from '@/components/BaseInput'
import BaseSelect from '@/components/BaseSelect'

interface Props {
  config: Record<string, any>
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
  ...props.config
})

const contentTypeOptions = [
  { label: 'HTML', value: 'html' },
  { label: '纯文本', value: 'plain' }
]

watch(() => props.config, (newVal) => {
  localConfig.value = { ...localConfig.value, ...newVal }
}, { deep: true })

const emitUpdate = () => {
  emit('update:config', localConfig.value)
}
</script>
