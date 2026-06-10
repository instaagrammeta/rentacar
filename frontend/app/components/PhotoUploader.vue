<template>
  <div>
    <span class="label">{{ label || 'Фотографии' }}</span>
    <div class="flex flex-wrap gap-3">
      <div
        v-for="(p, i) in photos"
        :key="i"
        class="relative h-24 w-24 overflow-hidden rounded-xl border border-surface-border"
      >
        <img :src="api.fileUrl(p)" class="h-full w-full object-cover" alt="" />
        <button
          type="button"
          class="absolute right-1 top-1 rounded-full bg-white/90 p-1 text-red-500 shadow"
          @click="removeAt(i)"
        >
          <AppIcon name="close" size="14" />
        </button>
      </div>

      <label class="flex h-24 w-24 cursor-pointer flex-col items-center justify-center gap-1 rounded-xl border border-dashed border-surface-border text-ink-muted hover:border-primary-300 hover:text-primary-500">
        <AppIcon name="plus" />
        <span class="text-xs">{{ uploading ? '…' : 'Добавить' }}</span>
        <input type="file" accept="image/*" multiple class="hidden" @change="onSelect" />
      </label>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useUiStore } from '~/stores/ui'

defineProps<{ label?: string; subfolder?: string }>()
const photos = defineModel<string[]>({ default: () => [] })

const api = useApi()
const ui = useUiStore()
const uploading = ref(false)

async function onSelect(e: Event) {
  const input = e.target as HTMLInputElement
  if (!input.files?.length) return
  uploading.value = true
  try {
    const form = new FormData()
    Array.from(input.files).forEach((f) => form.append('files', f))
    const res = await api.uploads.upload(form)
    photos.value = [...photos.value, ...res.paths]
  } catch {
    ui.error('Не удалось загрузить фото')
  } finally {
    uploading.value = false
    input.value = ''
  }
}

function removeAt(i: number) {
  photos.value = photos.value.filter((_, idx) => idx !== i)
}
</script>
