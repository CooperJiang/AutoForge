<template>
  <button type="button" @click="handleToggle" :title="getTooltipText()" class="theme-toggle-btn">
    
    <Sun v-if="appliedTheme === 'light'" class="w-5 h-5" />

    
    <Moon v-else class="w-5 h-5" />
  </button>
</template>

<script setup lang="ts">
import { Sun, Moon } from 'lucide-vue-next'
import { useTheme } from '@/composables/useTheme'

const { currentTheme, appliedTheme, toggleTheme } = useTheme()

// 切换主题
const handleToggle = () => {
  toggleTheme()
}

// 获取提示文本
const getTooltipText = () => {
  const themeNames = {
    light: '亮色模式',
    dark: '暗色模式',
    auto: '跟随系统',
  }
  return `当前：${themeNames[currentTheme.value]} | 点击切换`
}
</script>

<style scoped>
.theme-toggle-btn {
  @apply relative inline-flex items-center justify-center;
  @apply w-9 h-9 rounded-lg;
  @apply text-text-secondary hover:text-text-primary;
  @apply hover:bg-bg-hover active:bg-bg-active;
  @apply transition-all duration-200;
  @apply focus:outline-none;
}

.theme-toggle-btn svg {
  @apply transition-transform duration-300;
}

.theme-toggle-btn:hover svg {
  @apply scale-110;
}

.theme-toggle-btn:active svg {
  @apply scale-95;
}
</style>
