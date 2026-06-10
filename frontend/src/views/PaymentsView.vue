<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { paymentApi, clientApi, rentalApi } from '@/api'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'

const { t } = useI18n()
const app = useAppStore()
const auth = useAuthStore()

const items = ref([])
const clients = ref([])
const rentals = ref([])
const loading = ref(false)
const dialog = ref(false)
const form = ref({})
const saving = ref(false)

const methodOptions = [
  { value: 'cash', title: t('payment.methods.cash') },
  { value: 'bank_transfer', title: t('payment.methods.bank_transfer') },
  { value: 'card', title: t('payment.methods.card') },
]
const typeOptions = [
  { value: 'deposit', title: t('payment.types.deposit') },
  { value: 'rental', title: t('payment.types.rental') },
  { value: 'penalty', title: t('payment.types.penalty') },
  { value: 'damage', title: t('payment.types.damage') },
  { value: 'refund', title: t('payment.types.refund') },
]

const headers = [
  { title: t('payment.receipt'), key: 'receipt_number' },
  { title: t('payment.client'), key: 'client_name' },
  { title: t('payment.amount'), key: 'amount' },
  { title: t('payment.method'), key: 'method_label' },
  { title: t('payment.type'), key: 'payment_type_label' },
  { title: t('payment.date'), key: 'paid_at' },
  { title: t('app.actions'), key: 'actions', sortable: false },
]

async function load() {
  loading.value = true
  try {
    const { data } = await paymentApi.list({ per_page: 100 })
    items.value = data.items
  } finally {
    loading.value = false
  }
}

async function openCreate() {
  form.value = { method: 'cash', payment_type: 'rental', amount: 0 }
  const [c, r] = await Promise.all([clientApi.list({ per_page: 200 }), rentalApi.list({ per_page: 200 })])
  clients.value = c.data.items.map((x) => ({ value: x.id, title: `${x.full_name} (${x.client_code})` }))
  rentals.value = r.data.items.map((x) => ({ value: x.id, title: x.contract_number }))
  dialog.value = true
}

async function save() {
  saving.value = true
  try {
    await paymentApi.create(form.value)
    app.notify(t('app.created'))
    dialog.value = false
    await load()
  } finally {
    saving.value = false
  }
}

async function downloadReceipt(item) {
  const res = await fetch(paymentApi.receiptUrl(item.id), {
    headers: { Authorization: `Bearer ${auth.token}` },
  })
  const blob = await res.blob()
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${item.receipt_number}.pdf`
  a.click()
  URL.revokeObjectURL(url)
}

onMounted(load)
</script>

<template>
  <div>
    <div class="d-flex align-center mb-4">
      <h1 class="text-h4 font-weight-bold">{{ t('menu.payments') }}</h1>
      <v-spacer />
      <v-btn color="primary" prepend-icon="mdi-plus" @click="openCreate">{{ t('payment.add') }}</v-btn>
    </div>

    <v-card class="pa-4">
      <v-data-table :headers="headers" :items="items" :loading="loading" items-per-page="20">
        <template #[`item.amount`]="{ item }">{{ item.amount }} ₽</template>
        <template #[`item.actions`]="{ item }">
          <v-btn icon="mdi-receipt-text" size="small" variant="text" color="primary" :title="t('payment.downloadReceipt')" @click="downloadReceipt(item)" />
        </template>
      </v-data-table>
    </v-card>

    <v-dialog v-model="dialog" max-width="600">
      <v-card>
        <v-card-title class="bg-primary text-white">{{ t('payment.add') }}</v-card-title>
        <v-card-text class="pt-4">
          <v-select v-model="form.client_id" :items="clients" :label="t('payment.client')" />
          <v-select v-model="form.rental_id" :items="rentals" :label="t('payment.rental')" clearable />
          <v-text-field v-model.number="form.amount" :label="t('payment.amount')" type="number" />
          <v-select v-model="form.method" :items="methodOptions" :label="t('payment.method')" />
          <v-select v-model="form.payment_type" :items="typeOptions" :label="t('payment.type')" />
          <v-textarea v-model="form.notes" :label="t('app.notes')" rows="2" />
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
