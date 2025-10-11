import { message } from './message'

/**
 * 复制文本到剪贴板
 * @param text 要复制的文本
 * @param successMessage 成功提示消息，默认"已复制到剪贴板"
 */
export async function copyToClipboard(
  text: string,
  successMessage = '已复制到剪贴板'
): Promise<boolean> {
  try {
    await navigator.clipboard.writeText(text)
    if (successMessage) {
      message.success(successMessage)
    }
    return true
  } catch (error) {
    console.error('复制失败:', error)
    message.error('复制失败')
    return false
  }
}

/**
 * 从剪贴板读取文本
 */
export async function readFromClipboard(): Promise<string | null> {
  try {
    const text = await navigator.clipboard.readText()
    return text
  } catch (error) {
    console.error('读取剪贴板失败:', error)
    message.error('读取剪贴板失败')
    return null
  }
}
