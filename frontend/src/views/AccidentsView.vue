<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { accidentApi, carApi, clientApi, uploadApi } from '@/api'
import { useAppStore } from '@/stores/app'

const { t } = useI18n()
const app = useAppStore()

const items = ref([])
const cars = ref([])
const clients = ref([])
const loading = ref(false)
const dialog = ref(false)
const editing = ref(null)
const form = ref({})
const saving = ref(false)

const headers = [
  { title: t('accident.date'), key: 'accident_date' },
  { title: t('accident.car'), key: 'car_name' },
  { title: t('accident.client'), key: 'client_name' },
  { title: t('accident.repairCost'), key: 'repair_cost' },
  { title: t('app.actions'), key: 'actions', sortable: false },
]

async function load() {
  loading.value = true
  try {
    const { data } = await accidentApi.list({ per_page: 100 })
    items.value = data.items
  } finally {
    loading.value = false
  }
}

async function openCreate() {
  editing.value = null
  form.value = { accident_date: new Date().toISOString().slice(0, 10), repair_cost: 0, photos: [] }
  await loadOptions()
  dialog.value = true
}

function openEdit(item) {
  editing.value = item
  form.value = { ...item, photos: [...(item.photos || [])] }
  loadOptions()
  dialog.value = true
}

async function loadOptions() {
  const [c, cl] = await Promise.all([carApi.list({ per_page: 200 }), clientApi.list({ per_page: 200 })])
  cars.value = c.data.items.map((x) => ({ value: x.id, title: x.display_name }))
  clients.value = cl.data.items.map((x) => ({ value: x.id, title: `${x.full_name} (${x.client_code})` }))
}

async function save() {
  saving.value = true
  try {
    if (editing.value) await accidentApi.update(editing.value.id, form.value)
    else await accidentApi.create(form.value)
    app.notify(t('app.saved'))
    dialog.value = false
    await load()
  } finally {
    saving.value = false
  }
}

async function remove(item) {
  if (!confirm(t('app.confirmDelete'))) return
  await accidentApi.remove(item.id)
  app.notify(t('app.deleted'))
  await load()
}

async function onUploadPhotos(event) {
  const files = Array.from(event.target.files || [])
  if (!files.length) return
  const fd = new FormData()
  files.forEach((f) => fd.append('files', f))
  fd.append('subfolder', 'accidents')
  const { data } = await uploadApi.upload(fd)
  form.value.photos = [...(form.value.photos || []), ...data.paths]
}

function fileUrl(path) {
  return uploadApi.fileUrl(path)
}

onMounted(load)
</script>

<template>
  <div>
    <div class="d-flex align-center mb-4">
      <h1 class="text-h4 font-weight-bold">{{ t('accident.title') }}</h1>
      <v-spacer />
      <v-btn color="primary" prepend-icon="mdi-plus" @click="openCreate">{{ t('accident.add') }}</v-btn>
    </div>

    <v-card class="pa-4">
      <v-data-table :headers="headers" :items="items" :loading="loading" items-per-page="20">
        <template #[`item.repair_cost`]="{ item }">{{ item.repair_cost }} ₽</template>
        <template #[`item.actions`]="{ item }">
          <v-btn icon="mdi-pencil" size="small" variant="text" @click="openEdit(item)" />
          <v-btn icon="mdi-delete" size="small" variant="text" color="error" @click="remove(item)" />
        </template>
      </v-data-table>
    </v-card>

    <v-dialog v-model="dialog" max-width="700">
      <v-card>
        <v-card-title class="bg-primary text-white">{{ t('accident.add') }}</v-card-title>
        <v-card-text class="pt-4">
          <v-select v-model="form.car_id" :items="cars" :label="t('accident.car')" />
          <v-select v-model="form.client_id" :items="clients" :label="t('accident.client')" clearable />
          <v-text-field v-model="form.accident_date" :label="t('accident.date')" type="date" />
          <v-text-field v-model.number="form.repair_cost" :label="t('accident.repairCost')" type="number" />
          <v-textarea v-model="form.description" :label="t('accident.description')" rows="3" />
          <v-file-input :label="t('accident.photos')" accept="image/*" multiple @change="onUploadPhotos" />
          <div class="d-flex ga-2 flex-wrap">
            <v-img v-for="(p, i) in form.photos" :key="i" :src="fileUrl(p)" max-width="100" max-height="70" cover class="rounded" />
          </div>
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
