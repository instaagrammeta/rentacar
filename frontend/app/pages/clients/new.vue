<template>
  <div class="mx-auto max-w-3xl">
    <PageHeader title="Новый клиент" subtitle="Регистрация клиента в системе">
      <template #actions>
        <NuxtLink to="/clients" class="btn-secondary"><AppIcon name="back" size="18" /> Назад</NuxtLink>
      </template>
    </PageHeader>

    <form class="card p-6" @submit.prevent="submit">
      <ClientForm v-model="form" />
      <div class="mt-6 flex justify-end gap-2">
        <NuxtLink to="/clients" class="btn-secondary">Отмена</NuxtLink>
        <button type="submit" class="btn-primary" :disabled="saving">{{ saving ? 'Сохранение…' : 'Создать клиента' }}</button>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { useUiStore } from '~/stores/ui'

const api = useApi()
const ui = useUiStore()

const form = ref<Record<string, any>>({
  first_name: '',
  last_name: '',
  phone: '',
  email: '',
  date_of_birth: '',
  driver_experience_years: 0,
  passport_number: '',
  driver_license_number: '',
  driver_license_issue_date: '',
  status: 'pending_verification',
  notes: '',
  is_vip: false,
})
const saving = ref(false)

async function submit() {
  saving.value = true
  try {
    const client = await api.clients.create(form.value)
    ui.success('Клиент создан')
    await navigateTo(`/clients/${client.id}`)
  } catch {
    /* handled centrally */
  } finally {
    saving.value = false
  }
}
</script>
