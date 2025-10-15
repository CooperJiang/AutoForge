<template>
  <div
    class="bg-bg-elevated rounded-2xl shadow-lg border border-border-primary overflow-hidden sticky top-6"
  >
    <!-- 用户信息 -->
    <div class="bg-gradient-to-br from-primary to-accent p-6 text-white">
      <div class="flex flex-col items-center">
        <div
          class="w-20 h-20 rounded-full bg-bg-elevated/20 backdrop-blur-sm flex items-center justify-center text-white text-3xl font-bold shadow-lg mb-4 ring-4 ring-white/30"
        >
          {{ userInitial }}
        </div>
        <h2 class="text-lg font-bold mb-1">{{ userName }}</h2>
        <p class="text-sm text-primary-text mb-3 truncate max-w-full px-4">{{ userEmail }}</p>
        <span
          class="inline-flex items-center px-3 py-1 rounded-full text-xs font-semibold bg-bg-elevated/20 backdrop-blur-sm"
        >
          {{ roleText }}
        </span>
      </div>
    </div>

    <!-- 菜单列表 -->
    <nav class="p-4 space-y-1">
      <button
        v-for="item in menuItems"
        :key="item.id"
        @click="$emit('tab-change', item.id)"
        :class="[
          'w-full flex items-center gap-3 px-4 py-2.5 rounded-lg text-sm font-medium transition-all',
          activeTab === item.id
            ? 'bg-primary-light text-primary shadow-sm'
            : 'text-text-secondary hover:bg-bg-hover',
        ]"
      >
        <component :is="item.icon" class="w-5 h-5"></component>
        <span>{{ item.label }}</span>
      </button>
    </nav>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { User, Lock } from 'lucide-vue-next'

interface Props {
  activeTab: string
  userName: string
  userEmail: string
  roleText: string
}

const props = defineProps<Props>()

defineEmits<{
  'tab-change': [tabId: string]
}>()

// 用户名首字母（大写）
const userInitial = computed(() => {
  const name = props.userName
  return name ? name.charAt(0).toUpperCase() : '?'
})

// 菜单项
const menuItems = [
  { id: 'profile', label: '个人资料', icon: User },
  { id: 'password', label: '修改密码', icon: Lock },
]
</script>
