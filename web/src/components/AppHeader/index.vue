<template>
  <header class="bg-bg-elevated flex-shrink-0 backdrop-blur-xl">
    <div class="max-w-[1920px] mx-auto px-6 py-3.5 flex items-center justify-between">
      
      <Breadcrumb v-if="breadcrumbItems.length > 0" :items="breadcrumbItems" />
      <div v-else class="flex-1" />

      
      <div class="flex items-center gap-3">
        
        <ThemeToggle />
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import Breadcrumb from '@/components/Breadcrumb'
import ThemeToggle from '@/components/ThemeToggle/index.vue'
import { useBreadcrumb } from '@/composables/useBreadcrumb'

const { breadcrumbItems } = useBreadcrumb()
</script>
