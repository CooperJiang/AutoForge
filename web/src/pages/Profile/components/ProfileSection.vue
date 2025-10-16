<template>
  <div class="bg-bg-elevated rounded-2xl shadow-lg border border-border-primary overflow-hidden">
    <div class="border-b border-border-primary px-6 py-4">
      <h3 class="text-lg font-semibold text-text-primary">个人资料</h3>
      <p class="text-sm text-text-secondary">管理您的用户名和邮箱信息</p>
    </div>

    <div class="p-8">
      <UsernameForm :user-name="userName" @update="$emit('update-username', $event)" />

      <div class="my-8 border-b-2 border-border-primary"></div>

      <EmailForm :user-email="userEmail" @update="$emit('update-email', $event)" />
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
