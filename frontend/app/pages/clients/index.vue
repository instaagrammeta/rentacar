<template>
  <div>
    <PageHeader title="Клиенты" subtitle="База клиентов проката">
      <template #actions>
        <NuxtLink to="/clients/new" class="btn-primary">
          <AppIcon name="plus" size="18" /> Новый клиент
        </NuxtLink>
      </template>
    </PageHeader>

    <div class="card mb-4 flex flex-wrap items-center gap-3 p-4">
      <div class="relative flex-1 min-w-[220px]">
        <span class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-ink-muted">
          <AppIcon name="search" size="18" />
        </span>
        <input v-model="search" class="input pl-10" placeholder="Поиск по имени, телефону, коду…" />
      </div>
      <select v-model="status" class="input w-auto min-w-[180px]">
        <option value="">Все статусы</option>
        <option value="active">Активен</option>
        <option value="pending_verification">Ожидает проверки</option>
        <option value="blacklisted">В чёрном списке</option>
      </select>
    </div>

    <DataTable :columns="columns" :rows="data.items" :loading="loading" empty="Клиенты не найдены">
      <template #full_name="{ row }">
        <div class="flex items-center gap-2">
          <span class="font-medium text-ink">{{ row.full_name }}</span>
          <span v-if="row.is_vip" class="badge bg-amber-50 text-amber-700">VIP</span>
        </div>
      </template>
      <template #status="{ row }">
        <StatusBadge :status="row.status" :label="row.status_label" />
      </template>
      <template #actions="{ row }">
        <NuxtLink :to="`/clients/${row.id}`" class="btn-ghost !px-2 !py-1.5" title="Открыть">
          <AppIcon name="eye" size="18" />
        </NuxtLink>
        <button class="btn-ghost !px-2 !py-1.5 text-red-500" title="Удалить" @click="remove(row)">
          <AppIcon name="trash" size="18" />
        </button>
      </template>
    </DataTable>

    <AppPagination :total="data.total" :page="page" :per-page="data.per_page || 20" @update:page="(p) => (page = p)" />
  </div>
</template>

<script setup lang="ts">
import type { Client, Paginated } from '~/types'
import { useUiStore } from '~/stores/ui'

const api = useApi()
const ui = useUiStore()

const columns = [
  { key: 'client_code', label: 'Код' },
  { key: 'full_name', label: 'ФИО' },
  { key: 'phone', label: 'Телефон' },
  { key: 'driver_experience_years', label: 'Стаж, лет' },
  { key: 'status', label: 'Статус' },
]

const data = ref<Paginated<Client>>({ items: [], total: 0, page: 1, per_page: 20 })
const loading = ref(false)
const search = ref('')
const status = ref('')
const page = ref(1)
let timer: any

async function load() {
  loading.value = true
  try {
    data.value = await api.clients.list({ search: search.value, status: status.value, page: page.value, per_page: 20 })
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

async function remove(row: Client) {
  if (!confirm(`Удалить клиента «${row.full_name}»?`)) return
  try {
    await api.clients.remove(row.id)
    ui.success('Клиент удалён')
    load()
  } catch {
    /* handled centrally */
  }
}

onMounted(load)
</script>
