<template>
  <aside
    class="fixed inset-y-0 left-0 z-40 w-64 transform bg-night text-gray-300 transition-transform duration-200 lg:translate-x-0"
    :class="ui.sidebarOpen ? 'translate-x-0' : '-translate-x-full'"
  >
    <div class="flex h-16 items-center justify-between border-b border-white/10 px-5">
      <AppLogo dark />
      <button class="lg:hidden text-gray-400" @click="ui.closeSidebar()">
        <AppIcon name="close" />
      </button>
    </div>

    <nav class="flex flex-col gap-1 p-3">
      <NuxtLink
        v-for="item in visibleItems"
        :key="item.to"
        :to="item.to"
        class="flex items-center gap-3 rounded-xl px-3.5 py-2.5 text-sm font-medium transition"
        :class="isActive(item.to)
          ? 'bg-primary-500 text-white shadow-lg shadow-primary-500/20'
          : 'text-gray-400 hover:bg-white/5 hover:text-white'"
        @click="ui.closeSidebar()"
      >
        <AppIcon :name="item.icon" />
        <span>{{ item.label }}</span>
      </NuxtLink>
    </nav>

    <div class="absolute inset-x-0 bottom-0 border-t border-white/10 p-3">
      <button
        class="flex w-full items-center gap-3 rounded-xl px-3.5 py-2.5 text-left text-sm font-medium text-gray-400 transition hover:bg-white/5 hover:text-white"
        @click="onLogout"
      >
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
  { to: '/sms', label: 'SMS-рассылка', icon: 'bell', perm: 'clients' },
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
