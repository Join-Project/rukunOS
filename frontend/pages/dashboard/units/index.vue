<template>
  <div class="max-w-7xl mx-auto">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 tracking-tight">Manajemen Unit & Rumah</h1>
        <p class="text-gray-500">Daftar unit, status hunian, dan data pemilik.</p>
      </div>
      <button class="px-4 py-2 bg-primary-600 text-white rounded-xl font-medium hover:bg-primary-700 transition-colors shadow-lg shadow-primary-600/20 flex items-center">
        <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path></svg>
        Tambah Unit
      </button>
    </div>

    <!-- Stats -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
      <div class="bg-white p-4 rounded-xl border border-gray-100 shadow-sm">
        <p class="text-sm text-gray-500">Total Unit</p>
        <p class="text-2xl font-bold text-gray-900">150</p>
      </div>
      <div class="bg-white p-4 rounded-xl border border-gray-100 shadow-sm">
        <p class="text-sm text-gray-500">Terisi (Occupied)</p>
        <p class="text-2xl font-bold text-green-600">142</p>
      </div>
      <div class="bg-white p-4 rounded-xl border border-gray-100 shadow-sm">
        <p class="text-sm text-gray-500">Kosong (Vacant)</p>
        <p class="text-2xl font-bold text-gray-400">8</p>
      </div>
      <div class="bg-white p-4 rounded-xl border border-gray-100 shadow-sm">
        <p class="text-sm text-gray-500">Ruko/Usaha</p>
        <p class="text-2xl font-bold text-blue-600">12</p>
      </div>
    </div>

    <!-- Units Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div v-for="unit in units" :key="unit.id" class="bg-white rounded-2xl border border-gray-100 shadow-sm p-6 hover:shadow-md transition-shadow">
        <div class="flex justify-between items-start mb-4">
          <div class="flex items-center">
            <div class="w-12 h-12 bg-gray-100 rounded-xl flex items-center justify-center text-gray-700 font-bold text-lg">
              {{ unit.code }}
            </div>
            <div class="ml-4">
              <h3 class="font-bold text-gray-900">{{ unit.address }}</h3>
              <p class="text-sm text-gray-500">{{ unit.type }}</p>
            </div>
          </div>
          <span 
            :class="[
              'px-2.5 py-0.5 rounded-full text-xs font-medium border',
              unit.status === 'Occupied' ? 'bg-green-50 text-green-700 border-green-100' : 'bg-gray-50 text-gray-500 border-gray-100'
            ]"
          >
            {{ unit.status === 'Occupied' ? 'Dihuni' : 'Kosong' }}
          </span>
        </div>
        
        <div class="space-y-3">
          <div class="flex items-center text-sm">
            <svg class="w-4 h-4 text-gray-400 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path></svg>
            <span class="text-gray-600">{{ unit.owner || '-' }}</span>
          </div>
          <div class="flex items-center text-sm">
            <svg class="w-4 h-4 text-gray-400 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"></path></svg>
            <span class="text-gray-600">{{ unit.phone || '-' }}</span>
          </div>
        </div>

        <div class="mt-6 pt-4 border-t border-gray-50 flex justify-end gap-2">
          <button class="text-sm text-gray-500 hover:text-gray-900 font-medium px-3 py-1.5 rounded-lg hover:bg-gray-50">Edit</button>
          <button class="text-sm text-primary-600 hover:text-primary-700 font-medium px-3 py-1.5 rounded-lg hover:bg-primary-50">Detail</button>
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
// Mock Data
const units = [
  { id: 1, code: 'A-01', address: 'Jl. Merpati No. 1', type: 'Rumah', status: 'Occupied', owner: 'Budi Santoso', phone: '08123456789' },
  { id: 2, code: 'A-02', address: 'Jl. Merpati No. 2', type: 'Rumah', status: 'Occupied', owner: 'Siti Aminah', phone: '08123456788' },
  { id: 3, code: 'A-03', address: 'Jl. Merpati No. 3', type: 'Rumah', status: 'Vacant', owner: '-', phone: '-' },
  { id: 4, code: 'B-01', address: 'Jl. Elang No. 1', type: 'Ruko', status: 'Occupied', owner: 'Toko Maju Jaya', phone: '08123456777' },
  { id: 5, code: 'B-02', address: 'Jl. Elang No. 2', type: 'Rumah', status: 'Occupied', owner: 'Ahmad Dhani', phone: '08123456766' },
  { id: 6, code: 'C-05', address: 'Jl. Kenari No. 5', type: 'Rumah', status: 'Occupied', owner: 'Rina Nose', phone: '08123456755' },
]
</script>
