<template>
  <div>
    <PageHeader title="Автомобили" subtitle="Автопарк проката">
      <template #actions>
        <NuxtLink to="/cars/new" class="btn-primary"><AppIcon name="plus" size="18" /> Добавить авто</NuxtLink>
      </template>
    </PageHeader>

    <!-- Availability filter tabs -->
    <div class="mb-4 flex flex-wrap gap-2">
      <button
        v-for="t in tabs"
        :key="t.value"
        class="flex items-center gap-2 rounded-xl border px-3.5 py-2 text-sm font-medium transition"
        :class="status === t.value
          ? 'border-primary-500 bg-primary-500 text-white'
          : 'border-surface-border bg-white text-ink-soft hover:bg-surface-muted'"
        @click="setStatus(t.value)"
      >
        {{ t.label }}
        <span
          class="rounded-full px-1.5 py-0.5 text-xs"
          :class="status === t.value ? 'bg-white/20' : 'bg-surface-muted text-ink-soft'"
        >{{ t.count }}</span>
      </button>
    </div>

    <div class="card mb-4 flex flex-wrap items-center gap-3 p-4">
      <div class="relative flex-1 min-w-[220px]">
        <span class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-ink-muted"><AppIcon name="search" size="18" /></span>
        <input v-model="search" class="input pl-10" placeholder="Поиск: марка, модель, номер…" />
      </div>
    </div>

    <DataTable :columns="columns" :rows="data.items" :loading="loading" empty="Автомобили не найдены">
      <template #photo="{ row }">
        <div class="h-12 w-16 overflow-hidden rounded-lg border border-surface-border bg-surface-muted">
          <img v-if="row.photos && row.photos.length" :src="api.fileUrl(row.photos[0])" class="h-full w-full object-cover" alt="" />
          <div v-else class="flex h-full w-full items-center justify-center text-ink-muted"><AppIcon name="car" size="18" /></div>
        </div>
      </template>
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
import type { Car, DashboardSummary, Paginated } from '~/types'
import { useUiStore } from '~/stores/ui'

const api = useApi()
const ui = useUiStore()
const fmt = useFormat()

const columns = [
  { key: 'photo', label: 'Фото' },
  { key: 'plate_number', label: 'Гос. номер' },
  { key: 'display_name', label: 'Автомобиль' },
  { key: 'mileage', label: 'Пробег, км' },
  { key: 'daily_price', label: 'Цена/сутки' },
  { key: 'status', label: 'Статус' },
]

const data = ref<Paginated<Car>>({ items: [], total: 0, page: 1, per_page: 20 })
const summary = ref<DashboardSummary | null>(null)
const loading = ref(false)
const search = ref('')
const status = ref('')
const page = ref(1)
let timer: any

const tabs = computed(() => [
  { value: '', label: 'Все', count: totalCars.value },
  { value: 'available', label: 'Свободные', count: summary.value?.cars_available ?? 0 },
  { value: 'reserved', label: 'Забронированы', count: summary.value?.cars_reserved ?? 0 },
  { value: 'rented', label: 'В аренде', count: summary.value?.cars_rented ?? 0 },
  { value: 'maintenance', label: 'Обслуживание', count: summary.value?.cars_maintenance ?? 0 },
])

const totalCars = computed(() => {
  const s = summary.value
  if (!s) return 0
  return s.cars_available + s.cars_reserved + s.cars_rented + s.cars_maintenance
})

function setStatus(v: string) {
  status.value = v
  page.value = 1
}

async function loadSummary() {
  try {
    summary.value = await api.dashboard.summary()
  } catch {
    /* counts are optional */
  }
}

async function load() {
  loading.value = true
  try {
    data.value = await api.cars.list({ search: search.value, status: status.value, page: page.value, per_page: 20 })
  } finally {
    loading.value = false
  }
}

watch(search, () => {
  page.value = 1
  clearTimeout(timer)
  timer = setTimeout(load, 300)
})
watch([status, page], load)

async function remove(row: Car) {
  if (!confirm(`Удалить «${row.display_name}»?`)) return
  try {
    await api.cars.remove(row.id)
    ui.success('Автомобиль удалён')
    load()
    loadSummary()
  } catch {
    /* handled centrally */
  }
}

onMounted(() => {
  load()
  loadSummary()
})
</script>
