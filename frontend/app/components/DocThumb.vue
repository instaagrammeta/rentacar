<template>
  <div>
    <a
      v-if="path"
      :href="api.fileUrl(path)"
      target="_blank"
      rel="noopener"
      class="block h-24 overflow-hidden rounded-lg border border-surface-border"
    >
      <img :src="api.fileUrl(path)" class="h-full w-full object-cover transition hover:scale-105" alt="" />
    </a>
    <div
      v-else
      class="flex h-24 items-center justify-center rounded-lg border border-dashed border-surface-border bg-surface-muted text-xs text-ink-muted"
    >
      Нет фото
    </div>
    <p class="mt-1 truncate text-center text-xs text-ink-muted">{{ label }}</p>
  </div>
</template>

<script setup lang="ts">
defineProps<{ path?: string | null; label: string }>()
const api = useApi()
</script>
