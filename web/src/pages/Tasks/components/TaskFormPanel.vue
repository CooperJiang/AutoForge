<template>
  <div class="bg-bg-elevated border-2 border-border-primary rounded-lg shadow-sm p-4 sticky top-4">
    <div class="mb-3">
      <div class="flex items-center justify-between">
        <h2 class="text-base font-semibold text-text-primary">
          {{ editingTask ? '编辑任务' : '添加定时任务' }}
        </h2>
      </div>
    </div>

    <form @submit.prevent="handleSubmit" class="space-y-3 max-h-[calc(100vh-10rem)] overflow-y-auto pr-1">
      <BaseInput
        v-model="taskForm.name"
        label="任务名称"
        placeholder="例如：每日签到"
        required
      />

      <ScheduleSelector
        v-model:type="taskForm.scheduleType"
        v-model:value="taskForm.scheduleValue"
      />

      <div class="pt-3 border-t-2 border-border-primary space-y-2">
        <h3 class="text-xs font-semibold text-text-secondary">工具配置</h3>

        <BaseSelect
          v-model="taskForm.tool_code"
          :options="toolOptions"
          label="选择工具"
          placeholder="请选择工具"
          required
          @change="$emit('tool-change')"
        />

        <div v-if="taskForm.tool_code" class="space-y-2">
          <BaseButton
            variant="secondary"
            type="button"
            @click="$emit('config-click')"
            :full-width="true"
          >
            {{ isConfigured ? '✓ 已配置 - 点击修改' : '配置工具参数' }}
          </BaseButton>

          <BaseButton
            v-if="isConfigured"
            variant="ghost"
            type="button"
            @click="$emit('test-click')"
            :full-width="true"
            :disabled="testing"
          >
            {{ testing ? '测试中...' : '🧪 测试配置' }}
          </BaseButton>
        </div>
      </div>

      <div class="flex gap-2 pt-2">
        <BaseButton
          variant="primary"
          type="submit"
          :full-width="true"
          :disabled="!isConfigured"
        >
          {{ editingTask ? '保存修改' : '创建任务' }}
        </BaseButton>
        <BaseButton
          v-if="editingTask"
          variant="secondary"
          type="button"
          @click="$emit('cancel')"
          :full-width="true"
        >
          取消
        </BaseButton>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import BaseInput from '@/components/BaseInput'
import BaseSelect from '@/components/BaseSelect'
import BaseButton from '@/components/BaseButton'
import ScheduleSelector from './ScheduleSelector.vue'
import type { TaskFormData } from '@/composables/useTaskForm'
import type { Task } from '@/api/task'

defineProps<{
  taskForm: TaskFormData
  toolOptions: { label: string; value: string }[]
  isConfigured: boolean
  testing: boolean
  editingTask: Task | null
}>()

const emit = defineEmits<{
  submit: []
  cancel: []
  'tool-change': []
  'config-click': []
  'test-click': []
}>()

const handleSubmit = () => {
  emit('submit')
}
</script>
