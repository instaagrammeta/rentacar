<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { settingsApi, backupApi, uploadApi } from '@/api'
import { useAppStore } from '@/stores/app'

const { t } = useI18n()
const app = useAppStore()

const form = ref({})
const backups = ref([])
const saving = ref(false)

async function load() {
  const { data } = await settingsApi.get()
  form.value = data
  await loadBackups()
}

async function loadBackups() {
  const { data } = await backupApi.list()
  backups.value = data
}

async function save() {
  saving.value = true
  try {
    await settingsApi.update(form.value)
    app.notify(t('app.saved'))
  } finally {
    saving.value = false
  }
}

async function onLogo(event) {
  const file = event.target.files?.[0]
  if (!file) return
  const fd = new FormData()
  fd.append('file', file)
  fd.append('subfolder', 'logo')
  const { data } = await uploadApi.upload(fd)
  form.value.logo_path = data.paths[0]
  app.notify(t('app.saved'))
}

async function createBackup() {
  await backupApi.create()
  app.notify(t('app.saved'))
  await loadBackups()
}

function fileUrl(path) {
  return uploadApi.fileUrl(path)
}

onMounted(load)
</script>

<template>
  <div>
    <h1 class="text-h4 font-weight-bold mb-4">{{ t('settings.title') }}</h1>

    <v-row>
      <v-col cols="12" md="7">
        <v-card class="pa-4">
          <div class="text-subtitle-1 font-weight-bold mb-3">Информация о компании</div>
          <v-text-field v-model="form.company_name" :label="t('settings.companyName')" />
          <v-text-field v-model="form.address" :label="t('settings.address')" />
          <v-row>
            <v-col cols="12" md="6"><v-text-field v-model="form.phone" :label="t('settings.phone')" /></v-col>
            <v-col cols="12" md="6"><v-text-field v-model="form.email" :label="t('settings.email')" /></v-col>
          </v-row>
          <v-text-field v-model="form.currency" :label="t('settings.currency')" />
          <v-textarea v-model="form.contract_terms" :label="t('settings.contractTerms')" rows="3" />
          <v-file-input :label="t('settings.logo')" accept="image/*" @change="onLogo" />
          <v-img v-if="form.logo_path" :src="fileUrl(form.logo_path)" max-width="160" class="mb-3" />
          <v-btn color="primary" :loading="saving" prepend-icon="mdi-content-save" @click="save">{{ t('app.save') }}</v-btn>
        </v-card>
      </v-col>

      <v-col cols="12" md="5">
        <v-card class="pa-4">
          <div class="d-flex align-center mb-3">
            <div class="text-subtitle-1 font-weight-bold">{{ t('settings.backups') }}</div>
            <v-spacer />
            <v-btn color="secondary" size="small" prepend-icon="mdi-database-plus" @click="createBackup">{{ t('settings.createBackup') }}</v-btn>
          </div>
          <v-list density="compact">
            <v-list-item
              v-for="b in backups"
              :key="b.name"
              :title="b.name"
              :subtitle="`${(b.size / 1024).toFixed(1)} КБ — ${new Date(b.created_at).toLocaleString('ru-RU')}`"
              prepend-icon="mdi-database"
            />
            <v-list-item v-if="!backups.length" :title="t('app.noData')" />
          </v-list>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>
