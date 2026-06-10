<script setup>
import { ref, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { rentalApi } from '@/api'
import { useAppStore } from '@/stores/app'

const { t } = useI18n()
const app = useAppStore()

const activeRentals = ref([])
const loading = ref(false)
const dialog = ref(false)
const selected = ref(null)
const form = ref({})
const preview = ref(null)
const saving = ref(false)

const headers = [
  { title: t('rental.contractNumber'), key: 'contract_number' },
  { title: t('rental.client'), key: 'client_name' },
  { title: t('rental.car'), key: 'car_name' },
  { title: t('rental.end'), key: 'rental_end' },
  { title: t('app.actions'), key: 'actions', sortable: false },
]

async function load() {
  loading.value = true
  try {
    const { data } = await rentalApi.list({ status: 'active', per_page: 100 })
    activeRentals.value = data.items
  } finally {
    loading.value = false
  }
}

function openReturn(item) {
  selected.value = item
  form.value = {
    return_date: new Date().toISOString().slice(0, 10),
    mileage: 0,
    fuel_level: 100,
    damages: '',
    damage_cost: 0,
    penalties: 0,
  }
  preview.value = null
  dialog.value = true
  updatePreview()
}

async function updatePreview() {
  if (!selected.value) return
  const { data } = await rentalApi.previewReturn(selected.value.id, {
    return_date: form.value.return_date,
    damage_cost: form.value.damage_cost,
    penalties: form.value.penalties,
  })
  preview.value = data
}

watch(() => [form.value.return_date, form.value.damage_cost, form.value.penalties], updatePreview)

async function submit() {
  saving.value = true
  try {
    await rentalApi.createReturn(selected.value.id, form.value)
    app.notify(t('app.saved'))
    dialog.value = false
    await load()
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>

<template>
  <div>
    <h1 class="text-h4 font-weight-bold mb-4">{{ t('menu.returns') }}</h1>

    <v-card class="pa-4">
      <v-data-table :headers="headers" :items="activeRentals" :loading="loading" items-per-page="20">
        <template #[`item.actions`]="{ item }">
          <v-btn size="small" color="primary" prepend-icon="mdi-keyboard-return" @click="openReturn(item)">{{ t('ret.process') }}</v-btn>
        </template>
      </v-data-table>
    </v-card>

    <v-dialog v-model="dialog" max-width="700">
      <v-card v-if="selected">
        <v-card-title class="bg-primary text-white">{{ t('ret.title') }} — {{ selected.contract_number }}</v-card-title>
        <v-card-text class="pt-4">
          <v-row>
            <v-col cols="12" md="6"><v-text-field v-model="form.return_date" :label="t('ret.returnDate')" type="date" /></v-col>
            <v-col cols="12" md="6"><v-text-field v-model.number="form.mileage" :label="t('ret.mileage')" type="number" /></v-col>
            <v-col cols="12" md="6"><v-text-field v-model.number="form.fuel_level" :label="t('ret.fuelLevel')" type="number" /></v-col>
            <v-col cols="12" md="6"><v-text-field v-model.number="form.damage_cost" :label="t('ret.damageCost')" type="number" /></v-col>
            <v-col cols="12" md="6"><v-text-field v-model.number="form.penalties" :label="t('ret.penalties')" type="number" /></v-col>
            <v-col cols="12"><v-textarea v-model="form.damages" :label="t('ret.damages')" rows="2" /></v-col>
          </v-row>

          <v-card v-if="preview" variant="tonal" color="primary" class="pa-3 mt-2">
            <div class="d-flex justify-space-between"><span>{{ t('ret.extraDays') }}:</span><b>{{ preview.extra_days }}</b></div>
            <div class="d-flex justify-space-between"><span>{{ t('ret.lateFee') }}:</span><b>{{ preview.late_fee }} ₽</b></div>
            <div class="d-flex justify-space-between"><span>{{ t('ret.damageCost') }}:</span><b>{{ preview.damage_cost }} ₽</b></div>
            <div class="d-flex justify-space-between"><span>{{ t('ret.penalties') }}:</span><b>{{ preview.penalties }} ₽</b></div>
            <v-divider class="my-2" />
            <div class="d-flex justify-space-between text-h6">
              <span>{{ t('ret.finalPayment') }}:</span>
              <b :class="preview.final_payment >= 0 ? 'text-error' : 'text-success'">{{ preview.final_payment }} ₽</b>
            </div>
            <div class="text-caption">{{ preview.final_payment >= 0 ? 'Клиент доплачивает' : 'Возврат клиенту' }}</div>
          </v-card>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="dialog = false">{{ t('app.cancel') }}</v-btn>
          <v-btn color="primary" :loading="saving" @click="submit">{{ t('ret.process') }}</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>
