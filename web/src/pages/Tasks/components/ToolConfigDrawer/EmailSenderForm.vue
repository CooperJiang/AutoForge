<template>
  <div class="space-y-4">
    <!-- 提示 -->
    <div class="bg-primary-light border-l-4 border-primary p-3 mb-4">
      <p class="text-sm text-primary">
        <svg class="inline-block w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
          <path
            fill-rule="evenodd"
            d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z"
            clip-rule="evenodd"
          />
        </svg>
        邮件发送使用系统配置，只需填写收件人和邮件内容
      </p>
    </div>

    <!-- 收件人 -->
    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2">
        收件人 <span class="text-red-500">*</span>
      </label>
      <BaseInput
        v-model="config.to"
        placeholder="recipient@example.com, another@example.com"
        required
      />
      <p class="text-xs text-text-tertiary mt-1">多个收件人用逗号分隔</p>
    </div>

    <!-- 抄送人 -->
    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2"> 抄送人 </label>
      <BaseInput v-model="config.cc" placeholder="cc@example.com" />
    </div>

    <!-- 邮件主题 -->
    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2">
        邮件主题 <span class="text-red-500">*</span>
      </label>
      <BaseInput v-model="config.subject" placeholder="定时任务执行通知" required />
    </div>

    <!-- 邮件正文 -->
    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2">
        邮件正文 <span class="text-red-500">*</span>
      </label>
      <textarea
        v-model="config.body"
        class="w-full px-3 py-2 border-2 border-border-primary rounded-lg focus:outline-none focus:border-primary bg-bg-elevated text-text-primary font-mono text-sm"
        rows="8"
        placeholder="尊敬的用户，您好！&#10;&#10;您正在使用【自动任务系统】进行身份验证，您的验证码为：&#10;&#10;      123456&#10;&#10;验证码有效期为 10 分钟，请勿泄露给他人。&#10;如非本人操作，请忽略此邮件。&#10;&#10;感谢您的使用！&#10;&#10;---&#10;【自动任务系统】&#10;support@yourdomain.com"
        required
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

    <!-- 内容类型 -->
    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2"> 内容类型 </label>
      <BaseSelect v-model="config.content_type" :options="contentTypeOptions" />
    </div>
  </div>
</template>

<script setup lang="ts">
import BaseInput from '@/components/BaseInput'
import BaseSelect from '@/components/BaseSelect'
import type { EmailSenderConfig, SelectOption } from './types'

interface Props {
  config: EmailSenderConfig
  contentTypeOptions: SelectOption[]
}

defineProps<Props>()
</script>
