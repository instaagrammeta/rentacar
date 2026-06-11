<template>
  <div class="min-h-screen bg-gradient-to-br from-emerald-50 via-white to-emerald-100 px-4 py-10">
    <div class="mx-auto w-full max-w-md">
      <!-- Brand header -->
      <div class="mb-6 flex items-center justify-center gap-2 text-emerald-800">
        <AppIcon name="car" size="28" />
        <span class="text-xl font-bold tracking-tight">Wheelzie</span>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="card p-8 text-center text-ink-muted">
        <p>Загрузка…</p>
      </div>

      <!-- Not found / error -->
      <div v-else-if="error" class="card p-8 text-center">
        <div class="mx-auto mb-3 flex h-14 w-14 items-center justify-center rounded-full bg-red-100 text-red-600">
          <AppIcon name="close" size="28" />
        </div>
        <h1 class="text-lg font-semibold text-ink">Аренда не найдена</h1>
        <p class="mt-2 text-sm text-ink-muted">Ссылка недействительна или срок её действия истёк.</p>
      </div>

      <!-- Content -->
      <div v-else-if="data" class="space-y-4">
        <!-- Countdown card -->
        <div class="card overflow-hidden">
          <div
            class="px-6 py-5 text-center text-white"
            :class="headerClass"
          >
            <p class="text-sm/none opacity-90">{{ statusHeadline }}</p>
            <p v-if="isActive" class="mt-3 font-mono text-4xl font-bold tabular-nums tracking-tight">
              {{ countdown }}
            </p>
            <p v-else class="mt-3 text-2xl font-bold">{{ data.status_label }}</p>
          </div>

          <div class="grid grid-cols-2 gap-px bg-surface-border">
            <div class="bg-white px-4 py-3 text-center">
              <p class="text-xs text-ink-muted">Выдан</p>
              <p class="mt-1 text-sm font-semibold text-ink">{{ fmt.dateTime(data.pickup_at) }}</p>
            </div>
            <div class="bg-white px-4 py-3 text-center">
              <p class="text-xs text-ink-muted">Вернуть до</p>
              <p class="mt-1 text-sm font-semibold" :class="data.is_overdue ? 'text-red-600' : 'text-ink'">
                {{ fmt.dateTime(data.due_at) }}
              </p>
            </div>
          </div>
        </div>

        <!-- Details card -->
        <div class="card p-6">
          <dl class="space-y-3 text-sm">
            <div class="flex justify-between gap-4">
              <dt class="text-ink-muted">Автомобиль</dt>
              <dd class="text-right font-medium text-ink">{{ data.car_name || '—' }}</dd>
            </div>
            <div class="flex justify-between gap-4">
              <dt class="text-ink-muted">Договор</dt>
              <dd class="text-right font-medium text-ink">{{ data.contract_number }}</dd>
            </div>
            <div v-if="data.client_name" class="flex justify-between gap-4">
              <dt class="text-ink-muted">Клиент</dt>
              <dd class="text-right font-medium text-ink">{{ data.client_name }}</dd>
            </div>
            <div class="flex justify-between gap-4">
              <dt class="text-ink-muted">Период</dt>
              <dd class="text-right font-medium text-ink">
                {{ fmt.date(data.rental_start) }} — {{ fmt.date(data.rental_end) }}
              </dd>
            </div>
          </dl>
        </div>

        <p v-if="data.is_overdue && isActive" class="rounded-xl bg-red-50 px-4 py-3 text-center text-sm font-medium text-red-700">
          Срок аренды истёк. Пожалуйста, верните автомобиль как можно скорее.
        </p>

        <p class="text-center text-xs text-ink-muted">Wheelzie — аренда автомобилей</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { PublicRentalView } from '~/types'

// Standalone public page — no app chrome, no authentication required.
definePageMeta({ layout: false })

const api = useApi()
const route = useRoute()
const fmt = useFormat()
const token = String(route.params.token)

const data = ref<PublicRentalView | null>(null)
const loading = ref(true)
const error = ref(false)

// Local countdown driven by the server-provided remaining seconds, so client
// clock skew never affects accuracy.
const remaining = ref(0)
let timer: ReturnType<typeof setInterval> | null = null

const isActive = computed(() => data.value?.status === 'active')

const statusHeadline = computed(() => {
  if (!data.value) return ''
  if (data.value.status !== 'active') return 'Статус аренды'
  return data.value.is_overdue && remaining.value <= 0 ? 'Просрочено на' : 'Осталось времени'
})

const headerClass = computed(() => {
  if (!data.value || data.value.status !== 'active') return 'bg-slate-500'
  if (remaining.value <= 0) return 'bg-red-600'
  if (remaining.value <= 3600) return 'bg-amber-500'
  return 'bg-emerald-600'
})

// Format the absolute remaining (or overdue) time as Дд ЧЧ:ММ:СС.
const countdown = computed(() => {
  const total = Math.abs(remaining.value)
  const days = Math.floor(total / 86400)
  const hours = Math.floor((total % 86400) / 3600)
  const minutes = Math.floor((total % 3600) / 60)
  const seconds = Math.floor(total % 60)
  const pad = (n: number) => String(n).padStart(2, '0')
  const hms = `${pad(hours)}:${pad(minutes)}:${pad(seconds)}`
  return days > 0 ? `${days}д ${hms}` : hms
})

function tick() {
  remaining.value -= 1
}

async function load() {
  loading.value = true
  error.value = false
  try {
    const view = await api.public.rental(token)
    data.value = view
    remaining.value = view.remaining_seconds
    if (timer) clearInterval(timer)
    if (view.status === 'active') {
      timer = setInterval(tick, 1000)
    }
  } catch {
    error.value = true
  } finally {
    loading.value = false
  }
}

onMounted(load)
onBeforeUnmount(() => {
  if (timer) clearInterval(timer)
})
</script>
