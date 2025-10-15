<template>
  <div
    class="bg-bg-elevated border-2 border-border-primary rounded-lg shadow-sm hover:shadow-md transition-shadow"
  >
    <!-- 卡片头部 -->
    <div class="p-4 border-b border-border-primary">
      <div class="flex items-start justify-between">
        <div class="flex-1 min-w-0">
          <h3 class="text-lg font-semibold text-text-primary mb-1 truncate">{{ task.name }}</h3>
          <div class="text-sm text-text-secondary break-all">{{ getTaskUrl() }}</div>
        </div>
        <div class="flex gap-1.5 ml-4 flex-shrink-0">
          <BaseButton size="xs" variant="secondary" @click="$emit('test')"> ⚡ 测试 </BaseButton>
          <BaseButton size="xs" variant="secondary" @click="$emit('edit')"> 编辑 </BaseButton>
          <BaseButton size="xs" variant="secondary" @click="$emit('clone')"> 克隆 </BaseButton>
          <BaseButton v-if="task.enabled" size="xs" variant="danger" @click="$emit('disable')">
            禁用
          </BaseButton>
          <BaseButton v-else size="xs" variant="primary" @click="$emit('enable')">
            启用
          </BaseButton>
          <BaseButton size="xs" variant="danger" @click="confirmDelete">
            <Trash2 :size="12" />
          </BaseButton>
        </div>
      </div>
    </div>

    <!-- 卡片内容 -->
    <div class="px-4 py-3 bg-bg-hover">
      <div class="space-y-1.5 text-xs">
        <div class="flex items-center">
          <span class="text-text-tertiary w-16">方式:</span>
          <span class="font-medium text-text-secondary">{{ getMethod() }}</span>
        </div>
        <div class="flex items-center">
          <span class="text-text-tertiary w-16">状态:</span>
          <span
            :class="[
              'inline-flex items-center gap-1 px-2 py-0.5 rounded text-xs font-medium border',
              task.enabled
                ? 'bg-green-500/10 text-green-600 dark:text-green-400 border-green-500/20'
                : 'bg-slate-500/10 text-slate-600 dark:text-slate-400 border-slate-500/20',
            ]"
          >
            {{ task.enabled ? '✓ 已启用' : '✗ 已禁用' }}
          </span>
        </div>
        <div class="flex items-center">
          <span class="text-text-tertiary w-16">工具:</span>
          <span class="font-medium text-text-secondary">{{ task.tool_code }}</span>
        </div>
        <div class="flex items-center">
          <span class="text-text-tertiary w-16">调度:</span>
          <span class="font-medium text-text-secondary">{{ formatSchedule() }}</span>
        </div>
        <div class="flex items-center">
          <span class="text-text-tertiary w-16">下次执行:</span>
          <span class="font-medium text-primary">
            {{ task.next_run_time ? formatNextRunTime(task.next_run_time) : '-' }}
          </span>
        </div>
      </div>
    </div>

    <!-- 执行记录区域 -->
    <div class="border-t border-border-primary">
      <button
        @click="toggleExecutions"
        class="w-full px-4 py-2.5 flex items-center justify-between text-xs text-text-secondary hover:bg-bg-hover transition-colors"
      >
        <span class="flex items-center gap-2 font-medium">
          <RefreshCw :size="14" class="text-emerald-600" />
          {{ executionsExpanded ? '▼' : '▶' }} 执行记录
        </span>
        <span class="flex items-center gap-3">
          <button
            @click.stop="$emit('refresh-executions')"
            class="text-xs text-emerald-600 hover:text-emerald-700 font-medium transition-colors"
          >
            刷新
          </button>
          <button
            @click.stop="confirmDeleteAll"
            class="text-xs text-rose-600 hover:text-rose-700 font-medium transition-colors"
          >
            移除全部
          </button>
        </span>
      </button>

      <!-- 执行记录列表 -->
      <Transition
        enter-active-class="transition-all duration-200 ease-out"
        enter-from-class="opacity-0 max-h-0"
        enter-to-class="opacity-100 max-h-[400px]"
        leave-active-class="transition-all duration-200 ease-in"
        leave-from-class="opacity-100 max-h-[400px]"
        leave-to-class="opacity-0 max-h-0"
      >
        <div
          v-show="executionsExpanded"
          class="border-t border-border-primary bg-bg-elevated overflow-hidden"
        >
          <div
            v-if="executions.length === 0"
            class="px-4 py-8 text-center text-xs text-text-placeholder"
          >
            暂无执行记录
          </div>
          <div v-else class="max-h-96 overflow-y-auto">
            <div
              v-for="log in executions"
              :key="log.id"
              class="group px-4 py-2.5 border-b border-border-primary hover:bg-bg-hover transition-colors relative"
            >
              <div class="flex items-start justify-between gap-3">
                <div class="flex-1 min-w-0">
                  <div class="flex items-center gap-2 text-xs">
                    <span
                      :class="[
                        'inline-flex items-center px-1.5 py-0.5 rounded text-xs font-medium border',
                        log.status === 'success'
                          ? 'bg-green-500/10 text-green-600 dark:text-green-400 border-green-500/20'
                          : log.status === 'failed'
                            ? 'bg-red-500/10 text-red-600 dark:text-red-400 border-red-500/20'
                            : 'bg-yellow-500/10 text-yellow-600 dark:text-yellow-400 border-yellow-500/20',
                      ]"
                    >
                      {{
                        log.status === 'success'
                          ? '✓ 成功'
                          : log.status === 'failed'
                            ? '✗ 失败'
                            : '⏱ 超时'
                      }}
                    </span>
                    <span class="text-text-placeholder">|</span>
                    <span class="text-text-secondary">
                      {{ formatTimestamp(log.started_at) }}
                    </span>
                    <span v-if="log.response_status" class="text-text-placeholder">|</span>
                    <span v-if="log.response_status" class="text-text-secondary">
                      状态码: {{ log.response_status }}
                    </span>
                    <span v-if="log.duration_ms" class="text-text-placeholder">|</span>
                    <span v-if="log.duration_ms" class="text-text-secondary">
                      耗时: {{ log.duration_ms }}ms
                    </span>
                  </div>
                  <div
                    v-if="log.error_message"
                    class="mt-1.5 p-2 bg-rose-50 rounded text-xs text-rose-700"
                  >
                    {{ log.error_message }}
                  </div>
                  <div v-if="log.response_body" class="mt-2">
                    <details class="group/details">
                      <summary
                        class="cursor-pointer text-xs font-medium text-emerald-600 hover:text-emerald-700 select-none list-none"
                      >
                        <span class="inline-flex items-center gap-1">
                          <span
                            class="group-open/details:rotate-90 transition-transform inline-block"
                            >▶</span
                          >
                          查看响应内容 ({{ (log.response_body.length / 1024).toFixed(2) }} KB)
                        </span>
                      </summary>
                      <div class="mt-2 animate-in fade-in slide-in-from-top-1 duration-200">
                        <JsonViewer :content="log.response_body" />
                      </div>
                    </details>
                  </div>
                </div>
                <div class="opacity-0 group-hover:opacity-100 transition-opacity flex-shrink-0">
                  <button
                    @click="confirmDeleteExecution(log.id)"
                    class="p-1 text-text-placeholder hover:text-rose-600 hover:bg-rose-50 rounded transition-colors"
                  >
                    <Trash2 :size="12" />
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </Transition>
    </div>
  </div>

  <!-- 删除任务确认对话框 -->
  <Dialog
    v-model="showDeleteDialog"
    title="确认删除"
    :message="`确定要删除任务「${task.name}」吗？此操作不可恢复。`"
    confirm-text="删除"
    cancel-text="取消"
    confirm-variant="danger"
    @confirm="handleDelete"
  />

  <!-- 删除执行记录确认对话框 -->
  <Dialog
    v-model="showDeleteExecutionDialog"
    title="确认删除"
    message="确定要删除这条执行记录吗？此操作不可恢复。"
    confirm-text="删除"
    cancel-text="取消"
    confirm-variant="danger"
    @confirm="handleDeleteExecution"
  />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { RefreshCw, Trash2 } from 'lucide-vue-next'
