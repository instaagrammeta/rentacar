<template>
  <svg
    xmlns="http://www.w3.org/2000/svg"
    :width="size"
    :height="size"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    :stroke-width="strokeWidth"
    stroke-linecap="round"
    stroke-linejoin="round"
    aria-hidden="true"
  >
    <path :d="path" />
  </svg>
</template>

<script setup lang="ts">
const props = withDefaults(
  defineProps<{ name: string; size?: number | string; strokeWidth?: number | string }>(),
  { size: 20, strokeWidth: 1.8 },
)

// Minimal stroke-icon set (heroicons-style path data).
const ICONS: Record<string, string> = {
  dashboard: 'M4 5a1 1 0 0 1 1-1h5v7H4V5Zm0 8h6v6H5a1 1 0 0 1-1-1v-5Zm10 6v-7h6v6a1 1 0 0 1-1 1h-5Zm0-15h5a1 1 0 0 1 1 1v5h-6V4Z',
  clients: 'M16 14a4 4 0 1 0-8 0M12 11a3 3 0 1 0 0-6 3 3 0 0 0 0 6ZM3 20a6 6 0 0 1 12 0M17 14a5 5 0 0 1 4 6',
  cars: 'M5 16a2 2 0 1 0 0 .01M17 16a2 2 0 1 0 0 .01M3 16v-3l2-5h11l3 4 2 1v3M5 16h2m6 0h6',
  reservations: 'M8 3v3m8-3v3M4 8h16M5 6h14a1 1 0 0 1 1 1v12a1 1 0 0 1-1 1H5a1 1 0 0 1-1-1V7a1 1 0 0 1 1-1Zm3 8h2m4 0h2',
  rentals: 'M8 4h8a1 1 0 0 1 1 1v15l-5-3-5 3V5a1 1 0 0 1 1-1Zm1 4h6m-6 3h6',
  payments: 'M3 7a1 1 0 0 1 1-1h16a1 1 0 0 1 1 1v10a1 1 0 0 1-1 1H4a1 1 0 0 1-1-1V7Zm0 3h18M7 15h3',
  blacklist: 'M12 3a9 9 0 1 0 0 18 9 9 0 0 0 0-18Zm-6 3 12 12',
  accidents: 'M12 3 2 20h20L12 3Zm0 6v5m0 3v.01',
  reports: 'M4 19V5m0 14h16M8 16V9m4 7V6m4 10v-4',
  users: 'M9 13a3 3 0 1 0 0-6 3 3 0 0 0 0 6Zm-6 7a6 6 0 0 1 12 0M16 7a3 3 0 0 1 0 6m1 7a6 6 0 0 0-3-5',
  settings: 'M12 9a3 3 0 1 0 0 6 3 3 0 0 0 0-6Zm8 3a8 8 0 0 0-.2-1.8l2-1.5-2-3.4-2.3 1a8 8 0 0 0-3-1.7L14 1h-4l-.5 2.6a8 8 0 0 0-3 1.7l-2.3-1-2 3.4 2 1.5A8 8 0 0 0 4 12c0 .6.1 1.2.2 1.8l-2 1.5 2 3.4 2.3-1a8 8 0 0 0 3 1.7L10 23h4l.5-2.6a8 8 0 0 0 3-1.7l2.3 1 2-3.4-2-1.5c.1-.6.2-1.2.2-1.8Z',
  audit: 'M5 4h14a1 1 0 0 1 1 1v14a1 1 0 0 1-1 1H5a1 1 0 0 1-1-1V5a1 1 0 0 1 1-1Zm3 5h8M8 13h8M8 17h5',
  logout: 'M15 12H3m0 0 4-4m-4 4 4 4M9 4h8a1 1 0 0 1 1 1v14a1 1 0 0 1-1 1H9',
  menu: 'M4 6h16M4 12h16M4 18h16',
  search: 'M11 4a7 7 0 1 0 0 14 7 7 0 0 0 0-14Zm10 17-5-5',
  bell: 'M6 9a6 6 0 1 1 12 0c0 5 2 6 2 6H4s2-1 2-6Zm3 9a3 3 0 0 0 6 0',
  plus: 'M12 5v14M5 12h14',
  download: 'M12 3v12m0 0 4-4m-4 4-4-4M5 21h14',
  trash: 'M4 7h16M9 7V5a1 1 0 0 1 1-1h4a1 1 0 0 1 1 1v2m-8 0 1 13h6l1-13',
  edit: 'M4 20h4L19 9l-4-4L4 16v4ZM14 6l4 4',
  close: 'M6 6l12 12M18 6 6 18',
  check: 'M5 13l4 4L19 7',
  eye: 'M2 12s3.5-7 10-7 10 7 10 7-3.5 7-10 7-10-7-10-7Zm10 3a3 3 0 1 0 0-6 3 3 0 0 0 0 6Z',
  back: 'M15 18l-6-6 6-6',
  car: 'M5 16a2 2 0 1 0 0 .01M17 16a2 2 0 1 0 0 .01M3 16v-3l2-5h11l3 4 2 1v3M5 16h2m6 0h6',
  money: 'M12 1v22M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6',
  clock: 'M12 7v5l3 2M12 3a9 9 0 1 0 0 18 9 9 0 0 0 0-18Z',
}

const path = computed(() => ICONS[props.name] ?? ICONS.dashboard)
</script>
