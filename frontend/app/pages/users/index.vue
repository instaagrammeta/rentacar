<template>
  <div>
    <PageHeader title="Пользователи" subtitle="Сотрудники и доступы">
      <template #actions>
        <button class="btn-primary" @click="openCreate"><AppIcon name="plus" size="18" /> Добавить</button>
      </template>
    </PageHeader>

    <DataTable :columns="columns" :rows="rows" :loading="loading" empty="Пользователей нет">
      <template #is_active="{ value }">
        <span class="badge" :class="value ? 'bg-emerald-50 text-emerald-700' : 'bg-red-50 text-red-700'">
          {{ value ? 'Активен' : 'Отключён' }}
        </span>
      </template>
      <template #actions="{ row }">
        <button class="btn-ghost !px-2 !py-1.5" title="Сменить пароль" @click="openPassword(row)"><AppIcon name="edit" size="18" /></button>
      </template>
    </DataTable>

    <AppModal :open="showCreate" title="Новый пользователь" @close="showCreate = false">
      <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <FormField v-model="form.username" label="Логин" required />
        <FormField v-model="form.full_name" label="ФИО" />
        <FormField v-model="form.email" type="email" label="E-mail" />
        <FormField v-model="form.role" type="select" label="Роль" :options="roleOptions" />
        <div class="sm:col-span-2"><FormField v-model="form.password" type="password" label="Пароль" required /></div>
      </div>
      <template #footer>
        <button class="btn-secondary" @click="showCreate = false">Отмена</button>
        <button class="btn-primary" :disabled="saving" @click="create">Создать</button>
      </template>
    </AppModal>

    <AppModal :open="showPassword" title="Смена пароля" size="sm" @close="showPassword = false">
      <FormField v-model="newPassword" type="password" label="Новый пароль" required />
      <template #footer>
        <button class="btn-secondary" @click="showPassword = false">Отмена</button>
        <button class="btn-primary" :disabled="saving" @click="changePassword">Сохранить</button>
      </template>
    </AppModal>
  </div>
</template>

<script setup lang="ts">
import type { User, Option } from '~/types'
import { useUiStore } from '~/stores/ui'

const api = useApi()
const ui = useUiStore()

const columns = [
  { key: 'username', label: 'Логин' },
  { key: 'full_name', label: 'ФИО' },
  { key: 'role_label', label: 'Роль' },
  { key: 'is_active', label: 'Статус' },
]

const rows = ref<User[]>([])
const loading = ref(false)
const roleOptions = ref<Option[]>([])

const showCreate = ref(false)
const showPassword = ref(false)
const saving = ref(false)
const form = ref<Record<string, any>>({ username: '', full_name: '', email: '', role: 'operator', password: '' })
const newPassword = ref('')
const selected = ref<User | null>(null)

async function load() {
  loading.value = true
  try {
    rows.value = await api.auth.listUsers()
    roleOptions.value = await api.auth.roles()
  } finally {
    loading.value = false
  }
}

function openCreate() {
  form.value = { username: '', full_name: '', email: '', role: 'operator', password: '' }
  showCreate.value = true
}

async function create() {
  saving.value = true
  try {
    await api.auth.createUser(form.value)
    ui.success('Пользователь создан')
    showCreate.value = false
    load()
  } catch {
    /* handled centrally */
  } finally {
    saving.value = false
  }
}

function openPassword(u: User) {
  selected.value = u
  newPassword.value = ''
  showPassword.value = true
}

async function changePassword() {
  if (!selected.value) return
  saving.value = true
  try {
    await api.auth.changePassword(selected.value.id, newPassword.value)
    ui.success('Пароль обновлён')
    showPassword.value = false
  } catch {
    /* handled centrally */
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>
