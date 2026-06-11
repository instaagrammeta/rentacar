<template>
  <div>
    <PageHeader :title="rental?.contract_number || 'Договор'" subtitle="Детали аренды">
      <template #actions>
        <NuxtLink to="/rentals" class="btn-secondary"><AppIcon name="back" size="18" /> Назад</NuxtLink>
        <button class="btn-secondary" @click="downloadContract"><AppIcon name="download" size="18" /> Договор</button>
      </template>
    </PageHeader>

    <div v-if="rental" class="grid grid-cols-1 gap-5 lg:grid-cols-3">
      <!-- Info -->
      <div class="card p-6 lg:col-span-2">
        <div class="mb-4 flex items-center justify-between">
          <h3 class="font-semibold text-ink">Информация</h3>
          <StatusBadge :status="rental.status" :label="rental.status_label" />
        </div>
        <dl class="grid grid-cols-1 gap-x-6 gap-y-3 sm:grid-cols-2">
          <Info label="Клиент" :value="rental.client_name" />
          <Info label="Автомобиль" :value="rental.car_name" />
          <Info label="Период" :value="`${fmt.date(rental.rental_start)} — ${fmt.date(rental.rental_end)}`" />
          <Info label="Время выдачи" :value="fmt.dateTime(rental.pickup_at)" />
          <Info label="Срок возврата" :value="fmt.dateTime(rental.due_at)" />
          <Info label="Цена/сутки" :value="fmt.money(rental.daily_price)" />
          <Info label="Депозит" :value="fmt.money(rental.deposit)" />
          <Info label="Итого" :value="fmt.money(rental.total_price)" />
          <Info label="Сотрудник" :value="rental.employee_name" />
        </dl>

        <div v-if="rental.vehicle_return" class="mt-6 rounded-xl bg-surface-muted p-4">
          <h4 class="mb-2 font-semibold text-ink">Возврат оформлен</h4>
          <dl class="grid grid-cols-2 gap-3 text-sm sm:grid-cols-3">
            <Info label="Доп. дни" :value="rental.vehicle_return.extra_days" />
            <Info label="Пеня за просрочку" :value="fmt.money(rental.vehicle_return.late_fee)" />
            <Info label="Ущерб" :value="fmt.money(rental.vehicle_return.damage_cost)" />
            <Info label="Штрафы" :value="fmt.money(rental.vehicle_return.penalties)" />
            <Info label="К оплате" :value="fmt.money(rental.vehicle_return.final_payment)" />
          </dl>
        </div>

        <!-- QR code for the client -->
        <div class="mt-6 rounded-xl border border-surface-border p-4">
          <div class="flex items-center justify-between">
            <div>
              <h4 class="font-semibold text-ink">QR-код для клиента</h4>
              <p class="text-sm text-ink-muted">
                Клиент сканирует код и видит, сколько времени осталось до возврата.
              </p>
            </div>
            <button v-if="!qr" class="btn-secondary" :disabled="qrLoading" @click="loadQr">
              <AppIcon name="eye" size="18" /> Показать
            </button>
          </div>

          <div v-if="qr" class="mt-4 flex flex-col items-center gap-3">
            <img
              v-if="qr.qr_code_path"
              :src="api.fileUrl(qr.qr_code_path)"
              alt="QR-код аренды"
              class="h-48 w-48 rounded-lg border border-surface-border bg-white p-2"
            >
            <a
              :href="qr.public_url"
              target="_blank"
              rel="noopener"
              class="break-all text-center text-sm text-primary-600 underline"
            >
              {{ qr.public_url }}
            </a>
            <div class="flex gap-2">
              <button class="btn-secondary" @click="copyLink"><AppIcon name="check" size="18" /> Копировать ссылку</button>
              <button class="btn-secondary" @click="printQr"><AppIcon name="download" size="18" /> Печать</button>
            </div>
          </div>
        </div>
      </div>

      <!-- Return form -->
      <div class="card p-6">
        <h3 class="mb-4 font-semibold text-ink">Оформить возврат</h3>
        <template v-if="rental.has_return">
          <p class="rounded-xl bg-emerald-50 px-4 py-3 text-sm text-emerald-700">Возврат уже оформлен.</p>
        </template>
        <template v-else-if="rental.status === 'active'">
          <div class="space-y-3">
            <FormField v-model="ret.return_date" type="date" label="Дата возврата" />
            <FormField v-model="ret.mileage" type="number" label="Пробег при возврате" />
            <FormField v-model="ret.fuel_level" type="number" label="Уровень топлива, %" />
            <FormField v-model="ret.damage_cost" type="number" step="0.01" label="Стоимость ущерба" />
            <FormField v-model="ret.penalties" type="number" step="0.01" label="Штрафы" />
            <FormField v-model="ret.damages" type="textarea" label="Описание повреждений" />

            <div v-if="preview" class="rounded-xl bg-surface-muted p-3 text-sm">
              <p class="flex justify-between"><span>Доп. дни:</span><b>{{ preview.extra_days }}</b></p>
              <p class="flex justify-between"><span>Пеня:</span><b>{{ fmt.money(preview.late_fee) }}</b></p>
              <p class="flex justify-between border-t border-surface-border pt-1 mt-1"><span>Итого к оплате:</span><b>{{ fmt.money(preview.final_payment) }}</b></p>
            </div>

            <div class="flex gap-2">
              <button class="btn-secondary flex-1" @click="doPreview">Рассчитать</button>
              <button class="btn-primary flex-1" :disabled="saving" @click="submitReturn">Оформить</button>
            </div>
          </div>
        </template>
        <template v-else>
          <p class="text-sm text-ink-muted">Возврат недоступен для этого договора.</p>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Rental, VehicleReturn } from '~/types'
