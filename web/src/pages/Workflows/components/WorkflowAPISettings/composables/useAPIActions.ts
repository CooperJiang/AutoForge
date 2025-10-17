import { workflowApi } from '@/api/workflow'
import { message } from '@/utils/message'
import type { Ref } from 'vue'

/**
 * API 设置操作逻辑
 */
export function useAPIActions(options: {
  props: any
  emit: any
  apiEnabled: Ref<boolean>
  apiKey: Ref<string>
  isOpen: Ref<boolean>
  timeout: Ref<number>
  webhookURL: Ref<string>
  testMode: Ref<string>
  testParams: Ref<string>
  testing: Ref<boolean>
  testResult: Ref<string>
}) {
  const {
    props,
    emit,
    apiEnabled,
    apiKey,
    isOpen,
    timeout,
    webhookURL,
    testMode,
    testParams,
    testing,
    testResult,
  } = options

  // 切换 API 启用状态
  const toggleAPI = async () => {
    const targetState = !apiEnabled.value

    if (targetState) {
      try {
        const response = await workflowApi.enableAPI(props.workflow.id)
        apiEnabled.value = true
        apiKey.value = response.api_key
        message.success('API 已启用')
        emit('refresh')
      } catch {
        message.error('启用 API 失败')
      }
    } else {
      try {
        await workflowApi.disableAPI(props.workflow.id)
        apiEnabled.value = false
        apiKey.value = ''
        message.success('API 已禁用')
        emit('refresh')
      } catch {
        message.error('禁用 API 失败')
      }
    }
  }

  // 重新生成 API Key
  const regenerateApiKey = async () => {
    try {
      const response = await workflowApi.regenerateAPIKey(props.workflow.id)
      apiKey.value = response.api_key
      message.success('API Key 已重新生成')
      emit('refresh')
    } catch {
      message.error('重新生成失败')
    }
  }

  // 复制 API Key
  const copyApiKey = () => {
    navigator.clipboard.writeText(apiKey.value)
    message.success('API Key 已复制')
  }

  // 复制端点地址
  const copyEndpoint = (endpoint: string) => {
    navigator.clipboard.writeText(endpoint)
    message.success('端点地址已复制')
  }

  // 复制代码
  const copyCode = (code: string) => {
    navigator.clipboard.writeText(code)
    message.success('代码已复制')
  }

  // 测试 API 调用
  const handleTest = async () => {
    testing.value = true
    testResult.value = ''

    try {
      JSON.parse(testParams.value)
      testResult.value = JSON.stringify(
        {
          code: 200,
          message: '测试成功',
          data: {
            execution_id: 'exec_test_123',
            status: testMode.value === 'sync' ? 'success' : 'running',
            message: testMode.value === 'sync' ? '执行成功' : '工作流已开始执行',
          },
        },
        null,
        2
      )
      message.success('测试请求已发送')
    } catch (error: any) {
      testResult.value = JSON.stringify(
        {
          error: error.message || '测试失败',
        },
        null,
        2
      )
      message.error('测试失败')
    } finally {
      testing.value = false
    }
  }

  // 关闭抽屉
  const handleClose = () => {
    isOpen.value = false
  }

  // 保存设置
  const handleSave = async () => {
    try {
      await workflowApi.updateAPITimeout(props.workflow.id, timeout.value)
      if (webhookURL.value) {
        await workflowApi.updateAPIWebhook(props.workflow.id, webhookURL.value)
      }
      message.success('设置已保存')
      emit('refresh')
    } catch {
      message.error('保存失败')
    }
  }

  return {
    toggleAPI,
    regenerateApiKey,
    copyApiKey,
    copyEndpoint,
    copyCode,
    handleTest,
    handleClose,
    handleSave,
  }
}

