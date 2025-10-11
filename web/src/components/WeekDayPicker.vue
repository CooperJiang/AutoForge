<template>
  <div class="space-y-2">
    <label class="block text-sm font-medium text-slate-700">选择星期</label>
    <div class="grid grid-cols-7 gap-1">
      <button
        v-for="day in weekDays"
        :key="day.value"
        type="button"
        @click="toggleDay(day.value)"
        class="px-2 py-1.5 text-xs font-medium rounded border-2 transition-colors"
        :class="selectedDays.includes(day.value)
          ? 'bg-blue-500 text-white border-blue-500'
          : 'bg-white text-slate-700 border-slate-200 hover:border-blue-300'"
      >
        {{ day.label }}
      </button>
    </div>
    <TimePicker v-model="time" hint="执行时间" />
    <div v-if="hint" class="text-xs text-slate-500">{{ hint }}</div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import TimePicker from './TimePicker.vue'

interface Props {
  modelValue: string
  hint?: string
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const weekDays = [
  { label: '周日', value: 0 },
  { label: '周一', value: 1 },
  { label: '周二', value: 2 },
  { label: '周三', value: 3 },
  { label: '周四', value: 4 },
  { label: '周五', value: 5 },
  { label: '周六', value: 6 }
]

const selectedDays = ref<number[]>([])
const time = ref('09:00:00')

// 解析 modelValue: "1,3,5:09:00:00" -> days=[1,3,5], time="09:00:00"
const parseValue = (value: string) => {
  if (!value) {
    selectedDays.value = [1]
    time.value = '09:00:00'
    return
  }

  const parts = value.split(':')
  if (parts.length >= 4) {
    // days:HH:MM:SS
    const days = parts[0].split(',').map(d => parseInt(d)).filter(d => !isNaN(d))
    selectedDays.value = days.length > 0 ? days : [1]
    time.value = `${parts[1]}:${parts[2]}:${parts[3]}`
  } else {
    selectedDays.value = [1]
    time.value = '09:00:00'
  }
}

// 初始化
parseValue(props.modelValue)

const toggleDay = (day: number) => {
  const index = selectedDays.value.indexOf(day)
  if (index > -1) {
    if (selectedDays.value.length > 1) {
      selectedDays.value.splice(index, 1)
    }
  } else {
    selectedDays.value.push(day)
    selectedDays.value.sort((a, b) => a - b)
  }
  emitValue()
}

const emitValue = () => {
  const value = `${selectedDays.value.join(',')}:${time.value}`
  emit('update:modelValue', value)
}

watch(() => time.value, emitValue)
watch(() => props.modelValue, (newVal) => {
  parseValue(newVal)
})
</script>
