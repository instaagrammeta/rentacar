<template>
  <div>
    <PageHeader title="Настройки" subtitle="Параметры компании и данные системы" />

    <div class="grid grid-cols-1 gap-5 lg:grid-cols-3">
      <!-- Company settings -->
      <form class="card p-6 lg:col-span-2" @submit.prevent="save">
        <h3 class="mb-4 font-semibold text-ink">Компания</h3>
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <FormField v-model="form.company_name" label="Название компании" />
          <FormField v-model="form.currency" label="Валюта" />
          <FormField v-model="form.phone" label="Телефон" />
          <FormField v-model="form.email" type="email" label="E-mail" />
          <div class="sm:col-span-2"><FormField v-model="form.address" type="textarea" label="Адрес" /></div>
          <div class="sm:col-span-2"><FormField v-model="form.contract_terms" type="textarea" label="Условия договора" :rows="5" /></div>
        </div>
        <div class="mt-6 flex justify-end">
          <button type="submit" class="btn-primary" :disabled="saving">{{ saving ? 'Сохранение…' : 'Сохранить' }}</button>
        </div>
      </form>

      <!-- Data management -->
      <div class="space-y-5">
        <div class="card p-5">
          <h3 class="mb-3 font-semibold text-ink">Резервные копии</h3>
          <button class="btn-secondary w-full" :disabled="busy" @click="createBackup">Создать копию</button>
          <ul class="mt-3 space-y-1.5 text-sm">
            <li v-for="b in backups" :key="b.name" class="flex items-center justify-between text-ink-soft">
              <span class="truncate">{{ b.name }}</span>
            </li>
            <li v-if="!backups.length" class="text-ink-muted">Копий нет</li>
          </ul>
        </div>

        <div class="card p-5">
          <h3 class="mb-3 font-semibold text-ink">Проект</h3>
          <button class="btn-secondary mb-2 w-full" @click="exportProject"><AppIcon name="download" size="18" /> Экспорт (.rentacar)</button>
          <label class="btn-secondary w-full cursor-pointer">
            Импорт проекта
            <input type="file" accept=".rentacar" class="hidden" @change="importProject" />
          </label>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useUiStore } from '~/stores/ui'

const api = useApi()
const ui = useUiStore()

const form = ref<Record<string, any>>({})
const backups = ref<any[]>([])
const saving = ref(false)
const busy = ref(false)

async function load() {
  form.value = await api.settings.get()
  try {
    backups.value = await api.backups.list()
  } catch {
    backups.value = []
  }
}

async function save() {
  saving.value = true
  try {
    form.value = await api.settings.update(form.value)
    ui.success('Настройки сохранены')
  } catch {
    /* handled centrally */
  } finally {
    saving.value = false
  }
}

async function createBackup() {
  busy.value = true
  try {
    await api.backups.create()
    ui.success('Резервная копия создана')
    backups.value = await api.backups.list()
  } catch {
    /* handled centrally */
  } finally {
    busy.value = false
  }
}

async function exportProject() {
  try {
    await api.project.exportProject()
  } catch {
    ui.error('Не удалось экспортировать проект')
  }
}

async function importProject(e: Event) {
  const input = e.target as HTMLInputElement
  if (!input.files?.length) return
  if (!confirm('Импорт заменит все текущие данные. Продолжить?')) return
  const fd = new FormData()
  fd.append('file', input.files[0])
  try {
    await api.project.import(fd)
    ui.success('Проект импортирован')
    await load()
  } catch {
    /* handled centrally */
  } finally {
    input.value = ''
  }
}

onMounted(load)
</script>
