<script setup lang="ts">
import { ref } from "vue"

const emit = defineEmits(["clickSubmit", "clickCancel"])
const props = defineProps<{
	titleText: string
	registerRoute: RegisterRoute 
	roles: Role[]
}>()

const titleText = ref<string>("Register Route")
const rawPath = ref<string>("")
const registerRouteRoles = ref<string[]>(new Array<string>())
const registerRouteRolesChecked = ref<boolean[]>(new Array<boolean>())

function saveClick() {
	const registerRoute = {raw_path: rawPath.value, roles_selected: registerRouteRolesChecked.value}
	emit("clickSubmit", registerRoute)
	blankLocalRef()
}

function cancelClick() {
	blankLocalRef()
	emit("clickCancel")
}

function roleChange(index: number) {
	registerRouteRolesChecked.value = registerRouteRolesChecked.value.map((r, idx) => 
		index === idx ? !r : r
	)
}

function blankLocalRef() {
	rawPath.value = ""
	registerRouteRoles.value = []
	registerRouteRolesChecked.value = []
}

watch(() => props.registerRoute, (newRegisterRoute) => {
		if (newRegisterRoute) {
			rawPath.value = newRegisterRoute.raw_path
			registerRouteRoles.value = newRegisterRoute.roles
			registerRouteRolesChecked.value = newRegisterRoute.roles_selected
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
								<label class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">{{ rawPath }}</label>
								<p class="text-base mt-4 mb-2">Roles</p>
								<ul class="list-none pl-0 mt-0">
									<div v-for="(role, index) in props.roles">
										<li :key="index">
											<input :id="role.id.toString()" type="checkbox" :value="role.name" :checked="registerRouteRolesChecked![index]" class="mr-2" @change="roleChange(index)"/>
											<label :htmlFor="role.id">{{ role.name }}</label>
										</li>
									</div>
								</ul>
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