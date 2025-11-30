<template>
  <Teleport to="body">
    <Transition name="fade-fast">
      <div
        v-if="isLoading"
        class="fixed inset-0 z-[9999] flex items-center justify-center bg-white/60"
      >
        <div class="text-center">
          <!-- Spinner -->
          <div class="relative w-12 h-12 mx-auto mb-3">
            <div class="absolute inset-0 border-3 border-primary-200 rounded-full"></div>
            <div class="absolute inset-0 border-3 border-primary-600 rounded-full border-t-transparent animate-spin"></div>
          </div>
          <!-- Text -->
          <p class="text-xs font-medium text-gray-600">{{ message || 'Memuat...' }}</p>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
defineProps({
  isLoading: {
    type: Boolean,
    default: false
  },
  message: {
    type: String,
    default: 'Memuat...'
  }
})
</script>

<style scoped>
.fade-fast-enter-active,
.fade-fast-leave-active {
  transition: opacity 0.15s ease;
}

.fade-fast-enter-from,
.fade-fast-leave-to {
  opacity: 0;
}

.border-3 {
  border-width: 3px;
}
</style>

