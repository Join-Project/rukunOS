<template>
  <div>
    <!-- Page Header -->
    <div class="mb-8 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <h1 class="text-3xl font-bold text-gray-900 mb-1">Manajemen Pengguna</h1>
        <p class="text-gray-600">Kelola pengguna dan akses dalam komunitas</p>
      </div>
      <button
        @click="showCreateModal = true"
        class="px-4 py-2.5 bg-primary-600 text-white rounded-xl font-medium hover:bg-primary-700 transition-all shadow-lg shadow-primary-600/20 hover:shadow-primary-600/30 hover:-translate-y-0.5 flex items-center gap-2"
      >
        <Icon name="heroicons:plus" class="w-5 h-5" />
        Tambah Warga
      </button>
    </div>

    <!-- Filters -->
    <div class="card p-5 mb-6">
      <div class="grid grid-cols-1 md:grid-cols-5 gap-4">
        <div class="md:col-span-2">
          <input
            v-model="filters.search"
            type="text"
            placeholder="Cari nama atau email..."
            class="input px-4 py-2.5 w-full"
            @keyup.enter="loadUsers"
          />
        </div>
        <div>
          <select
            v-model="filters.role_id"
            class="input px-4 py-2.5 w-full"
          >
            <option value="">Semua Role</option>
            <option v-for="role in roles" :key="role.id" :value="role.id">
              {{ role.name }}
            </option>
          </select>
        </div>
        <div>
          <select
            v-model="filters.unit_id"
            class="input px-4 py-2.5 w-full"
          >
            <option value="">Semua Unit</option>
            <option value="unassigned">Belum Ter-assign</option>
            <option v-for="unit in units" :key="unit.id" :value="unit.id">
              {{ unit.code }}
            </option>
          </select>
        </div>
        <div>
          <button
            @click="loadUsers"
            class="btn-secondary w-full"
          >
            Filter
          </button>
        </div>
      </div>
    </div>

    <!-- Users Table -->
    <div class="card overflow-hidden">
      <div class="overflow-x-auto">
        <table class="table">
          <thead>
            <tr>
              <th>Pengguna</th>
              <th>Email</th>
              <th>Role</th>
              <th>Unit</th>
              <th class="text-right">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="user in users" :key="user.id">
              <td>
                <div class="flex items-center">
                  <div class="w-10 h-10 bg-gradient-to-br from-primary-500 to-primary-600 rounded-lg flex items-center justify-center mr-3 flex-shrink-0">
                    <span class="text-white font-semibold text-sm">
                      {{ user.full_name?.charAt(0).toUpperCase() }}
                    </span>
                  </div>
                  <div>
                    <div class="font-semibold text-gray-900">{{ user.full_name }}</div>
                  </div>
                </div>
              </td>
              <td>
                <div class="text-sm text-gray-600">{{ user.email }}</div>
              </td>
              <td>
                <span v-if="user.role" class="badge-primary">
                  {{ user.role.name }}
                </span>
                <span v-else class="text-sm text-gray-400">-</span>
              </td>
              <td>
                <span v-if="user.unit" class="text-sm text-gray-900 font-medium">{{ user.unit.code }}</span>
                <span v-else class="text-sm text-gray-400 italic">Belum ter-assign</span>
              </td>
              <td class="text-right">
                <button
                  @click="editUser(user)"
                  class="text-primary-600 hover:text-primary-700 font-medium text-sm"
                >
                  Edit
                </button>
              </td>
            </tr>
            <tr v-if="!loading && users.length === 0">
              <td colspan="5" class="text-center py-12">
                <div class="text-gray-400">
                  <Icon name="heroicons:users" class="w-12 h-12 mx-auto mb-3 opacity-50" />
                  <p class="text-sm">Belum ada pengguna</p>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div v-if="pagination.total_pages > 1" class="px-6 py-4 border-t border-gray-200 flex items-center justify-between bg-gray-50">
        <div class="text-sm text-gray-700">
          Menampilkan <span class="font-medium">{{ (pagination.page - 1) * pagination.limit + 1 }}</span> sampai
          <span class="font-medium">{{ Math.min(pagination.page * pagination.limit, pagination.total) }}</span> dari
          <span class="font-medium">{{ pagination.total }}</span> pengguna
        </div>
        <div class="flex space-x-2">
          <button
            @click="changePage(pagination.page - 1)"
            :disabled="pagination.page === 1"
            class="btn-ghost px-3 py-2 text-sm disabled:opacity-50"
          >
            Sebelumnya
          </button>
          <button
            @click="changePage(pagination.page + 1)"
            :disabled="pagination.page >= pagination.total_pages"
            class="btn-ghost px-3 py-2 text-sm disabled:opacity-50"
          >
            Selanjutnya
          </button>
        </div>
      </div>
    </div>

    <!-- Create User Modal -->
    <div
      v-if="showCreateModal"
      class="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4"
      @click.self="closeCreateModal"
    >
      <div class="bg-white rounded-2xl shadow-xl max-w-md w-full p-6 max-h-[90vh] overflow-y-auto">
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-xl font-semibold text-gray-900">Tambah Warga Baru</h3>
          <button
            @click="closeCreateModal"
            class="p-2 rounded-lg text-gray-400 hover:text-gray-500 hover:bg-gray-100"
          >
            <Icon name="heroicons:x-mark" class="w-5 h-5" />
          </button>
        </div>
        
        <form @submit.prevent="createUser" class="space-y-5">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1.5">Nama Lengkap</label>
            <input
              v-model="createForm.full_name"
              type="text"
              class="input px-4 py-2.5 w-full"
              placeholder="Nama lengkap warga"
              required
            />
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1.5">Email</label>
            <input
              v-model="createForm.email"
              type="email"
              class="input px-4 py-2.5 w-full"
              placeholder="email@example.com"
              required
            />
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1.5">Password</label>
            <input
              v-model="createForm.password"
              type="password"
              class="input px-4 py-2.5 w-full"
              placeholder="Minimal 8 karakter"
              required
              minlength="8"
            />
            <p class="mt-1 text-xs text-gray-500">Warga akan menggunakan password ini untuk login</p>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1.5">Role</label>
            <select
              v-model="createForm.role_id"
              class="input px-4 py-2.5 w-full"
              required
            >
              <option value="">Pilih Role</option>
              <option v-for="role in roles" :key="role.id" :value="role.id">
                {{ role.name }}
              </option>
            </select>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1.5">Unit Rumah</label>
            <select
              v-model="createForm.unit_id"
              class="input px-4 py-2.5 w-full"
            >
              <option value="">Tidak Ter-assign (Nanti)</option>
              <option v-for="unit in units" :key="unit.id" :value="unit.id">
                {{ unit.code }} - {{ unit.type }}
              </option>
            </select>
            <p class="mt-1 text-xs text-gray-500">Pilih unit rumah untuk warga ini (opsional)</p>
          </div>
          
          <div class="flex space-x-3 pt-4">
            <button
              type="button"
              @click="closeCreateModal"
              class="px-6 py-2.5 bg-white border border-gray-300 text-gray-700 font-medium rounded-xl hover:bg-gray-50 transition-colors shadow-sm flex-1"
            >
              Batal
            </button>
            <button
              type="submit"
              :disabled="loading"
              class="px-6 py-2.5 bg-primary-600 text-white font-medium rounded-xl hover:bg-primary-700 transition-all shadow-lg shadow-primary-600/20 hover:shadow-primary-600/30 disabled:opacity-50 disabled:cursor-not-allowed flex-1"
            >
              {{ loading ? 'Menyimpan...' : 'Simpan' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Edit Modal -->
    <div
      v-if="editingUser"
      class="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4"
      @click.self="closeModal"
    >
      <div class="bg-white rounded-2xl shadow-xl max-w-md w-full p-6">
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-xl font-semibold text-gray-900">Edit Pengguna</h3>
          <button
            @click="closeModal"
            class="p-2 rounded-lg text-gray-400 hover:text-gray-500 hover:bg-gray-100"
          >
            <Icon name="heroicons:x-mark" class="w-5 h-5" />
          </button>
        </div>
        
        <form @submit.prevent="saveUser" class="space-y-5">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1.5">Role</label>
            <select
              v-model="form.role_id"
              class="input px-4 py-2.5 w-full"
            >
              <option value="">Pilih Role</option>
              <option v-for="role in roles" :key="role.id" :value="role.id">
                {{ role.name }}
              </option>
            </select>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1.5">Unit Rumah</label>
            <select
              v-model="form.unit_id"
              class="input px-4 py-2.5 w-full"
            >
              <option value="">Tidak Ter-assign</option>
              <option v-for="unit in units" :key="unit.id" :value="unit.id">
                {{ unit.code }} - {{ unit.type }}
              </option>
            </select>
            <p class="mt-1 text-xs text-gray-500">Pilih unit rumah untuk warga ini</p>
          </div>
          
          <div class="flex space-x-3 pt-4">
            <button
              type="button"
              @click="closeModal"
              class="btn-secondary flex-1"
            >
              Batal
            </button>
            <button
              type="submit"
              :disabled="loading"
              class="btn-primary flex-1"
            >
              {{ loading ? 'Menyimpan...' : 'Simpan' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
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

const editingUser = ref(null)
const showCreateModal = ref(false)
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

const changePage = (page: number) => {
  pagination.value.page = page
  loadUsers()
}

const editUser = async (user: any) => {
  editingUser.value = user
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
}

const saveUser = async () => {
  loading.value = true
  try {
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
  } catch (error: any) {
    console.error('Failed to save user:', error)
    const errorMessage = error.data?.error || error.message || 'Gagal menyimpan data'
    showError(errorMessage)
  } finally {
    loading.value = false
  }
}

const closeModal = () => {
  editingUser.value = null
  form.value = { role_id: '', unit_id: '' }
}

const closeCreateModal = () => {
  showCreateModal.value = false
  createForm.value = {
    full_name: '',
    email: '',
    password: '',
    role_id: '',
    unit_id: '',
  }
}

const createUser = async () => {
  loading.value = true
  try {
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
    
    closeCreateModal()
    loadUsers()
  } catch (error: any) {
    console.error('Failed to create user:', error)
    const errorMessage = error.data?.error || error.message || 'Gagal membuat warga baru'
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
