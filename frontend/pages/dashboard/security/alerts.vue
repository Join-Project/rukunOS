<template>
  <div class="max-w-6xl mx-auto space-y-8">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 tracking-tight">Panic Alerts Monitoring</h1>
        <p class="text-gray-500">Pantau sinyal darurat dari warga secara real-time</p>
      </div>
      <div v-if="activeAlertsCount > 0" class="flex items-center gap-2 px-4 py-2 bg-red-50 text-red-700 rounded-xl border border-red-100 animate-pulse">
        <span class="w-2 h-2 bg-red-600 rounded-full"></span>
        <span class="text-sm font-bold">{{ activeAlertsCount }} Alert Aktif</span>
      </div>
      <div v-else class="flex items-center gap-2 px-4 py-2 bg-green-50 text-green-700 rounded-xl border border-green-100">
        <span class="w-2 h-2 bg-green-600 rounded-full"></span>
        <span class="text-sm font-bold">Tidak Ada Alert Aktif</span>
      </div>
    </div>

    <!-- Active Alerts -->
    <div v-if="!loading && activeAlerts.length > 0" class="space-y-4">
      <h3 class="text-lg font-bold text-gray-900">Alert Aktif (Belum Ditangani)</h3>
      
      <!-- Alert Card -->
      <div 
        v-for="alert in activeAlerts" 
        :key="alert.id"
        class="bg-red-50 border-l-4 border-red-600 rounded-r-2xl p-6 shadow-sm flex items-center justify-between animate-[pulse_2s_infinite]"
      >
        <div class="flex items-center gap-6">
          <div class="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center text-red-600 animate-bounce">
            <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path></svg>
          </div>
          <div>
            <h4 class="text-xl font-bold text-red-700">DARURAT{{ alert.unit_code ? ': ' + alert.unit_code : '' }}</h4>
            <p class="text-red-600 font-medium">{{ alert.user_name || 'Warga' }}{{ alert.user_phone ? ' - ' + alert.user_phone : '' }}</p>
            <p class="text-sm text-red-500 mt-1">Diterima: {{ formatTimeAgo(alert.created_at) }}</p>
          </div>
        </div>
        <button 
          @click="handleAlert(alert.id)"
          :disabled="handlingAlert === alert.id"
          :class="['bg-red-600 hover:bg-red-700 text-white px-6 py-3 rounded-xl font-bold shadow-lg shadow-red-600/30 transition-all transform hover:scale-105', handlingAlert === alert.id ? 'opacity-50 cursor-not-allowed' : '']"
        >
          {{ handlingAlert === alert.id ? 'Menangani...' : 'Tangani Sekarang' }}
        </button>
      </div>
    </div>
    <div v-else-if="!loading" class="text-center py-12 text-gray-500">
      <p>Tidak ada alert aktif saat ini</p>
    </div>

    <!-- Alert History -->
    <div class="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden mt-8">
      <div class="p-6 border-b border-gray-100 flex justify-between items-center">
        <h3 class="text-lg font-bold text-gray-900">Riwayat Alert</h3>
        <select 
          v-model="filterStatus" 
          @change="loadAlerts" 
          class="rounded-lg border-gray-200 text-sm focus:border-primary-500 focus:ring-primary-500"
        >
          <option value="">Semua Status</option>
          <option value="active">Aktif</option>
          <option value="responded">Ditangani</option>
          <option value="resolved">Selesai</option>
        </select>
      </div>
      <div v-if="!loading && allAlerts.length === 0" class="p-12 text-center text-gray-500">Belum ada riwayat alert</div>
      <div v-else class="overflow-x-auto">
        <table class="w-full text-left text-sm">
          <thead class="bg-gray-50 text-gray-500 font-medium">
            <tr>
              <th class="px-6 py-4">Waktu</th>
              <th class="px-6 py-4">Unit / Warga</th>
              <th class="px-6 py-4">Status</th>
              <th class="px-6 py-4">Petugas</th>
              <th class="px-6 py-4 text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            <tr v-for="alert in allAlerts" :key="alert.id" class="hover:bg-gray-50 transition-colors">
              <td class="px-6 py-4 text-gray-500">{{ formatDateTime(alert.created_at) }}</td>
              <td class="px-6 py-4">
                <div class="font-medium text-gray-900">{{ alert.unit_code || '-' }}</div>
                <div class="text-xs text-gray-500">{{ alert.user_name || '-' }}{{ alert.user_phone ? ' - ' + alert.user_phone : '' }}</div>
              </td>
              <td class="px-6 py-4">
                <span 
                  :class="[
                    'px-2 py-1 rounded-full text-xs font-bold',
                    alert.status === 'active' ? 'bg-red-50 text-red-700' :
                    alert.status === 'responded' ? 'bg-yellow-50 text-yellow-700' :
                    'bg-green-50 text-green-700'
                  ]"
                >
                  {{ alert.status === 'active' ? 'Aktif' : alert.status === 'responded' ? 'Ditangani' : 'Selesai' }}
                </span>
              </td>
              <td class="px-6 py-4">{{ alert.responder_name || '-' }}</td>
              <td class="px-6 py-4 text-right">
                <div class="flex items-center justify-end gap-2">
                  <button 
                    v-if="alert.status === 'active'"
                    @click="handleAlert(alert.id)"
                    :disabled="handlingAlert === alert.id"
                    class="text-red-600 hover:text-red-700 font-medium text-xs border border-red-200 hover:bg-red-50 px-3 py-1.5 rounded-lg transition-colors"
                  >
                    {{ handlingAlert === alert.id ? 'Menangani...' : 'Tangani' }}
                  </button>
                  <button 
                    v-if="alert.status === 'responded'"
                    @click="resolveAlert(alert.id)"
                    :disabled="resolvingAlert === alert.id"
                    class="text-green-600 hover:text-green-700 font-medium text-xs border border-green-200 hover:bg-green-50 px-3 py-1.5 rounded-lg transition-colors"
                  >
                    {{ resolvingAlert === alert.id ? 'Menyelesaikan...' : 'Selesaikan' }}
                  </button>
                  <button 
                    @click="viewAlertDetails(alert)"
                    class="text-gray-400 hover:text-primary-600 transition-colors"
                  >
                    Detail
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Alert Detail Modal -->
    <div v-if="selectedAlert" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50" @click.self="selectedAlert = null">
      <div class="bg-white rounded-2xl p-6 max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto">
        <div class="flex justify-between items-center mb-6">
          <h3 class="text-xl font-bold text-gray-900">Detail Alert</h3>
          <button @click="selectedAlert = null" class="text-gray-400 hover:text-gray-600">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
          </button>
        </div>
        
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Status</label>
            <span 
              :class="[
                'inline-block px-3 py-1 rounded-full text-sm font-bold',
                selectedAlert.status === 'active' ? 'bg-red-50 text-red-700' :
                selectedAlert.status === 'responded' ? 'bg-yellow-50 text-yellow-700' :
                'bg-green-50 text-green-700'
              ]"
            >
              {{ selectedAlert.status === 'active' ? 'Aktif' : selectedAlert.status === 'responded' ? 'Ditangani' : 'Selesai' }}
            </span>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Warga</label>
            <p class="text-gray-900">{{ selectedAlert.user_name || '-' }}</p>
            <p class="text-sm text-gray-500">{{ selectedAlert.user_email || '' }}</p>
            <p class="text-sm text-gray-500">{{ selectedAlert.user_phone || '' }}</p>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Unit</label>
            <p class="text-gray-900">{{ selectedAlert.unit_code || '-' }}</p>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Waktu Alert</label>
            <p class="text-gray-900">{{ formatDateTime(selectedAlert.created_at) }}</p>
          </div>
          
          <div v-if="selectedAlert.responded_at">
            <label class="block text-sm font-medium text-gray-700 mb-1">Ditangani Pada</label>
            <p class="text-gray-900">{{ formatDateTime(selectedAlert.responded_at) }}</p>
            <p class="text-sm text-gray-500">Oleh: {{ selectedAlert.responder_name || '-' }}</p>
          </div>
          
          <div v-if="selectedAlert.resolved_at">
            <label class="block text-sm font-medium text-gray-700 mb-1">Selesai Pada</label>
            <p class="text-gray-900">{{ formatDateTime(selectedAlert.resolved_at) }}</p>
          </div>
          
          <div v-if="selectedAlert.location">
            <label class="block text-sm font-medium text-gray-700 mb-1">Lokasi</label>
            <p class="text-gray-900 text-xs font-mono">{{ selectedAlert.location }}</p>
          </div>
          
          <div v-if="selectedAlert.notes">
            <label class="block text-sm font-medium text-gray-700 mb-1">Catatan</label>
            <p class="text-gray-900 whitespace-pre-wrap">{{ selectedAlert.notes }}</p>
          </div>
        </div>
        
        <div class="mt-6 flex justify-end gap-3">
          <button 
            v-if="selectedAlert.status === 'active'"
            @click="handleAlert(selectedAlert.id); selectedAlert = null"
            :disabled="handlingAlert === selectedAlert.id"
            class="px-4 py-2 bg-red-600 text-white rounded-lg font-medium hover:bg-red-700 transition-colors"
          >
            Tangani Sekarang
          </button>
          <button 
            v-if="selectedAlert.status === 'responded'"
            @click="resolveAlert(selectedAlert.id); selectedAlert = null"
            :disabled="resolvingAlert === selectedAlert.id"
            class="px-4 py-2 bg-green-600 text-white rounded-lg font-medium hover:bg-green-700 transition-colors"
          >
            Tandai Selesai
          </button>
          <button 
            @click="selectedAlert = null"
            class="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg font-medium hover:bg-gray-50 transition-colors"
          >
            Tutup
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useGlobalLoading } from '~/composables/useGlobalLoading'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

