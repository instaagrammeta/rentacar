<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { Bar, Doughnut } from 'vue-chartjs'
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  BarElement,
  CategoryScale,
  LinearScale,
  ArcElement,
} from 'chart.js'
import { dashboardApi } from '@/api'
import StatCard from '@/components/StatCard.vue'

ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale, ArcElement)

const { t } = useI18n()

const summary = ref({})
const revenueChart = ref(null)
const topCarsChart = ref(null)
const statsChart = ref(null)
const loading = ref(true)

function currency(value) {
  return new Intl.NumberFormat('ru-RU').format(value || 0) + ' ₽'
}

onMounted(async () => {
  try {
    const [s, rev, top, stats] = await Promise.all([
      dashboardApi.summary(),
      dashboardApi.revenueByMonth(12),
      dashboardApi.topCars(5),
      dashboardApi.rentalStatistics(),
    ])
    summary.value = s.data

    revenueChart.value = {
      labels: rev.data.map((r) => r.month),
      datasets: [
        { label: t('dashboard.revenueByMonth'), backgroundColor: '#2E7D32', data: rev.data.map((r) => r.revenue) },
      ],
    }
    topCarsChart.value = {
      labels: top.data.map((c) => c.car),
      datasets: [
        { label: t('dashboard.topCars'), backgroundColor: '#43A047', data: top.data.map((c) => c.rentals) },
      ],
    }
    statsChart.value = {
      labels: [t('rental.statuses.active'), t('rental.statuses.completed'), t('rental.statuses.cancelled')],
      datasets: [
        {
          backgroundColor: ['#2E7D32', '#1565C0', '#C62828'],
          data: [stats.data.active || 0, stats.data.completed || 0, stats.data.cancelled || 0],
        },
      ],
    }
  } finally {
    loading.value = false
  }
})

const chartOptions = { responsive: true, maintainAspectRatio: false }
</script>

<template>
  <div>
    <h1 class="text-h4 font-weight-bold mb-6">{{ t('dashboard.title') }}</h1>

    <v-row>
      <v-col cols="12" sm="6" md="3"><StatCard :title="t('dashboard.carsAvailable')" :value="summary.cars_available" icon="mdi-car" color="success" /></v-col>
      <v-col cols="12" sm="6" md="3"><StatCard :title="t('dashboard.carsRented')" :value="summary.cars_rented" icon="mdi-car-key" color="info" /></v-col>
      <v-col cols="12" sm="6" md="3"><StatCard :title="t('dashboard.carsReserved')" :value="summary.cars_reserved" icon="mdi-calendar-clock" color="warning" /></v-col>
      <v-col cols="12" sm="6" md="3"><StatCard :title="t('dashboard.carsMaintenance')" :value="summary.cars_maintenance" icon="mdi-wrench" color="error" /></v-col>
    </v-row>

    <v-row>
      <v-col cols="12" sm="6" md="3"><StatCard :title="t('dashboard.todayRevenue')" :value="currency(summary.today_revenue)" icon="mdi-cash" color="primary" /></v-col>
      <v-col cols="12" sm="6" md="3"><StatCard :title="t('dashboard.monthlyRevenue')" :value="currency(summary.monthly_revenue)" icon="mdi-cash-multiple" color="primary" /></v-col>
      <v-col cols="12" sm="6" md="3"><StatCard :title="t('dashboard.activeRentals')" :value="summary.active_rentals" icon="mdi-file-document" color="secondary" /></v-col>
      <v-col cols="12" sm="6" md="3"><StatCard :title="t('dashboard.upcomingReturns')" :value="summary.upcoming_returns" icon="mdi-keyboard-return" color="warning" /></v-col>
    </v-row>

    <v-row class="mt-2">
      <v-col cols="12" md="8">
        <v-card class="pa-4" height="360">
          <div class="text-subtitle-1 font-weight-bold mb-2">{{ t('dashboard.revenueByMonth') }}</div>
          <div style="height: 280px">
            <Bar v-if="revenueChart" :data="revenueChart" :options="chartOptions" />
          </div>
        </v-card>
      </v-col>
      <v-col cols="12" md="4">
        <v-card class="pa-4" height="360">
          <div class="text-subtitle-1 font-weight-bold mb-2">{{ t('dashboard.rentalStats') }}</div>
          <div style="height: 280px">
            <Doughnut v-if="statsChart" :data="statsChart" :options="chartOptions" />
          </div>
        </v-card>
      </v-col>
    </v-row>

    <v-row class="mt-2">
      <v-col cols="12">
        <v-card class="pa-4" height="340">
          <div class="text-subtitle-1 font-weight-bold mb-2">{{ t('dashboard.topCars') }}</div>
          <div style="height: 260px">
            <Bar v-if="topCarsChart" :data="topCarsChart" :options="chartOptions" />
          </div>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>
