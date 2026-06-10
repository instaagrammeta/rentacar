<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { rentalApi, clientApi, carApi } from '@/api'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'

const { t } = useI18n()
const app = useAppStore()
const auth = useAuthStore()

const items = ref([])
const clients = ref([])
const cars = ref([])
const loading = ref(false)
const dialog = ref(false)
const form = ref({})
const saving = ref(false)

const headers = [
  { title: t('rental.contractNumber'), key: 'contract_number' },
  { title: t('rental.client'), key: 'client_name' },
  { title: t('rental.car'), key: 'car_name' },
  { title: t('rental.start'), key: 'rental_start' },
  { title: t('rental.end'), key: 'rental_end' },
  { title: t('rental.totalPrice'), key: 'total_price' },
  { title: t('app.status'), key: 'status_label' },
  { title: t('app.actions'), key: 'actions', sortable: false },
]

function statusColor(status) {
  return { active: 'info', completed: 'success', cancelled: 'error' }[status]
}

async function load() {
  loading.value = true
  try {
    const { data } = await rentalApi.list({ per_page: 100 })
    items.value = data.items
  } finally {
    loading.value = false
  }
}

async function openCreate() {
  form.value = { rental_start: '', rental_end: '', deposit: 0 }
  const [c, cr] = await Promise.all([clientApi.list({ per_page: 200 }), carApi.available()])
  clients.value = c.data.items.map((x) => ({ value: x.id, title: `${x.full_name} (${x.client_code})` }))
  cars.value = cr.data.map((x) => ({ value: x.id, title: x.display_name, daily: x.daily_price, deposit: x.deposit_amount }))
  dialog.value = true
}

function onCarSelected(carId) {
  const car = cars.value.find((c) => c.value === carId)
  if (car) form.value.deposit = car.deposit
}

async function save() {
  saving.value = true
  try {
    await rentalApi.create(form.value)
    app.notify(t('app.created'))
    dialog.value = false
    await load()
  } finally {
    saving.value = false
  }
}

async function cancel(item) {
  if (!confirm(t('app.confirm') + '?')) return
  await rentalApi.cancel(item.id)
  app.notify(t('app.saved'))
  await load()
}

async function downloadContract(item) {
  const res = await fetch(rentalApi.contractUrl(item.id), {
    headers: { Authorization: `Bearer ${auth.token}` },
  })
  const blob = await res.blob()
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${item.contract_number}.pdf`
  a.click()
  URL.revokeObjectURL(url)
}

onMounted(load)
</script>

<template>
  <div>
    <div class="d-flex align-center mb-4">
      <h1 class="text-h4 font-weight-bold">{{ t('menu.rentals') }}</h1>
      <v-spacer />
      <v-btn color="primary" prepend-icon="mdi-plus" @click="openCreate">{{ t('rental.add') }}</v-btn>
    </div>

    <v-card class="pa-4">
      <v-data-table :headers="headers" :items="items" :loading="loading" items-per-page="20">
        <template #[`item.total_price`]="{ item }">{{ item.total_price }} ₽</template>
        <template #[`item.status_label`]="{ item }">
          <v-chip size="small" :color="statusColor(item.status)">{{ item.status_label }}</v-chip>
        </template>
        <template #[`item.actions`]="{ item }">
          <v-btn icon="mdi-file-pdf-box" size="small" variant="text" color="primary" :title="t('rental.downloadContract')" @click="downloadContract(item)" />
          <v-btn v-if="item.status === 'active'" icon="mdi-cancel" size="small" variant="text" color="error" @click="cancel(item)" />
        </template>
      </v-data-table>
    </v-card>

    <v-dialog v-model="dialog" max-width="600">
      <v-card>
        <v-card-title class="bg-primary text-white">{{ t('rental.add') }}</v-card-title>
        <v-card-text class="pt-4">
          <v-select v-model="form.client_id" :items="clients" :label="t('rental.client')" />
          <v-select v-model="form.car_id" :items="cars" :label="t('rental.car')" @update:model-value="onCarSelected" />
          <v-text-field v-model="form.rental_start" :label="t('rental.start')" type="date" />
          <v-text-field v-model="form.rental_end" :label="t('rental.end')" type="date" />
          <v-text-field v-model.number="form.deposit" :label="t('rental.deposit')" type="number" />
          <v-text-field v-model.number="form.start_mileage" :label="t('rental.startMileage')" type="number" />
          <v-textarea v-model="form.notes" :label="t('app.notes')" rows="2" />
          <v-alert type="info" variant="tonal" density="compact">
            Итоговая стоимость рассчитывается автоматически по тарифам автомобиля.
          </v-alert>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="dialog = false">{{ t('app.cancel') }}</v-btn>
          <v-btn color="primary" :loading="saving" @click="save">{{ t('app.create') }}</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>
