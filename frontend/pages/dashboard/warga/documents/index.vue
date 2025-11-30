<template>
  <div class="max-w-4xl mx-auto">
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
        <div>
          <h1 class="text-2xl font-bold text-gray-900 tracking-tight">Layanan Surat</h1>
          <p class="text-gray-500">Ajukan surat pengantar RT/RW secara online.</p>
        </div>
        <button @click="showForm = true" class="px-4 py-2 bg-primary-600 text-white rounded-xl font-medium hover:bg-primary-700 transition-colors shadow-lg shadow-primary-600/20 flex items-center">
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path></svg>
          Ajukan Surat
        </button>
      </div>

      <!-- Request Form -->
      <div v-if="showForm" class="bg-white rounded-2xl border border-gray-100 shadow-sm p-6 mb-8 animate-fade-in">
        <div class="flex justify-between items-center mb-6">
          <h3 class="text-lg font-bold text-gray-900">Formulir Pengajuan Surat</h3>
          <button @click="closeForm" class="text-gray-400 hover:text-gray-600">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
          </button>
        </div>
        
        <form @submit.prevent="submitRequest" class="space-y-6">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Jenis Surat</label>
            <select 
              v-model="form.document_type" 
              class="block w-full border-gray-300 rounded-xl shadow-sm focus:ring-primary-500 focus:border-primary-500 py-3 px-4"
            >
              <option value="">Pilih Jenis Surat</option>
              <option value="surat_pengantar_ktp">Surat Pengantar KTP</option>
              <option value="surat_pengantar_kk">Surat Pengantar KK</option>
              <option value="surat_keterangan_domisili">Surat Keterangan Domisili</option>
              <option value="surat_izin_keramaian">Surat Izin Keramaian</option>
            </select>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Keperluan</label>
            <textarea 
              v-model="form.purpose" 
              rows="3" 
              class="block w-full border-gray-300 rounded-xl shadow-sm focus:ring-primary-500 focus:border-primary-500 py-3 px-4" 
              placeholder="Jelaskan keperluan pengajuan surat ini..."
            ></textarea>
          </div>

          <div class="flex justify-end pt-4 border-t border-gray-50">
            <button 
              type="button"
              @click="closeForm" 
              class="mr-3 px-4 py-2 border border-gray-300 rounded-xl text-gray-700 hover:bg-gray-50 font-medium"
            >
              Batal
            </button>
            <button 
              type="submit"
              :disabled="saving || !form.document_type || !form.purpose"
              class="px-6 py-2 bg-primary-600 text-white rounded-xl font-medium hover:bg-primary-700 shadow-lg shadow-primary-600/20 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {{ saving ? 'Mengirim...' : 'Kirim Pengajuan' }}
            </button>
          </div>
        </form>
      </div>

      <!-- Requests List -->
      <div class="space-y-4">
        <h3 class="text-lg font-bold text-gray-900">Riwayat Pengajuan</h3>
        
        <div v-if="!loading && requests.length === 0 && pagination.total === 0" class="text-center py-12 text-gray-500">
          <p>Belum ada pengajuan surat</p>
        </div>
        <template v-else>
          <div v-if="requests.length > 0" class="space-y-4 mb-4">
            <div 
              v-for="req in requests" 
              :key="req.id" 
              class="bg-white rounded-2xl border border-gray-100 shadow-sm p-6 hover:shadow-md transition-shadow"
            >
              <div class="flex justify-between items-start">
                <div>
                  <div class="flex items-center gap-3 mb-2">
                    <h3 class="font-bold text-gray-900">{{ getDocumentTypeLabel(req.document_type) }}</h3>
                    <span 
                      :class="[
                        'px-2.5 py-0.5 rounded-full text-xs font-bold uppercase tracking-wide',
                        req.status === 'completed' ? 'bg-green-100 text-green-700' : 
                        req.status === 'approved' ? 'bg-blue-100 text-blue-700' : 
                        req.status === 'rejected' ? 'bg-red-100 text-red-700' :
                        'bg-yellow-100 text-yellow-700'
                      ]"
                    >
                      {{ getStatusLabel(req.status) }}
                    </span>
                  </div>
                  <p class="text-gray-600 text-sm mb-3">{{ req.purpose }}</p>
                  <p class="text-xs text-gray-400">Diajukan pada {{ formatDate(req.created_at) }}</p>
                  <div v-if="req.rejected_reason" class="mt-2 p-2 bg-red-50 rounded-lg">
                    <p class="text-sm text-red-800"><strong>Alasan Ditolak:</strong> {{ req.rejected_reason }}</p>
                  </div>
                </div>
                <button 
                  v-if="req.status === 'completed'" 
                  class="text-primary-600 hover:text-primary-700 font-medium text-sm flex items-center"
                >
                  <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"></path></svg>
                  Download PDF
                </button>
              </div>
            </div>
          </div>
          
          <!-- Pagination -->
          <div v-if="pagination.total > 0 && pagination.total_pages > 0" class="flex justify-center items-center gap-2 mt-4">
            <button
              @click="pagination.page > 1 && (pagination.page--, loadRequests())"
              :disabled="pagination.page === 1"
              :class="['px-4 py-2 rounded-lg border transition-colors', pagination.page === 1 ? 'border-gray-200 text-gray-400 cursor-not-allowed bg-gray-50' : 'border-gray-300 text-gray-700 hover:bg-gray-50 bg-white']"
            >
              Sebelumnya
            </button>
            <span class="text-sm text-gray-600 px-4">
              Halaman {{ pagination.page }} dari {{ pagination.total_pages }}
            </span>
            <button
              @click="pagination.page < pagination.total_pages && (pagination.page++, loadRequests())"
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
useGlobalLoading(loading, 'Memuat pengajuan surat...')
const saving = ref(false)
const showForm = ref(false)
const requests = ref([])

