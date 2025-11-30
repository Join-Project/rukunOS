<template>
  <div>
      <!-- Page Header -->
      <div class="mb-8 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 class="text-3xl font-bold text-gray-900 mb-1">Role & Izin</h1>
          <p class="text-gray-600">Kelola role dan permissions untuk pengguna</p>
        </div>
        <button
          @click="showCreateModal = true"
          class="btn btn-primary inline-flex items-center px-4 py-2.5"
        >
          <Icon name="heroicons:plus" class="w-5 h-5 mr-2" />
          Tambah Role
        </button>
      </div>

      <!-- Roles Grid -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div
          v-for="role in roles"
          :key="role.id"
          class="card card-hover p-6"
        >
          <div class="flex items-start justify-between mb-4">
            <div class="flex-1">
              <div class="flex items-center space-x-2 mb-2">
                <h3 class="text-lg font-semibold text-gray-900">{{ role.name }}</h3>
                <span
                  v-if="role.is_system"
                  class="badge badge-primary text-xs"
                >
                  System
                </span>
              </div>
              <p v-if="role.description" class="text-sm text-gray-600 mb-4">{{ role.description }}</p>
            </div>
          </div>

          <div class="mb-4">
            <p class="text-xs font-medium text-gray-500 mb-2 uppercase tracking-wide">
              Permissions ({{ role.permissions?.length || 0 }})
            </p>
            <div class="flex flex-wrap gap-1.5">
              <span
                v-for="perm in role.permissions?.slice(0, 4)"
                :key="perm"
                class="inline-flex items-center px-2 py-1 rounded-md text-xs font-medium bg-gray-100 text-gray-700 border border-gray-200"
              >
                {{ perm }}
              </span>
              <span
                v-if="role.permissions && role.permissions.length > 4"
                class="inline-flex items-center px-2 py-1 rounded-md text-xs font-medium bg-gray-100 text-gray-700 border border-gray-200"
              >
                +{{ role.permissions.length - 4 }} lainnya
              </span>
              <span
                v-if="!role.permissions || role.permissions.length === 0"
                class="inline-flex items-center px-2 py-1 rounded-md text-xs font-medium text-gray-400 italic"
              >
                Tidak ada permissions
              </span>
            </div>
          </div>

          <div class="flex space-x-2 pt-4 border-t border-gray-100">
            <button
              @click="editRole(role)"
              class="btn btn-secondary flex-1 text-sm px-3 py-2"
            >
              <Icon name="heroicons:pencil" class="w-4 h-4 mr-1.5" />
              Edit
            </button>
            <button
              v-if="!role.is_system"
              @click="deleteRole(role.id)"
              class="btn btn-ghost text-sm px-3 py-2 text-red-600 hover:text-red-700 hover:bg-red-50"
            >
              <Icon name="heroicons:trash" class="w-4 h-4" />
            </button>
          </div>
        </div>

        <!-- Empty State -->
        <div v-if="!loading && roles.length === 0" class="col-span-full">
          <div class="card p-12 text-center">
            <Icon name="heroicons:shield-check" class="w-16 h-16 mx-auto mb-4 text-gray-300" />
            <h3 class="text-lg font-semibold text-gray-900 mb-2">Belum ada Role</h3>
            <p class="text-gray-600 mb-6">Mulai dengan membuat role pertama Anda</p>
            <button
              @click="showCreateModal = true"
              class="btn btn-primary inline-flex items-center"
            >
              <Icon name="heroicons:plus" class="w-5 h-5 mr-2" />
              Tambah Role
            </button>
          </div>
        </div>
      </div>

      <!-- Create/Edit Modal -->
      <div
        v-if="showCreateModal || editingRole"
        class="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4"
        @click.self="closeModal"
      >
        <div class="bg-white rounded-2xl shadow-xl max-w-2xl w-full p-6 max-h-[90vh] overflow-y-auto">
          <div class="flex items-center justify-between mb-6">
            <h3 class="text-xl font-semibold text-gray-900">
              {{ editingRole ? 'Edit Role' : 'Tambah Role Baru' }}
            </h3>
            <button
              @click="closeModal"
              class="p-2 rounded-lg text-gray-400 hover:text-gray-500 hover:bg-gray-100 transition-colors"
            >
              <Icon name="heroicons:x-mark" class="w-5 h-5" />
            </button>
          </div>
          
          <form @submit.prevent="saveRole" class="space-y-5">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1.5">Nama Role *</label>
              <input
                v-model="form.name"
                type="text"
                required
                class="input"
                placeholder="Contoh: Pengurus Kebersihan"
              />
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1.5">Deskripsi</label>
              <textarea
                v-model="form.description"
                rows="2"
                class="input"
                placeholder="Deskripsi role..."
              ></textarea>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Permissions *</label>
              <div class="border border-gray-200 rounded-xl p-4 max-h-64 overflow-y-auto bg-gray-50">
                <div v-if="permissions.length === 0" class="text-center py-8 text-gray-500">
                  <Icon name="heroicons:arrow-path" class="w-8 h-8 mx-auto mb-2 animate-spin" />
                  <p class="text-sm">Memuat permissions...</p>
                </div>
                <div v-else class="space-y-2">
                  <div
                    v-for="permission in permissions"
                    :key="permission.id"
                    class="flex items-start space-x-3 p-2 rounded-lg hover:bg-white transition-colors"
                  >
                    <input
                      :id="`perm-${permission.id}`"
                      v-model="form.permissions"
                      type="checkbox"
                      :value="permission.key"
                      class="mt-1 w-4 h-4 text-primary-600 border-gray-300 rounded focus:ring-primary-500 focus:ring-2"
                    />
                    <label :for="`perm-${permission.id}`" class="flex-1 cursor-pointer">
                      <div class="text-sm font-medium text-gray-900">{{ permission.name }}</div>
                      <div class="text-xs text-gray-500 mt-0.5 font-mono">{{ permission.key }}</div>
                    </label>
                  </div>
                </div>
              </div>
              <p v-if="form.permissions.length > 0" class="text-xs text-gray-500 mt-2">
                {{ form.permissions.length }} permission{{ form.permissions.length > 1 ? 's' : '' }} dipilih
              </p>
            </div>
            
            <div class="flex space-x-3 pt-4 border-t border-gray-200">
              <button
                type="button"
                @click="closeModal"
                class="btn btn-secondary flex-1"
              >
                Batal
              </button>
              <button
                type="submit"
                :disabled="loading || form.permissions.length === 0"
                class="btn btn-primary flex-1 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <Icon v-if="loading" name="heroicons:arrow-path" class="w-4 h-4 mr-2 animate-spin" />
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

