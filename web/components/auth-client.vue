<script setup lang="ts">
import { ref } from "vue"
import ClipboardJS from "clipboard"

const emit = defineEmits(["passClientId"])
const clients = ref<AuthClient[]>()
const client = ref<AuthClient>(BlankAuthClient())
const clientId = ref<string>("")
const modeAdd = ref<boolean>(true)
const { getAuthHeader} = useAuth()
const fetch = apiFetch()

function loadClients() {
	fetch("/auth-client/search", {methods: "POST", headers: getAuthHeader()})
	.then(response => {
		clients.value = response.data
	})
}

function addClick() {
	modeAdd.value = true
	modalHideShow(true, "auth-client-modal")
}

function editClick(clientId: string) {
	client.value = clients.value?.find(c => c.id === clientId)!
	modeAdd.value = false
	modalHideShow(true, "auth-client-modal")
}

function deleteClick(clientId: string) {
	client.value = clients.value?.find(c => c.id === clientId)!
	modalHideShow(true, "auth-client-delete-modal")
}

// 'emit' from client modal
function submitClick(client: Client) {
	const body = {...client}
	if (client.id === "") {
		fetch("/auth-client", {method: "POST", body: body, headers: getAuthHeader()})
		.then(() => {
			loadClients()
		})
	} else {
		fetch("/auth-client", {method: "PATCH", body: body, headers: getAuthHeader()})
		.then(() => {
			loadClients()
		})
	}
	modalHideShow(false, "auth-client-modal")
}

// 'emit' from client modal
function cancelClick() {
	client.value = BlankAuthClient()
	modalHideShow(false, "auth-client-modal")
}

// click from the 'Details' link
function detailsClick(clientIdIn: string) {
	clientId.value = clientIdIn
	// emit("passClientId", clientId)
}

// 'emit' from client delete modal
function submitDeleteClick() {
	const clientId = client.value.id
	if (!clientId || clientId === "") {
		console.log("Client id is not set, something went wrong") // TODO: alert
		return
	}
	fetch("/auth-client/" + clientId, {method: "DELETE", headers: getAuthHeader()})
	.then(() => {
		client.value = BlankAuthClient()
		loadClients()
	})
	modalHideShow(false, "auth-client-delete-modal")
}

// 'emit' from client delete modal
function cancelDeleteClick() {
	client.value = BlankAuthClient()
	modalHideShow(false, "auth-client-delete-modal")
}

function clientIdToClipboard() {
	const copyText = document.querySelector("#client-id")!
	navigator.clipboard.writeText(copyText.innerHTML.toString())
}

onMounted(() => {
	loadClients()
	const scripts = [
		"https://cdn.jsdelivr.net/npm/clipboard@2.0.11/dist/clipboard.min.js"
	]
	scripts.forEach(script => {
		let tag = document.head.querySelector(`[src="${ script }"`)!
		if (!tag) {
			tag = document.createElement("script")
			tag.setAttribute("src", script)
			tag.setAttribute("type", "text/javascript")
			document.head.appendChild(tag)
		}
	})
	new ClipboardJS(".btn")
})

</script>

<template>
	<p class="mt-8 mb-2 text-2xl">Clients</p>
	<button @click="addClick()" class="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 me-2 mb-2 dark:bg-blue-600 dark:hover:bg-blue-700 focus:outline-none dark:focus:ring-blue-800">Add Client</button>
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
								Homepage Url	
							</th>
							<th scope="col" class="px-6 py-3">
								Action
							</th>
						</tr>
				</thead>
				<tbody>
						<tr class="odd:bg-white odd:dark:bg-gray-900 even:bg-gray-200 even:dark:bg-gray-800 border-b dark:border-gray-700" v-for="client in clients">
							<th scope="row" class="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white">
								{{ client.name }}
							</th>
							<td class="px-6 py-4">
								{{ client.description }}
							</td>
							<td class="px-6 py-4">
								{{ client.homepage_url }}
							</td>
							<td class="px-6 py-4">
								<a href="#" @click="detailsClick(client.id)" class="font-medium text-blue-600 dark:text-blue-500 hover:underline">Details</a>
								|
								<a href="#" @click="editClick(client.id)" class="font-medium text-blue-600 dark:text-blue-500 hover:underline">Edit</a>
								|
								<a href="#" @click="deleteClick(client.id)" class="font-medium text-blue-600 dark:text-blue-500 hover:underline">Delete</a>
							</td>
						</tr>
				</tbody>
		</table>
	</div>
	<div class="ml-2 flex flex-col">
		<div class="mt-4 space-y-8">
			<div>
				<label htmlFor="client_id" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">ClientId</label>
				<div class="flex flex-row">
					<input id="client-id" type="text" readonly class="bg-gray-200 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-2.5 dark:bg-gray-600 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500 w-96" :value="clientId" />
					<button data-clipboard-target="#client-id" class="btn ml-2 flex select-none items-center gap-2 rounded-lg bg-gray-900 px-6 py-3 text-center align-middle font-sans text-xs font-bold uppercase text-white shadow-md shadow-gray-900/10 transition-all hover:shadow-lg hover:shadow-gray-900/20 focus:opacity-[0.85] focus:shadow-none active:opacity-[0.85] active:shadow-none disabled:pointer-events-none disabled:opacity-50 disabled:shadow-none" type="button">
						<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" aria-hidden="true" class="w-4 h-4 text-white">
							<path stroke-linecap="round" stroke-linejoin="round" d="M15.75 17.25v3.375c0 .621-.504 1.125-1.125 1.125h-9.75a1.125 1.125 0 01-1.125-1.125V7.875c0-.621.504-1.125 1.125-1.125H6.75a9.06 9.06 0 011.5.124m7.5 10.376h3.375c.621 0 1.125-.504 1.125-1.125V11.25c0-4.46-3.243-8.161-7.5-8.876a9.06 9.06 0 00-1.5-.124H9.375c-.621 0-1.125.504-1.125 1.125v3.5m7.5 10.375H9.375a1.125 1.125 0 01-1.125-1.125v-9.25m12 6.625v-1.875a3.375 3.375 0 00-3.375-3.375h-1.5a1.125 1.125 0 01-1.125-1.125v-1.5a3.375 3.375 0 00-3.375-3.375H9.75"></path>
						</svg>
						Copy
					</button>
				</div>
			</div>
		</div>
		<div class="mt-4 flex flex-row space-x-8 place-content-stretch">
			<AuthSecret :clientId="clientId" class="w-3/6"/>
			<AuthCallback :clientId="clientId" class="w-3/6"/>
		</div>
	</div>
	<DeleteModal id="auth-client-delete-modal" class="hidden" recordType="Client" :name="client.name" @clickSubmit="submitDeleteClick" @clickCancel="cancelDeleteClick" />
	<AuthClientModal id="auth-client-modal" class="hidden" :client="client" @clickSubmit="submitClick" @clickCancel="cancelClick" :titleText="(modeAdd ? 'Add' : 'Edit') + ' Client'" />
</template>