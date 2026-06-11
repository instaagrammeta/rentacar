import { useAuthStore } from '~/stores/auth'

// Global route guard: protects every page except the login screen and the
// public rental status pages (e.g. /r/<token> opened via a scanned QR code).
export default defineNuxtRouteMiddleware(async (to) => {
  const auth = useAuthStore()

  // Public rental pages are always accessible without a session.
  if (to.path.startsWith('/r/')) return

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
