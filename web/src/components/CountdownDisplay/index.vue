<template>
  <span>{{ displayText }}</span>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useCountdown } from '@/composables/useCountdown'

interface Props {
  timestamp: number
  onFinish?: () => void
}

const props = defineProps<Props>()
const emit = defineEmits<{
  finish: []
}>()

const { countdown, isExpired } = useCountdown(() => props.timestamp, {
  autoStart: true,
  onFinish: () => {
    emit('finish')
    props.onFinish?.()
  },
})

const displayText = computed(() => {
  if (countdown.value === '已过期' || countdown.value === '' || isExpired.value) {
    return '等待重新计算...'
  }
  return countdown.value
})
</script>
