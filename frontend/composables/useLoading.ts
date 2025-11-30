export const useLoading = () => {
  const loadingStore = useLoadingStore()

  const showLoading = (message?: string) => {
    loadingStore.show(message)
  }

  const hideLoading = () => {
    loadingStore.hide()
  }

  const withLoading = async <T>(
    fn: () => Promise<T>,
    message?: string
  ): Promise<T> => {
    try {
      showLoading(message)
      return await fn()
    } finally {
      hideLoading()
    }
  }

  return {
    isLoading: computed(() => loadingStore.isLoading),
    message: computed(() => loadingStore.message),
    showLoading,
    hideLoading,
    withLoading
  }
}

