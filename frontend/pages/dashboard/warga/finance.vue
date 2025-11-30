<template>
  <div class="max-w-6xl mx-auto space-y-8">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-bold text-gray-900 tracking-tight">Laporan Keuangan RT</h1>
          <p class="text-gray-500">Transparansi dana kas warga RT 05 / RW 03</p>
        </div>
        <div class="flex items-center gap-2 px-4 py-2 bg-green-50 text-green-700 rounded-xl border border-green-100">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
          <span class="text-sm font-bold">Data Terverifikasi</span>
        </div>
      </div>

      <!-- Financial Summary Cards -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div class="bg-white p-6 rounded-2xl border border-gray-100 shadow-sm">
          <p class="text-sm text-gray-500 font-medium mb-1">Total Kas Saat Ini</p>
          <h3 class="text-3xl font-bold text-gray-900">Rp 45.250.000</h3>
          <div class="mt-4 flex items-center text-sm text-green-600">
            <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"></path></svg>
            <span>+12% dari bulan lalu</span>
          </div>
        </div>
        <div class="bg-white p-6 rounded-2xl border border-gray-100 shadow-sm">
          <p class="text-sm text-gray-500 font-medium mb-1">Pemasukan (Bulan Ini)</p>
          <h3 class="text-3xl font-bold text-green-600">Rp 8.500.000</h3>
          <p class="text-xs text-gray-400 mt-2">Iuran Warga & Donasi</p>
        </div>
        <div class="bg-white p-6 rounded-2xl border border-gray-100 shadow-sm">
          <p class="text-sm text-gray-500 font-medium mb-1">Pengeluaran (Bulan Ini)</p>
          <h3 class="text-3xl font-bold text-red-600">Rp 3.200.000</h3>
          <p class="text-xs text-gray-400 mt-2">Kebersihan & Keamanan</p>
        </div>
      </div>

      <!-- Charts Section (Mock) -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div class="bg-white p-6 rounded-2xl border border-gray-100 shadow-sm">
          <h3 class="text-lg font-bold text-gray-900 mb-6">Arus Kas 6 Bulan Terakhir</h3>
          <div class="h-64 flex items-end justify-between gap-2 px-2">
            <!-- Mock Bar Chart -->
            <div class="w-full bg-gray-50 rounded-t-lg relative group h-[60%]">
              <div class="absolute bottom-0 w-full bg-blue-500 rounded-t-lg transition-all duration-500 group-hover:bg-blue-600 h-full"></div>
              <span class="absolute -bottom-6 left-1/2 -translate-x-1/2 text-xs text-gray-500">Sep</span>
            </div>
            <div class="w-full bg-gray-50 rounded-t-lg relative group h-[70%]">
              <div class="absolute bottom-0 w-full bg-blue-500 rounded-t-lg transition-all duration-500 group-hover:bg-blue-600 h-full"></div>
              <span class="absolute -bottom-6 left-1/2 -translate-x-1/2 text-xs text-gray-500">Okt</span>
            </div>
            <div class="w-full bg-gray-50 rounded-t-lg relative group h-[65%]">
              <div class="absolute bottom-0 w-full bg-blue-500 rounded-t-lg transition-all duration-500 group-hover:bg-blue-600 h-full"></div>
              <span class="absolute -bottom-6 left-1/2 -translate-x-1/2 text-xs text-gray-500">Nov</span>
            </div>
            <div class="w-full bg-gray-50 rounded-t-lg relative group h-[80%]">
              <div class="absolute bottom-0 w-full bg-blue-500 rounded-t-lg transition-all duration-500 group-hover:bg-blue-600 h-full"></div>
              <span class="absolute -bottom-6 left-1/2 -translate-x-1/2 text-xs text-gray-500">Des</span>
            </div>
            <div class="w-full bg-gray-50 rounded-t-lg relative group h-[75%]">
              <div class="absolute bottom-0 w-full bg-blue-500 rounded-t-lg transition-all duration-500 group-hover:bg-blue-600 h-full"></div>
              <span class="absolute -bottom-6 left-1/2 -translate-x-1/2 text-xs text-gray-500">Jan</span>
            </div>
            <div class="w-full bg-gray-50 rounded-t-lg relative group h-[85%]">
              <div class="absolute bottom-0 w-full bg-blue-500 rounded-t-lg transition-all duration-500 group-hover:bg-blue-600 h-full"></div>
              <span class="absolute -bottom-6 left-1/2 -translate-x-1/2 text-xs text-gray-500">Feb</span>
            </div>
          </div>
        </div>

        <div class="bg-white p-6 rounded-2xl border border-gray-100 shadow-sm">
          <h3 class="text-lg font-bold text-gray-900 mb-6">Alokasi Pengeluaran</h3>
          <div class="flex items-center justify-center h-64">
            <div class="relative w-48 h-48 rounded-full border-[16px] border-gray-100 flex items-center justify-center">
              <!-- Simple CSS Pie Chart Representation -->
              <div class="absolute inset-0 rounded-full border-[16px] border-primary-500 border-r-transparent border-b-transparent rotate-45"></div>
              <div class="absolute inset-0 rounded-full border-[16px] border-yellow-400 border-l-transparent border-t-transparent border-b-transparent rotate-[135deg]"></div>
              <div class="text-center">
                <span class="block text-2xl font-bold text-gray-900">Feb</span>
                <span class="text-xs text-gray-500">2025</span>
              </div>
            </div>
          </div>
          <div class="flex justify-center gap-4 mt-4">
            <div class="flex items-center gap-2 text-sm text-gray-600">
              <span class="w-3 h-3 bg-primary-500 rounded-full"></span> Keamanan
            </div>
            <div class="flex items-center gap-2 text-sm text-gray-600">
              <span class="w-3 h-3 bg-yellow-400 rounded-full"></span> Kebersihan
            </div>
            <div class="flex items-center gap-2 text-sm text-gray-600">
              <span class="w-3 h-3 bg-gray-200 rounded-full"></span> Lainnya
            </div>
          </div>
        </div>
      </div>

      <!-- Recent Transactions -->
      <div class="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden">
        <div class="p-6 border-b border-gray-100">
          <h3 class="text-lg font-bold text-gray-900">Transaksi Terakhir</h3>
        </div>
        <div class="overflow-x-auto">
          <table class="w-full text-left text-sm">
            <thead class="bg-gray-50 text-gray-500 font-medium">
              <tr>
                <th class="px-6 py-4">Tanggal</th>
                <th class="px-6 py-4">Keterangan</th>
                <th class="px-6 py-4">Kategori</th>
                <th class="px-6 py-4 text-right">Jumlah</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
              <tr class="hover:bg-gray-50 transition-colors">
                <td class="px-6 py-4 text-gray-500">12 Feb 2025</td>
                <td class="px-6 py-4 font-medium text-gray-900">Pembayaran Gaji Satpam (3 Orang)</td>
                <td class="px-6 py-4"><span class="px-2 py-1 bg-red-50 text-red-700 rounded-full text-xs font-bold">Pengeluaran</span></td>
                <td class="px-6 py-4 text-right font-bold text-red-600">- Rp 4.500.000</td>
              </tr>
              <tr class="hover:bg-gray-50 transition-colors">
                <td class="px-6 py-4 text-gray-500">10 Feb 2025</td>
                <td class="px-6 py-4 font-medium text-gray-900">Iuran Warga (50 Unit)</td>
                <td class="px-6 py-4"><span class="px-2 py-1 bg-green-50 text-green-700 rounded-full text-xs font-bold">Pemasukan</span></td>
                <td class="px-6 py-4 text-right font-bold text-green-600">+ Rp 5.000.000</td>
              </tr>
              <tr class="hover:bg-gray-50 transition-colors">
                <td class="px-6 py-4 text-gray-500">05 Feb 2025</td>
                <td class="px-6 py-4 font-medium text-gray-900">Perbaikan Lampu Jalan Blok A</td>
                <td class="px-6 py-4"><span class="px-2 py-1 bg-red-50 text-red-700 rounded-full text-xs font-bold">Pengeluaran</span></td>
                <td class="px-6 py-4 text-right font-bold text-red-600">- Rp 750.000</td>
              </tr>
            </tbody>
          </table>
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
