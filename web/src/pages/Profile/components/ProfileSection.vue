<template>
  <div class="bg-white rounded-2xl shadow-lg border border-slate-200 overflow-hidden">
    <div class="border-b border-slate-200 px-6 py-4">
      <h3 class="text-lg font-semibold text-slate-900">个人资料</h3>
      <p class="text-sm text-slate-600">管理您的用户名和邮箱信息</p>
    </div>

    <div class="p-8">
      <!-- 修改用户名 -->
      <UsernameForm
        :user-name="userName"
        @update="$emit('update-username', $event)"
      />

      <!-- 分隔线 -->
      <div class="my-8 border-b-2 border-slate-200"></div>

      <!-- 修改邮箱 -->
      <EmailForm
        :user-email="userEmail"
        @update="$emit('update-email', $event)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import UsernameForm from './UsernameForm.vue'
import EmailForm from './EmailForm.vue'

interface Props {
  userName: string
  userEmail: string
}

defineProps<Props>()

defineEmits<{
  'update-username': [username: string]
  'update-email': [data: { email: string; code: string }]
}>()
</script>
