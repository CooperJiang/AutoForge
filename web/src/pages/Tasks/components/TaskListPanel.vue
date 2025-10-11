<template>
  <div>
    <div class="mb-4 flex gap-2 items-center">
      <div class="flex-1">
        <BaseInput
          :model-value="searchKeyword"
          @update:model-value="$emit('update:searchKeyword', $event)"
          placeholder="搜索任务名称或工具..."
        />
      </div>
      <div class="w-40">
        <BaseSelect
          :model-value="statusFilter"
          @update:model-value="$emit('update:statusFilter', $event)"
          :options="statusFilterOptions"
        />
      </div>
    </div>

    <div
      v-if="filteredTasks.length === 0"
      class="bg-white border-2 border-slate-200 rounded-lg p-8 text-center"
    >
      <p class="text-slate-500">{{ tasks.length === 0 ? '暂无任务，请添加一个定时任务' : '没有符合条件的任务' }}</p>
    </div>

    <div v-else class="space-y-3">
      <TaskCard
        v-for="task in filteredTasks"
        :key="task.id"
        :task="task"
        :executions="selectedTaskId === task.id ? executions : []"
        @test="$emit('test-task', task)"
        @edit="$emit('edit-task', task)"
        @clone="$emit('clone-task', task)"
        @enable="$emit('toggle-task', task)"
        @disable="$emit('toggle-task', task)"
        @delete="$emit('delete-task', task)"
        @refresh-executions="$emit('refresh-executions', task.id)"
        @delete-all-executions="$emit('delete-all-executions', task.id)"
        @delete-execution="$emit('delete-execution', $event)"
        @toggle-expanded="$emit('toggle-expanded', task.id)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import BaseInput from '@/components/BaseInput.vue'
import BaseSelect from '@/components/BaseSelect.vue'
import TaskCard from './TaskCard.vue'
import type { Task, TaskExecution } from '@/api/task'

defineProps<{
  tasks: Task[]
  filteredTasks: Task[]
  executions: TaskExecution[]
  selectedTaskId: string | null
  searchKeyword: string
  statusFilter: 'all' | 'enabled' | 'disabled'
}>()

defineEmits<{
  'update:searchKeyword': [value: string]
  'update:statusFilter': [value: 'all' | 'enabled' | 'disabled']
  'test-task': [task: Task]
  'edit-task': [task: Task]
  'clone-task': [task: Task]
  'toggle-task': [task: Task]
  'delete-task': [task: Task]
  'refresh-executions': [taskId: string]
  'delete-all-executions': [taskId: string]
  'delete-execution': [executionId: string]
  'toggle-expanded': [taskId: string]
}>()

const statusFilterOptions = [
  { label: '全部状态', value: 'all' },
  { label: '已启用', value: 'enabled' },
  { label: '已禁用', value: 'disabled' }
]
</script>
