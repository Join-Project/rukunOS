<template>
  <div class="max-w-4xl mx-auto space-y-8">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-bold text-gray-900 tracking-tight">Lapor Masalah</h1>
          <p class="text-gray-500">Laporkan keluhan atau masalah di lingkungan Anda</p>
        </div>
        <button @click="showForm = !showForm" class="bg-primary-600 hover:bg-primary-700 text-white px-4 py-2 rounded-xl font-medium transition-colors">
          {{ showForm ? 'Tutup Form' : 'Buat Laporan Baru' }}
        </button>
      </div>

      <!-- Complaint Form -->
      <div v-if="showForm" class="bg-white rounded-2xl border border-gray-100 p-6 shadow-sm animate-fade-in">
        <h3 class="text-lg font-bold text-gray-900 mb-4">Formulir Laporan</h3>
        <form @submit.prevent="submitComplaint" class="space-y-4">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Kategori</label>
              <select v-model="form.category" class="w-full rounded-xl border-gray-300 focus:border-primary-500 focus:ring-primary-500">
                <option value="">Pilih Kategori</option>
                <option value="keamanan">Keamanan</option>
                <option value="kebersihan">Kebersihan / Sampah</option>
                <option value="fasilitas_umum">Fasilitas Umum</option>
                <option value="lainnya">Lainnya</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Prioritas</label>
              <select v-model="form.priority" class="w-full rounded-xl border-gray-300 focus:border-primary-500 focus:ring-primary-500">
                <option value="normal">Normal</option>
                <option value="penting">Penting</option>
                <option value="darurat">Darurat</option>
              </select>
            </div>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Judul Laporan</label>
            <input 
              v-model="form.title" 
              type="text" 
              class="w-full rounded-xl border-gray-300 focus:border-primary-500 focus:ring-primary-500" 
              placeholder="Contoh: Lampu jalan mati di Blok A"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Deskripsi Detail</label>
            <textarea 
              v-model="form.description" 
              rows="4" 
              class="w-full rounded-xl border-gray-300 focus:border-primary-500 focus:ring-primary-500" 
              placeholder="Jelaskan masalah secara detail..."
            ></textarea>
          </div>

          <div class="flex justify-end pt-4">
            <button 
              type="submit" 
              :disabled="saving || !form.category || !form.title || !form.description"
              class="bg-primary-600 hover:bg-primary-700 text-white px-6 py-2.5 rounded-xl font-bold shadow-lg shadow-primary-600/20 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {{ saving ? 'Mengirim...' : 'Kirim Laporan' }}
            </button>
          </div>
        </form>
      </div>

      <!-- History List -->
      <div class="space-y-4">
        <h3 class="text-lg font-bold text-gray-900">Riwayat Laporan</h3>
        
        <div v-if="!loading && complaints.length === 0 && pagination.total === 0" class="text-center py-12 text-gray-500">
          <p>Belum ada laporan</p>
        </div>
        <template v-else>
          <div v-if="complaints.length > 0" class="space-y-4 mb-4">
            <div 
              v-for="complaint in complaints" 
              :key="complaint.id"
              class="bg-white rounded-2xl border border-gray-100 p-6 shadow-sm hover:shadow-md transition-shadow"
            >
              <div class="flex items-start justify-between">
                <div class="flex gap-4">
                  <div :class="[
                    'w-12 h-12 rounded-xl flex items-center justify-center shrink-0',
                    complaint.status === 'resolved' ? 'bg-green-100 text-green-600' :
                    complaint.status === 'in_progress' ? 'bg-blue-100 text-blue-600' :
                    'bg-yellow-100 text-yellow-600'
                  ]">
                    <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"></path></svg>
                  </div>
                  <div>
                    <div class="flex items-center gap-2 mb-1">
                      <span class="px-2 py-0.5 bg-gray-100 text-gray-600 text-xs font-bold rounded">{{ getCategoryLabel(complaint.category) }}</span>
                      <span class="text-xs text-gray-400">{{ formatTimeAgo(complaint.created_at) }}</span>
                    </div>
                    <h4 class="font-bold text-gray-900 text-lg">{{ complaint.title }}</h4>
                    <p class="text-gray-600 mt-1">{{ complaint.description }}</p>
                    <div v-if="complaint.resolution_notes" class="mt-2 p-2 bg-green-50 rounded-lg">
                      <p class="text-sm text-green-800"><strong>Resolusi:</strong> {{ complaint.resolution_notes }}</p>
                    </div>
                  </div>
                </div>
                <span :class="[
                  'px-3 py-1 text-sm font-bold rounded-full',
                  complaint.status === 'resolved' ? 'bg-green-50 text-green-700' :
                  complaint.status === 'in_progress' ? 'bg-blue-50 text-blue-700' :
                  'bg-yellow-50 text-yellow-700'
                ]">
                  {{ getStatusLabel(complaint.status) }}
                </span>
              </div>
            </div>
          </div>
          
          <!-- Pagination -->
          <div v-if="pagination.total > 0 && pagination.total_pages > 0" class="flex justify-center items-center gap-2 mt-4">
            <button
              @click="pagination.page > 1 && (pagination.page--, loadComplaints())"
              :disabled="pagination.page === 1"
              :class="['px-4 py-2 rounded-lg border transition-colors', pagination.page === 1 ? 'border-gray-200 text-gray-400 cursor-not-allowed bg-gray-50' : 'border-gray-300 text-gray-700 hover:bg-gray-50 bg-white']"
            >
              Sebelumnya
            </button>
            <span class="text-sm text-gray-600 px-4">
              Halaman {{ pagination.page }} dari {{ pagination.total_pages }}
            </span>
            <button
              @click="pagination.page < pagination.total_pages && (pagination.page++, loadComplaints())"
              :disabled="pagination.page >= pagination.total_pages"
              :class="['px-4 py-2 rounded-lg border transition-colors', pagination.page >= pagination.total_pages ? 'border-gray-200 text-gray-400 cursor-not-allowed bg-gray-50' : 'border-gray-300 text-gray-700 hover:bg-gray-50 bg-white']"
            >
              Selanjutnya
            </button>
          </div>
        </template>
      </div>
    </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useApi } from '~/composables/useApi'
