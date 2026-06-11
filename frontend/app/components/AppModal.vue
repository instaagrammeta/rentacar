<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="open" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <div class="absolute inset-0 bg-ink/50" @click="$emit('close')" />
        <div
          class="relative z-10 w-full overflow-hidden rounded-2xl bg-white shadow-card"
          :class="sizeClass"
        >
          <div class="flex items-center justify-between border-b border-surface-border px-5 py-4">
            <h3 class="text-base font-semibold text-ink">{{ title }}</h3>
            <button class="text-ink-muted hover:text-ink" @click="$emit('close')">
              <AppIcon name="close" />
            </button>
          </div>
          <div class="max-h-[70vh] overflow-y-auto px-5 py-4">
            <slot />
          </div>
          <div v-if="$slots.footer" class="flex justify-end gap-2 border-t border-surface-border px-5 py-4">
            <slot name="footer" />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{ open: boolean; title?: string; size?: 'sm' | 'md' | 'lg' }>(), {
  size: 'md',
})
defineEmits<{ close: [] }>()

const sizeClass = computed(
  () => ({ sm: 'max-w-md', md: 'max-w-xl', lg: 'max-w-3xl' })[props.size],
)
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s ease;
}
.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}
</style>
