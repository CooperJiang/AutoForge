<template>
  <Teleport to="body">
    <div class="fixed top-4 left-1/2 transform -translate-x-1/2 z-50 flex flex-col gap-2 pointer-events-none">
      <TransitionGroup name="message">
        <div
          v-for="msg in messages"
          :key="msg.id"
          class="message-item px-4 py-2.5 rounded-lg shadow-lg border-2 flex items-center gap-2 min-w-[300px] max-w-[500px] pointer-events-auto"
          :class="getMessageClass(msg.type)"
        >
          <!-- Icon -->
          <div class="flex-shrink-0">
            <CheckCircle v-if="msg.type === 'success'" :size="20" />
            <XCircle v-else-if="msg.type === 'error'" :size="20" />
            <AlertTriangle v-else-if="msg.type === 'warning'" :size="20" />
            <Info v-else :size="20" />
          </div>

          <!-- Content -->
          <div class="flex-1 text-sm font-medium">{{ msg.content }}</div>

          <!-- Close button -->
          <button
            @click="removeMessage(msg.id)"
            class="flex-shrink-0 hover:opacity-70 transition-opacity"
          >
            <X :size="16" />
          </button>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { CheckCircle, XCircle, AlertTriangle, Info, X } from 'lucide-vue-next'

export interface MessageItem {
  id: number
  type: 'success' | 'error' | 'warning' | 'info'
  content: string
  duration: number
}

const messages = ref<MessageItem[]>()
let messageId = 0

const getMessageClass = (type: string) => {
  switch (type) {
    case 'success':
      return 'bg-emerald-50 border-emerald-400 text-emerald-700'
    case 'error':
      return 'bg-rose-50 border-rose-400 text-rose-700'
    case 'warning':
      return 'bg-amber-50 border-amber-400 text-amber-700'
    default:
      return 'bg-blue-50 border-blue-400 text-blue-700'
  }
}

const addMessage = (msg: Omit<MessageItem, 'id'>) => {
  const id = messageId++
  const message = { ...msg, id }

  if (!messages.value) {
    messages.value = []
  }

  messages.value.push(message)

  if (msg.duration > 0) {
    setTimeout(() => {
      removeMessage(id)
    }, msg.duration)
  }

  return id
}

const removeMessage = (id: number) => {
  if (!messages.value) return
  const index = messages.value.findIndex(m => m.id === id)
  if (index > -1) {
    messages.value.splice(index, 1)
  }
}

defineExpose({
  addMessage,
  removeMessage
})
</script>

<style scoped>
.message-enter-active,
.message-leave-active {
  transition: all 0.3s ease;
}

.message-enter-from {
  opacity: 0;
  transform: translateY(-20px);
}

.message-leave-to {
  opacity: 0;
  transform: translateX(100px);
}
</style>
