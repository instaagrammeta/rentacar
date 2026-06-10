// Centralised API client. Every backend endpoint is wrapped in a typed method.
import http from './http'

export const authApi = {
  login: (payload) => http.post('/auth/login', payload),
  me: () => http.get('/auth/me'),
  roles: () => http.get('/auth/roles'),
  listUsers: () => http.get('/auth/users'),
  createUser: (payload) => http.post('/auth/users', payload),
  changePassword: (id, password) => http.post(`/auth/users/${id}/password`, { password }),
}

export const clientApi = {
  list: (params) => http.get('/clients', { params }),
  get: (id) => http.get(`/clients/${id}`),
  history: (id) => http.get(`/clients/${id}/history`),
  search: (params) => http.get('/clients/search', { params }),
  create: (payload) => http.post('/clients', payload),
  update: (id, payload) => http.put(`/clients/${id}`, payload),
  remove: (id) => http.delete(`/clients/${id}`),
}

export const carApi = {
  list: (params) => http.get('/cars', { params }),
  available: () => http.get('/cars/available'),
  get: (id) => http.get(`/cars/${id}`),
  create: (payload) => http.post('/cars', payload),
  update: (id, payload) => http.put(`/cars/${id}`, payload),
  remove: (id) => http.delete(`/cars/${id}`),
}

export const reservationApi = {
  list: (params) => http.get('/reservations', { params }),
  create: (payload) => http.post('/reservations', payload),
  setStatus: (id, status) => http.post(`/reservations/${id}/status`, { status }),
  cancel: (id) => http.post(`/reservations/${id}/cancel`),
}

export const rentalApi = {
  list: (params) => http.get('/rentals', { params }),
  get: (id) => http.get(`/rentals/${id}`),
  create: (payload) => http.post('/rentals', payload),
  cancel: (id) => http.post(`/rentals/${id}/cancel`),
  contractUrl: (id) => `/api/rentals/${id}/contract`,
  previewReturn: (id, payload) => http.post(`/rentals/${id}/return/preview`, payload),
  createReturn: (id, payload) => http.post(`/rentals/${id}/return`, payload),
}

export const paymentApi = {
  list: (params) => http.get('/payments', { params }),
  create: (payload) => http.post('/payments', payload),
  receiptUrl: (id) => `/api/payments/${id}/receipt`,
}

export const blacklistApi = {
  list: () => http.get('/blacklist'),
  reasons: () => http.get('/blacklist/reasons'),
  add: (payload) => http.post('/blacklist', payload),
  remove: (clientId) => http.delete(`/blacklist/${clientId}`),
}

export const accidentApi = {
  list: (params) => http.get('/accidents', { params }),
  create: (payload) => http.post('/accidents', payload),
  update: (id, payload) => http.put(`/accidents/${id}`, payload),
  remove: (id) => http.delete(`/accidents/${id}`),
}

export const dashboardApi = {
  summary: () => http.get('/dashboard/summary'),
  revenueByMonth: (months = 12) => http.get('/dashboard/revenue-by-month', { params: { months } }),
  topCars: (limit = 5) => http.get('/dashboard/top-cars', { params: { limit } }),
  rentalStatistics: () => http.get('/dashboard/rental-statistics'),
}

export const reportApi = {
  daily: (date) => http.get('/reports/daily-revenue', { params: { date } }),
  monthly: (year, month) => http.get('/reports/monthly-revenue', { params: { year, month } }),
  yearly: (year) => http.get('/reports/yearly-revenue', { params: { year } }),
  profitableCars: () => http.get('/reports/profitable-cars'),
  activeRentals: () => http.get('/reports/active-rentals'),
  debtors: () => http.get('/reports/debtors'),
  clientStats: () => http.get('/reports/client-statistics'),
  exportUrl: (type, params = {}) => {
    const qs = new URLSearchParams(params).toString()
    return `/api/reports/export/${type}${qs ? `?${qs}` : ''}`
  },
}

export const settingsApi = {
  get: () => http.get('/settings'),
  update: (payload) => http.put('/settings', payload),
}

export const backupApi = {
  list: () => http.get('/backups'),
  create: () => http.post('/backups'),
}

export const projectApi = {
  exportUrl: () => '/api/project/export',
  import: (formData) => http.post('/project/import', formData),
}

export const uploadApi = {
  upload: (formData) => http.post('/uploads', formData),
  fileUrl: (path) => `/api/uploads/${path}`,
}

export const auditApi = {
  list: (params) => http.get('/audit', { params }),
}
