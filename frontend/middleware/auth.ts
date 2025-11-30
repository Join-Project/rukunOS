export default defineNuxtRouteMiddleware((to, from) => {
  const authStore = useAuthStore()

  // If trying to access protected route without auth, redirect to login
  if (!authStore.isAuthenticated && to.path !== '/login' && to.path !== '/register' && to.path !== '/') {
    return navigateTo('/login')
  }

  // If already authenticated and trying to access login/register, redirect to dashboard
  if (authStore.isAuthenticated && (to.path === '/login' || to.path === '/register')) {
    return navigateTo('/dashboard')
  }
})
