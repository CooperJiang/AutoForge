<template>
  <div class="space-y-2">
    <label class="block text-sm font-medium text-text-secondary">选择日期</label>
    <div class="grid grid-cols-7 gap-1">
      <button
        v-for="day in 31"
        :key="day"
        type="button"
        @click="selectDay(day)"
        class="px-2 py-1.5 text-xs font-medium rounded border-2 transition-colors"
        :class="selectedDay === day
          ? 'bg-primary text-white border-primary'
          : 'bg-bg-elevated text-text-secondary border-border-primary hover:border-primary'"
      >
        {{ day }}
      </button>
    </div>
    <TimePicker v-model="time" hint="执行时间" />
    <div v-if="hint" class="text-xs text-text-tertiary">{{ hint }}</div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import TimePicker from '../TimePicker/index.vue'

interface Props {
  modelValue: string
  hint?: string
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const selectedDay = ref(1)
const time = ref('09:00:00')

// 解析 modelValue: "15:09:00:00" -> day=15, time="09:00:00"
const parseValue = (value: string) => {
  if (!value) {
    selectedDay.value = 1
    time.value = '09:00:00'
    return
  }

  const parts = value.split(':')
  if (parts.length >= 4) {
    // day:HH:MM:SS
    const day = parseInt(parts[0])
    selectedDay.value = !isNaN(day) && day >= 1 && day <= 31 ? day : 1
    time.value = `${parts[1]}:${parts[2]}:${parts[3]}`
  } else {
    selectedDay.value = 1
    time.value = '09:00:00'
  }
}

// 初始化
parseValue(props.modelValue)

const selectDay = (day: number) => {
  selectedDay.value = day
  emitValue()
}

const emitValue = () => {
  const value = `${selectedDay.value}:${time.value}`
  emit('update:modelValue', value)
}

watch(() => time.value, emitValue)
watch(() => props.modelValue, (newVal) => {
  parseValue(newVal)
})
</script>
