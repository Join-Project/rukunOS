<template>
  <div class="max-w-4xl mx-auto space-y-6">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-bold text-gray-900 tracking-tight">Data Keluarga</h1>
          <p class="text-gray-500">Kelola data anggota keluarga{{ unitCode ? ' di ' + unitCode : '' }}</p>
        </div>
      </div>

      <div v-if="!loading && familyMembers.length === 0" class="text-center py-12 text-gray-500">
        <p>Belum ada data keluarga</p>
      </div>
      <div v-else-if="familyMembers.length > 0" class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- Family Member Card -->
        <div 
          v-for="(member, index) in familyMembers" 
          :key="member.id"
          class="bg-white rounded-2xl border border-gray-100 p-6 shadow-sm relative overflow-hidden group"
        >
          <div class="absolute top-0 right-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity">
            <svg class="w-24 h-24 text-primary-600" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clip-rule="evenodd"></path></svg>
          </div>
          <div class="relative z-10">
            <div class="flex items-start justify-between mb-4">
              <div :class="[
                'w-16 h-16 rounded-full flex items-center justify-center text-2xl font-bold border-4 border-white shadow-sm',
                index === 0 ? 'bg-primary-100 text-primary-600' :
                index === 1 ? 'bg-pink-100 text-pink-600' :
                'bg-blue-100 text-blue-600'
              ]">
                {{ getInitials(member.full_name) }}
              </div>
              <span :class="[
                'px-3 py-1 text-xs font-bold rounded-full uppercase tracking-wider',
                index === 0 ? 'bg-primary-50 text-primary-700' :
                index === 1 ? 'bg-pink-50 text-pink-700' :
                'bg-blue-50 text-blue-700'
              ]">
                {{ member.relationship || 'Anggota' }}
              </span>
            </div>
            <h3 class="text-xl font-bold text-gray-900">{{ member.full_name }}</h3>
            <p class="text-gray-500 text-sm mb-4" v-if="member.email">{{ member.email }}</p>
            
            <div class="space-y-2 text-sm text-gray-600">
              <div class="flex items-center gap-2" v-if="member.phone">
                <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"></path></svg>
                {{ member.phone }}
              </div>
              <div class="flex items-center gap-2" v-if="member.role_name">
                <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 13.255A23.931 23.931 0 0112 15c-3.183 0-6.22-.62-9-1.745M16 6V4a2 2 0 00-2-2h-4a2 2 0 00-2 2v2m4 6h.01M5 20h14a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"></path></svg>
                {{ member.role_name }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useApi } from '~/composables/useApi'
import { useToast } from '~/composables/useToast'
import { useGlobalLoading } from '~/composables/useGlobalLoading'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

const { fetch } = useApi()
const { showError } = useToast()

const loading = ref(false)
const familyMembers = ref([])
const unitCode = ref('')

// Sync loading state dengan global loading spinner
useGlobalLoading(loading, 'Memuat data keluarga...')

const loadFamily = async () => {
  loading.value = true
  try {
    const response = await fetch('/api/family')
    familyMembers.value = response.family_members || []
    unitCode.value = response.unit_code || ''
  } catch (error: any) {
    console.error('Failed to load family:', error)
    showError(error.message || 'Gagal memuat data keluarga')
  } finally {
    loading.value = false
  }
}

const getInitials = (name: string) => {
  if (!name) return '?'
  const parts = name.split(' ')
  if (parts.length >= 2) {
    return (parts[0][0] + parts[1][0]).toUpperCase()
  }
  return name.substring(0, 2).toUpperCase()
}

onMounted(() => {
  loadFamily()
})
</script>
