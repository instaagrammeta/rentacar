<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { reportApi } from '@/api'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'

const { t } = useI18n()
const app = useAppStore()
const auth = useAuthStore()

const activeReport = ref(null)
const rows = ref([])
const headers = ref([])
const loading = ref(false)
const currentType = ref(null)

const year = ref(new Date().getFullYear())
const month = ref(new Date().getMonth() + 1)
const day = ref(new Date().toISOString().slice(0, 10))

const reportList = [
  { type: 'daily_revenue', title: t('reports.dailyRevenue'), icon: 'mdi-calendar-today' },
  { type: 'monthly_revenue', title: t('reports.monthlyRevenue'), icon: 'mdi-calendar-month' },
  { type: 'yearly_revenue', title: t('reports.yearlyRevenue'), icon: 'mdi-calendar' },
  { type: 'most_profitable_cars', title: t('reports.profitableCars'), icon: 'mdi-car-sports' },
  { type: 'active_rentals', title: t('reports.activeRentals'), icon: 'mdi-file-document' },
  { type: 'debtors', title: t('reports.debtors'), icon: 'mdi-account-alert' },
  { type: 'client_statistics', title: t('reports.clientStats'), icon: 'mdi-chart-pie' },
]

async function runReport(type) {
  currentType.value = type
  loading.value = true
  try {
    let data
    if (type === 'daily_revenue') {
      data = (await reportApi.daily(day.value)).data
      headers.value = [
        { title: t('payment.receipt'), key: 'receipt_number' },
        { title: t('payment.client'), key: 'client_name' },
        { title: t('payment.amount'), key: 'amount' },
      ]
      rows.value = data.payments
      activeReport.value = `${t('reports.dailyRevenue')}: ${data.total} ₽`
    } else if (type === 'yearly_revenue') {
      data = (await reportApi.yearly(year.value)).data
      headers.value = [{ title: 'Месяц', key: 'month' }, { title: 'Выручка', key: 'total' }]
      rows.value = data.months
      activeReport.value = `${t('reports.yearlyRevenue')} ${year.value}: ${data.total} ₽`
    } else if (type === 'most_profitable_cars') {
      data = (await reportApi.profitableCars()).data
      headers.value = [{ title: t('rental.car'), key: 'car' }, { title: 'Выручка', key: 'revenue' }]
      rows.value = data
      activeReport.value = t('reports.profitableCars')
    } else if (type === 'active_rentals') {
      data = (await reportApi.activeRentals()).data
      headers.value = [
        { title: t('rental.contractNumber'), key: 'contract_number' },
        { title: t('rental.client'), key: 'client_name' },
        { title: t('rental.car'), key: 'car_name' },
        { title: t('rental.totalPrice'), key: 'total_price' },
      ]
      rows.value = data
      activeReport.value = t('reports.activeRentals')
    } else if (type === 'debtors') {
      data = (await reportApi.debtors()).data
      headers.value = [
        { title: t('rental.client'), key: 'client' },
        { title: t('rental.contractNumber'), key: 'contract_number' },
        { title: 'Долг', key: 'amount' },
      ]
      rows.value = data
      activeReport.value = t('reports.debtors')
    } else if (type === 'client_statistics') {
      data = (await reportApi.clientStats()).data
      headers.value = [
        { title: t('rental.client'), key: 'client' },
        { title: t('client.code'), key: 'code' },
        { title: 'Кол-во аренд', key: 'rentals' },
      ]
      rows.value = data
      activeReport.value = t('reports.clientStats')
    } else if (type === 'monthly_revenue') {
      data = (await reportApi.monthly(year.value, month.value)).data
      headers.value = [{ title: 'Год', key: 'year' }, { title: 'Месяц', key: 'month' }, { title: 'Выручка', key: 'total' }]
      rows.value = [data]
      activeReport.value = `${t('reports.monthlyRevenue')}: ${data.total} ₽`
    }
  } finally {
    loading.value = false
  }
}

async function exportExcel() {
  if (!currentType.value) return
  const params = { date: day.value, year: year.value, month: month.value }
  const res = await fetch(reportApi.exportUrl(currentType.value, params), {
    headers: { Authorization: `Bearer ${auth.token}` },
  })
  const blob = await res.blob()
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${currentType.value}.xlsx`
  a.click()
  URL.revokeObjectURL(url)
  app.notify(t('app.export'))
}
</script>

<template>
  <div>
    <h1 class="text-h4 font-weight-bold mb-4">{{ t('reports.title') }}</h1>

    <v-row>
      <v-col cols="12" md="4">
        <v-card class="pa-4">
          <v-row class="mb-2">
            <v-col cols="12"><v-text-field v-model="day" label="Дата" type="date" hide-details /></v-col>
            <v-col cols="6"><v-text-field v-model.number="year" label="Год" type="number" hide-details /></v-col>
            <v-col cols="6"><v-text-field v-model.number="month" label="Месяц" type="number" hide-details /></v-col>
          </v-row>
          <v-list>
            <v-list-item
              v-for="r in reportList"
              :key="r.type"
              :prepend-icon="r.icon"
              :title="r.title"
              :active="currentType === r.type"
              @click="runReport(r.type)"
            />
          </v-list>
        </v-card>
      </v-col>

      <v-col cols="12" md="8">
        <v-card class="pa-4">
          <div class="d-flex align-center mb-3">
            <div class="text-subtitle-1 font-weight-bold">{{ activeReport || 'Выберите отчёт' }}</div>
            <v-spacer />
            <v-btn v-if="currentType" color="success" prepend-icon="mdi-file-excel" @click="exportExcel">{{ t('reports.exportExcel') }}</v-btn>
          </div>
          <v-data-table v-if="headers.length" :headers="headers" :items="rows" :loading="loading" items-per-page="25" />
          <div v-else class="text-center text-medium-emphasis pa-8">{{ t('app.noData') }}</div>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>
