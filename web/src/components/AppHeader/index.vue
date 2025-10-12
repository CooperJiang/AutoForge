<template>
  <header class="bg-white border-b-2 border-slate-200 flex-shrink-0">
    <div class="px-6 py-3">
      <!-- 面包屑导航 -->
      <Breadcrumb v-if="breadcrumbItems.length > 0" :items="breadcrumbItems" />
    </div>
  </header>
</template>

<script setup lang="ts">
import Breadcrumb from '@/components/Breadcrumb.vue'
import { useBreadcrumb } from '@/composables/useBreadcrumb'

const { breadcrumbItems } = useBreadcrumb()
</script>
