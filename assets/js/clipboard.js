function copyToClipboard(e) {
	navigator.clipboard.writeText(e.dataset.value);
}

async function decToHex() {
	const v = await navigator.clipboard.readText();
	const r = Number(v).toString(16).toUpperCase();
	navigator.clipboard.writeText(r);
}
