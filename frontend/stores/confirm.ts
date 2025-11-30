import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useConfirmStore = defineStore('confirm', () => {
  const isOpen = ref(false)
  const title = ref('Konfirmasi')
  const message = ref('')
  const confirmText = ref('Ya')
  const cancelText = ref('Batal')
  const type = ref<'danger' | 'warning' | 'info'>('info')
  let resolvePromise: ((value: boolean) => void) | null = null

  const show = (
    msg: string,
    options?: {
      title?: string
      confirmText?: string
      cancelText?: string
      type?: 'danger' | 'warning' | 'info'
    }
  ): Promise<boolean> => {
    message.value = msg
    title.value = options?.title ?? 'Konfirmasi'
    confirmText.value = options?.confirmText ?? 'Ya'
    cancelText.value = options?.cancelText ?? 'Batal'
    type.value = options?.type ?? 'info'
    isOpen.value = true

    return new Promise((resolve) => {
      resolvePromise = resolve
    })
  }

  const handleConfirm = () => {
    isOpen.value = false
    if (resolvePromise) {
      resolvePromise(true)
      resolvePromise = null
    }
  }

  const handleCancel = () => {
    isOpen.value = false
    if (resolvePromise) {
      resolvePromise(false)
      resolvePromise = null
    }
  }

  return {
    isOpen,
    title,
    message,
    confirmText,
    cancelText,
    type,
    show,
    handleConfirm,
    handleCancel,
  }
})






