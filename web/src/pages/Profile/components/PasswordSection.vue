<template>
  <div class="bg-white rounded-2xl shadow-lg border border-slate-200 overflow-hidden">
    <div class="border-b border-slate-200 px-6 py-4">
      <h3 class="text-lg font-semibold text-slate-900">修改密码</h3>
      <p class="text-sm text-slate-600">设置一个安全性更高的密码来保护您的账户</p>
    </div>

    <div class="p-6">
      <form @submit.prevent="handleSubmit" class="space-y-4">
        <BaseInput
          v-model="form.oldPassword"
          label="当前密码"
          type="password"
          placeholder="输入当前密码"
          required
        />
        <BaseInput
          v-model="form.newPassword"
          label="新密码"
          type="password"
          placeholder="输入新密码（6-20位）"
          required
        />
        <BaseInput
          v-model="form.confirmPassword"
          label="确认新密码"
          type="password"
          placeholder="再次输入新密码"
          required
        />
        <div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
          <div class="flex items-start gap-3">
            <div class="w-8 h-8 rounded-lg bg-blue-100 flex items-center justify-center flex-shrink-0">
              <svg class="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <p class="text-sm text-blue-900 flex-1">
              <strong>提示：</strong>修改密码成功后，您需要重新登录。
            </p>
          </div>
        </div>
        <div class="flex justify-end gap-3">
          <BaseButton
            type="button"
            variant="secondary"
            @click="handleReset"
          >
            重置
          </BaseButton>
          <BaseButton
            type="submit"
            variant="primary"
            :disabled="changing"
          >
            {{ changing ? '修改中...' : '修改密码' }}
          </BaseButton>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import BaseButton from '@/components/BaseButton'
import BaseInput from '@/components/BaseInput'
import SecureStorage, { STORAGE_KEYS } from '@/utils/storage'
import { message } from '@/utils/message'
import * as userApi from '@/api/user'

const router = useRouter()

const form = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})
const changing = ref(false)

const handleReset = () => {
  form.value = {
    oldPassword: '',
    newPassword: '',
    confirmPassword: ''
  }
}

const handleSubmit = async () => {
  if (!form.value.oldPassword || !form.value.newPassword || !form.value.confirmPassword) {
    message.error('请填写完整信息')
    return
  }

  if (form.value.newPassword.length < 6 || form.value.newPassword.length > 20) {
    message.error('新密码长度必须在6-20位之间')
    return
  }

  if (form.value.newPassword !== form.value.confirmPassword) {
    message.error('两次输入的新密码不一致')
    return
  }

  changing.value = true
  try {
    await userApi.changePassword({
      oldPassword: form.value.oldPassword,
      newPassword: form.value.newPassword
    })

    message.success('密码修改成功，请重新登录')
    handleReset()

    // 延迟后退出登录
    setTimeout(() => {
      SecureStorage.removeItem(STORAGE_KEYS.AUTH_TOKEN)
      SecureStorage.removeItem(STORAGE_KEYS.AUTH_USER)
      router.push('/auth')
    }, 1500)
  } catch (error: any) {
    message.error(error.response?.data?.message || '修改失败')
  } finally {
    changing.value = false
  }
}
</script>
