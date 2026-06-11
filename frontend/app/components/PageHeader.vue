<template>
  <div class="mb-5 flex flex-wrap items-center justify-between gap-3">
    <div>
      <h2 class="text-xl font-bold text-ink">{{ title }}</h2>
      <p v-if="subtitle" class="text-sm text-ink-muted">{{ subtitle }}</p>
    </div>
    <div class="flex items-center gap-2">
      <slot name="actions" />
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{ title: string; subtitle?: string }>()
</script>
