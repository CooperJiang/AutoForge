<template>
  <header class="bg-white border-b-2 border-slate-200 flex-shrink-0">
    <div class="max-w-7xl mx-auto px-4 py-2.5">
      <div class="flex justify-between items-center">
        <div class="flex items-center gap-3 cursor-pointer" @click="handleLogoClick">
          <img src="/logo.png" alt="Logo" class="h-10 w-10 object-contain" />
          <div>
            <h1 class="text-lg font-bold text-slate-900">定时任务系统</h1>
            <p class="text-xs text-slate-500">{{ subtitle }}</p>
          </div>
        </div>
        <div class="flex items-center gap-3">
          <!-- 工具箱入口 -->
          <button
            @click="router.push('/tools')"
            class="w-9 h-9 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center text-white hover:shadow-lg transition-all"
            title="工具箱"
          >
            <Wrench class="h-5 w-5" />
          </button>

          <!-- 管理员入口 (仅管理员可见) -->
          <button
            v-if="isAdmin"
            @click="router.push('/admin')"
            class="w-9 h-9 rounded-full bg-gradient-to-br from-green-500 to-cyan-600 flex items-center justify-center text-white hover:shadow-lg transition-all"
            title="管理后台"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
          </button>

          <!-- 用户头像下拉菜单 -->
          <div class="relative" @mouseenter="showDropdown = true" @mouseleave="showDropdown = false">
            <button
              class="w-9 h-9 rounded-full bg-gradient-to-br from-green-400 to-cyan-500 flex items-center justify-center text-white font-bold text-sm shadow-sm hover:shadow-lg transition-all overflow-hidden"
              :title="`${userName} - ${roleText}`"
            >
              <img
                v-if="userAvatar"
                :src="userAvatar"
                :alt="userName"
                class="w-full h-full object-cover"
                @error="handleAvatarError"
              />
              <span v-else>{{ userInitial }}</span>
            </button>

            <!-- 下拉菜单 -->
            <Transition
              enter-active-class="transition ease-out duration-100"
              enter-from-class="transform opacity-0 scale-95"
              enter-to-class="transform opacity-100 scale-100"
              leave-active-class="transition ease-in duration-75"
              leave-from-class="transform opacity-100 scale-100"
              leave-to-class="transform opacity-0 scale-95"
            >
              <div
                v-show="showDropdown"
                class="absolute right-0 mt-2 w-56 bg-white rounded-lg shadow-lg border-2 border-slate-200 py-2 z-50"
              >
                <!-- 用户信息 -->
                <div class="px-4 py-3 border-b-2 border-slate-100">
                  <p class="text-sm font-semibold text-slate-900">{{ userName }}</p>
                  <p class="text-xs text-slate-600 truncate">{{ user?.email }}</p>
                  <span :class="['inline-block mt-2 px-2 py-1 rounded-full text-xs font-semibold', roleBadgeClass]">
                    {{ roleText }}
                  </span>
                </div>

                <!-- 菜单项 -->
                <div class="py-2">
                  <button
                    @click="handleProfileClick"
                    class="w-full flex items-center gap-3 px-4 py-2 text-sm text-slate-700 hover:bg-slate-50 transition-colors"
                  >
                    <User class="w-4 h-4" />
                    <span>个人中心</span>
                  </button>
                  <button
                    @click="handleHomeClick"
                    class="w-full flex items-center gap-3 px-4 py-2 text-sm text-slate-700 hover:bg-slate-50 transition-colors"
                  >
                    <Home class="w-4 h-4" />
                    <span>返回首页</span>
                  </button>
                </div>

                <!-- 退出登录 -->
                <div class="border-t-2 border-slate-100 pt-2">
                  <button
                    @click="handleLogout"
                    class="w-full flex items-center gap-3 px-4 py-2 text-sm text-rose-600 hover:bg-rose-50 transition-colors"
                  >
                    <LogOut class="w-4 h-4" />
                    <span>退出登录</span>
                  </button>
                </div>
              </div>
            </Transition>
          </div>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { User, Home, LogOut, Wrench } from 'lucide-vue-next'
import SecureStorage, { STORAGE_KEYS } from '@/utils/storage'
import { message } from '@/utils/message'

const router = useRouter()

// 角色常量
const USER_ROLE = {
  SUPER_ADMIN: 1,
  ADMIN: 2,
  USER: 3
}

// 头像加载失败标志
const avatarLoadError = ref(false)

// 从 SecureStorage 获取用户信息
const user = computed(() => SecureStorage.getItem<{
  username: string
  id: string
  email: string
  role: number
  avatar?: string
}>(STORAGE_KEYS.AUTH_USER))

// 用户名
const userName = computed(() => user.value?.username || '未知用户')

// 用户头像 URL（如果加载失败则返回空）
const userAvatar = computed(() => {
  if (avatarLoadError.value) return ''
  return user.value?.avatar || ''
})

// 用户名首字母（大写）
const userInitial = computed(() => {
  const name = userName.value
  return name ? name.charAt(0).toUpperCase() : '?'
})

// 是否是管理员（包括超级管理员和管理员）
const isAdmin = computed(() => {
  const role = user.value?.role
  return role === USER_ROLE.SUPER_ADMIN || role === USER_ROLE.ADMIN
})

// 角色文本
const roleText = computed(() => {
  const role = user.value?.role
  if (role === USER_ROLE.SUPER_ADMIN) return '超级管理员'
  if (role === USER_ROLE.ADMIN) return '管理员'
  return '普通用户'
})

// 角色徽章样式
const roleBadgeClass = computed(() => {
  const role = user.value?.role
  if (role === USER_ROLE.SUPER_ADMIN) return 'bg-purple-100 text-purple-700'
  if (role === USER_ROLE.ADMIN) return 'bg-blue-100 text-blue-700'
  return 'bg-emerald-100 text-emerald-700'
})

// 副标题
const subtitle = computed(() => {
  return `欢迎, ${userName.value}`
})

// 下拉菜单显示状态
const showDropdown = ref(false)

// Logo 点击
const handleLogoClick = () => {
  router.push('/tasks')
}

// 个人中心点击
const handleProfileClick = () => {
  showDropdown.value = false
  router.push('/profile')
}

// 返回首页点击
const handleHomeClick = () => {
  showDropdown.value = false
  router.push('/')
}

// 退出登录
const handleLogout = () => {
  showDropdown.value = false
  SecureStorage.removeItem(STORAGE_KEYS.AUTH_TOKEN)
  SecureStorage.removeItem(STORAGE_KEYS.AUTH_USER)
  message.success('已退出登录')
  router.push('/auth')
}

// 头像加载失败处理
const handleAvatarError = () => {
  avatarLoadError.value = true
}
</script>
