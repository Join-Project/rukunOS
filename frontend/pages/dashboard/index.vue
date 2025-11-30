<template>
  <div>
    <!-- Admin Dashboard -->
    <template v-if="isAdmin">
      <!-- Page Header -->
      <div class="mb-6 md:mb-8">
        <h1 class="text-2xl md:text-3xl font-bold text-gray-900 mb-2">Dashboard Admin</h1>
        <p class="text-sm md:text-base text-gray-600">Ringkasan aktivitas dan statistik komunitas</p>
      </div>

      <!-- Stats Grid -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 md:gap-6 mb-6 md:mb-8">
        <div class="card p-4 md:p-6">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-xs md:text-sm font-medium text-gray-600 mb-1">Total Unit</p>
              <p class="text-2xl md:text-3xl font-bold text-gray-900">{{ stats.totalUnits || 0 }}</p>
              <p class="text-xs text-gray-500 mt-1">Unit properti</p>
            </div>
            <div class="w-12 h-12 md:w-14 md:h-14 bg-primary-50 rounded-xl flex items-center justify-center">
              <Icon name="heroicons:building-office" class="w-6 h-6 md:w-7 md:h-7 text-primary-600" />
            </div>
          </div>
        </div>

        <div class="card p-4 md:p-6">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-xs md:text-sm font-medium text-gray-600 mb-1">Total Pengguna</p>
              <p class="text-2xl md:text-3xl font-bold text-gray-900">{{ stats.totalUsers || 0 }}</p>
              <p class="text-xs text-gray-500 mt-1">Warga terdaftar</p>
            </div>
            <div class="w-12 h-12 md:w-14 md:h-14 bg-green-50 rounded-xl flex items-center justify-center">
              <Icon name="heroicons:users" class="w-6 h-6 md:w-7 md:h-7 text-green-600" />
            </div>
          </div>
        </div>

        <div class="card p-4 md:p-6">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-xs md:text-sm font-medium text-gray-600 mb-1">Role Aktif</p>
              <p class="text-2xl md:text-3xl font-bold text-gray-900">{{ stats.totalRoles || 0 }}</p>
              <p class="text-xs text-gray-500 mt-1">Role tersedia</p>
            </div>
            <div class="w-12 h-12 md:w-14 md:h-14 bg-yellow-50 rounded-xl flex items-center justify-center">
              <Icon name="heroicons:shield-check" class="w-6 h-6 md:w-7 md:h-7 text-yellow-600" />
            </div>
          </div>
        </div>

        <div class="card p-4 md:p-6">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-xs md:text-sm font-medium text-gray-600 mb-1">Komunitas</p>
              <p class="text-base md:text-lg font-bold text-gray-900 truncate">
                {{ authStore.user?.tenant_name || '-' }}
              </p>
              <p class="text-xs text-gray-500 mt-1">Nama RT/RW</p>
            </div>
            <div class="w-12 h-12 md:w-14 md:h-14 bg-purple-50 rounded-xl flex items-center justify-center">
              <Icon name="heroicons:home" class="w-6 h-6 md:w-7 md:h-7 text-purple-600" />
            </div>
          </div>
        </div>
      </div>

      <!-- Content Grid -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-4 md:gap-6">
        <!-- Quick Actions -->
        <div class="lg:col-span-2">
          <div class="card">
            <div class="p-4 md:p-6 border-b border-gray-200">
              <h2 class="text-base md:text-lg font-semibold text-gray-900">Akses Cepat</h2>
              <p class="text-xs md:text-sm text-gray-600 mt-1">Fitur-fitur yang sering digunakan</p>
            </div>
            <div class="p-4 md:p-6">
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 md:gap-4">
                <NuxtLink
                  to="/units"
                  class="group p-4 md:p-6 border-2 border-gray-200 rounded-xl hover:border-primary-500 hover:bg-primary-50 transition-all duration-200"
                >
                  <div class="w-10 h-10 md:w-12 md:h-12 bg-primary-100 rounded-xl flex items-center justify-center mb-3 md:mb-4 group-hover:bg-primary-200 transition-colors">
                    <Icon name="heroicons:building-office" class="w-5 h-5 md:w-6 md:h-6 text-primary-600" />
                  </div>
                  <h3 class="text-sm md:text-base font-semibold text-gray-900 mb-1">Unit</h3>
                  <p class="text-xs md:text-sm text-gray-600">Kelola unit properti</p>
                </NuxtLink>

                <NuxtLink
                  to="/users"
                  class="group p-4 md:p-6 border-2 border-gray-200 rounded-xl hover:border-primary-500 hover:bg-primary-50 transition-all duration-200"
                >
                  <div class="w-10 h-10 md:w-12 md:h-12 bg-green-100 rounded-xl flex items-center justify-center mb-3 md:mb-4 group-hover:bg-green-200 transition-colors">
                    <Icon name="heroicons:users" class="w-5 h-5 md:w-6 md:h-6 text-green-600" />
                  </div>
                  <h3 class="text-sm md:text-base font-semibold text-gray-900 mb-1">Pengguna</h3>
                  <p class="text-xs md:text-sm text-gray-600">Kelola warga</p>
                </NuxtLink>

                <NuxtLink
                  to="/roles"
                  class="group p-4 md:p-6 border-2 border-gray-200 rounded-xl hover:border-primary-500 hover:bg-primary-50 transition-all duration-200"
                >
                  <div class="w-10 h-10 md:w-12 md:h-12 bg-yellow-100 rounded-xl flex items-center justify-center mb-3 md:mb-4 group-hover:bg-yellow-200 transition-colors">
                    <Icon name="heroicons:shield-check" class="w-5 h-5 md:w-6 md:h-6 text-yellow-600" />
                  </div>
                  <h3 class="text-sm md:text-base font-semibold text-gray-900 mb-1">Role</h3>
                  <p class="text-xs md:text-sm text-gray-600">Kelola role & izin</p>
                </NuxtLink>

                <NuxtLink
                  to="/settings"
                  class="group p-4 md:p-6 border-2 border-gray-200 rounded-xl hover:border-primary-500 hover:bg-primary-50 transition-all duration-200"
                >
                  <div class="w-10 h-10 md:w-12 md:h-12 bg-purple-100 rounded-xl flex items-center justify-center mb-3 md:mb-4 group-hover:bg-purple-200 transition-colors">
                    <Icon name="heroicons:cog-6-tooth" class="w-5 h-5 md:w-6 md:h-6 text-purple-600" />
                  </div>
                  <h3 class="text-sm md:text-base font-semibold text-gray-900 mb-1">Pengaturan</h3>
                  <p class="text-xs md:text-sm text-gray-600">Konfigurasi sistem</p>
                </NuxtLink>
              </div>
            </div>
          </div>
        </div>

        <!-- Recent Activity -->
        <div class="card">
          <div class="p-4 md:p-6 border-b border-gray-200">
            <h2 class="text-base md:text-lg font-semibold text-gray-900">Aktivitas Terkini</h2>
          </div>
          <div class="p-4 md:p-6">
            <div class="space-y-4">
              <div v-for="i in 5" :key="i" class="flex items-start space-x-3">
                <div class="flex-shrink-0 w-2 h-2 bg-primary-500 rounded-full mt-2"></div>
                <div class="flex-1 min-w-0">
                  <p class="text-sm text-gray-900">Aktivitas sistem {{ i }}</p>
                  <p class="text-xs text-gray-500 mt-0.5">{{ i }} jam yang lalu</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- Warga Dashboard -->
    <template v-else-if="isResident">
      <div class="max-w-5xl mx-auto space-y-8">
        <!-- Header Section -->
        <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div>
            <h1 class="text-xl md:text-2xl font-bold text-gray-900 tracking-tight">Dashboard Warga</h1>
            <p class="text-sm md:text-base text-gray-500 mt-1">{{ greeting }}, {{ authStore.user?.full_name || 'Warga' }}</p>
          </div>
          <div class="flex items-center gap-3">
            <div class="px-3 py-1 bg-green-50 text-green-700 rounded-full text-sm font-medium border border-green-100 flex items-center shadow-sm">
              <span class="w-2 h-2 bg-green-500 rounded-full mr-2 animate-pulse"></span>
              {{ authStore.user?.tenant_name || 'Komunitas' }}
            </div>
          </div>
        </div>

        <!-- Stats Cards -->
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 md:gap-6">
          <div class="card p-4 md:p-6">
            <div class="flex items-center justify-between">
              <div>
                <p class="text-xs md:text-sm font-medium text-gray-600 mb-1">Tagihan Saya</p>
                <p class="text-xl md:text-2xl font-bold text-gray-900">{{ stats.myBills || 0 }}</p>
                <p class="text-xs text-gray-500 mt-1">Tagihan aktif</p>
              </div>
              <div class="w-10 h-10 md:w-12 md:h-12 bg-blue-50 rounded-xl flex items-center justify-center">
                <Icon name="heroicons:document-text" class="w-5 h-5 md:w-6 md:h-6 text-blue-600" />
              </div>
            </div>
          </div>

          <div class="card p-4 md:p-6">
            <div class="flex items-center justify-between">
              <div>
                <p class="text-xs md:text-sm font-medium text-gray-600 mb-1">Pengumuman</p>
                <p class="text-xl md:text-2xl font-bold text-gray-900">{{ stats.announcements || 0 }}</p>
                <p class="text-xs text-gray-500 mt-1">Pengumuman baru</p>
              </div>
              <div class="w-10 h-10 md:w-12 md:h-12 bg-purple-50 rounded-xl flex items-center justify-center">
                <Icon name="heroicons:megaphone" class="w-5 h-5 md:w-6 md:h-6 text-purple-600" />
              </div>
            </div>
          </div>

          <div class="card p-4 md:p-6">
            <div class="flex items-center justify-between">
              <div>
                <p class="text-xs md:text-sm font-medium text-gray-600 mb-1">Status</p>
                <p class="text-base md:text-lg font-bold text-green-600">Aktif</p>
                <p class="text-xs text-gray-500 mt-1">Anggota komunitas</p>
              </div>
              <div class="w-10 h-10 md:w-12 md:h-12 bg-green-50 rounded-xl flex items-center justify-center">
                <Icon name="heroicons:check-circle" class="w-5 h-5 md:w-6 md:h-6 text-green-600" />
              </div>
            </div>
          </div>
        </div>

        <!-- Quick Actions -->
        <div class="card">
          <div class="p-6 border-b border-gray-200">
            <h2 class="text-lg font-semibold text-gray-900">Akses Cepat</h2>
            <p class="text-sm text-gray-600 mt-1">Layanan untuk warga</p>
          </div>
          <div class="p-6">
            <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
              <NuxtLink
                to="/dashboard/warga/bills"
                class="group p-4 border-2 border-gray-200 rounded-xl hover:border-primary-500 hover:bg-primary-50 transition-all duration-200 text-center"
              >
                <div class="w-10 h-10 bg-blue-100 rounded-xl flex items-center justify-center mb-3 mx-auto group-hover:bg-blue-200 transition-colors">
                  <Icon name="heroicons:document-text" class="w-5 h-5 text-blue-600" />
                </div>
                <h3 class="font-semibold text-gray-900 mb-1 text-sm">Tagihan</h3>
              </NuxtLink>

              <NuxtLink
                to="/dashboard/warga/announcements"
                class="group p-4 border-2 border-gray-200 rounded-xl hover:border-primary-500 hover:bg-primary-50 transition-all duration-200 text-center"
              >
                <div class="w-10 h-10 bg-purple-100 rounded-xl flex items-center justify-center mb-3 mx-auto group-hover:bg-purple-200 transition-colors">
                  <Icon name="heroicons:megaphone" class="w-5 h-5 text-purple-600" />
                </div>
                <h3 class="font-semibold text-gray-900 mb-1 text-sm">Pengumuman</h3>
              </NuxtLink>

              <NuxtLink
                to="/settings"
                class="group p-4 border-2 border-gray-200 rounded-xl hover:border-primary-500 hover:bg-primary-50 transition-all duration-200 text-center"
              >
                <div class="w-10 h-10 bg-gray-100 rounded-xl flex items-center justify-center mb-3 mx-auto group-hover:bg-gray-200 transition-colors">
                  <Icon name="heroicons:cog-6-tooth" class="w-5 h-5 text-gray-600" />
                </div>
                <h3 class="font-semibold text-gray-900 mb-1 text-sm">Pengaturan</h3>
              </NuxtLink>

              <button
                class="group p-4 border-2 border-red-200 rounded-xl hover:border-red-500 hover:bg-red-50 transition-all duration-200 text-center"
              >
                <div class="w-10 h-10 bg-red-100 rounded-xl flex items-center justify-center mb-3 mx-auto group-hover:bg-red-200 transition-colors">
                  <Icon name="heroicons:exclamation-triangle" class="w-5 h-5 text-red-600" />
                </div>
                <h3 class="font-semibold text-red-600 mb-1 text-sm">Panic</h3>
              </button>
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- Other Roles (Bendahara, Sekretariat, Satpam) -->
    <template v-else>
      <div class="text-center py-12">
        <h1 class="text-2xl font-bold text-gray-900 mb-2">Dashboard</h1>
        <p class="text-gray-600">Selamat datang, {{ authStore.user?.full_name || 'User' }}</p>
        <p class="text-sm text-gray-500 mt-2">Role: {{ authStore.user?.role_name || 'Unknown' }}</p>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useAuth } from '@/composables/useAuth'
