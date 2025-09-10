export default defineNuxtPlugin(async (nuxtApp) => {
  const { initializeAuth, isLoggedIn } = useAuth()

  initializeAuth()

  if (import.meta.client) {
    window.addEventListener('storage', (event) => {
      if (event.key === 'isLoggedIn') {
        isLoggedIn.value = event.newValue === 'true'
        console.log('Auth state updated from storage event:', isLoggedIn.value)
      }
    })
  }
})