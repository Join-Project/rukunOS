<template>
  <div class="min-h-screen bg-gray-50 flex font-sans">
    <!-- Sidebar -->
    <aside class="w-64 bg-white border-r border-gray-100 hidden md:flex flex-col fixed h-full z-20">
      <div class="h-16 flex items-center px-6 border-b border-gray-100">
        <div class="w-8 h-8 bg-primary-600 rounded-lg flex items-center justify-center text-white font-bold mr-3">R</div>
        <span class="text-xl font-bold text-gray-900 tracking-tight">RukunOS</span>
      </div>

      <div class="flex-1 overflow-y-auto py-6 px-4 space-y-8">
        <!-- DEFAULT MENU (Fallback jika role tidak terdeteksi) -->
        <template v-if="!isAdmin && !isBendahara && !isSecretariat && !isSecurity && !isResident">
          <div>
            <h3 class="px-2 text-xs font-semibold text-gray-400 uppercase tracking-wider mb-3">Menu</h3>
            <nav class="space-y-1">
              <NuxtLink to="/dashboard" class="flex items-center px-2 py-2 text-sm font-medium rounded-lg transition-colors" :class="$route.path === '/dashboard' ? 'bg-primary-50 text-primary-700' : 'text-gray-600 hover:bg-gray-50 hover:text-primary-600'">
                <svg class="mr-3 h-5 w-5" :class="$route.path === '/dashboard' ? 'text-primary-500' : 'text-gray-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"></path></svg>
                Dashboard
              </NuxtLink>
            </nav>
          </div>
        </template>

        <!-- ADMIN MENU -->
        <template v-if="isAdmin">
          <!-- Menu Group: Main -->
          <div>
            <h3 class="px-2 text-xs font-semibold text-gray-400 uppercase tracking-wider mb-3">Main</h3>
            <nav class="space-y-1">
              <NuxtLink to="/dashboard" class="flex items-center px-2 py-2 text-sm font-medium rounded-lg transition-colors" :class="$route.path === '/dashboard' ? 'bg-primary-50 text-primary-700' : 'text-gray-600 hover:bg-gray-50 hover:text-primary-600'">
                <svg class="mr-3 h-5 w-5" :class="$route.path === '/dashboard' ? 'text-primary-500' : 'text-gray-400 group-hover:text-primary-500'" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z"></path></svg>
                Dashboard
              </NuxtLink>
            </nav>
          </div>

          <!-- Menu Group: Administrasi -->
          <div>
            <h3 class="px-2 text-xs font-semibold text-gray-400 uppercase tracking-wider mb-3">Administrasi</h3>
            <nav class="space-y-1">
              <NuxtLink to="/units" class="flex items-center px-2 py-2 text-sm font-medium rounded-lg transition-colors" :class="$route.path.startsWith('/units') ? 'bg-primary-50 text-primary-700' : 'text-gray-600 hover:bg-gray-50 hover:text-primary-600'">
                <svg class="mr-3 h-5 w-5" :class="$route.path.startsWith('/units') ? 'text-primary-500' : 'text-gray-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"></path></svg>
                Unit / Rumah
              </NuxtLink>
              <NuxtLink to="/users" class="flex items-center px-2 py-2 text-sm font-medium rounded-lg transition-colors" :class="$route.path.startsWith('/users') ? 'bg-primary-50 text-primary-700' : 'text-gray-600 hover:bg-gray-50 hover:text-primary-600'">
                <svg class="mr-3 h-5 w-5" :class="$route.path.startsWith('/users') ? 'text-primary-500' : 'text-gray-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"></path></svg>
                Warga
              </NuxtLink>
              <NuxtLink to="/roles" class="flex items-center px-2 py-2 text-sm font-medium rounded-lg transition-colors" :class="$route.path.startsWith('/roles') ? 'bg-primary-50 text-primary-700' : 'text-gray-600 hover:bg-gray-50 hover:text-primary-600'">
                <svg class="mr-3 h-5 w-5" :class="$route.path.startsWith('/roles') ? 'text-primary-500' : 'text-gray-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"></path></svg>
                Role & Izin
              </NuxtLink>
            </nav>
          </div>
        </template>

        <!-- KEUANGAN MENU (Admin OR Bendahara) -->
        <template v-if="isAdmin || isBendahara">
          <div>
            <h3 class="px-2 text-xs font-semibold text-gray-400 uppercase tracking-wider mb-3">Keuangan</h3>
            <nav class="space-y-1">
              <NuxtLink to="/dashboard/billing/reports" class="flex items-center px-2 py-2 text-sm font-medium rounded-lg transition-colors" :class="$route.path === '/dashboard/billing/reports' ? 'bg-primary-50 text-primary-700' : 'text-gray-600 hover:bg-gray-50 hover:text-primary-600'">
                <svg class="mr-3 h-5 w-5" :class="$route.path === '/dashboard/billing/reports' ? 'text-primary-500' : 'text-gray-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path></svg>
                Dashboard Keuangan
              </NuxtLink>
              <NuxtLink to="/dashboard/billing" class="flex items-center px-2 py-2 text-sm font-medium rounded-lg transition-colors" :class="$route.path === '/dashboard/billing' || ($route.path.startsWith('/dashboard/billing') && !$route.path.includes('/create') && !$route.path.includes('/reports')) ? 'bg-primary-50 text-primary-700' : 'text-gray-600 hover:bg-gray-50 hover:text-primary-600'">
                <svg class="mr-3 h-5 w-5" :class="$route.path === '/dashboard/billing' || ($route.path.startsWith('/dashboard/billing') && !$route.path.includes('/create') && !$route.path.includes('/reports')) ? 'text-primary-500' : 'text-gray-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
                Manajemen Tagihan
              </NuxtLink>
              <NuxtLink to="/dashboard/billing/create" class="flex items-center px-2 py-2 text-sm font-medium rounded-lg transition-colors" :class="$route.path.includes('/dashboard/billing/create') ? 'bg-primary-50 text-primary-700' : 'text-gray-600 hover:bg-gray-50 hover:text-primary-600'">
                <svg class="mr-3 h-5 w-5" :class="$route.path.includes('/dashboard/billing/create') ? 'text-primary-500' : 'text-gray-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path></svg>
                Buat Tagihan
              </NuxtLink>
            </nav>
          </div>
        </template>

        <!-- KOMUNIKASI MENU (Admin OR Sekretariat) -->
        <template v-if="isAdmin || isSecretariat">
          <div>
            <h3 class="px-2 text-xs font-semibold text-gray-400 uppercase tracking-wider mb-3">Komunikasi</h3>
            <nav class="space-y-1">
              <NuxtLink to="/dashboard/announcements" class="flex items-center px-2 py-2 text-sm font-medium rounded-lg transition-colors" :class="$route.path.startsWith('/dashboard/announcements') ? 'bg-primary-50 text-primary-700' : 'text-gray-600 hover:bg-gray-50 hover:text-primary-600'">
                <svg class="mr-3 h-5 w-5" :class="$route.path.startsWith('/dashboard/announcements') ? 'text-primary-500' : 'text-gray-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5.882V19.24a1.76 1.76 0 01-3.417.592l-2.147-6.15M18 13a3 3 0 100-6M5.436 13.683A4.001 4.001 0 017 6h1.832c4.1 0 7.625-1.234 9.168-3v14c-1.543-1.766-5.067-3-9.168-3H7a3.988 3.988 0 01-1.564-.317z"></path></svg>
                Pengumuman
              </NuxtLink>
            </nav>
          </div>
        </template>

        <!-- KEAMANAN MENU (Admin OR Satpam) -->
        <template v-if="isAdmin || isSecurity">
          <div>
            <h3 class="px-2 text-xs font-semibold text-gray-400 uppercase tracking-wider mb-3">Keamanan</h3>
            <nav class="space-y-1">
              <NuxtLink to="/dashboard/security/visitors" class="flex items-center px-2 py-2 text-sm font-medium rounded-lg transition-colors" :class="$route.path.startsWith('/dashboard/security/visitors') ? 'bg-primary-50 text-primary-700' : 'text-gray-600 hover:bg-gray-50 hover:text-primary-600'">
                <svg class="mr-3 h-5 w-5" :class="$route.path.startsWith('/dashboard/security/visitors') ? 'text-primary-500' : 'text-gray-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"></path></svg>
                Buku Tamu (VMS)
              </NuxtLink>
              <NuxtLink to="/dashboard/security/alerts" class="flex items-center px-2 py-2 text-sm font-medium rounded-lg transition-colors" :class="$route.path.startsWith('/dashboard/security/alerts') ? 'bg-primary-50 text-primary-700' : 'text-gray-600 hover:bg-gray-50 hover:text-primary-600'">
                <svg class="mr-3 h-5 w-5" :class="$route.path.startsWith('/dashboard/security/alerts') ? 'text-primary-500' : 'text-gray-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path></svg>
                Panic Alerts
              </NuxtLink>
            </nav>
          </div>
        </template>

        <!-- WARGA MENU (Resident Only) -->
        <template v-if="isResident">
          <div>
            <h3 class="px-2 text-xs font-semibold text-gray-400 uppercase tracking-wider mb-3">Main</h3>
            <nav class="space-y-1">
              <NuxtLink to="/dashboard" class="flex items-center px-2 py-2 text-sm font-medium rounded-lg transition-colors" :class="$route.path === '/dashboard' ? 'bg-primary-50 text-primary-700' : 'text-gray-600 hover:bg-gray-50 hover:text-primary-600'">
                <svg class="mr-3 h-5 w-5" :class="$route.path === '/dashboard' ? 'text-primary-500' : 'text-gray-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"></path></svg>
                Dashboard
              </NuxtLink>
            </nav>
          </div>

          <div>
            <h3 class="px-2 text-xs font-semibold text-gray-400 uppercase tracking-wider mb-3">Layanan</h3>
            <nav class="space-y-1">
              <NuxtLink to="/dashboard/warga/bills" class="flex items-center px-2 py-2 text-sm font-medium rounded-lg transition-colors" :class="$route.path.startsWith('/dashboard/warga/bills') ? 'bg-primary-50 text-primary-700' : 'text-gray-600 hover:bg-gray-50 hover:text-primary-600'">
                <svg class="mr-3 h-5 w-5" :class="$route.path.startsWith('/dashboard/warga/bills') ? 'text-primary-500' : 'text-gray-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01"></path></svg>
                Tagihan Saya
              </NuxtLink>
              <NuxtLink to="/dashboard/warga/announcements" class="flex items-center px-2 py-2 text-sm font-medium rounded-lg transition-colors" :class="$route.path.startsWith('/dashboard/warga/announcements') ? 'bg-primary-50 text-primary-700' : 'text-gray-600 hover:bg-gray-50 hover:text-primary-600'">
                <svg class="mr-3 h-5 w-5" :class="$route.path.startsWith('/dashboard/warga/announcements') ? 'text-primary-500' : 'text-gray-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5.882V19.24a1.76 1.76 0 01-3.417.592l-2.147-6.15M18 13a3 3 0 100-6M5.436 13.683A4.001 4.001 0 017 6h1.832c4.1 0 7.625-1.234 9.168-3v14c-1.543-1.766-5.067-3-9.168-3H7a3.988 3.988 0 01-1.564-.317z"></path></svg>
                Pengumuman
              </NuxtLink>
              <NuxtLink to="/dashboard/warga/complaints" class="flex items-center px-2 py-2 text-sm font-medium rounded-lg transition-colors" :class="$route.path.startsWith('/dashboard/warga/complaints') ? 'bg-primary-50 text-primary-700' : 'text-gray-600 hover:bg-gray-50 hover:text-primary-600'">
                <svg class="mr-3 h-5 w-5" :class="$route.path.startsWith('/dashboard/warga/complaints') ? 'text-primary-500' : 'text-gray-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"></path></svg>
                Lapor Masalah
              </NuxtLink>
              <NuxtLink to="/dashboard/warga/documents" class="flex items-center px-2 py-2 text-sm font-medium rounded-lg transition-colors" :class="$route.path.startsWith('/dashboard/warga/documents') ? 'bg-primary-50 text-primary-700' : 'text-gray-600 hover:bg-gray-50 hover:text-primary-600'">
                <svg class="mr-3 h-5 w-5" :class="$route.path.startsWith('/dashboard/warga/documents') ? 'text-primary-500' : 'text-gray-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z"></path></svg>
                Surat Pengantar
              </NuxtLink>
              <NuxtLink to="/dashboard/warga/family" class="flex items-center px-2 py-2 text-sm font-medium rounded-lg transition-colors" :class="$route.path.startsWith('/dashboard/warga/family') ? 'bg-primary-50 text-primary-700' : 'text-gray-600 hover:bg-gray-50 hover:text-primary-600'">
                <svg class="mr-3 h-5 w-5" :class="$route.path.startsWith('/dashboard/warga/family') ? 'text-primary-500' : 'text-gray-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z"></path></svg>
                Data Keluarga
              </NuxtLink>
            </nav>
          </div>

          <div>
            <h3 class="px-2 text-xs font-semibold text-gray-400 uppercase tracking-wider mb-3">Darurat</h3>
            <nav class="space-y-1">
              <NuxtLink to="/dashboard/warga/panic" class="flex items-center px-2 py-2 text-sm font-medium rounded-lg transition-colors" :class="$route.path.startsWith('/dashboard/warga/panic') ? 'bg-red-50 text-red-700' : 'text-gray-600 hover:bg-red-50 hover:text-red-600'">
                <svg class="mr-3 h-5 w-5" :class="$route.path.startsWith('/dashboard/warga/panic') ? 'text-red-500' : 'text-gray-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path></svg>
                Panic Alert
              </NuxtLink>
            </nav>
          </div>
        </template>
      </div>

      <div class="p-4 border-t border-gray-100 space-y-2">
        <NuxtLink to="/settings" class="flex items-center px-2 py-2 text-sm font-medium text-gray-600 rounded-lg hover:bg-gray-50 hover:text-primary-600 group transition-colors">
          <svg class="mr-3 h-5 w-5 text-gray-400 group-hover:text-primary-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path></svg>
          Pengaturan
        </NuxtLink>
        
        <!-- Logout Button -->
        <button 
          @click="logout"
          class="w-full flex items-center px-2 py-2 text-xs font-medium text-red-600 bg-red-50 rounded-lg hover:bg-red-100 transition-colors"
        >
          <svg class="mr-3 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"></path></svg>
          Keluar ({{ user?.full_name || 'User' }})
        </button>
      </div>
    </aside>

    <!-- Mobile Header -->
    <div class="md:hidden fixed w-full bg-white border-b border-gray-100 z-20 px-4 h-16 flex items-center justify-between">
      <div class="flex items-center gap-2">
        <div class="w-8 h-8 bg-primary-600 rounded-lg flex items-center justify-center text-white font-bold">R</div>
        <span class="text-xl font-bold text-gray-900 tracking-tight">RukunOS</span>
      </div>
      <button class="text-gray-500 hover:text-gray-900">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path></svg>
      </button>
    </div>

    <!-- Main Content -->
    <main class="flex-1 md:ml-64 p-4 md:p-8 pt-20 md:pt-8">
      <slot />
    </main>
  </div>
</template>

<script setup lang="ts">
const { isAdmin, isResident, isSecurity, isBendahara, isSecretariat, logout, user } = useAuth()
const authStore = useAuthStore()

// Debug: Log user info
onMounted(() => {
  console.log('=== Dashboard Layout Debug ===')
  console.log('User:', user.value)
  console.log('Role:', user.value?.role_name)
  console.log('Permissions:', user.value?.permissions)
  console.log('isAdmin:', isAdmin.value)
  console.log('isBendahara:', isBendahara.value)
  console.log('isSecretariat:', isSecretariat.value)
  console.log('isSecurity:', isSecurity.value)
  console.log('isResident:', isResident.value)
  console.log('=============================')
})
</script>
