<template>
  <NuxtLayout name="auth">
    <div class="space-y-6">
      <div class="text-center">
        <h3 class="text-xl font-bold text-gray-900">Buat Akun Baru</h3>
        <p class="mt-1 text-sm text-gray-500">
          Sudah punya akun? <NuxtLink to="/login" class="font-medium text-primary-600 hover:text-primary-500">Masuk disini</NuxtLink>
        </p>
      </div>

      <!-- Registration Type Selection -->
      <div class="space-y-4">
        <div class="text-sm font-medium text-gray-700 mb-3">Saya ingin:</div>
        <div class="grid grid-cols-2 gap-4">
          <button
            @click="registrationType = 'tenant'"
            :class="[
              'p-4 border-2 rounded-xl text-left transition-all',
              registrationType === 'tenant'
                ? 'border-primary-500 bg-primary-50'
                : 'border-gray-200 hover:border-gray-300'
            ]"
          >
            <div class="flex items-center space-x-2 mb-2">
              <div :class="[
                'w-5 h-5 rounded-full border-2 flex items-center justify-center',
                registrationType === 'tenant' ? 'border-primary-500 bg-primary-500' : 'border-gray-300'
              ]">
                <svg v-if="registrationType === 'tenant'" class="w-3 h-3 text-white" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                </svg>
              </div>
              <span class="font-semibold text-gray-900">Daftar sebagai Tenant</span>
            </div>
            <p class="text-xs text-gray-600">Membuat komunitas baru (RT/RW/Perumahan)</p>
          </button>

          <button
            @click="registrationType = 'warga'"
            :class="[
              'p-4 border-2 rounded-xl text-left transition-all',
              registrationType === 'warga'
                ? 'border-primary-500 bg-primary-50'
                : 'border-gray-200 hover:border-gray-300'
            ]"
          >
            <div class="flex items-center space-x-2 mb-2">
              <div :class="[
                'w-5 h-5 rounded-full border-2 flex items-center justify-center',
                registrationType === 'warga' ? 'border-primary-500 bg-primary-500' : 'border-gray-300'
              ]">
                <svg v-if="registrationType === 'warga'" class="w-3 h-3 text-white" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                </svg>
              </div>
              <span class="font-semibold text-gray-900">Daftar sebagai Warga</span>
            </div>
            <p class="text-xs text-gray-600">Bergabung dengan komunitas yang sudah ada</p>
          </button>
        </div>
      </div>

      <form class="space-y-6" @submit.prevent="handleRegister">
        <UiInput
          id="fullname"
          v-model="fullname"
          label="Nama Lengkap"
          type="text"
          placeholder="Nama Lengkap Anda"
          required
          :error="errors.fullname"
        />

        <UiInput
          id="email"
          v-model="email"
          label="Email Address"
          type="email"
          placeholder="nama@email.com"
          required
          :error="errors.email"
        />

        <!-- Tenant Registration Fields -->
        <template v-if="registrationType === 'tenant'">
          <UiInput
            id="tenant_name"
            v-model="tenantName"
            label="Nama Komunitas"
            type="text"
            placeholder="Contoh: RT 001 RW 05"
            required
            :error="errors.tenant_name"
          />

          <UiInput
            id="tenant_code"
            v-model="tenantCode"
            label="Kode Komunitas"
            type="text"
            placeholder="RT001 (unik, akan digunakan untuk invite warga)"
            required
            :error="errors.tenant_code"
            help="Kode ini akan digunakan warga untuk bergabung ke komunitas Anda"
          />

          <UiInput
            id="tenant_address"
            v-model="tenantAddress"
            label="Alamat Komunitas"
            type="text"
            placeholder="Alamat lengkap komunitas"
            :error="errors.tenant_address"
          />
        </template>

        <!-- Warga Registration Fields -->
        <template v-if="registrationType === 'warga'">
          <UiInput
            id="tenant_code_join"
            v-model="tenantCodeJoin"
            label="Kode Komunitas"
            type="text"
            placeholder="Masukkan kode komunitas yang diberikan"
            required
            :error="errors.tenant_code_join"
            help="Minta kode komunitas kepada pengurus RT/RW Anda"
          />
        </template>

        <UiInput
          id="password"
          v-model="password"
          label="Password"
          type="password"
          placeholder="••••••••"
          required
          :error="errors.password"
          help="Minimal 8 karakter"
        />

        <UiInput
          id="confirmPassword"
          v-model="confirmPassword"
          label="Konfirmasi Password"
          type="password"
          placeholder="••••••••"
          required
          :error="errors.confirmPassword"
        />

        <div class="flex items-center">
          <input
            id="terms"
            name="terms"
            type="checkbox"
            required
            class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
          />
          <label for="terms" class="ml-2 block text-sm text-gray-900">
            Saya setuju dengan <a href="#" class="text-primary-600 hover:text-primary-500">Syarat & Ketentuan</a>
          </label>
        </div>

        <div>
          <UiButton
            type="submit"
            class="w-full"
            :loading="loading"
          >
            {{ registrationType === 'tenant' ? 'Daftar sebagai Tenant' : 'Daftar sebagai Warga' }}
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
                Registrasi Gagal
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

