<script setup lang="ts">
import { ref } from "vue"

const emit = defineEmits(["clickSubmit", "clickCancel"])
const props = defineProps<{
	titleText: string
	secret: Secret 
}>()

const titleText = ref<string>("Secret")
const id = ref<string>("")
const secretVal = ref<string>("")
const active = ref<boolean>(true)

function saveClick() {
	const secret = {id: id.value, secret: secretVal.value, active: active.value}
	emit("clickSubmit", secret)
	blankLocalRef()
}

function cancelClick() {
	blankLocalRef()
	emit("clickCancel")
}

function blankLocalRef() {
	secretVal.value = ""
	active.value = false
}

watch(() => props.secret, (newSecret) => {
		if (newSecret) {
			console.log("document:", document)
			id.value = newSecret.id
			secretVal.value = newSecret.secret
			active.value = newSecret.active
			active.value ? document.getElementById("active-input")?.setAttribute("checked", "checked") : document.getElementById("active-input")?.removeAttribute("checked")
		}
	},
	{ immediate: true }
)

watch(() => props.titleText, (newTitleText) => {
		titleText.value = newTitleText
	},
	{ immediate: true}
)
</script>

<template>
	<div class="fixed inset-0 transition-opacity bg-gray-50 bg-opacity-75">
		<div tabindex="-1" class="fixed z-50 top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-auto bg-white dark:bg-gray-900 rounded-md px-4 py-2 space-y-5 drop-shadow-lg">
			<div class="relative p-4 w-full max-w-2xl max-h-full">
				<div class="relative bg-white dark:bg-gray-900 rounded-lg shadow dark:bg-black">
						<div class="flex items-center justify-between p-4 md:p-5 border-b rounded-t dark:bg-gray-900 dark:border-gray-600">
								<p class="text-xl font-semibold text-gray-900 dark:text-white">{{ titleText }}</p>
								<button type="button" class="text-gray-400 bg-transparent hover:bg-gray-50 hover:text-gray-900 rounded-lg text-sm w-8 h-8 ms-auto inline-flex justify-center items-center dark:hover:bg-gray-600 dark:hover:text-white" @click="cancelClick()">
										<svg class="w-3 h-3" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 14 14">
												<path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="m1 1 6 6m0 0 6 6M7 7l6-6M7 7l-6 6"/>
										</svg>
										<span class="sr-only">Close modal</span>
								</button>
						</div>
						<div class="p-4 md:p-5 space-y-8 dark:bg-gray-900">
							<div>
								<label for="secret" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Secret</label>
								<textarea id="secret" rows="6" v-model="secretVal" class="block p-2.5 max-w-screen-lg w-96 text-sm text-gray-900 bg-gray-200 rounded-lg border border-gray-300 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"></textarea>
							</div>
							<div>
								<label class="inline-flex items-center cursor-pointer">
									<input id="active-input" type="checkbox" v-model="active" class="sr-only peer" checked>
									<div class="relative w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600"></div>
									<span class="ms-3 text-sm font-medium text-gray-900 dark:text-gray-300">Active</span>
								</label>
							</div>
						</div>
						<div class="flex items-center p-4 md:p-5 border-t border-gray-200 rounded-b dark:border-gray-600 dark:bg-gray-900">
								<button @click="saveClick()" type="button" class="text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800">Submit</button>
								<button @click="cancelClick()" type="button" class="py-2.5 px-5 ms-3 text-sm font-medium text-gray-900 focus:outline-none bg-white rounded-lg border border-gray-200 hover:bg-gray-100 hover:text-blue-700 focus:z-10 focus:ring-4 focus:ring-gray-100 dark:focus:ring-gray-700 dark:bg-gray-800 dark:text-gray-400 dark:border-gray-600 dark:hover:text-white dark:hover:bg-gray-700">Cancel</button>
						</div>
				</div>
			</div>
		</div>
  </div>
</template>