<template>
  <div class="max-w-7xl mx-auto">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-6 md:mb-8">
      <div>
        <h1 class="text-xl md:text-2xl font-bold text-gray-900 tracking-tight">Manajemen Tagihan</h1>
        <p class="text-sm md:text-base text-gray-500 mt-1">Kelola tagihan warga, pantau status pembayaran, dan kirim pengingat.</p>
      </div>
      <div class="flex flex-col sm:flex-row gap-2 sm:gap-3">
        <NuxtLink to="/dashboard/billing/reports" class="px-4 py-2 bg-white border border-gray-200 text-gray-700 rounded-xl font-medium hover:bg-gray-50 transition-colors shadow-sm text-sm md:text-base text-center">
          Laporan Keuangan
        </NuxtLink>
        <NuxtLink to="/dashboard/billing/create" class="px-4 py-2 bg-primary-600 text-white rounded-xl font-medium hover:bg-primary-700 transition-colors shadow-lg shadow-primary-600/20 flex items-center justify-center text-sm md:text-base">
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path></svg>
          Buat Tagihan
        </NuxtLink>
      </div>
    </div>

    <!-- Filters & Stats -->
    <div class="bg-white rounded-2xl border border-gray-100 shadow-sm p-1 mb-6 flex flex-col md:flex-row">
      <button 
        v-for="tab in tabs" 
        :key="tab.id"
        @click="activeTab = tab.id"
        :class="[
          'flex-1 py-3 px-4 rounded-xl text-sm font-medium transition-all flex items-center justify-center gap-2',
          activeTab === tab.id 
            ? 'bg-gray-50 text-gray-900 shadow-sm border border-gray-200' 
            : 'text-gray-500 hover:text-gray-700 hover:bg-gray-50'
        ]"
      >
        {{ tab.label }}
        <span 
          :class="[
            'px-2 py-0.5 rounded-full text-xs',
            activeTab === tab.id ? 'bg-gray-200 text-gray-800' : 'bg-gray-100 text-gray-500'
          ]"
        >
          {{ tab.count }}
        </span>
      </button>
    </div>

    <!-- Search & Filter Bar -->
    <div class="flex flex-col md:flex-row gap-4 mb-6">
      <div class="flex-1 relative">
        <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
          <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path></svg>
        </div>
        <input 
          v-model="searchQuery"
          @keyup.enter="loadBills"
          type="text" 
          class="block w-full pl-10 pr-3 py-2.5 border border-gray-200 rounded-xl leading-5 bg-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent sm:text-sm transition-shadow" 
          placeholder="Cari nama warga, unit, atau kategori tagihan..."
        >
      </div>
      <div class="w-full md:w-48">
        <select class="block w-full pl-3 pr-10 py-2.5 text-base border border-gray-200 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent sm:text-sm rounded-xl bg-white">
          <option>Januari 2025</option>
          <option>Desember 2024</option>
          <option>November 2024</option>
        </select>
      </div>
      <div class="w-full md:w-48">
        <select class="block w-full pl-3 pr-10 py-2.5 text-base border border-gray-200 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent sm:text-sm rounded-xl bg-white">
          <option value="">Semua Unit</option>
          <option value="A">Blok A</option>
          <option value="B">Blok B</option>
        </select>
      </div>
    </div>

    <!-- Bills Table - Desktop -->
    <div class="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden hidden md:block">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">No. Tagihan</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Unit</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Periode</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Jumlah</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Jatuh Tempo</th>
              <th scope="col" class="relative px-6 py-3">
                <span class="sr-only">Actions</span>
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-if="!loading && bills.length === 0">
              <td colspan="7" class="px-6 py-12 text-center">
                <div class="text-gray-400">
                  <Icon name="heroicons:document-text" class="w-12 h-12 mx-auto mb-3 opacity-50" />
                  <p class="text-sm">Belum ada tagihan</p>
                </div>
              </td>
            </tr>
            <tr v-else v-for="bill in bills" :key="bill.id" class="hover:bg-gray-50 transition-colors group">
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm font-mono font-medium text-gray-900">
                  {{ bill.bill_number || '-' }}
                </div>
                <div v-if="bill.bill_number" class="text-xs text-gray-500 mt-0.5">ID: {{ bill.id.substring(0, 8) }}...</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center">
                  <div class="flex-shrink-0 h-10 w-10 bg-primary-100 rounded-lg flex items-center justify-center text-primary-700 font-bold text-sm">
                    {{ bill.unit_code || '-' }}
                  </div>
                  <div class="ml-4">
                    <div class="text-sm font-medium text-gray-900">{{ bill.unit_code || '-' }}</div>
                    <div class="text-sm text-gray-500">{{ bill.unit_type || '-' }}</div>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm text-gray-900">{{ bill.period }}</div>
                <div class="text-sm text-gray-500">{{ bill.category }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm font-bold text-gray-900">Rp {{ (bill.total_amount || bill.amount || 0).toLocaleString('id-ID') }}</div>
                <div v-if="bill.late_fee > 0" class="text-xs text-red-500">+ Denda Rp {{ bill.late_fee.toLocaleString('id-ID') }}</div>
                <div v-if="bill.total_amount && bill.total_amount !== bill.amount" class="text-xs text-gray-500">
                  Pokok: Rp {{ (bill.amount || 0).toLocaleString('id-ID') }}
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span 
                  :class="[
                    'px-2.5 py-0.5 rounded-full text-xs font-medium border',
                    statusClasses[bill.status] || 'bg-gray-100 text-gray-800 border-gray-200'
                  ]"
                >
                  {{ statusLabels[bill.status] || bill.status }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ formatDate(bill.due_date) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                <div class="relative inline-block text-left" @click.stop>
                  <button
                    @click.stop="toggleActionsMenu(bill.id)"
                    type="button"
                    class="text-gray-400 hover:text-primary-600 transition-colors focus:outline-none"
                  >
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"></path></svg>
                  </button>
                  <!-- Dropdown Menu -->
                  <Transition
                    enter-active-class="transition ease-out duration-100"
                    enter-from-class="transform opacity-0 scale-95"
                    enter-to-class="transform opacity-100 scale-100"
                    leave-active-class="transition ease-in duration-75"
                    leave-from-class="transform opacity-100 scale-100"
                    leave-to-class="transform opacity-0 scale-95"
                  >
                    <div
                      v-if="showActionsMenu === bill.id"
                      class="absolute right-0 mt-2 w-48 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5 z-50"
                      @click.stop
                    >
                      <div class="py-1">
                        <button
                          @click.stop="viewBill(bill.id)"
                          type="button"
                          class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 transition-colors"
                        >
                          <svg class="w-4 h-4 inline mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path></svg>
                          Lihat Detail
                        </button>
                        <button
                          @click.stop="editBill(bill.id)"
                          type="button"
                          class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 transition-colors"
                        >
                          <svg class="w-4 h-4 inline mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path></svg>
                          Edit
                        </button>
                        <button
                          @click.stop="deleteBill(bill.id)"
                          type="button"
                          class="block w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-red-50 transition-colors"
                        >
                          <svg class="w-4 h-4 inline mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
                          Hapus
                        </button>
                      </div>
                    </div>
                  </Transition>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <!-- Pagination -->
      <div v-if="pagination.total_pages > 1" class="bg-white px-4 py-3 border-t border-gray-200 flex flex-col sm:flex-row items-center justify-between gap-4 sm:px-6">
        <div class="flex-1 flex flex-col sm:flex-row items-center justify-between gap-4 w-full">
          <div class="text-center sm:text-left">
            <p class="text-xs sm:text-sm text-gray-700">
              Menampilkan <span class="font-medium">{{ (pagination.page - 1) * pagination.limit + 1 }}</span> sampai
              <span class="font-medium">{{ Math.min(pagination.page * pagination.limit, pagination.total) }}</span> dari
              <span class="font-medium">{{ pagination.total }}</span> tagihan
            </p>
          </div>
          <div>
            <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px" aria-label="Pagination">
              <button
                @click="pagination.page > 1 && (pagination.page--, loadBills())"
                :disabled="pagination.page === 1"
                class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <span class="sr-only">Previous</span>
                <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                  <path fill-rule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clip-rule="evenodd" />
                </svg>
              </button>
              <button
                v-for="page in Math.min(5, pagination.total_pages)"
                :key="page"
                @click="pagination.page = page; loadBills()"
                :class="[
                  'relative inline-flex items-center px-4 py-2 border text-sm font-medium',
                  pagination.page === page
                    ? 'z-10 bg-primary-50 border-primary-500 text-primary-600'
                    : 'bg-white border-gray-300 text-gray-500 hover:bg-gray-50'
                ]"
              >
                {{ page }}
              </button>
              <button
                @click="pagination.page < pagination.total_pages && (pagination.page++, loadBills())"
                :disabled="pagination.page >= pagination.total_pages"
                class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <span class="sr-only">Next</span>
                <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                  <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
                </svg>
              </button>
            </nav>
          </div>
        </div>
      </div>
    </div>

    <!-- Bills Cards - Mobile -->
    <div class="md:hidden space-y-4">
      <div v-if="!loading && bills.length === 0" class="bg-white rounded-2xl border border-gray-100 shadow-sm p-12 text-center">
        <div class="text-gray-400">
          <Icon name="heroicons:document-text" class="w-12 h-12 mx-auto mb-3 opacity-50" />
          <p class="text-sm">Belum ada tagihan</p>
        </div>
      </div>
      <div
        v-for="bill in bills"
        :key="bill.id"
        class="bg-white rounded-2xl border border-gray-100 shadow-sm p-4 space-y-3"
      >
        <div class="flex items-start justify-between">
          <div class="flex-1">
            <div class="flex items-center gap-2 mb-2">
              <div class="flex-shrink-0 h-10 w-10 bg-primary-100 rounded-lg flex items-center justify-center text-primary-700 font-bold text-sm">
                {{ bill.unit_code || '-' }}
              </div>
              <div>
                <div class="text-sm font-medium text-gray-900">{{ bill.unit_code || '-' }}</div>
                <div class="text-xs text-gray-500">{{ bill.unit_type || '-' }}</div>
              </div>
            </div>
            <div class="text-xs font-mono text-gray-500 mb-1">No: {{ bill.bill_number || '-' }}</div>
          </div>
          <div class="relative" @click.stop>
            <button
              @click.stop="toggleActionsMenu(bill.id)"
              type="button"
              class="text-gray-400 hover:text-primary-600 transition-colors focus:outline-none p-1"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"></path></svg>
            </button>
            <Transition
              enter-active-class="transition ease-out duration-100"
              enter-from-class="transform opacity-0 scale-95"
              enter-to-class="transform opacity-100 scale-100"
              leave-active-class="transition ease-in duration-75"
              leave-from-class="transform opacity-100 scale-100"
              leave-to-class="transform opacity-0 scale-95"
            >
              <div
                v-if="showActionsMenu === bill.id"
                class="absolute right-0 mt-2 w-48 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5 z-50"
                @click.stop
              >
                <div class="py-1">
                  <button
                    @click.stop="viewBill(bill.id)"
                    type="button"
                    class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 transition-colors"
                  >
                    <svg class="w-4 h-4 inline mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path></svg>
                    Lihat Detail
                  </button>
                  <button
                    @click.stop="editBill(bill.id)"
                    type="button"
                    class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 transition-colors"
                  >
                    <svg class="w-4 h-4 inline mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path></svg>
                    Edit
                  </button>
                  <button
                    @click.stop="deleteBill(bill.id)"
                    type="button"
                    class="block w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-red-50 transition-colors"
                  >
                    <svg class="w-4 h-4 inline mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
                    Hapus
                  </button>
                </div>
              </div>
            </Transition>
          </div>
        </div>
        
        <div class="grid grid-cols-2 gap-3 pt-2 border-t border-gray-100">
          <div>
            <div class="text-xs text-gray-500 mb-1">Periode</div>
            <div class="text-sm font-medium text-gray-900">{{ bill.period }}</div>
            <div class="text-xs text-gray-500">{{ bill.category }}</div>
          </div>
          <div>
            <div class="text-xs text-gray-500 mb-1">Jatuh Tempo</div>
            <div class="text-sm text-gray-900">{{ formatDate(bill.due_date) }}</div>
          </div>
        </div>
        
        <div class="flex items-center justify-between pt-2 border-t border-gray-100">
          <div>
            <div class="text-xs text-gray-500 mb-1">Jumlah</div>
            <div class="text-base font-bold text-gray-900">Rp {{ (bill.total_amount || bill.amount || 0).toLocaleString('id-ID') }}</div>
            <div v-if="bill.late_fee > 0" class="text-xs text-red-500">+ Denda Rp {{ bill.late_fee.toLocaleString('id-ID') }}</div>
          </div>
          <span 
            :class="[
              'px-2.5 py-1 rounded-full text-xs font-medium border',
              statusClasses[bill.status] || 'bg-gray-100 text-gray-800 border-gray-200'
            ]"
          >
            {{ statusLabels[bill.status] || bill.status }}
          </span>
        </div>
      </div>
      
      <!-- Mobile Pagination -->
      <div v-if="pagination.total_pages > 1" class="bg-white rounded-2xl border border-gray-100 shadow-sm p-4">
        <div class="flex flex-col items-center gap-4">
          <p class="text-xs text-gray-700 text-center">
            Halaman <span class="font-medium">{{ pagination.page }}</span> dari
            <span class="font-medium">{{ pagination.total_pages }}</span>
          </p>
          <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px" aria-label="Pagination">
            <button
              @click="pagination.page > 1 && (pagination.page--, loadBills())"
              :disabled="pagination.page === 1"
              class="relative inline-flex items-center px-3 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span class="sr-only">Previous</span>
              <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                <path fill-rule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clip-rule="evenodd" />
              </svg>
            </button>
            <button
              v-for="page in Math.min(5, pagination.total_pages)"
              :key="page"
              @click="pagination.page = page; loadBills()"
              :class="[
                'relative inline-flex items-center px-4 py-2 border text-sm font-medium',
                pagination.page === page
                  ? 'z-10 bg-primary-50 border-primary-500 text-primary-600'
                  : 'bg-white border-gray-300 text-gray-500 hover:bg-gray-50'
              ]"
            >
              {{ page }}
            </button>
            <button
              @click="pagination.page < pagination.total_pages && (pagination.page++, loadBills())"
              :disabled="pagination.page >= pagination.total_pages"
              class="relative inline-flex items-center px-3 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span class="sr-only">Next</span>
              <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
              </svg>
            </button>
          </nav>
        </div>
      </div>
    </div>

    <!-- View Bill Modal -->
    <div v-if="showViewModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4" @click.self="showViewModal = false">
      <div class="bg-white rounded-2xl p-4 md:p-6 max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        <div class="flex justify-between items-center mb-6">
          <h3 class="text-lg font-bold text-gray-900">Detail Tagihan</h3>
          <button @click="showViewModal = false" class="text-gray-400 hover:text-gray-600">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
          </button>
        </div>
        
        <div v-if="viewingBill" class="space-y-4">
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label class="text-sm font-medium text-gray-500">Unit</label>
              <p class="text-sm text-gray-900 mt-1">{{ viewingBill.unit_code || '-' }} ({{ viewingBill.unit_type || '-' }})</p>
            </div>
            <div>
              <label class="text-sm font-medium text-gray-500">Status</label>
              <p class="mt-1">
                <span 
                  :class="[
                    'px-2.5 py-0.5 rounded-full text-xs font-medium border',
                    statusClasses[viewingBill.status] || 'bg-gray-100 text-gray-800 border-gray-200'
                  ]"
                >
                  {{ statusLabels[viewingBill.status] || viewingBill.status }}
                </span>
              </p>
            </div>
            <div>
              <label class="text-sm font-medium text-gray-500">Kategori</label>
              <p class="text-sm text-gray-900 mt-1">{{ viewingBill.category || '-' }}</p>
            </div>
            <div>
              <label class="text-sm font-medium text-gray-500">Periode</label>
              <p class="text-sm text-gray-900 mt-1">{{ viewingBill.period || '-' }}</p>
            </div>
            <div>
              <label class="text-sm font-medium text-gray-500">Jumlah Tagihan</label>
              <p class="text-sm font-bold text-gray-900 mt-1">Rp {{ (viewingBill.amount || 0).toLocaleString('id-ID') }}</p>
            </div>
            <div>
              <label class="text-sm font-medium text-gray-500">Denda</label>
              <p class="text-sm text-gray-900 mt-1">Rp {{ (viewingBill.late_fee || 0).toLocaleString('id-ID') }}</p>
            </div>
            <div>
              <label class="text-sm font-medium text-gray-500">Jatuh Tempo</label>
              <p class="text-sm text-gray-900 mt-1">{{ formatDate(viewingBill.due_date) || '-' }}</p>
            </div>
            <div v-if="viewingBill.paid_at">
              <label class="text-sm font-medium text-gray-500">Tanggal Bayar</label>
              <p class="text-sm text-gray-900 mt-1">{{ formatDate(viewingBill.paid_at) }}</p>
            </div>
            <div v-if="viewingBill.payment_method">
              <label class="text-sm font-medium text-gray-500">Metode Pembayaran</label>
              <p class="text-sm text-gray-900 mt-1">{{ viewingBill.payment_method }}</p>
            </div>
            <div v-if="viewingBill.payment_reference">
              <label class="text-sm font-medium text-gray-500">Referensi Pembayaran</label>
              <p class="text-sm text-gray-900 mt-1">{{ viewingBill.payment_reference }}</p>
            </div>
          </div>
          <div v-if="viewingBill.notes">
            <label class="text-sm font-medium text-gray-500">Catatan</label>
            <p class="text-sm text-gray-900 mt-1">{{ viewingBill.notes }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Edit Bill Modal -->
    <div v-if="showEditModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4" @click.self="showEditModal = false">
      <div class="bg-white rounded-2xl p-4 md:p-6 max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        <div class="flex justify-between items-center mb-6">
          <h3 class="text-lg font-bold text-gray-900">Edit Tagihan</h3>
          <button @click="showEditModal = false" class="text-gray-400 hover:text-gray-600">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
          </button>
        </div>
        
        <div v-if="editingBill" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Unit</label>
            <select
              v-model="editForm.unit_id"
              class="w-full border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
            >
              <option value="">Pilih Unit</option>
              <option v-for="unit in allUnits" :key="unit.id" :value="unit.id">
                {{ unit.code }} - {{ unit.type }}
              </option>
            </select>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Kategori</label>
            <input
              v-model="editForm.category"
              type="text"
              class="w-full border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
            >
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Periode</label>
            <input
              v-model="editForm.period"
              type="text"
              class="w-full border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
            >
          </div>
          
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Jumlah Tagihan</label>
              <input
                v-model.number="editForm.amount"
                type="number"
                class="w-full border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                min="0"
              >
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Denda</label>
              <input
                v-model.number="editForm.late_fee"
                type="number"
                class="w-full border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                min="0"
              >
            </div>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Jatuh Tempo</label>
            <input
              v-model="editForm.due_date"
              type="date"
              class="w-full border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
            >
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Status</label>
            <select
              v-model="editForm.status"
              class="w-full border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
            >
              <option value="pending">Belum Bayar</option>
              <option value="paid">Lunas</option>
              <option value="overdue">Terlambat</option>
              <option value="cancelled">Dibatalkan</option>
            </select>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Catatan</label>
            <textarea
              v-model="editForm.notes"
              rows="3"
              class="w-full border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
            ></textarea>
          </div>
          
          <div class="flex flex-col sm:flex-row gap-3 mt-6">
            <button
              @click="showEditModal = false"
              class="flex-1 px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors"
            >
              Batal
            </button>
            <button
              @click="saveBill"
              :disabled="saving"
              :class="['flex-1 px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors', saving ? 'opacity-50 cursor-not-allowed' : '']"
            >
              {{ saving ? 'Menyimpan...' : 'Simpan' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useGlobalLoading } from '~/composables/useGlobalLoading'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

const { fetch } = useApi()
const { showSuccess, showError, showWarning } = useToast()
const confirm = useConfirm()

const activeTab = ref('all')
const bills = ref([])
const pagination = ref({
  page: 1,
  limit: 20,
  total: 0,
  total_pages: 0,
})
const loading = ref(false)

// Sync loading state dengan global loading spinner
useGlobalLoading(loading, 'Memuat data tagihan...')
const searchQuery = ref('')
const showActionsMenu = ref(null)
const selectedBillId = ref(null)
const showViewModal = ref(false)
const showEditModal = ref(false)
const viewingBill = ref(null)
const editingBill = ref(null)
const saving = ref(false)
const allUnits = ref([])

const editForm = ref({
  unit_id: '',
  category: '',
  period: '',
  amount: 0,
  late_fee: 0,
  due_date: '',
  status: 'pending',
  notes: '',
})

const tabs = computed(() => {
  const all = bills.value.length
  const pending = bills.value.filter((b: any) => b.status === 'pending').length
  const paid = bills.value.filter((b: any) => b.status === 'paid').length
  const overdue = bills.value.filter((b: any) => b.status === 'overdue').length
  
  return [
    { id: 'all', label: 'Semua', count: all },
    { id: 'pending', label: 'Belum Bayar', count: pending },
    { id: 'paid', label: 'Lunas', count: paid },
    { id: 'overdue', label: 'Terlambat', count: overdue },
  ]
})

const statusClasses = {
  'paid': 'bg-green-100 text-green-800 border-green-200',
  'pending': 'bg-yellow-100 text-yellow-800 border-yellow-200',
  'overdue': 'bg-red-100 text-red-800 border-red-200',
  'cancelled': 'bg-gray-100 text-gray-800 border-gray-200',
}

const statusLabels = {
  'paid': 'Lunas',
  'pending': 'Belum Bayar',
  'overdue': 'Terlambat',
  'cancelled': 'Dibatalkan',
}

const loadBills = async () => {
  loading.value = true
  try {
    const params = new URLSearchParams({
      page: pagination.value.page.toString(),
      limit: pagination.value.limit.toString(),
    })
    
    if (activeTab.value !== 'all') {
      params.append('status', activeTab.value)
    }
    
    if (searchQuery.value) {
      params.append('search', searchQuery.value)
    }

    const response = await fetch(`/api/billing?${params.toString()}`)
    bills.value = (response.bills || []).map((bill: any) => ({
      ...bill,
      unit: bill.unit_code || '-',
      type: bill.unit_type || '-',
      lateFee: bill.late_fee || 0,
    }))
    
    if (response.pagination) {
      pagination.value = response.pagination
    }
  } catch (error) {
    console.error('Failed to load bills:', error)
  } finally {
    loading.value = false
  }
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric'
  })
}

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
  }).format(amount)
}

