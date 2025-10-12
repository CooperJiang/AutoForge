<template>
  <div class="table-container">
    <div v-if="loading" class="flex justify-center items-center py-20">
      <div class="text-slate-500">加载中...</div>
    </div>
    <div v-else class="overflow-x-auto border-2 border-slate-200 rounded-lg">
      <table class="w-full table-fixed">
        <thead class="bg-slate-50 border-b-2 border-slate-200">
          <tr>
            <slot name="header"></slot>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-200">
          <slot name="body"></slot>
        </tbody>
      </table>
      <div v-if="(!data || data.length === 0) && !loading" class="text-center py-10 text-slate-400">
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
  @apply px-4 py-3 text-left text-xs font-semibold text-slate-700 uppercase tracking-wider;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

td {
  @apply px-4 py-3 text-sm text-slate-900;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
