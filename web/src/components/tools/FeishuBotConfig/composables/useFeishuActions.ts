import type { Ref } from 'vue'
import type { FeishuConfig } from '../types'

/**
 * 飞书机器人配置操作逻辑
 */
export function useFeishuActions(localConfig: Ref<FeishuConfig>) {
  // 处理消息类型变更
  const handleMsgTypeChange = (msgType: string) => {
    // 清空其他类型的字段
    if (msgType === 'text') {
      localConfig.value.post_content = ''
      localConfig.value.image_url = ''
      localConfig.value.card_content = ''
    } else if (msgType === 'post') {
      localConfig.value.content = ''
      localConfig.value.image_url = ''
      localConfig.value.card_content = ''
    } else if (msgType === 'image') {
      localConfig.value.content = ''
      localConfig.value.post_content = ''
      localConfig.value.card_content = ''
    } else if (msgType === 'interactive') {
      localConfig.value.content = ''
      localConfig.value.post_content = ''
      localConfig.value.image_url = ''
    }
  }

  // 插入字段变量
  const insertFieldVariable = (nodeId: string, fieldName: string, textareaRef: any) => {
    if (!textareaRef) return

    const textarea = textareaRef
    const variable = `{{nodes.${nodeId}.${fieldName}}}`
    const start = textarea.selectionStart
    const end = textarea.selectionEnd
    const text = textarea.value

    const newText = text.substring(0, start) + variable + text.substring(end)

    // 根据当前消息类型更新对应字段
    if (localConfig.value.msg_type === 'text') {
      localConfig.value.content = newText
    } else if (localConfig.value.msg_type === 'post') {
      localConfig.value.post_content = newText
    } else if (localConfig.value.msg_type === 'interactive') {
      if (localConfig.value.card_template === 'custom') {
        localConfig.value.card_custom_json = newText
      } else {
        localConfig.value.card_content = newText
      }
    }

    // 设置光标位置
    setTimeout(() => {
      textarea.selectionStart = textarea.selectionEnd = start + variable.length
      textarea.focus()
    }, 0)
  }

  // 插入节点变量
  const insertNodeVariable = (nodeId: string, textareaRef: any) => {
    if (!textareaRef) return

    const textarea = textareaRef
    const variable = `{{nodes.${nodeId}}}`
    const start = textarea.selectionStart
    const end = textarea.selectionEnd
    const text = textarea.value

    const newText = text.substring(0, start) + variable + text.substring(end)

    if (localConfig.value.msg_type === 'text') {
      localConfig.value.content = newText
    } else if (localConfig.value.msg_type === 'post') {
      localConfig.value.post_content = newText
    } else if (localConfig.value.msg_type === 'interactive') {
      if (localConfig.value.card_template === 'custom') {
        localConfig.value.card_custom_json = newText
      } else {
        localConfig.value.card_content = newText
      }
    }

    setTimeout(() => {
      textarea.selectionStart = textarea.selectionEnd = start + variable.length
      textarea.focus()
    }, 0)
  }

  // 插入环境变量
  const insertEnvVariable = (key: string, textareaRef: any) => {
    if (!textareaRef) return

    const textarea = textareaRef
    const variable = `{{env.${key}}}`
    const start = textarea.selectionStart
    const end = textarea.selectionEnd
    const text = textarea.value

    const newText = text.substring(0, start) + variable + text.substring(end)

    if (localConfig.value.msg_type === 'text') {
      localConfig.value.content = newText
    } else if (localConfig.value.msg_type === 'post') {
      localConfig.value.post_content = newText
    } else if (localConfig.value.msg_type === 'interactive') {
      if (localConfig.value.card_template === 'custom') {
        localConfig.value.card_custom_json = newText
      } else {
        localConfig.value.card_content = newText
      }
    }

    setTimeout(() => {
      textarea.selectionStart = textarea.selectionEnd = start + variable.length
      textarea.focus()
    }, 0)
  }

  return {
    handleMsgTypeChange,
    insertFieldVariable,
    insertNodeVariable,
    insertEnvVariable,
  }
}

