import { ref } from 'vue'

export interface ConfirmOptions {
  title?: string
  message: string
  confirmText?: string
  cancelText?: string
  variant?: 'info' | 'warning' | 'danger' | 'question'
}

interface ConfirmState extends ConfirmOptions {
  show: boolean
  resolve?: (value: boolean) => void
}

const state = ref<ConfirmState>({
  show: false,
  message: ''
})

export function useConfirm() {
  const confirm = (options: ConfirmOptions): Promise<boolean> => {
    return new Promise((resolve) => {
      state.value = {
        ...options,
        show: true,
        resolve
      }
    })
  }

  const handleConfirm = () => {
    state.value.resolve?.(true)
    state.value.show = false
  }

  const handleCancel = () => {
    state.value.resolve?.(false)
    state.value.show = false
  }

  return {
    state,
    confirm,
    handleConfirm,
    handleCancel
  }
}
