import { defineStore } from 'pinia'

export type ToastType = 'success' | 'error' | 'warning' | 'info'

export interface Toast {
  id: number
  message: string
  type: ToastType
}

let counter = 0

export const useUiStore = defineStore('ui', {
  state: () => ({
    toasts: [] as Toast[],
    sidebarOpen: false,
  }),
  actions: {
    notify(message: string, type: ToastType = 'info') {
      const id = ++counter
      this.toasts.push({ id, message, type })
      setTimeout(() => this.dismiss(id), 4000)
    },
    success(message: string) {
      this.notify(message, 'success')
    },
    error(message: string) {
      this.notify(message, 'error')
    },
    dismiss(id: number) {
      this.toasts = this.toasts.filter((t) => t.id !== id)
    },
    toggleSidebar() {
      this.sidebarOpen = !this.sidebarOpen
    },
    closeSidebar() {
      this.sidebarOpen = false
    },
  },
})
