<template>
  <div class="flex h-56 items-end gap-2">
    <div v-for="(d, i) in data" :key="i" class="flex flex-1 flex-col items-center gap-2">
      <div class="flex w-full flex-1 items-end justify-center">
        <div
          class="w-full max-w-[36px] rounded-t-lg transition-all"
          :class="i === highlightIndex ? 'bg-primary-500' : 'bg-primary-100'"
          :style="{ height: barHeight(d.value) }"
          :title="`${d.label}: ${d.value}`"
        />
      </div>
      <span class="truncate text-[10px] text-ink-muted">{{ d.label }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{ data: { label: string; value: number }[]; highlightIndex?: number }>()

const max = computed(() => Math.max(...props.data.map((d) => d.value), 1))
function barHeight(v: number) {
  return `${Math.max(4, (v / max.value) * 100)}%`
}
</script>
