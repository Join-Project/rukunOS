<template>
  <div class="max-w-4xl mx-auto space-y-6">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-bold text-gray-900 tracking-tight">Pengumuman</h1>
          <p class="text-gray-500">Informasi terbaru seputar lingkungan RT/RW</p>
        </div>
      </div>

      <!-- Search & Filter -->
      <div class="flex gap-4">
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
          v-model="filterCategory" 
          @change="loadAnnouncements" 
          class="rounded-xl border-gray-200 focus:border-primary-500 focus:ring-primary-500"
        >
          <option value="">Semua Kategori</option>
          <option value="Penting">Penting</option>
          <option value="Kegiatan">Kegiatan</option>
          <option value="Info">Info</option>
        </select>
      </div>

      <!-- Announcement List -->
      <div v-if="!loading && announcements.length === 0" class="text-center py-12 text-gray-500">Belum ada pengumuman</div>
      <div v-else class="space-y-4">
        <!-- Pinned/Important -->
        <div 
          v-for="item in announcements" 
          :key="item.id"
          :class="[
            'rounded-2xl p-6 shadow-sm hover:shadow-md transition-shadow',
            item.is_pinned ? 'bg-white border-l-4 border-red-500' : 'bg-white border border-gray-100',
            item.priority === 'high' ? 'border-l-4 border-red-500' : ''
          ]"
        >
          <div class="flex items-start justify-between mb-2">
            <div class="flex items-center gap-2">
              <span 
                :class="[
                  'px-2 py-1 text-xs font-bold rounded',
                  item.priority === 'high' ? 'bg-red-50 text-red-600' : 
                  item.priority === 'medium' ? 'bg-yellow-50 text-yellow-600' : 
                  'bg-blue-50 text-blue-600'
                ]"
              >
                {{ item.priority === 'high' ? 'PENTING' : item.priority === 'medium' ? 'PENTING' : 'INFO' }}
              </span>
              <span class="text-xs text-gray-400">{{ formatDate(item.created_at) }}</span>
            </div>
            <svg v-if="!item.is_read" class="w-5 h-5 text-primary-500" fill="currentColor" viewBox="0 0 20 20"><path d="M10 2a6 6 0 00-6 6v3.586l-.707.707A1 1 0 004 14h12a1 1 0 00.707-1.707L16 11.586V8a6 6 0 00-6-6zM10 18a3 3 0 01-3-3h6a3 3 0 01-3 3z"></path></svg>
          </div>
          <h3 class="text-xl font-bold text-gray-900 mb-2">{{ item.title }}</h3>
          <p class="text-gray-600 mb-4 whitespace-pre-wrap">{{ item.content }}</p>
          <div v-if="item.category || item.author_name" class="flex items-center gap-4 text-sm text-gray-500 border-t border-gray-50 pt-4">
            <span v-if="item.category" class="flex items-center gap-1">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z"></path></svg>
              {{ item.category }}
            </span>
            <span v-if="item.author_name" class="flex items-center gap-1">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path></svg>
              {{ item.author_name }}
            </span>
          </div>
        </div>
      </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useApi } from '~/composables/useApi'
import { useToast } from '~/composables/useToast'
import { useGlobalLoading } from '~/composables/useGlobalLoading'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

const { fetch } = useApi()
const { showError } = useToast()

const loading = ref(false)

// Sync loading state dengan global loading spinner
useGlobalLoading(loading, 'Memuat pengumuman...')
const announcements = ref([])
const searchQuery = ref('')
const filterCategory = ref('')

const loadAnnouncements = async () => {
  loading.value = true
  try {
    const params = new URLSearchParams({
      page: '1',
      limit: '50',
    })
    if (searchQuery.value) {
      params.append('search', searchQuery.value)
    }
    if (filterCategory.value) {
      params.append('category', filterCategory.value)
    }

    const response = await fetch(`/api/announcements?${params.toString()}`)
    announcements.value = response.announcements || []
  } catch (error: any) {
    console.error('Failed to load announcements:', error)
    showError(error.message || 'Gagal memuat pengumuman')
  } finally {
    loading.value = false
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

onMounted(() => {
  loadAnnouncements()
})
</script>
