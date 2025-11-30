<template>
  <NuxtLayout name="auth">
    <div class="space-y-6">
      <div class="text-center">
        <h3 class="text-xl font-bold text-gray-900">Masuk ke Akun Anda</h3>
        <p class="mt-1 text-sm text-gray-500">
          Atau <NuxtLink to="/register" class="font-medium text-primary-600 hover:text-primary-500">daftar akun baru</NuxtLink>
        </p>
      </div>

      <form class="space-y-6" @submit.prevent="handleLogin">
        <UiInput
              id="email"
              v-model="email"
          label="Email Address"
              type="email"
          placeholder="nama@email.com"
              required
          :error="errors.email"
            />

        <UiInput
              id="password"
              v-model="password"
          label="Password"
              type="password"
          placeholder="••••••••"
              required
          :error="errors.password"
        />

        <div class="flex items-center justify-between">
          <div class="flex items-center">
            <input
              id="remember-me"
              name="remember-me"
              type="checkbox"
              class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
            />
            <label for="remember-me" class="ml-2 block text-sm text-gray-900">
              Ingat saya
            </label>
          </div>

          <div class="text-sm">
            <a href="#" class="font-medium text-primary-600 hover:text-primary-500">
              Lupa password?
            </a>
          </div>
        </div>

        <div>
          <UiButton
            type="submit"
            class="w-full"
            :loading="loading"
          >
            Masuk
          </UiButton>
        </div>

        <div v-if="errorMessage" class="rounded-md bg-red-50 p-4">
          <div class="flex">
            <div class="flex-shrink-0">
              <svg class="h-5 w-5 text-red-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
              </svg>
            </div>
            <div class="ml-3">
              <h3 class="text-sm font-medium text-red-800">
                Login Gagal
              </h3>
              <div class="mt-2 text-sm text-red-700">
                <p>{{ errorMessage }}</p>
              </div>
            </div>
          </div>
        </div>
      </form>
    </div>
  </NuxtLayout>
</template>

<script setup lang="ts">
definePageMeta({
  layout: false,
  middleware: 'auth'
})

const { fetch } = useApi()
const authStore = useAuthStore()
const router = useRouter()

const email = ref('')
const password = ref('')
const loading = ref(false)
const errorMessage = ref('')
const errors = ref({
  email: '',
  password: ''
})

const handleLogin = async () => {
  loading.value = true
  errorMessage.value = ''
  errors.value = { email: '', password: '' }

  // Validation
  if (!email.value) {
    errors.value.email = 'Email harus diisi'
    loading.value = false
    return
  }
  if (!password.value) {
    errors.value.password = 'Password harus diisi'
    loading.value = false
    return
  }

  try {
    const response = await fetch<{ token: string; user: any; tenant_id: string }>('/api/auth/login', {
      method: 'POST',
      body: JSON.stringify({
        email: email.value,
        password: password.value
      })
    })

    // Set token
    authStore.setToken(response.token)

    // Map user data to match User interface
    const userData = {
      id: response.user.id,
      email: response.user.email,
      full_name: response.user.full_name,
      phone: response.user.phone,
      tenant_id: response.tenant_id,
      tenant_name: response.user.tenant_name,
      role_id: response.user.role_id,
      role_name: response.user.role_name,
      unit_id: response.user.unit_id,
      permissions: response.user.permissions || [],
    }
    authStore.setUser(userData)

    // Fetch current user to get full details with permissions
    // Note: We already have role_name from login response, so this is optional
    try {
      const currentUser = await authStore.fetchCurrentUser()
      // Ensure role_name is preserved
      if (currentUser && !currentUser.role_name && userData.role_name) {
        authStore.setUser({ ...currentUser, role_name: userData.role_name })
      }
    } catch (err) {
      console.warn('Failed to fetch current user details, using basic user data:', err)
    }

    // Redirect to dashboard
    await navigateTo('/dashboard')
  } catch (error: any) {
    console.error('Login error:', error)
    const errorMsg = error.data?.error || error.message || 'Terjadi kesalahan saat login. Silakan coba lagi.'
    errorMessage.value = errorMsg
  } finally {
    loading.value = false
  }
}
</script>
