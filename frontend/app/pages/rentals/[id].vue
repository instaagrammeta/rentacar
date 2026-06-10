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

onMounted(load)
</script>
