<template>
  <div class="mx-auto max-w-2xl">
    <PageHeader title="Отправка SMS" subtitle="Выберите клиента и отправьте сообщение" />

    <div class="card p-6">
      <!-- Client selection -->
      <label class="block">
        <span class="label">Клиент<span class="text-primary-500"> *</span></span>
        <div class="relative">
          <span class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-ink-muted"><AppIcon name="search" size="18" /></span>
          <input v-model="search" class="input pl-10" placeholder="Поиск по имени, телефону, коду…" @input="onSearch" />
        </div>
      </label>

      <!-- Results -->
      <div v-if="!selected" class="mt-2 max-h-64 overflow-y-auto rounded-xl border border-surface-border">
        <button
          v-for="c in clients"
          :key="c.id"
          type="button"
          class="flex w-full items-center justify-between border-b border-surface-border px-4 py-2.5 text-left text-sm last:border-0 hover:bg-surface-muted"
          @click="select(c)"
        >
          <span class="font-medium text-ink">{{ c.full_name }}</span>
          <span class="text-ink-muted">{{ c.phone }}</span>
        </button>
        <p v-if="!clients.length" class="px-4 py-6 text-center text-sm text-ink-muted">
          {{ loading ? 'Загрузка…' : 'Клиенты не найдены' }}
        </p>
      </div>

      <!-- Selected client -->
      <div v-else class="mt-3 flex items-center justify-between rounded-xl bg-surface-muted px-4 py-3">
        <div>
          <p class="font-semibold text-ink">{{ selected.full_name }}</p>
          <p class="text-sm text-ink-muted">{{ selected.phone }}</p>
        </div>
        <button class="btn-ghost !px-2 !py-1.5" title="Изменить" @click="clearSelection"><AppIcon name="close" size="18" /></button>
      </div>

      <!-- Templates -->
      <div class="mt-5">
        <span class="label">Быстрые шаблоны</span>
        <div class="flex flex-wrap gap-2">
          <button v-for="(t, i) in templates" :key="i" type="button" class="btn-secondary !px-3 !py-1.5 text-xs" @click="message = t">
            {{ t.length > 32 ? t.slice(0, 32) + '…' : t }}
          </button>
        </div>
      </div>

      <!-- Message -->
      <div class="mt-4">
        <FormField v-model="message" type="textarea" label="Сообщение" :rows="5" placeholder="Введите текст сообщения…" />
        <p class="mt-1 text-right text-xs text-ink-muted">{{ message.length }} симв.</p>
      </div>

      <div class="mt-4 flex justify-end gap-2">
        <button class="btn-primary" :disabled="!canSend || sending" @click="send">
          <AppIcon name="bell" size="18" /> {{ sending ? 'Отправка…' : 'Отправить SMS' }}
        </button>
      </div>
    </div>

    <p class="mt-4 rounded-xl bg-surface-muted px-4 py-3 text-center text-xs text-ink-muted">
      Автоматические напоминания отправляются клиентам за 1 час и за 30 минут до окончания аренды.
    </p>
  </div>
</template>

<script setup lang="ts">
import type { Client } from '~/types'
import { useUiStore } from '~/stores/ui'

const api = useApi()
const ui = useUiStore()

const clients = ref<Client[]>([])
const selected = ref<Client | null>(null)
const search = ref('')
const message = ref('')
const loading = ref(false)
const sending = ref(false)
let timer: any

const templates = [
  'Уважаемый клиент! Напоминаем о сроке возврата автомобиля.',
  'Ваш автомобиль готов к выдаче. Ждём вас!',
  'Благодарим за аренду автомобиля в нашем сервисе!',
]

const canSend = computed(() => !!selected.value && message.value.trim().length > 0)

async function load() {
  loading.value = true
  try {
    const res = await api.clients.list({ search: search.value, per_page: 50 })
    clients.value = res.items
  } finally {
    loading.value = false
  }
}

function onSearch() {
  clearTimeout(timer)
  timer = setTimeout(load, 300)
}

function select(c: Client) {
  selected.value = c
}

function clearSelection() {
  selected.value = null
}

async function send() {
  if (!canSend.value || !selected.value) return
  sending.value = true
  try {
    await api.sms.sendToClient(selected.value.id, message.value.trim())
    ui.success('SMS отправлено')
    message.value = ''
  } catch {
    /* error toast handled centrally */
  } finally {
    sending.value = false
  }
}

onMounted(load)
</script>
