<template>
  <div class="min-h-screen flex bg-white">
    <!-- Left Side: Branding & Visuals (Hidden on mobile) -->
    <div class="hidden lg:flex lg:w-1/2 relative overflow-hidden bg-gray-900">
      <!-- Abstract Background -->
      <div class="absolute inset-0">
         <div class="absolute top-0 -left-4 w-[800px] h-[800px] bg-primary-900 rounded-full mix-blend-multiply filter blur-3xl opacity-50 animate-blob"></div>
         <div class="absolute top-0 -right-4 w-[800px] h-[800px] bg-blue-900 rounded-full mix-blend-multiply filter blur-3xl opacity-50 animate-blob animation-delay-2000"></div>
         <div class="absolute -bottom-8 left-20 w-[800px] h-[800px] bg-indigo-900 rounded-full mix-blend-multiply filter blur-3xl opacity-50 animate-blob animation-delay-4000"></div>
         <div class="absolute inset-0 bg-gradient-to-t from-gray-900 via-transparent to-transparent opacity-80"></div>
      </div>

      <!-- Content Overlay -->
      <div class="relative z-10 flex flex-col justify-between w-full p-12 text-white">
        <div>
          <span class="text-2xl font-bold tracking-tight text-white">RukunOS</span>
        </div>
        
        <div class="space-y-6 max-w-lg">
          <blockquote class="text-2xl font-medium leading-relaxed">
            "Platform ini benar-benar mengubah cara kami mengelola lingkungan. Transparan, modern, dan sangat mudah digunakan oleh semua warga."
          </blockquote>
          <div class="flex items-center space-x-4">
            <div class="h-12 w-12 rounded-full bg-white/20 flex items-center justify-center text-lg font-bold">
              BS
            </div>
            <div>
              <div class="font-semibold">Budi Santoso</div>
              <div class="text-gray-400 text-sm">Ketua RW 05, Jakarta Selatan</div>
            </div>
          </div>
        </div>

        <div class="text-sm text-gray-500">
          &copy; 2025 RukunOS. All rights reserved.
        </div>
      </div>
    </div>

    <!-- Right Side: Form -->
    <div class="flex-1 flex flex-col justify-center py-12 px-4 sm:px-6 lg:px-20 xl:px-24 bg-white">
      <div class="mx-auto w-full max-w-sm lg:w-96">
        <div class="lg:hidden mb-10">
          <span class="text-2xl font-bold tracking-tight text-primary-900">RukunOS</span>
        </div>
        
        <slot />
      </div>
    </div>
  </div>
</template>

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
