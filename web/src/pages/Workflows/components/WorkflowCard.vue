<template>
  <div
    class="group rounded-xl border-2 border-border-secondary hover:border-primary transition-all duration-200 bg-bg-elevated shadow-md hover:shadow-lg"
  >
    <div class="p-3.5 pb-3">
      <div class="flex items-start justify-between gap-2 mb-2">
        <h3 class="text-sm font-bold text-text-primary truncate flex-1">
          {{ workflow.name }}
        </h3>
        <div
          :class="[
            'flex-shrink-0 px-2 py-0.5 text-xs font-semibold rounded-full',
            workflow.enabled
              ? 'bg-success-light text-success-text'
              : 'bg-bg-tertiary text-text-disabled',
          ]"
        >
          {{ workflow.enabled ? '启用' : '禁用' }}
        </div>
      </div>
      <p class="text-xs text-text-tertiary line-clamp-1">
        {{ workflow.description || '暂无描述' }}
      </p>
    </div>

    <div class="px-3.5 pb-2.5 flex items-center gap-4">
      <span class="flex items-center gap-1.5 text-xs text-text-secondary">
        <Box class="w-3.5 h-3.5" />
        <span class="font-medium">{{ workflow.nodes.length }}</span>
        <span class="text-text-placeholder">节点</span>
      </span>
      <span class="flex items-center gap-1.5 text-xs text-text-secondary">
        <GitBranch class="w-3.5 h-3.5" />
        <span class="font-medium">{{ workflow.edges.length }}</span>
        <span class="text-text-placeholder">连接</span>
      </span>
    </div>

    <div class="px-3.5 pb-3">
      <div
        v-if="hasSchedule"
        class="flex items-center gap-2 text-xs bg-info-light text-info-text px-2.5 py-2 rounded-md border border-info"
      >
        <Clock class="w-3.5 h-3.5 flex-shrink-0" />
        <span class="truncate font-medium">
          <CountdownDisplay
            :timestamp="workflow.next_run_time || 0"
            @finish="handleCountdownFinish"
          />
        </span>
      </div>
      <div
        v-else
        class="text-xs text-text-placeholder bg-bg-hover px-2.5 py-2 rounded-md flex items-center gap-2 border border-border-primary"
      >
        <Clock class="w-3.5 h-3.5" />
        {{ getNoScheduleText(workflow) }}
      </div>
    </div>

    <div class="px-3.5 pb-3.5 flex items-center gap-2">
      <button
        @click="$emit('edit', workflow)"
        class="flex-1 px-3 py-2 text-xs font-semibold text-text-secondary bg-bg-hover hover:bg-bg-active rounded-md transition-colors flex items-center justify-center gap-1.5"
      >
        <Edit3 class="w-3.5 h-3.5" />
        编辑
      </button>
      <button
        @click="$emit('executions', workflow)"
        class="flex-1 px-3 py-2 text-xs font-semibold text-text-secondary bg-bg-hover hover:bg-bg-active rounded-md transition-colors flex items-center justify-center gap-1.5"
      >
        <History class="w-3.5 h-3.5" />
        历史
      </button>
      <button
        @click="$emit('execute', workflow)"
        :disabled="workflow.nodes.length === 0"
        class="flex-1 px-3 py-2 text-xs font-semibold text-text-secondary bg-bg-hover hover:bg-primary hover:text-primary-text rounded-md transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-1.5 border border-border-primary hover:border-primary"
      >
        <Play class="w-3.5 h-3.5" />
        执行
      </button>
    </div>

    <div
      class="px-3.5 py-2.5 bg-bg-secondary border-t-2 border-border-primary flex items-center justify-between rounded-b-xl"
    >
      <div class="flex items-center gap-1.5">
        <Power :class="['w-3.5 h-3.5', workflow.enabled ? 'text-success' : 'text-text-disabled']" />
        <span
          :class="[
            'text-xs font-semibold',
            workflow.enabled ? 'text-success-text' : 'text-text-disabled',
          ]"
        >
          {{ workflow.enabled ? '运行中' : '已停止' }}
        </span>
      </div>
      <div class="flex items-center gap-1">
        <button
          @click="$emit('toggle', workflow)"
          :class="[
            'px-2.5 py-1 text-xs font-semibold rounded transition-colors',
            workflow.enabled
              ? 'text-text-secondary hover:bg-bg-active'
              : 'text-success-text hover:bg-success-light',
          ]"
        >
          {{ workflow.enabled ? '停止' : '启动' }}
        </button>
        <button
          @click="handleDelete"
          class="p-1.5 text-text-placeholder hover:text-error hover:bg-error-light rounded transition-colors"
          title="删除"
        >
          <Trash2 class="w-3.5 h-3.5" />
        </button>
      </div>
    </div>
  </div>

  <Dialog
    v-model="showDeleteDialog"
    title="确认删除"
    :message="`确定要删除工作流「${workflow.name}」吗？此操作不可恢复。`"
    confirm-text="删除"
    cancel-text="取消"
    confirm-variant="danger"
    @confirm="confirmDelete"
  />
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Box, GitBranch, Edit3, Trash2, History, Play, Power, Clock } from 'lucide-vue-next'
import CountdownDisplay from '@/components/CountdownDisplay'
import Dialog from '@/components/Dialog'
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

// 删除确认对话框
const showDeleteDialog = ref(false)

// 判断是否有调度配置
const hasSchedule = computed(
  () =>
    props.workflow.enabled &&
    props.workflow.schedule_type &&
    props.workflow.schedule_type !== 'manual' &&
    props.workflow.next_run_time
)

const handleDelete = () => {
  showDeleteDialog.value = true
}

const confirmDelete = () => {
  emit('delete', props.workflow)
  showDeleteDialog.value = false
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
