<template>
  <div class="max-w-3xl mx-auto">
      <!-- Header -->
      <div class="mb-8">
        <NuxtLink to="/dashboard/billing" class="text-sm text-gray-500 hover:text-gray-900 flex items-center mb-4">
          <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path></svg>
          Kembali ke Tagihan
        </NuxtLink>
        <h1 class="text-2xl font-bold text-gray-900 tracking-tight">Buat Tagihan Baru</h1>
        <p class="text-gray-500">Generate tagihan massal untuk warga berdasarkan template.</p>
      </div>

      <!-- Wizard Steps -->
      <div class="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden">
        <!-- Progress Bar -->
        <div class="bg-gray-50 border-b border-gray-100 px-8 py-4">
          <div class="flex items-center justify-between">
            <div class="flex items-center">
              <div :class="['w-8 h-8 rounded-full flex items-center justify-center text-sm font-bold', step >= 1 ? 'bg-primary-600 text-white' : 'bg-gray-200 text-gray-500']">1</div>
              <span :class="['ml-3 text-sm font-medium', step >= 1 ? 'text-gray-900' : 'text-gray-500']">Pilih Template</span>
            </div>
            <div class="h-0.5 w-16 bg-gray-200 mx-4"></div>
            <div class="flex items-center">
              <div :class="['w-8 h-8 rounded-full flex items-center justify-center text-sm font-bold', step >= 2 ? 'bg-primary-600 text-white' : 'bg-gray-200 text-gray-500']">2</div>
              <span :class="['ml-3 text-sm font-medium', step >= 2 ? 'text-gray-900' : 'text-gray-500']">Detail & Periode</span>
            </div>
            <div class="h-0.5 w-16 bg-gray-200 mx-4"></div>
            <div class="flex items-center">
              <div :class="['w-8 h-8 rounded-full flex items-center justify-center text-sm font-bold', step >= 3 ? 'bg-primary-600 text-white' : 'bg-gray-200 text-gray-500']">3</div>
              <span :class="['ml-3 text-sm font-medium', step >= 3 ? 'text-gray-900' : 'text-gray-500']">Konfirmasi</span>
            </div>
          </div>
        </div>

        <div class="p-8">
          <!-- Step 1: Select Template -->
          <div v-if="step === 1" class="space-y-6">
            <div class="flex justify-between items-center">
              <h2 class="text-lg font-semibold text-gray-900">Pilih Template Tagihan</h2>
              <button
                @click="showCreateTemplateModal = true"
                class="px-4 py-2 text-sm bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors flex items-center gap-2"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path></svg>
                Template Baru
              </button>
            </div>
            <div v-if="loadingTemplates" class="text-center py-8 text-gray-500">
              Memuat template...
            </div>
            <div v-else-if="templates.length === 0" class="text-center py-8 text-gray-500">
              Belum ada template. Buat template baru untuk memulai.
            </div>
            <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div 
                v-for="template in templates" 
                :key="template.id"
                @click="selectedTemplate = template"
                :class="['cursor-pointer p-4 rounded-xl border-2 transition-all relative', selectedTemplate?.id === template.id ? 'border-primary-500 bg-primary-50' : 'border-gray-100 hover:border-gray-200']"
              >
                <div v-if="!template.is_system" class="absolute top-2 right-2 flex gap-1">
                  <button
                    @click.stop="editTemplate(template)"
                    class="p-1 text-gray-400 hover:text-primary-600 transition-colors"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path></svg>
                  </button>
                  <button
                    @click.stop="deleteTemplate(template)"
                    class="p-1 text-gray-400 hover:text-red-600 transition-colors"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
                  </button>
                </div>
                <div class="flex justify-between items-start mb-2">
                  <div class="p-2 bg-white rounded-lg shadow-sm text-primary-600">
                    <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"></path></svg>
                  </div>
                  <span class="px-2 py-1 bg-gray-100 text-gray-600 text-xs font-medium rounded">{{ template.type }}</span>
                </div>
                <h3 class="font-bold text-gray-900">{{ template.name }}</h3>
                <p class="text-sm text-gray-500 mt-1">{{ template.description || '-' }}</p>
                <div class="mt-2 flex items-center gap-2">
                  <span v-if="template.is_active" class="px-2 py-0.5 bg-green-100 text-green-700 text-xs font-medium rounded">Aktif</span>
                  <span v-else class="px-2 py-0.5 bg-gray-100 text-gray-600 text-xs font-medium rounded">Tidak Aktif</span>
                  <span v-if="template.recurring_type === 'monthly'" class="px-2 py-0.5 bg-blue-100 text-blue-700 text-xs font-medium rounded">Bulanan</span>
                  <span v-else-if="template.recurring_type === 'yearly'" class="px-2 py-0.5 bg-purple-100 text-purple-700 text-xs font-medium rounded">Tahunan</span>
                </div>
                <div class="mt-4 text-sm font-medium text-gray-900">
                  <span v-if="template.amount_rules && template.amount_rules.length > 0">
                    Tarif: 
                    <span v-for="(rule, idx) in template.amount_rules" :key="idx" class="ml-1">
                      {{ rule.unit_type }}: Rp {{ rule.amount.toLocaleString('id-ID') }}
                      <span v-if="idx < template.amount_rules.length - 1">, </span>
                    </span>
                  </span>
                  <span v-else>
                    Mulai Rp {{ (template.amount || 0).toLocaleString('id-ID') }}
                  </span>
                </div>
                <div v-if="template.due_day" class="mt-2 text-xs text-gray-500">
                  Jatuh tempo: Tanggal {{ template.due_day }} setiap periode
                </div>
              </div>
            </div>
          </div>

          <!-- Step 2: Details -->
          <div v-if="step === 2" class="space-y-6">
            <h2 class="text-lg font-semibold text-gray-900">Detail Tagihan</h2>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Kategori Tagihan</label>
              <input 
                v-model="formData.category"
                type="text" 
                class="block w-full border-gray-300 rounded-xl shadow-sm focus:ring-primary-500 focus:border-primary-500 py-3"
                placeholder="Contoh: IPL Bulanan"
              >
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Periode Tagihan</label>
              <!-- For Bulanan template: show month selector -->
              <div v-if="selectedTemplate?.type === 'Bulanan'" class="space-y-3">
                <div class="flex items-center gap-2 mb-2">
                  <input
                    type="checkbox"
                    id="selectAllMonths"
                    :checked="selectedMonths.length === 12"
                    @change="toggleAllMonths"
                    class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
                  >
                  <label for="selectAllMonths" class="text-sm font-medium text-gray-700 cursor-pointer">
                    Pilih Semua Bulan
                  </label>
                </div>
                <div class="grid grid-cols-3 gap-2 max-h-48 overflow-y-auto border border-gray-200 rounded-lg p-3">
                  <label
                    v-for="(month, index) in months"
                    :key="index"
                    class="flex items-center p-2 hover:bg-gray-50 rounded cursor-pointer"
                  >
                    <input
                      v-model="selectedMonths"
                      type="checkbox"
                      :value="index + 1"
                      class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
                    >
                    <span class="ml-2 text-sm text-gray-700">{{ month }}</span>
                  </label>
                </div>
                <div v-if="selectedMonths.length === 0" class="text-sm text-red-500">
                  Pilih minimal satu bulan
                </div>
              </div>
              <!-- For other template types: show text input -->
              <input 
                v-else
                v-model="formData.period"
                type="text" 
                class="block w-full border-gray-300 rounded-xl shadow-sm focus:ring-primary-500 focus:border-primary-500 py-3"
                placeholder="Contoh: Januari 2025"
              >
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Jumlah Tagihan</label>
              <input 
                v-model.number="formData.amount"
                type="number" 
                class="block w-full border-gray-300 rounded-xl shadow-sm focus:ring-primary-500 focus:border-primary-500 py-3"
                placeholder="150000"
                min="0"
              >
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Denda Keterlambatan (Opsional)</label>
              <input 
                v-model.number="formData.lateFee"
                type="number" 
                class="block w-full border-gray-300 rounded-xl shadow-sm focus:ring-primary-500 focus:border-primary-500 py-3"
                placeholder="25000"
                min="0"
              >
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Jatuh Tempo
                <span class="text-gray-400 text-xs font-normal">(Opsional)</span>
              </label>
              <!-- For Bulanan template: show day selector -->
              <div v-if="selectedTemplate?.type === 'Bulanan'" class="space-y-2">
                <div class="flex items-center gap-3">
                  <span class="text-sm text-gray-600">Tanggal setiap bulan:</span>
                  <select
                    v-model.number="dueDateDay"
                    class="border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500 py-2 px-3"
                  >
                    <option :value="null">Pilih tanggal</option>
                    <option v-for="day in 31" :key="day" :value="day">{{ day }}</option>
                  </select>
                </div>
                <p class="text-xs text-gray-500">
                  Tagihan akan jatuh tempo pada tanggal {{ dueDateDay || '-' }} setiap bulan
                </p>
              </div>
              <!-- For other template types: show date picker -->
              <input 
                v-else
                v-model="formData.dueDate"
                type="date" 
                class="block w-full border-gray-300 rounded-xl shadow-sm focus:ring-primary-500 focus:border-primary-500 py-3"
              >
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Catatan (Opsional)</label>
              <textarea 
                v-model="formData.notes"
                rows="3"
                class="block w-full border-gray-300 rounded-xl shadow-sm focus:ring-primary-500 focus:border-primary-500 py-3"
                placeholder="Catatan tambahan untuk tagihan..."
              ></textarea>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Target Unit</label>
              <div class="space-y-2">
                <label class="flex items-center p-3 border border-gray-200 rounded-xl cursor-pointer hover:bg-gray-50">
                  <input 
                    v-model="targetMode" 
                    type="radio" 
                    value="all" 
                    class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300"
                    @change="loadUnits"
                  >
                  <span class="ml-3 block text-sm font-medium text-gray-700">
                    Semua Unit Aktif ({{ allUnits.length }} Unit)
                  </span>
                </label>
                <label class="flex items-center p-3 border border-gray-200 rounded-xl cursor-pointer hover:bg-gray-50">
                  <input 
                    v-model="targetMode" 
                    type="radio" 
                    value="manual" 
                    class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300"
                    @change="loadUnits"
                  >
                  <span class="ml-3 block text-sm font-medium text-gray-700">Pilih Manual</span>
                </label>
              </div>

              <!-- Manual Unit Selection -->
              <div v-if="targetMode === 'manual'" class="mt-4 max-h-64 overflow-y-auto border border-gray-200 rounded-xl p-4 space-y-2">
                <div v-if="loadingUnits" class="text-center text-gray-500 py-4">
                  Memuat unit...
                </div>
                <div v-else-if="allUnits.length === 0" class="text-center text-gray-500 py-4">
                  Belum ada unit
                </div>
                <label 
                  v-else
                  v-for="unit in allUnits" 
                  :key="unit.id"
                  class="flex items-center p-2 hover:bg-gray-50 rounded cursor-pointer"
                >
                  <input 
                    v-model="selectedUnitIds" 
                    type="checkbox" 
                    :value="unit.id"
                    class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
                  >
                  <span class="ml-3 text-sm text-gray-700">
                    {{ unit.code }} - {{ unit.type }} 
                    <span v-if="unit.owner_name" class="text-gray-500">({{ unit.owner_name }})</span>
                  </span>
                </label>
              </div>
            </div>
          </div>

          <!-- Step 3: Confirmation -->
          <div v-if="step === 3" class="space-y-6">
            <div class="bg-green-50 border border-green-100 rounded-xl p-4 flex items-start">
              <svg class="w-6 h-6 text-green-500 mt-0.5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
              <div>
                <h3 class="text-sm font-medium text-green-800">Siap Generate Tagihan</h3>
                <p class="text-sm text-green-700 mt-1">Sistem akan membuat 150 tagihan baru dan mengirimkan notifikasi WhatsApp ke warga.</p>
              </div>
            </div>

            <div class="bg-gray-50 rounded-xl p-6 space-y-4">
              <div class="flex justify-between">
                <span class="text-gray-500">Template</span>
                <span class="font-medium text-gray-900">{{ selectedTemplate?.name }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-500">Kategori</span>
                <span class="font-medium text-gray-900">{{ formData.category || '-' }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-500">Periode</span>
                <span class="font-medium text-gray-900">
                  <span v-if="selectedTemplate?.type === 'Bulanan'">
                    {{ selectedMonths.length === 12 ? 'Semua Bulan' : `${selectedMonths.length} Bulan Terpilih` }}
                  </span>
                  <span v-else>{{ formData.period || '-' }}</span>
                </span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-500">Jatuh Tempo</span>
                <span class="font-medium text-gray-900">
                  <span v-if="selectedTemplate?.type === 'Bulanan'">
                    {{ dueDateDay ? `Tanggal ${dueDateDay} setiap bulan` : 'Tidak diatur' }}
                  </span>
                  <span v-else>{{ formatDate(formData.dueDate) || 'Tidak diatur' }}</span>
                </span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-500">Jumlah per Unit</span>
                <span class="font-medium text-gray-900">Rp {{ (formData.amount || 0).toLocaleString('id-ID') }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-500">Total Unit</span>
                <span class="font-medium text-gray-900">{{ selectedUnitIds.length || allUnits.length }} Unit</span>
              </div>
              <div v-if="selectedTemplate?.type === 'Bulanan'" class="flex justify-between">
                <span class="text-gray-500">Jumlah Periode</span>
                <span class="font-medium text-gray-900">{{ selectedMonths.length }} Bulan</span>
              </div>
              <div class="pt-4 border-t border-gray-200 flex justify-between">
                <span class="text-gray-500">Estimasi Total</span>
                <span class="font-bold text-lg text-primary-600">
                  Rp {{ ((formData.amount || 0) * (selectedUnitIds.length || allUnits.length) * (selectedTemplate?.type === 'Bulanan' ? selectedMonths.length : 1)).toLocaleString('id-ID') }}
                </span>
              </div>
            </div>
          </div>

          <!-- Actions -->
          <div class="mt-8 flex justify-between pt-6 border-t border-gray-100">
            <button 
              v-if="step > 1"
              @click="step--"
              class="px-6 py-2.5 border border-gray-300 text-gray-700 font-medium rounded-xl hover:bg-gray-50 transition-colors"
            >
              Kembali
            </button>
            <div v-else></div>

            <button 
              v-if="step < 3"
              @click="handleNext"
              :disabled="!canProceed"
              :class="['px-6 py-2.5 bg-primary-600 text-white font-medium rounded-xl transition-colors shadow-lg shadow-primary-600/20', !canProceed ? 'opacity-50 cursor-not-allowed' : 'hover:bg-primary-700']"
            >
              Lanjut
            </button>
            <button 
              v-else
              @click="generateBills"
              :disabled="generating"
              :class="['px-6 py-2.5 bg-primary-600 text-white font-medium rounded-xl transition-colors shadow-lg shadow-primary-600/20', generating ? 'opacity-50 cursor-not-allowed' : 'hover:bg-primary-700']"
            >
              {{ generating ? 'Membuat Tagihan...' : 'Generate Tagihan' }}
            </button>
          </div>
        </div>
      </div>

      <!-- Create/Edit Template Modal -->
      <div v-if="showCreateTemplateModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click.self="showCreateTemplateModal = false">
        <div class="bg-white rounded-2xl p-6 max-w-md w-full mx-4 max-h-[90vh] overflow-y-auto">
          <h3 class="text-lg font-bold text-gray-900 mb-4">
            {{ editingTemplate ? 'Edit Template' : 'Template Baru' }}
          </h3>
          
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Nama Template *</label>
              <input
                v-model="templateForm.name"
                type="text"
                class="w-full border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                placeholder="Contoh: IPL Bulanan"
              >
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Kategori *</label>
              <input
                v-model="templateForm.category"
                type="text"
                class="w-full border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                placeholder="Contoh: IPL Bulanan"
              >
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Tipe *</label>
              <select
                v-model="templateForm.type"
                class="w-full border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
              >
                <option value="Bulanan">Bulanan</option>
                <option value="Tahunan">Tahunan</option>
                <option value="One-time">One-time</option>
              </select>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Deskripsi</label>
              <textarea
                v-model="templateForm.description"
                rows="2"
                class="w-full border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                placeholder="Deskripsi template..."
              ></textarea>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Jumlah Tagihan *</label>
              <input
                v-model.number="templateForm.amount"
                type="number"
                class="w-full border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                placeholder="150000"
                min="0"
              >
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Tipe Recurring</label>
              <select
                v-model="templateForm.recurringType"
                class="w-full border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
              >
                <option value="one-time">One-time</option>
                <option value="monthly">Monthly (Bulanan)</option>
                <option value="yearly">Yearly (Tahunan)</option>
              </select>
            </div>

            <div v-if="templateForm.recurringType !== 'one-time'">
              <label class="block text-sm font-medium text-gray-700 mb-2">Tanggal Jatuh Tempo (1-31)</label>
              <input
                v-model.number="templateForm.dueDay"
                type="number"
                class="w-full border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                placeholder="5"
                min="1"
                max="31"
              >
              <p class="text-xs text-gray-500 mt-1">Tagihan akan jatuh tempo pada tanggal ini setiap periode</p>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Denda Keterlambatan (Per Hari)</label>
              <div class="space-y-3">
                <div>
                  <label class="block text-xs text-gray-600 mb-1">Tipe Denda</label>
                  <select
                    v-model="templateForm.lateFeeType"
                    class="w-full border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500 text-sm"
                  >
                    <option value="fixed">Fixed (Nominal per hari)</option>
                    <option value="percentage">Percentage (% dari tagihan per hari)</option>
                  </select>
                </div>
                <div v-if="templateForm.lateFeeType === 'fixed'">
                  <input
                    v-model.number="templateForm.lateFee"
                    type="number"
                    class="w-full border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                    placeholder="5000"
                    min="0"
                  >
                  <p class="text-xs text-gray-500 mt-1">Contoh: Rp 5.000 per hari keterlambatan</p>
                </div>
                <div v-else>
                  <div class="space-y-2">
                    <input
                      v-model.number="templateForm.lateFeePercentage"
                      type="number"
                      class="w-full border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                      placeholder="0.5"
                      min="0"
                      max="100"
                      step="0.1"
                    >
                    <p class="text-xs text-gray-500">Contoh: 0.5% dari tagihan per hari</p>
                    <div v-if="templateForm.lateFeePercentage">
                      <label class="block text-xs text-gray-600 mb-1">Maksimal Denda (Opsional)</label>
                      <input
                        v-model.number="templateForm.lateFeeMax"
                        type="number"
                        class="w-full border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                        placeholder="100000"
                        min="0"
                      >
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Tarif per Tipe Unit (Opsional)
                <span class="text-gray-400 text-xs font-normal">- Kosongkan jika sama untuk semua unit</span>
              </label>
              <div class="space-y-2 border border-gray-200 rounded-lg p-3">
                <div v-for="(rule, index) in templateForm.amountRules" :key="index" class="flex gap-2 items-center">
                  <select
                    v-model="rule.unit_type"
                    class="flex-1 border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500 text-sm"
                  >
                    <option value="rumah">Rumah</option>
                    <option value="ruko">Ruko</option>
                    <option value="kios">Kios</option>
                  </select>
                  <input
                    v-model.number="rule.amount"
                    type="number"
                    class="flex-1 border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500 text-sm"
                    placeholder="Jumlah"
                    min="0"
                  >
                  <button
                    @click="templateForm.amountRules.splice(index, 1)"
                    class="p-2 text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                    </svg>
                  </button>
                </div>
                <button
                  @click="templateForm.amountRules.push({ unit_type: 'rumah', amount: 0 })"
                  class="w-full px-3 py-2 text-sm border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors flex items-center justify-center gap-2"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
                  </svg>
                  Tambah Tarif Unit
                </button>
                <p class="text-xs text-gray-500 mt-2">
                  Jika tidak diisi, akan menggunakan jumlah default di atas (Rp {{ (templateForm.amount || 0).toLocaleString('id-ID') }})
                </p>
              </div>
            </div>

            <div>
              <label class="flex items-center">
                <input
                  v-model="templateForm.isActive"
                  type="checkbox"
                  class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
                >
                <span class="ml-2 text-sm text-gray-700">Template Aktif</span>
              </label>
              <p class="text-xs text-gray-500 mt-1">Template yang tidak aktif tidak dapat digunakan untuk generate tagihan</p>
            </div>
          </div>
          
          <div class="flex gap-3 mt-6">
            <button
              @click="showCreateTemplateModal = false; editingTemplate = null; resetTemplateForm()"
              class="flex-1 px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors"
            >
              Batal
            </button>
            <button
              @click="saveTemplate"
              :disabled="!templateForm.name || !templateForm.category || !templateForm.amount"
              :class="['flex-1 px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors', (!templateForm.name || !templateForm.category || !templateForm.amount) ? 'opacity-50 cursor-not-allowed' : '']"
            >
              Simpan
            </button>
          </div>
        </div>
      </div>
    </div>
</template>

<script setup lang="ts">
const { fetch } = useApi()
const { showSuccess, showError, showWarning } = useToast()
const confirm = useConfirm()

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

const step = ref(1)
const selectedTemplate = ref(null)
const targetMode = ref('all')
const allUnits = ref([])
const selectedUnitIds = ref([])
const loadingUnits = ref(false)
const generating = ref(false)

const formData = ref({
  category: '',
  period: '',
  amount: 0,
  lateFee: 0,
  dueDate: '',
  notes: '',
})

const selectedMonths = ref([])
const dueDateDay = ref(null)
const currentYear = new Date().getFullYear()
const months = [
  'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
  'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
]

const toggleAllMonths = () => {
  if (selectedMonths.value.length === 12) {
    selectedMonths.value = []
  } else {
    selectedMonths.value = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12]
  }
}

const templates = ref([])
const loadingTemplates = ref(false)
const showCreateTemplateModal = ref(false)
const editingTemplate = ref(null)
const templateForm = ref({
  name: '',
  category: '',
  type: 'Bulanan',
  description: '',
  amount: 0,
  lateFee: 0,
  dueDay: null as number | null,
  recurringType: 'one-time' as 'monthly' | 'yearly' | 'one-time',
  lateFeeType: 'fixed' as 'fixed' | 'percentage',
  lateFeePercentage: null as number | null,
  lateFeeMax: null as number | null,
  isActive: true,
  amountRules: [] as Array<{ unit_type: string; amount: number }>,
})

const canProceed = computed(() => {
  if (step.value === 1) {
    return selectedTemplate.value !== null
  }
  if (step.value === 2) {
    const hasPeriod = selectedTemplate.value?.type === 'Bulanan' 
      ? selectedMonths.value.length > 0
      : formData.value.period
    return formData.value.category && 
           hasPeriod && 
           formData.value.amount > 0 &&
           (targetMode.value === 'all' || selectedUnitIds.value.length > 0)
  }
  return true
})

const loadUnits = async () => {
  if (allUnits.value.length === 0) {
    loadingUnits.value = true
    try {
      // Load units with pagination - for template selection, we need all active units
      // Using a reasonable limit, can be increased if needed
      const response = await fetch('/api/units?limit=100&status=active')
      allUnits.value = response.units || []
      // If there are more units, load them
      if (response.pagination && response.pagination.total > response.pagination.limit) {
        const totalPages = response.pagination.total_pages
        for (let page = 2; page <= totalPages; page++) {
          const nextResponse = await fetch(`/api/units?limit=100&page=${page}&status=active`)
          const nextData = await nextResponse.json()
          if (nextData.units) {
            allUnits.value.push(...nextData.units)
          }
        }
      }
    } catch (error) {
      console.error('Failed to load units:', error)
    } finally {
      loadingUnits.value = false
    }
  }
}

const loadTemplates = async () => {
  loadingTemplates.value = true
  try {
    const response = await fetch('/api/billing/templates')
    templates.value = response.templates || []
  } catch (error) {
    console.error('Failed to load templates:', error)
  } finally {
    loadingTemplates.value = false
  }
}

const handleNext = async () => {
  if (step.value === 1 && selectedTemplate.value) {
    // Pre-fill form with template data
    formData.value.category = selectedTemplate.value.category || ''
    formData.value.amount = selectedTemplate.value.amount || 0
    if (selectedTemplate.value.late_fee) {
      formData.value.lateFee = selectedTemplate.value.late_fee
    }
    // If Bulanan template, select all months by default
    if (selectedTemplate.value.type === 'Bulanan') {
      selectedMonths.value = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12]
      dueDateDay.value = null // Reset due date day
    } else {
      selectedMonths.value = []
      dueDateDay.value = null
    }
    await loadUnits()
  }
  if (step.value === 2) {
    // Validate before proceeding
    if (!canProceed.value) {
      showWarning('Harap lengkapi semua field yang wajib diisi')
      return
    }
  }
  step.value++
}

const saveTemplate = async () => {
  try {
    const url = editingTemplate.value 
      ? `/api/billing/templates/${editingTemplate.value.id}`
      : '/api/billing/templates'
    const method = editingTemplate.value ? 'PUT' : 'POST'
    
    // Prepare request body
    const requestBody: any = {
      name: templateForm.value.name,
      category: templateForm.value.category,
      type: templateForm.value.type,
      description: templateForm.value.description || undefined,
      amount: templateForm.value.amount,
      late_fee: templateForm.value.lateFee || undefined,
      due_day: templateForm.value.dueDay || undefined,
      recurring_type: templateForm.value.recurringType,
      late_fee_type: templateForm.value.lateFeeType,
      late_fee_percentage: templateForm.value.lateFeePercentage || undefined,
      late_fee_max: templateForm.value.lateFeeMax || undefined,
      is_active: templateForm.value.isActive,
    }

    // Add amount rules if any
    if (templateForm.value.amountRules.length > 0) {
      requestBody.amount_rules = templateForm.value.amountRules.filter(rule => rule.amount > 0)
    }
    
    const response = await fetch(url, {
      method,
      body: JSON.stringify(requestBody),
    })
    
    await loadTemplates()
    showCreateTemplateModal.value = false
    editingTemplate.value = null
    resetTemplateForm()
    showSuccess('Template berhasil disimpan')
  } catch (error: any) {
    showError(error.message || 'Terjadi kesalahan', 'Gagal menyimpan template')
  }
}

const resetTemplateForm = () => {
  templateForm.value = {
    name: '',
    category: '',
    type: 'Bulanan',
    description: '',
    amount: 0,
    lateFee: 0,
    dueDay: null,
    recurringType: 'one-time',
    lateFeeType: 'fixed',
    lateFeePercentage: null,
    lateFeeMax: null,
    isActive: true,
    amountRules: [],
  }
}

const editTemplate = async (template: any) => {
  editingTemplate.value = template
  // Load full template details including amount rules
  try {
    const response = await fetch(`/api/billing/templates/${template.id}`)
    const fullTemplate = response
    
    templateForm.value = {
      name: fullTemplate.name,
      category: fullTemplate.category,
      type: fullTemplate.type,
      description: fullTemplate.description || '',
      amount: fullTemplate.amount,
      lateFee: fullTemplate.late_fee || 0,
      dueDay: fullTemplate.due_day || null,
      recurringType: fullTemplate.recurring_type || 'one-time',
      lateFeeType: fullTemplate.late_fee_type || 'fixed',
      lateFeePercentage: fullTemplate.late_fee_percentage || null,
      lateFeeMax: fullTemplate.late_fee_max || null,
      isActive: fullTemplate.is_active !== undefined ? fullTemplate.is_active : true,
      amountRules: fullTemplate.amount_rules || [],
    }
    showCreateTemplateModal.value = true
  } catch (error: any) {
    showError(error.message || 'Gagal memuat detail template', 'Error')
  }
}

const deleteTemplate = async (template: any) => {
  const result = await confirm.show(
    `Yakin ingin menghapus template "${template.name}"?`,
    {
      title: 'Hapus Template',
      confirmText: 'Ya, Hapus',
      cancelText: 'Batal',
      type: 'danger'
    }
  )
  
  if (!result) {
    return
  }
  
  try {
    await fetch(`/api/billing/templates/${template.id}`, {
      method: 'DELETE',
    })
    await loadTemplates()
    showSuccess('Template berhasil dihapus')
  } catch (error: any) {
    showError(error.message || 'Terjadi kesalahan', 'Gagal menghapus template')
  }
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'long',
    year: 'numeric'
  })
}

