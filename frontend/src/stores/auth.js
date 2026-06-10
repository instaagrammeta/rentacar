// Authentication store: token persistence, current user and role helpers.
import { defineStore } from 'pinia'
import { authApi } from '@/api'

// Module access matrix mirrors the backend ROLE_PERMISSIONS so the UI hides
// menu items the current role cannot use.
const ROLE_PERMISSIONS = {
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
  audit: ['administrator'],
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('rentacar_token') || null,
    user: JSON.parse(localStorage.getItem('rentacar_user') || 'null'),
  }),
  getters: {
    isAuthenticated: (state) => !!state.token,
    role: (state) => state.user?.role || null,
  },
  actions: {
    async login(username, password) {
      const { data } = await authApi.login({ username, password })
      this.token = data.access_token
      this.user = data.user
      localStorage.setItem('rentacar_token', data.access_token)
      localStorage.setItem('rentacar_user', JSON.stringify(data.user))
      return data.user
    },
    logout() {
      this.token = null
      this.user = null
      localStorage.removeItem('rentacar_token')
      localStorage.removeItem('rentacar_user')
    },
    can(permission) {
      if (!this.role) return false
      const allowed = ROLE_PERMISSIONS[permission]
      return allowed ? allowed.includes(this.role) : true
    },
  },
})
