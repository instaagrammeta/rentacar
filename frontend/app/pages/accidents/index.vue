<template>
  <div>
    <PageHeader title="ДТП и повреждения" subtitle="Учёт происшествий с автомобилями">
      <template #actions>
        <button class="btn-primary" @click="openCreate"><AppIcon name="plus" size="18" /> Зарегистрировать</button>
      </template>
    </PageHeader>

    <DataTable :columns="columns" :rows="data.items" :loading="loading" empty="Записей нет">
      <template #accident_date="{ row }">{{ fmt.date(row.accident_date) }}</template>
      <template #repair_cost="{ row }">{{ fmt.money(row.repair_cost) }}</template>
      <template #actions="{ row }">
        <button class="btn-ghost !px-2 !py-1.5 text-red-500" @click="remove(row)"><AppIcon name="trash" size="18" /></button>
      </template>
    </DataTable>

    <AppPagination :total="data.total" :page="page" :per-page="data.per_page || 20" @update:page="(p) => (page = p)" />

    <AppModal :open="showCreate" title="Регистрация ДТП" @close="showCreate = false">
      <div class="space-y-4">
        <FormField v-model="form.car_id" type="select" label="Автомобиль" :options="carOptions" placeholder="Выберите авто" />
        <FormField v-model="form.client_id" type="select" label="Клиент (если есть)" :options="clientOptions" placeholder="—" />
        <FormField v-model="form.accident_date" type="date" label="Дата происшествия" />
        <FormField v-model="form.repair_cost" type="number" step="0.01" label="Стоимость ремонта" />
        <FormField v-model="form.description" type="textarea" label="Описание" />
        <PhotoUploader v-model="photos" label="Фото повреждений" subfolder="accidents" />
      </div>
      <template #footer>
        <button class="btn-secondary" @click="showCreate = false">Отмена</button>
        <button class="btn-primary" :disabled="saving" @click="create">Сохранить</button>
      </template>
    </AppModal>
  </div>
</template>

<script setup lang="ts">
import type { Accident, Paginated, Option } from '~/types'
import { useUiStore } from '~/stores/ui'

const api = useApi()
const ui = useUiStore()
const fmt = useFormat()

const columns = [
  { key: 'id', label: '#' },
  { key: 'car_name', label: 'Автомобиль' },
  { key: 'client_name', label: 'Клиент' },
  { key: 'accident_date', label: 'Дата' },
  { key: 'repair_cost', label: 'Ремонт' },
]

const data = ref<Paginated<Accident>>({ items: [], total: 0, page: 1, per_page: 20 })
const loading = ref(false)
const page = ref(1)

const showCreate = ref(false)
const saving = ref(false)
const form = ref<Record<string, any>>({ car_id: '', client_id: '', accident_date: '', repair_cost: 0, description: '' })
const photos = ref<string[]>([])
const carOptions = ref<Option[]>([])
const clientOptions = ref<Option[]>([])

async function load() {
  loading.value = true
  try {
    data.value = await api.accidents.list({ page: page.value, per_page: 20 })
  } finally {
    loading.value = false
  }
}

async function openCreate() {
  showCreate.value = true
  form.value = { car_id: '', client_id: '', accident_date: new Date().toISOString().slice(0, 10), repair_cost: 0, description: '' }
  photos.value = []
  const [cars, clients] = await Promise.all([api.cars.list({ per_page: 200 }), api.clients.list({ per_page: 200 })])
  carOptions.value = cars.items.map((c) => ({ value: String(c.id), label: c.display_name }))
  clientOptions.value = clients.items.map((c) => ({ value: String(c.id), label: c.full_name }))
}

async function create() {
  saving.value = true
  try {
    const payload: Record<string, any> = { ...form.value, car_id: Number(form.value.car_id), photos: photos.value }
    if (form.value.client_id) payload.client_id = Number(form.value.client_id)
    else delete payload.client_id
    await api.accidents.create(payload)
    ui.success('Запись о ДТП создана')
    showCreate.value = false
    load()
  } catch {
    /* handled centrally */
  } finally {
    saving.value = false
  }
}

async function remove(row: Accident) {
  if (!confirm('Удалить запись о ДТП?')) return
  try {
    await api.accidents.remove(row.id)
    ui.success('Запись удалена')
    load()
  } catch {
    /* handled centrally */
  }
}

watch(page, load)
onMounted(load)
</script>
