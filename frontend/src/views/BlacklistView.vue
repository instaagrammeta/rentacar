<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { blacklistApi, clientApi } from '@/api'
import { useAppStore } from '@/stores/app'

const { t } = useI18n()
const app = useAppStore()

const items = ref([])
const clients = ref([])
const loading = ref(false)
const dialog = ref(false)
const form = ref({})
const saving = ref(false)

const reasonOptions = [
  { value: 'fraud', title: t('blacklist.reasons.fraud') },
  { value: 'vehicle_damage', title: t('blacklist.reasons.vehicle_damage') },
  { value: 'non_payment', title: t('blacklist.reasons.non_payment') },
  { value: 'serious_violation', title: t('blacklist.reasons.serious_violation') },
]

const headers = [
  { title: t('blacklist.client'), key: 'client_name' },
  { title: t('blacklist.reason'), key: 'reason_label' },
  { title: t('blacklist.comment'), key: 'comment' },
  { title: t('app.actions'), key: 'actions', sortable: false },
]

async function load() {
  loading.value = true
  try {
    const { data } = await blacklistApi.list()
    items.value = data
  } finally {
    loading.value = false
  }
}

async function openCreate() {
  form.value = { reason: 'fraud' }
  const { data } = await clientApi.list({ per_page: 200 })
  clients.value = data.items.map((x) => ({ value: x.id, title: `${x.full_name} (${x.client_code})` }))
  dialog.value = true
}

async function save() {
  saving.value = true
  try {
    await blacklistApi.add(form.value)
    app.notify(t('app.created'))
    dialog.value = false
    await load()
  } finally {
    saving.value = false
  }
}

async function remove(item) {
  if (!confirm(t('blacklist.remove') + '?')) return
  await blacklistApi.remove(item.client_id)
  app.notify(t('app.saved'))
  await load()
}

onMounted(load)
</script>

<template>
  <div>
    <div class="d-flex align-center mb-4">
      <h1 class="text-h4 font-weight-bold">{{ t('blacklist.title') }}</h1>
      <v-spacer />
      <v-btn color="error" prepend-icon="mdi-account-cancel" @click="openCreate">{{ t('blacklist.add') }}</v-btn>
    </div>

    <v-card class="pa-4">
      <v-data-table :headers="headers" :items="items" :loading="loading" items-per-page="20">
        <template #[`item.reason_label`]="{ item }">
          <v-chip size="small" color="error">{{ item.reason_label }}</v-chip>
        </template>
        <template #[`item.actions`]="{ item }">
          <v-btn size="small" variant="text" color="success" @click="remove(item)">{{ t('blacklist.remove') }}</v-btn>
        </template>
      </v-data-table>
    </v-card>

    <v-dialog v-model="dialog" max-width="500">
      <v-card>
        <v-card-title class="bg-error text-white">{{ t('blacklist.add') }}</v-card-title>
        <v-card-text class="pt-4">
          <v-select v-model="form.client_id" :items="clients" :label="t('blacklist.client')" />
          <v-select v-model="form.reason" :items="reasonOptions" :label="t('blacklist.reason')" />
          <v-textarea v-model="form.comment" :label="t('blacklist.comment')" rows="3" />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="dialog = false">{{ t('app.cancel') }}</v-btn>
          <v-btn color="error" :loading="saving" @click="save">{{ t('app.save') }}</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>