const { fetch } = useApi()
const { showSuccess, showError } = useToast()

const loading = ref(false)

// Sync loading state dengan global loading spinner
useGlobalLoading(loading, 'Memuat data alert...')
const allAlerts = ref([])
const filterStatus = ref('')
const handlingAlert = ref<string | null>(null)
const resolvingAlert = ref<string | null>(null)
const selectedAlert = ref<any>(null)

let refreshInterval: any = null

const activeAlerts = computed(() => {
  return allAlerts.value.filter((alert: any) => alert.status === 'active')
})

const activeAlertsCount = computed(() => {
  return activeAlerts.value.length
})

const loadAlerts = async () => {
  loading.value = true
  try {
    const params = new URLSearchParams({
      page: '1',
      limit: '100',
    })
    if (filterStatus.value) {
      params.append('status', filterStatus.value)
    }

    const response = await fetch(`/api/panic-alerts?${params.toString()}`)
    allAlerts.value = response.alerts || []
  } catch (error: any) {
    console.error('Failed to load alerts:', error)
    showError(error.message || 'Gagal memuat alert')
  } finally {
    loading.value = false
  }
}

const handleAlert = async (alertId: string) => {
  handlingAlert.value = alertId
  try {
    await fetch(`/api/panic-alerts/${alertId}`, {
      method: 'PUT',
      body: JSON.stringify({
        status: 'responded'
      }),
    })
    showSuccess('Alert berhasil ditandai sebagai ditangani')
    await loadAlerts()
    if (selectedAlert.value && selectedAlert.value.id === alertId) {
      // Refresh selected alert
      const updated = allAlerts.value.find((a: any) => a.id === alertId)
      if (updated) {
        selectedAlert.value = updated
      }
    }
  } catch (error: any) {
    console.error('Failed to handle alert:', error)
    showError(error.message || 'Gagal menangani alert')
  } finally {
    handlingAlert.value = null
  }
}

