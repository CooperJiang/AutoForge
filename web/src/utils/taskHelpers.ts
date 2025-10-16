/**
 * 任务相关的工具函数
 */

/**
 * 获取调度类型的显示名称
 */
export const getScheduleTypeName = (type: string): string => {
  const typeMap: Record<string, string> = {
    daily: '每天',
    weekly: '每周',
    monthly: '每月',
    hourly: '每小时',
    interval: '间隔',
    cron: 'Cron',
  }
  return typeMap[type] || type
}

/**
 * 格式化调度值显示
 */
export const formatScheduleValue = (type: string, value: string): string => {
  switch (type) {
    case 'daily': {
      return `每天 ${value}`
    }
    case 'weekly': {
      const parts = value.split(':')
      if (parts.length >= 4) {
        const days = parts[0].split(',').map((d) => parseInt(d))
        const time = `${parts[1]}:${parts[2]}:${parts[3]}`
        const dayNames = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
        const dayText = days.map((d) => dayNames[d] || d).join('、')
        return `每周${dayText} ${time}`
      }
      return value
    }
    case 'monthly': {
      const parts = value.split(':')
      if (parts.length >= 4) {
        const day = parts[0]
        const time = `${parts[1]}:${parts[2]}:${parts[3]}`
        return `每月${day}日 ${time}`
      }
      return value
    }
    case 'hourly':
      return `每小时 ${value}`
    case 'interval':
      return `每隔 ${value} 秒`
    case 'cron':
      return value
    default:
      return value
  }
}

/**
 * 格式化时间戳
 */
export const formatTimestamp = (timestamp: number): string => {
  const date = new Date(timestamp * 1000)
  return date.toLocaleString('zh-CN')
}

/**
 * 格式化下次运行时间
 */
export const formatNextRunTime = (timestamp: number): string => {
  const now = new Date()
  const nextRun = new Date(timestamp * 1000)
  const diff = nextRun.getTime() - now.getTime()

  if (diff < 0) {
    return '即将执行'
  }

  const seconds = Math.floor(diff / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)

  if (seconds < 60) {
    return `${seconds}秒后`
  }

  if (minutes < 60) {
    return `${minutes}分钟后`
  }

  if (hours < 24) {
    return `${hours}小时${minutes % 60}分钟后`
  }

  const year = nextRun.getFullYear()
  const month = String(nextRun.getMonth() + 1).padStart(2, '0')
  const day = String(nextRun.getDate()).padStart(2, '0')
  const hour = String(nextRun.getHours()).padStart(2, '0')
  const minute = String(nextRun.getMinutes()).padStart(2, '0')

  if (year === now.getFullYear()) {
    return `${month}-${day} ${hour}:${minute}`
  }
  return `${year}-${month}-${day} ${hour}:${minute}`
}

/**
 * 格式化文件大小
 */
export const formatSize = (bytes: number): string => {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}

/**
 * 获取调度值的默认值
 */
export const getDefaultScheduleValue = (type: string): string => {
  const defaultValues: Record<string, string> = {
    daily: '09:00:00',
    weekly: '1:09:00:00',
    monthly: '1:09:00:00',
    hourly: '05:00',
    interval: '300',
    cron: '0 0 * * * *',
  }
  return defaultValues[type] || ''
}
