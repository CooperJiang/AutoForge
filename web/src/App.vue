<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script setup lang="ts">
import { useTheme } from '@/composables/useTheme'

// 初始化主题（确保应用启动时立即应用主题）
useTheme()
</script>

<style>
#app {
  width: 100%;
  min-height: 100vh;
}
</style>
