<template>
  <div>
    <svg :viewBox="`0 0 ${W} ${H}`" class="h-56 w-full" preserveAspectRatio="none">
      <defs>
        <linearGradient :id="gid" x1="0" y1="0" x2="0" y2="1">
          <stop offset="0%" stop-color="#f5452b" stop-opacity="0.25" />
          <stop offset="100%" stop-color="#f5452b" stop-opacity="0" />
        </linearGradient>
      </defs>

      <!-- gridlines -->
      <line v-for="i in 4" :key="i" :x1="0" :x2="W" :y1="(H / 4) * i" :y2="(H / 4) * i"
        stroke="#eceef3" stroke-width="1" />

      <template v-if="points.length">
        <path :d="areaPath" :fill="`url(#${gid})`" />
        <polyline :points="linePoints" fill="none" stroke="#f5452b" stroke-width="2.5"
          stroke-linejoin="round" stroke-linecap="round" />
        <circle v-for="(p, i) in points" :key="i" :cx="p.x" :cy="p.y" r="3" fill="#fff"
          stroke="#f5452b" stroke-width="2" />
      </template>
    </svg>

    <div class="mt-2 flex justify-between text-[10px] text-ink-muted">
      <span v-for="(d, i) in data" :key="i" class="flex-1 text-center">{{ shortMonth(d.month) }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{ data: { month: string; revenue: number }[] }>()

const W = 300
const H = 160
const gid = 'lg-' + Math.random().toString(36).slice(2, 8)

const points = computed(() => {
  const d = props.data
  if (!d.length) return [] as { x: number; y: number }[]
  const max = Math.max(...d.map((x) => x.revenue), 1)
  const stepX = d.length > 1 ? W / (d.length - 1) : W
  return d.map((x, i) => ({
    x: d.length > 1 ? i * stepX : W / 2,
    y: H - (x.revenue / max) * (H - 20) - 8,
  }))
})

const linePoints = computed(() => points.value.map((p) => `${p.x},${p.y}`).join(' '))
const areaPath = computed(() => {
  if (!points.value.length) return ''
  const first = points.value[0]
  const last = points.value[points.value.length - 1]
  return `M ${first.x},${H} L ` + points.value.map((p) => `${p.x},${p.y}`).join(' L ') + ` L ${last.x},${H} Z`
})

function shortMonth(m: string) {
  const part = (m || '').split('-')[1] || ''
  return ['', 'Янв', 'Фев', 'Мар', 'Апр', 'Май', 'Июн', 'Июл', 'Авг', 'Сен', 'Окт', 'Ноя', 'Дек'][Number(part)] || m
}
</script>
