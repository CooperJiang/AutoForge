<template>
  <div
    class="group bg-bg-elevated rounded-lg border border-border-primary hover:border-primary/40 transition-all duration-200 overflow-hidden"
  >
    <!-- 头部 -->
    <div class="p-3">
      <div class="flex items-start justify-between gap-2 mb-1.5">
        <h3 class="text-sm font-semibold text-text-primary truncate flex-1">
          {{ workflow.name }}
        </h3>
        <div
          :class="[
            'flex-shrink-0 w-1.5 h-1.5 rounded-full mt-1',
            workflow.enabled ? 'bg-green-500 shadow-sm shadow-green-500/50' : 'bg-slate-400',
          ]"
          :title="workflow.enabled ? '已启用' : '已禁用'"
        ></div>
      </div>
      <p class="text-xs text-text-tertiary line-clamp-1">
        {{ workflow.description || '暂无描述' }}
      </p>
    </div>

    <!-- 统计信息和时间 -->
    <div class="px-3 pb-2 space-y-2">
      <!-- 统计 -->
      <div class="flex items-center gap-3 text-xs text-text-tertiary">
        <span class="flex items-center gap-1">
          <Box class="w-3 h-3" />
          {{ workflow.nodes.length }}
        </span>
        <span class="flex items-center gap-1">
          <GitBranch class="w-3 h-3" />
          {{ workflow.edges.length }}
        </span>
      </div>

      <!-- 下次执行时间 -->
      <div
        v-if="hasSchedule"
        class="flex items-center gap-1.5 text-xs bg-bg-hover px-2 py-1.5 rounded"
      >
        <Clock class="w-3 h-3 text-blue-500" />
        <span class="truncate text-text-secondary">
          <CountdownDisplay
            :timestamp="workflow.next_run_time || 0"
            @finish="handleCountdownFinish"
          />
        </span>
      </div>
      <div v-else class="text-xs text-text-placeholder px-2 py-1.5 flex items-center gap-1.5">
        <Clock class="w-3 h-3" />
        {{ getNoScheduleText(workflow) }}
      </div>
    </div>

    <!-- 操作按钮 -->
    <div class="px-3 pb-3 flex items-center gap-1.5">
      <button
        @click="$emit('edit', workflow)"
        class="flex-1 px-2 py-1.5 text-xs font-medium text-text-secondary hover:text-primary hover:bg-bg-hover rounded transition-colors flex items-center justify-center gap-1"
        title="编辑"
      >
        <Edit3 class="w-3 h-3" />
        编辑
      </button>
      <button
        @click="$emit('executions', workflow)"
        class="flex-1 px-2 py-1.5 text-xs font-medium text-text-secondary hover:text-primary hover:bg-bg-hover rounded transition-colors flex items-center justify-center gap-1"
        title="历史"
      >
        <History class="w-3 h-3" />
        历史
      </button>
      <button
        @click="$emit('execute', workflow)"
        :disabled="workflow.nodes.length === 0"
        class="flex-1 px-2 py-1.5 text-xs font-medium text-text-secondary hover:text-primary hover:bg-primary/10 rounded transition-colors disabled:opacity-40 disabled:cursor-not-allowed flex items-center justify-center gap-1 border border-border-primary hover:border-primary/30"
        title="执行"
      >
        <Play class="w-3 h-3" />
        执行
      </button>
    </div>

    <!-- 底部操作栏 -->
    <div class="px-3 py-2 bg-bg-hover/50 border-t border-border-primary flex items-center gap-1.5">
      <button
        @click="$emit('toggle', workflow)"
        :class="[
          'flex-1 px-2 py-1 text-xs font-medium rounded transition-all flex items-center justify-center gap-1',
          workflow.enabled
            ? 'text-green-600 dark:text-green-400 hover:bg-green-500/10'
            : 'text-slate-500 hover:bg-slate-500/10',
        ]"
      >
        <Power class="w-3 h-3" />
        {{ workflow.enabled ? '启用中' : '已禁用' }}
      </button>
      <button
        @click="handleDelete"
        class="p-1 text-text-placeholder hover:text-red-500 hover:bg-red-500/10 rounded transition-colors"
        title="删除"
      >
        <Trash2 class="w-3.5 h-3.5" />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Box, GitBranch, Edit3, Trash2, History, Play, Power, Clock } from 'lucide-vue-next'
import CountdownDisplay from '@/components/CountdownDisplay'
import type { Workflow } from '@/types/workflow'

interface Props {
  workflow: Workflow
}

const props = defineProps<Props>()

const emit = defineEmits<{
  edit: [workflow: Workflow]
  executions: [workflow: Workflow]
  execute: [workflow: Workflow]
  delete: [workflow: Workflow]
  toggle: [workflow: Workflow]
  refresh: [workflow: Workflow] // 到达执行时间时触发刷新
}>()

// 判断是否有调度配置
const hasSchedule = computed(
  () =>
    props.workflow.enabled &&
    props.workflow.schedule_type &&
    props.workflow.schedule_type !== 'manual' &&
    props.workflow.next_run_time
)

const handleDelete = () => {
  if (confirm('确定要删除这个工作流吗？')) {
    emit('delete', props.workflow)
  }
}

const handleCountdownFinish = () => {
  // 到达执行时间，延迟 1 秒后触发刷新（让后端有时间更新）
  setTimeout(() => {
    emit('refresh', props.workflow)
  }, 1000)
}

const getNoScheduleText = (workflow: Workflow) => {
  if (!workflow.enabled) {
    return '工作流未启用'
  }
  if (!workflow.schedule_type || workflow.schedule_type === 'manual') {
    return '手动触发'
  }
  if (workflow.schedule_type === 'webhook') {
    return 'Webhook 触发'
  }
  return '暂无下次执行时间'
}
</script>
