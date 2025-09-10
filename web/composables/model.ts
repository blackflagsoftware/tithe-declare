export interface Login {
	id: string,
	first_name: string,
	last_name: string,
	email_address: string,
	roles: string[],
	roles_selected: boolean[]
}

export interface Role {
	id: string,
	name: string,
	description: string
}

export interface AuthClient {
	id: string,
	name: string,
	description: string,
	homepage_url: string
}

export interface Secret {
	id: string
	client_id: string
	secret: string
	active: boolean
}

export interface Callback {
	client_id: string
	callback_url: string
}

export interface RegisterRoute {
	raw_path: string
	roles: Role[]
	roles_selected: boolean[]
}

export interface BulkRoles {
	add_roles: boolean[]
	remove_roles: boolean[]
}

export function BlankLogin(roles: Role[]): Login {
	return {id: "", first_name: "", last_name: "", email_address: "", roles: [], roles_selected: roles.map(r => {return false})}
}

export function BlankRole(): Role {
	return {id: "", name: "", description: ""}
}

export function BlankAuthClient(): AuthClient {
	return {id: "", name: "", description: "", homepage_url: ""}
}

export function BlankSecret(): Secret {
	return {id: "", client_id: "", secret: "", active: true}
}

export function BlankCallback(): Callback {
	return {client_id: "", callback_url: ""}
}

export function BlankRegisterRoute(roles: Role[]): RegisterRoute {
	return {raw_path: "", roles: [], roles_selected: roles.map(r => {return false})}
}
