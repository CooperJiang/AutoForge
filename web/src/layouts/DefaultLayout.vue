<template>
  <div class="h-screen flex bg-bg-secondary">
    <!-- Sidebar with border -->
    <div class="border-r-2 border-border-primary">
      <Sidebar />
    </div>

    <!-- Main Area -->
    <div class="flex-1 flex flex-col overflow-hidden">
      <!-- Fixed Header with border -->
      <div class="border-b-2 border-border-primary">
        <AppHeader />
      </div>

      <!-- Scrollable Content Area -->
      <div class="flex-1 overflow-y-auto bg-bg-secondary">
        <div class="max-w-[1920px] mx-auto p-6">
          <router-view />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import AppHeader from '@/components/AppHeader'
import Sidebar from '@/components/Sidebar'
</script>
