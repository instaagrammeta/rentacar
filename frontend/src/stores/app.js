// Global UI state: snackbar notifications and a loading flag.
import { defineStore } from 'pinia'

export const useAppStore = defineStore('app', {
  state: () => ({
    snackbar: { show: false, text: '', color: 'success' },
    loading: false,
  }),
  actions: {
    notify(text, color = 'success') {
      this.snackbar = { show: true, text, color }
    },
    setLoading(value) {
      this.loading = value
    },
  },
})
