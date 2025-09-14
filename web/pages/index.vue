<script setup lang="ts">
import type { TitheDeclareDate, ValidDateAndTimes } from "~/types/types"

const route = useRoute()
const daySelected = ref("")
const timeSelected = ref("")
const dateOptions = ref<string[]>(new Array<string>())
const timeOptions = ref<string[]>(new Array<string>())
const validDateAndTimes = ref<ValidDateAndTimes>({})
const { getAuthHeader } = useAuth()
const fetch = apiFetch()
const date = ref<string>("")
const timeStart = ref<string>("")
const timeEnd = ref<string>("")
const showUserForm = ref<boolean>(false)
const showConfirmation = ref<boolean>(false)
const name = ref<string>("")
const email = ref<string>("")
const phone = ref<string>("")
const msg = ref<string>("")
const upcomingDeclarations = ref<TitheDeclareDate[]>([])
const emailReminder = ref<string>("")
const knownEmails = ref<string[]>([])

function loadCurrentDays() {
	fetch("/td-date/current-days", {method: "GET", headers: getAuthHeader()})
	.then(response => {
		validDateAndTimes.value = response.data as ValidDateAndTimes
		dateOptions.value = Object.keys(validDateAndTimes.value)
		if (daySelected.value !== "") {
			const timeArray = validDateAndTimes.value[daySelected.value] || []
			timeOptions.value = timeArray
		}
	})
	.catch(error => {
		console.error("Error fetching current days:", error)
	})
}

function onDayChange() {
	const timeArray = validDateAndTimes.value[daySelected.value] || []
	timeOptions.value = timeArray
	timeSelected.value = ""
	showUserForm.value = false
}

function onTimeChange() {
	const body = JSON.stringify({
		date: daySelected.value,
		time: timeSelected.value
	})
	fetch("/td-date/check-hold-time", {method: "POST", body: body, headers: getAuthHeader()})
	.then(() => {
		showUserForm.value = true
	})
	.catch(error => {
		error.response && error.response.status === 423 ? alert("This time is already held. Please select another time.") : console.error("Error checking hold time:", error)
		timeSelected.value = ""
		loadCurrentDays()
	})
}

function onConfirmClick() {
	const body = JSON.stringify({
		date: daySelected.value,
		time: timeSelected.value,
		name: name.value,
		email: email.value,
		phone: phone.value
	})
	fetch("/td-date/confirm", {method: "POST", body: body, headers: getAuthHeader()})
	.then(() => {
		msg.value = `Your declaration date/time: ${dateFormat(daySelected.value)} ${timeSelected.value}. You will receive a reminder if you provided an email.`
		daySelected.value = ""
		timeSelected.value = ""
		name.value = ""
		email.value = ""
		phone.value = ""
		showUserForm.value = false
		showConfirmation.value = true
		loadCurrentDays()
	})
	.catch(error => {
		console.error("Error confirming date and time:", error)
	})
}

function dateFormat(d: string) {
	// change the format of the date from 2024-08-30 to August 30, 2024
	const options: Intl.DateTimeFormatOptions = { year: 'numeric', month: 'long', day: 'numeric' };
	const date = new Date(d);
	date.setUTCHours(12); 
	return date.toLocaleDateString(undefined, options);
}

const disabledTime = computed(() => {
	return daySelected.value === ""
})

const hideUserForm = computed(() => {
	return !showUserForm.value
})

const disabledConfirm = computed(() => {
	return name.value === ""
})


// admin functionality
const amAdmin = computed(() => {
	return route.query.admin === useRuntimeConfig().public.adminPwd ? true : false
})

function addDate() {
	fetch("/td-date/block", {
		method: "POST",
		headers: {
			...getAuthHeader(),
		},
		body: JSON.stringify({
			new_date: date.value,
			start_time: timeStart.value,
			end_time: timeEnd.value
		})
	})
	.catch(error => {
		console.error("Error adding date:", error)
	})	
}

function SendMessage(id: number) {
	// fetch(`/td-date/send-reminder/${id}`, {
	// 	method: "POST",
	// 	headers: {
	// 		...getAuthHeader(),
	// 	},
	// })
	// .catch(error => {
	// 	console.error("Error sending reminder:", error)
	// })	
	alert(`Pretend sending reminder for id ${id}`)
}

function AddEmail() {
	if (emailReminder.value === "") {
		alert("Please enter an email address")
		return
	}
	fetch("/email-reminder", {
		method: "POST",
		headers: {
			...getAuthHeader(),
		},
		body: JSON.stringify({
			email: emailReminder.value,
		})
	})
	.then(() => {
		if (!knownEmails.value.includes(emailReminder.value)) {
			knownEmails.value.push(emailReminder.value)
		}
		emailReminder.value = ""
	})
	.catch(error => {
		console.error("Error adding email reminder:", error)
	})
}

