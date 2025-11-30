import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface Toast {
  id: string
  type: 'success' | 'error' | 'warning' | 'info'
  title?: string
  message: string
  duration?: number
  progress: number
}

export const useToastStore = defineStore('toast', () => {
  const toasts = ref<Toast[]>([])
  const timers = new Map<string, NodeJS.Timeout>()
  const progressTimers = new Map<string, NodeJS.Timeout>()

  const add = (toast: Omit<Toast, 'id' | 'progress'>) => {
    const id = Date.now().toString() + Math.random().toString(36).substr(2, 9)
    const duration = toast.duration ?? 3000
    const newToast: Toast = {
      ...toast,
      id,
      progress: 100,
    }

    toasts.value.push(newToast)

    if (duration > 0) {
      // Progress bar animation
      const progressInterval = 50
      const progressStep = (progressInterval / duration) * 100
      
      const progressTimer = setInterval(() => {
        const index = toasts.value.findIndex(t => t.id === id)
        if (index !== -1) {
          toasts.value[index].progress = Math.max(0, toasts.value[index].progress - progressStep)
        }
      }, progressInterval)
      
      progressTimers.set(id, progressTimer)

      // Auto remove
      const timer = setTimeout(() => {
        remove(id)
      }, duration)
      
      timers.set(id, timer)
    }

    return id
  }

  const remove = (id: string) => {
    const index = toasts.value.findIndex(t => t.id === id)
    if (index !== -1) {
      toasts.value.splice(index, 1)
    }

    // Clear timers
    const timer = timers.get(id)
    if (timer) {
      clearTimeout(timer)
      timers.delete(id)
    }

    const progressTimer = progressTimers.get(id)
    if (progressTimer) {
      clearInterval(progressTimer)
      progressTimers.delete(id)
    }
  }

  const clear = () => {
    toasts.value.forEach(toast => {
      const timer = timers.get(toast.id)
      if (timer) clearTimeout(timer)
      const progressTimer = progressTimers.get(toast.id)
      if (progressTimer) clearInterval(progressTimer)
    })
    timers.clear()
    progressTimers.clear()
    toasts.value = []
  }

  return {
    toasts,
    add,
    remove,
    clear,
  }
})
