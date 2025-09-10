<script setup lang="ts">
const props = defineProps<{
	clientId: string
}>()

const secrets = ref<Secret[]>()
const modeAdd = ref<boolean>(false)
const secret = ref<Secret>(BlankSecret())
const { getAuthHeader } = useAuth()
const fetch = apiFetch()

function loadSecrets() {
	const body = {search: {filters: [{column: "client_id", value: props.clientId, compare: "="}]}}
	fetch("/auth-client-secret/search", {method: "POST", body: body, headers: getAuthHeader()})
	.then(response => {
		secrets.value = response.data
	})
}

function addClick() {
	modeAdd.value = true
	modalHideShow(true, "auth-secret-modal")
}

function editClick(secretIn: string) {
	secret.value = secrets.value?.find(s => s.secret === secretIn)!
	modeAdd.value = false
	modalHideShow(true, "auth-secret-modal")
}

function deleteClick(secretIn: string) {
	secret.value = secrets.value?.find(s => s.secret === secretIn)!
	modalHideShow(true, "auth-secret-delete-modal")
}

function secretSecret(secret: string): string {
	return (secret.length > 4 ? secret.substring(0, 4) : secret) + "**********"
}

// 'emit' from secret modal
function submitClick(secret: Secret) {
	modalHideShow(false, "auth-secret-modal")
	const body = {...secret, client_id: props.clientId}
	if (modeAdd.value) {
		fetch("/auth-client-secret", {method: "POST", body: body, headers: getAuthHeader()})
		.then(() => {
			loadSecrets()
		})
	} else {
		fetch("/auth-client-secret", {method: "PATCH", body: body, headers: getAuthHeader()})
		.then(() => {
			loadSecrets()
		})
	}
}

// 'emit' from secret modal
function cancelClick(secret: Secret) {
	modalHideShow(false, "auth-secret-modal")
}

// 'emit' from secret delete modal
function submitDeleteClick() {
	const secretId = secret.value.id
	if (!secretId || secretId === "") {
		console.log("Secret id is not set, something went wrong") // TODO: alert
		return
	}
	fetch("/auth-client-secret/" + secretId, {method: "DELETE", headers: getAuthHeader()})
	.then(() => {
		secret.value = BlankSecret()
		loadSecrets()
	})
	modalHideShow(false, "auth-secret-delete-modal")
}

// 'emit' from secret delete modal
function cancelDeleteClick() {
	secret.value = BlankSecret()
	modalHideShow(false, "auth-secret-delete-modal")
}

const noClientId = computed(() => {
	return props.clientId === ""
})

watch(() => props.clientId, (newClient) => {
	if (newClient !== "") {
		loadSecrets()
	} else {
		secrets.value = []
	}
})

</script>

<template>
	<div>
		<p class="mb-2 text-xl">Secrets</p>
		<button id="add" :disabled="noClientId" @click="addClick" class="mt-4 mb-4 text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 me-2 dark:bg-blue-600 dark:hover:bg-blue-700 focus:outline-none dark:focus:ring-blue-800 disabled:bg-gray-600 disabled:hover:bg-gray-500">Add Secret</button>
		<div class="overflow-x-auto shadow-md sm:rounded-lg">
			<table class="w-full text-sm text-left rtl:text-right text-gray-600 dark:text-gray-400">
					<thead class="text-xs text-gray-700 uppercase bg-gray-400 dark:bg-gray-500 dark:text-gray-200">
						<tr>
							<th scope="col" class="px-6 py-3">
								Secret
							</th>
							<th scope="col" class="px-6 py-3">
								Active	
							</th>
							<th scope="col" class="px-6 py-3">
								Action
							</th>
						</tr>
				</thead>
				<tbody>
						<tr class="odd:bg-white odd:dark:bg-gray-900 even:bg-gray-200 even:dark:bg-gray-800 border-b dark:border-gray-700" v-for="secret in secrets">
							<th scope="row" class="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white">
								{{ secretSecret(secret.secret) }}
							</th>
							<td class="px-6 py-4">
								{{ secret.active }}
							</td>
							<td class="px-6 py-4">
								<a href="#" @click="editClick(secret.secret)" class="font-medium text-blue-600 dark:text-blue-500 hover:underline">Edit</a>
								|
								<a href="#" @click="deleteClick(secret.secret)" class="font-medium text-blue-600 dark:text-blue-500 hover:underline">Delete</a>
							</td>
						</tr>
				</tbody>
			</table>
		</div>
		<AuthSecretModal id="auth-secret-modal" class="hidden" @clickSubmit="submitClick" @clickCancel="cancelClick" :secret="secret" :titleText="(modeAdd ? 'Add' : 'Edit') + ' Secret'" />
		<DeleteModal id="auth-secret-delete-modal" class="hidden" recordType="Secret" :name="secret.secret.substring(0, 4) + '******'" @clickSubmit="submitDeleteClick" @clickCancel="cancelDeleteClick" />
	</div>
</template>