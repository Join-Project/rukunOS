<template>
  <div class="max-w-7xl mx-auto">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-6 md:mb-8">
      <div>
        <h1 class="text-xl md:text-2xl font-bold text-gray-900 tracking-tight">Manajemen Warga & Pengguna</h1>
        <p class="text-sm md:text-base text-gray-500 mt-1">Kelola data warga, peran (role), dan status akun.</p>
      </div>
      <button @click="openModal()" class="px-4 py-2 bg-primary-600 text-white rounded-xl font-medium hover:bg-primary-700 transition-colors shadow-lg shadow-primary-600/20 flex items-center justify-center text-sm md:text-base">
        <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path></svg>
        Tambah Warga
      </button>
    </div>

    <!-- Filters -->
    <div class="bg-white rounded-2xl border border-gray-100 shadow-sm p-4 mb-6 flex flex-col md:flex-row gap-4">
      <div class="flex-1 relative">
        <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
          <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path></svg>
        </div>
        <input 
          type="text" 
          v-model="filters.search"
          @keyup.enter="loadUsers"
          class="block w-full pl-10 pr-3 py-2.5 border border-gray-200 rounded-xl leading-5 bg-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent sm:text-sm transition-shadow" 
          placeholder="Cari nama, email, atau unit..."
        >
      </div>
      <div class="w-full md:w-48">
        <select v-model="filters.role_id" @change="loadUsers" class="block w-full pl-3 pr-10 py-2.5 text-base border border-gray-200 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent sm:text-sm rounded-xl bg-white">
          <option value="">Semua Role</option>
          <option v-for="role in roles" :key="role.id" :value="role.id">
            {{ role.name }}
          </option>
        </select>
      </div>
      <div class="w-full md:w-48">
        <select v-model="filters.unit_id" @change="loadUsers" class="block w-full pl-3 pr-10 py-2.5 text-base border border-gray-200 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent sm:text-sm rounded-xl bg-white">
          <option value="">Semua Unit</option>
          <option value="unassigned">Belum Ter-assign</option>
          <option v-for="unit in units" :key="unit.id" :value="unit.id">
            {{ unit.code }}
          </option>
        </select>
      </div>
    </div>

    <!-- Users Table - Desktop -->
    <div class="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden hidden md:block">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Nama & Email</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Role</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Unit</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Bergabung</th>
              <th scope="col" class="relative px-6 py-3">
                <span class="sr-only">Actions</span>
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="user in users" :key="user.id" class="hover:bg-gray-50 transition-colors">
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center">
                  <div class="flex-shrink-0 h-10 w-10 rounded-full bg-gradient-to-br from-primary-500 to-primary-600 flex items-center justify-center text-white font-bold">
                    {{ user.full_name?.charAt(0).toUpperCase() }}
                  </div>
                  <div class="ml-4">
                    <div class="text-sm font-medium text-gray-900">{{ user.full_name }}</div>
                    <div class="text-sm text-gray-500">{{ user.email }}</div>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span v-if="user.role" class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-blue-100 text-blue-800">
                  {{ user.role.name }}
                </span>
                <span v-else class="text-sm text-gray-400">-</span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm">
                <span v-if="user.unit" class="text-gray-900 font-medium">{{ user.unit.code }}</span>
                <span v-else class="text-gray-400 italic">Belum ter-assign</span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span 
                  :class="[
                    'px-2 inline-flex text-xs leading-5 font-semibold rounded-full',
                    user.status === 'active' ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'
                  ]"
                >
                  {{ user.status === 'active' ? 'Active' : 'Inactive' }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ formatDate(user.created_at) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                <button @click="openModal(user)" class="text-primary-600 hover:text-primary-900 mr-4">Edit</button>
              </td>
            </tr>
            <tr v-if="!loading && users.length === 0">
              <td colspan="6" class="px-6 py-12 text-center">
                <div class="text-gray-400">
                  <Icon name="heroicons:users" class="w-12 h-12 mx-auto mb-3 opacity-50" />
                  <p class="text-sm">Belum ada pengguna</p>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Users Cards - Mobile -->
    <div class="md:hidden space-y-4">
      <div v-if="!loading && users.length === 0" class="bg-white rounded-2xl border border-gray-100 shadow-sm p-12 text-center">
        <div class="text-gray-400">
          <Icon name="heroicons:users" class="w-12 h-12 mx-auto mb-3 opacity-50" />
          <p class="text-sm">Belum ada pengguna</p>
        </div>
      </div>
      <div
        v-for="user in users"
        :key="user.id"
        class="bg-white rounded-2xl border border-gray-100 shadow-sm p-4 space-y-3"
      >
        <div class="flex items-start justify-between">
          <div class="flex items-center gap-3 flex-1">
            <div class="flex-shrink-0 h-12 w-12 rounded-full bg-gradient-to-br from-primary-500 to-primary-600 flex items-center justify-center text-white font-bold">
              {{ user.full_name?.charAt(0).toUpperCase() }}
            </div>
            <div class="flex-1 min-w-0">
              <div class="text-sm font-medium text-gray-900 truncate">{{ user.full_name }}</div>
              <div class="text-xs text-gray-500 truncate">{{ user.email }}</div>
            </div>
          </div>
          <button @click="openModal(user)" class="text-primary-600 hover:text-primary-900 text-sm font-medium px-2">
            Edit
          </button>
        </div>
        
        <div class="grid grid-cols-2 gap-3 pt-2 border-t border-gray-100">
          <div>
            <div class="text-xs text-gray-500 mb-1">Role</div>
            <span v-if="user.role" class="inline-flex text-xs leading-5 font-semibold rounded-full bg-blue-100 text-blue-800 px-2 py-0.5">
              {{ user.role.name }}
            </span>
            <span v-else class="text-xs text-gray-400">-</span>
          </div>
          <div>
            <div class="text-xs text-gray-500 mb-1">Status</div>
            <span 
              :class="[
                'inline-flex text-xs leading-5 font-semibold rounded-full px-2 py-0.5',
                user.status === 'active' ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'
              ]"
            >
              {{ user.status === 'active' ? 'Active' : 'Inactive' }}
            </span>
          </div>
        </div>
        
        <div class="pt-2 border-t border-gray-100">
          <div class="text-xs text-gray-500 mb-1">Unit</div>
          <span v-if="user.unit" class="text-sm text-gray-900 font-medium">{{ user.unit.code }}</span>
          <span v-else class="text-sm text-gray-400 italic">Belum ter-assign</span>
        </div>
        
        <div class="pt-2 border-t border-gray-100">
          <div class="text-xs text-gray-500 mb-1">Bergabung</div>
          <div class="text-sm text-gray-900">{{ formatDate(user.created_at) }}</div>
        </div>
      </div>
    </div>
  </div>

  <!-- Add/Edit Modal -->
  <div v-if="isModalOpen" class="fixed inset-0 z-50 overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
    <div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
      <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" aria-hidden="true" @click="closeModal"></div>
      <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>
      <div class="inline-block align-bottom bg-white rounded-2xl text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full w-full max-h-[90vh] overflow-y-auto">
        <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
          <div class="sm:flex sm:items-start">
            <div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left w-full">
              <h3 class="text-lg leading-6 font-medium text-gray-900" id="modal-title">
                {{ editingUser ? 'Edit Warga' : 'Tambah Warga Baru' }}
              </h3>
              <div class="mt-4 space-y-4">
                <div v-if="!editingUser">
                  <label class="block text-sm font-medium text-gray-700">Nama Lengkap</label>
                  <input v-model="createForm.full_name" type="text" class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-primary-500 focus:border-primary-500 sm:text-sm p-2 border" required>
                </div>
                <div v-if="!editingUser">
                  <label class="block text-sm font-medium text-gray-700">Email</label>
                  <input v-model="createForm.email" type="email" class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-primary-500 focus:border-primary-500 sm:text-sm p-2 border" required>
                </div>
                <div v-if="!editingUser">
                  <label class="block text-sm font-medium text-gray-700">Password</label>
                  <input v-model="createForm.password" type="password" class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-primary-500 focus:border-primary-500 sm:text-sm p-2 border" required minlength="8">
                  <p class="mt-1 text-xs text-gray-500">Minimal 8 karakter</p>
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700">Role</label>
                  <select v-model="form.role_id" class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-primary-500 focus:border-primary-500 sm:text-sm p-2 border">
                    <option value="">Pilih Role</option>
                    <option v-for="role in roles" :key="role.id" :value="role.id">
                      {{ role.name }}
                    </option>
                  </select>
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700">Unit Rumah</label>
                  <select v-model="form.unit_id" class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-primary-500 focus:border-primary-500 sm:text-sm p-2 border">
                    <option value="">Tidak Ter-assign</option>
                    <option v-for="unit in units" :key="unit.id" :value="unit.id">
                      {{ unit.code }} - {{ unit.type }}
                    </option>
                  </select>
                  <p class="mt-1 text-xs text-gray-500">Pilih unit rumah untuk warga ini</p>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
          <button type="button" :disabled="loading" class="w-full inline-flex justify-center rounded-xl border border-transparent shadow-sm px-4 py-2 bg-primary-600 text-base font-medium text-white hover:bg-primary-700 focus:outline-none sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50 disabled:cursor-not-allowed" @click="saveUser">
            {{ loading ? 'Menyimpan...' : 'Simpan' }}
          </button>
          <button type="button" class="mt-3 w-full inline-flex justify-center rounded-xl border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm" @click="closeModal">
            Batal
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useGlobalLoading } from '~/composables/useGlobalLoading'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