import { useToast } from '~/composables/useToast'
import { useGlobalLoading } from '~/composables/useGlobalLoading'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

const { fetch } = useApi()
const { showSuccess, showError } = useToast()

const loading = ref(false)

// Sync loading state dengan global loading spinner
useGlobalLoading(loading, 'Memuat laporan...')
const saving = ref(false)
const showForm = ref(false)
const complaints = ref([])

const form = reactive({
  category: '',
  priority: 'normal',
  title: '',
  description: '',
})

const pagination = ref({
  page: 1,
  limit: 20,
  total: 0,
  total_pages: 0
})

const loadComplaints = async () => {
  loading.value = true
  try {
    const response = await fetch(`/api/complaints?page=${pagination.value.page}&limit=${pagination.value.limit}`)
    complaints.value = response.complaints || []
    if (response.pagination) {
      pagination.value = {
        page: response.pagination.page || pagination.value.page,
        limit: response.pagination.limit || pagination.value.limit,
        total: response.pagination.total || 0,
        total_pages: response.pagination.total_pages || 0
      }
    }
  } catch (error: any) {
    console.error('Failed to load complaints:', error)
    showError(error.message || 'Gagal memuat laporan')
  } finally {
    loading.value = false
  }
}

const submitComplaint = async () => {
  if (!form.category || !form.title || !form.description) {
    showError('Semua field wajib diisi')
    return
  }

  saving.value = true
  try {
    await fetch('/api/complaints', {
      method: 'POST',
      body: JSON.stringify({
        category: form.category,
        priority: form.priority,
        title: form.title,
        description: form.description,
      }),
    })
    showSuccess('Laporan berhasil dikirim! Petugas akan segera menindaklanjuti.')
    showForm.value = false
    form.category = ''
    form.priority = 'normal'
    form.title = ''
    form.description = ''
    await loadComplaints()
  } catch (error: any) {
    console.error('Failed to submit complaint:', error)
    showError(error.message || 'Gagal mengirim laporan')
  } finally {
    saving.value = false
  }
}

const getCategoryLabel = (category: string) => {
  const labels: Record<string, string> = {
    keamanan: 'Keamanan',
    kebersihan: 'Kebersihan',
    fasilitas_umum: 'Fasilitas Umum',
    lainnya: 'Lainnya',
  }
  return labels[category] || category
}

const getStatusLabel = (status: string) => {
  const labels: Record<string, string> = {
    pending: 'Menunggu',
    in_progress: 'Diproses',
    resolved: 'Selesai',
    rejected: 'Ditolak',
  }
  return labels[status] || status
}

const formatTimeAgo = (dateStr: string) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  
  if (days === 0) return 'Hari ini'
  if (days === 1) return 'Kemarin'
  if (days < 7) return `${days} hari yang lalu`
  
  return date.toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric'
  })
}

onMounted(() => {
  loadComplaints()
})
</script>
