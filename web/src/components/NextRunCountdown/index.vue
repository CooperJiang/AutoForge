<template>
  <div class="text-text-secondary">
    <div v-if="!nextRunTime || nextRunTime === 0" class="text-text-tertiary">--</div>
    <div v-else-if="countdown <= 0" class="text-orange-600 font-medium text-sm">即将执行</div>
    <div v-else>
      <span class="text-sm font-medium" :class="timeColorClass" :title="formattedTime">{{
        formattedCountdown
      }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { formatCountdown, formatTime } from '@/utils/format'

const props = defineProps<{
  nextRunTime: number | null | undefined
}>()

const countdown = ref(0)
let timer: number | null = null

/**
 * 更新倒计时（秒）
 */
const updateCountdown = () => {
  if (!props.nextRunTime) {
    countdown.value = 0
    return
  }

  const now = Math.floor(Date.now() / 1000)
  countdown.value = props.nextRunTime - now
}

/**
 * 格式化倒计时显示
 */
const formattedCountdown = computed(() => {
  return formatCountdown(countdown.value)
})

/**
 * 格式化具体时间
 */
const formattedTime = computed(() => {
  return formatTime(props.nextRunTime, 'datetime')
})

/**
 * 根据剩余时间设置颜色
 */
const timeColorClass = computed(() => {
  if (countdown.value <= 0) return 'text-orange-600'
  if (countdown.value < 60) return 'text-red-600'
  if (countdown.value < 300) return 'text-orange-600'
  if (countdown.value < 3600) return 'text-yellow-600'
  return 'text-green-600'
})

onMounted(() => {
  updateCountdown()
  // 每秒更新倒计时
  timer = window.setInterval(updateCountdown, 1000)
})

onUnmounted(() => {
  if (timer) {
    clearInterval(timer)
    timer = null
  }
})
</script>
