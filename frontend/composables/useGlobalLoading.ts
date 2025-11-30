import { watch } from 'vue'
import type { Ref } from 'vue'

/**
 * Composable untuk sync loading state lokal dengan global loading store
 * Gunakan ini di setiap halaman untuk menampilkan loading spinner global
 */
export const useGlobalLoading = (localLoading: Ref<boolean>, message?: string) => {
  const loadingStore = useLoadingStore()

  // Watch local loading state dan sync dengan global store
  watch(
    localLoading,
    (isLoading) => {
      if (isLoading) {
        loadingStore.show(message || 'Memuat...')
      } else {
        loadingStore.hide()
      }
    },
    { immediate: true }
  )

  return {
    loadingStore
  }
}

