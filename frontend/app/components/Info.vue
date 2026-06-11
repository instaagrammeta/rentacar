<template>
  <div>
    <dt class="text-xs uppercase tracking-wide text-ink-muted">{{ label }}</dt>
    <dd class="mt-0.5 text-sm font-medium text-ink">{{ value || value === 0 ? value : '—' }}</dd>
  </div>
</template>

<script setup lang="ts">
defineProps<{ label: string; value: string | number | null | undefined }>()
</script>
