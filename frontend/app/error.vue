<template>
  <div class="flex min-h-screen flex-col items-center justify-center gap-4 bg-surface-muted p-6 text-center">
    <AppLogo />
    <h1 class="text-5xl font-extrabold text-primary-500">{{ error?.statusCode || 500 }}</h1>
    <p class="text-ink-soft">{{ error?.statusMessage || 'Что-то пошло не так' }}</p>
    <button class="btn-primary" @click="handleError">На главную</button>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{ error: { statusCode?: number; statusMessage?: string } }>()
void props
function handleError() {
  clearError({ redirect: '/' })
}
</script>