const roles = ref([])
const permissions = ref([])
const showCreateModal = ref(false)
const editingRole = ref(null)
const loading = ref(false)

// Sync loading state dengan global loading spinner
useGlobalLoading(loading, 'Memuat data role...')

const form = ref({
  name: '',
  description: '',
  permissions: [] as string[],
})

const loadRoles = async () => {
  try {
    const response = await fetch('/api/roles')
    roles.value = response.roles || []
  } catch (error) {
    console.error('Failed to load roles:', error)
  }
}

const loadPermissions = async () => {
  try {
    const response = await fetch('/api/roles/permissions')
    permissions.value = response.permissions || []
  } catch (error) {
    console.error('Failed to load permissions:', error)
  }
}

const saveRole = async () => {
  loading.value = true
  try {
    if (editingRole.value) {
      await fetch(`/api/roles/${editingRole.value.id}`, {
        method: 'PUT',
        body: JSON.stringify({
          name: form.value.name,
          description: form.value.description || null,
          permissions: form.value.permissions,
        }),
      })
    } else {
      await fetch('/api/roles', {
        method: 'POST',
        body: JSON.stringify({
          name: form.value.name,
          description: form.value.description || null,
          permissions: form.value.permissions,
        }),
      })
    }
    closeModal()
    loadRoles()
    showSuccess('Role berhasil disimpan')
  } catch (error) {
    console.error('Failed to save role:', error)
    showError('Gagal menyimpan role. Silakan coba lagi.')
  } finally {
    loading.value = false
  }
}

const editRole = (role: any) => {
  editingRole.value = role
  form.value = {
    name: role.name,
    description: role.description || '',
    permissions: role.permissions || [],
  }
  showCreateModal.value = true
}

const deleteRole = async (id: string) => {
  const result = await confirm.show(
    'Apakah Anda yakin ingin menghapus role ini?',
    {
      title: 'Hapus Role',
      confirmText: 'Ya, Hapus',
      cancelText: 'Batal',
      type: 'danger'
    }
  )
  if (!result) return
  try {
    await fetch(`/api/roles/${id}`, { method: 'DELETE' })
    loadRoles()
    showSuccess('Role berhasil dihapus')
  } catch (error) {
    console.error('Failed to delete role:', error)
    showError('Gagal menghapus role. Silakan coba lagi.')
  }
}

const closeModal = () => {
  showCreateModal.value = false
  editingRole.value = null
  form.value = {
    name: '',
    description: '',
    permissions: [],
  }
}

onMounted(() => {
  loadRoles()
  loadPermissions()
})
</script>
