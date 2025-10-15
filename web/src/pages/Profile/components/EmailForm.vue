<template>
  <div>
    <h4 class="text-sm font-semibold text-text-primary mb-4 flex items-center gap-2">
      <div class="w-8 h-8 rounded-lg bg-primary-light flex items-center justify-center">
        <Mail class="w-4 h-4 text-primary"></Mail>
      </div>
      修改邮箱
    </h4>
    <form @submit.prevent="handleSubmit" class="space-y-4">
      <div>
        <label class="block text-sm font-medium text-text-secondary mb-2">当前邮箱</label>
        <div
          class="px-4 py-2.5 bg-bg-hover border border-border-primary rounded-lg text-text-primary"
        >
          {{ userEmail }}
        </div>
      </div>
      <BaseInput
        v-model="form.email"
        label="新邮箱地址"
        type="email"
        placeholder="输入新的邮箱地址"
      />
      <div>
        <label class="block text-sm font-medium text-text-secondary mb-2">邮箱验证码</label>
        <div class="flex gap-3">
          <BaseInput v-model="form.code" placeholder="输入6位验证码" class="flex-1" />
          <BaseButton
            type="button"
            variant="secondary"
            @click="handleSendCode"
            :disabled="sendingCode || countdown > 0"
            class="whitespace-nowrap"
            size="sm"
          >
            {{ countdown > 0 ? `${countdown}秒` : sendingCode ? '发送中...' : '发送验证码' }}
          </BaseButton>
        </div>
      </div>
      <div class="flex justify-end gap-3">
        <BaseButton type="button" variant="secondary" @click="handleReset" size="sm">
          重置
        </BaseButton>
        <BaseButton type="submit" variant="primary" :disabled="updating" size="sm">
          {{ updating ? '保存中...' : '保存修改' }}
        </BaseButton>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Mail } from 'lucide-vue-next'
import BaseButton from '@/components/BaseButton'
import BaseInput from '@/components/BaseInput'
import { message } from '@/utils/message'
import * as userApi from '@/api/user'

interface Props {
  userEmail: string
}

defineProps<Props>()

const emit = defineEmits<{
  update: [data: { email: string; code: string }]
}>()

const form = ref({
  email: '',
  code: '',
})
const updating = ref(false)
const sendingCode = ref(false)
const countdown = ref(0)

const handleReset = () => {
  form.value.email = ''
  form.value.code = ''
}

const handleSendCode = async () => {
  if (!form.value.email.trim()) {
    message.error('请输入邮箱地址')
    return
  }

  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.value.email)) {
    message.error('请输入有效的邮箱地址')
    return
  }

  sendingCode.value = true
  try {
    await userApi.sendChangeEmailCode(form.value.email)
    message.success('验证码已发送，请查收邮箱')

    // 启动倒计时
    countdown.value = 60
    const timer = setInterval(() => {
      countdown.value--
      if (countdown.value <= 0) {
        clearInterval(timer)
      }
    }, 1000)
  } catch (error: any) {
    message.error(error.response?.data?.message || '发送失败')
  } finally {
    sendingCode.value = false
  }
}

const handleSubmit = async () => {
  if (!form.value.email.trim() || !form.value.code.trim()) {
    message.error('请填写完整信息')
    return
  }

  updating.value = true
  try {
    await userApi.updateProfile({
      email: form.value.email,
      code: form.value.code,
    })
    message.success('邮箱修改成功')
    emit('update', { email: form.value.email, code: form.value.code })
    handleReset()
  } catch (error: any) {
    message.error(error.response?.data?.message || '修改失败')
  } finally {
    updating.value = false
  }
}
</script>
