<template>
  <div class="min-h-screen bg-white font-sans selection:bg-primary-100 selection:text-primary-900">
    <!-- Navbar -->
    <nav class="fixed w-full z-50 bg-white/80 backdrop-blur-lg border-b border-gray-100 transition-all duration-300">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-20 items-center">
          <div class="flex items-center gap-2">
            <div class="w-8 h-8 bg-primary-600 rounded-lg flex items-center justify-center text-white font-bold">R</div>
            <span class="text-xl font-bold text-gray-900 tracking-tight">RukunOS</span>
          </div>
          <div class="hidden md:flex space-x-8">
            <a href="#features" class="text-sm font-medium text-gray-600 hover:text-primary-600 transition-colors">Fitur</a>
            <a href="#solutions" class="text-sm font-medium text-gray-600 hover:text-primary-600 transition-colors">Solusi</a>
            <a href="#pricing" class="text-sm font-medium text-gray-600 hover:text-primary-600 transition-colors">Harga</a>
          </div>
          <div class="flex items-center space-x-4">
            <NuxtLink to="/login" class="text-sm font-medium text-gray-900 hover:text-primary-600 transition-colors">
              Masuk
            </NuxtLink>
            <NuxtLink to="/register" class="bg-gray-900 text-white px-5 py-2.5 rounded-lg text-sm font-medium hover:bg-gray-800 transition-all shadow-lg shadow-gray-900/20 hover:shadow-gray-900/30 hover:-translate-y-0.5">
              Daftar Sekarang
            </NuxtLink>
          </div>
        </div>
      </div>
    </nav>

    <!-- Hero Section -->
    <section class="pt-32 pb-20 lg:pt-48 lg:pb-32 overflow-hidden relative">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
        <div class="text-center max-w-4xl mx-auto">
          <div class="inline-flex items-center px-4 py-1.5 rounded-full bg-primary-50 border border-primary-100 text-primary-700 text-sm font-medium mb-8 animate-fade-in">
            <span class="flex h-2 w-2 rounded-full bg-primary-600 mr-2 animate-pulse"></span>
            Platform Manajemen Komunitas Modern
          </div>
          <h1 class="text-5xl lg:text-7xl font-bold text-gray-900 tracking-tight mb-8 leading-[1.1] animate-slide-up">
            Kelola Lingkungan <br>
            <span class="text-transparent bg-clip-text bg-gradient-to-r from-primary-600 to-indigo-600">Tanpa Ribet.</span>
          </h1>
          <p class="text-xl text-gray-600 mb-10 leading-relaxed max-w-2xl mx-auto animate-slide-up" style="animation-delay: 0.1s">
            Sistem operasi lengkap untuk RT/RW dan perumahan. Otomatisasi tagihan, transparansi keuangan, dan komunikasi warga dalam satu aplikasi elegan.
          </p>
          <div class="flex flex-col sm:flex-row justify-center gap-4 animate-slide-up" style="animation-delay: 0.2s">
            <NuxtLink to="/register" class="px-8 py-4 bg-primary-600 text-white rounded-xl font-semibold text-lg hover:bg-primary-700 transition-all shadow-xl shadow-primary-500/30 hover:shadow-primary-500/40 hover:-translate-y-1">
              Mulai Gratis
            </NuxtLink>
            <a href="#demo" class="px-8 py-4 bg-white text-gray-700 border border-gray-200 rounded-xl font-semibold text-lg hover:bg-gray-50 transition-all hover:border-gray-300">
              Lihat Demo
            </a>
          </div>
        </div>
      </div>
      
      <!-- Background Gradients -->
      <div class="absolute top-0 left-1/2 -translate-x-1/2 -z-10 w-full h-full overflow-hidden pointer-events-none">
        <div class="absolute top-[-10%] left-[-10%] w-[50%] h-[50%] bg-blue-100/50 rounded-full blur-3xl opacity-60 mix-blend-multiply filter animate-blob"></div>
        <div class="absolute top-[-10%] right-[-10%] w-[50%] h-[50%] bg-indigo-100/50 rounded-full blur-3xl opacity-60 mix-blend-multiply filter animate-blob animation-delay-2000"></div>
        <div class="absolute bottom-[-20%] left-[20%] w-[60%] h-[60%] bg-purple-100/50 rounded-full blur-3xl opacity-60 mix-blend-multiply filter animate-blob animation-delay-4000"></div>
      </div>
    </section>

    <!-- Social Proof -->
    <section class="py-10 border-y border-gray-100 bg-gray-50/50">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
        <p class="text-sm font-medium text-gray-500 mb-6 uppercase tracking-wider">Dipercaya oleh 500+ Komunitas Modern</p>
        <div class="flex flex-wrap justify-center gap-8 md:gap-16 opacity-60 grayscale hover:grayscale-0 transition-all duration-500">
          <!-- Mock Logos -->
          <div class="h-8 font-bold text-xl text-gray-800 flex items-center gap-2"><div class="w-6 h-6 bg-gray-800 rounded"></div>Griya Asri</div>
          <div class="h-8 font-bold text-xl text-gray-800 flex items-center gap-2"><div class="w-6 h-6 bg-gray-800 rounded"></div>Royal Residence</div>
          <div class="h-8 font-bold text-xl text-gray-800 flex items-center gap-2"><div class="w-6 h-6 bg-gray-800 rounded"></div>Green Valley</div>
          <div class="h-8 font-bold text-xl text-gray-800 flex items-center gap-2"><div class="w-6 h-6 bg-gray-800 rounded"></div>Metro Park</div>
        </div>
      </div>
    </section>

    <!-- Features Section (Bento Grid) -->
    <section id="features" class="py-24 bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-16">
          <h2 class="text-3xl md:text-4xl font-bold text-gray-900 mb-4 tracking-tight">Satu Platform, Segala Kebutuhan</h2>
          <p class="text-gray-600 max-w-2xl mx-auto text-lg">
            Kami menggabungkan fitur finansial, administrasi, dan sosial dalam antarmuka yang indah.
          </p>
        </div>
        
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6 auto-rows-[300px]">
          <!-- Feature 1: Large Card -->
          <div class="md:col-span-2 row-span-1 bg-gray-50 rounded-3xl p-8 border border-gray-100 relative overflow-hidden group hover:shadow-lg transition-all duration-300">
            <div class="relative z-10 max-w-md">
              <div class="w-12 h-12 bg-white rounded-xl shadow-sm flex items-center justify-center mb-6 text-primary-600">
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
              </div>
              <h3 class="text-2xl font-bold text-gray-900 mb-3">Smart Billing System</h3>
              <p class="text-gray-600 text-lg">
                Tagihan IPL otomatis via WhatsApp. Integrasi QRIS untuk pembayaran instan. Rekonsiliasi otomatis tanpa pusing.
              </p>
            </div>
            <!-- Decorative UI Element -->
            <div class="absolute right-0 bottom-0 w-1/2 h-3/4 bg-white rounded-tl-3xl shadow-2xl border border-gray-100 transform translate-y-8 translate-x-8 group-hover:translate-x-4 transition-transform duration-500 p-6">
              <div class="space-y-4">
                <div class="flex justify-between items-center border-b border-gray-50 pb-2">
                  <div class="h-2 w-20 bg-gray-200 rounded"></div>
                  <div class="h-2 w-10 bg-green-200 rounded"></div>
                </div>
                <div class="flex justify-between items-center border-b border-gray-50 pb-2">
                  <div class="h-2 w-24 bg-gray-200 rounded"></div>
                  <div class="h-2 w-10 bg-green-200 rounded"></div>
                </div>
                <div class="flex justify-between items-center">
                  <div class="h-2 w-16 bg-gray-200 rounded"></div>
                  <div class="h-2 w-10 bg-red-200 rounded"></div>
                </div>
              </div>
            </div>
          </div>

          <!-- Feature 2: Tall Card -->
          <div class="md:col-span-1 row-span-2 bg-gray-900 rounded-3xl p-8 relative overflow-hidden group hover:shadow-xl transition-all duration-300 text-white">
            <div class="relative z-10">
              <div class="w-12 h-12 bg-white/10 rounded-xl backdrop-blur-sm flex items-center justify-center mb-6 text-white">
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5.882V19.24a1.76 1.76 0 01-3.417.592l-2.147-6.15M18 13a3 3 0 100-6M5.436 13.683A4.001 4.001 0 017 6h1.832c4.1 0 7.625-1.234 9.168-3v14c-1.543-1.766-5.067-3-9.168-3H7a3.988 3.988 0 01-1.564-.317z"></path></svg>
              </div>
              <h3 class="text-2xl font-bold mb-3">Pusat Informasi Digital</h3>
              <p class="text-gray-400 text-lg mb-8">
                Broadcast pengumuman penting langsung ke HP warga. Surat menyurat digital tanpa kertas.
              </p>
              <div class="space-y-3">
                <div class="bg-white/10 p-4 rounded-xl backdrop-blur-sm border border-white/5 hover:bg-white/20 transition-colors">
                  <div class="flex items-center gap-3 mb-2">
                    <div class="w-2 h-2 rounded-full bg-red-500"></div>
                    <span class="text-sm font-medium">Penting</span>
                  </div>
                  <p class="text-sm text-gray-300">Kerja Bakti Minggu Ini</p>
                </div>
                <div class="bg-white/10 p-4 rounded-xl backdrop-blur-sm border border-white/5 hover:bg-white/20 transition-colors">
                  <div class="flex items-center gap-3 mb-2">
                    <div class="w-2 h-2 rounded-full bg-blue-500"></div>
                    <span class="text-sm font-medium">Info</span>
                  </div>
                  <p class="text-sm text-gray-300">Jadwal Keamanan Baru</p>
                </div>
              </div>
            </div>
            <!-- Gradient Glow -->
            <div class="absolute bottom-0 right-0 w-64 h-64 bg-primary-600 rounded-full blur-[100px] opacity-20"></div>
          </div>

          <!-- Feature 3: Standard Card -->
          <div class="md:col-span-1 row-span-1 bg-white rounded-3xl p-8 border border-gray-100 shadow-sm hover:shadow-lg transition-all duration-300">
            <div class="w-12 h-12 bg-green-50 rounded-xl flex items-center justify-center mb-6 text-green-600">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"></path></svg>
            </div>
            <h3 class="text-xl font-bold text-gray-900 mb-2">Database Warga</h3>
            <p class="text-gray-600">
              Data penghuni, kendaraan, dan hewan peliharaan terpusat dan aman.
            </p>
          </div>

          <!-- Feature 4: Standard Card -->
          <div class="md:col-span-1 row-span-1 bg-white rounded-3xl p-8 border border-gray-100 shadow-sm hover:shadow-lg transition-all duration-300">
            <div class="w-12 h-12 bg-purple-50 rounded-xl flex items-center justify-center mb-6 text-purple-600">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"></path></svg>
            </div>
            <h3 class="text-xl font-bold text-gray-900 mb-2">Keamanan 24/7</h3>
            <p class="text-gray-600">
              Panic button, manajemen tamu (VMS), dan laporan insiden real-time.
            </p>
          </div>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="bg-white border-t border-gray-100 py-12">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex flex-col md:flex-row justify-between items-center">
          <div class="flex items-center gap-2 mb-4 md:mb-0">
            <div class="w-6 h-6 bg-gray-900 rounded flex items-center justify-center text-white text-xs font-bold">R</div>
            <span class="text-lg font-bold text-gray-900 tracking-tight">RukunOS</span>
          </div>
          <div class="text-sm text-gray-500">
            &copy; 2025 RukunOS. All rights reserved.
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: false
})
</script>

<style scoped>
.animate-blob {
  animation: blob 10s infinite;
}
.animation-delay-2000 {
  animation-delay: 2s;
}
.animation-delay-4000 {
  animation-delay: 4s;
}
@keyframes blob {
  0% { transform: translate(0px, 0px) scale(1); }
  33% { transform: translate(30px, -50px) scale(1.1); }
  66% { transform: translate(-20px, 20px) scale(0.9); }
  100% { transform: translate(0px, 0px) scale(1); }
}
</style>
