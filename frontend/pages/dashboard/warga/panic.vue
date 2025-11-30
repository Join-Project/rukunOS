<template>
  <div class="max-w-xl mx-auto text-center pt-8">
      <h1 class="text-3xl font-bold text-gray-900 tracking-tight mb-2">Tombol Darurat</h1>
      <p class="text-gray-500 mb-12">Tekan tombol di bawah ini hanya dalam keadaan darurat. Satpam dan pengurus akan segera diberitahu lokasi Anda.</p>

      <!-- Panic Button -->
      <div class="relative inline-block group">
        <div class="absolute inset-0 bg-red-500 rounded-full blur-xl opacity-20 group-hover:opacity-40 transition-opacity animate-pulse"></div>
        <button 
          @click="activatePanic"
          :disabled="isCountingDown || sending"
          :class="[
            'relative w-64 h-64 rounded-full bg-gradient-to-b from-red-500 to-red-600 text-white shadow-2xl shadow-red-600/50 flex flex-col items-center justify-center transform transition-transform border-8 border-red-100',
            (isCountingDown || sending) ? 'opacity-50 cursor-not-allowed' : 'active:scale-95 group-hover:scale-105'
          ]"
        >
          <svg class="w-20 h-20 mb-2 animate-bounce" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path></svg>
          <span class="text-2xl font-bold uppercase tracking-widest">PANIC</span>
          <span class="text-xs font-medium opacity-80 mt-1">Tekan untuk Bantuan</span>
        </button>
      </div>

      <!-- Emergency Contacts -->
      <div class="mt-16 grid grid-cols-2 gap-4">
        <a href="tel:110" class="flex items-center justify-center p-4 bg-white border border-gray-100 rounded-2xl shadow-sm hover:bg-gray-50 transition-colors">
          <div class="w-10 h-10 bg-blue-100 rounded-full flex items-center justify-center text-blue-600 mr-3">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"></path></svg>
          </div>
          <div class="text-left">
            <p class="font-bold text-gray-900">Polisi</p>
            <p class="text-xs text-gray-500">110</p>
          </div>
        </a>
        <a href="tel:113" class="flex items-center justify-center p-4 bg-white border border-gray-100 rounded-2xl shadow-sm hover:bg-gray-50 transition-colors">
          <div class="w-10 h-10 bg-red-100 rounded-full flex items-center justify-center text-red-600 mr-3">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 18.657A8 8 0 016.343 7.343S7 9 9 10c0-2 .5-5 2.986-7C14 5 16.09 5.777 17.656 7.343A7.975 7.975 0 0120 13a7.975 7.975 0 01-2.343 5.657z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.879 16.121A3 3 0 1012.015 11L11 14H9c0 .768.293 1.536.879 2.121z"></path></svg>
          </div>
          <div class="text-left">
            <p class="font-bold text-gray-900">Pemadam</p>
            <p class="text-xs text-gray-500">113</p>
          </div>
        </a>
      </div>
    </div>

    <!-- Confirmation Modal -->
    <div v-if="isCountingDown" class="fixed inset-0 z-50 flex items-center justify-center bg-red-600 bg-opacity-90 backdrop-blur-sm">
      <div class="text-center text-white p-8 max-w-md w-full">
        <div class="text-6xl font-bold mb-4 animate-ping">{{ countdown }}</div>
        <h2 class="text-2xl font-bold mb-2">Mengirim Sinyal Darurat</h2>
        <p class="mb-8 opacity-90">Bantuan akan segera dikirim ke lokasi Anda.</p>
        
        <button 
          @click="cancelPanic"
          class="w-full py-4 bg-white text-red-600 rounded-2xl font-bold text-lg hover:bg-gray-100 transition-colors shadow-lg"
        >
          BATALKAN
        </button>
      </div>
    </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

const { fetch } = useApi()
const { showSuccess, showError } = useToast()

const isCountingDown = ref(false)
const sending = ref(false)
const countdown = ref(5)
let timer: any = null

const activatePanic = async () => {
  isCountingDown.value = true
  countdown.value = 5
  
  timer = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0) {
      clearInterval(timer)
      sendPanicAlert()
    }
  }, 1000)
}

const sendPanicAlert = async () => {
  sending.value = true
  try {
    // Try to get location (with timeout)
    let location: string | null = null
    if (navigator.geolocation) {
      try {
        const position = await new Promise<GeolocationPosition>((resolve, reject) => {
          navigator.geolocation.getCurrentPosition(resolve, reject, {
            timeout: 5000,
            maximumAge: 0,
          })
        })
        location = JSON.stringify({
          lat: position.coords.latitude,
          lng: position.coords.longitude,
        })
      } catch (err) {
        console.log('Location access denied or failed:', err)
        // Continue without location
      }
    }

    const data: any = {}
    if (location) {
      data.location = location
    }

    await fetch('/api/panic-alerts', {
      method: 'POST',
      body: JSON.stringify(data),
    })

    showSuccess('Sinyal Darurat Terkirim! Satpam sedang menuju lokasi.', 'Darurat Terkirim')
    isCountingDown.value = false
  } catch (error: any) {
    console.error('Failed to send panic alert:', error)
    showError(error.message || 'Gagal mengirim sinyal darurat')
    isCountingDown.value = false
  } finally {
    sending.value = false
  }
}

const cancelPanic = () => {
  clearInterval(timer)
  isCountingDown.value = false
  countdown.value = 5
}
</script>
