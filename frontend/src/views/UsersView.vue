<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { authApi } from '@/api'
import { useAppStore } from '@/stores/app'

const { t } = useI18n()
const app = useAppStore()

const items = ref([])
const loading = ref(false)
const dialog = ref(false)
const pwdDialog = ref(false)
const form = ref({})
const pwdForm = ref({})
const saving = ref(false)

const roleOptions = [
  { value: 'administrator', title: t('users.roles.administrator') },
  { value: 'rental_manager', title: t('users.roles.rental_manager') },
  { value: 'cashier', title: t('users.roles.cashier') },
  { value: 'operator', title: t('users.roles.operator') },
]

const headers = [
  { title: t('users.username'), key: 'username' },
  { title: t('users.fullName'), key: 'full_name' },
  { title: t('users.role'), key: 'role_label' },
  { title: t('client.email'), key: 'email' },
  { title: t('app.actions'), key: 'actions', sortable: false },
]

async function load() {
  loading.value = true
  try {
    const { data } = await authApi.listUsers()
    items.value = data
  } finally {
    loading.value = false
  }
}

function openCreate() {
  form.value = { role: 'operator' }
  dialog.value = true
}

async function save() {
  saving.value = true
  try {
    await authApi.createUser(form.value)
    app.notify(t('app.created'))
    dialog.value = false
    await load()
  } finally {
    saving.value = false
  }
}

function openPassword(item) {
  pwdForm.value = { id: item.id, password: '' }
  pwdDialog.value = true
}

async function savePassword() {
  await authApi.changePassword(pwdForm.value.id, pwdForm.value.password)
  app.notify(t('app.saved'))
  pwdDialog.value = false
}

onMounted(load)
</script>

<template>
  <div>
    <div class="d-flex align-center mb-4">
      <h1 class="text-h4 font-weight-bold">{{ t('users.title') }}</h1>
      <v-spacer />
      <v-btn color="primary" prepend-icon="mdi-plus" @click="openCreate">{{ t('users.add') }}</v-btn>
    </div>

    <v-card class="pa-4">
      <v-data-table :headers="headers" :items="items" :loading="loading" items-per-page="20">
        <template #[`item.role_label`]="{ item }">
          <v-chip size="small" color="primary">{{ item.role_label }}</v-chip>
        </template>
        <template #[`item.actions`]="{ item }">
          <v-btn icon="mdi-lock-reset" size="small" variant="text" :title="t('users.changePassword')" @click="openPassword(item)" />
        </template>
      </v-data-table>
    </v-card>

    <v-dialog v-model="dialog" max-width="500">
      <v-card>
        <v-card-title class="bg-primary text-white">{{ t('users.add') }}</v-card-title>
        <v-card-text class="pt-4">
          <v-text-field v-model="form.username" :label="t('users.username')" />
          <v-text-field v-model="form.full_name" :label="t('users.fullName')" />
          <v-text-field v-model="form.email" :label="t('client.email')" />
          <v-text-field v-model="form.password" :label="t('users.password')" type="password" />
          <v-select v-model="form.role" :items="roleOptions" :label="t('users.role')" />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="dialog = false">{{ t('app.cancel') }}</v-btn>
          <v-btn color="primary" :loading="saving" @click="save">{{ t('app.create') }}</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="pwdDialog" max-width="400">
      <v-card>
        <v-card-title class="bg-primary text-white">{{ t('users.changePassword') }}</v-card-title>
        <v-card-text class="pt-4">
          <v-text-field v-model="pwdForm.password" :label="t('users.password')" type="password" />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="pwdDialog = false">{{ t('app.cancel') }}</v-btn>
          <v-btn color="primary" @click="savePassword">{{ t('app.save') }}</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>
