<template>
  <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
    <FormField v-model="form.last_name" label="Фамилия" required />
    <FormField v-model="form.first_name" label="Имя" required />
    <FormField v-model="form.phone" label="Телефон" placeholder="+7 (___) ___-__-__" required />
    <FormField v-model="form.email" type="email" label="E-mail" />
    <FormField v-model="form.date_of_birth" type="date" label="Дата рождения" />
    <FormField v-model="form.driver_experience_years" type="number" label="Стаж вождения, лет" />
    <FormField v-model="form.passport_number" label="Номер паспорта" />
    <FormField v-model="form.driver_license_number" label="Вод. удостоверение" />
    <FormField v-model="form.driver_license_issue_date" type="date" label="Дата выдачи ВУ" />
    <FormField v-model="form.status" type="select" label="Статус" :options="statusOptions" />
    <div class="sm:col-span-2">
      <FormField v-model="form.notes" type="textarea" label="Примечание" />
    </div>
    <FormField v-model="form.is_vip" type="checkbox" checkbox-label="VIP-клиент" />
  </div>
</template>

<script setup lang="ts">
import type { Option } from '~/types'

const form = defineModel<Record<string, any>>({ required: true })

const statusOptions: Option[] = [
  { value: 'active', label: 'Активен' },
  { value: 'pending_verification', label: 'Ожидает проверки' },
  { value: 'blacklisted', label: 'В чёрном списке' },
]
</script>
