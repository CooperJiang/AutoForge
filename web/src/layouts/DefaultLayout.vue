<template>
  <div class="h-screen flex bg-slate-50">
    <!-- Sidebar -->
    <Sidebar />

    <!-- Main Area -->
    <div class="flex-1 flex flex-col overflow-hidden">
      <!-- Fixed Header -->
      <AppHeader />

      <!-- Scrollable Content Area -->
      <div class="flex-1 overflow-y-auto bg-slate-50">
        <router-view />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import AppHeader from '@/components/AppHeader'
import Sidebar from '@/components/Sidebar.vue'
</script>
