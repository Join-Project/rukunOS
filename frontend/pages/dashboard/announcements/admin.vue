<template>
  <div class="max-w-5xl mx-auto">
      <!-- Header -->
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
        <div>
          <h1 class="text-2xl font-bold text-gray-900 tracking-tight">Pengumuman & Informasi</h1>
          <p class="text-gray-500">Kelola pengumuman untuk warga. Kirim via WhatsApp & Aplikasi.</p>
        </div>
        <button @click="showCreateForm = true" class="px-4 py-2 bg-primary-600 text-white rounded-xl font-medium hover:bg-primary-700 transition-colors shadow-lg shadow-primary-600/20 flex items-center">
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path></svg>
          Buat Pengumuman
        </button>
      </div>

      <!-- Create Form (Collapsible) -->
      <div v-if="showCreateForm" class="bg-white rounded-2xl border border-gray-100 shadow-sm p-6 mb-8 animate-fade-in">
        <div class="flex justify-between items-center mb-6">
          <h3 class="text-lg font-bold text-gray-900">{{ editingAnnouncement ? 'Edit Pengumuman' : 'Buat Pengumuman Baru' }}</h3>
          <button @click="closeForm" class="text-gray-400 hover:text-gray-600">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
          </button>
        </div>
        
        <div class="space-y-6">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Judul Pengumuman *</label>
            <input 
              v-model="form.value.title" 
              type="text" 
              class="block w-full border-gray-300 rounded-xl shadow-sm focus:ring-primary-500 focus:border-primary-500 py-3 px-4" 
              placeholder="Contoh: Kerja Bakti Minggu Ini"
            >
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Isi Pengumuman *</label>
            <textarea 
              v-model="form.value.content" 
              rows="4" 
              class="block w-full border-gray-300 rounded-xl shadow-sm focus:ring-primary-500 focus:border-primary-500 py-3 px-4" 
              placeholder="Tulis detail pengumuman di sini..."
            ></textarea>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Prioritas</label>
              <select v-model="form.value.priority" class="block w-full border-gray-300 rounded-xl shadow-sm focus:ring-primary-500 focus:border-primary-500 py-3 px-4">
                <option value="low">Normal (Info Biasa)</option>
                <option value="medium">Penting (Perlu Perhatian)</option>
                <option value="high">Darurat / Sangat Penting</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Kategori</label>
              <input 
                v-model="form.value.category" 
                type="text" 
                class="block w-full border-gray-300 rounded-xl shadow-sm focus:ring-primary-500 focus:border-primary-500 py-3 px-4" 
                placeholder="Contoh: Kegiatan, Info, Penting"
              >
            </div>
            <div>
              <label class="flex items-center">
                <input 
                  v-model="form.value.is_pinned" 
                  type="checkbox" 
                  class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
                >
                <span class="ml-2 text-sm text-gray-700">Pin pengumuman (tampilkan di atas)</span>
              </label>
            </div>
          </div>

          <div class="flex justify-end pt-4 border-t border-gray-50">
            <button @click="closeForm" class="mr-3 px-4 py-2 border border-gray-300 rounded-xl text-gray-700 hover:bg-gray-50 font-medium">Batal</button>
            <button 
              @click.stop="saveAnnouncement" 
              type="button"
              :disabled="saving || !form.value.title || !form.value.content"
              :class="['px-6 py-2 bg-primary-600 text-white rounded-xl font-medium hover:bg-primary-700 shadow-lg shadow-primary-600/20 transition-all', (saving || !form.value.title || !form.value.content) ? 'opacity-50 cursor-not-allowed' : '']"
            >
              {{ saving ? 'Menyimpan...' : (editingAnnouncement ? 'Update Pengumuman' : 'Kirim Pengumuman') }}
            </button>
          </div>
        </div>
      </div>

      <!-- Search & Filter -->
      <div class="flex gap-4 mb-6">
        <div class="relative flex-1">
          <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path></svg>
          <input 
            v-model="searchQuery" 
            @input="loadAnnouncements" 
            type="text" 
            placeholder="Cari pengumuman..." 
            class="w-full pl-10 pr-4 py-2.5 rounded-xl border-gray-200 focus:border-primary-500 focus:ring-primary-500"
          >
        </div>
        <select 
          v-model="filterPriority" 
          @change="loadAnnouncements" 
          class="rounded-xl border-gray-200 focus:border-primary-500 focus:ring-primary-500"
        >
          <option value="">Semua Prioritas</option>
          <option value="high">Darurat</option>
          <option value="medium">Penting</option>
          <option value="low">Normal</option>
        </select>
      </div>

      <!-- Announcements List -->
      <div v-if="!loading && announcements.length === 0" class="text-center py-12 text-gray-500">Belum ada pengumuman. Buat yang pertama!</div>
      <div v-else class="space-y-4">
        <div v-for="item in announcements" :key="item.id" class="bg-white rounded-2xl border border-gray-100 shadow-sm p-6 hover:shadow-md transition-shadow relative overflow-hidden">
          <div :class="['absolute left-0 top-0 bottom-0 w-1.5', priorityColor(item.priority)]"></div>
          
          <div class="flex justify-between items-start mb-2 pl-4">
            <div class="flex items-center gap-3">
              <span v-if="item.is_pinned" class="px-2 py-0.5 bg-yellow-100 text-yellow-700 text-xs font-bold rounded">PINNED</span>
              <span :class="['px-2.5 py-0.5 rounded-full text-xs font-bold uppercase tracking-wide', priorityBadge(item.priority)]">
                {{ item.priority === 'high' ? 'Darurat' : item.priority === 'medium' ? 'Penting' : 'Normal' }}
              </span>
              <span class="text-sm text-gray-500">{{ formatDate(item.created_at) }}</span>
            </div>
            <div class="flex items-center gap-2">
              <button @click="editAnnouncement(item)" class="text-gray-400 hover:text-primary-600 transition-colors">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path></svg>
              </button>
              <button @click="deleteAnnouncement(item)" class="text-gray-400 hover:text-red-600 transition-colors">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
              </button>
            </div>
          </div>

          <div class="pl-4">
            <h3 class="text-lg font-bold text-gray-900 mb-2">{{ item.title }}</h3>
            <p class="text-gray-600 mb-4 whitespace-pre-wrap">{{ item.content }}</p>
            
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <span v-if="item.category" class="px-2 py-1 bg-gray-100 text-gray-600 text-xs rounded">{{ item.category }}</span>
                <span v-if="item.author_name" class="text-xs text-gray-500">Oleh: {{ item.author_name }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div v-if="pagination.total_pages > 1" class="flex justify-center items-center gap-2 mt-8">
        <button 
          @click="pagination.page > 1 && (pagination.page--, loadAnnouncements())"
          :disabled="pagination.page === 1"
          :class="['px-4 py-2 rounded-lg border', pagination.page === 1 ? 'border-gray-200 text-gray-400 cursor-not-allowed' : 'border-gray-300 text-gray-700 hover:bg-gray-50']"
        >
          Sebelumnya
        </button>
        <span class="px-4 py-2 text-sm text-gray-600">
          Halaman {{ pagination.page }} dari {{ pagination.total_pages }}
        </span>
        <button 
          @click="pagination.page < pagination.total_pages && (pagination.page++, loadAnnouncements())"
          :disabled="pagination.page === pagination.total_pages"
          :class="['px-4 py-2 rounded-lg border', pagination.page === pagination.total_pages ? 'border-gray-200 text-gray-400 cursor-not-allowed' : 'border-gray-300 text-gray-700 hover:bg-gray-50']"
        >
          Selanjutnya
        </button>
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
const confirm = useConfirm()

const showCreateForm = ref(false)
const loading = ref(false)

// Sync loading state dengan global loading spinner
useGlobalLoading(loading, 'Memuat pengumuman...')
const saving = ref(false)
const announcements = ref([])
const editingAnnouncement = ref(null)
const searchQuery = ref('')
const filterPriority = ref('')

const pagination = ref({
  page: 1,
  limit: 20,
  total: 0,
  total_pages: 0,
})

const form = ref({
  title: '',
  content: '',
  priority: 'medium',
  category: '',
  is_pinned: false,
})

const loadAnnouncements = async () => {
  loading.value = true
  try {
    const params = new URLSearchParams({
      page: pagination.value.page.toString(),
      limit: pagination.value.limit.toString(),
    })
    if (searchQuery.value) {
      params.append('search', searchQuery.value)
    }
    if (filterPriority.value) {
      params.append('priority', filterPriority.value)
    }

    console.log('Loading announcements with params:', params.toString())
    const response = await fetch(`/api/announcements?${params.toString()}`)
    console.log('Announcements response:', response)
    announcements.value = response.announcements || []
    pagination.value = {
      page: response.pagination?.page || 1,
      limit: response.pagination?.limit || 20,
      total: response.pagination?.total || 0,
      total_pages: response.pagination?.total_pages || 0,
    }
  } catch (error: any) {
    console.error('Failed to load announcements:', error)
    const errorMsg = error.data?.error || error.message || 'Gagal memuat pengumuman'
    showError(errorMsg)
  } finally {
    loading.value = false
  }
}

const saveAnnouncement = async (e?: Event) => {
  if (e) {
    e.preventDefault()
    e.stopPropagation()
  }
  
  console.log('=== saveAnnouncement called ===')
  console.log('Form value:', form.value)
  console.log('Form title:', form.value.title)
  console.log('Form content:', form.value.content)
  
  if (!form.value.title || !form.value.content) {
    console.log('Validation failed:', { title: form.value.title, content: form.value.content })
    showError('Judul dan isi pengumuman wajib diisi')
    return
  }

  console.log('Starting save process...')
  saving.value = true
  
  try {
    const data: any = {
      title: form.value.title,
      content: form.value.content,
      priority: form.value.priority,
      is_pinned: form.value.is_pinned,
    }
    if (form.value.category) {
      data.category = form.value.category
    }

    console.log('Saving announcement with data:', data)

    if (editingAnnouncement.value) {
      console.log('Updating announcement:', editingAnnouncement.value.id)
      const response = await fetch(`/api/announcements/${editingAnnouncement.value.id}`, {
        method: 'PUT',
        body: JSON.stringify(data),
      })
      console.log('Update response:', response)
      showSuccess('Pengumuman berhasil diperbarui')
    } else {
      console.log('Creating new announcement...')
      const response = await fetch('/api/announcements', {
        method: 'POST',
        body: JSON.stringify(data),
      })
      console.log('Create response:', response)
      showSuccess('Pengumuman berhasil dibuat')
    }

    console.log('Save successful, closing form and reloading...')
    closeForm()
    await loadAnnouncements()
  } catch (error: any) {
    console.error('Failed to save announcement:', error)
    console.error('Error details:', {
      message: error.message,
      data: error.data,
      status: error.status
    })
    const errorMsg = error.data?.error || error.message || 'Gagal menyimpan pengumuman'
    showError(errorMsg)
  } finally {
    saving.value = false
    console.log('Save process completed')
  }
}

const editAnnouncement = (announcement: any) => {
  editingAnnouncement.value = announcement
  form.value = {
    title: announcement.title,
    content: announcement.content,
    priority: announcement.priority,
    category: announcement.category || '',
    is_pinned: announcement.is_pinned || false,
  }
  showCreateForm.value = true
}

const deleteAnnouncement = async (announcement: any) => {
  const result = await confirm.show(
    `Yakin ingin menghapus pengumuman "${announcement.title}"?`,
    {
      title: 'Hapus Pengumuman',
      confirmText: 'Ya, Hapus',
      cancelText: 'Batal',
      type: 'danger'
    }
  )

  if (!result) return

  try {
    await fetch(`/api/announcements/${announcement.id}`, {
      method: 'DELETE',
    })
    showSuccess('Pengumuman berhasil dihapus')
    await loadAnnouncements()
  } catch (error: any) {
    console.error('Failed to delete announcement:', error)
    showError(error.message || 'Gagal menghapus pengumuman')
  }
}

const closeForm = () => {
  showCreateForm.value = false
  editingAnnouncement.value = null
  form.value = {
    title: '',
    content: '',
    priority: 'medium',
    category: '',
    is_pinned: false,
  }
}

const formatDate = (dateStr: string) => {
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
    month: 'long',
    year: 'numeric'
  })
}

const priorityColor = (p: string) => {
  if (p === 'high') return 'bg-red-500'
  if (p === 'medium') return 'bg-yellow-500'
  return 'bg-blue-500'
}

const priorityBadge = (p: string) => {
  if (p === 'high') return 'bg-red-100 text-red-700'
  if (p === 'medium') return 'bg-yellow-100 text-yellow-700'
  return 'bg-blue-100 text-blue-700'
}

onMounted(() => {
  loadAnnouncements()
})
</script>
