<template>
  <div class="space-y-5">
    <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
      <FormField v-model="form.brand" label="Марка" required />
      <FormField v-model="form.model" label="Модель" required />
      <FormField v-model="form.year" type="number" label="Год выпуска" />
      <FormField v-model="form.color" label="Цвет" />
      <FormField v-model="form.plate_number" label="Гос. номер" required />
      <FormField v-model="form.vin" label="VIN" />
      <FormField v-model="form.mileage" type="number" label="Пробег, км" />
      <FormField v-model="form.status" type="select" label="Статус" :options="statusOptions" />
      <FormField v-model="form.daily_price" type="number" step="0.01" label="Цена / сутки" />
      <FormField v-model="form.weekly_price" type="number" step="0.01" label="Цена / неделя" />
      <FormField v-model="form.monthly_price" type="number" step="0.01" label="Цена / месяц" />
      <FormField v-model="form.deposit_amount" type="number" step="0.01" label="Депозит" />
    </div>
    <PhotoUploader v-model="photos" subfolder="cars" />
  </div>
</template>

<script setup lang="ts">
import type { Option } from '~/types'

const form = defineModel<Record<string, any>>('form', { required: true })
const photos = defineModel<string[]>('photos', { default: () => [] })

const statusOptions: Option[] = [
  { value: 'available', label: 'Доступен' },
  { value: 'reserved', label: 'Забронирован' },
  { value: 'rented', label: 'В аренде' },
  { value: 'maintenance', label: 'На обслуживании' },
]
</script>
