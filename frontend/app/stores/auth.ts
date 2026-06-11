import { defineStore } from 'pinia'
import type { User } from '~/types'

// Mirror of the backend ROLE_PERMISSIONS matrix, used to drive menu visibility
// and client-side guards. The backend remains the source of truth.
const ROLE_PERMISSIONS: Record<string, string[]> = {
  dashboard: ['administrator', 'rental_manager', 'cashier', 'operator'],
  clients: ['administrator', 'rental_manager'],
  cars: ['administrator', 'rental_manager'],
  reservations: ['administrator', 'rental_manager', 'operator'],
  rentals: ['administrator', 'rental_manager'],
  returns: ['administrator', 'rental_manager'],
  payments: ['administrator', 'cashier'],
  blacklist: ['administrator', 'rental_manager'],
  accidents: ['administrator', 'rental_manager'],
  reports: ['administrator'],
  settings: ['administrator'],
  users: ['administrator'],
  backups: ['administrator'],
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null as User | null,
    token: import.meta.client ? localStorage.getItem('token') : null,
  }),

  getters: {
    isAuthenticated: (s) => !!s.token,
    role: (s) => s.user?.role ?? '',
    initials: (s) => {
      const name = s.user?.full_name || s.user?.username || ''
      return name
        .split(' ')
        .filter(Boolean)
        .slice(0, 2)
        .map((p) => p[0]?.toUpperCase())
        .join('')
    },
  },

  actions: {
    can(permission: string): boolean {
      const allowed = ROLE_PERMISSIONS[permission]
      if (!allowed) return true
      return allowed.includes(this.user?.role ?? '')
    },

    setToken(token: string | null) {
      this.token = token
      if (import.meta.client) {
        if (token) localStorage.setItem('token', token)
        else localStorage.removeItem('token')
      }
    },

    async login(username: string, password: string) {
      const api = useApi()
      const res = await api.auth.login({ username, password })
      this.setToken(res.access_token)
      this.user = res.user
      return res.user
    },

    async fetchMe() {
      if (!this.token) return null
      const api = useApi()
      try {
        this.user = await api.auth.me()
      } catch {
        this.logout()
      }
      return this.user
    },

    logout() {
      const api = useApi()
      if (this.token) {
        api.auth.logout().catch(() => {})
      }
      this.user = null
      this.setToken(null)
    },
  },
})
