<template>
  <div>
    <!-- 搜索筛选栏 -->
    <div class="flex gap-3 mb-6 items-center">
      <div class="flex-shrink-0" style="width: 200px;">
        <BaseInput
          v-model="filters.user_id"
          placeholder="用户ID"
        />
      </div>
      <div class="flex-shrink-0" style="width: 150px;">
        <BaseSelect
          v-model="filters.status"
          :options="statusOptions"
          placeholder="全部状态"
        />
      </div>
      <div class="flex-1">
        <BaseInput
          v-model="filters.keyword"
          placeholder="搜索任务名称或URL"
        />
      </div>
      <BaseButton
        @click="$emit('search', filters)"
        variant="primary"
        class="flex-shrink-0"
      >
        搜索
      </BaseButton>
    </div>

    <!-- 表格 -->
    <div class="overflow-x-auto">
      <table class="w-full">
        <thead>
          <tr class="border-b-2 border-slate-200 text-left">
            <th class="pb-3 text-sm font-semibold text-slate-700 w-24">ID</th>
            <th class="pb-3 text-sm font-semibold text-slate-700 w-28">用户ID</th>
            <th class="pb-3 text-sm font-semibold text-slate-700 w-32">任务名称</th>
            <th class="pb-3 text-sm font-semibold text-slate-700 w-28">工具</th>
            <th class="pb-3 text-sm font-semibold text-slate-700 w-32">调度规则</th>
            <th class="pb-3 text-sm font-semibold text-slate-700 w-20">状态</th>
            <th class="pb-3 text-sm font-semibold text-slate-700 w-28">下次执行</th>
            <th class="pb-3 text-sm font-semibold text-slate-700 w-40">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="task in tasks"
            :key="task.id"
            class="border-b border-slate-100 hover:bg-slate-50 transition-colors"
          >
            <td class="py-3 text-sm text-slate-900 font-mono group relative w-24 align-middle">
              <div class="flex items-center gap-2">
                <span class="truncate" :title="task.id">{{ truncateId(task.id) }}</span>
                <button
                  @click="$emit('copy-id', task.id)"
                  class="opacity-0 group-hover:opacity-100 p-1 hover:bg-slate-200 rounded transition-all"
                  title="复制完整ID"
                >
                  <Copy :size="14" />
                </button>
              </div>
            </td>
            <td class="py-3 text-sm text-slate-900 font-mono w-28 truncate align-middle" :title="task.user_id">
              {{ maskUserId(task.user_id) }}
            </td>
            <td class="py-3 text-sm text-slate-900 font-medium w-32 truncate align-middle" :title="task.name">
              {{ task.name }}
            </td>
            <td class="py-3 text-sm text-slate-600 w-28 align-middle">
              <span class="inline-flex px-2 py-0.5 text-xs font-medium rounded bg-slate-100 text-slate-700">
                {{ task.tool_code || 'N/A' }}
              </span>
            </td>
            <td class="py-3 text-sm text-slate-600 w-32 truncate align-middle" :title="getScheduleText(task)">
              {{ getScheduleText(task) }}
            </td>
            <td class="py-3 w-20 align-middle">
              <span
                :class="[
                  'inline-block px-2 py-1 text-xs font-medium rounded-full whitespace-nowrap',
                  task.enabled ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-700'
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
                  class="p-1.5 bg-purple-100 hover:bg-purple-200 text-purple-700 rounded-lg transition-colors"
                  title="查看详情"
                >
                  <Eye :size="16" />
                </button>
                <button
                  @click="$emit('toggle-status', task)"
                  :class="[
                    'p-1.5 rounded-lg transition-colors',
                    task.enabled
                      ? 'bg-gray-100 hover:bg-gray-200 text-gray-700'
                      : 'bg-green-100 hover:bg-green-200 text-green-700'
                  ]"
                  :title="task.enabled ? '禁用' : '启用'"
                >
                  <Pause v-if="task.enabled" :size="16" />
                  <Play v-else :size="16" />
                </button>
                <button
                  @click="$emit('execute', task)"
                  class="p-1.5 bg-blue-100 hover:bg-blue-200 text-blue-700 rounded-lg transition-colors"
                  title="立即执行"
                >
                  <Zap :size="16" />
                </button>
                <button
                  @click="$emit('delete', task)"
                  class="p-1.5 bg-red-100 hover:bg-red-200 text-red-700 rounded-lg transition-colors"
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

    <!-- 分页 -->
    <div class="flex justify-between items-center mt-6">
      <div class="text-sm text-slate-600">
        共 {{ total }} 条记录
      </div>
      <div class="flex gap-2">
        <button
          @click="$emit('prev-page')"
          :disabled="currentPage === 1"
          class="px-4 py-2 bg-slate-100 text-slate-700 text-sm font-medium rounded-lg hover:bg-slate-200 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          上一页
        </button>
        <span class="px-4 py-2 text-sm text-slate-700">
          {{ currentPage }} / {{ totalPages }}
        </span>
        <button
          @click="$emit('next-page')"
          :disabled="currentPage >= totalPages"
          class="px-4 py-2 bg-slate-100 text-slate-700 text-sm font-medium rounded-lg hover:bg-slate-200 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
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
  keyword: ''
})

const statusOptions = [
  { label: '全部状态', value: '' },
  { label: '已启用', value: 'enabled' },
  { label: '已禁用', value: 'disabled' }
]

const getScheduleText = (task: Task) => {
  const typeMap: Record<string, string> = {
    daily: '每天',
    weekly: '每周',
    monthly: '每月',
    hourly: '每小时',
    interval: '间隔',
    cron: 'Cron'
  }
  const type = typeMap[task.schedule_type] || task.schedule_type
  return `${type} ${task.schedule_value}`
}
</script>
