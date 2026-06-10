<template>
  <div class="min-h-screen bg-surface-muted">
    <AppSidebar />

    <!-- Mobile overlay -->
    <div
      v-if="ui.sidebarOpen"
      class="fixed inset-0 z-30 bg-ink/40 lg:hidden"
      @click="ui.closeSidebar()"
    />

    <div class="flex min-h-screen flex-col lg:pl-64">
      <AppTopbar />
      <main class="flex-1 p-4 sm:p-6">
        <slot />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useUiStore } from '~/stores/ui'
const ui = useUiStore()
</script>
