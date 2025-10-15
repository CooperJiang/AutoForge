<template>
  <aside class="w-56 bg-bg-primary flex flex-col flex-shrink-0 h-screen">
    
    <div class="px-4 py-4 border-b border-border-primary">
      <div class="flex items-center gap-2.5 cursor-pointer group" @click="handleLogoClick">
        <div
          class="w-8 h-8 flex items-center justify-center group-hover:scale-105 transition-transform"
        >
          <img src="/logo.png" alt="Logo" class="h-8 w-8 object-contain" />
        </div>
        <div>
          <h1
            class="text-sm font-bold text-text-primary group-hover:text-primary transition-colors"
          >
            AutoForage
          </h1>
          <p class="text-xs text-text-tertiary">AutoForage Workflow</p>
        </div>
      </div>
    </div>

    
    <nav class="flex-1 overflow-y-auto py-2.5 px-2.5">
      <div class="space-y-0.5">
        <router-link
          v-for="item in menuItems"
          :key="item.path"
          :to="item.path"
          :class="[
            'flex items-center gap-2.5 px-3 py-2 rounded-lg text-sm font-medium transition-all duration-200',
            isActive(item.path)
              ? 'bg-[var(--color-primary)] text-white shadow-lg shadow-primary/30'
              : 'text-text-secondary hover:bg-bg-hover hover:text-text-primary',
          ]"
        >
          <component :is="item.icon" class="w-4 h-4 flex-shrink-0" />
          <span>{{ item.label }}</span>
        </router-link>
      </div>

      
      <div v-if="isAdmin" class="mt-5">
        <div class="px-3 py-1.5 mb-1 text-xs font-bold text-text-tertiary uppercase tracking-wider">
          管理功能
        </div>
        <div class="space-y-0.5">
          <router-link
            v-for="item in adminMenuItems"
            :key="item.path"
            :to="item.path"
            :class="[
              'flex items-center gap-2.5 px-3 py-2 rounded-lg text-sm font-medium transition-all duration-200',
              isActive(item.path)
                ? 'bg-warning text-white shadow-lg shadow-warning/30'
                : 'text-text-secondary hover:bg-bg-hover hover:text-text-primary',
            ]"
          >
            <component :is="item.icon" class="w-4 h-4 flex-shrink-0" />
            <span>{{ item.label }}</span>
          </router-link>
        </div>
      </div>
    </nav>

    
    <div class="border-t border-border-primary p-3">
      <div class="flex items-center gap-2.5">
        <div
          class="w-9 h-9 rounded-lg bg-gradient-primary flex items-center justify-center text-white font-bold text-sm shadow-md flex-shrink-0"
        >
          {{ userInitial }}
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-semibold text-text-primary truncate">{{ userName }}</p>
          <p class="text-xs text-text-tertiary truncate">{{ user?.email }}</p>
        </div>
        <button
          @click="handleLogout"
          class="flex-shrink-0 w-8 h-8 flex items-center justify-center rounded-lg text-text-tertiary hover:text-error hover:bg-error-light transition-all"
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
import { ListTodo, Wrench, Workflow, Settings, User, LogOut } from 'lucide-vue-next'
import SecureStorage, { STORAGE_KEYS } from '@/utils/storage'
import { message } from '@/utils/message'

const router = useRouter()
const route = useRoute()

// 角色常量
const USER_ROLE = {
  SUPER_ADMIN: 1,
  ADMIN: 2,
  USER: 3,
}

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
const adminMenuItems = [{ path: '/admin', label: '管理后台', icon: Settings }]

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
