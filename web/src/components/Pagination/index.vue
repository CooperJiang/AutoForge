<template>
  <div class="flex items-center justify-between px-4 py-3 border-t-2 border-border-primary">
    <div class="text-sm text-text-secondary">
      共 <span class="font-semibold">{{ total }}</span> 条记录
    </div>
    <div class="flex items-center gap-2">
      <BaseButton
        variant="secondary"
        size="sm"
        :disabled="currentPage === 1"
        @click="handlePageChange(currentPage - 1)"
      >
        上一页
      </BaseButton>

      <div class="flex gap-1">
        <button
          v-for="page in visiblePages"
          :key="page"
          @click="page !== '...' && handlePageChange(page as number)"
          :disabled="page === '...'"
          class="px-3 py-1 text-sm border-2 rounded-md transition-colors"
          :class="page === currentPage
            ? 'bg-success-light0 text-white border-green-500'
            : page === '...'
            ? 'border-transparent text-text-tertiary cursor-default'
            : 'border-border-primary text-text-secondary hover:border-green-500 hover:text-green-600'"
        >
          {{ page }}
        </button>
      </div>

      <BaseButton
        variant="secondary"
        size="sm"
        :disabled="currentPage === totalPages"
        @click="handlePageChange(currentPage + 1)"
      >
        下一页
      </BaseButton>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import BaseButton from '../BaseButton/index.vue'

interface Props {
  currentPage: number
  pageSize: number
  total: number
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:currentPage': [page: number]
  'change': [page: number]
}>()

const totalPages = computed(() => Math.ceil(props.total / props.pageSize))

const visiblePages = computed(() => {
  const pages: (number | string)[] = []
  const current = props.currentPage
  const total = totalPages.value

  if (total <= 7) {
    for (let i = 1; i <= total; i++) {
      pages.push(i)
    }
  } else {
    pages.push(1)

    if (current > 3) {
      pages.push('...')
    }

    const start = Math.max(2, current - 1)
    const end = Math.min(total - 1, current + 1)

    for (let i = start; i <= end; i++) {
      pages.push(i)
    }

    if (current < total - 2) {
      pages.push('...')
    }

    pages.push(total)
  }

  return pages
})

const handlePageChange = (page: number) => {
  if (page === props.currentPage || page < 1 || page > totalPages.value) return
  emit('update:currentPage', page)
  emit('change', page)
}
</script>
