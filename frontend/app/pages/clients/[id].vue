<template>
  <div>
    <PageHeader :title="form.last_name ? `${form.last_name} ${form.first_name}` : 'Клиент'" subtitle="Карточка клиента">
      <template #actions>
        <button class="btn-secondary" @click="openSms"><AppIcon name="bell" size="18" /> SMS</button>
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

      <!-- QR + documents -->
      <div class="space-y-5">
        <div class="card p-5 text-center">
          <p class="text-sm text-ink-muted">Код клиента</p>
          <p class="text-lg font-bold text-ink">{{ client?.client_code }}</p>
          <img v-if="qrUrl" :src="qrUrl" alt="QR" class="mx-auto mt-3 h-40 w-40 rounded-xl border border-surface-border" />
        </div>

        <div class="card p-5">
          <h3 class="mb-3 font-semibold text-ink">Документы</h3>
          <div class="grid grid-cols-2 gap-3">
            <DocThumb :path="client?.passport_front" label="Паспорт (лицо)" />
            <DocThumb :path="client?.passport_back" label="Паспорт (оборот)" />
            <DocThumb :path="client?.driver_license_scan" label="Вод. удостоверение" />
            <DocThumb :path="client?.driver_photo" label="Фото водителя" />
          </div>
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

    <!-- SMS modal -->
    <AppModal :open="showSms" title="Отправить SMS клиенту" @close="showSms = false">
      <p class="mb-3 text-sm text-ink-muted">Получатель: <b class="text-ink">{{ client?.full_name }}</b> ({{ client?.phone }})</p>
      <FormField v-model="smsText" type="textarea" label="Сообщение" :rows="4" />
      <template #footer>
        <button class="btn-secondary" @click="showSms = false">Отмена</button>
        <button class="btn-primary" :disabled="sending" @click="sendSms">Отправить</button>
      </template>
    </AppModal>
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

const showSms = ref(false)
const smsText = ref('')
const sending = ref(false)

const qrUrl = computed(() => (client.value?.qr_code_path ? api.fileUrl(client.value.qr_code_path) : ''))

async function load() {
  try {
    client.value = await api.clients.get(id)
    form.value = { ...client.value }
  } catch {
    return
  }
  // History is non-critical: a failure here must not block editing.
  try {
    history.value = await api.clients.history(id)
  } catch {
    /* ignore */
  }
}

// Build a clean payload with only the editable fields (avoids sending
// computed/read-only fields back to the API).
function buildPayload() {
  const f = form.value
  return {
    first_name: f.first_name,
    last_name: f.last_name,
    phone: f.phone,
    email: f.email ?? '',
    date_of_birth: f.date_of_birth ?? '',
    passport_number: f.passport_number ?? '',
    driver_license_number: f.driver_license_number ?? '',
    driver_license_issue_date: f.driver_license_issue_date ?? '',
    driver_experience_years: Number(f.driver_experience_years) || 0,
    status: f.status,
    notes: f.notes ?? '',
    is_vip: !!f.is_vip,
    passport_front: f.passport_front ?? null,
    passport_back: f.passport_back ?? null,
    driver_license_scan: f.driver_license_scan ?? null,
    driver_photo: f.driver_photo ?? null,
  }
}

async function save() {
  saving.value = true
  try {
    client.value = await api.clients.update(id, buildPayload())
    form.value = { ...client.value }
    ui.success('Данные сохранены')
  } catch {
    /* handled centrally */
  } finally {
    saving.value = false
  }
}

function openSms() {
  smsText.value = ''
  showSms.value = true
}

async function sendSms() {
  if (!smsText.value.trim()) return
  sending.value = true
  try {
    await api.raw(`/clients/${id}/sms`, { method: 'POST', body: { message: smsText.value } })
    ui.success('SMS отправлено')
    showSms.value = false
  } catch {
    /* handled centrally */
  } finally {
    sending.value = false
  }
}

onMounted(load)
</script>
