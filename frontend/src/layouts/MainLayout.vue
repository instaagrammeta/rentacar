<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useAppStore } from '@/stores/app'
import { projectApi } from '@/api'

const { t } = useI18n()
const router = useRouter()
const auth = useAuthStore()
const app = useAppStore()

const drawer = ref(true)
const importInput = ref(null)

// Navigation items with the permission required to see them.
const allItems = [
  { title: 'menu.dashboard', icon: 'mdi-view-dashboard', to: '/dashboard', perm: 'dashboard' },
  { title: 'menu.clients', icon: 'mdi-account-group', to: '/clients', perm: 'clients' },
  { title: 'menu.cars', icon: 'mdi-car', to: '/cars', perm: 'cars' },
  { title: 'menu.reservations', icon: 'mdi-calendar-clock', to: '/reservations', perm: 'reservations' },
  { title: 'menu.rentals', icon: 'mdi-file-document-edit', to: '/rentals', perm: 'rentals' },
  { title: 'menu.returns', icon: 'mdi-keyboard-return', to: '/returns', perm: 'returns' },
  { title: 'menu.payments', icon: 'mdi-cash-multiple', to: '/payments', perm: 'payments' },
  { title: 'menu.blacklist', icon: 'mdi-account-cancel', to: '/blacklist', perm: 'blacklist' },
  { title: 'menu.accidents', icon: 'mdi-car-emergency', to: '/accidents', perm: 'accidents' },
  { title: 'menu.reports', icon: 'mdi-chart-bar', to: '/reports', perm: 'reports' },
  { title: 'menu.settings', icon: 'mdi-cog', to: '/settings', perm: 'settings' },
  { title: 'menu.users', icon: 'mdi-account-key', to: '/users', perm: 'users' },
  { title: 'menu.audit', icon: 'mdi-history', to: '/audit', perm: 'audit' },
]

const items = computed(() => allItems.filter((i) => auth.can(i.perm)))
const isAdmin = computed(() => auth.role === 'administrator')

function logout() {
  auth.logout()
  router.push('/login')
}

// Файл → Сохранить проект (download a .rentacar file)
async function saveProject() {
  try {
    const res = await fetch(projectApi.exportUrl(), {
      method: 'POST',
      headers: { Authorization: `Bearer ${auth.token}` },
    })
    if (!res.ok) throw new Error('export failed')
    const blob = await res.blob()
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `rentacar_project_${new Date().toISOString().slice(0, 10)}.rentacar`
    a.click()
    URL.revokeObjectURL(url)
    app.notify('Проект сохранён в файл .rentacar')
  } catch (e) {
    app.notify('Не удалось сохранить проект', 'error')
  }
}

// Файл → Открыть проект (upload + import a .rentacar file)
function openProject() {
  importInput.value?.click()
}

async function onImportFile(event) {
  const file = event.target.files?.[0]
  if (!file) return
  const formData = new FormData()
  formData.append('file', file)
  try {
    await projectApi.import(formData)
    app.notify('Проект успешно импортирован. Обновите страницу.')
  } finally {
    event.target.value = ''
  }
}
</script>

<template>
  <v-navigation-drawer v-model="drawer" color="primary" theme="dark">
    <div class="pa-4 text-center">
      <v-icon size="40" icon="mdi-car-key" />
      <div class="text-h6 font-weight-bold mt-1">Rentacar CRM</div>
    </div>
    <v-divider />
    <v-list nav density="comfortable">
      <v-list-item
        v-for="item in items"
        :key="item.to"
        :to="item.to"
        :prepend-icon="item.icon"
        :title="t(item.title)"
      />
    </v-list>
  </v-navigation-drawer>

  <v-app-bar color="surface" flat border>
    <v-app-bar-nav-icon @click="drawer = !drawer" />
    <v-toolbar-title class="font-weight-bold text-primary">Rentacar CRM</v-toolbar-title>
    <v-spacer />

    <!-- Файл menu -->
    <v-menu v-if="isAdmin">
      <template #activator="{ props }">
        <v-btn v-bind="props" prepend-icon="mdi-folder" variant="text">{{ t('menu.file') }}</v-btn>
      </template>
      <v-list>
        <v-list-item prepend-icon="mdi-content-save" :title="t('menu.saveProject')" @click="saveProject" />
        <v-list-item prepend-icon="mdi-folder-open" :title="t('menu.openProject')" @click="openProject" />
      </v-list>
    </v-menu>

    <v-chip class="mr-2" color="primary" variant="tonal">
      <v-icon start icon="mdi-account" />
      {{ auth.user?.full_name || auth.user?.username }} — {{ auth.user?.role_label }}
    </v-chip>
    <v-btn icon="mdi-logout" :title="t('menu.logout')" @click="logout" />

    <input ref="importInput" type="file" accept=".rentacar" class="d-none" @change="onImportFile" />
  </v-app-bar>

  <v-main>
    <v-container fluid class="pa-6">
      <router-view />
    </v-container>
  </v-main>
</template>
