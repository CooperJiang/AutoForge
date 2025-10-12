<template>
  <div>
    <h4 class="text-sm font-semibold text-slate-900 mb-4 flex items-center gap-2">
      <div class="w-8 h-8 rounded-lg bg-blue-50 flex items-center justify-center">
        <User class="w-4 h-4 text-blue-600"></User>
      </div>
      修改用户名
    </h4>
    <form @submit.prevent="handleSubmit" class="space-y-4">
      <div>
        <label class="block text-sm font-medium text-slate-700 mb-2">当前用户名</label>
        <div class="px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-lg text-slate-900">
          {{ userName }}
        </div>
      </div>
      <BaseInput
        v-model="form.username"
        label="新用户名"
        placeholder="输入新用户名（2-20个字符）"
      />
      <div class="flex justify-end gap-3">
        <BaseButton
          type="button"
          variant="secondary"
          @click="handleReset"
          size="sm"
        >
          重置
        </BaseButton>
        <BaseButton
          type="submit"
          variant="primary"
          :disabled="updating"
          size="sm"
        >
          {{ updating ? '保存中...' : '保存修改' }}
        </BaseButton>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { User } from 'lucide-vue-next'
import BaseButton from '@/components/BaseButton'
import BaseInput from '@/components/BaseInput'
import { message } from '@/utils/message'
import * as userApi from '@/api/user'

interface Props {
  userName: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
  update: [username: string]
}>()

const form = ref({
  username: ''
})
const updating = ref(false)

const handleReset = () => {
  form.value.username = ''
}

const handleSubmit = async () => {
  if (!form.value.username.trim()) {
    message.error('请输入用户名')
    return
  }

  if (form.value.username === props.userName) {
    message.warning('用户名未修改')
    return
  }

  updating.value = true
  try {
    await userApi.updateProfile({ username: form.value.username })
    message.success('用户名修改成功')
    emit('update', form.value.username)
    form.value.username = ''
  } catch (error: any) {
    message.error(error.response?.data?.message || '修改失败')
  } finally {
    updating.value = false
  }
}
</script>
