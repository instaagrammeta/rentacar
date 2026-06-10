<template>
  <span class="badge" :class="cls">
    <span class="h-1.5 w-1.5 rounded-full" :class="dot" />
    {{ label }}
  </span>
</template>

<script setup lang="ts">
const props = defineProps<{ status: string; label: string }>()

// Map machine status values to colour schemes.
const SCHEMES: Record<string, { bg: string; dot: string }> = {
  // green
  available: { bg: 'bg-emerald-50 text-emerald-700', dot: 'bg-emerald-500' },
  active: { bg: 'bg-emerald-50 text-emerald-700', dot: 'bg-emerald-500' },
  confirmed: { bg: 'bg-emerald-50 text-emerald-700', dot: 'bg-emerald-500' },
  completed: { bg: 'bg-blue-50 text-blue-700', dot: 'bg-blue-500' },
  // amber
  reserved: { bg: 'bg-amber-50 text-amber-700', dot: 'bg-amber-500' },
  pending_verification: { bg: 'bg-amber-50 text-amber-700', dot: 'bg-amber-500' },
  maintenance: { bg: 'bg-amber-50 text-amber-700', dot: 'bg-amber-500' },
  // blue / violet
  rented: { bg: 'bg-violet-50 text-violet-700', dot: 'bg-violet-500' },
  // red
  blacklisted: { bg: 'bg-red-50 text-red-700', dot: 'bg-red-500' },
  cancelled: { bg: 'bg-red-50 text-red-700', dot: 'bg-red-500' },
}

const scheme = computed(() => SCHEMES[props.status] || { bg: 'bg-surface-muted text-ink-soft', dot: 'bg-ink-muted' })
const cls = computed(() => scheme.value.bg)
const dot = computed(() => scheme.value.dot)
</script>
