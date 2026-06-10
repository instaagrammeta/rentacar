// Vue Router configuration with authentication & permission guards.
import { createRouter, createWebHashHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes = [
  { path: '/login', name: 'login', component: () => import('@/views/LoginView.vue'), meta: { public: true } },
  {
    path: '/',
    component: () => import('@/layouts/MainLayout.vue'),
    children: [
      { path: '', redirect: '/dashboard' },
      { path: 'dashboard', name: 'dashboard', component: () => import('@/views/DashboardView.vue'), meta: { permission: 'dashboard' } },
      { path: 'clients', name: 'clients', component: () => import('@/views/ClientsView.vue'), meta: { permission: 'clients' } },
      { path: 'clients/:id', name: 'client-detail', component: () => import('@/views/ClientDetailView.vue'), meta: { permission: 'clients' } },
      { path: 'cars', name: 'cars', component: () => import('@/views/CarsView.vue'), meta: { permission: 'cars' } },
      { path: 'reservations', name: 'reservations', component: () => import('@/views/ReservationsView.vue'), meta: { permission: 'reservations' } },
      { path: 'rentals', name: 'rentals', component: () => import('@/views/RentalsView.vue'), meta: { permission: 'rentals' } },
      { path: 'returns', name: 'returns', component: () => import('@/views/ReturnsView.vue'), meta: { permission: 'returns' } },
      { path: 'payments', name: 'payments', component: () => import('@/views/PaymentsView.vue'), meta: { permission: 'payments' } },
      { path: 'blacklist', name: 'blacklist', component: () => import('@/views/BlacklistView.vue'), meta: { permission: 'blacklist' } },
      { path: 'accidents', name: 'accidents', component: () => import('@/views/AccidentsView.vue'), meta: { permission: 'accidents' } },
      { path: 'reports', name: 'reports', component: () => import('@/views/ReportsView.vue'), meta: { permission: 'reports' } },
      { path: 'settings', name: 'settings', component: () => import('@/views/SettingsView.vue'), meta: { permission: 'settings' } },
      { path: 'users', name: 'users', component: () => import('@/views/UsersView.vue'), meta: { permission: 'users' } },
      { path: 'audit', name: 'audit', component: () => import('@/views/AuditView.vue'), meta: { permission: 'audit' } },
    ],
  },
  { path: '/:pathMatch(.*)*', redirect: '/dashboard' },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  if (to.meta.public) {
    return auth.isAuthenticated && to.name === 'login' ? { name: 'dashboard' } : true
  }
  if (!auth.isAuthenticated) {
    return { name: 'login' }
  }
  if (to.meta.permission && !auth.can(to.meta.permission)) {
    return { name: 'dashboard' }
  }
  return true
})

export default router