import { useUiStore } from '~/stores/ui'

const api = useApi()
const ui = useUiStore()
const route = useRoute()
const fmt = useFormat()
const id = Number(route.params.id)

const rental = ref<Rental | null>(null)
const preview = ref<VehicleReturn | null>(null)
const saving = ref(false)
const qr = ref<{ qr_code_path: string | null; public_token: string | null; public_url: string } | null>(null)
const qrLoading = ref(false)
const ret = ref<Record<string, any>>({
  return_date: new Date().toISOString().slice(0, 10),
  mileage: 0,
  fuel_level: 100,
  damage_cost: 0,
  penalties: 0,
  damages: '',
})

async function load() {
  rental.value = await api.rentals.get(id)
}

async function doPreview() {
  try {
    preview.value = await api.rentals.previewReturn(id, ret.value)
  } catch {
    /* handled centrally */
  }
}

async function submitReturn() {
  if (!confirm('Оформить возврат автомобиля?')) return
  saving.value = true
  try {
    await api.rentals.createReturn(id, ret.value)
    ui.success('Возврат оформлен')
    await load()
  } catch {
    /* handled centrally */
  } finally {
    saving.value = false
  }
}

async function downloadContract() {
  if (!rental.value) return
  try {
    await api.rentals.downloadContract(id, rental.value.contract_number)
  } catch {
    ui.error('Не удалось скачать договор')
  }
}

async function loadQr() {
  qrLoading.value = true
  try {
    qr.value = await api.rentals.qr(id)
  } catch {
    /* handled centrally */
  } finally {
    qrLoading.value = false
  }
}

async function copyLink() {
  if (!qr.value?.public_url) return
  try {
    await navigator.clipboard.writeText(qr.value.public_url)
    ui.success('Ссылка скопирована')
  } catch {
    ui.error('Не удалось скопировать ссылку')
  }
}

function printQr() {
  if (!qr.value?.qr_code_path) return
  const url = api.fileUrl(qr.value.qr_code_path)
  const w = window.open('', '_blank')
  if (!w) return
  w.document.write(
    `<html><head><title>QR — ${rental.value?.contract_number || ''}</title></head>` +
      `<body style="display:flex;flex-direction:column;align-items:center;justify-content:center;height:100vh;margin:0;font-family:sans-serif">` +
      `<h2 style="margin-bottom:8px">${rental.value?.contract_number || ''}</h2>` +
      `<img src="${url}" style="width:320px;height:320px" onload="window.print()">` +
      `<p style="margin-top:8px;color:#555">Сканируйте, чтобы увидеть оставшееся время аренды</p>` +
      `</body></html>`,
  )
  w.document.close()
}

onMounted(load)
</script>