const toggleActionsMenu = (billId: string) => {
  if (showActionsMenu.value === billId) {
    showActionsMenu.value = null
  } else {
    showActionsMenu.value = billId
  }
}

const loadUnits = async () => {
  try {
    const response = await fetch('/api/units?limit=1000')
    allUnits.value = response.units || []
  } catch (error) {
    console.error('Failed to load units:', error)
  }
}

const viewBill = async (billId: string) => {
  showActionsMenu.value = null
  try {
    const response = await fetch(`/api/billing/${billId}`)
    viewingBill.value = response
    showViewModal.value = true
  } catch (error: any) {
    showError(error.message || 'Terjadi kesalahan', 'Gagal memuat detail tagihan')
  }
}

const editBill = async (billId: string) => {
  showActionsMenu.value = null
  try {
    const response = await fetch(`/api/billing/${billId}`)
    editingBill.value = response
    
    // Format due_date for date input
    let dueDate = ''
    if (response.due_date) {
      const date = new Date(response.due_date)
      dueDate = date.toISOString().split('T')[0]
    }
    
    editForm.value = {
      unit_id: response.unit_id || '',
      category: response.category || '',
      period: response.period || '',
      amount: response.amount || 0,
      late_fee: response.late_fee || 0,
      due_date: dueDate,
      status: response.status || 'pending',
      notes: response.notes || '',
    }
    
    showEditModal.value = true
  } catch (error: any) {
    showError(error.message || 'Terjadi kesalahan', 'Gagal memuat data tagihan')
  }
}

