/**
 * 统一错误处理工具
 */

/**
 * 错误级别
 */
export enum ErrorLevel {
  /** 信息 */
  INFO = 'info',
  /** 警告 */
  WARN = 'warn',
  /** 错误 */
  ERROR = 'error',
  /** 致命错误 */
  FATAL = 'fatal',
}

/**
 * 错误上下文信息
 */
export interface ErrorContext {
  /** 错误发生的模块/组件名称 */
  module?: string
  /** 错误发生时的操作 */
  action?: string
  /** 额外的上下文数据 */
  data?: Record<string, unknown>
}

/**
 * 处理错误
 * @param error 错误对象
 * @param context 错误上下文
 * @param level 错误级别
 */
export function handleError(
  error: unknown,
  context?: ErrorContext,
  level: ErrorLevel = ErrorLevel.ERROR
) {
  const errorMessage = error instanceof Error ? error.message : String(error)
  const errorStack = error instanceof Error ? error.stack : undefined

  const logData = {
    level,
    message: errorMessage,
    stack: errorStack,
    context,
    timestamp: new Date().toISOString(),
  }

  // 开发环境：输出到控制台
  if (import.meta.env.DEV) {
    const logMethod = level === ErrorLevel.FATAL || level === ErrorLevel.ERROR ? 'error' : 'warn'
    console[logMethod]('[ErrorHandler]', logData)
  }

  if (import.meta.env.PROD && (level === ErrorLevel.ERROR || level === ErrorLevel.FATAL)) {
    // reportToMonitoring(logData)
  }
}

/**
 * 创建错误处理器（柯里化）
 * @param context 固定的错误上下文
 * @returns 错误处理函数
 *
 * @example
 * ```ts
 * const handleUploadError = createErrorHandler({ module: 'FileUpload', action: 'upload' })
 *
 * try {
 *   await uploadFile(file)
 * } catch (error) {
 *   handleUploadError(error, ErrorLevel.ERROR)
 * }
 * ```
 */
export function createErrorHandler(context: ErrorContext) {
  return (error: unknown, level: ErrorLevel = ErrorLevel.ERROR) => {
    handleError(error, context, level)
  }
}

/**
 * 异步函数错误包装器
 * @param fn 异步函数
 * @param context 错误上下文
 * @returns 包装后的函数
 *
 * @example
 * ```ts
 * const safeUpload = withErrorHandler(
 *   async (file: File) => {
 *     return await uploadFile(file)
 *   },
 *   { module: 'FileUpload', action: 'upload' }
 * )
 *
 * const result = await safeUpload(file)
 * ```
 */
export function withErrorHandler<T extends (...args: any[]) => Promise<any>>(
  fn: T,
  context: ErrorContext
): T {
  return (async (...args: Parameters<T>) => {
    try {
      return await fn(...args)
    } catch (error) {
      handleError(error, context)
      throw error
    }
  }) as T
}

