<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Sidebar -->
    <aside
      :class="[
        'fixed inset-y-0 left-0 z-50 w-72 bg-white border-r border-gray-200 transform transition-transform duration-300 ease-in-out',
        sidebarOpen ? 'translate-x-0' : '-translate-x-full',
        'lg:translate-x-0'
      ]"
    >
      <div class="flex flex-col h-full">
        <!-- Logo -->
        <div class="flex items-center h-20 px-6 border-b border-gray-200">
          <div class="flex items-center space-x-3 flex-1">
            <div class="w-10 h-10 bg-gradient-to-br from-primary-600 to-primary-700 rounded-xl flex items-center justify-center shadow-lg">
              <span class="text-white font-bold text-xl">R</span>
            </div>
            <div class="flex-1 min-w-0">
              <h1 class="text-lg font-bold text-gray-900 truncate">RukunOS</h1>
              <p class="text-xs text-gray-500 truncate" v-if="authStore.user?.tenant_name">
                {{ authStore.user.tenant_name }}
              </p>
            </div>
          </div>
          <button
            @click="sidebarOpen = false"
            class="lg:hidden p-2 rounded-lg text-gray-400 hover:text-gray-500 hover:bg-gray-100 transition-colors"
          >
            <Icon name="heroicons:x-mark" class="w-5 h-5" />
          </button>
        </div>

        <!-- Navigation -->
        <nav class="flex-1 px-4 py-6 space-y-1 overflow-y-auto">
          <NuxtLink
            v-for="item in navigation"
            :key="item.name"
            :to="item.to"
            class="flex items-center px-4 py-3 text-sm font-medium rounded-xl transition-all duration-200"
            :class="[
              $route.path === item.to
                ? 'bg-primary-50 text-primary-700 shadow-sm'
                : 'text-gray-700 hover:bg-gray-50'
            ]"
          >
            <Icon :name="item.icon" class="w-5 h-5 mr-3 flex-shrink-0" />
            <span>{{ item.name }}</span>
          </NuxtLink>
        </nav>

        <!-- User Menu -->
        <div class="p-4 border-t border-gray-200 bg-gray-50">
          <div class="flex items-center space-x-3 mb-3 p-3 rounded-xl bg-white">
            <div class="w-10 h-10 bg-gradient-to-br from-primary-500 to-primary-600 rounded-lg flex items-center justify-center flex-shrink-0">
              <span class="text-white font-semibold text-sm">
                {{ authStore.user?.full_name?.charAt(0).toUpperCase() }}
              </span>
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-sm font-semibold text-gray-900 truncate">
                {{ authStore.user?.full_name }}
              </p>
              <p class="text-xs text-gray-500 truncate">
                {{ authStore.user?.email }}
              </p>
            </div>
          </div>
          <button
            @click="authStore.logout()"
            class="w-full flex items-center px-4 py-2.5 text-sm font-medium text-gray-700 hover:bg-white rounded-xl transition-colors"
          >
            <Icon name="heroicons:arrow-right-on-rectangle" class="w-5 h-5 mr-3" />
            Keluar
          </button>
        </div>
      </div>
    </aside>

    <!-- Overlay for mobile -->
    <div
      v-if="sidebarOpen"
      @click="sidebarOpen = false"
      class="fixed inset-0 bg-black/50 z-40 lg:hidden backdrop-blur-sm"
    ></div>

    <!-- Main Content -->
    <div class="lg:pl-72">
      <!-- Top Bar -->
      <header class="sticky top-0 z-30 bg-white/80 backdrop-blur-md border-b border-gray-200">
        <div class="flex items-center justify-between h-16 px-4 sm:px-6 lg:px-8">
          <button
            @click="sidebarOpen = true"
            class="lg:hidden p-2 rounded-lg text-gray-400 hover:text-gray-500 hover:bg-gray-100 transition-colors"
          >
            <Icon name="heroicons:bars-3" class="w-6 h-6" />
          </button>

          <div class="flex-1 flex items-center justify-end space-x-4">
            <!-- Notifications -->
            <button class="relative p-2 rounded-lg text-gray-400 hover:text-gray-500 hover:bg-gray-100 transition-colors">
              <Icon name="heroicons:bell" class="w-5 h-5" />
              <span class="absolute top-1 right-1 w-2 h-2 bg-red-500 rounded-full ring-2 ring-white"></span>
            </button>
          </div>
        </div>
      </header>

      <!-- Page Content -->
      <main class="p-4 sm:p-6 lg:p-8">
        <slot />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
const authStore = useAuthStore()
const sidebarOpen = ref(false)

const navigation = [
  { name: 'Dashboard', to: '/', icon: 'heroicons:home' },
  { name: 'Unit', to: '/units', icon: 'heroicons:building-office' },
  { name: 'Pengguna', to: '/users', icon: 'heroicons:users' },
  { name: 'Role & Izin', to: '/roles', icon: 'heroicons:shield-check' },
  { name: 'Pengaturan', to: '/settings', icon: 'heroicons:cog-6-tooth' },
]

// Close sidebar on route change (mobile)
watch(() => useRoute().path, () => {
  sidebarOpen.value = false
})
</script>
