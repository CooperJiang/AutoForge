<template>
  <header class="bg-bg-elevated border-b-2 border-border-primary flex-shrink-0">
    <div class="flex items-center h-16 px-6">
      <!-- Logo -->
      <router-link to="/" class="flex items-center gap-2.5 mr-8 hover:opacity-80 transition-opacity">
        <div class="w-8 h-8 flex items-center justify-center">
          <img src="/logo.png" alt="Logo" class="h-8 w-8 object-contain" />
        </div>
        <div>
          <h1 class="text-sm font-bold text-text-primary">AutoForage</h1>
          <p class="text-xs text-text-tertiary">管理后台</p>
        </div>
      </router-link>

      <!-- Navigation Tabs -->
      <nav class="flex-1 flex items-center space-x-1 overflow-x-auto">
        <button
          v-for="tab in tabs"
          :key="tab.value"
          @click="$emit('change-tab', tab.value)"
          :class="[
            'px-4 py-2 font-medium rounded-lg transition-all whitespace-nowrap',
            activeTab === tab.value
              ? 'bg-green-500/10 text-green-600 dark:text-green-400'
              : 'text-text-secondary hover:text-text-primary hover:bg-bg-hover',
          ]"
        >
          {{ tab.label }}
        </button>
      </nav>

      <!-- Right Actions -->
      <div class="flex items-center space-x-3 ml-4">
        <ThemeToggle />
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import ThemeToggle from '@/components/ThemeToggle/index.vue'

interface Tab {
  label: string
  value: string
}

defineProps<{
  activeTab: string
  tabs: Tab[]
}>()

defineEmits<{
  'change-tab': [value: string]
}>()
</script>
