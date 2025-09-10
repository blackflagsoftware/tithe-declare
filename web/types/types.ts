export interface TitheDeclareDate {
	id: number
	date_value: string
	hold: string
	confirm: string
	name: string
	phone: string
	email: string
}

export interface ValidDateAndTimes {
  [date: string]: string[];
}