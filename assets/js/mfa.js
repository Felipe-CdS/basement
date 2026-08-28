function ato32(stringSecret) {
	const decode32Dict = {
		"A": "00000",
		"B": "00001",
		"C": "00010",
		"D": "00011",
		"E": "00100",
		"F": "00101",
		"G": "00110",
		"H": "00111",
		"I": "01000",
		"J": "01001",
		"K": "01010",
		"L": "01011",
		"M": "01100",
		"N": "01101",
		"O": "01110",
		"P": "01111",
		"Q": "10000",
		"R": "10001",
		"S": "10010",
		"T": "10011",
		"U": "10100",
		"V": "10101",
		"W": "10110",
		"X": "10111",
		"Y": "11000",
		"Z": "11001",
		"2": "11010",
		"3": "11011",
		"4": "11100",
		"5": "11101",
		"6": "11110",
		"7": "11111",
	}

	var stringSecret5bitsConcat = "";

	for(let i = 0; i < stringSecret.length; i++){
		var nextChar = stringSecret[i];
		stringSecret5bitsConcat = `${stringSecret5bitsConcat}${decode32Dict[nextChar]}`
	}

	var nextByte = "";
	var result8bitsArray = [];
	var c = 0;

	for(let i = 0; i < stringSecret5bitsConcat.length; i++){
		if((i != 0 && i % 8 == 0)){
			result8bitsArray.push(parseInt(nextByte, 2));
			nextByte = "";
			c++;
		}

		nextByte = `${nextByte}${stringSecret5bitsConcat[i]}`
	}

	if(nextByte !== "0"){
			result8bitsArray.push(parseInt(nextByte, 2));
	}

	return result8bitsArray;
}

function truncate(buf) {
	const view = new Uint8Array(buf)
	const offset = view[view.length - 1] & 0xf

	const code = (view[offset] &0x7f) << 24 | 
		(view[offset+1] &0xff) << 16 |
		(view[offset+2] &0xff) << 8 |
		(view[offset+3] &0xff) 

	return String(code % 1000000).padStart(6, '0')
}

async function getAuth(secretKey) {
	if(!secretKey) return ''
	const decodedSecret = ato32(secretKey);

	// https://developer.mozilla.org/en-US/docs/Web/API/SubtleCrypto/importKey
  const key = await crypto.subtle.importKey(
    "raw", 														// format
    new Uint8Array(decodedSecret),		// keyData
    { name: "HMAC", hash: "SHA-1" },	// algorithm
    false,														// extractable
    ["sign"]													// keyUsages
  );

	const buffer = new ArrayBuffer(8);
	const unixSeconds = Math.floor((Date.now() /1000) / 30);
	new DataView(buffer).setBigUint64(0, BigInt(unixSeconds), false);

	// https://developer.mozilla.org/en-US/docs/Web/API/SubtleCrypto/sign
  const signedBuffer = await crypto.subtle.sign(
    "HMAC",				// algorithm
    key,					// key
    buffer,				// data
  );

	return truncate(signedBuffer);
}

async function fillTokens(){
	const items = document.querySelectorAll('.token-label');
	const master = document.getElementById("master-input").value;
	const tokenSecrets = []

	if(!master) return;

	try {
		const encryptedEntry = event.detail.xhr.response.trim()
		const decryptedEntry = await decodeTokenSecrets(encryptedEntry, master);
		const entryLines = decryptedEntry.split("\n")

		for(let i = 0; i < entryLines.length; i++){
			if(entryLines[i] == "") continue;

			var lineSplit = entryLines[i].split(":")
			tokenSecrets[lineSplit[0]] = lineSplit[1]
		}
	} catch(e) {
		console.log(e.message)
	}

	for(let i = 0; i < items.length; i++){
		let item = items[i]
		item.innerHTML = "xxx xxx"

		let result = await getAuth(tokenSecrets[item.id])

		if(result != ""){
			item.innerHTML = `${result.slice(0,3)} ${result.slice(3,6)}`
		} 
	};
}


