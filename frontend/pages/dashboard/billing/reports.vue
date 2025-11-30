<template>
  <div class="max-w-7xl mx-auto">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
      <div>
        <NuxtLink to="/dashboard/billing" class="text-sm text-gray-500 hover:text-gray-900 flex items-center mb-2">
          <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path></svg>
          Kembali ke Tagihan
        </NuxtLink>
        <h1 class="text-2xl font-bold text-gray-900 tracking-tight">Dashboard Keuangan</h1>
        <p class="text-gray-500 text-sm mt-1">Pantau koleksi, tunggakan, dan tren keuangan real-time</p>
      </div>
      <div class="flex gap-3">
        <select 
          v-model="selectedPeriod" 
          @change="loadDashboard"
          class="block w-full md:w-48 pl-3 pr-10 py-2.5 text-base border border-gray-200 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent sm:text-sm rounded-xl bg-white shadow-sm"
        >
          <option value="">Bulan Ini</option>
          <option v-for="month in availableMonths" :key="month.value" :value="month.value">
            {{ month.label }}
          </option>
        </select>
        <button
          @click="loadDashboard"
          :disabled="loading"
          class="px-4 py-2.5 bg-primary-600 text-white rounded-xl font-medium hover:bg-primary-700 transition-colors shadow-sm disabled:opacity-50"
        >
          <svg class="w-5 h-5" :class="{ 'animate-spin': loading }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
          </svg>
        </button>
      </div>
    </div>

    <!-- Key Metrics -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
      <div class="bg-white p-6 rounded-2xl border border-gray-100 shadow-sm">
        <p class="text-sm font-medium text-gray-500 mb-1">Tingkat Koleksi</p>
        <h3 class="text-3xl font-bold text-gray-900">{{ formatPercentage(dashboardData.summary?.collection_rate || 0) }}%</h3>
        <div class="mt-2 flex items-center text-sm" :class="getCollectionRateColor(dashboardData.summary?.collection_rate || 0)">
          <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"></path>
          </svg>
          {{ getCollectionRateText(dashboardData.summary?.collection_rate || 0) }}
        </div>
      </div>
      <div class="bg-white p-6 rounded-2xl border border-gray-100 shadow-sm">
        <p class="text-sm font-medium text-gray-500 mb-1">Total Belum Bayar</p>
        <h3 class="text-3xl font-bold text-gray-900">{{ formatCurrency(dashboardData.summary?.total_pending?.amount || 0) }}</h3>
        <p class="text-sm text-gray-500 mt-1">{{ dashboardData.summary?.total_pending?.count || 0 }} tagihan</p>
      </div>
      <div class="bg-white p-6 rounded-2xl border border-gray-100 shadow-sm">
        <p class="text-sm font-medium text-gray-500 mb-1">Total Lunas</p>
        <h3 class="text-3xl font-bold text-green-600">{{ formatCurrency(dashboardData.summary?.total_paid?.amount || 0) }}</h3>
        <p class="text-sm text-gray-500 mt-1">{{ dashboardData.summary?.total_paid?.count || 0 }} tagihan</p>
      </div>
      <div class="bg-white p-6 rounded-2xl border border-gray-100 shadow-sm">
        <p class="text-sm font-medium text-gray-500 mb-1">Total Terlambat</p>
        <h3 class="text-3xl font-bold text-red-600">{{ formatCurrency(dashboardData.summary?.total_overdue?.amount || 0) }}</h3>
        <p class="text-sm text-gray-500 mt-1">{{ dashboardData.summary?.total_overdue?.count || 0 }} tagihan</p>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-8">
      <!-- Monthly Trend Chart -->
      <div class="bg-white p-6 rounded-2xl border border-gray-100 shadow-sm overflow-visible">
        <h3 class="text-lg font-bold text-gray-900 mb-6">Tren Bulanan (12 Bulan Terakhir)</h3>
        <div v-if="loading" class="h-80 flex items-center justify-center">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
        </div>
        <div v-else-if="dashboardData.monthly_trend?.length === 0" class="h-80 flex items-center justify-center text-gray-400">
          <p>Tidak ada data</p>
        </div>
        <div v-else class="relative w-full overflow-visible" style="height: 360px; min-height: 360px; padding-bottom: 20px;">
          <svg class="w-full h-full" :viewBox="`0 0 ${svgWidth} ${svgHeight}`" preserveAspectRatio="xMidYMid meet" style="overflow: visible;">
            <defs>
              <linearGradient id="paidGradient" x1="0%" y1="0%" x2="0%" y2="100%">
                <stop offset="0%" style="stop-color:#3b82f6;stop-opacity:0.2" />
                <stop offset="100%" style="stop-color:#3b82f6;stop-opacity:0" />
              </linearGradient>
              <linearGradient id="pendingGradient" x1="0%" y1="0%" x2="0%" y2="100%">
                <stop offset="0%" style="stop-color:#eab308;stop-opacity:0.2" />
                <stop offset="100%" style="stop-color:#eab308;stop-opacity:0" />
              </linearGradient>
              <linearGradient id="overdueGradient" x1="0%" y1="0%" x2="0%" y2="100%">
                <stop offset="0%" style="stop-color:#ef4444;stop-opacity:0.2" />
                <stop offset="100%" style="stop-color:#ef4444;stop-opacity:0" />
              </linearGradient>
            </defs>
            
            <!-- Horizontal grid lines -->
            <g v-for="i in 5" :key="i">
              <line 
                :x1="chartPadding.left" 
                :y1="chartPadding.top + (i - 1) * (chartHeight / 4)" 
                :x2="chartPadding.left + chartWidth" 
                :y2="chartPadding.top + (i - 1) * (chartHeight / 4)" 
                stroke="#e5e7eb" 
                stroke-width="1" 
                stroke-dasharray="2,2"
              />
            </g>
            
            <!-- Area under curves (gradient fill) -->
            <path 
              v-if="paidAreaPath" 
              :d="paidAreaPath" 
              fill="url(#paidGradient)" 
              class="transition-opacity hover:opacity-80"
            />
            <path 
              v-if="pendingAreaPath" 
              :d="pendingAreaPath" 
              fill="url(#pendingGradient)" 
              class="transition-opacity hover:opacity-80"
            />
            <path 
              v-if="overdueAreaPath" 
              :d="overdueAreaPath" 
              fill="url(#overdueGradient)" 
              class="transition-opacity hover:opacity-80"
            />
            
            <!-- Line paths -->
            <path 
              v-if="paidLinePath" 
              :d="paidLinePath" 
              fill="none" 
              stroke="#3b82f6" 
              stroke-width="3" 
              stroke-linecap="round" 
              stroke-linejoin="round"
              class="transition-all"
            />
            <path 
              v-if="pendingLinePath" 
              :d="pendingLinePath" 
              fill="none" 
              stroke="#eab308" 
              stroke-width="3" 
              stroke-linecap="round" 
              stroke-linejoin="round"
              class="transition-all"
            />
            <path 
              v-if="overdueLinePath" 
              :d="overdueLinePath" 
              fill="none" 
              stroke="#ef4444" 
              stroke-width="3" 
              stroke-linecap="round" 
              stroke-linejoin="round"
              class="transition-all"
            />
            
            <!-- Data points -->
            <g v-for="(month, index) in dashboardData.monthly_trend" :key="index">
              <!-- Paid point -->
              <circle 
                v-if="(month.paid_amount || 0) > 0"
                :cx="chartPadding.left + (index * chartStep)" 
                :cy="getYPosition(month.paid_amount || 0, maxAmount)" 
                r="5" 
                fill="#3b82f6" 
                stroke="white" 
                stroke-width="2"
                class="cursor-pointer transition-all hover:r-7"
              >
                <title>Lunas: {{ formatCurrency(month.paid_amount || 0) }}</title>
              </circle>
              <!-- Pending point -->
              <circle 
                v-if="(month.pending_amount || 0) > 0"
                :cx="chartPadding.left + (index * chartStep)" 
                :cy="getYPosition(month.pending_amount || 0, maxAmount)" 
                r="5" 
                fill="#eab308" 
                stroke="white" 
                stroke-width="2"
                class="cursor-pointer transition-all hover:r-7"
              >
                <title>Belum Bayar: {{ formatCurrency(month.pending_amount || 0) }}</title>
              </circle>
              <!-- Overdue point -->
              <circle 
                v-if="(month.overdue_amount || 0) > 0"
                :cx="chartPadding.left + (index * chartStep)" 
                :cy="getYPosition(month.overdue_amount || 0, maxAmount)" 
                r="5" 
                fill="#ef4444" 
                stroke="white" 
                stroke-width="2"
                class="cursor-pointer transition-all hover:r-7"
              >
                <title>Terlambat: {{ formatCurrency(month.overdue_amount || 0) }}</title>
              </circle>
            </g>
          </svg>
          
          <!-- X-axis labels -->
          <div class="absolute bottom-0 left-0 right-0" style="height: 70px; padding-left: 90px; padding-right: 50px;">
            <div class="relative h-full w-full">
              <span 
                v-for="(month, index) in dashboardData.monthly_trend" 
                :key="index"
                class="absolute text-xs text-gray-500 whitespace-nowrap"
                :style="{ 
                  left: `${(index / Math.max(dashboardData.monthly_trend.length - 1, 1)) * 100}%`,
                  transform: 'translateX(-50%)',
                  bottom: '20px'
                }"
              >
                {{ formatMonthLabel(month.month) || 'N/A' }}
              </span>
            </div>
          </div>
          
          <!-- Y-axis labels -->
          <div class="absolute left-0 top-0" style="width: 85px; padding-left: 8px; padding-top: 20px; height: 180px; z-index: 10;">
            <div class="flex flex-col justify-between h-full text-xs text-gray-500 bg-white/80 backdrop-blur-sm rounded-r pr-2">
              <span class="truncate" :title="formatCurrency(maxAmount)">{{ formatCurrency(maxAmount) }}</span>
              <span class="truncate" :title="formatCurrency(maxAmount * 0.75)">{{ formatCurrency(maxAmount * 0.75) }}</span>
              <span class="truncate" :title="formatCurrency(maxAmount * 0.5)">{{ formatCurrency(maxAmount * 0.5) }}</span>
              <span class="truncate" :title="formatCurrency(maxAmount * 0.25)">{{ formatCurrency(maxAmount * 0.25) }}</span>
              <span>Rp 0</span>
            </div>
          </div>
        </div>
        <div class="flex gap-4 mt-6 justify-center">
          <div class="flex items-center text-sm text-gray-600">
            <div class="w-3 h-3 bg-primary-500 rounded-full mr-2"></div>Lunas
          </div>
          <div class="flex items-center text-sm text-gray-600">
            <div class="w-3 h-3 bg-yellow-500 rounded-full mr-2"></div>Belum Bayar
          </div>
          <div class="flex items-center text-sm text-gray-600">
            <div class="w-3 h-3 bg-red-500 rounded-full mr-2"></div>Terlambat
          </div>
        </div>
      </div>

      <!-- Top Overdue -->
      <div class="bg-white p-6 rounded-2xl border border-gray-100 shadow-sm">
        <div class="flex justify-between items-center mb-6">
          <h3 class="text-lg font-bold text-gray-900">Tunggakan Tertinggi</h3>
          <NuxtLink to="/dashboard/billing?status=overdue" class="text-sm text-primary-600 hover:text-primary-700 font-medium">
            Lihat Semua
          </NuxtLink>
        </div>
        <div v-if="loading" class="space-y-4">
          <div v-for="i in 5" :key="i" class="animate-pulse">
            <div class="h-16 bg-gray-100 rounded-xl"></div>
          </div>
        </div>
        <div v-else-if="dashboardData.top_overdue?.length === 0" class="text-center py-12 text-gray-400">
          <Icon name="heroicons:check-circle" class="w-12 h-12 mx-auto mb-3 opacity-50" />
          <p class="text-sm">Tidak ada tunggakan</p>
        </div>
        <div v-else class="space-y-4">
          <div 
            v-for="(unit, index) in dashboardData.top_overdue" 
            :key="index"
            class="flex items-center justify-between p-3 hover:bg-gray-50 rounded-xl transition-colors"
          >
            <div class="flex items-center">
              <div class="w-10 h-10 bg-red-50 text-red-600 rounded-lg flex items-center justify-center font-bold text-sm">
                {{ unit.unit_code }}
              </div>
              <div class="ml-3">
                <p class="text-sm font-medium text-gray-900">Unit {{ unit.unit_code }}</p>
                <p class="text-xs text-gray-500">{{ unit.overdue_count }} tagihan terlambat</p>
              </div>
            </div>
            <div class="text-right">
              <p class="text-sm font-bold text-gray-900">{{ formatCurrency(unit.total_amount) }}</p>
              <NuxtLink 
                :to="`/dashboard/billing?unit=${unit.unit_code}&status=overdue`"
                class="text-xs text-primary-600 hover:underline"
              >
                Detail
              </NuxtLink>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Summary Stats -->
    <div class="bg-white rounded-2xl border border-gray-100 shadow-sm p-6">
      <h3 class="text-lg font-bold text-gray-900 mb-6">Ringkasan</h3>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div>
          <p class="text-sm text-gray-500 mb-2">Total Tagihan</p>
          <p class="text-2xl font-bold text-gray-900">{{ formatCurrency(dashboardData.summary?.total_bills?.amount || 0) }}</p>
          <p class="text-sm text-gray-500 mt-1">{{ dashboardData.summary?.total_bills?.count || 0 }} tagihan</p>
        </div>
        <div>
          <p class="text-sm text-gray-500 mb-2">Rata-rata per Tagihan</p>
          <p class="text-2xl font-bold text-gray-900">
            {{ formatCurrency(getAverageBillAmount()) }}
          </p>
        </div>
        <div>
          <p class="text-sm text-gray-500 mb-2">Periode</p>
          <p class="text-2xl font-bold text-gray-900">{{ formatPeriodLabel(selectedPeriod) }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const { fetch } = useApi()
const { showError } = useToast()

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

const loading = ref(false)
const selectedPeriod = ref('')
const dashboardData = ref<any>({
  summary: {},
  monthly_trend: [],
  top_overdue: [],
  period: ''
})

// Generate available months (last 12 months)
const availableMonths = computed(() => {
  const months = []
  const now = new Date()
  for (let i = 0; i < 12; i++) {
    const date = new Date(now.getFullYear(), now.getMonth() - i, 1)
    const value = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`
    const label = date.toLocaleDateString('id-ID', { year: 'numeric', month: 'long' })
    months.push({ value, label })
  }
  return months
})

// Calculate max amount for chart scaling
const maxAmount = computed(() => {
  if (!dashboardData.value.monthly_trend?.length) return 1
  const amounts = dashboardData.value.monthly_trend.map((m: any) => 
    Math.max(m.paid_amount || 0, m.pending_amount || 0, m.overdue_amount || 0)
  )
  const max = Math.max(...amounts, 1)
  // Round up to nearest nice number for better chart display
  return Math.ceil(max / 100000) * 100000
})

// Chart dimensions and calculations
const svgWidth = 900
const svgHeight = 240
const chartWidth = 750
const chartHeight = 180
const chartPadding = { top: 20, right: 50, bottom: 60, left: 90 }

const chartStep = computed(() => {
  const dataLength = dashboardData.value.monthly_trend?.length || 1
  return dataLength > 1 ? chartWidth / (dataLength - 1) : chartWidth
})

// Get Y position for a value (inverted because SVG Y increases downward)
const getYPosition = (value: number, max: number) => {
  if (max === 0) return chartPadding.top + chartHeight
  const percentage = value / max
  return chartPadding.top + chartHeight - (percentage * chartHeight)
}

// Generate line path for paid amounts
const paidLinePath = computed(() => {
  if (!dashboardData.value.monthly_trend?.length) return ''
  const points = dashboardData.value.monthly_trend.map((m: any, index: number) => {
    const x = chartPadding.left + (index * chartStep.value)
    const y = getYPosition(m.paid_amount || 0, maxAmount.value)
    return `${index === 0 ? 'M' : 'L'} ${x} ${y}`
  })
  return points.join(' ')
})

// Generate area path for paid amounts (for gradient fill)
const paidAreaPath = computed(() => {
  if (!paidLinePath.value) return ''
  const firstX = chartPadding.left
  const lastX = chartPadding.left + ((dashboardData.value.monthly_trend?.length || 1) - 1) * chartStep.value
  const bottomY = chartHeight + chartPadding.top
  return `${paidLinePath.value} L ${lastX} ${bottomY} L ${firstX} ${bottomY} Z`
})

// Generate line path for pending amounts
const pendingLinePath = computed(() => {
  if (!dashboardData.value.monthly_trend?.length) return ''
  const points = dashboardData.value.monthly_trend.map((m: any, index: number) => {
    const x = chartPadding.left + (index * chartStep.value)
    const y = getYPosition(m.pending_amount || 0, maxAmount.value)
    return `${index === 0 ? 'M' : 'L'} ${x} ${y}`
  })
  return points.join(' ')
})

// Generate area path for pending amounts
const pendingAreaPath = computed(() => {
  if (!pendingLinePath.value) return ''
  const firstX = chartPadding.left
  const lastX = chartPadding.left + ((dashboardData.value.monthly_trend?.length || 1) - 1) * chartStep.value
  const bottomY = chartHeight + chartPadding.top
  return `${pendingLinePath.value} L ${lastX} ${bottomY} L ${firstX} ${bottomY} Z`
})

// Generate line path for overdue amounts
const overdueLinePath = computed(() => {
  if (!dashboardData.value.monthly_trend?.length) return ''
  const points = dashboardData.value.monthly_trend.map((m: any, index: number) => {
    const x = chartPadding.left + (index * chartStep.value)
    const y = getYPosition(m.overdue_amount || 0, maxAmount.value)
    return `${index === 0 ? 'M' : 'L'} ${x} ${y}`
  })
  return points.join(' ')
})

// Generate area path for overdue amounts
const overdueAreaPath = computed(() => {
  if (!overdueLinePath.value) return ''
  const firstX = chartPadding.left
  const lastX = chartPadding.left + ((dashboardData.value.monthly_trend?.length || 1) - 1) * chartStep.value
  const bottomY = chartHeight + chartPadding.top
  return `${overdueLinePath.value} L ${lastX} ${bottomY} L ${firstX} ${bottomY} Z`
})

const loadDashboard = async () => {
  loading.value = true
  try {
    const params = selectedPeriod.value ? `?period=${selectedPeriod.value}` : ''
    const response = await fetch(`/api/billing/dashboard${params}`)
    dashboardData.value = response
    
    // Debug: log response data
    console.log('Dashboard response:', response)
    if (response.monthly_trend) {
      console.log('Monthly trend data:', response.monthly_trend)
    }
    if (response.summary) {
      console.log('Summary data:', response.summary)
    }
  } catch (error: any) {
    console.error('Dashboard load error:', error)
    showError(error.message || 'Gagal memuat dashboard', 'Error')
  } finally {
    loading.value = false
  }
}

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(amount)
}

const formatPercentage = (value: number) => {
  return value.toFixed(1)
}

const formatMonthLabel = (monthStr: string) => {
  if (!monthStr) return ''
  
  // Try to parse YYYY-MM format
  const parts = monthStr.split('-')
  if (parts.length === 2) {
    const year = parseInt(parts[0])
    const month = parseInt(parts[1])
    
    // Validate year and month
    if (!isNaN(year) && !isNaN(month) && month >= 1 && month <= 12 && year > 2000 && year < 2100) {
      try {
        const date = new Date(year, month - 1, 1)
        // Check if date is valid
        if (date.getFullYear() === year && date.getMonth() === month - 1) {
          return date.toLocaleDateString('id-ID', { month: 'short' })
        }
      } catch (e) {
        console.warn('Invalid date format:', monthStr, e)
      }
    }
  }
  
  // Fallback: try to extract month name if format is different
  // If it's already a readable format, return first 3 chars
  if (monthStr.length > 3) {
    // Try to find month name in Indonesian
    const monthNames = ['Jan', 'Feb', 'Mar', 'Apr', 'Mei', 'Jun', 'Jul', 'Agu', 'Sep', 'Okt', 'Nov', 'Des']
    for (const name of monthNames) {
      if (monthStr.includes(name)) {
        return name
      }
    }
    // Return first 3 characters as fallback
    return monthStr.substring(0, 3)
  }
  
  return monthStr
}

const formatPeriodLabel = (period: string) => {
  if (!period) {
    return 'Semua Periode'
  }
  const [year, month] = period.split('-')
  if (month) {
    const date = new Date(parseInt(year), parseInt(month) - 1, 1)
    return date.toLocaleDateString('id-ID', { year: 'numeric', month: 'long' })
  }
  return `Tahun ${year}`
}


const getCollectionRateColor = (rate: number) => {
  if (rate >= 80) return 'text-green-600'
  if (rate >= 60) return 'text-yellow-600'
  return 'text-red-600'
}

const getCollectionRateText = (rate: number) => {
  if (rate >= 80) return 'Sangat Baik'
  if (rate >= 60) return 'Baik'
  return 'Perlu Perhatian'
}

const getAverageBillAmount = () => {
  const total = dashboardData.value.summary?.total_bills?.amount || 0
  const count = dashboardData.value.summary?.total_bills?.count || 0
  return count > 0 ? total / count : 0
}

// Auto-refresh every 30 seconds
let refreshInterval: NodeJS.Timeout | null = null

onMounted(() => {
  loadDashboard()
  refreshInterval = setInterval(() => {
    loadDashboard()
  }, 30000) // 30 seconds
})

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})
</script>
