<script setup lang="ts">
const props = defineProps<{
	clientId: string
}>()

const callbacks = ref<Callback[]>()
const modeAdd = ref<boolean>(false)
const callback = ref<Callback>(BlankCallback())
const deleteCallbackUrl = ref<string>("")
const { getAuthHeader } = useAuth()
const fetch = apiFetch()

function loadCallbacks() {
	const body = {search: {filters: [{column: "client_id", value: props.clientId, compare: "="}]}}
	fetch("/auth-client-callback/search", {method: "POST", body: body, headers: getAuthHeader()})
	.then(response => {
		callbacks.value = response.data
	})
}

function addClick() {
	modeAdd.value = true
	modalHideShow(true, "auth-callback-modal")
}

function editClick(callback_urlIn: string) {
	deleteCallbackUrl.value = callback_urlIn // this will be used if 'editing'
	callback.value = callbacks.value?.find(c => c.callback_url === callback_urlIn)!
	modeAdd.value = false
	modalHideShow(true, "auth-callback-modal")
}

function deleteClick(callbackUrlIn: string) {
	callback.value = callbacks.value?.find(s => s.callback_url === callbackUrlIn)!
	modalHideShow(true, "auth-callback-delete-modal")
}

// 'emit' from callback modal
function submitClick(callback: Callback) {
	modalHideShow(false, "auth-callback-modal")
	if (!modeAdd.value) {
		// delete the old one, add new one for 'edit'
		fetch("/auth-client-callback/" + props.clientId + "/callback_url/" + deleteCallbackUrl.value, {method: "DELETE", headers: getAuthHeader()})
		.then(() => {
			deleteCallbackUrl.value = ""
		})
	}
	const body = {...callback, client_id: props.clientId}
	fetch("/auth-client-callback", {methods: "POST", body: body, headers: getAuthHeader()})
	.then(() => {
		loadCallbacks()
	})
}

// 'emit' from callback modal
function cancelClick(callback: Callback) {
	modalHideShow(false, "auth-callback-modal")
}

// 'emit' from callback delete modal
function submitDeleteClick() {
	const callbackUrl = callback.value.callback_url
	if (!callbackUrl || callbackUrl === "") {
		console.log("CallbackUrl is not set, something went wrong") // TODO: alert
		return
	}
	fetch("/auth-client-callback/" + props.clientId + "/callback_url/" + callbackUrl, {method: "DELETE", headlers: getAuthHeader()})
	.then(() => {
		callback.value = BlankCallback()
		loadCallbacks()
	})
	modalHideShow(false, "auth-callback-delete-modal")
}

// 'emit' from callback delete modal
function cancelDeleteClick() {
	callback.value = BlankCallback()
	deleteCallbackUrl.value = ""
	modalHideShow(false, "auth-callback-delete-modal")
}

const noClientId = computed(() => {
	return props.clientId === ""
})

watch(() => props.clientId, (newClient) => {
	if (newClient !== "") {
		loadCallbacks()
	} else {
		callbacks.value = []
	}
})

</script>

<template>
	<div>
		<p class="mb-2 text-xl">Callbacks</p>
		<button id="add" :disabled="noClientId" @click="addClick" class="mt-4 mb-4 text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 me-2 dark:bg-blue-600 dark:hover:bg-blue-700 focus:outline-none dark:focus:ring-blue-800 disabled:bg-gray-600 disabled:hover:bg-gray-500">Add Callback</button>
		<div class="overflow-x-auto shadow-md sm:rounded-lg">
			<table class="w-full text-sm text-left rtl:text-right text-gray-600 dark:text-gray-400">
					<thead class="text-xs text-gray-700 uppercase bg-gray-400 dark:bg-gray-500 dark:text-gray-200">
						<tr>
							<th scope="col" class="px-6 py-3">
								Callback
							</th>
							<th scope="col" class="px-6 py-3">
								Action
							</th>
						</tr>
				</thead>
				<tbody>
						<tr class="odd:bg-white odd:dark:bg-gray-900 even:bg-gray-200 even:dark:bg-gray-800 border-b dark:border-gray-700" v-for="callback in callbacks">
							<th scope="row" class="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white">
								{{ callback.callback_url }}
							</th>
							<td class="px-6 py-4">
								<a href="#" @click="editClick(callback.callback_url)" class="font-medium text-blue-600 dark:text-blue-500 hover:underline">Edit</a>
								|
								<a href="#" @click="deleteClick(callback.callback_url)" class="font-medium text-blue-600 dark:text-blue-500 hover:underline">Delete</a>
							</td>
						</tr>
				</tbody>
			</table>
		</div>
		<AuthCallbackModal id="auth-callback-modal" class="hidden" @clickSubmit="submitClick" @clickCancel="cancelClick" :callback="callback" :titleText="(modeAdd ? 'Add' : 'Edit') + ' Callback'" />
		<DeleteModal id="auth-callback-delete-modal" class="hidden" recordType="Callback" :name="callback.callback_url" @clickSubmit="submitDeleteClick" @clickCancel="cancelDeleteClick" />
	</div>
</template>