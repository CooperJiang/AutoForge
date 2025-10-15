<template>
  <nav class="flex items-center gap-2">
    <template v-for="(item, index) in items" :key="index">
      <ChevronRight v-if="index > 0" class="w-4 h-4 text-text-tertiary flex-shrink-0" />
      <a
        :href="item.to || 'javascript:;'"
        @click.prevent="item.to && handleClick(item.to)"
        :class="[
          'text-sm transition-all duration-200 rounded-md px-2 py-1',
          index === items.length - 1
            ? 'text-text-primary font-semibold cursor-default bg-bg-hover'
            : 'text-text-secondary hover:text-primary hover:bg-primary-light cursor-pointer font-medium',
        ]"
      >
        {{ item.label }}
      </a>
    </template>
  </nav>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { ChevronRight } from 'lucide-vue-next'

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
