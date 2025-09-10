<script setup lang="ts">
import { ref } from "vue"

definePageMeta({
  layout: "login"
})

const auth = useAuth()
const email = ref("")
const password = ref("")
const isPwd = ref(true)

function loginClick(event: Event) {
  event.preventDefault()
	if (email && password) {
    document.getElementById("loginError")!.style.display = 'none'
    const body = {email_address: email.value, password: password.value}
    fetch("/login/sign-in", {method: "POST", body: body, headers: auth.getUnAuthHeader()})
    .then(response => {
      sessionStorage.setItem("token", response.data.token)
      auth.login("/")
    })
    return false
	}
  return false
}

function openModal() {
  var dialog = document.getElementById("forgot-modal")!
  dialog.classList.remove("hidden")
}

function resetClick() {
    fetch("/login/forgot-password/"+email.value, {method: "GET", headers: auth.getUnAuthHeader()})
    .then(() => {
      var dialog = document.getElementById("forgot-modal")!
      dialog.classList.add("hidden")
    })
}

function closeModal() {
  var dialog = document.getElementById("forgot-modal")!
  dialog.classList.add("hidden")
}

function formListener() {
  var loginForm = document.getElementById("login-form")!
  loginForm.addEventListener("submit", loginClick)
}

onMounted(() => formListener())

</script>

<template>
  <div class="title-div">
    <p>Change to img tag</p>
    <p class="h-2">Login</p>
  </div>	
  <div>
    <form id="login-form" class="group" novalidate>
      <div class="mb-5">
        <label for="email" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Email</label>
        <input type="email" id="email" v-model="email" class="bg-gray-200 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500 invalid:[&:not(:placeholder-shown):not(:focus)]:border-red-500 peer" placeholder=" " required pattern="[^@\s]+@[^@\s]+\.[^@\s]+" autocomplete="email" />
        <p class="mt-2 hidden text-sm text-red-600 dark:text-red-500 peer-[&:not(:placeholder-shown):not(:focus):invalid]:block">Email is required!</p>
      </div>
      <div class="mb-5">
        <label for="password" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Password</label>
        <input type="password" id="password" v-model="password" class="bg-gray-200 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500 invalid:[&:not(:placeholder-shown):not(:focus)]:border-red-500 peer" placeholder=" " pattern=".{8,}" required autocomplete="password" />
        <p class="mt-2 hidden text-sm text-red-600 dark:text-red-500 peer-[&:not(:placeholder-shown):not(:focus):invalid]:block">Password is required, add at least 8 charaters in length!</p>
      </div>
      <div id="submit" class="button-div">
        <button type="submit"
          class="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 me-2 mb-2 dark:bg-blue-600 dark:hover:bg-blue-700 focus:outline-none dark:focus:ring-blue-800 group-invalid:pointer-events-none group-invalid:opacity-30"
        >Login</button>
      </div>
      <div class="error-div">
        <p id="loginError" class="text-red-500" style="display: none">Email/Password combination invalid</p>
      </div>
      <div class="flex flex-col items-center mt-2">
        <button type="button" @click="openModal()" class="font-medium text-blue-600 dark:text-blue-500">Forgot Password?</button>
      </div>
    </form>
  </div>
  <div id="forgot-modal" aria-hidden="true" class="hidden fixed inset-0 transition-opacity bg-gray-200 bg-opacity-75">
    <div tabindex="-1" class="fixed z-50 top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-auto bg-white dark:bg-gray-900 rounded-md px-4 py-2 space-y-5 drop-shadow-lg">
      <div class="relative p-4 w-full max-w-2xl max-h-full">
          <!-- Modal content -->
          <div class="relative bg-white dark:bg-gray-900 rounded-lg shadow dark:bg-black">
              <!-- Modal header -->
              <div class="flex items-center justify-between p-4 md:p-5 border-b rounded-t dark:bg-gray-900 dark:border-gray-600">
                  <h3 class="text-xl font-semibold text-gray-900 dark:text-white">
                    Forgot Password?
                  </h3>
                  <button type="button" class="text-gray-400 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm w-8 h-8 ms-auto inline-flex justify-center items-center dark:hover:bg-gray-600 dark:hover:text-white" @click="closeModal()">
                      <svg class="w-3 h-3" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 14 14">
                          <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="m1 1 6 6m0 0 6 6M7 7l6-6M7 7l-6 6"/>
                      </svg>
                      <span class="sr-only">Close modal</span>
                  </button>
              </div>
              <!-- Modal body -->
              <div class="p-4 md:p-5 space-y-8 dark:bg-gray-900">
                <p class="text-sm">Enter your email, if the email address is in our system, an email with instructions will be sent.</p>
                <input
                  type="email"
                  id="email-modal"
                  v-model="email"
                  class="bg-gray-200 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500 invalid:[&:not(:placeholder-shown):not(:focus)]:border-red-500 peer"
                  placeholder="email"
                  required
                  pattern="[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$">
                </input>
                <p class="text-xs">(Please check your spam folder, just in case you don't recieve this email within a few minutes)</p>
              </div>
              <!-- Modal footer -->
              <div class="flex items-center p-4 md:p-5 border-t border-gray-200 rounded-b dark:border-gray-600 dark:bg-gray-900">
                  <button @click="resetClick()" type="button" class="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800">Submit</button>
                  <button @click="closeModal()" type="button" class="py-2.5 px-5 ms-3 text-sm font-medium text-gray-900 focus:outline-none bg-white rounded-lg border border-gray-200 hover:bg-gray-100 hover:text-blue-700 focus:z-10 focus:ring-4 focus:ring-gray-100 dark:focus:ring-gray-700 dark:bg-gray-800 dark:text-gray-400 dark:border-gray-600 dark:hover:text-white dark:hover:bg-gray-700">Close</button>
              </div>
          </div>
      </div>
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