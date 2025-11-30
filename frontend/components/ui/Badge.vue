<template>
  <span :class="['badge', `badge-${variant}`]">
    <slot />
  </span>
</template>

<script setup lang="ts">
defineProps({
  variant: {
    type: String,
    default: 'primary',
    validator: (value: string) => ['primary', 'secondary', 'success', 'warning', 'danger', 'gray'].includes(value)
  }
})
</script>
