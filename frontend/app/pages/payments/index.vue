<template>
  <div>
    <PageHeader title="Платежи" subtitle="Приём оплат и квитанции">
      <template #actions>
        <button class="btn-primary" @click="openCreate"><AppIcon name="plus" size="18" /> Новый платёж</button>
      </template>
    </PageHeader>

    <DataTable :columns="columns" :rows="data.items" :loading="loading" empty="Платежей нет">
      <template #amount="{ row }"><span class="font-semibold text-ink">{{ fmt.money(row.amount) }}</span></template>
      <template #paid_at="{ row }">{{ fmt.dateTime(row.paid_at) }}</template>
      <template #actions="{ row }">
        <button class="btn-ghost !px-2 !py-1.5" title="Квитанция" @click="downloadReceipt(row)"><AppIcon name="download" size="18" /></button>
      </template>
    </DataTable>

    <AppPagination :total="data.total" :page="page" :per-page="data.per_page || 20" @update:page="(p) => (page = p)" />

    <AppModal :open="showCreate" title="Новый платёж" @close="showCreate = false">
      <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <FormField v-model="form.client_id" type="select" label="Клиент" :options="clientOptions" placeholder="Выберите клиента" />
        <FormField v-model="form.amount" type="number" step="0.01" label="Сумма" required />
        <FormField v-model="form.payment_type" type="select" label="Тип платежа" :options="typeOptions" />
        <FormField v-model="form.method" type="select" label="Способ оплаты" :options="methodOptions" />
        <div class="sm:col-span-2"><FormField v-model="form.notes" type="textarea" label="Примечание" /></div>
      </div>
      <template #footer>
        <button class="btn-secondary" @click="showCreate = false">Отмена</button>
        <button class="btn-primary" :disabled="saving" @click="create">Принять платёж</button>
      </template>
    </AppModal>
  </div>
</template>

<script setup lang="ts">
import type { Paginated, Payment, Option } from '~/types'
import { useUiStore } from '~/stores/ui'

const api = useApi()
const ui = useUiStore()
const fmt = useFormat()

const columns = [
  { key: 'receipt_number', label: 'Квитанция' },
  { key: 'client_name', label: 'Клиент' },
  { key: 'payment_type_label', label: 'Тип' },
  { key: 'method_label', label: 'Способ' },
  { key: 'amount', label: 'Сумма' },
  { key: 'paid_at', label: 'Дата' },
]

const typeOptions: Option[] = [
  { value: 'rental', label: 'Оплата аренды' },
  { value: 'deposit', label: 'Депозит' },
  { value: 'penalty', label: 'Штраф' },
  { value: 'damage', label: 'Возмещение ущерба' },
  { value: 'refund', label: 'Возврат средств' },
]
const methodOptions: Option[] = [
  { value: 'cash', label: 'Наличные' },
  { value: 'card', label: 'Карта' },
  { value: 'bank_transfer', label: 'Банковский перевод' },
]

const data = ref<Paginated<Payment>>({ items: [], total: 0, page: 1, per_page: 20 })
const loading = ref(false)
const page = ref(1)

const showCreate = ref(false)
const saving = ref(false)
const form = ref<Record<string, any>>({ client_id: '', amount: 0, payment_type: 'rental', method: 'cash', notes: '' })
const clientOptions = ref<Option[]>([])

async function load() {
  loading.value = true
  try {
    data.value = await api.payments.list({ page: page.value, per_page: 20 })
  } finally {
    loading.value = false
  }
}

async function openCreate() {
  showCreate.value = true
  form.value = { client_id: '', amount: 0, payment_type: 'rental', method: 'cash', notes: '' }
  const clients = await api.clients.list({ per_page: 200 })
  clientOptions.value = clients.items.map((c) => ({ value: String(c.id), label: `${c.full_name} (${c.client_code})` }))
}

async function create() {
  saving.value = true
  try {
    await api.payments.create({ ...form.value, client_id: Number(form.value.client_id), amount: Number(form.value.amount) })
    ui.success('Платёж принят')
    showCreate.value = false
    load()
  } catch {
    /* handled centrally */
  } finally {
    saving.value = false
  }
}

async function downloadReceipt(row: Payment) {
  try {
    await api.payments.downloadReceipt(row.id, row.receipt_number)
  } catch {
    ui.error('Не удалось скачать квитанцию')
  }
}

watch(page, load)
onMounted(load)
</script>
