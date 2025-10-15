<template>
  <div class="space-y-2">
    <BaseSelect
      :model-value="type"
      @update:model-value="handleTypeChange"
      :options="scheduleOptions"
      label="执行规则"
      required
    />

    <TimePicker
      v-if="type === 'daily'"
      :model-value="value"
      @update:model-value="handleValueChange"
      hint="每天在指定时间执行"
    />

    <BaseInput
      v-if="type === 'hourly'"
      :model-value="value"
      @update:model-value="$emit('update:value', $event)"
      placeholder="分:秒 (例如: 30:00)"
      hint="每小时的第N分N秒执行（最小间隔5分钟）"
      required
    />

    <BaseInput
      v-if="type === 'interval'"
      :model-value="value"
      @update:model-value="$emit('update:value', $event)"
      type="number"
      placeholder="秒数"
      :min="300"
      hint="每隔N秒执行一次（最小300秒，即5分钟）"
      required
    />

    <WeekDayPicker
      v-if="type === 'weekly'"
      :model-value="value"
      @update:model-value="$emit('update:value', $event)"
      hint="每周在选定的星期几执行"
    />

    <MonthDayPicker
      v-if="type === 'monthly'"
      :model-value="value"
      @update:model-value="$emit('update:value', $event)"
      hint="每月在选定的日期执行"
    />

    <BaseInput
      v-if="type === 'cron'"
      :model-value="value"
      @update:model-value="$emit('update:value', $event)"
      placeholder="0 0 * * * *"
      hint="Cron表达式: 秒 分 时 日 月 周"
      required
    />
  </div>
</template>

<script setup lang="ts">
import BaseInput from '@/components/BaseInput'
import BaseSelect from '@/components/BaseSelect'
import TimePicker from '@/components/TimePicker'
import WeekDayPicker from '@/components/WeekDayPicker'
import MonthDayPicker from '@/components/MonthDayPicker'

defineProps<{
  type: string
  value: string
}>()

const emit = defineEmits<{
  'update:type': [value: string]
  'update:value': [value: string]
}>()

const handleTypeChange = (newType: string) => {
  emit('update:type', newType)
}

const handleValueChange = (newValue: string) => {
  emit('update:value', newValue)
}

const scheduleOptions = [
  { label: '每天', value: 'daily' },
  { label: '每周', value: 'weekly' },
  { label: '每月', value: 'monthly' },
  { label: '每小时', value: 'hourly' },
  { label: '间隔', value: 'interval' },
  { label: 'Cron表达式', value: 'cron' },
]
</script>
