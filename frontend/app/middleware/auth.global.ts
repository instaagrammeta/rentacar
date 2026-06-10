import { useAuthStore } from '~/stores/auth'

// Global route guard: protects every page except the login screen.
export default defineNuxtRouteMiddleware(async (to) => {
  const auth = useAuthStore()

  if (to.path === '/login') {
    if (auth.isAuthenticated) return navigateTo('/')
    return
  }

  if (!auth.isAuthenticated) {
    return navigateTo('/login')
  }

  if (!auth.user) {
    await auth.fetchMe()
    if (!auth.isAuthenticated) return navigateTo('/login')
  }
})
