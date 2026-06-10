<template>
  <div class="flex items-center gap-5">
    <svg viewBox="0 0 36 36" class="h-32 w-32 flex-shrink-0">
      <circle cx="18" cy="18" r="15.915" fill="none" stroke="#f1f3f7" stroke-width="3.5" />
      <circle
        v-for="(seg, i) in computedSegments"
        :key="i"
        cx="18"
        cy="18"
        r="15.915"
        fill="none"
        :stroke="seg.color"
        stroke-width="3.5"
        :stroke-dasharray="`${seg.percent} ${100 - seg.percent}`"
        :stroke-dashoffset="seg.offset"
        stroke-linecap="round"
        transform="rotate(-90 18 18)"
      />
      <text x="18" y="17" text-anchor="middle" class="fill-ink text-[5px] font-bold">{{ total }}</text>
      <text x="18" y="22" text-anchor="middle" class="fill-current text-[2.6px]" style="fill:#8b93a3">всего</text>
    </svg>

    <ul class="flex-1 space-y-1.5">
      <li v-for="(seg, i) in segments" :key="i" class="flex items-center justify-between text-sm">
        <span class="flex items-center gap-2 text-ink-soft">
          <span class="h-2.5 w-2.5 rounded-full" :style="{ background: seg.color }" />
          {{ seg.label }}
        </span>
        <span class="font-semibold text-ink">{{ seg.value }}</span>
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{ segments: { label: string; value: number; color: string }[] }>()

const total = computed(() => props.segments.reduce((a, s) => a + s.value, 0))

const computedSegments = computed(() => {
  const t = total.value || 1
  let acc = 0
  return props.segments.map((s) => {
    const percent = (s.value / t) * 100
    const offset = 100 - acc + 25 // start at top (12 o'clock)
    acc += percent
    return { ...s, percent, offset }
  })
})
</script>