import type { Task, TaskExecution } from '@/api/task'
import {
  formatScheduleValue,
  formatTimestamp,
  formatNextRunTime,
  getScheduleTypeName,
} from '@/utils/taskHelpers'
import BaseButton from '@/components/BaseButton'
import JsonViewer from '@/components/JsonViewer'
import Dialog from '@/components/Dialog'

interface Props {
  task: Task
  executions: TaskExecution[]
}

const props = defineProps<Props>()

const emit = defineEmits<{
  test: []
  edit: []
  clone: []
  enable: []
  disable: []
  delete: []
  'refresh-executions': []
  'delete-all-executions': []
  'delete-execution': [id: string]
  'toggle-expanded': []
}>()

const executionsExpanded = ref(false)
const showDeleteDialog = ref(false)
const showDeleteExecutionDialog = ref(false)
const deletingExecutionId = ref<string | null>(null)

const getTaskUrl = () => {
  try {
    const config = JSON.parse(props.task.config)
    return config.url || '-'
  } catch {
    return '-'
  }
}

const getMethod = () => {
  try {
    const config = JSON.parse(props.task.config)
    return config.method || 'GET'
  } catch {
    return 'GET'
  }
}

const formatSchedule = () => {
  const typeName = getScheduleTypeName(props.task.schedule_type)
  const value = formatScheduleValue(props.task.schedule_type, props.task.schedule_value)
  return `${typeName}: ${value}`
}

const toggleExecutions = () => {
  executionsExpanded.value = !executionsExpanded.value
  if (executionsExpanded.value) {
    emit('refresh-executions')
  }
}

const confirmDeleteAll = () => {
  if (confirm('确定要删除该任务的所有执行记录吗？')) {
    emit('delete-all-executions')
  }
}

const confirmDelete = () => {
  showDeleteDialog.value = true
}

const handleDelete = () => {
  emit('delete')
  showDeleteDialog.value = false
}

const confirmDeleteExecution = (executionId: string) => {
  deletingExecutionId.value = executionId
  showDeleteExecutionDialog.value = true
}

const handleDeleteExecution = () => {
  if (deletingExecutionId.value) {
    emit('delete-execution', deletingExecutionId.value)
  }
  showDeleteExecutionDialog.value = false
  deletingExecutionId.value = null
}
</script>
