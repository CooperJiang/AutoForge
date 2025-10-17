<template>
  <div class="h-screen flex bg-bg-secondary">
    <div class="border-r-2 border-border-primary">
      <Sidebar />
    </div>

    <div class="flex-1 flex flex-col overflow-hidden">
      <div class="border-b-2 border-border-primary">
        <AppHeader />
      </div>

      <div class="flex-1 overflow-y-auto bg-bg-secondary p-6">
        <router-view />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import AppHeader from '@/components/AppHeader'
import Sidebar from '@/components/Sidebar'
</script>
