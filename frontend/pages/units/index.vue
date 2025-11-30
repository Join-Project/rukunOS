<template>
  <div>
    <!-- Page Header -->
    <div class="mb-8 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <h1 class="text-3xl font-bold text-gray-900 mb-1">Manajemen Unit</h1>
        <p class="text-gray-600">Kelola unit properti di komunitas Anda</p>
      </div>
      <button
        @click="showCreateModal = true"
        class="px-4 py-2.5 bg-primary-600 text-white rounded-xl font-medium hover:bg-primary-700 transition-all shadow-lg shadow-primary-600/20 hover:shadow-primary-600/30 hover:-translate-y-0.5 flex items-center gap-2"
      >
        <Icon name="heroicons:plus" class="w-5 h-5" />
        Tambah Unit
      </button>
    </div>

    <!-- Filters -->
    <div class="card p-5 mb-6">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div class="md:col-span-2">
          <input
            v-model="filters.search"
            type="text"
            placeholder="Cari unit atau nama pemilik..."
            class="input px-4 py-2.5 w-full"
            @keyup.enter="loadUnits"
          />
        </div>
        <div>
          <select
            v-model="filters.type"
            class="input px-4 py-2.5 w-full"
          >
            <option value="">Semua Tipe</option>
            <option value="rumah">Rumah</option>
            <option value="ruko">Ruko</option>
            <option value="kios">Kios</option>
          </select>
        </div>
        <div>
          <button
            @click="loadUnits"
            class="px-4 py-2.5 bg-white border border-gray-200 text-gray-700 rounded-xl font-medium hover:bg-gray-50 transition-colors shadow-sm w-full"
          >
            Filter
          </button>
        </div>
      </div>
    </div>

    <!-- Units Table -->
    <div class="card overflow-hidden">
      <div class="overflow-x-auto">
        <table class="table">
          <thead>
            <tr>
              <th>Kode Unit</th>
              <th>Tipe</th>
              <th>Pemilik</th>
              <th>Warga</th>
              <th>Kontak</th>
              <th>Status</th>
              <th class="text-right">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="unit in units" :key="unit.id">
              <td>
                <div class="font-semibold text-gray-900">{{ unit.code }}</div>
              </td>
              <td>
                <span class="badge-primary capitalize">{{ unit.type }}</span>
              </td>
              <td>
                <div class="text-gray-900">{{ unit.owner_name || '-' }}</div>
              </td>
              <td>
                <div v-if="unit.users && unit.users.length > 0" class="space-y-1">
                  <div
                    v-for="user in unit.users"
                    :key="user.id"
                    class="text-sm text-gray-900 font-medium"
                  >
                    {{ user.full_name }}
                  </div>
                </div>
                <span v-else class="text-sm text-gray-400 italic">Belum ada warga</span>
              </td>
              <td>
                <div class="text-sm text-gray-600">{{ unit.owner_phone || '-' }}</div>
              </td>
              <td>
                <span
                  :class="[
                    'badge',
                    unit.status === 'active' ? 'badge-success' : 'badge-gray'
                  ]"
                >
                  {{ unit.status }}
                </span>
              </td>
              <td class="text-right">
                <div class="flex items-center justify-end space-x-2">
                  <button
                    @click="editUnit(unit)"
                    class="text-primary-600 hover:text-primary-700 font-medium text-sm"
                  >
                    Edit
                  </button>
                  <span class="text-gray-300">|</span>
                  <button
                    @click="deleteUnit(unit.id)"
                    class="text-red-600 hover:text-red-700 font-medium text-sm"
                  >
                    Hapus
                  </button>
                </div>
              </td>
            </tr>
            <tr v-if="!loading && units.length === 0">
              <td colspan="7" class="text-center py-12">
                <div class="text-gray-400">
                  <Icon name="heroicons:building-office" class="w-12 h-12 mx-auto mb-3 opacity-50" />
                  <p class="text-sm">Belum ada unit</p>
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
          <span class="font-medium">{{ pagination.total }}</span> unit
        </div>
        <div class="flex space-x-2">
          <button
            @click="changePage(pagination.page - 1)"
            :disabled="pagination.page === 1"
            class="px-3 py-2 text-sm text-gray-600 hover:bg-gray-100 rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:bg-transparent"
          >
            Sebelumnya
          </button>
          <button
            @click="changePage(pagination.page + 1)"
            :disabled="pagination.page >= pagination.total_pages"
            class="px-3 py-2 text-sm text-gray-600 hover:bg-gray-100 rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:bg-transparent"
          >
            Selanjutnya
          </button>
        </div>
      </div>
    </div>

    <!-- Create/Edit Modal -->
    <div
      v-if="showCreateModal || editingUnit"
      class="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4"
      @click.self="closeModal"
    >
      <div class="bg-white rounded-2xl shadow-xl max-w-md w-full p-6 max-h-[90vh] overflow-y-auto">
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-xl font-semibold text-gray-900">
            {{ editingUnit ? 'Edit Unit' : 'Tambah Unit Baru' }}
          </h3>
          <button
            @click="closeModal"
            class="p-2 rounded-lg text-gray-400 hover:text-gray-500 hover:bg-gray-100"
          >
            <Icon name="heroicons:x-mark" class="w-5 h-5" />
          </button>
        </div>
        
        <form @submit.prevent="saveUnit" class="space-y-5">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1.5">Kode Unit *</label>
            <input
              v-model="form.code"
              type="text"
              required
              class="input px-4 py-2.5"
              placeholder="A-01"
            />
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1.5">Tipe *</label>
            <select
              v-model="form.type"
              required
              class="input px-4 py-2.5"
            >
              <option value="">Pilih Tipe</option>
              <option value="rumah">Rumah</option>
              <option value="ruko">Ruko</option>
              <option value="kios">Kios</option>
            </select>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1.5">Nama Pemilik</label>
            <input
              v-model="form.owner_name"
              type="text"
              class="input px-4 py-2.5"
              placeholder="Nama lengkap pemilik"
            />
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1.5">No. Telepon</label>
            <input
              v-model="form.owner_phone"
              type="tel"
              class="input px-4 py-2.5"
              placeholder="081234567890"
            />
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1.5">Email</label>
            <input
              v-model="form.owner_email"
              type="email"
              class="input px-4 py-2.5"
              placeholder="email@example.com"
            />
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1.5">Alamat</label>
            <textarea
              v-model="form.address"
              rows="2"
              class="input px-4 py-2.5"
              placeholder="Alamat unit"
            ></textarea>
          </div>

          <div class="flex space-x-3 pt-4">
            <button
              type="button"
              @click="closeModal"
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
  </div>
