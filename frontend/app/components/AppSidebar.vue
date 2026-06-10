<template>
  <aside
    class="fixed inset-y-0 left-0 z-40 w-64 transform border-r border-surface-border bg-white transition-transform duration-200 lg:translate-x-0"
    :class="ui.sidebarOpen ? 'translate-x-0' : '-translate-x-full'"
  >
    <div class="flex h-16 items-center justify-between border-b border-surface-border px-5">
      <AppLogo />
      <button class="lg:hidden text-ink-muted" @click="ui.closeSidebar()">
        <AppIcon name="close" />
      </button>
    </div>

    <nav class="flex flex-col gap-1 p-3">
      <NuxtLink
        v-for="item in visibleItems"
        :key="item.to"
        :to="item.to"
        class="nav-link"
        :class="{ 'nav-link-active': isActive(item.to) }"
        @click="ui.closeSidebar()"
      >
        <AppIcon :name="item.icon" />
        <span>{{ item.label }}</span>
      </NuxtLink>
    </nav>

    <div class="absolute inset-x-0 bottom-0 border-t border-surface-border p-3">
      <button class="nav-link w-full text-left" @click="onLogout">
        <AppIcon name="logout" />
        <span>Выйти</span>
      </button>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'
import { useUiStore } from '~/stores/ui'

const auth = useAuthStore()
const ui = useUiStore()
const route = useRoute()

interface NavItem {
  to: string
  label: string
  icon: string
  perm: string
}

const items: NavItem[] = [
  { to: '/', label: 'Дашборд', icon: 'dashboard', perm: 'dashboard' },
  { to: '/clients', label: 'Клиенты', icon: 'clients', perm: 'clients' },
  { to: '/cars', label: 'Автомобили', icon: 'cars', perm: 'cars' },
  { to: '/reservations', label: 'Брони', icon: 'reservations', perm: 'reservations' },
  { to: '/rentals', label: 'Аренды', icon: 'rentals', perm: 'rentals' },
  { to: '/payments', label: 'Платежи', icon: 'payments', perm: 'payments' },
  { to: '/blacklist', label: 'Чёрный список', icon: 'blacklist', perm: 'blacklist' },
  { to: '/accidents', label: 'ДТП', icon: 'accidents', perm: 'accidents' },
  { to: '/reports', label: 'Отчёты', icon: 'reports', perm: 'reports' },
  { to: '/users', label: 'Пользователи', icon: 'users', perm: 'users' },
  { to: '/settings', label: 'Настройки', icon: 'settings', perm: 'settings' },
  { to: '/audit', label: 'Аудит', icon: 'audit', perm: 'users' },
]

const visibleItems = computed(() => items.filter((i) => auth.can(i.perm)))

function isActive(to: string) {
  if (to === '/') return route.path === '/'
  return route.path === to || route.path.startsWith(to + '/')
}

function onLogout() {
  auth.logout()
  navigateTo('/login')
}
</script>
