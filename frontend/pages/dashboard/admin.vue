<template>
  <div class="space-y-6">
      <!-- Page Header -->
      <div class="sm:flex sm:items-center sm:justify-between">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">Dashboard</h1>
          <p class="mt-1 text-sm text-gray-500">Overview kondisi keuangan dan komunitas Anda.</p>
        </div>
        <div class="mt-4 sm:mt-0 flex space-x-3">
          <UiButton variant="secondary" size="sm">
            <span class="mr-2">📥</span> Export Report
          </UiButton>
          <UiButton variant="primary" size="sm">
            <span class="mr-2">➕</span> Create Bill
          </UiButton>
        </div>
      </div>

      <!-- Stats Cards -->
      <div class="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4">
        <UiCard class="relative overflow-hidden">
          <dt>
            <div class="absolute rounded-md bg-primary-500 p-3">
              <svg class="h-6 w-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <p class="ml-16 truncate text-sm font-medium text-gray-500">Total Tagihan</p>
          </dt>
          <dd class="ml-16 flex items-baseline pb-1 sm:pb-2">
            <p class="text-2xl font-semibold text-gray-900">150</p>
            <p class="ml-2 flex items-baseline text-sm font-semibold text-green-600">
              <span class="sr-only">Increased by</span>
              12%
            </p>
          </dd>
        </UiCard>

        <UiCard class="relative overflow-hidden">
          <dt>
            <div class="absolute rounded-md bg-green-500 p-3">
              <svg class="h-6 w-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <p class="ml-16 truncate text-sm font-medium text-gray-500">Collection Rate</p>
          </dt>
          <dd class="ml-16 flex items-baseline pb-1 sm:pb-2">
            <p class="text-2xl font-semibold text-gray-900">75.5%</p>
            <p class="ml-2 flex items-baseline text-sm font-semibold text-green-600">
              <span class="sr-only">Increased by</span>
              5.4%
            </p>
          </dd>
        </UiCard>

        <UiCard class="relative overflow-hidden">
          <dt>
            <div class="absolute rounded-md bg-red-500 p-3">
              <svg class="h-6 w-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <p class="ml-16 truncate text-sm font-medium text-gray-500">Overdue Bills</p>
          </dt>
          <dd class="ml-16 flex items-baseline pb-1 sm:pb-2">
            <p class="text-2xl font-semibold text-gray-900">25</p>
            <p class="ml-2 flex items-baseline text-sm font-semibold text-red-600">
              <span class="sr-only">Increased by</span>
              2
            </p>
          </dd>
        </UiCard>
        
        <UiCard class="relative overflow-hidden">
          <dt>
            <div class="absolute rounded-md bg-blue-500 p-3">
              <svg class="h-6 w-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
              </svg>
            </div>
            <p class="ml-16 truncate text-sm font-medium text-gray-500">Total Warga</p>
          </dt>
          <dd class="ml-16 flex items-baseline pb-1 sm:pb-2">
            <p class="text-2xl font-semibold text-gray-900">342</p>
            <p class="ml-2 flex items-baseline text-sm font-semibold text-gray-500">
              Active
            </p>
          </dd>
        </UiCard>
      </div>

      <!-- Recent Activity & Top Overdue -->
      <div class="grid grid-cols-1 gap-5 lg:grid-cols-2">
        <!-- Recent Activity -->
        <UiCard title="Aktivitas Terbaru">
          <div class="flow-root">
            <ul role="list" class="-mb-8">
              <li v-for="(activity, activityIdx) in activities" :key="activity.id">
                <div class="relative pb-8">
                  <span v-if="activityIdx !== activities.length - 1" class="absolute top-4 left-4 -ml-px h-full w-0.5 bg-gray-200" aria-hidden="true"></span>
                  <div class="relative flex space-x-3">
                    <div>
                      <span :class="['h-8 w-8 rounded-full flex items-center justify-center ring-8 ring-white', activity.iconBackground]">
                        <component :is="activity.icon" class="h-5 w-5 text-white" aria-hidden="true" />
                      </span>
                    </div>
                    <div class="min-w-0 flex-1 pt-1.5 flex justify-between space-x-4">
                      <div>
                        <p class="text-sm text-gray-500">
                          {{ activity.content }} <span class="font-medium text-gray-900">{{ activity.target }}</span>
                        </p>
                      </div>
                      <div class="text-right text-sm whitespace-nowrap text-gray-500">
                        <time :datetime="activity.datetime">{{ activity.date }}</time>
                      </div>
                    </div>
                  </div>
                </div>
              </li>
            </ul>
          </div>
          <template #footer>
            <div class="text-center">
              <a href="#" class="text-sm font-medium text-primary-600 hover:text-primary-500">Lihat semua aktivitas</a>
            </div>
          </template>
        </UiCard>

        <!-- Top Overdue -->
        <UiCard title="Top Tunggakan (Overdue)">
          <div class="overflow-hidden">
            <ul role="list" class="divide-y divide-gray-200">
              <li v-for="item in overdueItems" :key="item.id" class="py-4 flex items-center justify-between">
                <div class="flex items-center">
                  <div class="h-10 w-10 rounded-full bg-gray-200 flex items-center justify-center text-gray-500 font-bold">
                    {{ item.unit }}
                  </div>
                  <div class="ml-3">
                    <p class="text-sm font-medium text-gray-900">{{ item.name }}</p>
                    <p class="text-sm text-gray-500">{{ item.days }} hari terlambat</p>
                  </div>
                </div>
                <div class="text-right">
                  <p class="text-sm font-semibold text-gray-900">{{ item.amount }}</p>
                  <UiBadge variant="danger">Overdue</UiBadge>
                </div>
              </li>
            </ul>
          </div>
          <template #footer>
            <div class="text-center">
              <a href="#" class="text-sm font-medium text-primary-600 hover:text-primary-500">Lihat semua tunggakan</a>
            </div>
          </template>
        </UiCard>
      </div>
    </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

// Mock data for UI demo
const activities = [
  {
    id: 1,
    content: 'Pembayaran diterima dari',
    target: 'Unit A-02 (Budi)',
    date: 'Baru saja',
    datetime: '2025-01-29T10:00',
    icon: 'div', // Placeholder for icon component
    iconBackground: 'bg-green-500',
  },
  {
    id: 2,
    content: 'Tagihan baru dibuat untuk',
    target: 'Periode Februari 2025',
    date: '1 jam yang lalu',
    datetime: '2025-01-29T09:00',
    icon: 'div',
    iconBackground: 'bg-blue-500',
  },
  {
    id: 3,
    content: 'Pengumuman dipublish',
    target: 'Kerja Bakti Minggu Ini',
    date: '2 jam yang lalu',
    datetime: '2025-01-29T08:00',
    icon: 'div',
    iconBackground: 'bg-yellow-500',
  },
]

const overdueItems = [
  { id: 1, unit: 'A-05', name: 'Jane Doe', days: 15, amount: 'Rp 150.000' },
  { id: 2, unit: 'B-12', name: 'John Doe', days: 10, amount: 'Rp 100.000' },
  { id: 3, unit: 'C-08', name: 'Bob Smith', days: 8, amount: 'Rp 75.000' },
  { id: 4, unit: 'D-01', name: 'Alice', days: 5, amount: 'Rp 50.000' },
]
</script>
