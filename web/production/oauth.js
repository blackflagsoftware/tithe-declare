function StartOAuth(clientId, clientName, scope, authId) {
	if (clientId === "") {
		return {result: false, message: "client id is not set"}
	}
	if (clientName === "") {
		return {result: false, message: "client name is not set"}
	}
	if (authId === "") {
		return {result: false, message: "auth id is not set"}
	}
	const params = new URLSearchParams({
		client_id: clientId,
		client_name: clientName,
		scope: scope,
		auth_id: authId,
	})
	
	const consentWindow = 
	window.onmessage = function (e) {
		if (e.data) {
			return {result: true, message: ""}
		} else {
			return {result: false, message: "operation cancelled"}
		}
	}
	var strWindowFeatures = "location=yes,height=570,width=520,scrollbars=yes,status=yes";
	window.open("https://localhost/consent?" + params.toString(), "_blank", strWindowFeatures)
	console.log("opened window")
} 

// function async openWindow() Promise<void> {

// }