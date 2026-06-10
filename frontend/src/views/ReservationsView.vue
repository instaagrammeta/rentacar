<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { reservationApi, clientApi, carApi } from '@/api'
import { useAppStore } from '@/stores/app'

const { t } = useI18n()
const app = useAppStore()

const items = ref([])
const clients = ref([])
const cars = ref([])
const loading = ref(false)
const dialog = ref(false)
const form = ref({})
const saving = ref(false)

const headers = [
  { title: t('reservation.client'), key: 'client_name' },
  { title: t('reservation.car'), key: 'car_name' },
  { title: t('reservation.startDate'), key: 'start_date' },
  { title: t('reservation.endDate'), key: 'end_date' },
  { title: t('reservation.deposit'), key: 'deposit' },
  { title: t('app.status'), key: 'status_label' },
  { title: t('app.actions'), key: 'actions', sortable: false },
]

function statusColor(status) {
  return { reserved: 'warning', confirmed: 'success', cancelled: 'error' }[status]
}

async function load() {
  loading.value = true
  try {
    const { data } = await reservationApi.list({ per_page: 100 })
    items.value = data.items
  } finally {
    loading.value = false
  }
}

async function openCreate() {
  form.value = { start_date: '', end_date: '', deposit: 0 }
  const [c, cr] = await Promise.all([clientApi.list({ per_page: 200 }), carApi.available()])
  clients.value = c.data.items.map((x) => ({ value: x.id, title: `${x.full_name} (${x.client_code})` }))
  cars.value = cr.data.map((x) => ({ value: x.id, title: x.display_name }))
  dialog.value = true
}

async function save() {
  saving.value = true
  try {
    await reservationApi.create(form.value)
    app.notify(t('app.created'))
    dialog.value = false
    await load()
  } finally {
    saving.value = false
  }
}

async function setStatus(item, status) {
  await reservationApi.setStatus(item.id, status)
  app.notify(t('app.saved'))
  await load()
}

onMounted(load)
</script>

<template>
  <div>
    <div class="d-flex align-center mb-4">
      <h1 class="text-h4 font-weight-bold">{{ t('menu.reservations') }}</h1>
      <v-spacer />
      <v-btn color="primary" prepend-icon="mdi-plus" @click="openCreate">{{ t('reservation.add') }}</v-btn>
    </div>

    <v-card class="pa-4">
      <v-data-table :headers="headers" :items="items" :loading="loading" items-per-page="20">
        <template #[`item.status_label`]="{ item }">
          <v-chip size="small" :color="statusColor(item.status)">{{ item.status_label }}</v-chip>
        </template>
        <template #[`item.actions`]="{ item }">
          <v-btn v-if="item.status === 'reserved'" size="small" variant="text" color="success" @click="setStatus(item, 'confirmed')">{{ t('reservation.confirm') }}</v-btn>
          <v-btn v-if="item.status !== 'cancelled'" size="small" variant="text" color="error" @click="setStatus(item, 'cancelled')">{{ t('reservation.cancel') }}</v-btn>
        </template>
      </v-data-table>
    </v-card>

    <v-dialog v-model="dialog" max-width="600">
      <v-card>
        <v-card-title class="bg-primary text-white">{{ t('reservation.add') }}</v-card-title>
        <v-card-text class="pt-4">
          <v-select v-model="form.client_id" :items="clients" :label="t('reservation.client')" />
          <v-select v-model="form.car_id" :items="cars" :label="t('reservation.car')" />
          <v-text-field v-model="form.start_date" :label="t('reservation.startDate')" type="date" />
          <v-text-field v-model="form.end_date" :label="t('reservation.endDate')" type="date" />
          <v-text-field v-model.number="form.deposit" :label="t('reservation.deposit')" type="number" />
          <v-textarea v-model="form.notes" :label="t('app.notes')" rows="2" />
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
