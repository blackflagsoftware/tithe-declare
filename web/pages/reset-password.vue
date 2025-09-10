<script setup lang="ts">
import { ref } from "vue"

definePageMeta({
  layout: "login"
})

const route = useRoute()
const password = ref("")
const passwordconfirm = ref("")
const { getUnAuthHeader } = useAuth()
const fetch = apiFetch()

function resetClick() {
	if (route.query.token === "") {
    console.log("Token not set") //  TODO: change to an alert
		return false
	}
  if (password.value !== passwordconfirm.value) {
    console.log("passwords don't match") //  TODO: change to an alert
		return false
  }
  const body = {reset_token: route.query.token, email_address: route.query.email, password: password.value, confirm_password: passwordconfirm.value}
  fetch("/login/reset/pwd", {method: "POST", body: body, headers: getUnAuthHeader()})
  .then(() => {
    // navigate to login page
    navigateTo("/")
  })
  return false
}

function cancelClick() {
  navigateTo("/login")
}

const disableResetBtn = computed(() => {
    return (password.value === "" || passwordconfirm.value === "") ? true : false
})

</script>

<template>
  <div class="title-div">
    <p>Change to img tag</p>
    <p class="text-2xl">Reset Password</p>
  </div>	
  <p class="mb-4">For email: {{ route.query.email }}</p>
  <form>

  <div class="mb-5">
    <label for="password" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Password</label>
    <input
      id="password"
      type="password"
      v-model="password"
      class="bg-gray-200 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
      placeholder=" "
      pattern=".{8,}"
      required
      autocomplete="password"
    />
    <p class="mt-2 hidden text-sm text-red-600 dark:text-red-500">Password is required, add at least 8 charaters in length!</p>
  </div>
  <div class="mb-5">
    <label for="passwordconfitm" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Confirm Password</label>
    <input
      id="passwordconfirm"
      type="password"
      v-model="passwordconfirm"
      class="bg-gray-200 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
      placeholder=" "
      pattern=".{8,}"
      required
      autocomplete="password"
    />
    <p class="mt-2 hidden text-sm text-red-600 dark:text-red-500">Password is required, add at least 8 charaters in length!</p>
  </div>
  </form>
  <div id="submit" class="button-div">
		<button type="button" @click="cancelClick()" class="text-black dark:text-white bg-gray-300 hover:bg-gray-400 focus:ring-4 focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 me-2 mb-2 dark:bg-gray-600 dark:hover:bg-gray-700 focus:outline-none dark:focus:ring-blue-80">Cancel</button>
    <button
      type="submit"
      @click="resetClick()"
      class="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 mb-2 dark:bg-blue-600 dark:hover:bg-blue-700 focus:outline-none dark:focus:ring-blue-800 disabled:bg-gray-700 disabled:hover:bg-gray-700"
      :disabled="disableResetBtn"
    >Reset</button>
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
  justify-content: space-between;
  padding-top: 10px;
}

</style>