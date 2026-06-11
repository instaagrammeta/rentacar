<template>
  <div class="pointer-events-none fixed bottom-4 right-4 z-50 flex w-80 max-w-[calc(100vw-2rem)] flex-col gap-2">
    <TransitionGroup name="toast">
      <div
        v-for="t in ui.toasts"
        :key="t.id"
        class="pointer-events-auto flex items-start gap-3 rounded-xl border bg-white px-4 py-3 shadow-card"
        :class="borderClass(t.type)"
      >
        <span class="mt-0.5 h-2.5 w-2.5 flex-shrink-0 rounded-full" :class="dotClass(t.type)" />
        <p class="flex-1 text-sm text-ink">{{ t.message }}</p>
        <button class="text-ink-muted hover:text-ink" @click="ui.dismiss(t.id)">
          <AppIcon name="close" size="16" />
        </button>
      </div>
    </TransitionGroup>
  </div>
</template>

<script setup lang="ts">
import { useUiStore } from '~/stores/ui'
const ui = useUiStore()

function dotClass(type: string) {
  return {
    success: 'bg-emerald-500',
    error: 'bg-red-500',
    warning: 'bg-amber-500',
    info: 'bg-primary-500',
  }[type] as string
}
function borderClass(type: string) {
  return {
    success: 'border-emerald-200',
    error: 'border-red-200',
    warning: 'border-amber-200',
    info: 'border-surface-border',
  }[type] as string
}
</script>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.25s ease;
}
.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateX(20px);
}
</style>
