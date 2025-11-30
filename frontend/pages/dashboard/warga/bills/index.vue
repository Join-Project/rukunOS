<template>
  <div class="max-w-3xl mx-auto">
      <h1 class="text-2xl font-bold text-gray-900 tracking-tight mb-6">Tagihan Saya</h1>

      <!-- Active Bill Card -->
      <div v-if="activeBill" class="bg-white rounded-3xl border border-gray-100 shadow-lg overflow-hidden mb-8 relative">
        <div class="absolute top-0 right-0 w-64 h-64 bg-primary-50 rounded-full -translate-y-1/2 translate-x-1/2 opacity-50"></div>
        
        <div class="p-8 relative z-10">
          <div class="flex justify-between items-start mb-6">
            <div>
              <span :class="[
                'px-3 py-1 rounded-full text-xs font-bold uppercase tracking-wide',
                activeBill.status === 'pending' ? 'bg-red-100 text-red-700' :
                activeBill.status === 'paid' ? 'bg-green-100 text-green-700' :
                'bg-yellow-100 text-yellow-700'
              ]">
                {{ activeBill.status === 'pending' ? 'Belum Dibayar' : activeBill.status === 'paid' ? 'Lunas' : 'Menunggu Konfirmasi' }}
              </span>
              <h2 class="text-3xl font-bold text-gray-900 mt-3">Rp {{ formatCurrency(activeBill.amount) }}</h2>
              <p class="text-gray-500">{{ activeBill.category }} - {{ activeBill.period }}</p>
            </div>
            <div class="text-right" v-if="activeBill.due_date">
              <p class="text-sm text-gray-500">Jatuh Tempo</p>
              <p class="font-bold text-gray-900">{{ formatDate(activeBill.due_date) }}</p>
            </div>
            <div class="text-right" v-else>
              <p class="text-sm text-gray-500">Tidak ada jatuh tempo</p>
            </div>
          </div>

          <div class="space-y-3 mb-8">
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">Tagihan Dasar</span>
              <span class="font-medium">Rp {{ formatCurrency(activeBill.amount) }}</span>
            </div>
            <div class="flex justify-between text-sm" v-if="activeBill.late_fee && activeBill.late_fee > 0">
              <span class="text-gray-600">Denda Keterlambatan</span>
              <span class="font-medium text-red-600">Rp {{ formatCurrency(activeBill.late_fee) }}</span>
            </div>
            <div class="h-px bg-gray-100 my-2"></div>
            <div class="flex justify-between text-base font-bold">
              <span class="text-gray-900">Total Pembayaran</span>
              <span class="text-primary-600">Rp {{ formatCurrency(activeBill.amount + (activeBill.late_fee || 0)) }}</span>
            </div>
          </div>

          <button 
            v-if="activeBill.status === 'pending'"
            @click="showPaymentModal = true" 
            class="w-full py-4 bg-primary-600 text-white rounded-xl font-bold text-lg hover:bg-primary-700 transition-all shadow-xl shadow-primary-600/20 hover:shadow-primary-600/30 hover:-translate-y-0.5"
          >
            Bayar Sekarang
          </button>
        </div>
      </div>
      <div v-else-if="!loading && !activeBill" class="text-center py-12 text-gray-500">
        <p>Tidak ada tagihan aktif</p>
      </div>

      <!-- History -->
      <h3 class="text-lg font-bold text-gray-900 mb-4">Riwayat Pembayaran</h3>
      <div v-if="!loading && history.length === 0 && pagination.total === 0" class="text-center py-12 text-gray-500">
        <p>Belum ada riwayat pembayaran</p>
      </div>
      <template v-else>
        <div v-if="history.length > 0" class="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden mb-4">
          <div v-for="bill in history" :key="bill.id" class="p-4 border-b border-gray-50 last:border-0 hover:bg-gray-50 transition-colors flex items-center justify-between">
            <div class="flex items-center">
              <div :class="[
                'w-10 h-10 rounded-full flex items-center justify-center mr-4',
                bill.status === 'paid' ? 'bg-green-100 text-green-600' : 'bg-gray-100 text-gray-600'
              ]">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>
              </div>
              <div>
                <p class="font-medium text-gray-900">{{ bill.category }} - {{ bill.period }}</p>
                <p class="text-xs text-gray-500" v-if="bill.paid_at">Dibayar pada {{ formatDateTime(bill.paid_at) }}</p>
                <p class="text-xs text-gray-500" v-else-if="bill.due_date">Jatuh tempo {{ formatDate(bill.due_date) }}</p>
              </div>
            </div>
            <div class="text-right">
              <p class="font-bold text-gray-900">Rp {{ formatCurrency(bill.amount) }}</p>
              <span :class="[
                'text-xs font-medium',
                bill.status === 'paid' ? 'text-green-600' : 'text-red-600'
              ]">
                {{ bill.status === 'paid' ? 'Lunas' : 'Belum Dibayar' }}
              </span>
            </div>
          </div>
        </div>
        
        <!-- Pagination -->
        <div v-if="pagination.total > 0 && pagination.total_pages > 0" class="flex justify-center items-center gap-2 mt-4">
          <button
            @click="pagination.page > 1 && (pagination.page--, loadBills())"
            :disabled="pagination.page === 1"
            :class="['px-4 py-2 rounded-lg border transition-colors', pagination.page === 1 ? 'border-gray-200 text-gray-400 cursor-not-allowed bg-gray-50' : 'border-gray-300 text-gray-700 hover:bg-gray-50 bg-white']"
          >
            Sebelumnya
          </button>
          <span class="text-sm text-gray-600 px-4">
            Halaman {{ pagination.page }} dari {{ pagination.total_pages }}
          </span>
          <button
            @click="pagination.page < pagination.total_pages && (pagination.page++, loadBills())"
            :disabled="pagination.page >= pagination.total_pages"
            :class="['px-4 py-2 rounded-lg border transition-colors', pagination.page >= pagination.total_pages ? 'border-gray-200 text-gray-400 cursor-not-allowed bg-gray-50' : 'border-gray-300 text-gray-700 hover:bg-gray-50 bg-white']"
          >
            Selanjutnya
          </button>
        </div>
      </template>
    </div>

    <!-- Payment Modal -->
    <div v-if="showPaymentModal" class="fixed inset-0 z-50 overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
      <div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div class="fixed inset-0 bg-gray-900 bg-opacity-75 transition-opacity backdrop-blur-sm" @click="showPaymentModal = false"></div>
        <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>
        <div class="inline-block align-bottom bg-white rounded-3xl text-left overflow-hidden shadow-2xl transform transition-all sm:my-8 sm:align-middle sm:max-w-md sm:w-full">
          <div class="bg-white p-6">
            <div class="flex justify-between items-center mb-6">
              <h3 class="text-xl font-bold text-gray-900">Pembayaran</h3>
              <button @click="showPaymentModal = false" class="text-gray-400 hover:text-gray-600">
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
              </button>
            </div>

            <div class="text-center mb-8">
              <p class="text-gray-500 mb-1">Total Tagihan</p>
              <p class="text-3xl font-bold text-gray-900">Rp {{ formatCurrency(activeBill ? (activeBill.amount + (activeBill.late_fee || 0)) : 0) }}</p>
            </div>

            <!-- QRIS -->
            <div class="bg-gray-50 p-6 rounded-2xl border-2 border-dashed border-gray-200 flex flex-col items-center justify-center mb-6">
              <div class="w-48 h-48 bg-white p-2 rounded-lg shadow-sm mb-4">
                <!-- Mock QR Code -->
                <div class="w-full h-full bg-gray-900 flex items-center justify-center text-white text-xs">
                  [QRIS CODE IMAGE]
                </div>
              </div>
              <p class="text-sm font-medium text-gray-900">Scan QRIS untuk membayar</p>
              <p class="text-xs text-gray-500 mt-1">Mendukung GoPay, OVO, Dana, BCA, dll</p>
            </div>

            <div class="space-y-3">
              <button class="w-full py-3 bg-gray-100 text-gray-700 rounded-xl font-medium hover:bg-gray-200 transition-colors">
                Salin Nomor Virtual Account
              </button>
              <button @click="simulatePayment" class="w-full py-3 bg-green-600 text-white rounded-xl font-bold hover:bg-green-700 transition-colors shadow-lg shadow-green-600/20">
                Simulasi Bayar Berhasil
              </button>
            </div>
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
const { showSuccess, showError } = useToast()

