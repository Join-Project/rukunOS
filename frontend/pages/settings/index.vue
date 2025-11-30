<template>
  <div class="max-w-4xl mx-auto space-y-8">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 tracking-tight">Pengaturan</h1>
        <p class="text-gray-500">Kelola profil dan preferensi aplikasi Anda</p>
      </div>

      <!-- Profile Section -->
      <div class="bg-white rounded-2xl border border-gray-100 p-8 shadow-sm">
        <h3 class="text-lg font-bold text-gray-900 mb-6 flex items-center gap-2">
          <svg class="w-5 h-5 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path></svg>
          Profil Saya
        </h3>
        
        <div class="flex flex-col md:flex-row gap-8 items-start">
          <!-- Avatar -->
          <div class="flex flex-col items-center gap-4">
            <div class="w-24 h-24 bg-primary-100 rounded-full flex items-center justify-center text-3xl font-bold text-primary-600 border-4 border-white shadow-md">
              BS
            </div>
            <button class="text-sm font-medium text-primary-600 hover:text-primary-700 hover:underline">Ubah Foto</button>
          </div>

          <!-- Form -->
          <div class="flex-1 w-full space-y-4">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Nama Lengkap</label>
                <input type="text" value="Budi Santoso" class="w-full rounded-xl border-gray-300 focus:border-primary-500 focus:ring-primary-500">
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Email</label>
                <input type="email" value="budi@rukunos.com" class="w-full rounded-xl border-gray-300 focus:border-primary-500 focus:ring-primary-500">
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Nomor Telepon</label>
                <input type="tel" value="0812-3456-7890" class="w-full rounded-xl border-gray-300 focus:border-primary-500 focus:ring-primary-500">
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Unit / Rumah</label>
                <input type="text" value="A-02" disabled class="w-full rounded-xl border-gray-200 bg-gray-50 text-gray-500 cursor-not-allowed">
              </div>
            </div>
            
            <div class="pt-4 border-t border-gray-100 mt-4">
              <h4 class="text-sm font-bold text-gray-900 mb-3">Ganti Password</h4>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <input type="password" placeholder="Password Baru" class="w-full rounded-xl border-gray-300 focus:border-primary-500 focus:ring-primary-500">
                <input type="password" placeholder="Konfirmasi Password" class="w-full rounded-xl border-gray-300 focus:border-primary-500 focus:ring-primary-500">
              </div>
            </div>

            <div class="flex justify-end pt-4">
              <button class="bg-primary-600 hover:bg-primary-700 text-white px-6 py-2.5 rounded-xl font-bold shadow-lg shadow-primary-600/20 transition-all">
                Simpan Perubahan
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- App Preferences -->
      <div class="bg-white rounded-2xl border border-gray-100 p-8 shadow-sm">
        <h3 class="text-lg font-bold text-gray-900 mb-6 flex items-center gap-2">
          <svg class="w-5 h-5 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path></svg>
          Preferensi Aplikasi
        </h3>

        <div class="space-y-4">
          <div class="flex items-center justify-between py-3 border-b border-gray-50">
            <div>
              <h4 class="font-medium text-gray-900">Notifikasi WhatsApp</h4>
              <p class="text-sm text-gray-500">Terima info penting via WhatsApp</p>
            </div>
            <label class="relative inline-flex items-center cursor-pointer">
              <input type="checkbox" class="sr-only peer" checked>
              <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary-600"></div>
            </label>
          </div>

          <div class="flex items-center justify-between py-3 border-b border-gray-50">
            <div>
              <h4 class="font-medium text-gray-900">Notifikasi Email</h4>
              <p class="text-sm text-gray-500">Terima tagihan bulanan via email</p>
            </div>
            <label class="relative inline-flex items-center cursor-pointer">
              <input type="checkbox" class="sr-only peer" checked>
              <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary-600"></div>
            </label>
          </div>

          <div class="flex items-center justify-between py-3">
            <div>
              <h4 class="font-medium text-gray-900">Bahasa</h4>
              <p class="text-sm text-gray-500">Pilih bahasa tampilan aplikasi</p>
            </div>
            <select class="rounded-lg border-gray-200 text-sm focus:border-primary-500 focus:ring-primary-500">
              <option>Bahasa Indonesia</option>
              <option>English</option>
            </select>
          </div>
        </div>
      </div>
    </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})
</script>
