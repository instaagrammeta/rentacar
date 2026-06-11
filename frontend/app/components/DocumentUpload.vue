<template>
  <div>
    <span class="label">{{ label }}</span>
    <div
      class="relative flex h-36 items-center justify-center overflow-hidden rounded-xl border border-dashed border-surface-border bg-surface-muted"
    >
      <template v-if="path">
        <img :src="api.fileUrl(path)" class="h-full w-full object-contain" alt="" />
        <button
          type="button"
          class="absolute right-2 top-2 rounded-full bg-white/90 p-1.5 text-red-500 shadow"
          @click="clear"
        >
          <AppIcon name="trash" size="16" />
        </button>
      </template>
      <label
        v-else
        class="flex h-full w-full cursor-pointer flex-col items-center justify-center gap-1 text-ink-muted hover:text-primary-500"
      >
        <AppIcon name="plus" size="22" />
        <span class="text-xs">{{ uploading ? 'Загрузка…' : 'Загрузить фото' }}</span>
        <input type="file" accept="image/*" class="hidden" @change="onSelect" />
      </label>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useUiStore } from '~/stores/ui'

defineProps<{ label: string; subfolder?: string }>()
const path = defineModel<string | null>({ default: null })

const api = useApi()
const ui = useUiStore()
const uploading = ref(false)

async function onSelect(e: Event) {
  const input = e.target as HTMLInputElement
  if (!input.files?.length) return
  uploading.value = true
  try {
    const form = new FormData()
    form.append('files', input.files[0])
    const res = await api.uploads.upload(form)
    if (res.paths?.length) path.value = res.paths[0]
  } catch {
    ui.error('Не удалось загрузить фото')
  } finally {
    uploading.value = false
    input.value = ''
  }
}

function clear() {
  path.value = null
}
</script>