const generateBills = async () => {
  if (!canProceed.value || !selectedTemplate.value) {
    await showWarning('Harap lengkapi semua field yang wajib diisi')
    return
  }

  generating.value = true
  try {
    // Use new generate from template endpoint if template has recurring_type
    if (selectedTemplate.value.recurring_type && selectedTemplate.value.recurring_type !== 'one-time') {
      // Use template generate endpoint
      let totalCreated = 0
      const errors: string[] = []

      if (selectedTemplate.value.type === 'Bulanan' && selectedMonths.value.length > 0) {
        // Generate for each selected month
        for (const monthNum of selectedMonths.value) {
          const year = currentYear
          const month = String(monthNum).padStart(2, '0')
          const period = `${year}-${month}`
          
          try {
            const requestData: any = {
              template_id: selectedTemplate.value.id,
              period: period,
            }

            // Add unit_ids if manual selection
            if (targetMode.value === 'manual' && selectedUnitIds.value.length > 0) {
              requestData.unit_ids = selectedUnitIds.value
            }

            const response = await fetch(`/api/billing/templates/${selectedTemplate.value.id}/generate`, {
              method: 'POST',
              body: JSON.stringify(requestData),
            })

            totalCreated += response.generated_count || 0
            if (response.errors && response.errors.length > 0) {
              errors.push(...response.errors)
            }
          } catch (err: any) {
            errors.push(`Gagal generate untuk periode ${period}: ${err.message || 'Unknown error'}`)
          }
        }
      } else {
        // Single period
        const period = formData.value.period || `${currentYear}-${String(new Date().getMonth() + 1).padStart(2, '0')}`
        
        const requestData: any = {
          template_id: selectedTemplate.value.id,
          period: period,
        }

        if (targetMode.value === 'manual' && selectedUnitIds.value.length > 0) {
          requestData.unit_ids = selectedUnitIds.value
        }

        const response = await fetch(`/api/billing/templates/${selectedTemplate.value.id}/generate`, {
          method: 'POST',
          body: JSON.stringify(requestData),
        })

        totalCreated = response.generated_count || 0
        if (response.errors && response.errors.length > 0) {
          errors.push(...response.errors)
        }
      }

      if (errors.length > 0) {
        showWarning(`Berhasil membuat ${totalCreated} tagihan, namun ada ${errors.length} error.`, 'Peringatan')
      } else {
        showSuccess(
          `Berhasil membuat ${totalCreated} tagihan!`,
          'Tagihan Berhasil Dibuat'
        )
      }
      navigateTo('/dashboard/billing')
      return
    }

    // Fallback to old method for one-time or manual creation
    // Get unit IDs to create bills for
    let unitIds: string[] = []
    if (targetMode.value === 'all') {
      if (allUnits.value.length === 0) {
        await loadUnits()
      }
      unitIds = allUnits.value.map((u: any) => u.id)
    } else {
      unitIds = selectedUnitIds.value
    }

    if (unitIds.length === 0) {
      showWarning('Tidak ada unit yang dipilih')
      generating.value = false
      return
    }

    // Determine periods to generate
    let periods: string[] = []
    if (selectedTemplate.value?.type === 'Bulanan' && selectedMonths.value.length > 0) {
      // Generate period for each selected month (format: YYYY-MM)
      periods = selectedMonths.value.map((month: number) => {
        const monthStr = String(month).padStart(2, '0')
        return `${currentYear}-${monthStr}`
      })
    } else {
      // Use single period from form (convert to YYYY-MM if needed)
      let period = formData.value.period
      // Try to parse if it's in different format
      if (period && !period.match(/^\d{4}-\d{2}$/)) {
        // Assume it's text format, use current month
        const now = new Date()
        period = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`
      }
      periods = [period || `${currentYear}-${String(new Date().getMonth() + 1).padStart(2, '0')}`]
    }

    // Generate bills for each period
    let totalCreated = 0
    for (let i = 0; i < periods.length; i++) {
      const period = periods[i]
      const monthIndex = selectedTemplate.value?.type === 'Bulanan' ? selectedMonths.value[i] - 1 : null
      
      // Calculate due date
      let dueDate: string | null = null
      if (selectedTemplate.value?.type === 'Bulanan' && dueDateDay.value && monthIndex !== null) {
        // Calculate due date based on selected day and month
        const year = currentYear
        const month = monthIndex + 1
        const day = dueDateDay.value
        
        // Handle months with fewer days (e.g., Feb 30 -> Feb 28/29)
        const daysInMonth = new Date(year, month, 0).getDate()
        const finalDay = Math.min(day, daysInMonth)
        
        dueDate = `${year}-${String(month).padStart(2, '0')}-${String(finalDay).padStart(2, '0')}`
      } else if (selectedTemplate.value?.type !== 'Bulanan' && formData.value.dueDate) {
        dueDate = formData.value.dueDate
      }

      const requestData: any = {
        category: formData.value.category,
        period: period,
        amount: formData.value.amount,
        unit_ids: unitIds,
      }

      // Only add due_date if it's set
      if (dueDate) {
        requestData.due_date = dueDate
      }

      if (formData.value.lateFee > 0) {
        requestData.late_fee = formData.value.lateFee
      }

      if (formData.value.notes) {
        requestData.notes = formData.value.notes
      }

      // Call bulk create API for this period
      try {
        const response = await fetch('/api/billing/bulk', {
          method: 'POST',
          body: JSON.stringify(requestData),
        })

        totalCreated += response.created_count || unitIds.length
      } catch (err: any) {
        // If one period fails, log but continue with others
        console.error(`Failed to create bills for period ${period}:`, err)
        throw err // Re-throw to stop the loop
      }
    }

    showSuccess(
      `Berhasil membuat ${totalCreated} tagihan untuk ${periods.length} periode!`,
      'Tagihan Berhasil Dibuat'
    )
    navigateTo('/dashboard/billing')
  } catch (error: any) {
    console.error('Failed to generate bills:', error)
    const errorMessage = error.data?.error || error.message || 'Terjadi kesalahan saat membuat tagihan'
    showError(errorMessage, 'Gagal membuat tagihan')
  } finally {
    generating.value = false
  }
}

// Load templates and units on mount
onMounted(() => {
  loadTemplates()
  loadUnits()
})

// Watch targetMode to load units when changed
watch(targetMode, () => {
  loadUnits()
})
</script>
