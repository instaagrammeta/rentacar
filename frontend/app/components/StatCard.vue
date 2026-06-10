<template>
  <div class="card flex items-center gap-4 p-5">
    <span class="flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-xl" :class="accentBg">
      <AppIcon :name="icon" size="22" :class="accentText" />
    </span>
    <div class="min-w-0">
      <p class="truncate text-sm text-ink-muted">{{ title }}</p>
      <p class="text-2xl font-bold text-ink">{{ value }}</p>
      <p v-if="subtitle" class="text-xs text-ink-muted">{{ subtitle }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = withDefaults(
  defineProps<{ title: string; value: string | number; icon: string; subtitle?: string; accent?: string }>(),
  { accent: 'primary' },
)

const MAP: Record<string, { bg: string; text: string }> = {
  primary: { bg: 'bg-primary-100', text: 'text-primary-600' },
  green: { bg: 'bg-emerald-100', text: 'text-emerald-600' },
  blue: { bg: 'bg-blue-100', text: 'text-blue-600' },
  amber: { bg: 'bg-amber-100', text: 'text-amber-600' },
  violet: { bg: 'bg-violet-100', text: 'text-violet-600' },
}

const accentBg = computed(() => (MAP[props.accent] || MAP.primary).bg)
const accentText = computed(() => (MAP[props.accent] || MAP.primary).text)
</script>
