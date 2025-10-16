<template>
  <div class="h-screen flex flex-col bg-bg-secondary">
    <AdminHeader :active-tab="activeTab" :tabs="tabs" @change-tab="handleTabChange" />

    <div class="flex-1 overflow-y-auto">
      <router-view v-slot="{ Component }">
        <component :is="Component" :active-tab="activeTab" @update:active-tab="activeTab = $event" />
      </router-view>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import AdminHeader from '@/components/AdminHeader/index.vue'

const route = useRoute()
const router = useRouter()

// 从 URL query 参数中恢复 tab 状态，默认为 executions
const activeTab = ref((route.query.tab as string) || 'executions')

const tabs = [
  { label: '执行记录', value: 'executions' },
  { label: '任务管理', value: 'tasks' },
  { label: '用户管理', value: 'users' },
  { label: '模板管理', value: 'workflows' },
  { label: '模板分类', value: 'categories' },
]

const handleTabChange = (value: string) => {
  activeTab.value = value
  // 更新 URL query 参数，保持 tab 状态
  router.push({ query: { tab: value } })
}

onMounted(() => {
  // 确保初始状态也更新 URL
  if (!route.query.tab) {
    router.replace({ query: { tab: activeTab.value } })
  }
})
</script>
