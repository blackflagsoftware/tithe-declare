export function modalHideShow(open: boolean, id: string) {
	var modal = document.getElementById(id)!
	open ? modal.classList.remove("hidden") : modal.classList.add("hidden")
}
