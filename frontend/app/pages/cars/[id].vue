<template>
  <div class="mx-auto max-w-3xl">
    <PageHeader :title="form.brand ? `${form.brand} ${form.model}` : 'Автомобиль'" subtitle="Редактирование авто">
      <template #actions>
        <NuxtLink to="/cars" class="btn-secondary"><AppIcon name="back" size="18" /> Назад</NuxtLink>
      </template>
    </PageHeader>

    <form class="card p-6" @submit.prevent="save">
      <CarForm v-model:form="form" v-model:photos="photos" />
      <div class="mt-6 flex justify-end gap-2">
        <NuxtLink to="/cars" class="btn-secondary">Отмена</NuxtLink>
        <button type="submit" class="btn-primary" :disabled="saving">{{ saving ? 'Сохранение…' : 'Сохранить' }}</button>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { useUiStore } from '~/stores/ui'

const api = useApi()
const ui = useUiStore()
const route = useRoute()
const id = Number(route.params.id)

const form = ref<Record<string, any>>({})
const photos = ref<string[]>([])
const saving = ref(false)

async function load() {
  const car = await api.cars.get(id)
  form.value = { ...car }
  photos.value = car.photos || []
}

async function save() {
  saving.value = true
  try {
    await api.cars.update(id, { ...form.value, photos: photos.value })
    ui.success('Изменения сохранены')
    await navigateTo('/cars')
  } catch {
    /* handled centrally */
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>
