/**
 * 图片变换逻辑 Hook
 * 负责：缩放、旋转、平移计算
 */

import { ref, computed, type Ref } from 'vue'

interface UseImageTransformOptions {
  containerRef: Ref<HTMLElement | null>
  imageRef: Ref<HTMLImageElement | null>
  minScale?: number
  maxScale?: number
  scaleStep?: number
}

export function useImageTransform(options: UseImageTransformOptions) {
  const { containerRef, imageRef, minScale = 0.1, maxScale = 5, scaleStep = 0.2 } = options

  // 变换状态
  const scale = ref(1)
  const rotation = ref(0)
  const translateX = ref(0)
  const translateY = ref(0)

  // 计算样式
  const imageStyle = computed(() => ({
    // 先平移到中心，再缩放和旋转
    // 使用 translate(-50%, -50%) 将图片中心对齐到容器中心（left-1/2 top-1/2）
    // 然后应用自定义平移、缩放、旋转
    transform: `translate(-50%, -50%) translate(${translateX.value}px, ${translateY.value}px) scale(${scale.value}) rotate(${rotation.value}deg)`,
    transformOrigin: 'center center',
  }))

  // 放大
  const zoomIn = () => {
    const newScale = Math.min(scale.value + scaleStep, maxScale)
    scale.value = newScale
    centerImage()
  }

  // 缩小
  const zoomOut = () => {
    const newScale = Math.max(scale.value - scaleStep, minScale)
    scale.value = newScale
    centerImage()
  }

  // 左旋转
  const rotateLeft = () => {
    rotation.value -= 90
  }

  // 右旋转
  const rotateRight = () => {
    rotation.value += 90
  }

  // 重置
  const reset = () => {
    rotation.value = 0
    const fitScale = calculateFitScale()
    scale.value = fitScale
    centerImage()
  }

  // 设置缩放
  const setScale = (newScale: number) => {
    scale.value = Math.max(minScale, Math.min(newScale, maxScale))
    centerImage()
  }

  // 居中图片
  const centerImage = () => {
    // 由于使用了 translate(-50%, -50%)，图片已经自动居中
    // 只需将自定义平移重置为 0
    translateX.value = 0
    translateY.value = 0
  }

  // 计算适合容器的初始缩放比例
  const calculateFitScale = (): number => {
    if (!imageRef.value || !containerRef.value) return 1

    const img = imageRef.value
    if (img.naturalWidth === 0 || img.naturalHeight === 0) return 1

    const rect = containerRef.value.getBoundingClientRect()
    const containerWidth = rect.width
    const containerHeight = rect.height

    const maxWidth = containerWidth
    const maxHeight = containerHeight

    const imgWidth = img.naturalWidth
    const imgHeight = img.naturalHeight

    if (imgHeight > imgWidth) {
      return maxHeight / imgHeight
    } else {
      return maxWidth / imgWidth
    }
  }

  // 初始化
  const initialize = () => {
    if (!imageRef.value) return

    const img = imageRef.value
    if (img.complete && img.naturalWidth > 0) {
      // 图片已加载完成
      const fitScale = calculateFitScale()
      scale.value = fitScale
      centerImage()
    } else {
      // 等待图片加载
      img.addEventListener(
        'load',
        () => {
          const fitScale = calculateFitScale()
          scale.value = fitScale
          centerImage()
        },
        { once: true }
      )
    }
  }

  return {
    // 状态
    scale,
    rotation,
    translateX,
    translateY,
    imageStyle,

    // 方法
    zoomIn,
    zoomOut,
    rotateLeft,
    rotateRight,
    reset,
    setScale,
    centerImage,
    calculateFitScale,
    initialize,
  }
}
