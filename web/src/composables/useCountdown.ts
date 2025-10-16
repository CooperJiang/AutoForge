import { ref, computed, onUnmounted, watch } from 'vue'

export interface CountdownResult {
  countdown: string

  isExpired: boolean

  remainingMs: number

  remainingSeconds: number

  days: number
  hours: number
  minutes: number
  seconds: number
}

export interface UseCountdownOptions {
  autoStart?: boolean

  onFinish?: () => void

  onTick?: (result: CountdownResult) => void

  finishOffset?: number

  onlyUpdateOnSecondChange?: boolean
}

/**
 * 倒计时 Hook
 * @param targetTimestamp Unix 时间戳（秒）
 * @param options 配置选项
 */
export function useCountdown(
  targetTimestamp: number | (() => number),
  options: UseCountdownOptions = {}
) {
  const {
    autoStart = true,
    onFinish,
    onTick,
    finishOffset = 1000,
    onlyUpdateOnSecondChange = true,
  } = options

  const countdown = ref('')
  const isExpired = ref(false)
  const remainingMs = ref(0)
  const days = ref(0)
  const hours = ref(0)
  const minutes = ref(0)
  const seconds = ref(0)

  let timer: NodeJS.Timeout | null = null
  let hasTriggeredFinish = false
  let lastSecond = -1

  const getTargetTimestamp = () => {
    return typeof targetTimestamp === 'function' ? targetTimestamp() : targetTimestamp
  }

  const calculate = () => {
    const target = getTargetTimestamp()

    if (!target) {
      if (countdown.value !== '') {
        countdown.value = ''
        isExpired.value = false
        remainingMs.value = 0
        days.value = 0
        hours.value = 0
        minutes.value = 0
        seconds.value = 0
        lastSecond = -1
      }
      return
    }

    const now = Date.now()
    const targetMs = target * 1000
    const diff = targetMs - now
    const diffSeconds = Math.floor(diff / 1000)

    if (onlyUpdateOnSecondChange && diffSeconds === lastSecond) {
      return
    }

    lastSecond = diffSeconds

    remainingMs.value = diff

    if (diff <= 0) {
      if (countdown.value !== '已过期') {
        countdown.value = '已过期'
        isExpired.value = true
        days.value = 0
        hours.value = 0
        minutes.value = 0
        seconds.value = 0
      }

      if (!hasTriggeredFinish && onFinish) {
        hasTriggeredFinish = true
        onFinish()
      }

      return
    }

    if (diff <= finishOffset && !hasTriggeredFinish && onFinish) {
      hasTriggeredFinish = true

      setTimeout(() => {
        onFinish()
      }, diff)
    }

    days.value = Math.floor(diffSeconds / 86400)
    hours.value = Math.floor((diffSeconds % 86400) / 3600)
    minutes.value = Math.floor((diffSeconds % 3600) / 60)
    seconds.value = diffSeconds % 60

    const parts = []
    if (days.value > 0) parts.push(`${days.value}天`)
    if (hours.value > 0) parts.push(`${hours.value}小时`)
    if (minutes.value > 0) parts.push(`${minutes.value}分钟`)
    if (seconds.value > 0 || parts.length === 0) parts.push(`${seconds.value}秒`)

    const newCountdown = parts.join('') + '后'

    if (countdown.value !== newCountdown) {
      countdown.value = newCountdown
    }

    if (isExpired.value !== false) {
      isExpired.value = false
    }

    if (onTick) {
      onTick({
        countdown: countdown.value,
        isExpired: isExpired.value,
        remainingMs: remainingMs.value,
        remainingSeconds: diffSeconds,
        days: days.value,
        hours: hours.value,
        minutes: minutes.value,
        seconds: seconds.value,
      })
    }
  }

  const start = () => {
    if (timer) return

    hasTriggeredFinish = false
    calculate()
    timer = setInterval(calculate, 1000)
  }

  const stop = () => {
    if (timer) {
      clearInterval(timer)
      timer = null
    }
  }

  const reset = () => {
    stop()
    hasTriggeredFinish = false
    calculate()
  }

  const restart = () => {
    reset()
    start()
  }

  const remainingSeconds = computed(() => Math.floor(remainingMs.value / 1000))

  if (typeof targetTimestamp !== 'function') {
    watch(
      () => targetTimestamp,
      () => {
        if (timer) {
          restart()
        } else {
          calculate()
        }
      }
    )
  }

  if (autoStart) {
    start()
  }

  onUnmounted(() => {
    stop()
  })

  return {
    countdown,
    isExpired,
    remainingMs,
    remainingSeconds,
    days,
    hours,
    minutes,
    seconds,
    start,
    stop,
    reset,
    restart,
  }
}

/**
 * 格式化时间戳为本地时间字符串
 * @param timestamp Unix 时间戳（秒）
 */
export function formatTimestamp(timestamp: number): string {
  if (!timestamp) return '-'
  const date = new Date(timestamp * 1000)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  })
}

/**
 * 简化版本：只返回倒计时文本（用于简单场景）
 */
export function useSimpleCountdown(
  targetTimestamp: number | (() => number),
  onFinish?: () => void
) {
  const { countdown } = useCountdown(targetTimestamp, {
    autoStart: true,
    onFinish,
  })

  return countdown
}