</template>

<script setup lang="ts">
import { useGlobalLoading } from '~/composables/useGlobalLoading'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

const { fetch } = useApi()
const { showSuccess, showError } = useToast()
const confirm = useConfirm()

const units = ref([])
const pagination = ref({
  page: 1,
  limit: 20,
  total: 0,
  total_pages: 0,
})

const filters = ref({
  search: '',
  type: '',
})

const showCreateModal = ref(false)
const editingUnit = ref(null)
const loading = ref(false)

// Sync loading state dengan global loading spinner
useGlobalLoading(loading, 'Memuat data unit...')

const form = ref({
  code: '',
  type: '',
  owner_name: '',
  owner_phone: '',
  owner_email: '',
  address: '',
})

const loadUnits = async () => {
  try {
    const params = new URLSearchParams({
      page: pagination.value.page.toString(),
      limit: pagination.value.limit.toString(),
    })
    if (filters.value.search) params.append('search', filters.value.search)
    if (filters.value.type) params.append('type', filters.value.type)

    const response = await fetch(`/api/units?${params.toString()}`)
    const unitsList = response.units || []
    
    // Load users for each unit
    const unitsWithUsers = await Promise.all(
      unitsList.map(async (unit: any) => {
        try {
          const unitDetail = await fetch(`/api/units/${unit.id}`)
          return {
            ...unit,
            users: unitDetail.users || []
          }
        } catch (error) {
          console.error(`Failed to load users for unit ${unit.id}:`, error)
          return {
            ...unit,
            users: []
          }
        }
      })
    )
    
    units.value = unitsWithUsers
    pagination.value = response.pagination || pagination.value
  } catch (error) {
    console.error('Failed to load units:', error)
  }
}

const changePage = (page: number) => {
  pagination.value.page = page
  loadUnits()
}

const saveUnit = async () => {
  loading.value = true
  try {
    if (editingUnit.value) {
      await fetch(`/api/units/${editingUnit.value.id}`, {
        method: 'PUT',
        body: JSON.stringify(form.value),
      })
    } else {
      await fetch('/api/units', {
        method: 'POST',
        body: JSON.stringify(form.value),
      })
    }
    closeModal()
    loadUnits()
  } catch (error) {
    console.error('Failed to save unit:', error)
  } finally {
    loading.value = false
  }
}

const editUnit = (unit: any) => {
  editingUnit.value = unit
  form.value = {
    code: unit.code,
    type: unit.type,
    owner_name: unit.owner_name?.String || '',
    owner_phone: unit.owner_phone?.String || '',
    owner_email: unit.owner_email?.String || '',
    address: unit.address?.String || '',
  }
}

const deleteUnit = async (id: string) => {
  const result = await confirm.show(
    'Apakah Anda yakin ingin menghapus unit ini?',
    {
      title: 'Hapus Unit',
      confirmText: 'Ya, Hapus',
      cancelText: 'Batal',
      type: 'danger'
    }
  )
  if (!result) return
  try {
    await fetch(`/api/units/${id}`, { method: 'DELETE' })
    loadUnits()
    showSuccess('Unit berhasil dihapus')
  } catch (error) {
    console.error('Failed to delete unit:', error)
    showError('Gagal menghapus unit. Silakan coba lagi.')
  }
}

const closeModal = () => {
  showCreateModal.value = false
  editingUnit.value = null
  form.value = {
    code: '',
    type: '',
    owner_name: '',
    owner_phone: '',
    owner_email: '',
    address: '',
  }
}

onMounted(() => {
  loadUnits()
})
</script>