const { fetch } = useApi()
const { showError } = useToast()

const users = ref([])
const roles = ref([])
const units = ref([])
const pagination = ref({
  page: 1,
  limit: 20,
  total: 0,
  total_pages: 0,
})

const filters = ref({
  search: '',
  role_id: '',
  unit_id: '',
})

const isModalOpen = ref(false)
const editingUser = ref(null)
const loading = ref(false)

// Sync loading state dengan global loading spinner
useGlobalLoading(loading, 'Memuat data pengguna...')

const form = ref({
  role_id: '',
  unit_id: '',
})

const createForm = ref({
  full_name: '',
  email: '',
  password: '',
  role_id: '',
  unit_id: '',
})

const loadUsers = async () => {
  try {
    const params = new URLSearchParams({
      page: pagination.value.page.toString(),
      limit: pagination.value.limit.toString(),
    })
    if (filters.value.search) params.append('search', filters.value.search)
    if (filters.value.role_id) params.append('role_id', filters.value.role_id)
    if (filters.value.unit_id && filters.value.unit_id !== 'unassigned') {
      params.append('unit_id', filters.value.unit_id)
    }

    const response = await fetch(`/api/users?${params.toString()}`)
    let usersList = response.users || []
    
    // Filter unassigned users if needed
    if (filters.value.unit_id === 'unassigned') {
      usersList = usersList.filter((u: any) => !u.unit)
    }
    
    users.value = usersList
    pagination.value = response.pagination || pagination.value
  } catch (error) {
    console.error('Failed to load users:', error)
  }
}