onMounted(() => {
	loadCurrentDays()
	if (amAdmin.value) {
		// load upcoming declarations
		const body = JSON.stringify({
			search : {
				filters: [
					{
						column: "date_value",
						compare: ">",
						value: new Date().toISOString(),
					},
					{
						column: "confirm",
						compare: "NOT NULL",
						value: "",
					},
				]
			},
		})
		fetch("/td-date/search", {method: "POST", body: body, headers: getAuthHeader()})
		.then(response => {
			upcomingDeclarations.value = response.data as TitheDeclareDate[]
		})
		.catch(error => {
			console.error("Error fetching upcoming declarations:", error)
		})
		fetch("/email-reminder/search", {method: "POST", body: {}, headers: getAuthHeader()})
		.then(response => {
			knownEmails.value = response.data.map((er: {email: string}) => er.email) as string[]
		})
		.catch(error => {
			console.error("Error fetching known emails:", error)
		})
	}	
})
</script>
<template>
	<div class="flex flex-col items-center h-full">
		<h1 class="text-4xl font-bold mb-4 text-black dark:text-white">River Ridge 11th Ward Tithing Declaration Sign-up</h1>
		<div class="flex flex-row space-x-8">
			<div class="flex flex-col">
				<label for="date" class="mb-2 text-black dark:text-white">Date:</label>
				<select id="date" class="border border-gray-300 rounded p-2 text-gray-800" v-model="daySelected" @change="onDayChange">
					<option disabled value="">Select a Day</option>
					<option v-for="d in dateOptions" :key="d" :value="d">{{ dateFormat(d) }}</option>
				</select>
			</div>
			<div class="flex flex-col">
				<label for="time" class="mb-2 text-black dark:text-white">Time:</label>
				<select id="time" class="border border-gray-300 rounded p-2 text-gray-800" v-model="timeSelected" @change="onTimeChange" :disabled="disabledTime">
					<option disabled value="">Select a Time</option>
					<option v-for="t in timeOptions" :key="t" :value="t">{{ t }}</option>
				</select>
			</div>
		</div>
		<div v-if="showUserForm" class="mt-5 p-4 text-white rounded w-full max-w-md">
				<div class="mb-4">
					<label for="name" class="block text-black dark:text-white font-bold mb-2">Name:</label>
					<input id="name" type="text" class="w-full border border-gray-300 rounded p-2 text-gray-800" v-model="name" />
				</div>
				<p class="mb-4 text-sm text-white">If you want a reminder, please enter your email</p>
				<div class="mb-4">
					<label for="email" class="block text-black dark:text-white font-bold mb-2">Email</label>
					<input id="email" type="email" class="w-full border border-gray-300 rounded p-2 text-gray-800" v-model="email"/>
				</div>
				<!-- <div class="mb-4">
					<label for="phone" class="block text-black dark:text-white font-bold mb-2">Text</label>
					<input id="phone" type="text" class="w-full border border-gray-300 rounded p-2 text-gray-800" v-model="phone"/>
				</div> -->
				<button type="submit" class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 disabled:bg-gray-500" @click="onConfirmClick" :disabled="disabledConfirm">Confirm</button>
		</div>
		<div v-if="showConfirmation" class="mt-5 p-4 border border-gray-500 bg-gray-100 text-gray-700 rounded">
			<p>{{ msg }}</p>
			<button class="mt-2 bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600" @click="showConfirmation = false">Close</button>
		</div>
	</div>
	<div v-if="amAdmin" class="mt-5 p-4 border border-red-500 bg-red-100 text-red-700 rounded">
		<div class="flex flex-row">
			<div class="flex flex-col">
				<label for="date-block" class="mb-2 font-bold">Add Date</label>
				<input id="date-block" type="text" class="border border-gray-300 rounded p-2" v-model="date" placeholder="2025-08-30"/>
			</div>
			<div class="flex flex-col ml-4">
				<label for="time-start" class="mb-2 font-bold">Start Time</label>
				<input id="time-start" type="text" class="border border-gray-300 rounded p-2" v-model="timeStart" placeholder="12:00" />
			</div>
			<div class="flex flex-col ml-4">
				<label for="time-end" class="mb-2 font-bold">End Time</label>
				<input id="time-end" type="text" class="border border-gray-300 rounded p-2" v-model="timeEnd" placeholder="15:00" />
			</div>
			<button class="ml-4 bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 self-end" @click="addDate">Add</button>
		</div>
		<div>
			<p class="mt-2 text-lg">Upcoming appointments</p>
			<div v-for="appt in upcomingDeclarations" :key="appt.id" class="border-b border-gray-300 py-2">
				<div class="flex flex-row justify-between">
					<div class="flex flex-row gap-4">
					<p class="font-bold mr-2">ID: {{ appt.id }}</p>
					<p>{{ appt.date_value }}</p>
					<p v-if="appt.name">Name: {{ appt.name }}</p>
					<p v-if="appt.email">Email: {{ appt.email }}</p>
					<p v-if="appt.phone">Phone: {{ appt.phone }}</p>
					</div>
					<div>
						<button class="ml-4 bg-blue-500 text-white px-2 py-1 rounded hover:bg-blue-600 self-end" @click="SendMessage(appt.id)">Send Reminder</button>
					</div>
				</div>
			</div>
			<div class="mt-4">
				<p>Email Reminder</p>
				<input type="text" class="border border-gray-300 rounded p-2" placeholder="Email Address" v-model="emailReminder"/>
				<button class="ml-4 bg-blue-500 text-white px-2 py-1 rounded hover:bg-blue-600 self-end" @click="AddEmail">Add Email</button>
			</div>
			<div v-if="knownEmails.length > 0">
				<p class="mt-2">Known Emails:</p>
				<ul class="list-disc list-inside">
					<li v-for="email in knownEmails" :key="email">{{ email }}</li>
				</ul>
			</div>
		</div>
	</div>
</template>