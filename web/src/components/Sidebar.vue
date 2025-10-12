<template>
  <aside class="w-64 bg-white border-r-2 border-slate-200 flex flex-col flex-shrink-0">
    <!-- Logo -->
    <div class="px-6 py-4 border-b-2 border-slate-200">
      <div class="flex items-center gap-3 cursor-pointer" @click="handleLogoClick">
        <img src="/logo.png" alt="Logo" class="h-10 w-10 object-contain" />
        <div>
          <h1 class="text-lg font-bold text-slate-900">定时任务系统</h1>
          <p class="text-xs text-slate-500">自动化工作流</p>
        </div>
      </div>
    </div>

    <!-- Navigation -->
    <nav class="flex-1 overflow-y-auto py-4">
      <div class="px-3 space-y-1">
        <router-link
          v-for="item in menuItems"
          :key="item.path"
          :to="item.path"
          :class="[
            'flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-all',
            isActive(item.path)
              ? 'bg-blue-50 text-blue-700 shadow-sm'
              : 'text-slate-700 hover:bg-slate-50'
          ]"
        >
          <component :is="item.icon" class="w-5 h-5 flex-shrink-0" />
          <span>{{ item.label }}</span>
        </router-link>
      </div>

      <!-- Admin Section -->
      <div v-if="isAdmin" class="mt-6 px-3">
        <div class="px-3 py-2 text-xs font-semibold text-slate-500 uppercase tracking-wider">
          管理功能
        </div>
        <div class="space-y-1">
          <router-link
            v-for="item in adminMenuItems"
            :key="item.path"
            :to="item.path"
            :class="[
              'flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-all',
              isActive(item.path)
                ? 'bg-green-50 text-green-700 shadow-sm'
                : 'text-slate-700 hover:bg-slate-50'
            ]"
          >
            <component :is="item.icon" class="w-5 h-5 flex-shrink-0" />
            <span>{{ item.label }}</span>
          </router-link>
        </div>
      </div>
    </nav>

    <!-- User Info -->
    <div class="border-t-2 border-slate-200 p-4">
      <div class="flex items-center gap-3">
        <div class="w-10 h-10 rounded-full bg-gradient-to-br from-green-400 to-cyan-500 flex items-center justify-center text-white font-bold text-sm shadow-sm flex-shrink-0">
          {{ userInitial }}
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-semibold text-slate-900 truncate">{{ userName }}</p>
          <p class="text-xs text-slate-500 truncate">{{ user?.email }}</p>
        </div>
        <button
          @click="handleLogout"
          class="flex-shrink-0 w-8 h-8 flex items-center justify-center rounded-lg text-rose-600 hover:bg-rose-50 transition-colors"
          title="退出登录"
        >
          <LogOut class="w-4 h-4" />
        </button>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Home, ListTodo, Wrench, Workflow, Settings, User, LogOut } from 'lucide-vue-next'
import SecureStorage, { STORAGE_KEYS } from '@/utils/storage'
import { message } from '@/utils/message'

const router = useRouter()
const route = useRoute()

// 角色常量
const USER_ROLE = {
  SUPER_ADMIN: 1,
  ADMIN: 2,
  USER: 3
}

// 从 SecureStorage 获取用户信息
const user = computed(() => SecureStorage.getItem<{
  username: string
  id: string
  email: string
  role: number
}>(STORAGE_KEYS.AUTH_USER))

// 用户名
const userName = computed(() => user.value?.username || '未知用户')

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

// 菜单项
const menuItems = [
  { path: '/', label: '任务管理', icon: ListTodo },
  { path: '/workflows', label: '工作流', icon: Workflow },
  { path: '/tools', label: '工具箱', icon: Wrench },
  { path: '/profile', label: '个人中心', icon: User },
]

// 管理员菜单项
const adminMenuItems = [
  { path: '/admin', label: '管理后台', icon: Settings },
]

// 判断菜单项是否激活
const isActive = (path: string) => {
  if (path === '/') {
    return route.path === '/' || route.path === '/tasks'
  }
  return route.path.startsWith(path)
}

// Logo 点击
const handleLogoClick = () => {
  router.push('/')
}

// 退出登录
const handleLogout = () => {
  SecureStorage.removeItem(STORAGE_KEYS.AUTH_TOKEN)
  SecureStorage.removeItem(STORAGE_KEYS.AUTH_USER)
  message.success('已退出登录')
  router.push('/auth')
}
</script>
