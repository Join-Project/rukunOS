<template>
  <div class="max-w-6xl mx-auto space-y-8">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 tracking-tight">Visitor Management</h1>
        <p class="text-gray-500">Catat dan pantau tamu yang masuk keluar lingkungan</p>
      </div>
      <button @click="showForm = !showForm" class="bg-primary-600 hover:bg-primary-700 text-white px-4 py-2 rounded-xl font-medium transition-colors flex items-center gap-2">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path></svg>
        Catat Tamu Baru
      </button>
    </div>

    <!-- Stats Overview -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div class="bg-white p-6 rounded-2xl border border-gray-100 shadow-sm">
        <div class="flex items-center gap-4">
          <div class="w-12 h-12 bg-blue-50 rounded-xl flex items-center justify-center text-blue-600">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"></path></svg>
          </div>
          <div>
            <p class="text-sm text-gray-500 font-medium">Total Tamu Hari Ini</p>
            <h3 class="text-2xl font-bold text-gray-900">{{ stats.totalToday }}</h3>
          </div>
        </div>
      </div>
      <div class="bg-white p-6 rounded-2xl border border-gray-100 shadow-sm">
        <div class="flex items-center gap-4">
          <div class="w-12 h-12 bg-green-50 rounded-xl flex items-center justify-center text-green-600">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 16l-4-4m0 0l4-4m-4 4h14m-5 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h7a3 3 0 013 3v1"></path></svg>
          </div>
          <div>
            <p class="text-sm text-gray-500 font-medium">Sedang Berkunjung</p>
            <h3 class="text-2xl font-bold text-gray-900">{{ stats.checkedIn }}</h3>
          </div>
        </div>
      </div>
      <div class="bg-white p-6 rounded-2xl border border-gray-100 shadow-sm">
        <div class="flex items-center gap-4">
          <div class="w-12 h-12 bg-yellow-50 rounded-xl flex items-center justify-center text-yellow-600">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"></path></svg>
          </div>
          <div>
            <p class="text-sm text-gray-500 font-medium">Sudah Check-Out</p>
            <h3 class="text-2xl font-bold text-gray-900">{{ stats.checkedOut }}</h3>
          </div>
        </div>
      </div>
    </div>

    <!-- Entry Form -->
    <div v-if="showForm" class="bg-white rounded-2xl border border-gray-100 p-6 shadow-sm animate-fade-in">
      <h3 class="text-lg font-bold text-gray-900 mb-4">Catat Tamu Masuk</h3>
      <form @submit.prevent="submitVisitor" class="space-y-4">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Nama Pengunjung *</label>
            <input 
              v-model="form.visitor_name" 
              type="text" 
              class="w-full rounded-xl border-gray-300 focus:border-primary-500 focus:ring-primary-500" 
              placeholder="Nama lengkap sesuai KTP"
              required
            >
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Nomor Telepon</label>
            <input 
              v-model="form.visitor_phone" 
              type="tel" 
              class="w-full rounded-xl border-gray-300 focus:border-primary-500 focus:ring-primary-500" 
              placeholder="08xxx"
            >
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Unit Tujuan</label>
            <select 
              v-model="form.unit_id" 
              class="w-full rounded-xl border-gray-300 focus:border-primary-500 focus:ring-primary-500"
            >
              <option value="">Pilih Unit</option>
              <option v-for="unit in units" :key="unit.id" :value="unit.id">{{ unit.code }}</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Nama yang Dikunjungi</label>
            <input 
              v-model="form.host_name" 
              type="text" 
              class="w-full rounded-xl border-gray-300 focus:border-primary-500 focus:ring-primary-500" 
              placeholder="Nama penghuni unit"
            >
          </div>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Keperluan</label>
          <textarea 
            v-model="form.purpose" 
            rows="2" 
            class="w-full rounded-xl border-gray-300 focus:border-primary-500 focus:ring-primary-500" 
            placeholder="Jelaskan tujuan kunjungan..."
          ></textarea>
        </div>
        <div class="flex justify-end pt-4">
          <button 
            type="submit" 
            :disabled="saving || !form.visitor_name"
            :class="['bg-primary-600 hover:bg-primary-700 text-white px-6 py-2.5 rounded-xl font-bold shadow-lg shadow-primary-600/20 transition-all', (saving || !form.visitor_name) ? 'opacity-50 cursor-not-allowed' : '']"
          >
            {{ saving ? 'Menyimpan...' : 'Catat Masuk (Check-In)' }}
          </button>
        </div>
      </form>
    </div>

    <!-- Visitor List -->
    <div class="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden">
      <div class="p-6 border-b border-gray-100 flex justify-between items-center">
        <h3 class="text-lg font-bold text-gray-900">Daftar Pengunjung</h3>
        <div class="flex gap-2">
          <select 
            v-model="filterStatus" 
            @change="loadVisitors" 
            class="rounded-lg border-gray-200 text-sm focus:border-primary-500 focus:ring-primary-500"
          >
            <option value="">Semua Status</option>
            <option value="checked_in">Sedang Berkunjung</option>
            <option value="checked_out">Sudah Check-Out</option>
          </select>
          <div class="relative">
            <input 
              v-model="searchQuery" 
              @input="loadVisitors" 
              type="text" 
              placeholder="Cari nama/unit..." 
              class="pl-9 pr-4 py-2 rounded-lg border-gray-200 text-sm focus:border-primary-500 focus:ring-primary-500"
            >
            <svg class="w-4 h-4 text-gray-400 absolute left-3 top-1/2 -translate-y-1/2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path></svg>
          </div>
        </div>
      </div>
      <div v-if="!loading && visitors.length === 0" class="p-12 text-center text-gray-500">Belum ada data pengunjung</div>
      <div v-else class="overflow-x-auto">
        <table class="w-full text-left text-sm">
          <thead class="bg-gray-50 text-gray-500 font-medium">
            <tr>
              <th class="px-6 py-4">Nama</th>
              <th class="px-6 py-4">Unit Tujuan</th>
              <th class="px-6 py-4">Keperluan</th>
              <th class="px-6 py-4">Waktu Masuk</th>
              <th class="px-6 py-4">Waktu Keluar</th>
              <th class="px-6 py-4 text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            <tr v-for="visitor in visitors" :key="visitor.id" class="hover:bg-gray-50 transition-colors">
              <td class="px-6 py-4 font-medium text-gray-900">{{ visitor.visitor_name }}</td>
              <td class="px-6 py-4">{{ visitor.unit_code || '-' }}</td>
              <td class="px-6 py-4 text-gray-600">{{ visitor.purpose || '-' }}</td>
              <td class="px-6 py-4 text-gray-500">{{ formatDateTime(visitor.checked_in_at) }}</td>
              <td class="px-6 py-4 text-gray-500">{{ visitor.checked_out_at ? formatDateTime(visitor.checked_out_at) : '-' }}</td>
              <td class="px-6 py-4 text-right">
                <button 
                  v-if="!visitor.checked_out_at"
                  @click="checkOutVisitor(visitor.id)" 
                  class="text-red-600 hover:text-red-700 font-medium text-xs border border-red-200 hover:bg-red-50 px-3 py-1.5 rounded-lg transition-colors"
                >
                  Check-Out
                </button>
                <span v-else class="text-xs text-gray-400">Selesai</span>
              </td>
            </tr>
          </tbody>
        </table>
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
const { showSuccess, showError } = useToast()

