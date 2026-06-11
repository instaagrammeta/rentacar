import type {
  Accident,
  BlacklistEntry,
  Car,
  Client,
  DashboardSummary,
  Option,
  Paginated,
  Payment,
  PublicRentalView,
  Rental,
  Reservation,
  User,
  VehicleReturn,
} from '~/types'
import { useUiStore } from '~/stores/ui'

type Query = Record<string, any>

/**
 * useApi exposes a typed client for the Go backend. The bearer token is read
 * from the `token` cookie and a 401 response clears the session.
 */
// getStoredToken reads the JWT from localStorage (client-only SPA).
function getStoredToken(): string | null {
  return import.meta.client ? localStorage.getItem('token') : null
}

// resolveApiBase makes the SPA work out-of-the-box when it is opened from a
// remote server: if the configured base still points at localhost but the page
// itself is served from another host, derive the API URL from the current host
// (same hostname, port 5000). An explicit non-localhost NUXT_PUBLIC_API_BASE
// (e.g. a real domain behind a reverse proxy) always wins.
function resolveApiBase(configured: string): string {
  if (import.meta.client) {
    const isLocalConfig =
      !configured || configured.includes('localhost') || configured.includes('127.0.0.1')
    const host = window.location.hostname
    const onLocalhost = host === 'localhost' || host === '127.0.0.1'
    if (isLocalConfig && !onLocalhost) {
      return `${window.location.protocol}//${host}:5000/api`
    }
  }
  return configured || 'http://localhost:5000/api'
}

