<script setup lang="ts">

definePageMeta({
  layout: "login"
})

const route = useRoute()
const clientId = ref<string>(route.query.client_id ? route.query.client_id.toString() : "")
const scopes = ref<string>(route.query.scopes ? route.query.scopes.toString() : "")
const state = ref<string>(route.query.state ? route.query.state.toString() : "")
const redirectUri = ref<string>(route.query.redirect_uri ? route.query.redirect_uri.toString() : "")
const codeChallenge = ref<string>(route.query.code_challenge ? route.query.code_challenge.toString() : "")
const codeChallengeMethod = ref<string>(route.query.code_challenge_method ? route.query.code_challenge_method.toString() : "")

function pre() {
	const token = sessionStorage.getItem("token")
	console.log("out in the open")
	if (token !== "") {
		console.log("found token", token)
  	// const config = useRuntimeConfig()
    // const options = {headers: {"Authorization": "bearer "+token}}
		// 	console.log("config.ENV", config.public.apiURL)
    // axios.post(config.public.apiURL + "/login/verify", null, options)
		// .then(() => {
		// 	authShow.value = false
		// 	consentShow.value = true
		// })
    // .catch(error => {
		// 	console.log(error)
    // })
	} else {
		console.log("no token, go to login")
	}
}

onMounted(() => {
	pre()
})
</script>

<template>
	<p>client id: {{clientId}}</p>
</template>