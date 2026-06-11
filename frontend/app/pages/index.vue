<template>
  <div>
    <PageHeader title="Дашборд" subtitle="Обзор ключевых показателей проката" />

    <!-- Stat cards -->
    <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">
      <StatCard title="Выручка за месяц" :value="fmt.money(summary?.monthly_revenue)" icon="money" accent="green" />
      <StatCard title="Выручка за сегодня" :value="fmt.money(summary?.today_revenue)" icon="payments" accent="primary" />
      <StatCard title="Активные аренды" :value="summary?.active_rentals ?? 0" icon="rentals" accent="violet" />
      <StatCard title="Скоро возврат" :value="summary?.upcoming_returns ?? 0" icon="clock" accent="amber" />
    </div>

    <div class="mt-5 grid grid-cols-1 gap-4 xl:grid-cols-3">
      <!-- Revenue line chart -->
      <div class="card p-5 xl:col-span-2">
        <div class="mb-3 flex items-center justify-between">
          <h3 class="font-semibold text-ink">Выручка по месяцам</h3>
          <span class="text-sm text-ink-muted">за {{ revenue.length }} мес.</span>
        </div>
        <LineChart :data="revenue" />
      </div>

      <!-- Car availability donut -->
      <div class="card p-5">
        <h3 class="mb-4 font-semibold text-ink">Автопарк</h3>
        <DonutChart :segments="carSegments" />
      </div>
    </div>

    <div class="mt-5 grid grid-cols-1 gap-4 lg:grid-cols-2">
      <!-- Top cars -->
      <div class="card p-5">
        <h3 class="mb-4 font-semibold text-ink">Самые популярные автомобили</h3>
        <ul class="space-y-3">
          <li v-for="(c, i) in topCars" :key="i" class="flex items-center gap-3">
            <span class="flex h-8 w-8 items-center justify-center rounded-lg bg-primary-100 text-sm font-bold text-primary-600">
              {{ i + 1 }}
            </span>
            <span class="flex-1 truncate text-sm text-ink">{{ c.car }}</span>
            <span class="text-sm font-semibold text-ink-soft">{{ c.rentals }} аренд</span>
          </li>
          <li v-if="!topCars.length" class="py-6 text-center text-sm text-ink-muted">Нет данных</li>
        </ul>
      </div>

      <!-- Rental statistics -->
      <div class="card p-5">
        <h3 class="mb-4 font-semibold text-ink">Статусы аренд</h3>
        <BarChart :data="rentalStatsData" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { DashboardSummary } from '~/types'

const api = useApi()
const fmt = useFormat()

const summary = ref<DashboardSummary | null>(null)
const revenue = ref<{ month: string; revenue: number }[]>([])
const topCars = ref<{ car: string; rentals: number }[]>([])
const rentalStats = ref<Record<string, number>>({})

const carSegments = computed(() => [
  { label: 'Доступны', value: summary.value?.cars_available ?? 0, color: '#10b981' },
  { label: 'В аренде', value: summary.value?.cars_rented ?? 0, color: '#8b5cf6' },
  { label: 'Забронированы', value: summary.value?.cars_reserved ?? 0, color: '#f59e0b' },
  { label: 'Обслуживание', value: summary.value?.cars_maintenance ?? 0, color: '#94a3b8' },
])

const rentalStatsData = computed(() => [
  { label: 'Активные', value: rentalStats.value.active ?? 0 },
  { label: 'Завершён.', value: rentalStats.value.completed ?? 0 },
  { label: 'Отменён.', value: rentalStats.value.cancelled ?? 0 },
])

onMounted(async () => {
  try {
    const [s, r, t, rs] = await Promise.all([
      api.dashboard.summary(),
      api.dashboard.revenueByMonth(12),
      api.dashboard.topCars(5),
      api.dashboard.rentalStatistics(),
    ])
    summary.value = s
    revenue.value = r
    topCars.value = t
    rentalStats.value = rs
  } catch {
    /* handled centrally */
  }
})
</script>
