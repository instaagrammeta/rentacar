<template>
  <div>
    <PageHeader :title="form.last_name ? `${form.last_name} ${form.first_name}` : 'Клиент'" subtitle="Карточка клиента">
      <template #actions>
        <NuxtLink to="/clients" class="btn-secondary"><AppIcon name="back" size="18" /> Назад</NuxtLink>
      </template>
    </PageHeader>

    <div class="grid grid-cols-1 gap-5 lg:grid-cols-3">
      <!-- Edit form -->
      <form class="card p-6 lg:col-span-2" @submit.prevent="save">
        <h3 class="mb-4 font-semibold text-ink">Данные клиента</h3>
        <ClientForm v-model="form" />
        <div class="mt-6 flex justify-end">
          <button type="submit" class="btn-primary" :disabled="saving">{{ saving ? 'Сохранение…' : 'Сохранить' }}</button>
        </div>
      </form>

      <!-- QR + meta -->
      <div class="space-y-5">
        <div class="card p-5 text-center">
          <p class="text-sm text-ink-muted">Код клиента</p>
          <p class="text-lg font-bold text-ink">{{ client?.client_code }}</p>
          <img v-if="qrUrl" :src="qrUrl" alt="QR" class="mx-auto mt-3 h-40 w-40 rounded-xl border border-surface-border" />
        </div>
      </div>
    </div>

    <!-- History -->
    <div class="mt-6 grid grid-cols-1 gap-5 lg:grid-cols-2">
      <div class="card p-5">
        <h3 class="mb-3 font-semibold text-ink">История аренд</h3>
        <ul class="divide-y divide-surface-border">
          <li v-for="r in history.rentals" :key="r.id" class="flex items-center justify-between py-2.5 text-sm">
            <span>{{ r.contract_number }} — {{ r.car_name }}</span>
            <StatusBadge :status="r.status" :label="r.status_label" />
          </li>
          <li v-if="!history.rentals?.length" class="py-4 text-center text-sm text-ink-muted">Нет аренд</li>
        </ul>
      </div>

      <div class="card p-5">
        <h3 class="mb-3 font-semibold text-ink">Платежи</h3>
        <ul class="divide-y divide-surface-border">
          <li v-for="p in history.payments" :key="p.id" class="flex items-center justify-between py-2.5 text-sm">
            <span>{{ p.receipt_number }} · {{ p.payment_type_label }}</span>
            <span class="font-semibold">{{ fmt.money(p.amount) }}</span>
          </li>
          <li v-if="!history.payments?.length" class="py-4 text-center text-sm text-ink-muted">Нет платежей</li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Client } from '~/types'
import { useUiStore } from '~/stores/ui'

const api = useApi()
const ui = useUiStore()
const route = useRoute()
const fmt = useFormat()
const id = Number(route.params.id)

const client = ref<Client | null>(null)
const form = ref<Record<string, any>>({})
const history = ref<any>({ rentals: [], payments: [], accidents: [], penalties: [] })
const saving = ref(false)

const qrUrl = computed(() => (client.value?.qr_code_path ? api.fileUrl(client.value.qr_code_path) : ''))

async function load() {
  client.value = await api.clients.get(id)
  form.value = { ...client.value }
  history.value = await api.clients.history(id)
}

async function save() {
  saving.value = true
  try {
    client.value = await api.clients.update(id, form.value)
    ui.success('Данные сохранены')
  } catch {
    /* handled centrally */
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>