const loadRoles = async () => {
  try {
    const response = await fetch('/api/roles')
    roles.value = response.roles || []
  } catch (error) {
    console.error('Failed to load roles:', error)
  }
}

const loadUnits = async () => {
  try {
    const response = await fetch('/api/units?limit=1000')
    units.value = response.units || []
  } catch (error) {
    console.error('Failed to load units:', error)
  }
}

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('id-ID', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

const openModal = async (user = null) => {
  editingUser.value = user
  if (user) {
    // Load fresh user data to ensure we have latest unit info
    try {
      const userDetail = await fetch(`/api/users/${user.id}`)
      form.value = {
        role_id: userDetail.role?.id || user.role?.id || '',
        unit_id: userDetail.unit?.id || user.unit?.id || '',
      }
    } catch (error) {
      console.error('Failed to load user detail:', error)
      // Fallback to existing user data
      form.value = {
        role_id: user.role?.id || '',
        unit_id: user.unit?.id || '',
      }
    }
  } else {
    form.value = {
      role_id: '',
      unit_id: '',
    }
    createForm.value = {
      full_name: '',
      email: '',
      password: '',
      role_id: '',
      unit_id: '',
    }
  }
  isModalOpen.value = true
}

const closeModal = () => {
  isModalOpen.value = false
  editingUser.value = null
  form.value = { role_id: '', unit_id: '' }
  createForm.value = {
    full_name: '',
    email: '',
    password: '',
    role_id: '',
    unit_id: '',
  }
}

const saveUser = async () => {
  loading.value = true
  try {
    if (editingUser.value) {
      // Update existing user
      const updateBody: any = {}
      if (form.value.role_id) {
        updateBody.role_id = form.value.role_id
      }
      if (form.value.unit_id !== undefined) {
        // Allow empty string to unassign unit
        updateBody.unit_id = form.value.unit_id || null
      }
      
      if (Object.keys(updateBody).length > 0) {
        await fetch(`/api/users/${editingUser.value.id}`, {
          method: 'PUT',
          body: JSON.stringify(updateBody),
        })
        closeModal()
        loadUsers()
      } else {
        closeModal()
      }
    } else {
      // Create new user
      const authStore = useAuthStore()
      
      // Use new admin endpoint to create user (no need for tenant_code)
      const createBody: any = {
        email: createForm.value.email,
        password: createForm.value.password,
        full_name: createForm.value.full_name,
      }
      
      if (createForm.value.role_id) {
        createBody.role_id = createForm.value.role_id
      }
      if (createForm.value.unit_id) {
        createBody.unit_id = createForm.value.unit_id
      }
      
      await fetch('/api/users', {
        method: 'POST',
        body: JSON.stringify(createBody),
      })
      
      closeModal()
      loadUsers()
    }
  } catch (error: any) {
    console.error('Failed to save user:', error)
    const errorMessage = error.data?.error || error.message || 'Gagal menyimpan data'
    showError(errorMessage)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadUsers()
  loadRoles()
  loadUnits()
})
</script>
