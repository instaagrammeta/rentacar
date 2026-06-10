<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { clientApi, uploadApi } from '@/api'
import { useAppStore } from '@/stores/app'

const { t } = useI18n()
const router = useRouter()
const app = useAppStore()

const items = ref([])
const total = ref(0)
const loading = ref(false)
const search = ref('')
const dialog = ref(false)
const editing = ref(null)
const form = ref({})
const saving = ref(false)

const statusOptions = [
  { value: 'active', title: t('client.statuses.active') },
  { value: 'pending_verification', title: t('client.statuses.pending_verification') },
  { value: 'blacklisted', title: t('client.statuses.blacklisted') },
]

const headers = [
  { title: t('client.code'), key: 'client_code' },
  { title: t('client.lastName'), key: 'full_name' },
  { title: t('client.phone'), key: 'phone' },
  { title: t('client.experience'), key: 'driver_experience_years' },
  { title: t('app.status'), key: 'status_label' },
  { title: 'VIP', key: 'is_vip' },
  { title: t('app.actions'), key: 'actions', sortable: false },
]

async function load() {
  loading.value = true
  try {
    const { data } = await clientApi.list({ search: search.value, per_page: 100 })
    items.value = data.items
    total.value = data.total
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editing.value = null
  form.value = { status: 'pending_verification', driver_experience_years: 0, is_vip: false }
  dialog.value = true
}

function openEdit(item) {
  editing.value = item
  form.value = { ...item }
  dialog.value = true
}

async function save() {
  saving.value = true
  try {
    if (editing.value) {
      await clientApi.update(editing.value.id, form.value)
      app.notify(t('app.saved'))
    } else {
      await clientApi.create(form.value)
      app.notify(t('app.created'))
    }
    dialog.value = false
    await load()
  } finally {
    saving.value = false
  }
}

async function remove(item) {
  if (!confirm(t('app.confirmDelete'))) return
  await clientApi.remove(item.id)
  app.notify(t('app.deleted'))
  await load()
}

async function onUpload(event, field) {
  const file = event.target.files?.[0]
  if (!file) return
  const fd = new FormData()
  fd.append('file', file)
  fd.append('subfolder', 'clients')
  const { data } = await uploadApi.upload(fd)
  form.value[field] = data.paths[0]
  app.notify(t('app.saved'))
}

onMounted(load)
</script>

<template>
  <div>
    <div class="d-flex align-center mb-4">
      <h1 class="text-h4 font-weight-bold">{{ t('menu.clients') }}</h1>
      <v-spacer />
      <v-btn color="primary" prepend-icon="mdi-plus" @click="openCreate">{{ t('client.addClient') }}</v-btn>
    </div>

    <v-card class="pa-4">
      <v-text-field
        v-model="search"
        :label="t('client.quickSearch')"
        prepend-inner-icon="mdi-magnify"
        clearable
        class="mb-2"
        @update:model-value="load"
      />
      <v-data-table :headers="headers" :items="items" :loading="loading" items-per-page="20">
        <template #[`item.is_vip`]="{ item }">
          <v-icon v-if="item.is_vip" color="amber" icon="mdi-star" />
        </template>
        <template #[`item.status_label`]="{ item }">
          <v-chip size="small" :color="item.status === 'blacklisted' ? 'error' : item.status === 'active' ? 'success' : 'warning'">
            {{ item.status_label }}
          </v-chip>
        </template>
        <template #[`item.actions`]="{ item }">
          <v-btn icon="mdi-history" size="small" variant="text" @click="router.push(`/clients/${item.id}`)" />
          <v-btn icon="mdi-pencil" size="small" variant="text" @click="openEdit(item)" />
          <v-btn icon="mdi-delete" size="small" variant="text" color="error" @click="remove(item)" />
        </template>
      </v-data-table>
    </v-card>

    <v-dialog v-model="dialog" max-width="800" scrollable>
      <v-card>
        <v-card-title class="bg-primary text-white">
          {{ editing ? t('client.editClient') : t('client.addClient') }}
        </v-card-title>
        <v-card-text class="pt-4">
          <v-row>
            <v-col cols="12" md="6"><v-text-field v-model="form.first_name" :label="t('client.firstName')" /></v-col>
            <v-col cols="12" md="6"><v-text-field v-model="form.last_name" :label="t('client.lastName')" /></v-col>
            <v-col cols="12" md="6"><v-text-field v-model="form.phone" :label="t('client.phone')" /></v-col>
            <v-col cols="12" md="6"><v-text-field v-model="form.email" :label="t('client.email')" /></v-col>
            <v-col cols="12" md="6"><v-text-field v-model="form.date_of_birth" :label="t('client.dateOfBirth')" type="date" /></v-col>
            <v-col cols="12" md="6"><v-text-field v-model="form.passport_number" :label="t('client.passportNumber')" /></v-col>
            <v-col cols="12" md="6"><v-text-field v-model="form.driver_license_number" :label="t('client.licenseNumber')" /></v-col>
            <v-col cols="12" md="6"><v-text-field v-model="form.driver_license_issue_date" :label="t('client.licenseIssueDate')" type="date" /></v-col>
            <v-col cols="12" md="6"><v-text-field v-model.number="form.driver_experience_years" :label="t('client.experience')" type="number" /></v-col>
            <v-col cols="12" md="6"><v-select v-model="form.status" :items="statusOptions" :label="t('app.status')" /></v-col>
            <v-col cols="12" md="6">
              <v-file-input :label="t('client.passportScan')" accept="image/*,.pdf" @change="(e) => onUpload(e, 'passport_scan')" />
            </v-col>
            <v-col cols="12" md="6">
              <v-file-input :label="t('client.licenseScan')" accept="image/*,.pdf" @change="(e) => onUpload(e, 'driver_license_scan')" />
            </v-col>
            <v-col cols="12" md="6"><v-switch v-model="form.is_vip" :label="t('client.vip')" color="amber" /></v-col>
            <v-col cols="12"><v-textarea v-model="form.notes" :label="t('app.notes')" rows="2" /></v-col>
          </v-row>
          <v-alert type="info" variant="tonal" density="compact">
            Минимальный возраст — 21 год, минимальный стаж — 1 год. Иначе аренда невозможна.
          </v-alert>
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
