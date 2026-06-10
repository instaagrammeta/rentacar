<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { carApi, uploadApi } from '@/api'
import { useAppStore } from '@/stores/app'

const { t } = useI18n()
const app = useAppStore()

const items = ref([])
const loading = ref(false)
const search = ref('')
const statusFilter = ref(null)
const dialog = ref(false)
const editing = ref(null)
const form = ref({})
const saving = ref(false)

const statusOptions = [
  { value: 'available', title: t('car.statuses.available') },
  { value: 'reserved', title: t('car.statuses.reserved') },
  { value: 'rented', title: t('car.statuses.rented') },
  { value: 'maintenance', title: t('car.statuses.maintenance') },
]

const headers = [
  { title: t('car.brand'), key: 'brand' },
  { title: t('car.model'), key: 'model' },
  { title: t('car.year'), key: 'year' },
  { title: t('car.plate'), key: 'plate_number' },
  { title: t('car.dailyPrice'), key: 'daily_price' },
  { title: t('app.status'), key: 'status_label' },
  { title: t('app.actions'), key: 'actions', sortable: false },
]

function statusColor(status) {
  return { available: 'success', reserved: 'warning', rented: 'info', maintenance: 'error' }[status]
}

async function load() {
  loading.value = true
  try {
    const { data } = await carApi.list({ search: search.value, status: statusFilter.value || '', per_page: 100 })
    items.value = data.items
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editing.value = null
  form.value = { status: 'available', year: new Date().getFullYear(), photos: [] }
  dialog.value = true
}

function openEdit(item) {
  editing.value = item
  form.value = { ...item, photos: [...(item.photos || [])] }
  dialog.value = true
}

async function save() {
  saving.value = true
  try {
    if (editing.value) {
      await carApi.update(editing.value.id, form.value)
    } else {
      await carApi.create(form.value)
    }
    app.notify(t('app.saved'))
    dialog.value = false
    await load()
  } finally {
    saving.value = false
  }
}

async function remove(item) {
  if (!confirm(t('app.confirmDelete'))) return
  await carApi.remove(item.id)
  app.notify(t('app.deleted'))
  await load()
}

async function onUploadPhotos(event) {
  const files = Array.from(event.target.files || [])
  if (!files.length) return
  const fd = new FormData()
  files.forEach((f) => fd.append('files', f))
  fd.append('subfolder', 'cars')
  const { data } = await uploadApi.upload(fd)
  form.value.photos = [...(form.value.photos || []), ...data.paths]
  app.notify(t('app.saved'))
}

function fileUrl(path) {
  return uploadApi.fileUrl(path)
}

onMounted(load)
</script>

<template>
  <div>
    <div class="d-flex align-center mb-4">
      <h1 class="text-h4 font-weight-bold">{{ t('menu.cars') }}</h1>
      <v-spacer />
      <v-btn color="primary" prepend-icon="mdi-plus" @click="openCreate">{{ t('car.addCar') }}</v-btn>
    </div>

    <v-card class="pa-4">
      <v-row class="mb-2">
        <v-col cols="12" md="8"><v-text-field v-model="search" :label="t('app.search')" prepend-inner-icon="mdi-magnify" clearable @update:model-value="load" /></v-col>
        <v-col cols="12" md="4"><v-select v-model="statusFilter" :items="statusOptions" :label="t('app.status')" clearable @update:model-value="load" /></v-col>
      </v-row>
      <v-data-table :headers="headers" :items="items" :loading="loading" items-per-page="20">
        <template #[`item.daily_price`]="{ item }">{{ item.daily_price }} ₽</template>
        <template #[`item.status_label`]="{ item }">
          <v-chip size="small" :color="statusColor(item.status)">{{ item.status_label }}</v-chip>
        </template>
        <template #[`item.actions`]="{ item }">
          <v-btn icon="mdi-pencil" size="small" variant="text" @click="openEdit(item)" />
          <v-btn icon="mdi-delete" size="small" variant="text" color="error" @click="remove(item)" />
        </template>
      </v-data-table>
    </v-card>

    <v-dialog v-model="dialog" max-width="800" scrollable>
      <v-card>
        <v-card-title class="bg-primary text-white">{{ editing ? t('car.editCar') : t('car.addCar') }}</v-card-title>
        <v-card-text class="pt-4">
          <v-row>
            <v-col cols="12" md="6"><v-text-field v-model="form.brand" :label="t('car.brand')" /></v-col>
            <v-col cols="12" md="6"><v-text-field v-model="form.model" :label="t('car.model')" /></v-col>
            <v-col cols="6" md="3"><v-text-field v-model.number="form.year" :label="t('car.year')" type="number" /></v-col>
            <v-col cols="6" md="3"><v-text-field v-model="form.color" :label="t('car.color')" /></v-col>
            <v-col cols="12" md="6"><v-text-field v-model="form.vin" :label="t('car.vin')" /></v-col>
            <v-col cols="12" md="6"><v-text-field v-model="form.plate_number" :label="t('car.plate')" /></v-col>
            <v-col cols="12" md="6"><v-text-field v-model.number="form.mileage" :label="t('car.mileage')" type="number" /></v-col>
            <v-col cols="6" md="3"><v-text-field v-model.number="form.daily_price" :label="t('car.dailyPrice')" type="number" /></v-col>
            <v-col cols="6" md="3"><v-text-field v-model.number="form.weekly_price" :label="t('car.weeklyPrice')" type="number" /></v-col>
            <v-col cols="6" md="3"><v-text-field v-model.number="form.monthly_price" :label="t('car.monthlyPrice')" type="number" /></v-col>
            <v-col cols="6" md="3"><v-text-field v-model.number="form.deposit_amount" :label="t('car.deposit')" type="number" /></v-col>
            <v-col cols="6" md="3"><v-select v-model="form.status" :items="statusOptions" :label="t('app.status')" /></v-col>
            <v-col cols="12">
              <v-file-input :label="t('car.photos')" accept="image/*" multiple @change="onUploadPhotos" />
              <div class="d-flex ga-2 flex-wrap">
                <v-img v-for="(p, i) in form.photos" :key="i" :src="fileUrl(p)" max-width="100" max-height="70" cover class="rounded" />
              </div>
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="dialog = false">{{ t('app.cancel') }}</v-btn>
          <v-btn color="primary" :loading="saving" @click="save">{{ t('app.save') }}</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>
