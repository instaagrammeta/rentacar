<template>
  <div>
    <PageHeader title="Брони" subtitle="Бронирование автомобилей">
      <template #actions>
        <button class="btn-primary" @click="openCreate"><AppIcon name="plus" size="18" /> Новая бронь</button>
      </template>
    </PageHeader>

    <div class="card mb-4 flex flex-wrap items-center gap-3 p-4">
      <select v-model="status" class="input w-auto min-w-[180px]">
        <option value="">Все статусы</option>
        <option value="reserved">Забронировано</option>
        <option value="confirmed">Подтверждено</option>
        <option value="cancelled">Отменено</option>
      </select>
    </div>

    <DataTable :columns="columns" :rows="data.items" :loading="loading" empty="Броней нет">
      <template #start_date="{ row }">{{ fmt.date(row.start_date) }} — {{ fmt.date(row.end_date) }}</template>
      <template #deposit="{ row }">{{ fmt.money(row.deposit) }}</template>
      <template #status="{ row }"><StatusBadge :status="row.status" :label="row.status_label" /></template>
      <template #actions="{ row }">
        <button v-if="row.status !== 'cancelled'" class="btn-ghost !px-2 !py-1.5 text-red-500" title="Отменить" @click="cancel(row)">
          <AppIcon name="close" size="18" />
        </button>
      </template>
    </DataTable>

    <AppPagination :total="data.total" :page="page" :per-page="data.per_page || 20" @update:page="(p) => (page = p)" />

    <AppModal :open="showCreate" title="Новая бронь" @close="showCreate = false">
      <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <FormField v-model="form.client_id" type="select" label="Клиент" :options="clientOptions" placeholder="Выберите клиента" />
        <FormField v-model="form.car_id" type="select" label="Автомобиль" :options="carOptions" placeholder="Выберите авто" />
        <FormField v-model="form.start_date" type="date" label="Начало" required />
        <FormField v-model="form.end_date" type="date" label="Окончание" required />
        <FormField v-model="form.deposit" type="number" step="0.01" label="Депозит" />
      </div>
      <template #footer>
        <button class="btn-secondary" @click="showCreate = false">Отмена</button>
        <button class="btn-primary" :disabled="saving" @click="create">Создать</button>
      </template>
    </AppModal>
  </div>
</template>

<script setup lang="ts">
import type { Paginated, Reservation, Option } from '~/types'
import { useUiStore } from '~/stores/ui'

const api = useApi()
const ui = useUiStore()
const fmt = useFormat()

const columns = [
  { key: 'id', label: '#' },
  { key: 'client_name', label: 'Клиент' },
  { key: 'car_name', label: 'Автомобиль' },
  { key: 'start_date', label: 'Период' },
  { key: 'deposit', label: 'Депозит' },
  { key: 'status', label: 'Статус' },
]

const data = ref<Paginated<Reservation>>({ items: [], total: 0, page: 1, per_page: 20 })
const loading = ref(false)
const status = ref('')
const page = ref(1)

const showCreate = ref(false)
const saving = ref(false)
const form = ref<Record<string, any>>({ client_id: '', car_id: '', start_date: '', end_date: '', deposit: 0 })
const clientOptions = ref<Option[]>([])
const carOptions = ref<Option[]>([])

async function load() {
  loading.value = true
  try {
    data.value = await api.reservations.list({ status: status.value, page: page.value, per_page: 20 })
  } finally {
    loading.value = false
  }
}

async function openCreate() {
  showCreate.value = true
  form.value = { client_id: '', car_id: '', start_date: '', end_date: '', deposit: 0 }
  const [clients, cars] = await Promise.all([api.clients.list({ per_page: 200 }), api.cars.available()])
  clientOptions.value = clients.items.map((c) => ({ value: String(c.id), label: `${c.full_name} (${c.client_code})` }))
  carOptions.value = cars.map((c) => ({ value: String(c.id), label: c.display_name }))
}

async function create() {
  saving.value = true
  try {
    await api.reservations.create({
      ...form.value,
      client_id: Number(form.value.client_id),
      car_id: Number(form.value.car_id),
    })
    ui.success('Бронь создана')
    showCreate.value = false
    load()
  } catch {
    /* handled centrally */
  } finally {
    saving.value = false
  }
}

async function cancel(row: Reservation) {
  if (!confirm('Отменить бронь?')) return
  try {
    await api.reservations.cancel(row.id)
    ui.success('Бронь отменена')
    load()
  } catch {
    /* handled centrally */
  }
}

watch([status, page], load)
onMounted(load)
</script>
