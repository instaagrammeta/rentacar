<template>
  <div class="mx-auto max-w-3xl">
    <PageHeader title="Новый автомобиль" subtitle="Добавление авто в автопарк">
      <template #actions>
        <NuxtLink to="/cars" class="btn-secondary"><AppIcon name="back" size="18" /> Назад</NuxtLink>
      </template>
    </PageHeader>

    <form class="card p-6" @submit.prevent="submit">
      <CarForm v-model:form="form" v-model:photos="photos" />
      <div class="mt-6 flex justify-end gap-2">
        <NuxtLink to="/cars" class="btn-secondary">Отмена</NuxtLink>
        <button type="submit" class="btn-primary" :disabled="saving">{{ saving ? 'Сохранение…' : 'Создать' }}</button>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { useUiStore } from '~/stores/ui'

const api = useApi()
const ui = useUiStore()

const form = ref<Record<string, any>>({
  brand: '', model: '', year: new Date().getFullYear(), color: '', plate_number: '', vin: '',
  mileage: 0, status: 'available', daily_price: 0, weekly_price: 0, monthly_price: 0, deposit_amount: 0,
})
const photos = ref<string[]>([])
const saving = ref(false)

async function submit() {
  saving.value = true
  try {
    await api.cars.create({ ...form.value, photos: photos.value })
    ui.success('Автомобиль добавлен')
    await navigateTo('/cars')
  } catch {
    /* handled centrally */
  } finally {
    saving.value = false
  }
}
</script>
