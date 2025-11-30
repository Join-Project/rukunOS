import { useToastStore } from '~/stores/toast'

export const useToast = () => {
  const toastStore = useToastStore()

  const showSuccess = (message: string, title?: string, duration?: number) => {
    return toastStore.add({
      type: 'success',
      message,
      title,
      duration: duration ?? 3000,
    })
  }

  const showError = (message: string, title?: string, duration?: number) => {
    return toastStore.add({
      type: 'error',
      message,
      title: title ?? 'Terjadi Kesalahan',
      duration: duration ?? 5000,
    })
  }

  const showWarning = (message: string, title?: string, duration?: number) => {
    return toastStore.add({
      type: 'warning',
      message,
      title: title ?? 'Peringatan',
      duration: duration ?? 4000,
    })
  }

  const showInfo = (message: string, title?: string, duration?: number) => {
    return toastStore.add({
      type: 'info',
      message,
      title: title ?? 'Informasi',
      duration: duration ?? 3000,
    })
  }

  return {
    showSuccess,
    showError,
    showWarning,
    showInfo,
  }
}










