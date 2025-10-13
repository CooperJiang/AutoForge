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
        v-model="form.name"
        label="任务名称"
        placeholder="例如：每日签到"
        required
      />

      <div class="space-y-2">
        <BaseSelect
          v-model="form.scheduleType"
          :options="scheduleOptions"
          label="执行规则"
          required
        />

        <TimePicker
          v-if="form.scheduleType === 'daily'"
          v-model="form.scheduleValue"
          hint="每天在指定时间执行"
        />

        <BaseInput
          v-if="form.scheduleType === 'hourly'"
          v-model="form.scheduleValue"
          placeholder="分:秒 (例如: 30:00)"
          hint="每小时的第N分N秒执行（最小间隔5分钟）"
          required
        />

        <BaseInput
          v-if="form.scheduleType === 'interval'"
          v-model="form.scheduleValue"
          type="number"
          placeholder="秒数"
          :min="300"
          hint="每隔N秒执行一次（最小300秒，即5分钟）"
          required
        />

        <WeekDayPicker
          v-if="form.scheduleType === 'weekly'"
          v-model="form.scheduleValue"
          hint="每周在选定的星期几执行"
        />

        <MonthDayPicker
          v-if="form.scheduleType === 'monthly'"
          v-model="form.scheduleValue"
          hint="每月在选定的日期执行"
        />

        <BaseInput
          v-if="form.scheduleType === 'cron'"
          v-model="form.scheduleValue"
          placeholder="0 0 * * * *"
          hint="Cron表达式: 秒 分 时 日 月 周"
          required
        />
      </div>

      <div class="pt-3 border-t-2 border-border-primary space-y-2">
        <h3 class="text-xs font-semibold text-text-secondary">工具配置</h3>

        <BaseSelect
          v-model="form.tool_code"
          :options="toolOptions"
          label="选择工具"
          placeholder="请选择工具"
          required
          @change="$emit('tool-change')"
        />

        <div v-if="form.tool_code" class="space-y-2">
          <BaseButton
            variant="secondary"
            type="button"
            @click="$emit('open-config')"
            :full-width="true"
          >
            {{ isConfigured ? '✓ 已配置 - 点击修改' : '配置工具参数' }}
          </BaseButton>

          <BaseButton
            v-if="isConfigured"
            variant="ghost"
            type="button"
            @click="$emit('test-config')"
            :full-width="true"
          >
            🧪 测试配置
          </BaseButton>
        </div>
      </div>

      <div class="flex gap-2 pt-2">
        <BaseButton variant="primary" type="submit" :full-width="true" :disabled="!isConfigured">
          {{ editingTask ? '保存修改' : '创建任务' }}
        </BaseButton>
        <BaseButton v-if="editingTask" variant="secondary" type="button" @click="$emit('cancel')" :full-width="true">
          取消
        </BaseButton>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import BaseInput from '@/components/BaseInput'
import BaseButton from '@/components/BaseButton'
import BaseSelect from '@/components/BaseSelect'
import TimePicker from '@/components/TimePicker'
import WeekDayPicker from '@/components/WeekDayPicker'
import MonthDayPicker from '@/components/MonthDayPicker'

interface TaskForm {
  name: string
  scheduleType: string
  scheduleValue: string
  tool_code: string
  method: string
  url: string
  headers: { key: string; value: string }[]
  params: { key: string; value: string }[]
  body: string
}

defineProps<{
  form: TaskForm
  editingTask: any
  toolOptions: { label: string; value: string }[]
  isConfigured: boolean
}>()

defineEmits<{
  submit: []
  cancel: []
  'tool-change': []
  'open-config': []
  'test-config': []
}>()

const scheduleOptions = [
  { label: '每天', value: 'daily' },
  { label: '每周', value: 'weekly' },
  { label: '每月', value: 'monthly' },
  { label: '每小时', value: 'hourly' },
  { label: '间隔', value: 'interval' },
  { label: 'Cron表达式', value: 'cron' }
]

const handleSubmit = () => {
  // Emit submit event
}
</script>
