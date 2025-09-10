<script setup lang="ts">
import { onMounted, ref } from "vue"

const logins = ref<Login[]>(new Array<Login>())
const roles = ref<Role[]>(new Array<Role>())
const currentLogin = ref<Login>(BlankLogin(roles.value))
const currentRole = ref<Role>(BlankRole())
const modeAdd = ref<boolean>(true)
const deleteLogin = ref<Login>(BlankLogin(roles.value))
const { getAuthHeader } = useAuth()
const fetch = apiFetch()

function loadLogins() {
	fetch("/login/roles", {method: "GET", headers: getAuthHeader()})
	.then(response => {
		logins.value = response.data.data
	})
}

function loadRoles() {
	fetch("/role/search", {method: "POST", body: {}, headers: getAuthHeader()})
	.then(response => {
		roles.value = response.data
	})
}

/*
login logic
*/
// 'Add Login' button click
function addLoginClick() {
	currentLogin.value = BlankLogin(roles.value)
	modeAdd.value = true
	modalHideShow(true, "login-modal")
}

// 'Edit' link click
function editLoginClick(id: string) {
	currentLogin.value = logins.value?.find((l: Login) => l.id === id)!
	currentLogin.value.roles_selected = translateSelectedRoles(currentLogin.value.roles)
	modeAdd.value = false
	modalHideShow(true, "login-modal")
}

// comes from the 'emit save' click within the modal
function saveLoginModalClick(login: Login) {
	// depending on the mode: post / path with the currentLogin
	const roleIds = roles.value.reduce((roleIds: string[], r: Role, index: number) => {
		if (login.roles_selected[index]) {
			roleIds.push(r.id)
		}
		return roleIds
	}, [])

	const body = {...login}
	if (modeAdd.value) {
			fetch("/login", {method: "POST", body: body, headers: getAuthHeader()})
			.then(response => {
				updateLoginRoles(response.data.id, roleIds)
			})
		} else {
			fetch("/login", {method: "PATCH", body: body, headers: getAuthHeader()})
			.then(() => {
				updateLoginRoles(login.id, roleIds)
			})
		}
	// reload the login(s)
	modalHideShow(false, "login-modal")
}

// comes from the 'emit cancel' click within the modal
function cancelLoginModalClick() {
	currentLogin.value = BlankLogin(roles.value)
	modalHideShow(false, "login-modal")
}

function updateLoginRoles(loginId: string, roleIds: string[]) {
	const body = {login_id: loginId, role_ids: roleIds}
	fetch("/login-role", {method: "POST", body: body, headers: getAuthHeader()})
	.then(() => {
		loadLogins()
	})
}

// 'Delete' link click
function deleteLoginClick(loginId: string) {
	deleteLogin.value = logins.value.find((l: Login) => l.id === loginId)!
	if (!deleteLogin.value) {
		console.log("Unable to find login record, something went wrong") // TODO: alert
		return
	}
	modalHideShow(true, "login-delete-modal")
}

// comes from the 'emit cancel delete' click within the modal
function deleteLoginModalClick() {
	const loginId = deleteLogin.value.id
	if (!loginId || loginId === "") {
		console.log("Login id is not set, something went wrong") // TODO: alert
		return
	}
	fetch("/login/" + loginId, {method: "DELETE", headers: getAuthHeader()})
	.then(() => {
		deleteLogin.value = BlankLogin(roles.value)
		loadLogins()
	})
	modalHideShow(false, "login-delete-modal")
}

// comes from the 'emit cancel delete' click within the modal
function cancelDeleteLoginModalClick() {
	deleteLogin.value = BlankLogin(roles.value)
	modalHideShow(false, "login-delete-modal")
}

/*
role logic
*/
// 'Add Login' button click
function addRoleClick() {
	currentRole.value = BlankRole()
	modeAdd.value = true
	modalHideShow(true, "role-modal")
}

// 'Edit' link click
function editRoleClick(id: string) {
	currentRole.value = roles.value?.find((r: Role) => r.id === id)!
	modeAdd.value = false
	modalHideShow(true, "role-modal")
}

// comes from the 'emit save' click within the modal
function saveRoleClick(role: Role) {
	const body = {...role}
	let promise: Promise<any>
	if (modeAdd.value) {
		promise = fetch("/role", {method: "POST", body: body, headers: getAuthHeader()})
		.then(() => {
			loadRoles()
		})
	} else {
		promise = fetch("/role", {method: "PATCH", body: body, headers: getAuthHeader()})
		.then(() => {
			loadLogins()
			loadRoles()
		})
	}
	modalHideShow(false, "role-modal")
}

function cancelRoleClick() {
	currentRole.value = BlankRole()
	modalHideShow(false, "role-modal")
}

// 'Delete' link click
function deleteRoleClick(roleId: string) {
	currentRole.value = roles.value?.find((r: Role) => r.id === roleId)!
	if (!currentRole.value) {
		console.log("Unable to find role record, something went wrong") // TODO: alert
		return
	}
	modalHideShow(true, "role-delete-modal")
}