const resolveAlert = async (alertId: string) => {
  resolvingAlert.value = alertId
  try {
    await fetch(`/api/panic-alerts/${alertId}`, {
      method: 'PUT',
      body: JSON.stringify({
        status: 'resolved'
      }),
    })
    showSuccess('Alert berhasil ditandai sebagai selesai')
    await loadAlerts()
    if (selectedAlert.value && selectedAlert.value.id === alertId) {
      // Refresh selected alert
      const updated = allAlerts.value.find((a: any) => a.id === alertId)
      if (updated) {
        selectedAlert.value = updated
      }
    }
  } catch (error: any) {
    console.error('Failed to resolve alert:', error)
    showError(error.message || 'Gagal menyelesaikan alert')
  } finally {
    resolvingAlert.value = null
  }
}

const viewAlertDetails = (alert: any) => {
  selectedAlert.value = alert
}

const formatDateTime = (dateStr: string) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const minutes = Math.floor(diff / (1000 * 60))
  const hours = Math.floor(diff / (1000 * 60 * 60))
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  
  if (minutes < 1) return 'Baru saja'
  if (minutes < 60) return `${minutes} menit yang lalu`
  if (hours < 24) return `${hours} jam yang lalu`
  if (days === 1) return 'Kemarin'
  if (days < 7) return `${days} hari yang lalu`
  
  return date.toLocaleString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

const formatTimeAgo = (dateStr: string) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const seconds = Math.floor(diff / 1000)
  const minutes = Math.floor(diff / (1000 * 60))
  const hours = Math.floor(diff / (1000 * 60 * 60))
  
  if (seconds < 60) return `Baru saja (${seconds} detik yang lalu)`
  if (minutes < 60) return `${minutes} menit yang lalu`
  if (hours < 24) return `${hours} jam yang lalu`
  
  return date.toLocaleString('id-ID', {
    day: 'numeric',
    month: 'short',
    hour: '2-digit',
    minute: '2-digit',
  })
}

onMounted(() => {
  loadAlerts()
  // Auto-refresh setiap 10 detik untuk monitoring real-time
  refreshInterval = setInterval(() => {
    loadAlerts()
  }, 10000)
})

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})
</script>