const registrationType = ref<'tenant' | 'warga'>('warga')
const fullname = ref('')
const email = ref('')
const tenantName = ref('')
const tenantCode = ref('')
const tenantAddress = ref('')
const tenantCodeJoin = ref('')
const password = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const errorMessage = ref('')
const errors = ref({
  fullname: '',
  email: '',
  tenant_name: '',
  tenant_code: '',
  tenant_address: '',
  tenant_code_join: '',
  password: '',
  confirmPassword: ''
})

const handleRegister = async () => {
  loading.value = true
  errorMessage.value = ''
  errors.value = {
    fullname: '',
    email: '',
    tenant_name: '',
    tenant_code: '',
    tenant_address: '',
    tenant_code_join: '',
    password: '',
    confirmPassword: ''
  }

  // Validation
  if (!fullname.value) {
    errors.value.fullname = 'Nama lengkap harus diisi'
    loading.value = false
    return
  }
  if (!email.value) {
    errors.value.email = 'Email harus diisi'
    loading.value = false
    return
  }
  if (registrationType.value === 'tenant') {
    if (!tenantName.value) {
      errors.value.tenant_name = 'Nama komunitas harus diisi'
      loading.value = false
      return
    }
    if (!tenantCode.value) {
      errors.value.tenant_code = 'Kode komunitas harus diisi'
      loading.value = false
      return
    }
  } else {
    if (!tenantCodeJoin.value) {
      errors.value.tenant_code_join = 'Kode komunitas harus diisi'
      loading.value = false
      return
    }
  }
  if (!password.value) {
    errors.value.password = 'Password harus diisi'
    loading.value = false
    return
  }
  if (password.value.length < 8) {
    errors.value.password = 'Password minimal 8 karakter'
    loading.value = false
    return
  }
  if (password.value !== confirmPassword.value) {
    errors.value.confirmPassword = 'Password tidak sama'
    loading.value = false
    return
  }

  try {
    if (registrationType.value === 'tenant') {
      // Register as Tenant - create tenant and user as Admin
      const body = {
        email: email.value,
        password: password.value,
        full_name: fullname.value,
        registration_type: 'tenant',
        tenant_name: tenantName.value,
        tenant_code: tenantCode.value,
        tenant_address: tenantAddress.value || null,
      }

      const response = await fetch<{ token: string; user: any; tenant_id: string }>('/api/auth/register', {
        method: 'POST',
        body: JSON.stringify(body)
      })

      // Set token
      authStore.setToken(response.token)

      // Map user data
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

      // Fetch current user to get full details
      // Note: We already have role_name from register response, so this is optional
      try {
        const currentUser = await authStore.fetchCurrentUser()
        // Ensure role_name is preserved
        if (currentUser && !currentUser.role_name && userData.role_name) {
          authStore.setUser({ ...currentUser, role_name: userData.role_name })
        }
      } catch (err) {
        console.warn('Failed to fetch current user details:', err)
      }

      // Redirect to dashboard
      await navigateTo('/dashboard')
    } else {
      // Register as Warga - join existing tenant
      const body = {
        email: email.value,
        password: password.value,
        full_name: fullname.value,
        registration_type: 'warga',
        tenant_code_join: tenantCodeJoin.value,
      }

      const response = await fetch<{ token: string; user: any; tenant_id: string }>('/api/auth/register', {
        method: 'POST',
        body: JSON.stringify(body)
      })

      // Set token
      authStore.setToken(response.token)

      // Map user data
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

      // Fetch current user to get full details
      // Note: We already have role_name from register response, so this is optional
      try {
        const currentUser = await authStore.fetchCurrentUser()
        // Ensure role_name is preserved
        if (currentUser && !currentUser.role_name && userData.role_name) {
          authStore.setUser({ ...currentUser, role_name: userData.role_name })
        }
      } catch (err) {
        console.warn('Failed to fetch current user details:', err)
      }

      // Redirect to dashboard
      await navigateTo('/dashboard')
    }
  } catch (error: any) {
    console.error('Register error:', error)
    let errorMsg = 'Terjadi kesalahan saat mendaftar. Silakan coba lagi.'
    
    if (error.data) {
      if (typeof error.data === 'string') {
        errorMsg = error.data
      } else if (error.data.error) {
        errorMsg = error.data.error
      } else if (error.data.message) {
        errorMsg = error.data.message
      }
    } else if (error.message) {
      errorMsg = error.message
    }
    
    errorMessage.value = errorMsg
  } finally {
    loading.value = false
  }
}
</script>
