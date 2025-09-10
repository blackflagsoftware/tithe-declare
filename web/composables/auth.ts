import { Buffer } from "buffer"

export const loggedInState = () => useState<boolean>("isLoggedIn", () => false)

export const useAuth = () => {
	const isLoggedIn = loggedInState()

	const login = (redirectPath: string = "/") => {
		isLoggedIn.value = true
		if (import.meta.client) {
			sessionStorage.setItem("isLoggedIn", "true")
		}
		return navigateTo(redirectPath)
	}

	const logout = (redirectPath: string = "/login") => {
		isLoggedIn.value = false
		if (import.meta.client) {
			sessionStorage.removeItem("isLoggedIn")
		}
		return navigateTo(redirectPath)
	}

	const initializeAuth = () => {
		if (import.meta.client) {
			const isLoggedInLocal = sessionStorage.getItem("isLoggedIn")
			isLoggedIn.value = isLoggedInLocal === "true"
		}
	}

	const getUnAuthHeader = (version: string = "v1") => {
		
		const appName = useRuntimeConfig().public.appName
		return { Accept: "application/vnd." + appName + "." + version + "+json" }
	}

	const getAuthHeader = (version: string = "v1") => {
		const appName = useRuntimeConfig().public.appName
		const token = btoa("test:test") // dGVzdDp0ZXN0Cg==" // sessionStorage.getItem("token")
		return { Authorization: "Basic " + token, Accept: "application/vnd." + appName + "." + version + "+json" }
	}

	const checkPermission = (checkRole: string): boolean => {
		if (import.meta.server) return false
		try {
			const token = sessionStorage.getItem("token")
			if (token === null || token === undefined) {
				console.log("token is missing")
				return false
			}
			const tokenSplit = token?.split(".") ?? []
			const encodeToken = tokenSplit[1]
			const decode = (str: string): string => Buffer.from(str, "base64").toString("binary")
			const decodedToken = JSON.parse(decode(encodeToken))
			if (decodedToken.roles === null || decodedToken.roles === undefined) {
				return false
			}
			for (let i = 0; i < decodedToken.roles.length; i++) {
				if (checkRole === decodedToken.roles[i]) {
					return true
				}
			}
		} catch (error) {
			console.log(error)
		}
		return false
	}

	return {
		isLoggedIn,
		login,
		logout,
		initializeAuth,
		getUnAuthHeader,
		getAuthHeader,
		checkPermission,
	}
}