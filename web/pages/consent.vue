<script setup lang="ts">
import { ref } from "vue"

definePageMeta({
  layout: "login"
})

const route = useRoute()
const clientId = ref<string>(route.query.client_id ? route.query.client_id.toString() : "")
const clientName = ref<string>(route.query.client_name ? route.query.client_name.toString() : "")
const scopes = ref<string>(route.query.scopes ? route.query.scopes.toString() : "")
const state = ref<string>(route.query.state ? route.query.state.toString() : "")
const redirectUri = ref<string>(route.query.redirect_uri ? route.query.redirect_uri.toString() : "")
const codeChallenge = ref<string>(route.query.code_challenge ? route.query.code_challenge.toString() : "")
const codeChallengeMethod = ref<string>(route.query.code_challenge_method ? route.query.code_challenge_method.toString() : "")
const consentShow = ref<boolean>(true)
const authShow = ref<boolean>(false)
const email = ref<string>("")
const password = ref<string>("")
const authCode = ref<string>("")
const { getUnAuthHeader } = useAuthHeader()
const fetch = apiFetch()

function authCancelClick() {
	consentShow.value = false
	authShow.value = true
}

function authOkClick(event: Event) {
  	event.preventDefault()
	if (email && password) {
		document.getElementById("loginError")!.style.display = 'none'
		const body = {email_address: email.value, password: password.value, client_id: clientId.value, redirect_uri: redirectUri.value, code_challenge: codeChallenge.value, code_challenge_method: codeChallengeMethod.value}
		fetch("/auth/oauth2/sign-in", {method: "POST", body: body, headers: getUnAuthHeader()})
			.then(resp => {
				authCode.value = resp.data.code
				authShow.value = false
				consentShow.value = true
			})
    	return false
	}
  	return false
}

function consentCancelClick() {
	const redirectPath = redirectUri.value + "?" + "error=user_canceled_consent&state=" + state.value
	window.location.href = redirectPath
}

function consentAllowClick() {
	const redirectPath = redirectUri.value + "?" + "code=" + authCode.value + "&state=" + state.value
	window.location.href = redirectPath
}


function formListener() {
  var loginForm = document.getElementById("login-form")!
  loginForm.addEventListener("submit", authOkClick)
}

onMounted(() => {
	const token = sessionStorage.getItem("token")
	if (token !== "") {
		fetch("/login/verify", {method: "POST", headers: getUnAuthHeader()})
		.then(() => {
			authShow.value = false
			consentShow.value = true
		})
	} 
	formListener()
})

</script>

<template>
	<div :hidden="consentShow">
		<div class="flex flex-col items-center">
			<p class="text-xl text-blue-600 dark:text-blue-500">{{ clientName }}</p>
			<p class="text-xl mb-4">wants to access your Account</p>
			<p class="text-base">This will allow <span class="text-xl text-blue-600 dark:text-blue-500">{{ clientName }}</span> to have permission to do:</p>
			<p class="text-base">{{ scopes }}</p>
		</div>
		<div class="flex flex-row justify-between mt-8">
			<button type="button" @click="consentCancelClick()" class="text-black dark:text-white bg-gray-300 hover:bg-gray-400 focus:ring-4 focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 me-2 mb-2 dark:bg-gray-600 dark:hover:bg-gray-700 focus:outline-none dark:focus:ring-blue-80">Cancel</button>
			<button type="submit" @click="consentAllowClick()" class="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 me-2 mb-2 dark:bg-blue-600 dark:hover:bg-blue-700 focus:outline-none dark:focus:ring-blue-800 group-invalid:pointer-events-none group-invalid:opacity-30">Allow</button>
		</div>
	</div>
	<div :hidden="authShow">
		<div class="title-div">
			<p>Change to img tag</p>
			<p class="h-2">Login</p>
		</div>	
		<div>
			<form id="login-form" class="group" novalidate>
				<div class="mb-5">
					<label for="email" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Email</label>
					<input
						type="email"
						id="email"
						v-model="email"
						class="bg-gray-200 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500 invalid:[&:not(:placeholder-shown):not(:focus)]:border-red-500 peer"
						placeholder=" "
						required
						pattern="[^@\s]+@[^@\s]+\.[^@\s]+"
						autocomplete="email"
					/>
					<p class="mt-2 hidden text-sm text-red-600 dark:text-red-500 peer-[&:not(:placeholder-shown):not(:focus):invalid]:block">Email is required!</p>
				</div>
				<div class="mb-5">
					<label for="password" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Password</label>
					<input
						type="password"
						id="password"
						v-model="password"
						class="bg-gray-200 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500 invalid:[&:not(:placeholder-shown):not(:focus)]:border-red-500 peer"
						placeholder=" "
						pattern=".{8,}"
						required
						autocomplete="password"
					/>
					<p class="mt-2 hidden text-sm text-red-600 dark:text-red-500 peer-[&:not(:placeholder-shown):not(:focus):invalid]:block">Password is required, add at least 8 charaters in length!</p>
				</div>
				<div class="flex flex-row justify-between mt-8">
					<button type="button" @click="authCancelClick()" class="text-black dark:text-white bg-gray-300 hover:bg-gray-400 focus:ring-4 focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 me-2 mb-2 dark:bg-gray-600 dark:hover:bg-gray-700 focus:outline-none dark:focus:ring-blue-80">Cancel</button>
					<button type="submit"
						class="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 me-2 mb-2 dark:bg-blue-600 dark:hover:bg-blue-700 focus:outline-none dark:focus:ring-blue-800 group-invalid:pointer-events-none group-invalid:opacity-30"
					>Login</button>
				</div>
				<div class="error-div">
					<p id="loginError" class="text-red-500" style="display: none">Email/Password combination invalid</p>
				</div>
			</form>
		</div>
	</div>
</template>

<style>
.title-div {
	display: flex;
  flex-direction: column;
  align-items: center;
  padding-bottom: 10px;
}

.button-div {
  display: flex;
  flex-direction: row;
  justify-content: center;
  padding-top: 10px;
}

.error-div {
  display: flex;
  flex-direction: row;
  justify-content: center;
  padding-top: 10px;
}

</style>