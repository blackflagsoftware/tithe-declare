<script setup lang="ts">
import { onMounted, ref } from "vue"

const registerRoutes = ref<RegisterRoute[]>(new Array<RegisterRoute>())
const roles = ref<Role[]>(new Array<Role>())
const currentRegisterRoute = ref<RegisterRoute>(BlankRegisterRoute(roles.value))
const currentRole = ref<Role>(BlankRole())
const registerRouteTerm = ref<string>("")
const isLoading = ref<boolean>(false)
const modeAdd = ref<boolean>(true)
const bulkUpdateRouteArray = ref<String[]>([])
const deleteRegisterRoute = ref<RegisterRoute>(BlankRegisterRoute(roles.value))
const { getAuthHeader } = useAuth()
const fetch = apiFetch()

function loadRegisterRoutes() {
	fetch("/register-route/search", {method: "POST", headers: getAuthHeader()})
	.then(response => {
		registerRoutes.value = response.data.data
	})
}

function loadRoles() {
	fetch("/role/search", {method: "POST", body: {}, headers: getAuthHeader()})
	.then(response => {
		roles.value = response.data.data
	})
}

/*
registerRoute logic
*/
// 'Edit' link click
function editRegisterRouteClick(rawPath: string) {
	currentRegisterRoute.value = registerRoutes.value?.find((r: RegisterRoute) => r.raw_path === rawPath)!
	currentRegisterRoute.value.roles_selected = translateSelectedRoles(currentRegisterRoute.value.roles)
	modeAdd.value = false
	modalHideShow(true, "register-route-modal")
}

// comes from the 'emit save' click within the modal
function saveRegisterRouteModalClick(registerRoute: RegisterRoute) {
	// depending on the mode: post / path with the currentRegisterRoute
	const roleIds = roles.value.reduce((roleIds: string[], r: Role, index: number) => {
		if (registerRoute.roles_selected[index]) {
			roleIds.push(r.name)
		}
		return roleIds
	}, [])

	const body = {...registerRoute, roles: roleIds}
	fetch("/register-route", {method: "PATCH", body: body, headers: getAuthHeader()})
	loadRegisterRoutes()
	modalHideShow(false, "register-route-modal")
}

// comes from the 'emit cancel' click within the modal
function cancelRegisterRouteModalClick() {
	currentRegisterRoute.value = BlankRegisterRoute(roles.value)
	modalHideShow(false, "register-route-modal")
}

function updateRegisterRouteRoles(rawPath: string, roleIds: string[]) {
	const body = {raw_path: rawPath, roles: roleIds}
	fetch("/register-route-role", {method: "PATCH", body: body, headers: getAuthHeader()})
	.then(() => {
		loadRegisterRoutes()
	})
}

// comes from the 'emit cancel delete' click within the modal
function cancelDeleteRegisterRouteModalClick() {
	deleteRegisterRoute.value = BlankRegisterRoute(roles.value)
	modalHideShow(false, "register-route-delete-modal")
}

function translateSelectedRoles(loginRoles: Role[]): boolean[] {
	return roles.value.map(role => {
		const checked = loginRoles.includes(role)
		return checked
	})
}

function debounce<A extends any[], R>(fn: (...args: A) => R, delay: number): (...args: A) => void {
  let timeoutId: NodeJS.Timeout | number | undefined = undefined
  return (...args: A) => {
    clearTimeout(timeoutId)
	if (timeoutId) {
	  clearTimeout(timeoutId)
	}
    timeoutId = setTimeout(() => fn(...args), delay)
  }
}

const handleSearch = async () => {
	if (registerRouteTerm.value.length < 2) {
		loadRegisterRoutes()
		bulkUpdateRouteArray.value = []
		const element = document.getElementById("bulkCheckboxHeader");
		if (element instanceof HTMLInputElement) {
			element.checked = false;
		} else {
			console.error("The element 'bulkCheckboxHeader' is not a valid input element.");
		}
		return
	}
	isLoading.value = true
	try {
		const body = { search: {filters: [{column: "raw_path", compare: "LIKE", value: registerRouteTerm.value }]}}
		const response = await fetch("/register-route/search", {method: "POST", body: body, headers: getAuthHeader()})
		registerRoutes.value = response.data.data
		bulkUpdateRouteArray.value = []
		const element = document.getElementById("bulkCheckboxHeader");
		if (element instanceof HTMLInputElement) {
			element.checked = false;
		} else {
			console.error("The element 'bulkCheckboxHeader' is not a valid input element.");
		}
	} catch (error) {
		console.error("Error searching register routes:", error)
	} finally {
		isLoading.value = false
	}
}

function bulkUpdate() {
	modalHideShow(true, "register-route-bulk-modal")
}

