<template>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
  <UiToast />
  <UiConfirmDialog
    :is-open="confirmStore.isOpen"
    :title="confirmStore.title"
    :message="confirmStore.message"
    :confirm-text="confirmStore.confirmText"
    :cancel-text="confirmStore.cancelText"
    :type="confirmStore.type"
    @confirm="confirmStore.handleConfirm"
    @cancel="confirmStore.handleCancel"
  />
  <UiLoader
    :is-loading="loadingStore.isLoading"
    :message="loadingStore.message"
  />
</template>

<script setup>
const confirmStore = useConfirmStore()
const loadingStore = useLoadingStore()
</script>

