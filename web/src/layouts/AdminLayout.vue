<template>
  <div class="h-screen flex flex-col bg-slate-50">
    <!-- Fixed Header -->
    <AppHeader />

    <!-- Scrollable Content Area -->
    <div class="flex-1 overflow-y-auto bg-slate-50">
      <router-view />
    </div>
  </div>
</template>

<script setup lang="ts">
import AppHeader from '@/components/AppHeader'
</script>
