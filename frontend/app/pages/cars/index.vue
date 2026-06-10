<template>
  <div>
    <PageHeader title="Автомобили" subtitle="Автопарк проката">
      <template #actions>
        <NuxtLink to="/cars/new" class="btn-primary"><AppIcon name="plus" size="18" /> Добавить авто</NuxtLink>
      </template>
    </PageHeader>

    <div class="card mb-4 flex flex-wrap items-center gap-3 p-4">
      <div class="relative flex-1 min-w-[220px]">
        <span class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-ink-muted"><AppIcon name="search" size="18" /></span>
        <input v-model="search" class="input pl-10" placeholder="Поиск: марка, модель, номер…" />
      </div>
      <select v-model="status" class="input w-auto min-w-[180px]">
        <option value="">Все статусы</option>
        <option value="available">Доступен</option>
        <option value="reserved">Забронирован</option>
        <option value="rented">В аренде</option>
        <option value="maintenance">На обслуживании</option>
      </select>
    </div>

    <DataTable :columns="columns" :rows="data.items" :loading="loading" empty="Автомобили не найдены">
      <template #display_name="{ row }">
        <span class="font-medium text-ink">{{ row.brand }} {{ row.model }}</span>
        <span class="ml-1 text-ink-muted">{{ row.year }}</span>
      </template>
      <template #daily_price="{ row }">{{ fmt.money(row.daily_price) }}</template>
      <template #status="{ row }"><StatusBadge :status="row.status" :label="row.status_label" /></template>
      <template #actions="{ row }">
        <NuxtLink :to="`/cars/${row.id}`" class="btn-ghost !px-2 !py-1.5"><AppIcon name="edit" size="18" /></NuxtLink>
        <button class="btn-ghost !px-2 !py-1.5 text-red-500" @click="remove(row)"><AppIcon name="trash" size="18" /></button>
      </template>
    </DataTable>

    <AppPagination :total="data.total" :page="page" :per-page="data.per_page || 20" @update:page="(p) => (page = p)" />
  </div>
</template>

<script setup lang="ts">
import type { Car, Paginated } from '~/types'
import { useUiStore } from '~/stores/ui'

const api = useApi()
const ui = useUiStore()
const fmt = useFormat()

const columns = [
  { key: 'plate_number', label: 'Гос. номер' },
  { key: 'display_name', label: 'Автомобиль' },
  { key: 'mileage', label: 'Пробег, км' },
  { key: 'daily_price', label: 'Цена/сутки' },
  { key: 'status', label: 'Статус' },
]

const data = ref<Paginated<Car>>({ items: [], total: 0, page: 1, per_page: 20 })
const loading = ref(false)
const search = ref('')
const status = ref('')
const page = ref(1)
let timer: any

async function load() {
  loading.value = true
  try {
    data.value = await api.cars.list({ search: search.value, status: status.value, page: page.value, per_page: 20 })
  } finally {
    loading.value = false
  }
}

watch([search, status], () => {
  page.value = 1
  clearTimeout(timer)
  timer = setTimeout(load, 300)
})
watch(page, load)

async function remove(row: Car) {
  if (!confirm(`Удалить «${row.display_name}»?`)) return
  try {
    await api.cars.remove(row.id)
    ui.success('Автомобиль удалён')
    load()
  } catch {
    /* handled centrally */
  }
}

onMounted(load)
</script>
