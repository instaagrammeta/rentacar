// Axios instance with JWT injection and centralised error handling.
import axios from 'axios'
import { useAuthStore } from '@/stores/auth'
import { useAppStore } from '@/stores/app'

const http = axios.create({
  baseURL: '/api',
  timeout: 30000,
})

// Attach the bearer token to every request.
http.interceptors.request.use((config) => {
  const auth = useAuthStore()
  if (auth.token) {
    config.headers.Authorization = `Bearer ${auth.token}`
  }
  return config
})

// Handle auth errors and surface server messages to the user.
http.interceptors.response.use(
  (response) => response,
  (error) => {
    const app = useAppStore()
    const status = error.response?.status
    const message = error.response?.data?.error || 'Произошла ошибка'

    if (status === 401) {
      const auth = useAuthStore()
      // Avoid redirect loops on the login request itself.
      if (!error.config?.url?.includes('/auth/login')) {
        auth.logout()
        app.notify('Сессия истекла. Войдите снова.', 'warning')
      }
    } else if (status === 403) {
      app.notify('Недостаточно прав для выполнения операции', 'error')
    } else {
      app.notify(message, 'error')
    }
    return Promise.reject(error)
  },
)

export default http
