<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { clientApi, uploadApi } from '@/api'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()

const data = ref(null)
const loading = ref(true)

function fileUrl(path) {
  return uploadApi.fileUrl(path)
}

onMounted(async () => {
  try {
    const { data: history } = await clientApi.history(route.params.id)
    data.value = history
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div>
    <v-btn variant="text" prepend-icon="mdi-arrow-left" class="mb-4" @click="router.back()">
      {{ t('menu.clients') }}
    </v-btn>

    <v-progress-linear v-if="loading" indeterminate color="primary" />

    <template v-if="data">
      <v-row>
        <v-col cols="12" md="4">
          <v-card class="pa-4 text-center">
            <v-img
              v-if="data.client.qr_code_path"
              :src="fileUrl(data.client.qr_code_path)"
              max-width="180"
              class="mx-auto mb-3"
            />
            <h2 class="text-h6 font-weight-bold">{{ data.client.full_name }}</h2>
            <v-chip color="primary" class="my-2">{{ data.client.client_code }}</v-chip>
            <v-list density="compact">
              <v-list-item :title="data.client.phone" prepend-icon="mdi-phone" />
              <v-list-item :title="data.client.email || '—'" prepend-icon="mdi-email" />
              <v-list-item :title="`${t('client.experience')}: ${data.client.driver_experience_years}`" prepend-icon="mdi-card-account-details" />
              <v-list-item :title="data.client.status_label" prepend-icon="mdi-information" />
            </v-list>
            <div class="d-flex justify-center ga-2 mt-2">
              <v-img v-if="data.client.passport_scan" :src="fileUrl(data.client.passport_scan)" max-width="80" />
              <v-img v-if="data.client.driver_license_scan" :src="fileUrl(data.client.driver_license_scan)" max-width="80" />
            </div>
          </v-card>
        </v-col>

        <v-col cols="12" md="8">
          <v-card class="pa-4 mb-4">
            <div class="text-subtitle-1 font-weight-bold mb-2">{{ t('menu.rentals') }}</div>
            <v-table density="compact">
              <thead><tr><th>{{ t('rental.contractNumber') }}</th><th>{{ t('rental.car') }}</th><th>{{ t('rental.totalPrice') }}</th><th>{{ t('app.status') }}</th></tr></thead>
              <tbody>
                <tr v-for="r in data.rentals" :key="r.id">
                  <td>{{ r.contract_number }}</td><td>{{ r.car_name }}</td><td>{{ r.total_price }}</td><td>{{ r.status_label }}</td>
                </tr>
                <tr v-if="!data.rentals.length"><td colspan="4" class="text-center text-medium-emphasis">{{ t('app.noData') }}</td></tr>
              </tbody>
            </v-table>
          </v-card>

          <v-card class="pa-4 mb-4">
            <div class="text-subtitle-1 font-weight-bold mb-2">{{ t('menu.payments') }}</div>
            <v-table density="compact">
              <thead><tr><th>{{ t('payment.receipt') }}</th><th>{{ t('payment.amount') }}</th><th>{{ t('payment.type') }}</th><th>{{ t('payment.date') }}</th></tr></thead>
              <tbody>
                <tr v-for="p in data.payments" :key="p.id">
                  <td>{{ p.receipt_number }}</td><td>{{ p.amount }}</td><td>{{ p.payment_type_label }}</td><td>{{ p.paid_at }}</td>
                </tr>
                <tr v-if="!data.payments.length"><td colspan="4" class="text-center text-medium-emphasis">{{ t('app.noData') }}</td></tr>
              </tbody>
            </v-table>
          </v-card>

          <v-row>
            <v-col cols="12" md="6">
              <v-card class="pa-4">
                <div class="text-subtitle-1 font-weight-bold mb-2">Штрафы</div>
                <v-list density="compact">
                  <v-list-item v-for="(p, i) in data.penalties" :key="i" :title="p.contract_number" :subtitle="`${p.amount} ₽`" />
                  <v-list-item v-if="!data.penalties.length" :title="t('app.noData')" />
                </v-list>
              </v-card>
            </v-col>
            <v-col cols="12" md="6">
              <v-card class="pa-4">
                <div class="text-subtitle-1 font-weight-bold mb-2">{{ t('menu.accidents') }}</div>
                <v-list density="compact">
                  <v-list-item v-for="a in data.accidents" :key="a.id" :title="a.car_name" :subtitle="`${a.accident_date} — ${a.repair_cost} ₽`" />
                  <v-list-item v-if="!data.accidents.length" :title="t('app.noData')" />
                </v-list>
              </v-card>
            </v-col>
          </v-row>
        </v-col>
      </v-row>
    </template>
  </div>
</template>