// comes from the 'emit cancel delete' click within the modal
function deleteRoleModalClick() {
	const roleId = currentRole.value.id
	if (!roleId || roleId === "") {
		console.log("Role id is not set, something went wrong") // TODO: alert
		return
	}
	fetch("/role/" + roleId, {method: "DELETE", headers: getAuthHeader()})
	.then(() => {
		currentRole.value = BlankRole()
		loadLogins()
		loadRoles()
	})
	modalHideShow(false, "role-delete-modal")
}

// comes from the 'emit cancel delete' click within the modal
function cancelDeleteRoleModalClick() {
	currentRole.value = BlankRole()
	modalHideShow(false, "role-delete-modal")
}

function translateSelectedRoles(loginRoles: string[]): boolean[] {
	return roles.value.map(role => {
		const checked = loginRoles.includes(role.name)
		return checked
	})
}

onMounted(() => {
	loadLogins()
	loadRoles()
})

</script>

<template>
	<p class="text-4xl">Admin</p>
	<p class="mt-8 mb-2 text-2xl">Logins</p>
	<button @click="addLoginClick()" class="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 me-2 mb-2 dark:bg-blue-600 dark:hover:bg-blue-700 focus:outline-none dark:focus:ring-blue-800">Add Login</button>
	<div class="relative overflow-x-auto shadow-md sm:rounded-lg">
    <table class="w-full text-sm text-left rtl:text-right text-gray-600 dark:text-gray-400">
        <thead class="text-xs text-gray-700 uppercase bg-gray-400 dark:bg-gray-600 dark:text-gray-200">
            <tr>
							<th scope="col" class="px-6 py-3">
								Email
							</th>
							<th scope="col" class="px-6 py-3">
								Roles	
							</th>
							<th scope="col" class="px-6 py-3">
								Action
							</th>
            </tr>
        </thead>
        <tbody>
            <tr class="odd:bg-white odd:dark:bg-gray-900 even:bg-gray-200 even:dark:bg-gray-800 border-b dark:border-gray-700" v-for="login in logins">
							<th scope="row" class="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white">
								{{ login.email_address }}
							</th>
							<td class="px-6 py-4">
								{{ login.roles }}
							</td>
							<td class="px-6 py-4">
									<a href="#" @click="editLoginClick(login.id)" class="font-medium text-blue-600 dark:text-blue-500 hover:underline">Edit</a>
									|
									<a href="#" @click="deleteLoginClick(login.id)" class="font-medium text-blue-600 dark:text-blue-500 hover:underline">Delete</a>
							</td>
            </tr>
        </tbody>
    </table>
	</div>
	<p class="mt-8 mb-2 text-2xl">Roles</p>
	<button @click="addRoleClick()" class="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 me-2 mb-2 dark:bg-blue-600 dark:hover:bg-blue-700 focus:outline-none dark:focus:ring-blue-800">Add Role</button>
	<div class="relative overflow-x-auto shadow-md sm:rounded-lg">
    <table class="w-full text-sm text-left rtl:text-right text-gray-600 dark:text-gray-400">
        <thead class="text-xs text-gray-700 uppercase bg-gray-400 dark:bg-gray-500 dark:text-gray-200">
            <tr>
							<th scope="col" class="px-6 py-3">
								Name
							</th>
							<th scope="col" class="px-6 py-3">
								Description
							</th>
							<th scope="col" class="px-6 py-3">
								Action
							</th>
            </tr>
        </thead>
        <tbody>
            <tr class="odd:bg-white odd:dark:bg-gray-900 even:bg-gray-200 even:dark:bg-gray-800 border-b dark:border-gray-700" v-for="role in roles">
							<th scope="row" class="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white">
								{{ role.name }}
							</th>
							<th class="px-6 py-4">
								{{ role.description }}
							</th>
							<td class="px-6 py-4">
									<a href="#" @click="editRoleClick(role.id)" class="font-medium text-blue-600 dark:text-blue-500 hover:underline">Edit</a>
									|
									<a href="#" @click="deleteRoleClick(role.id)" class="font-medium text-blue-600 dark:text-blue-500 hover:underline">Delete</a>
							</td>
            </tr>
        </tbody>
    </table>
	</div>
	<LoginModal id="login-modal" class="hidden" :login="currentLogin" :roles="roles" @clickSubmit="saveLoginModalClick" @clickCancel="cancelLoginModalClick" :titleText="(modeAdd ? 'Add' : 'Edit') + ' Login'" />
	<DeleteModal id="login-delete-modal" class="hidden" @clickSubmit="deleteLoginModalClick" @clickCancel="cancelDeleteLoginModalClick" recordType="Login" :name="deleteLogin.email_address" />
	<RoleModal id="role-modal" class="hidden" :role="currentRole!" @clickSubmit="saveRoleClick" @clickCancel="cancelRoleClick" :titleText="(modeAdd ? 'Add' : 'Edit') + ' Role'" />
	<DeleteModal id="role-delete-modal" class="hidden" @clickSubmit="deleteRoleModalClick" @clickCancel="cancelDeleteRoleModalClick" recordType="Role" :name="currentRole.name" />

</template>