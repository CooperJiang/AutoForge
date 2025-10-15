/**
 * 格式化文件大小
 * @param bytes 字节数
 * @param decimals 小数位数
 */
export function formatFileSize(bytes: number, decimals = 2): string {
  if (bytes === 0) return '0 Bytes'

  const k = 1024
  const dm = decimals < 0 ? 0 : decimals
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']

  const i = Math.floor(Math.log(bytes) / Math.log(k))

  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i]
}

/**
 * 格式化日期
 * @param date 日期
 * @param format 格式
 */
export function formatDate(date: Date | string | number, format = 'YYYY-MM-DD HH:mm:ss'): string {
  const d = new Date(date)

  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const hours = String(d.getHours()).padStart(2, '0')
  const minutes = String(d.getMinutes()).padStart(2, '0')
  const seconds = String(d.getSeconds()).padStart(2, '0')

  return format
    .replace('YYYY', String(year))
    .replace('MM', month)
    .replace('DD', day)
    .replace('HH', hours)
    .replace('mm', minutes)
    .replace('ss', seconds)
}

/**
 * 格式化数字
 * @param num 数字
 * @param options 选项
 */
export function formatNumber(num: number, options?: Intl.NumberFormatOptions): string {
  return new Intl.NumberFormat('zh-CN', options).format(num)
}

/**
 * 格式化百分比
 * @param value 值 (0-1)
 * @param decimals 小数位数
 */
export function formatPercentage(value: number, decimals = 1): string {
  return (value * 100).toFixed(decimals) + '%'
}

/**
 * 格式化时长（毫秒转可读格式）
 * @param milliseconds 毫秒数
 */
export function formatDuration(milliseconds: number): string {
  const seconds = Math.floor(milliseconds / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)
  const days = Math.floor(hours / 24)

  if (days > 0) {
    return `${days}天${hours % 24}小时`
  }
  if (hours > 0) {
    return `${hours}小时${minutes % 60}分钟`
  }
  if (minutes > 0) {
    return `${minutes}分钟${seconds % 60}秒`
  }
  return `${seconds}秒`
}

/**
 * 遮挡用户ID中间部分
 * @param userId 用户ID
 * @param showStart 显示开头字符数，默认4
 * @param showEnd 显示结尾字符数，默认4
 */
export function maskUserId(userId: string, showStart = 4, showEnd = 4): string {
  if (!userId) return ''

  const len = userId.length

  // 如果ID太短，只遮挡中间部分
  if (len <= showStart + showEnd) {
    if (len <= 3) return userId // 太短则不遮挡
    const start = Math.floor(len / 3)
    const end = len - Math.floor(len / 3)
    return userId.slice(0, start) + '***' + userId.slice(end)
  }

  return userId.slice(0, showStart) + '***' + userId.slice(-showEnd)
}

/**
 * 格式化时间（Unix时间戳或ISO字符串）
 * @param time Unix时间戳（秒或毫秒）或ISO时间字符串
 * @param format 格式选项
 */
export function formatTime(
  time: number | string | null | undefined,
  format: 'full' | 'date' | 'time' | 'datetime' = 'datetime'
): string {
  if (!time) return '--'

  let timestamp: number
  if (typeof time === 'string') {
    timestamp = new Date(time).getTime()
  } else {
    // 如果是秒级时间戳，转换为毫秒
    timestamp = time < 10000000000 ? time * 1000 : time
  }

  const date = new Date(timestamp)

  const options: Intl.DateTimeFormatOptions = {
    year: format === 'full' ? 'numeric' : undefined,
    month: '2-digit',
    day: '2-digit',
    hour: format === 'date' ? undefined : '2-digit',
    minute: format === 'date' ? undefined : '2-digit',
    second: format === 'full' ? '2-digit' : undefined,
  }

  return date.toLocaleString('zh-CN', options)
}

/**
 * 格式化倒计时（秒转可读格式）
 * @param seconds 剩余秒数
 */
export function formatCountdown(seconds: number): string {
  if (seconds <= 0) return '0秒'

  if (seconds < 60) {
    return `${seconds}秒`
  } else if (seconds < 3600) {
    const mins = Math.floor(seconds / 60)
    const secs = seconds % 60
    return `${mins}分${secs}秒`
  } else if (seconds < 86400) {
    const hours = Math.floor(seconds / 3600)
    const mins = Math.floor((seconds % 3600) / 60)
    return `${hours}时${mins}分`
  } else {
    const days = Math.floor(seconds / 86400)
    const hours = Math.floor((seconds % 86400) / 3600)
    return `${days}天${hours}时`
  }
}

/**
 * 截断ID显示（保留前N位）
 * @param id ID字符串
 * @param length 保留长度，默认6
 */
export function truncateId(id: string, length = 6): string {
  if (!id || id.length <= length) return id
  return id.slice(0, length) + '...'
}
