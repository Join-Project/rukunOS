<template>
  <div class="max-w-5xl mx-auto space-y-8 animate-fade-in">
      <!-- Header Section -->
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 class="text-2xl font-bold text-gray-900 tracking-tight">Dashboard Warga</h1>
          <p class="text-gray-500">{{ greeting }}, {{ user?.full_name || 'Warga' }}</p>
        </div>
        <div class="flex items-center gap-3">
          <div v-if="summary.unit_code" class="px-3 py-1 bg-green-50 text-green-700 rounded-full text-sm font-medium border border-green-100 flex items-center shadow-sm">
            <span class="w-2 h-2 bg-green-500 rounded-full mr-2 animate-pulse"></span>
            {{ summary.unit_code }} (Aktif)
          </div>
          <button class="p-2 text-gray-400 hover:text-primary-600 transition-colors relative">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"></path></svg>
            <span class="absolute top-1.5 right-1.5 w-2.5 h-2.5 bg-red-500 rounded-full border-2 border-white"></span>
          </button>
        </div>
      </div>

      <!-- Main Fintech Card -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- Left Column: Financials -->
        <div class="lg:col-span-2 space-y-6">
          <!-- Balance/Bill Card -->
          <div v-if="summary.active_bill" class="bg-gray-900 rounded-3xl p-8 text-white relative overflow-hidden shadow-2xl shadow-primary-900/20 group hover:shadow-primary-900/30 transition-all duration-500">
            <div class="relative z-10">
              <div class="flex justify-between items-start mb-8">
                <div>
                  <p class="text-gray-400 text-sm font-medium mb-1">Total Tagihan Aktif</p>
                  <h2 class="text-4xl font-bold tracking-tight">Rp {{ formatCurrency(summary.total_pending_amount) }}</h2>
                </div>
                <div class="bg-white/10 p-2 rounded-lg backdrop-blur-sm border border-white/10">
                  <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path></svg>
                </div>
              </div>
              
              <div class="space-y-4">
                <div class="flex items-center justify-between p-4 bg-white/5 rounded-2xl border border-white/10 hover:bg-white/10 transition-colors cursor-pointer">
                  <div class="flex items-center gap-4">
                    <div class="w-10 h-10 rounded-full bg-blue-500/20 flex items-center justify-center text-blue-400">
                      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"></path></svg>
                    </div>
                    <div>
                      <p class="font-medium">{{ summary.active_bill.category }}</p>
                      <p class="text-sm text-gray-400">{{ summary.active_bill.period }}</p>
                    </div>
                  </div>
                  <div class="text-right">
                    <p class="font-bold">Rp {{ formatCurrency(summary.active_bill.amount) }}</p>
                    <p v-if="summary.active_bill.due_date" class="text-xs text-red-400 font-medium bg-red-400/10 px-2 py-0.5 rounded-full mt-1">Jatuh Tempo {{ formatDate(summary.active_bill.due_date) }}</p>
                  </div>
                </div>
              </div>

              <div class="mt-8 pt-6 border-t border-white/10 flex gap-4">
                <NuxtLink to="/dashboard/warga/bills" class="flex-1 bg-primary-600 hover:bg-primary-500 text-center text-white py-3 px-4 rounded-xl font-semibold transition-all shadow-lg shadow-primary-600/20 hover:shadow-primary-600/40 hover:-translate-y-0.5">
                  Bayar Sekarang
                </NuxtLink>
                <NuxtLink to="/dashboard/warga/bills" class="flex-1 bg-white/10 hover:bg-white/20 text-center text-white py-3 px-4 rounded-xl font-semibold transition-colors backdrop-blur-sm border border-white/5">
                  Riwayat
                </NuxtLink>
              </div>
            </div>
            
            <!-- Decorative Gradients -->
            <div class="absolute top-0 right-0 w-64 h-64 bg-primary-600 rounded-full blur-[80px] opacity-20 -translate-y-1/2 translate-x-1/2 group-hover:opacity-30 transition-opacity duration-500"></div>
            <div class="absolute bottom-0 left-0 w-64 h-64 bg-purple-600 rounded-full blur-[80px] opacity-20 translate-y-1/2 -translate-x-1/2 group-hover:opacity-30 transition-opacity duration-500"></div>
          </div>
          <div v-else class="bg-gray-900 rounded-3xl p-8 text-white relative overflow-hidden shadow-2xl shadow-primary-900/20">
            <div class="text-center py-8">
              <p class="text-gray-400 text-sm font-medium mb-2">Tidak ada tagihan aktif</p>
              <p class="text-2xl font-bold">Semua tagihan sudah lunas!</p>
            </div>
          </div>

          <!-- Quick Actions Grid -->
          <div>
            <h3 class="text-lg font-bold text-gray-900 mb-4">Aksi Cepat</h3>
            <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
              <NuxtLink to="/dashboard/warga/family" class="flex flex-col items-center justify-center p-6 bg-white rounded-2xl border border-gray-100 shadow-sm hover:shadow-md hover:-translate-y-0.5 transition-all group">
                <div class="w-12 h-12 bg-blue-50 rounded-xl flex items-center justify-center text-blue-600 mb-3 group-hover:bg-blue-100 transition-colors">
                  <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z"></path></svg>
                </div>
                <span class="text-sm font-medium text-gray-700 group-hover:text-gray-900">Data Keluarga</span>
              </NuxtLink>
              
              <NuxtLink to="/dashboard/warga/complaints" class="flex flex-col items-center justify-center p-6 bg-white rounded-2xl border border-gray-100 shadow-sm hover:shadow-md hover:-translate-y-0.5 transition-all group">
                <div class="w-12 h-12 bg-purple-50 rounded-xl flex items-center justify-center text-purple-600 mb-3 group-hover:bg-purple-100 transition-colors">
                  <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"></path></svg>
                </div>
                <span class="text-sm font-medium text-gray-700 group-hover:text-gray-900">Lapor Masalah</span>
              </NuxtLink>

              <NuxtLink to="/dashboard/warga/documents" class="flex flex-col items-center justify-center p-6 bg-white rounded-2xl border border-gray-100 shadow-sm hover:shadow-md hover:-translate-y-0.5 transition-all group">
                <div class="w-12 h-12 bg-yellow-50 rounded-xl flex items-center justify-center text-yellow-600 mb-3 group-hover:bg-yellow-100 transition-colors">
                  <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z"></path></svg>
                </div>
                <span class="text-sm font-medium text-gray-700 group-hover:text-gray-900">Surat Pengantar</span>
              </NuxtLink>

              <NuxtLink to="/dashboard/warga/panic" class="flex flex-col items-center justify-center p-6 bg-white rounded-2xl border border-gray-100 shadow-sm hover:shadow-md hover:-translate-y-0.5 transition-all group">
                <div class="w-12 h-12 bg-red-50 rounded-xl flex items-center justify-center text-red-600 mb-3 group-hover:bg-red-100 transition-colors">
                  <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path></svg>
                </div>
                <span class="text-sm font-medium text-gray-700 group-hover:text-gray-900">Darurat</span>
              </NuxtLink>
            </div>
          </div>

          <!-- Recent Activity -->
          <div>
            <h3 class="text-lg font-bold text-gray-900 mb-4">Aktivitas Terbaru</h3>
            <div v-if="!loading && summary.activities && summary.activities.length === 0" class="text-center py-12 text-gray-500">
              <p>Belum ada aktivitas</p>
            </div>
            <div v-else class="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden">
              <div 
                v-for="activity in summary.activities" 
                :key="activity.type + activity.date"
                class="p-4 border-b border-gray-50 last:border-0 flex items-center gap-4 hover:bg-gray-50 transition-colors"
              >
                <div :class="[
                  'w-10 h-10 rounded-full flex items-center justify-center',
                  activity.iconColor === 'green' ? 'bg-green-100 text-green-600' :
                  activity.iconColor === 'blue' ? 'bg-blue-100 text-blue-600' :
                  activity.iconColor === 'yellow' ? 'bg-yellow-100 text-yellow-600' :
                  'bg-gray-100 text-gray-600'
                ]">
                  <svg v-if="activity.icon === 'check'" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>
                  <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path></svg>
                </div>
                <div class="flex-1">
                  <p class="text-sm font-medium text-gray-900">{{ activity.title }}</p>
                  <p class="text-xs text-gray-500">{{ formatTimeAgo(activity.date) }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Right Column: Info & Events -->
        <div class="space-y-6">
          <div class="bg-white rounded-3xl p-6 border border-gray-100 shadow-sm">
            <h3 class="text-lg font-bold text-gray-900 mb-4">Pengumuman</h3>
            <div v-if="!loading && (!summary.announcements || summary.announcements.length === 0)" class="text-center py-8 text-gray-500 text-sm">
              <p>Tidak ada pengumuman</p>
            </div>
            <div v-else class="space-y-4">
              <div 
                v-for="announcement in summary.announcements" 
                :key="announcement.id"
                class="pb-4 border-b border-gray-50 last:border-0 last:pb-0 group cursor-pointer"
                @click="$router.push('/dashboard/warga/announcements')"
              >
                <div class="flex items-center justify-between mb-2">
                  <span :class="[
                    'px-2 py-1 text-xs font-bold rounded',
                    announcement.priority === 'high' ? 'bg-red-50 text-red-600' :
                    announcement.priority === 'medium' ? 'bg-yellow-50 text-yellow-600' :
                    'bg-blue-50 text-blue-600'
                  ]">
                    {{ announcement.priority === 'high' ? 'PENTING' : announcement.priority === 'medium' ? 'INFO' : 'UMUM' }}
                  </span>
                  <span class="text-xs text-gray-400">{{ formatTimeAgo(announcement.created_at) }}</span>
                </div>
                <h4 class="font-semibold text-gray-900 mb-1 group-hover:text-primary-600 transition-colors">{{ announcement.title }}</h4>
                <p class="text-sm text-gray-500 line-clamp-2">{{ announcement.content }}</p>
              </div>
            </div>
            <NuxtLink to="/dashboard/warga/announcements" class="block w-full mt-4 text-center text-sm font-medium text-primary-600 hover:text-primary-700 hover:underline">Lihat Semua</NuxtLink>
          </div>

          <div class="bg-gradient-to-br from-gray-50 to-white rounded-3xl p-6 border border-gray-100 shadow-sm">
            <h3 class="text-lg font-bold text-gray-900 mb-4">Kas Warga</h3>
            <div class="flex items-end gap-2 mb-2">
              <span class="text-3xl font-bold text-gray-900">Rp 15.4Jt</span>
              <span class="text-sm text-green-600 font-medium mb-1.5">▲ 12%</span>
            </div>
            <p class="text-sm text-gray-500 mb-4">Saldo kas RW saat ini</p>
            <div class="h-1.5 w-full bg-gray-100 rounded-full overflow-hidden">
              <div class="h-full bg-primary-500 w-[70%] rounded-full relative overflow-hidden">
                <div class="absolute inset-0 bg-white/20 animate-[shimmer_2s_infinite]"></div>
              </div>
            </div>
            <p class="text-xs text-gray-400 mt-2 text-right">Target: Rp 20Jt</p>
          </div>
        </div>
      </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useApi } from '~/composables/useApi'
import { useToast } from '~/composables/useToast'
import { useGlobalLoading } from '~/composables/useGlobalLoading'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

const { fetch } = useApi()
const { showError } = useToast()
const { user } = useAuth()

const loading = ref(false)

// Sync loading state dengan global loading spinner
useGlobalLoading(loading, 'Memuat dashboard...')
const summary = ref({
  unit_code: '',
  total_pending_amount: 0,
  pending_bills_count: 0,
  active_bill: null,
  announcements: [],
  activities: [],
})

const greeting = computed(() => {
  const hour = new Date().getHours()
  if (hour < 11) return 'Selamat Pagi'
  if (hour < 15) return 'Selamat Siang'
  if (hour < 18) return 'Selamat Sore'
  return 'Selamat Malam'
})

const loadSummary = async () => {
  loading.value = true
  try {
    const response = await fetch('/api/dashboard/warga/summary')
    summary.value = response
  } catch (error: any) {
    console.error('Failed to load dashboard summary:', error)
    showError(error.message || 'Gagal memuat data dashboard')
  } finally {
    loading.value = false
  }
}

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('id-ID').format(amount || 0)
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'short',
  })
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
  loadSummary()
})
</script>
