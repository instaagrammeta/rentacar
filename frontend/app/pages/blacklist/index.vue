<template>
  <div>
    <PageHeader title="Чёрный список" subtitle="Клиенты с ограничениями">
      <template #actions>
        <button class="btn-primary" @click="openAdd"><AppIcon name="plus" size="18" /> Добавить</button>
      </template>
    </PageHeader>

    <DataTable :columns="columns" :rows="rows" :loading="loading" empty="Чёрный список пуст">
      <template #created_at="{ row }">{{ fmt.date(row.created_at) }}</template>
      <template #actions="{ row }">
        <button class="btn-ghost !px-2 !py-1.5 text-emerald-600" title="Убрать из ЧС" @click="remove(row)"><AppIcon name="check" size="18" /></button>
      </template>
    </DataTable>

    <AppModal :open="showAdd" title="Добавить в чёрный список" @close="showAdd = false">
      <div class="space-y-4">
        <FormField v-model="form.client_id" type="select" label="Клиент" :options="clientOptions" placeholder="Выберите клиента" />
        <FormField v-model="form.reason" type="select" label="Причина" :options="reasonOptions" placeholder="Выберите причину" />
        <FormField v-model="form.comment" type="textarea" label="Комментарий" />
      </div>
      <template #footer>
        <button class="btn-secondary" @click="showAdd = false">Отмена</button>
        <button class="btn-danger" :disabled="saving" @click="add">Добавить</button>
      </template>
    </AppModal>
  </div>
</template>

<script setup lang="ts">
import type { BlacklistEntry, Option } from '~/types'
import { useUiStore } from '~/stores/ui'

const api = useApi()
const ui = useUiStore()
const fmt = useFormat()

const columns = [
  { key: 'client_name', label: 'Клиент' },
  { key: 'reason_label', label: 'Причина' },
  { key: 'comment', label: 'Комментарий' },
  { key: 'created_at', label: 'Дата' },
]

const rows = ref<BlacklistEntry[]>([])
const loading = ref(false)
const showAdd = ref(false)
const saving = ref(false)
const form = ref<Record<string, any>>({ client_id: '', reason: '', comment: '' })
const clientOptions = ref<Option[]>([])
const reasonOptions = ref<Option[]>([])

async function load() {
  loading.value = true
  try {
    rows.value = await api.blacklist.list()
  } finally {
    loading.value = false
  }
}

async function openAdd() {
  showAdd.value = true
  form.value = { client_id: '', reason: '', comment: '' }
  const [clients, reasons] = await Promise.all([api.clients.list({ per_page: 200 }), api.blacklist.reasons()])
  clientOptions.value = clients.items.map((c) => ({ value: String(c.id), label: `${c.full_name} (${c.client_code})` }))
  reasonOptions.value = reasons
}

async function add() {
  saving.value = true
  try {
    await api.blacklist.add({ ...form.value, client_id: Number(form.value.client_id) })
    ui.success('Клиент добавлен в чёрный список')
    showAdd.value = false
    load()
  } catch {
    /* handled centrally */
  } finally {
    saving.value = false
  }
}

async function remove(row: BlacklistEntry) {
  if (!confirm('Убрать клиента из чёрного списка?')) return
  try {
    await api.blacklist.remove(row.client_id)
    ui.success('Клиент убран из чёрного списка')
    load()
  } catch {
    /* handled centrally */
  }
}

onMounted(load)
</script>
