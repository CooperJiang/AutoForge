import { ref, computed, onUnmounted, watch } from 'vue'

export interface CountdownResult {
  // 格式化的倒计时文本（如：2天3小时5分10秒后）
  countdown: string
  // 是否已过期
  isExpired: boolean
  // 剩余毫秒数
  remainingMs: number
  // 剩余秒数
  remainingSeconds: number
  // 时间组件
  days: number
  hours: number
  minutes: number
  seconds: number
}

export interface UseCountdownOptions {
  // 是否自动启动
  autoStart?: boolean
  // 倒计时结束回调（到达目标时间时触发）
  onFinish?: () => void
  // 倒计时更新回调（每秒触发）
  onTick?: (result: CountdownResult) => void
  // 提前触发时间（毫秒），默认 1000ms（提前1秒触发 onFinish）
  finishOffset?: number
  // 是否只在秒数变化时更新（避免频繁重渲染）
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
  let lastSecond = -1 // 记录上一次的秒数，避免不必要的更新

  const getTargetTimestamp = () => {
    return typeof targetTimestamp === 'function' ? targetTimestamp() : targetTimestamp
  }

  const calculate = () => {
    const target = getTargetTimestamp()

    if (!target) {
      // 只在值真正改变时更新
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

    // 如果只在秒数变化时更新，且秒数未变化，则跳过
    if (onlyUpdateOnSecondChange && diffSeconds === lastSecond) {
      return
    }

    lastSecond = diffSeconds

    remainingMs.value = diff

    // 已过期
    if (diff <= 0) {
      if (countdown.value !== '已过期') {
        countdown.value = '已过期'
        isExpired.value = true
        days.value = 0
        hours.value = 0
        minutes.value = 0
        seconds.value = 0
      }

      // 触发完成回调（只触发一次）
      if (!hasTriggeredFinish && onFinish) {
        hasTriggeredFinish = true
        onFinish()
      }

      return
    }

    // 即将到达目标时间（提前触发）
    if (diff <= finishOffset && !hasTriggeredFinish && onFinish) {
      hasTriggeredFinish = true
      // 延迟到准确时间再触发
      setTimeout(() => {
        onFinish()
      }, diff)
    }

    days.value = Math.floor(diffSeconds / 86400)
    hours.value = Math.floor((diffSeconds % 86400) / 3600)
    minutes.value = Math.floor((diffSeconds % 3600) / 60)
    seconds.value = diffSeconds % 60

    // 构建倒计时字符串
    const parts = []
    if (days.value > 0) parts.push(`${days.value}天`)
    if (hours.value > 0) parts.push(`${hours.value}小时`)
    if (minutes.value > 0) parts.push(`${minutes.value}分钟`)
    if (seconds.value > 0 || parts.length === 0) parts.push(`${seconds.value}秒`)

    const newCountdown = parts.join('') + '后'

    // 只在文本真正改变时更新
    if (countdown.value !== newCountdown) {
      countdown.value = newCountdown
    }

    if (isExpired.value !== false) {
      isExpired.value = false
    }

    // 触发 tick 回调
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

  // 计算属性：剩余秒数
  const remainingSeconds = computed(() => Math.floor(remainingMs.value / 1000))

  // 监听目标时间变化，重新计算
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

  // 自动启动
  if (autoStart) {
    start()
  }

  // 清理
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
