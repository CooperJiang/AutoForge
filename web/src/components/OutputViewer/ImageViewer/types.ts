/**
 * ImageViewer 相关类型定义
 * 注意：基础类型已移至 @/types/image-viewer
 */

export type { ImageViewerProps, ImageTransformState } from '@/types/image-viewer'

export interface ImageDragState {
  isDragging: boolean
  startX: number
  startY: number
  initialTranslateX: number
  initialTranslateY: number
}

export interface ImageSize {
  width: number
  height: number
  naturalWidth: number
  naturalHeight: number
}

export interface ContainerSize {
  width: number
  height: number
}

export interface ImageTransformActions {
  zoomIn: () => void
  zoomOut: () => void
  rotateLeft: () => void
  rotateRight: () => void
  reset: () => void
  setScale: (scale: number) => void
  centerImage: () => void
}

export interface ImageInteractionActions {
  startDrag: (e: MouseEvent) => void
  onDrag: (e: MouseEvent) => void
  stopDrag: () => void
  onWheel: (e: WheelEvent) => void
}

