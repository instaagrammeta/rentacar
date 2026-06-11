<template>
  <div>
    <PageHeader title="Журнал аудита" subtitle="История действий пользователей" />

    <DataTable :columns="columns" :rows="rows" :loading="loading" empty="Записей нет">
      <template #created_at="{ value }">{{ fmt.dateTime(value) }}</template>
    </DataTable>
  </div>
</template>

<script setup lang="ts">
const api = useApi()
const fmt = useFormat()

const columns = [
  { key: 'created_at', label: 'Дата' },
  { key: 'username', label: 'Пользователь' },
  { key: 'action', label: 'Действие' },
  { key: 'entity', label: 'Объект' },
  { key: 'details', label: 'Детали' },
]

const rows = ref<any[]>([])
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    rows.value = await api.audit.list({ limit: 200 })
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>