export function useApi() {
  const config = useRuntimeConfig()
  const base = resolveApiBase(config.public.apiBase as string)
  const ui = useUiStore()

  async function request<T>(path: string, opts: Record<string, any> = {}): Promise<T> {
    const token = getStoredToken()
    return await $fetch<T>(path, {
      baseURL: base,
      headers: token ? { Authorization: `Bearer ${token}` } : {},
      onResponseError({ response }) {
        const status = response.status
        const message = (response._data && response._data.error) || 'Произошла ошибка'
        if (status === 401) {
          if (path.includes('/auth/login')) {
            ui.error(message || 'Неверное имя пользователя или пароль')
          } else {
            if (import.meta.client) localStorage.removeItem('token')
            ui.error('Сессия истекла. Войдите снова.')
            if (import.meta.client && !window.location.pathname.startsWith('/login')) {
              window.location.href = '/login'
            }
          }
        } else if (status === 403) {
          ui.error('Недостаточно прав для выполнения операции')
        } else if (status >= 400) {
          ui.error(message)
        }
      },
      onRequestError() {
        ui.error('Не удалось подключиться к серверу. Проверьте, что бэкенд запущен.')
      },
      ...opts,
    })
  }

  async function download(path: string, filename: string, query?: Query) {
    const blob = await request<Blob>(path, { method: 'GET', query, responseType: 'blob' })
    const url = URL.createObjectURL(blob as Blob)
    const a = document.createElement('a')
    a.href = url
    a.download = filename
    document.body.appendChild(a)
    a.click()
    a.remove()
    URL.revokeObjectURL(url)
  }

  return {
    raw: request,
    download,
    fileUrl: (p: string) => `${base}/uploads/${p}`,

    auth: {
      login: (body: { username: string; password: string }) =>
        request<{ access_token: string; refresh_token: string; user: User }>('/auth/login', {
          method: 'POST',
          body,
        }),
      me: () => request<User>('/auth/me'),
      logout: () => request('/auth/logout', { method: 'POST' }),
      roles: () => request<Option[]>('/auth/roles'),
      listUsers: () => request<User[]>('/auth/users'),
      createUser: (body: Query) => request<User>('/auth/users', { method: 'POST', body }),
      changePassword: (id: number, password: string) =>
        request(`/auth/users/${id}/password`, { method: 'POST', body: { password } }),
    },

    clients: {
      list: (query?: Query) => request<Paginated<Client>>('/clients', { query }),
      get: (id: number) => request<Client>(`/clients/${id}`),
      history: (id: number) => request<any>(`/clients/${id}/history`),
      search: (query: Query) => request<Client | Client[]>('/clients/search', { query }),
      create: (body: Query) => request<Client>('/clients', { method: 'POST', body }),
      update: (id: number, body: Query) => request<Client>(`/clients/${id}`, { method: 'PUT', body }),
      remove: (id: number) => request(`/clients/${id}`, { method: 'DELETE' }),
    },

    cars: {
      list: (query?: Query) => request<Paginated<Car>>('/cars', { query }),
      available: () => request<Car[]>('/cars-available'),
      get: (id: number) => request<Car>(`/cars/${id}`),
      create: (body: Query) => request<Car>('/cars', { method: 'POST', body }),
      update: (id: number, body: Query) => request<Car>(`/cars/${id}`, { method: 'PUT', body }),
      remove: (id: number) => request(`/cars/${id}`, { method: 'DELETE' }),
    },

    reservations: {
      list: (query?: Query) => request<Paginated<Reservation>>('/reservations', { query }),
      get: (id: number) => request<Reservation>(`/reservations/${id}`),
      create: (body: Query) => request<Reservation>('/reservations', { method: 'POST', body }),
      setStatus: (id: number, status: string) =>
        request<Reservation>(`/reservations/${id}/status`, { method: 'POST', body: { status } }),
      cancel: (id: number) => request<Reservation>(`/reservations/${id}/cancel`, { method: 'POST' }),
    },

    rentals: {
      list: (query?: Query) => request<Paginated<Rental>>('/rentals', { query }),
      get: (id: number) => request<Rental>(`/rentals/${id}`),
      create: (body: Query) => request<Rental>('/rentals', { method: 'POST', body }),
      cancel: (id: number) => request<Rental>(`/rentals/${id}/cancel`, { method: 'POST' }),
      downloadContract: (id: number, contract: string) =>
        download(`/rentals/${id}/contract`, `${contract}.pdf`),
      previewReturn: (id: number, body: Query) =>
        request<VehicleReturn>(`/rentals/${id}/return/preview`, { method: 'POST', body }),
      createReturn: (id: number, body: Query) =>
        request<VehicleReturn>(`/rentals/${id}/return`, { method: 'POST', body }),
      qr: (id: number) =>
        request<{ qr_code_path: string | null; public_token: string | null; public_url: string }>(
          `/rentals/${id}/qr`,
        ),
    },

    // Public, unauthenticated endpoints (rental status page reached via QR).
    public: {
      rental: (token: string) => request<PublicRentalView>(`/public/rentals/${token}`),
    },

    payments: {
      list: (query?: Query) => request<Paginated<Payment>>('/payments', { query }),
      get: (id: number) => request<Payment>(`/payments/${id}`),
      create: (body: Query) => request<Payment>('/payments', { method: 'POST', body }),
      downloadReceipt: (id: number, receipt: string) =>
        download(`/payments/${id}/receipt`, `${receipt}.pdf`),
    },

    blacklist: {
      list: () => request<BlacklistEntry[]>('/blacklist'),
      reasons: () => request<Option[]>('/blacklist/reasons'),
      add: (body: Query) => request<BlacklistEntry>('/blacklist', { method: 'POST', body }),
      remove: (clientId: number) => request(`/blacklist/${clientId}`, { method: 'DELETE' }),
    },

    accidents: {
      list: (query?: Query) => request<Paginated<Accident>>('/accidents', { query }),
      get: (id: number) => request<Accident>(`/accidents/${id}`),
      create: (body: Query) => request<Accident>('/accidents', { method: 'POST', body }),
      update: (id: number, body: Query) => request<Accident>(`/accidents/${id}`, { method: 'PUT', body }),
      remove: (id: number) => request(`/accidents/${id}`, { method: 'DELETE' }),
    },

    dashboard: {
      summary: () => request<DashboardSummary>('/dashboard/summary'),
      revenueByMonth: (months = 12) =>
        request<{ month: string; revenue: number }[]>('/dashboard/revenue-by-month', { query: { months } }),
      topCars: (limit = 5) =>
        request<{ car: string; rentals: number }[]>('/dashboard/top-cars', { query: { limit } }),
      rentalStatistics: () => request<Record<string, number>>('/dashboard/rental-statistics'),
    },

    reports: {
      daily: (date?: string) => request<any>('/reports/daily-revenue', { query: { date } }),
      monthly: (year: number, month: number) =>
        request<any>('/reports/monthly-revenue', { query: { year, month } }),
      yearly: (year: number) => request<any>('/reports/yearly-revenue', { query: { year } }),
      profitableCars: () => request<any[]>('/reports/profitable-cars'),
      activeRentals: () => request<Rental[]>('/reports/active-rentals'),
      debtors: () => request<any[]>('/reports/debtors'),
      clientStats: () => request<any[]>('/reports/client-statistics'),
      exportReport: (type: string, query: Query = {}) =>
        download(`/reports/export/${type}`, `${type}.xlsx`, query),
    },

    settings: {
      get: () => request<any>('/settings'),
      update: (body: Query) => request<any>('/settings', { method: 'PUT', body }),
    },

    backups: {
      list: () => request<any[]>('/backups'),
      create: () => request<any>('/backups', { method: 'POST' }),
    },

    project: {
      exportProject: () => download('/project/export', 'project.rentacar'),
      import: (form: FormData) => request<any>('/project/import', { method: 'POST', body: form }),
    },

    uploads: {
      upload: (form: FormData) => request<{ paths: string[] }>('/uploads', { method: 'POST', body: form }),
    },

    audit: {
      list: (query?: Query) => request<any[]>('/audit', { query }),
    },

    sms: {
      sendToClient: (clientId: number, message: string) =>
        request(`/clients/${clientId}/sms`, { method: 'POST', body: { message } }),
    },
  }
}