const loading = ref(false)

// Sync loading state dengan global loading spinner
useGlobalLoading(loading, 'Memuat tagihan...')
const bills = ref([])
const showPaymentModal = ref(false)

const activeBill = ref(null)
const history = ref([])

const pagination = ref({
  page: 1,
  limit: 20,
  total: 0,
  total_pages: 0
})

const loadBills = async () => {
  loading.value = true
  try {
    // Load pending bills for activeBill (limit to 1, should only be one active)
    const pendingResponse = await fetch('/api/billing?status=pending&limit=1')
    if (pendingResponse.bills && pendingResponse.bills.length > 0) {
      activeBill.value = pendingResponse.bills[0]
    } else {
      activeBill.value = null
    }

    // Load history with pagination (all bills, then filter for history)
    const historyResponse = await fetch(`/api/billing?page=${pagination.value.page}&limit=${pagination.value.limit}`)
    // Filter for history (paid, overdue, cancelled) and exclude pending
    const allBills = historyResponse.bills || []
    history.value = allBills
      .filter((bill: any) => bill.status !== 'pending')
      .sort((a: any, b: any) => {
        // Sort by paid_at if available, otherwise by created_at
        const dateA = a.paid_at ? new Date(a.paid_at).getTime() : new Date(a.created_at).getTime()
        const dateB = b.paid_at ? new Date(b.paid_at).getTime() : new Date(b.created_at).getTime()
        return dateB - dateA
      })
    
    if (historyResponse.pagination) {
      pagination.value = {
        page: historyResponse.pagination.page || pagination.value.page,
        limit: historyResponse.pagination.limit || pagination.value.limit,
        total: historyResponse.pagination.total || 0,
        total_pages: historyResponse.pagination.total_pages || 0
      }
    }
  } catch (error: any) {
    console.error('Failed to load bills:', error)
    showError(error.message || 'Gagal memuat tagihan')
  } finally {
    loading.value = false
  }
}

const simulatePayment = async () => {
  if (!activeBill.value) {
    showError('Tidak ada tagihan yang dipilih')
    return
  }

  try {
    await fetch(`/api/billing/${activeBill.value.id}/payment`, {
      method: 'POST',
      body: JSON.stringify({
        payment_method: 'qris',
        payment_reference: 'SIM-' + Date.now(),
      }),
    })
    showSuccess('Pembayaran Berhasil! Terima kasih.')
    showPaymentModal.value = false
    await loadBills()
  } catch (error: any) {
    console.error('Failed to process payment:', error)
    showError(error.message || 'Gagal memproses pembayaran')
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
    year: 'numeric'
  })
}

const formatDateTime = (dateStr: string) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

onMounted(() => {
  loadBills()
})
</script>
