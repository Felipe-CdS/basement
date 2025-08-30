function copyToClipboard(e) {
	navigator.clipboard.writeText(e.dataset.value);
	document.getElementById("display").innerHTML = e.dataset.name + ", 생일 축하합니다~";
}
