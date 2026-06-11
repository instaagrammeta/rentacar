<template>
  <div class="flex items-center gap-2.5">
    <span class="relative flex h-9 w-9 items-center justify-center rounded-full bg-primary-500">
      <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2">
        <circle cx="12" cy="12" r="9" />
        <circle cx="12" cy="12" r="2.5" fill="white" />
        <path d="M12 3v6M12 15v6M3 12h6M15 12h6" />
      </svg>
    </span>
    <span v-if="!compact" class="text-lg font-extrabold tracking-tight" :class="dark ? 'text-white' : 'text-ink'">
      Wheelzie
      <span class="ml-1 rounded-md bg-primary-500/15 px-1.5 py-0.5 align-middle text-[10px] font-bold text-primary-500">v2.0</span>
    </span>
  </div>
</template>

<script setup lang="ts">
defineProps<{ compact?: boolean; dark?: boolean }>()
</script>
