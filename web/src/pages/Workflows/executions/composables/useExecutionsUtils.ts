import { Clock, Webhook, MousePointerClick, Play } from 'lucide-vue-next'
import type { WorkflowExecution } from '../types'
import { formatTimestamp } from '@/composables/useCountdown'

/**
 * 执行历史工具函数
 */
export function useExecutionsUtils() {
  // 状态相关
  const getStatusClass = (status: string) => {
    const classes = {
      running: 'bg-primary-light text-primary',
      success: 'bg-green-100 text-green-700',
      failed: 'bg-red-100 text-red-700',
      cancelled: 'bg-bg-tertiary text-text-secondary',
    }
    return classes[status as keyof typeof classes] || classes.cancelled
  }

  const getStatusDotClass = (status: string) => {
    const classes = {
      running: 'bg-[var(--color-primary)] animate-pulse',
      success: 'bg-green-500',
      failed: 'bg-red-500',
      cancelled: 'bg-bg-hover0',
    }
    return classes[status as keyof typeof classes] || classes.cancelled
  }

  const getStatusText = (status: string) => {
    const texts = {
      running: '运行中',
      success: '成功',
      failed: '失败',
      cancelled: '已取消',
    }
    return texts[status as keyof typeof texts] || '未知'
  }

  // 触发器相关
  const getTriggerIcon = (type: string) => {
    const icons = {
      schedule: Clock,
      scheduled: Clock,
      webhook: Webhook,
      manual: MousePointerClick,
    }
    return icons[type as keyof typeof icons] || Play
  }

  const getTriggerText = (type: string) => {
    const texts = {
      schedule: '定时触发',
      scheduled: '定时触发',
      webhook: 'Webhook',
      manual: '手动触发',
    }
    return texts[type as keyof typeof texts] || '未知'
  }

  // 进度计算
  const getProgress = (execution: WorkflowExecution) => {
    const completed = execution.success_nodes + execution.failed_nodes + execution.skipped_nodes
    const total = execution.total_nodes
    return total > 0 ? Math.round((completed / total) * 100) : 0
  }

  // 时间格式化
  const formatTime = (timestamp?: number) => {
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

  // 格式化下次执行时间
  const formatNextRunTime = formatTimestamp

  // 获取调度描述
  const getScheduleDescription = (type: string, value: string) => {
    if (!type || !value) return ''

    switch (type) {
      case 'daily':
        return `每天 ${value}`
      case 'weekly':
        const [days, time] = value.split(':').slice(0, 2).join(':').split(':')
        return `每周 ${days} ${time || ''}`
      case 'monthly':
        const [day, ...timeParts] = value.split(':')
        return `每月 ${day}号 ${timeParts.join(':')}`
      case 'hourly':
        return `每小时 ${value}`
      case 'interval':
        return `每 ${value} 秒`
      case 'cron':
        return `Cron: ${value}`
      default:
        return type
    }
  }

  // 格式化持续时间
  const formatDurationMs = (ms: number) => {
    const seconds = Math.floor(ms / 1000)
    const minutes = Math.floor(seconds / 60)
    const hours = Math.floor(minutes / 60)

    if (hours > 0) {
      return `${hours}小时${minutes % 60}分钟${seconds % 60}秒`
    } else if (minutes > 0) {
      return `${minutes}分钟${seconds % 60}秒`
    } else if (seconds > 0) {
      return `${seconds}秒`
    } else {
      return `${ms}毫秒`
    }
  }

  return {
    getStatusClass,
    getStatusDotClass,
    getStatusText,
    getTriggerIcon,
    getTriggerText,
    getProgress,
    formatTime,
    formatNextRunTime,
    getScheduleDescription,
    formatDurationMs,
  }
}

