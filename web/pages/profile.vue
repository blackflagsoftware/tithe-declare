<script setup lang="ts">
import { onMounted, ref } from "vue"

const roles = ref<Role[]>(new Array<Role>())
const login = defineModel({ default: BlankLogin(new Array<Role>()) })
const loginId = ref<string>("")
const currentPwd = ref<string>("")
const newPwdOne = ref<string>("")
const newPwdTwo = ref<string>("")
const { getAuthHeader } = useAuth()
const fetch = apiFetch()

function loadLogin() {
  const token = sessionStorage.getItem("token")
  if (token?.length === 0) {
    console.log("error in getting token")
    return
  }
  const tokenParts = token?.split(".")!
  if (tokenParts?.length < 2) {
    console.log("error misformatted token")
    return
  }
  const claims = JSON.parse(atob(tokenParts[1]))
  loginId.value = claims.jti
  if (loginId.value.length === 0) {
    console.log("error unable to get login id")
    return
  }
    fetch("/login/" + loginId.value, {method: "GET", headers: getAuthHeader()})
    .then(response => {
      login.value = response.data.data
    })
}

function saveNameClick() {
  const body = { id: loginId.value, first_name: login.value.first_name, last_name: login.value.last_name }
  fetch("/login", {method: "PATCH", body: body, headers: getAuthHeader()})
}

function savePwdClick() {
  const body = { id: loginId.value, password: newPwdOne.value, confirm_password: newPwdTwo.value }
  fetch("/login/pwd", {method: "PATCH", body: body, headers: getAuthHeader()})
    .then(() => {
      newPwdOne.value = ""
      newPwdTwo.value = ""
    })
}

onMounted(() => {
  loadLogin()
})
</script>

<template>
  <p class="text-4xl mb-8">Profile</p>
  <div>
    <div class="mb-8">
      <div class="mb-4">
        <label for="firstName" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">First Name</label>
        <input id="firstName" type="text" v-model="login.first_name" class="bg-gray-200 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500 w-96" />
      </div>
      <div class="mb-4">
        <label for="lastName" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Last Name</label>
        <input id="lastName" type="text" v-model="login.last_name" class="bg-gray-200 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500 w-96" />
      </div>
      <div class="mb-4">
        <label for="email" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Email</label>
        <input id="email" type="text" v-model="login.email_address" class="bg-gray-200 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500 w-96" />
      </div>
      <div class="flex items-center">
        <button @click="saveNameClick()" type="button" class="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800">Save</button>
      </div>
    </div>
    <div class="mb-8">
      <div class="mb-4">
        <label for="newPwdOne" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">New Password</label>
        <input id="newPwdOne" type="password" v-model="newPwdOne" class="bg-gray-200 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500 w-96" />
      </div>
      <div class="mb-4">
        <label for="newPwdTwo" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">New Password Confirm</label>
        <input id="newPwdTwo" type="password" v-model="newPwdTwo" class="bg-gray-200 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500 w-96" />
      </div>
      <div class="flex items-center">
        <button @click="savePwdClick()" type="button" class="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800">Change</button>
      </div>
    </div>
  </div>
</template>
