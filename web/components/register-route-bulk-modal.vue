<script setup lang="ts">
import { ref } from "vue"

const emit = defineEmits(["clickSubmit", "clickCancel"])
const props = defineProps<{
	rawPathCount: number 
	roles: Role[]
}>()

const addRolesChecked = ref<boolean[]>(new Array<boolean>())
const removeRolesChecked = ref<boolean[]>(new Array<boolean>())

function saveClick() {
	const roles = { add_roles: addRolesChecked.value, remove_roles: removeRolesChecked.value }
	blankLocalRef()
	emit("clickSubmit", roles)
}

function cancelClick() {
	blankLocalRef()
	emit("clickCancel")
}

function roleAddChange(index: number) {
	addRolesChecked.value = addRolesChecked.value.map((r, idx) => 
		index === idx ? !r : r
	)
	if (addRolesChecked.value[index]) {
		removeRolesChecked.value[index] = false
	}
}

function roleRemoveChange(index: number) {
	removeRolesChecked.value = removeRolesChecked.value.map((r, idx) => 
		index === idx ? !r : r
	)
	if (removeRolesChecked.value[index]) {
		addRolesChecked.value[index] = false
	}
}

function blankLocalRef() {
	addRolesChecked.value = []
	removeRolesChecked.value = []
}

watch(() => props.rawPathCount, (newRawPathCount) => {
		addRolesChecked.value = new Array<boolean>(props.roles.length).fill(false)
		removeRolesChecked.value = new Array<boolean>(props.roles.length).fill(false)
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
							<p class="text-xl font-semibold text-gray-900 dark:text-white">Bulk Add/Remove Roles</p>
							<button type="button" class="text-gray-400 bg-transparent hover:bg-gray-50 hover:text-gray-900 rounded-lg text-sm w-8 h-8 ms-auto inline-flex justify-center items-center dark:hover:bg-gray-600 dark:hover:text-white" @click="cancelClick()">
								<svg class="w-3 h-3" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 14 14">
									<path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="m1 1 6 6m0 0 6 6M7 7l6-6M7 7l-6 6"/>
								</svg>
								<span class="sr-only">Close modal</span>
							</button>
						</div>
						<div class="p-4 md:p-5 space-y-8 dark:bg-gray-900">
							<div>
								<label class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Apply changes for ({{ props.rawPathCount }}) paths</label>
								<p class="text-base mt-4 mb-2">Add Roles</p>
								<ul class="list-none pl-0 mt-0">
									<div v-for="(role, index) in props.roles">
										<li :key="index">
											<input :id="role.id.toString()" type="checkbox" :value="role.name" :checked="addRolesChecked![index]" class="mr-2" @change="roleAddChange(index)"/>
											<label :htmlFor="role.id">{{ role.name }}</label>
										</li>
									</div>
								</ul>
								<p class="text-base mt-4 mb-2">Remove Roles</p>
								<ul class="list-none pl-0 mt-0">
									<div v-for="(role, index) in props.roles">
										<li :key="index">
											<input :id="role.id.toString()" type="checkbox" :value="role.name" :checked="removeRolesChecked![index]" class="mr-2" @change="roleRemoveChange(index)"/>
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