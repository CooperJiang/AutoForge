/**
 * 图片查看器相关类型定义
 */

/**
 * 图片缩略图组件属性
 */
export interface ImageThumbnailProps {
  /** 图片源地址 */
  src: string
  /** 图片替代文本 */
  alt?: string
  /** 图片描述 */
  description?: string
  /** 是否在图片中心显示描述 */
  showCenteredDescription?: boolean
}

/**
 * 图片缩略图组件事件
 */
export interface ImageThumbnailEmits {
  /** 点击事件 */
  (e: 'click'): void
  /** 加载错误事件 */
  (e: 'error', event: Event): void
}

/**
 * 图片工具栏组件属性
 */
export interface ImageToolbarProps {
  /** 当前缩放比例 */
  scale: number
}

/**
 * 图片工具栏组件事件
 */
export interface ImageToolbarEmits {
  /** 放大 */
  (e: 'zoom-in'): void
  /** 缩小 */
  (e: 'zoom-out'): void
  /** 逆时针旋转 */
  (e: 'rotate-left'): void
  /** 顺时针旋转 */
  (e: 'rotate-right'): void
  /** 重置 */
  (e: 'reset'): void
  /** 在新标签页打开 */
  (e: 'open-new-tab'): void
  /** 关闭 */
  (e: 'close'): void
}

/**
 * 图片查看器主组件属性
 */
export interface ImageViewerProps {
  /** 图片源地址 */
  src: string
  /** 图片替代文本 */
  alt?: string
  /** 图片描述 */
  description?: string
  /** 是否在图片中心显示描述 */
  showCenteredDescription?: boolean
}

/**
 * 图片变换状态
 */
export interface ImageTransformState {
  /** 缩放比例 */
  scale: number
  /** 旋转角度 */
  rotation: number
  /** X 轴偏移 */
  translateX: number
  /** Y 轴偏移 */
  translateY: number
}

