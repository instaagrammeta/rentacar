<template>
  <header class="sticky top-0 z-20 flex h-16 items-center gap-3 border-b border-surface-border bg-white/90 px-4 backdrop-blur sm:px-6">
    <button class="lg:hidden text-ink-soft" @click="ui.toggleSidebar()">
      <AppIcon name="menu" size="22" />
    </button>

    <h1 class="text-base font-semibold text-ink sm:text-lg">{{ title }}</h1>

    <div class="ml-auto flex items-center gap-2 sm:gap-3">
      <span class="hidden text-right sm:block">
        <span class="block text-sm font-semibold text-ink">{{ auth.user?.full_name || auth.user?.username }}</span>
        <span class="block text-xs text-ink-muted">{{ auth.user?.role_label }}</span>
      </span>
      <span class="flex h-9 w-9 items-center justify-center rounded-full bg-primary-100 text-sm font-bold text-primary-600">
        {{ auth.initials || 'U' }}
      </span>
    </div>
  </header>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'
import { useUiStore } from '~/stores/ui'

const auth = useAuthStore()
const ui = useUiStore()
const route = useRoute()

const TITLES: Record<string, string> = {
  '/': 'Дашборд',
  '/clients': 'Клиенты',
  '/cars': 'Автомобили',
  '/reservations': 'Брони',
  '/rentals': 'Аренды',
  '/payments': 'Платежи',
  '/blacklist': 'Чёрный список',
  '/accidents': 'ДТП',
  '/reports': 'Отчёты',
  '/users': 'Пользователи',
  '/settings': 'Настройки',
  '/audit': 'Аудит',
}

const title = computed(() => {
  const path = route.path
  if (TITLES[path]) return TITLES[path]
  const base = '/' + (path.split('/')[1] || '')
  return TITLES[base] || 'Wheelzie'
})
</script>