const saveBill = async () => {
  if (!editingBill.value) return
  
  saving.value = true
  try {
    const updateData: any = {}
    
    if (editForm.value.category !== editingBill.value.category) {
      updateData.category = editForm.value.category
    }
    if (editForm.value.period !== editingBill.value.period) {
      updateData.period = editForm.value.period
    }
    if (editForm.value.amount !== editingBill.value.amount) {
      updateData.amount = editForm.value.amount
    }
    if (editForm.value.late_fee !== (editingBill.value.late_fee || 0)) {
      updateData.late_fee = editForm.value.late_fee
    }
    if (editForm.value.status !== editingBill.value.status) {
      updateData.status = editForm.value.status
    }
    if (editForm.value.notes !== (editingBill.value.notes || '')) {
      updateData.notes = editForm.value.notes
    }
    
    // Handle due_date
    const currentDueDate = editingBill.value.due_date ? new Date(editingBill.value.due_date).toISOString().split('T')[0] : ''
    if (editForm.value.due_date !== currentDueDate) {
      updateData.due_date = editForm.value.due_date || null
    }
    
    await fetch(`/api/billing/${editingBill.value.id}`, {
      method: 'PUT',
      body: JSON.stringify(updateData),
    })
    
    showEditModal.value = false
    editingBill.value = null
    await loadBills()
    showSuccess('Tagihan berhasil diperbarui')
  } catch (error: any) {
    showError(error.message || 'Terjadi kesalahan', 'Gagal memperbarui tagihan')
  } finally {
    saving.value = false
  }
}

const deleteBill = async (billId: string) => {
  showActionsMenu.value = null
  
  const result = await confirm.show(
    'Yakin ingin menghapus tagihan ini? Tindakan ini tidak dapat dibatalkan.',
    {
      title: 'Hapus Tagihan',
      confirmText: 'Ya, Hapus',
      cancelText: 'Batal',
      type: 'danger'
    }
  )
  
  if (!result) {
    return
  }
  
  try {
    await fetch(`/api/billing/${billId}`, {
      method: 'DELETE',
    })
    
    await loadBills()
    showSuccess('Tagihan berhasil dihapus')
  } catch (error: any) {
    showError(error.message || 'Terjadi kesalahan', 'Gagal menghapus tagihan')
  }
}

// Close dropdown when clicking outside
onMounted(() => {
  loadBills()
  loadUnits()
  
  const handleClickOutside = (event: MouseEvent) => {
    const target = event.target as HTMLElement
    if (!target.closest('.relative.inline-block')) {
      showActionsMenu.value = null
    }
  }
  
  document.addEventListener('click', handleClickOutside)
  
  onUnmounted(() => {
    document.removeEventListener('click', handleClickOutside)
  })
})

watch(activeTab, () => {
  pagination.value.page = 1
  loadBills()
})
</script>
