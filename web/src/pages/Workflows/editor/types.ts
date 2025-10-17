/**
 * 工作流编辑器相关类型定义
 */

/**
 * 确认对话框状态
 */
export interface ConfirmDialogState {
  show: boolean
  title: string
  message: string
  variant?: 'danger' | 'warning' | 'info'
  confirmText?: string
  cancelText?: string
  resolve?: (value: boolean) => void
}

/**
 * 确认对话框选项
 */
export type ConfirmDialogOptions = Omit<ConfirmDialogState, 'show' | 'resolve'>

/**
 * 编辑器 UI 状态
 */
export interface EditorUIState {
  showConfigDrawer: boolean
  showEnvVarManager: boolean
  showAPISettings: boolean
  showImportDialog: boolean
  showExportDialog: boolean
  showExecuteDialog: boolean
  showPublishDialog: boolean
}

/**
 * 工作流验证结果
 */
export interface WorkflowValidation {
  valid: boolean
  message: string
}

/**
 * 用户角色
 */
export const USER_ROLE = {
  ADMIN: 'admin',
  USER: 'user',
} as const

export type UserRole = (typeof USER_ROLE)[keyof typeof USER_ROLE]

