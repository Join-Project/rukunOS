import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useLoadingStore = defineStore('loading', () => {
  const isLoading = ref(false)
  const message = ref('Memuat...')

  const setLoading = (loading: boolean, loadingMessage?: string) => {
    isLoading.value = loading
    if (loadingMessage) {
      message.value = loadingMessage
    } else if (!loading) {
      message.value = 'Memuat...'
    }
  }

  const show = (loadingMessage?: string) => {
    setLoading(true, loadingMessage)
  }

  const hide = () => {
    setLoading(false)
  }

  return {
    isLoading,
    message,
    setLoading,
    show,
    hide
  }
})

