<template>
  <div>
    <PageHeader title="Отчёты" subtitle="Аналитика и выгрузка в Excel" />

    <div class="card mb-4 flex flex-wrap items-end gap-3 p-4">
      <label class="flex-1 min-w-[220px]">
        <span class="label">Тип отчёта</span>
        <select v-model="type" class="input" @change="load">
          <option value="daily_revenue">Выручка за день</option>
          <option value="yearly_revenue">Годовая выручка</option>
          <option value="most_profitable_cars">Прибыльные автомобили</option>
          <option value="active_rentals">Активные аренды</option>
          <option value="debtors">Должники</option>
          <option value="client_statistics">Статистика клиентов</option>
        </select>
      </label>

      <label v-if="type === 'daily_revenue'">
        <span class="label">Дата</span>
        <input v-model="date" type="date" class="input" @change="load" />
      </label>
      <label v-if="type === 'yearly_revenue'">
        <span class="label">Год</span>
        <input v-model.number="year" type="number" class="input w-32" @change="load" />
      </label>

      <button class="btn-primary" @click="exportExcel"><AppIcon name="download" size="18" /> Excel</button>
    </div>

    <div v-if="type === 'daily_revenue'" class="mb-4 grid grid-cols-1 gap-4 sm:grid-cols-2">
      <StatCard title="Выручка за день" :value="fmt.money(daily?.total)" icon="money" accent="green" />
      <StatCard title="Кол-во платежей" :value="daily?.payments?.length || 0" icon="payments" accent="primary" />
    </div>
    <div v-if="type === 'yearly_revenue'" class="mb-4">
      <StatCard title="Выручка за год" :value="fmt.money(yearly?.total)" icon="money" accent="green" />
    </div>

    <DataTable :columns="columns" :rows="rows" :loading="loading" empty="Нет данных">
      <template #amount="{ value }">{{ fmt.money(value) }}</template>
      <template #revenue="{ value }">{{ fmt.money(value) }}</template>
      <template #total_price="{ value }">{{ fmt.money(value) }}</template>
    </DataTable>
  </div>
</template>

<script setup lang="ts">
const api = useApi()
const fmt = useFormat()

const type = ref('daily_revenue')
const date = ref(new Date().toISOString().slice(0, 10))
const year = ref(new Date().getFullYear())
const loading = ref(false)
const rows = ref<any[]>([])
const daily = ref<any>(null)
const yearly = ref<any>(null)

const COLUMNS: Record<string, { key: string; label: string }[]> = {
  daily_revenue: [
    { key: 'receipt_number', label: 'Квитанция' },
    { key: 'client_name', label: 'Клиент' },
    { key: 'payment_type_label', label: 'Тип' },
    { key: 'amount', label: 'Сумма' },
  ],
  yearly_revenue: [
    { key: 'month', label: 'Месяц' },
    { key: 'total', label: 'Выручка' },
  ],
  most_profitable_cars: [
    { key: 'car', label: 'Автомобиль' },
    { key: 'revenue', label: 'Выручка' },
  ],
  active_rentals: [
    { key: 'contract_number', label: 'Договор' },
    { key: 'client_name', label: 'Клиент' },
    { key: 'car_name', label: 'Автомобиль' },
    { key: 'total_price', label: 'Сумма' },
  ],
  debtors: [
    { key: 'client', label: 'Клиент' },
    { key: 'contract_number', label: 'Договор' },
    { key: 'amount', label: 'Долг' },
  ],
  client_statistics: [
    { key: 'client', label: 'Клиент' },
    { key: 'code', label: 'Код' },
    { key: 'rentals', label: 'Аренд' },
  ],
}

const columns = computed(() => {
  if (type.value === 'yearly_revenue') return COLUMNS.yearly_revenue.map((c) => (c.key === 'total' ? { key: 'revenue', label: 'Выручка' } : c))
  return COLUMNS[type.value] || []
})

async function load() {
  loading.value = true
  try {
    switch (type.value) {
      case 'daily_revenue':
        daily.value = await api.reports.daily(date.value)
        rows.value = daily.value.payments || []
        break
      case 'yearly_revenue':
        yearly.value = await api.reports.yearly(year.value)
        rows.value = (yearly.value.months || []).map((m: any) => ({ month: `${m.year}-${String(m.month).padStart(2, '0')}`, revenue: m.total }))
        break
      case 'most_profitable_cars':
        rows.value = await api.reports.profitableCars()
        break
      case 'active_rentals':
        rows.value = await api.reports.activeRentals()
        break
      case 'debtors':
        rows.value = await api.reports.debtors()
        break
      case 'client_statistics':
        rows.value = await api.reports.clientStats()
        break
    }
  } finally {
    loading.value = false
  }
}

async function exportExcel() {
  const params: Record<string, any> = {}
  if (type.value === 'daily_revenue') params.date = date.value
  if (type.value === 'yearly_revenue') params.year = year.value
  await api.reports.exportReport(type.value, params)
}

onMounted(load)
</script>