const form = reactive({
  document_type: '',
  purpose: '',
})

const pagination = ref({
  page: 1,
  limit: 20,
  total: 0,
  total_pages: 0
})

const loadRequests = async () => {
  loading.value = true
  try {
    const response = await fetch(`/api/document-requests?page=${pagination.value.page}&limit=${pagination.value.limit}`)
    requests.value = response.requests || []
    if (response.pagination) {
      pagination.value = {
        page: response.pagination.page || pagination.value.page,
        limit: response.pagination.limit || pagination.value.limit,
        total: response.pagination.total || 0,
        total_pages: response.pagination.total_pages || 0
      }
    }
  } catch (error: any) {
    console.error('Failed to load requests:', error)
    showError(error.message || 'Gagal memuat pengajuan')
  } finally {
    loading.value = false
  }
}

const submitRequest = async () => {
  if (!form.document_type || !form.purpose) {
    showError('Jenis surat dan keperluan wajib diisi')
    return
  }

  saving.value = true
  try {
    await fetch('/api/document-requests', {
      method: 'POST',
      body: JSON.stringify({
        document_type: form.document_type,
        purpose: form.purpose,
      }),
    })
    showSuccess('Pengajuan surat berhasil dikirim!')
    closeForm()
    await loadRequests()
  } catch (error: any) {
    console.error('Failed to submit request:', error)
    showError(error.message || 'Gagal mengirim pengajuan')
  } finally {
    saving.value = false
  }
}

const closeForm = () => {
  showForm.value = false
  form.document_type = ''
  form.purpose = ''
}

const getDocumentTypeLabel = (type: string) => {
  const labels: Record<string, string> = {
    surat_pengantar_ktp: 'Surat Pengantar KTP',
    surat_pengantar_kk: 'Surat Pengantar KK',
    surat_keterangan_domisili: 'Surat Keterangan Domisili',
    surat_izin_keramaian: 'Surat Izin Keramaian',
  }
  return labels[type] || type
}

const getStatusLabel = (status: string) => {
  const labels: Record<string, string> = {
    pending: 'Menunggu',
    approved: 'Disetujui',
    rejected: 'Ditolak',
    completed: 'Selesai',
  }
  return labels[status] || status
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric'
  })
}

onMounted(() => {
  loadRequests()
})
</script>
