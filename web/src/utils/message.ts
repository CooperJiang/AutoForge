import { createApp, h } from 'vue'
import MessageComponent from '@/components/Message.vue'

let messageInstance: any = null

const initMessage = () => {
  if (messageInstance) return messageInstance

  const container = document.createElement('div')
  document.body.appendChild(container)

  const app = createApp({
    render() {
      return h(MessageComponent, { ref: 'messageRef' })
    }
  })

  const instance = app.mount(container)
  messageInstance = (instance as any).$refs.messageRef
  return messageInstance
}

export const message = {
  success(content: string, duration = 3000) {
    const instance = initMessage()
    return instance.addMessage({ type: 'success', content, duration })
  },
  error(content: string, duration = 3000) {
    const instance = initMessage()
    return instance.addMessage({ type: 'error', content, duration })
  },
  warning(content: string, duration = 3000) {
    const instance = initMessage()
    return instance.addMessage({ type: 'warning', content, duration })
  },
  info(content: string, duration = 3000) {
    const instance = initMessage()
    return instance.addMessage({ type: 'info', content, duration })
  }
}
