export default defineNuxtRouteMiddleware((to, from) => {
	if (import.meta.server) return 
	// const { isLoggedIn } = useAuth()
  	// // These paths are accessible regardless of login status
  	// const publicPaths = ['/login', '/reset-password', '/consent']

	// // If trying to access a public path, allow it
	// if (publicPaths.some(path => to.path.startsWith(path))) {
	// 	// If logged in and trying to access /login, redirect to home or dashboard
	// 	if (isLoggedIn.value && to.path === '/login') {
	// 		return navigateTo('/') // Or your main dashboard route
	// 	}
	// 	return // Allow navigation
	// }

    // console.log(`Middleware: checking isLoggedIn ${to.path}.`)
	// // If not logged in and trying to access a protected route
	// if (!isLoggedIn.value) {
	// 	// Store the intended path to redirect after login (optional)
	// 	// if (to.fullPath !== '/') {
	// 	//   sessionStorage.setItem('redirectAfterLogin', to.fullPath)
	// 	// }
	// 	return navigateTo('/login')
	// }
})