<script setup lang="ts">

const { checkPermission, logout } = useAuth()

// function logout() {
// 	console.log("logging out")
// 	sessionStorage.removeItem("token")
// 	const loggedIn = useLoggedIn()
// 	loggedIn.value = false
// }

const isUser = computed(() => {
	return checkPermission("user")
})

const isAdmin = computed(() => {
	return true // checkPermission("admin")
})

</script>

<template>
	<div class="relative flex flex-row -top-0 -left-0 w-full h-full overflow-hidden bg-white dark:bg-gray-900">
		<div class="nav-w float-left bg-gray-400 h-full">
			<div class="flex flex-col items-center justify-center mt-10 mb-10">
				<img src="/Christ-Church-Symbol.png" alt="Logo" class="w-192 h-108 mt-2 mb-2"/>
			</div>
			<div class="flex flex-col items-start pl-10">
				<ul class="list-none pl-5 m-0 text-black">
					<!-- <li v-if="isUser || isAdmin"><NuxtLink to="/">Default</NuxtLink></li>
					<li v-if="isUser || isAdmin"><NuxtLink to="/profile">Profile</NuxtLink></li>
					<li v-if="isAdmin"><NuxtLink to="/register_routes">Register Routes</NuxtLink></li>
					<li v-if="isAdmin"><NuxtLink to="/admin">Admin</NuxtLink></li>
					<li v-if="isAdmin"><NuxtLink to="/auth">Authentication</NuxtLink></li>
					<li><NuxtLink @click.native="logout()" to="/login">Logout</NuxtLink></li> -->
				</ul>
			</div>
		</div>
		<div class="w-full p-5 bg-white dark:bg-gray-700 overflow-y-auto" id="main">
			<slot />
		</div>
	</div>
</template>