// comes from the 'emit save' click within the bulk modal
function saveRegisterRouteBulkModalClick(bulkRoles: BulkRoles) {
	const addRoles = roles.value.reduce((addRoles: string[], r: Role, index: number) => {
		if (bulkRoles.add_roles[index]) {
			addRoles.push(r.name)
		}
		return addRoles
	}, [])
	const removeRoles = roles.value.reduce((removeRoles: string[], r: Role, index: number) => {
		if (bulkRoles.remove_roles[index]) {
			removeRoles.push(r.name)
		}
		return removeRoles
	}, [])
	const body = {raw_paths: bulkUpdateRouteArray.value, add_roles: addRoles, remove_roles: removeRoles}
	fetch("/register-route/bulk", {method: "POST", body: body, headers: getAuthHeader()})
	.then(() => {
		handleSearch()
		bulkUpdateRouteArray.value = []
	})
	modalHideShow(false, "register-route-bulk-modal")
}

// comes from the 'emit cancel' click within the modal
function cancelRegisterRouteBulkModalClick() {
	modalHideShow(false, "register-route-bulk-modal")
}

function bulkCheckboxChangeHeader(e: Event) {
	const checked = (e.target as HTMLInputElement).checked
	if (checked) {
		registerRoutes.value.forEach((r: RegisterRoute) => {
			if (!bulkUpdateRouteArray.value.includes(r.raw_path)) {
				bulkUpdateRouteArray.value.push(r.raw_path)
			}
		})
	} else {
		bulkUpdateRouteArray.value = []
	}
}

watch(registerRouteTerm, debounce(handleSearch, 500))

const disabledUpdateButton = computed(() => {
	return bulkUpdateRouteArray.value.length === 0
})	

onMounted(() => {
	loadRegisterRoutes()
	loadRoles()
})

</script>

<template>
	<p class="text-4xl">Admin</p>
	<p class="mt-8 mb-2 text-2xl">Register Routes</p>
	<div class="flex flex-column mb-4">
		<input id="register_route_search" v-model="registerRouteTerm" class="bg-gray-200 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500 w-96" placeholder="Search..."/>
		<button id="bulk-update" :disabled="disabledUpdateButton" @click="bulkUpdate" class="ml-4 text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 me-2 dark:bg-blue-600 dark:hover:bg-blue-700 focus:outline-none dark:focus:ring-blue-800 disabled:bg-gray-600 disabled:hover:bg-gray-500">Bulk Update</button>
	</div>
	<p class="mt-4 mb-2 text-sm">No Roles are considered 'open' routes</p>
	<div class="relative overflow-x-auto shadow-md sm:rounded-lg">
		<table class="w-full text-sm text-left rtl:text-right text-gray-600 dark:text-gray-400">
			<thead class="text-xs text-gray-700 uppercase bg-gray-400 dark:bg-gray-600 dark:text-gray-200">
				<tr> 
					<th class="px-6 py-3">
						<input id="bulkCheckboxHeader" type="checkbox" class="form-checkbox h-4 w-4 text-blue-600" @change="bulkCheckboxChangeHeader" />
					</th>
					<th scope="col" class="px-6 py-3">
						Raw Path
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
				<tr class="odd:bg-white odd:dark:bg-gray-900 even:bg-gray-200 even:dark:bg-gray-800 border-b dark:border-gray-700" v-for="registerRoute in registerRoutes">
					<th class="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white">
						<input type="checkbox" class="form-checkbox h-4 w-4 text-blue-600" :value="registerRoute.raw_path" v-model="bulkUpdateRouteArray" />
					</th>
					<th scope="row" class="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white">
						{{ registerRoute.raw_path }}
					</th>
					<td class="px-6 py-4">
						{{ registerRoute.roles }}
					</td>
					<td class="px-6 py-4">
						<a href="#" @click="editRegisterRouteClick(registerRoute.raw_path)" class="font-medium text-blue-600 dark:text-blue-500 hover:underline">Edit</a>
					</td>
				</tr>
			</tbody>
		</table>
	</div>
	<RegisterRouteModal id="register-route-modal" class="hidden" :registerRoute="currentRegisterRoute" :roles="roles" @clickSubmit="saveRegisterRouteModalClick" @clickCancel="cancelRegisterRouteModalClick" :titleText="(modeAdd ? 'Add' : 'Edit') + ' RegisterRoute'" />
	<RegisterRouteBulkModal id="register-route-bulk-modal" class="hidden" :rawPathCount="bulkUpdateRouteArray.length" :roles="roles" @clickSubmit="saveRegisterRouteBulkModalClick" @clickCancel="cancelRegisterRouteBulkModalClick" />
</template>