<template>
  <div class="h-screen flex flex-col bg-bg-secondary">
    
    <AppHeader />

    
    <div class="flex-1 overflow-y-auto bg-bg-secondary">
      <router-view />
    </div>
  </div>
</template>

<script setup lang="ts">
import AppHeader from '@/components/AppHeader'
</script>
