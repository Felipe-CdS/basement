async function decodeTokenSecrets(encryptedEntry, masterPassword) {
	const encoder = new TextEncoder();
	const decoder = new TextDecoder("utf-8");

	const saltLen = 32
	const ivLen = 24
	const tagLen = 32

	const entryParts = {
		salt: encryptedEntry.slice(0, saltLen),
		iv: encryptedEntry.slice(saltLen, saltLen + ivLen),
		data: encryptedEntry.slice(saltLen + ivLen),
	};

	const keyMaterial = await window.crypto.subtle.importKey(
		"raw",
		encoder.encode(masterPassword),
		"PBKDF2",
		false,
		["deriveBits", "deriveKey"],
	);

	const key = await window.crypto.subtle.deriveKey(
		{ 
			name: "PBKDF2", 
			salt: Uint8Array.fromHex(entryParts.salt),
			iterations: 600000, 
			hash: "SHA-256" 
		},
		keyMaterial,
		{ name: "AES-GCM", length: 256 },
		true,
		["decrypt"],
	);

	try {

		const arrayBufferResult = await window.crypto.subtle.decrypt(
			{ name: "AES-GCM", iv: Uint8Array.fromHex(entryParts.iv) }, 
			key, 
			Uint8Array.fromHex(entryParts.data),
		);

		return JSON.parse(decoder.decode(arrayBufferResult))
	} catch(e){
		console.log("decrypt fail")
		return {}
	}
}
