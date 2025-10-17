import { ref, computed, watch, type Ref } from 'vue'
import type { FeishuConfig, SelectOption } from '../types'

/**
 * 飞书机器人配置状态管理
 */
export function useFeishuState(config: Ref<FeishuConfig>, envVars: Ref<any[]>) {
  const localConfig = ref<FeishuConfig>({ ...config.value })
  const showVariableHelper = ref(false)

  // 消息类型选项
  const msgTypeOptions: SelectOption[] = [
    { label: '文本消息', value: 'text' },
    { label: '富文本消息', value: 'post' },
    { label: '图片消息', value: 'image' },
    { label: '交互式卡片', value: 'interactive' },
  ]

  // 卡片模板选项
  const cardTemplateOptions: SelectOption[] = [
    { label: '通知卡片', value: 'notification' },
    { label: '更新卡片', value: 'update' },
    { label: '警告卡片', value: 'alert' },
    { label: '自定义 JSON', value: 'custom' },
  ]

  // 卡片状态选项
  const cardStatusOptions: SelectOption[] = [
    { label: '信息', value: 'info' },
    { label: '成功', value: 'success' },
    { label: '警告', value: 'warning' },
    { label: '错误', value: 'error' },
  ]

  // 格式化环境变量
  const formattedEnvVars = computed(() => {
    return envVars.value.map((v: any) => ({
      name: v.name,
      value: v.value,
    }))
  })

  // 监听配置变化
  watch(
    config,
    (newConfig) => {
      localConfig.value = { ...newConfig }
    },
    { deep: true }
  )

  // 监听本地配置变化，同步回父组件
  watch(
    localConfig,
    (newConfig) => {
      Object.assign(config.value, newConfig)
    },
    { deep: true }
  )

  return {
    localConfig,
    showVariableHelper,
    msgTypeOptions,
    cardTemplateOptions,
    cardStatusOptions,
    formattedEnvVars,
  }
}

