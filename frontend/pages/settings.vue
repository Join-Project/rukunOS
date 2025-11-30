<template>
  <div class="space-y-6">
      <!-- Page Header -->
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Pengaturan</h1>
        <p class="mt-1 text-sm text-gray-500">Kelola pengaturan akun dan komunitas Anda.</p>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <!-- Tenant Info -->
        <UiCard title="Informasi Komunitas">
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-500 mb-1">Nama Komunitas</label>
              <p class="text-gray-900 font-medium text-lg">{{ authStore.user?.tenant_name || 'Rukun Warga 01' }}</p>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-500 mb-1">Email Kontak</label>
              <p class="text-gray-900 font-medium">{{ authStore.user?.email || 'admin@rukunos.id' }}</p>
            </div>
          </div>
        </UiCard>

        <!-- User Info -->
        <UiCard title="Informasi Akun">
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-500 mb-1">Nama Lengkap</label>
              <p class="text-gray-900 font-medium text-lg">{{ authStore.user?.full_name || 'Admin User' }}</p>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-500 mb-1">Email</label>
              <p class="text-gray-900 font-medium">{{ authStore.user?.email || 'admin@rukunos.id' }}</p>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-500 mb-2">Izin Akses</label>
              <div class="flex flex-wrap gap-2">
                <!-- Mock permissions if store is empty for demo -->
                <UiBadge variant="primary">admin</UiBadge>
                <UiBadge variant="success">billing.manage</UiBadge>
                <UiBadge variant="warning">users.view</UiBadge>
              </div>
            </div>
          </div>
          <template #footer>
             <UiButton variant="secondary" size="sm">Edit Profil</UiButton>
          </template>
        </UiCard>
      </div>
    </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

const authStore = useAuthStore()
</script>
