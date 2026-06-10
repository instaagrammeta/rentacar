<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { auditApi } from '@/api'

const { t } = useI18n()

const items = ref([])
const loading = ref(false)

const headers = [
  { title: t('audit.date'), key: 'created_at' },
  { title: t('audit.user'), key: 'username' },
  { title: t('audit.action'), key: 'action' },
  { title: t('audit.entity'), key: 'entity' },
  { title: t('audit.details'), key: 'details' },
  { title: t('audit.ip'), key: 'ip_address' },
]

async function load() {
  loading.value = true
  try {
    const { data } = await auditApi.list({ limit: 500 })
    items.value = data
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<template>
  <div>
    <h1 class="text-h4 font-weight-bold mb-4">{{ t('audit.title') }}</h1>
    <v-card class="pa-4">
      <v-data-table :headers="headers" :items="items" :loading="loading" items-per-page="50">
        <template #[`item.created_at`]="{ item }">
          {{ new Date(item.created_at).toLocaleString('ru-RU') }}
        </template>
        <template #[`item.action`]="{ item }">
          <v-chip size="x-small" color="primary">{{ item.action }}</v-chip>
        </template>
      </v-data-table>
    </v-card>
  </div>
</template>
