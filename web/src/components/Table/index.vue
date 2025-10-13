<template>
  <div class="table-container">
    <div v-if="loading" class="flex justify-center items-center py-20">
      <div class="text-text-tertiary">加载中...</div>
    </div>
    <div v-else class="overflow-x-auto border-2 border-border-primary rounded-lg">
      <table class="w-full table-fixed">
        <thead class="bg-bg-hover border-b-2 border-border-primary">
          <tr>
            <slot name="header"></slot>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-200">
          <slot name="body"></slot>
        </tbody>
      </table>
      <div v-if="(!data || data.length === 0) && !loading" class="text-center py-10 text-text-tertiary">
        暂无数据
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  data?: any[]
  loading?: boolean
}

withDefaults(defineProps<Props>(), {
  loading: false
})
</script>

<style scoped>
th {
  @apply px-4 py-3 text-left text-xs font-semibold text-text-secondary uppercase tracking-wider;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

td {
  @apply px-4 py-3 text-sm text-text-primary;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
