<template>
  <nav class="flex items-center gap-2 text-sm">
    <template v-for="(item, index) in items" :key="index">
      <span v-if="index > 0" class="text-slate-400 mx-1">/</span>
      <a
        :href="item.to || 'javascript:;'"
        @click.prevent="item.to && handleClick(item.to)"
        :class="[
          'transition-colors',
          index === items.length - 1
            ? 'text-slate-900 font-semibold cursor-default'
            : 'text-slate-600 hover:text-slate-900 cursor-pointer'
        ]"
      >
        {{ item.label }}
      </a>
    </template>
  </nav>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

interface BreadcrumbItem {
  label: string
  to?: string
}

interface Props {
  items: BreadcrumbItem[]
}

defineProps<Props>()

const router = useRouter()

const handleClick = (to: string) => {
  router.push(to)
}
</script>
