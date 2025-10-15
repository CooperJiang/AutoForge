<template>
  <div>
    <div class="grid grid-cols-1 xl:grid-cols-12 gap-6">
      <!-- 左侧菜单 -->
      <div class="xl:col-span-3">
        <ProfileSidebar
          :active-tab="activeTab"
          :user-name="userName"
          :user-email="userEmail"
          :role-text="roleText"
          @tab-change="activeTab = $event"
        />
      </div>

      <!-- 右侧内容区 -->
      <div class="xl:col-span-9">
        <!-- 个人资料 -->
        <ProfileSection
          v-show="activeTab === 'profile'"
          :user-name="userName"
          :user-email="userEmail"
          @update-username="handleUpdateUsername"
          @update-email="handleUpdateEmail"
        />

        <!-- 修改密码 -->
        <PasswordSection v-show="activeTab === 'password'" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import ProfileSidebar from './components/ProfileSidebar.vue'
import ProfileSection from './components/ProfileSection.vue'
import PasswordSection from './components/PasswordSection.vue'
import SecureStorage, { STORAGE_KEYS } from '@/utils/storage'
import { message } from '@/utils/message'

const router = useRouter()

// 角色常量
const USER_ROLE = {
  SUPER_ADMIN: 1,
  ADMIN: 2,
  USER: 3,
}

// 当前激活的标签
const activeTab = ref('profile')

// 从 SecureStorage 获取用户信息
const user = computed(() =>
  SecureStorage.getItem<{
    username: string
    id: string
    email: string
    role: number
  }>(STORAGE_KEYS.AUTH_USER)
)

// 用户名
const userName = computed(() => user.value?.username || '未知用户')

// 邮箱
const userEmail = computed(() => user.value?.email || '未设置')

// 角色文本
const roleText = computed(() => {
  const role = user.value?.role
  if (role === USER_ROLE.SUPER_ADMIN) return '超级管理员'
  if (role === USER_ROLE.ADMIN) return '管理员'
  return '普通用户'
})

// 更新用户名
const handleUpdateUsername = (username: string) => {
  const currentUser = user.value
  if (currentUser) {
    currentUser.username = username
    SecureStorage.setItem(STORAGE_KEYS.AUTH_USER, currentUser)
  }
}

// 更新邮箱
const handleUpdateEmail = (data: { email: string; code: string }) => {
  const currentUser = user.value
  if (currentUser) {
    currentUser.email = data.email
    SecureStorage.setItem(STORAGE_KEYS.AUTH_USER, currentUser)
  }
}

// 检查是否已登录
onMounted(() => {
  if (!user.value) {
    message.warning('请先登录')
    router.push('/auth')
  }
})
</script>