import { useApi } from '@/composables/useApi'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

const authStore = useAuthStore()
const { isAdmin, isResident } = useAuth()
const { fetch } = useApi()
const router = useRouter()

const stats = ref({
  totalUnits: 0,
  totalUsers: 0,
  totalRoles: 0,
  myBills: 0,
  announcements: 0,
})

const greeting = computed(() => {
  const hour = new Date().getHours()
  if (hour < 12) return 'Selamat Pagi'
  if (hour < 17) return 'Selamat Siang'
  return 'Selamat Malam'
})

onMounted(async () => {
  // Redirect warga to their dashboard
  if (isResident.value) {
    router.replace('/dashboard/warga')
    return
  }

  try {
    if (isAdmin.value) {
      // Admin stats
      const [unitsRes, usersRes, rolesRes] = await Promise.all([
        fetch('/api/units?limit=1').catch(() => ({ pagination: { total: 0 } })),
        fetch('/api/users?limit=1').catch(() => ({ pagination: { total: 0 } })),
        fetch('/api/roles').catch(() => ({ roles: [] })),
      ])

      stats.value = {
        totalUnits: unitsRes.pagination?.total || 0,
        totalUsers: usersRes.pagination?.total || 0,
        totalRoles: Array.isArray(rolesRes.roles) ? rolesRes.roles.length : 0,
        myBills: 0,
        announcements: 0,
      }
    }
  } catch (error) {
    console.error('Failed to fetch stats:', error)
  }
})
</script>
