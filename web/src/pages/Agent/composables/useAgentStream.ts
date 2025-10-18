import { ref } from 'vue'
import type { AgentMessage, AgentPlan, AgentStep } from '@/api/agent'
import SecureStorage, { STORAGE_KEYS } from '@/utils/storage'

export interface StreamEvent {
  type: 'plan_start' | 'plan_step' | 'step_start' | 'step_end' | 'final' | 'error'
  data: any
}

export function useAgentStream() {
  const isStreaming = ref(false)
  const currentPlan = ref<AgentPlan | null>(null)
  const currentSteps = ref<AgentStep[]>([])
  const finalAnswer = ref('')
  const error = ref('')

  // 开始流式接收
  const startStream = (
    conversationId: string,
    messageContent: string,
    config: Record<string, any>,
    onComplete: (message: AgentMessage) => void
  ) => {
    isStreaming.value = true
    currentPlan.value = null
    currentSteps.value = []
    finalAnswer.value = ''
    error.value = ''

    const token = SecureStorage.getItem<string>(STORAGE_KEYS.AUTH_TOKEN)
    const baseURL = import.meta.env.VITE_API_BASE_URL || '/api/v1'
    const url = `${baseURL}/agent/conversations/${conversationId}/messages`

    // 使用 fetch 发送 POST 请求，并接收流式响应
    fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
        Accept: 'text/event-stream',
      },
      body: JSON.stringify({
        message: messageContent,
        config,
      }),
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`)
        }

        const reader = response.body?.getReader()
        const decoder = new TextDecoder()

        if (!reader) {
          throw new Error('无法获取响应流')
        }

        const readStream = () => {
          reader.read().then(({ done, value }) => {
            if (done) {
              isStreaming.value = false
              return
            }

            // 解析 SSE 数据
            const chunk = decoder.decode(value)
            const lines = chunk.split('\n')

            for (const line of lines) {
              if (line.startsWith('data: ')) {
                const data = line.substring(6).trim()
                if (data === '') continue

                try {
                  const event: StreamEvent = JSON.parse(data)
                  handleEvent(event, onComplete)
                } catch (e) {
                  console.error('解析事件失败:', e, data)
                }
              } else if (line.startsWith('event: done')) {
                isStreaming.value = false
                return
              }
            }

            readStream()
          })
        }

        readStream()
      })
      .catch((err) => {
        console.error('流式请求失败:', err)
        error.value = err.message || '请求失败'
        isStreaming.value = false
      })
  }

  // 处理事件
  const handleEvent = (event: StreamEvent, onComplete: (message: AgentMessage) => void) => {
    switch (event.type) {
      case 'plan_start':
        currentPlan.value = event.data.plan
        break

      case 'plan_step':
        if (currentPlan.value) {
          const stepIndex = event.data.step_index
          if (currentPlan.value.steps[stepIndex]) {
            currentPlan.value.steps[stepIndex].status = event.data.status
          }
        }
        break

      case 'step_start':
        // 步骤开始，可以显示加载动画
        break

      case 'step_end':
        // 步骤结束，添加到步骤列表
        const step = {
          step: event.data.step,
          observation: event.data.observation,
          tool_output: event.data.output,
          elapsed_ms: event.data.elapsed_ms,
          timestamp: new Date().toISOString(),
        }
        currentSteps.value.push(step)
        break

      case 'final':
        finalAnswer.value = event.data.answer
        isStreaming.value = false

        // 构建完整的 Agent 消息
        const agentMessage: AgentMessage = {
          id: '', // 将由后端填充
          conversation_id: '',
          role: 'agent',
          content: event.data.answer,
          trace: event.data.trace,
          token_usage: event.data.token_usage,
          status: 'completed',
          created_at: Date.now(),
        }

        onComplete(agentMessage)
        break

      case 'error':
        error.value = event.data.error
        isStreaming.value = false
        console.error('Agent 错误:', error.value)
        break
    }
  }

  // 停止流式接收
  const stopStream = () => {
    isStreaming.value = false
  }

  return {
    isStreaming,
    currentPlan,
    currentSteps,
    finalAnswer,
    error,
    startStream,
    stopStream,
  }
}