const showForm = ref(false)
const loading = ref(false)

// Sync loading state dengan global loading spinner
useGlobalLoading(loading, 'Memuat data pengunjung...')
const saving = ref(false)
const visitors = ref([])
const units = ref([])
const searchQuery = ref('')
const filterStatus = ref('')

const stats = computed(() => {
  const today = new Date().toISOString().split('T')[0]
  const todayVisitors = visitors.value.filter((v: any) => 
    v.checked_in_at && v.checked_in_at.startsWith(today)
  )
  const checkedIn = visitors.value.filter((v: any) => !v.checked_out_at)
  const checkedOut = visitors.value.filter((v: any) => v.checked_out_at)
  
  return {
    totalToday: todayVisitors.length,
    checkedIn: checkedIn.length,
    checkedOut: checkedOut.length,
  }
})

const form = ref({
  visitor_name: '',
  visitor_phone: '',
  unit_id: '',
  host_name: '',
  purpose: '',
})

const loadUnits = async () => {
  try {
    const response = await fetch('/api/units?limit=1000')
    units.value = response.units || []
  } catch (error) {
    console.error('Failed to load units:', error)
  }
}

const loadVisitors = async () => {
  loading.value = true
  try {
    const params = new URLSearchParams({
      page: '1',
      limit: '100',
    })
    if (searchQuery.value) {
      params.append('search', searchQuery.value)
    }
    if (filterStatus.value) {
      params.append('status', filterStatus.value)
    }

    const response = await fetch(`/api/visitors?${params.toString()}`)
    visitors.value = response.visitors || []
  } catch (error: any) {
    console.error('Failed to load visitors:', error)
    showError(error.message || 'Gagal memuat data pengunjung')
  } finally {
    loading.value = false
  }
}

const submitVisitor = async () => {
  if (!form.value.visitor_name) {
    showError('Nama pengunjung wajib diisi')
    return
  }

  saving.value = true
  try {
    const data: any = {
      visitor_name: form.value.visitor_name,
    }
    if (form.value.visitor_phone) {
      data.visitor_phone = form.value.visitor_phone
    }
    if (form.value.unit_id) {
      data.unit_id = form.value.unit_id
    }
    if (form.value.host_name) {
      data.host_name = form.value.host_name
    }
    if (form.value.purpose) {
      data.purpose = form.value.purpose
    }

    await fetch('/api/visitors', {
      method: 'POST',
      body: JSON.stringify(data),
    })

    showSuccess('Data pengunjung berhasil dicatat!', 'Data Tersimpan')
    showForm.value = false
    form.value = {
      visitor_name: '',
      visitor_phone: '',
      unit_id: '',
      host_name: '',
      purpose: '',
    }
    await loadVisitors()
  } catch (error: any) {
    console.error('Failed to create visitor log:', error)
    showError(error.message || 'Gagal mencatat pengunjung')
  } finally {
    saving.value = false
  }
}

const checkOutVisitor = async (visitorId: string) => {
  try {
    await fetch(`/api/visitors/${visitorId}/checkout`, {
      method: 'POST',
    })
    showSuccess('Pengunjung berhasil check-out')
    await loadVisitors()
  } catch (error: any) {
    console.error('Failed to check out visitor:', error)
    showError(error.message || 'Gagal melakukan check-out')
  }
}

const formatDateTime = (dateStr: string) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

onMounted(() => {
  loadUnits()
  loadVisitors()
})
</script>
