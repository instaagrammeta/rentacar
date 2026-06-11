<template>
  <div>
    <PageHeader title="Аренды" subtitle="Договоры аренды">
      <template #actions>
        <button class="btn-primary" @click="openCreate"><AppIcon name="plus" size="18" /> Новый договор</button>
      </template>
    </PageHeader>

    <div class="card mb-4 flex flex-wrap items-center gap-3 p-4">
      <select v-model="status" class="input w-auto min-w-[180px]">
        <option value="">Все статусы</option>
        <option value="active">Активна</option>
        <option value="completed">Завершена</option>
        <option value="cancelled">Отменена</option>
      </select>
    </div>

    <DataTable :columns="columns" :rows="data.items" :loading="loading" empty="Договоров нет">
      <template #rental_start="{ row }">{{ fmt.date(row.rental_start) }} — {{ fmt.date(row.rental_end) }}</template>
      <template #pickup_at="{ row }">{{ fmt.dateTime(row.pickup_at) }}</template>
      <template #total_price="{ row }">{{ fmt.money(row.total_price) }}</template>
      <template #status="{ row }"><StatusBadge :status="row.status" :label="row.status_label" /></template>
      <template #actions="{ row }">
        <NuxtLink :to="`/rentals/${row.id}`" class="btn-ghost !px-2 !py-1.5" title="Открыть"><AppIcon name="eye" size="18" /></NuxtLink>
        <button class="btn-ghost !px-2 !py-1.5" title="Договор PDF" @click="download(row)"><AppIcon name="download" size="18" /></button>
      </template>
    </DataTable>

    <AppPagination :total="data.total" :page="page" :per-page="data.per_page || 20" @update:page="(p) => (page = p)" />

    <AppModal :open="showCreate" title="Новый договор аренды" size="lg" @close="showCreate = false">
      <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <FormField v-model="form.client_id" type="select" label="Клиент" :options="clientOptions" placeholder="Выберите клиента" />
        <FormField v-model="form.car_id" type="select" label="Автомобиль" :options="carOptions" placeholder="Выберите авто" />
        <FormField v-model="form.rental_start" type="date" label="Начало аренды" required />
        <FormField v-model="form.rental_end" type="date" label="Окончание аренды" required />
        <FormField v-model="form.pickup_at" type="datetime-local" label="Время выдачи (когда принят автомобиль)" />
        <FormField v-model="form.due_at" type="datetime-local" label="Срок возврата (дата и время)" />
        <FormField v-model="form.deposit" type="number" step="0.01" label="Депозит" />
        <FormField v-model="form.total_price" type="number" step="0.01" label="Итоговая сумма (необязательно)" />
      </div>
      <template #footer>
        <button class="btn-secondary" @click="showCreate = false">Отмена</button>
        <button class="btn-primary" :disabled="saving" @click="create">Создать договор</button>
      </template>
    </AppModal>
  </div>
</template>

<script setup lang="ts">
import type { Paginated, Rental, Option } from '~/types'
import { useUiStore } from '~/stores/ui'

const api = useApi()
const ui = useUiStore()
const fmt = useFormat()

const columns = [
  { key: 'contract_number', label: 'Договор' },
  { key: 'client_name', label: 'Клиент' },
  { key: 'car_name', label: 'Автомобиль' },
  { key: 'rental_start', label: 'Период' },
  { key: 'pickup_at', label: 'Выдан' },
  { key: 'total_price', label: 'Сумма' },
  { key: 'status', label: 'Статус' },
]

const data = ref<Paginated<Rental>>({ items: [], total: 0, page: 1, per_page: 20 })
const loading = ref(false)
const status = ref('')
const page = ref(1)

const showCreate = ref(false)
const saving = ref(false)
const form = ref<Record<string, any>>({ client_id: '', car_id: '', rental_start: '', rental_end: '', pickup_at: '', due_at: '', deposit: 0, total_price: 0 })
const clientOptions = ref<Option[]>([])
const carOptions = ref<Option[]>([])

// Local datetime string "YYYY-MM-DDTHH:MM" for datetime-local inputs.
function nowLocal(): string {
  const d = new Date()
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`
}

async function load() {
  loading.value = true
  try {
    data.value = await api.rentals.list({ status: status.value, page: page.value, per_page: 20 })
  } finally {
    loading.value = false
  }
}

async function openCreate() {
  showCreate.value = true
  form.value = { client_id: '', car_id: '', rental_start: '', rental_end: '', pickup_at: nowLocal(), due_at: '', deposit: 0, total_price: 0 }
  const [clients, cars] = await Promise.all([api.clients.list({ per_page: 200 }), api.cars.available()])
  clientOptions.value = clients.items.map((c) => ({ value: String(c.id), label: `${c.full_name} (${c.client_code})` }))
  carOptions.value = cars.map((c) => ({ value: String(c.id), label: c.display_name }))
}

async function create() {
  saving.value = true
  try {
    const payload: Record<string, any> = {
      ...form.value,
      client_id: Number(form.value.client_id),
      car_id: Number(form.value.car_id),
    }
    if (!payload.total_price) delete payload.total_price
    const rental = await api.rentals.create(payload)
    ui.success('Договор создан')
    showCreate.value = false
    await navigateTo(`/rentals/${rental.id}`)
  } catch {
    /* handled centrally */
  } finally {
    saving.value = false
  }
}

async function download(row: Rental) {
  try {
    await api.rentals.downloadContract(row.id, row.contract_number)
  } catch {
    ui.error('Не удалось скачать договор')
  }
}

watch([status, page], load)
onMounted(load)
</script>
