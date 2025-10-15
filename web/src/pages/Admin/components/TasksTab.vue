<template>
  <div>
    
    <div class="flex gap-3 mb-6 items-center">
      <div class="flex-shrink-0" style="width: 200px">
        <BaseInput v-model="filters.user_id" placeholder="用户ID" />
      </div>
      <div class="flex-shrink-0" style="width: 150px">
        <BaseSelect v-model="filters.status" :options="statusOptions" placeholder="全部状态" />
      </div>
      <div class="flex-1">
        <BaseInput v-model="filters.keyword" placeholder="搜索任务名称或URL" />
      </div>
      <BaseButton @click="$emit('search', filters)" variant="primary" class="flex-shrink-0">
        搜索
      </BaseButton>
    </div>

    
    <div class="overflow-x-auto">
      <table class="w-full">
        <thead>
          <tr class="border-b-2 border-border-primary text-left">
            <th class="pb-3 text-sm font-semibold text-text-secondary w-24">ID</th>
            <th class="pb-3 text-sm font-semibold text-text-secondary w-28">用户ID</th>
            <th class="pb-3 text-sm font-semibold text-text-secondary w-32">任务名称</th>
            <th class="pb-3 text-sm font-semibold text-text-secondary w-28">工具</th>
            <th class="pb-3 text-sm font-semibold text-text-secondary w-32">调度规则</th>
            <th class="pb-3 text-sm font-semibold text-text-secondary w-20">状态</th>
            <th class="pb-3 text-sm font-semibold text-text-secondary w-28">下次执行</th>
            <th class="pb-3 text-sm font-semibold text-text-secondary w-40">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="task in tasks"
            :key="task.id"
            class="border-b border-border-primary hover:bg-bg-hover transition-colors"
          >
            <td class="py-3 text-sm text-text-primary font-mono group relative w-24 align-middle">
              <div class="flex items-center gap-2">
                <span class="truncate" :title="task.id">{{ truncateId(task.id) }}</span>
                <button
                  @click="$emit('copy-id', task.id)"
                  class="opacity-0 group-hover:opacity-100 p-1 hover:bg-bg-tertiary rounded transition-all"
                  title="复制完整ID"
                >
                  <Copy :size="14" />
                </button>
              </div>
            </td>
            <td
              class="py-3 text-sm text-text-primary font-mono w-28 truncate align-middle"
              :title="task.user_id"
            >
              {{ maskUserId(task.user_id) }}
            </td>
            <td
              class="py-3 text-sm text-text-primary font-medium w-32 truncate align-middle"
              :title="task.name"
            >
              {{ task.name }}
            </td>
            <td class="py-3 text-sm text-text-secondary w-28 align-middle">
              <span
                class="inline-flex px-2 py-0.5 text-xs font-medium rounded bg-bg-tertiary text-text-secondary"
              >
                {{ task.tool_code || 'N/A' }}
              </span>
            </td>
            <td
              class="py-3 text-sm text-text-secondary w-32 truncate align-middle"
              :title="getScheduleText(task)"
            >
              {{ getScheduleText(task) }}
            </td>
            <td class="py-3 w-20 align-middle">
              <span
                :class="[
                  'inline-block px-2 py-1 text-xs font-medium rounded-full whitespace-nowrap border',
                  task.enabled
                    ? 'bg-green-500/10 text-green-600 dark:text-green-400 border-green-500/20'
                    : 'bg-slate-500/10 text-slate-600 dark:text-slate-400 border-slate-500/20',
                ]"
              >
                {{ task.enabled ? '已启用' : '已禁用' }}
              </span>
            </td>
            <td class="py-3 text-sm w-28 align-middle">
              <NextRunCountdown :next-run-time="task.next_run_time" />
            </td>
            <td class="py-3 w-40 align-middle">
              <div class="flex gap-2">
                <button
                  @click="$emit('view', task)"
                  class="p-1.5 bg-purple-500/10 hover:bg-purple-500/20 text-purple-600 dark:text-purple-400 rounded-lg transition-colors border border-purple-500/20"
                  title="查看详情"
                >
                  <Eye :size="16" />
                </button>
                <button
                  @click="$emit('toggle-status', task)"
                  :class="[
                    'p-1.5 rounded-lg transition-colors border',
                    task.enabled
                      ? 'bg-slate-500/10 hover:bg-slate-500/20 text-slate-600 dark:text-slate-400 border-slate-500/20'
                      : 'bg-green-500/10 hover:bg-green-500/20 text-green-600 dark:text-green-400 border-green-500/20',
                  ]"
                  :title="task.enabled ? '禁用' : '启用'"
                >
                  <Pause v-if="task.enabled" :size="16" />
                  <Play v-else :size="16" />
                </button>
                <button
                  @click="$emit('execute', task)"
                  class="p-1.5 bg-blue-500/10 hover:bg-blue-500/20 text-blue-600 dark:text-blue-400 rounded-lg transition-colors border border-blue-500/20"
                  title="立即执行"
                >
                  <Zap :size="16" />
                </button>
                <button
                  @click="$emit('delete', task)"
                  class="p-1.5 bg-red-500/10 hover:bg-red-500/20 text-red-600 dark:text-red-400 rounded-lg transition-colors border border-red-500/20"
                  title="删除"
                >
                  <Trash2 :size="16" />
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    
    <div class="flex justify-between items-center mt-6">
      <div class="text-sm text-text-secondary">共 {{ total }} 条记录</div>
      <div class="flex gap-2">
        <button
          @click="$emit('prev-page')"
          :disabled="currentPage === 1"
          class="px-4 py-2 bg-bg-tertiary text-text-secondary text-sm font-medium rounded-lg hover:bg-bg-tertiary disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          上一页
        </button>
        <span class="px-4 py-2 text-sm text-text-secondary">
          {{ currentPage }} / {{ totalPages }}
        </span>
        <button
          @click="$emit('next-page')"
          :disabled="currentPage >= totalPages"
          class="px-4 py-2 bg-bg-tertiary text-text-secondary text-sm font-medium rounded-lg hover:bg-bg-tertiary disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          下一页
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Play, Pause, Zap, Trash2, Eye, Copy } from 'lucide-vue-next'
import type { Task } from '@/api/task'
import BaseSelect from '@/components/BaseSelect'
import BaseInput from '@/components/BaseInput'
import BaseButton from '@/components/BaseButton'
import NextRunCountdown from '@/components/NextRunCountdown'
import { maskUserId, truncateId } from '@/utils/format'

defineProps<{
  tasks: Task[]
  total: number
  currentPage: number
  totalPages: number
}>()

defineEmits<{
  search: [filters: { user_id: string; status: string; keyword: string }]
  'copy-id': [id: string]
  view: [task: Task]
  'toggle-status': [task: Task]
  execute: [task: Task]
  delete: [task: Task]
  'prev-page': []
  'next-page': []
}>()

const filters = ref({
  user_id: '',
  status: '',
  keyword: '',
})

const statusOptions = [
  { label: '全部状态', value: '' },
  { label: '已启用', value: 'enabled' },
  { label: '已禁用', value: 'disabled' },
]

const getScheduleText = (task: Task) => {
  const typeMap: Record<string, string> = {
    daily: '每天',
    weekly: '每周',
    monthly: '每月',
    hourly: '每小时',
    interval: '间隔',
    cron: 'Cron',
  }
  const type = typeMap[task.schedule_type] || task.schedule_type
  return `${type} ${task.schedule_value}`
}
</script>
