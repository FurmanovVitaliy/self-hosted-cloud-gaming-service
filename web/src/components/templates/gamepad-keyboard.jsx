const html = (
	<div className='canva'>
		<canvas id='gamepad-keyboard' width='250' height='250'></canvas>
	</div>
);

const canvas = html.querySelector("#gamepad-keyboard");
const ctx = canvas.getContext("2d");
let inputField = undefined;
const centerX = canvas.width / 2;
const centerY = canvas.height / 2;

const radius = canvas.width / 2;
const lettersRadius = radius * 1;
const numbersRadius = radius * 0.7;
const spetialCharsRadius = radius * 0.4;

let letters = "mwfgypbvkjxqzetaoinshrdlcu".split("");
const extraSpecialChars = "\"#$%&'()*+-/:;<=>?[\\]^_{}~";
const quantityOfLettersSectors = letters.length;
const letterSpacingAngle = (3 * Math.PI) / 180; // Gap in radians
const lettersSectorAngle = (2 * Math.PI - quantityOfLettersSectors * letterSpacingAngle) / quantityOfLettersSectors;
let selectedLetterSection = undefined;

const numbers = "6789012345".split("");
const quantityOfNumbersSectors = numbers.length;
const numbersSpacingAngle = (10 * Math.PI) / 180;
const numbersSectorAngle = (2 * Math.PI - quantityOfNumbersSectors * numbersSpacingAngle) / quantityOfNumbersSectors;
let selectedNumberSection = undefined;

const spetialChars = "!@. ,".split("");
const quantityOfSpetialCharsSectors = spetialChars.length;
const spetialCharsSpacingAngle = (20 * Math.PI) / 180;
const spetialCharsSectorAngle =
	(2 * Math.PI - quantityOfSpetialCharsSectors * spetialCharsSpacingAngle) / quantityOfSpetialCharsSectors;
let selectedSpetialCharsSection = undefined;

const sectorTypes = ["letters", "numbers", "spetialChars"];

const smoothingFactor = 0.5;
const sensitivity = 0.01;
const deadzone = 0.3;

let previousX = 0;
let previousY = 0;
let previousInputButtonState = false;
let previousDeleteButtonState = false;

let sectorSwicher = { previousState: false, sector: sectorTypes[0] };
let capslock = { previousState: false, active: false };
let letterChange = { previousState: false, active: false };

function loadGamepadKeyboard(element) {
	document.querySelector("#app").appendChild(html);
	inputField = element;
	handleGamepadInput();
}

function closeGamepadKeyboard() {
	html.remove();
}

function drawLetterSectors() {
	ctx.clearRect(0, 0, canvas.width, canvas.height);
	ctx.save();

	for (let i = 0; i < quantityOfLettersSectors; i++) {
		const startAngle = i * (lettersSectorAngle + letterSpacingAngle);
		const endAngle = startAngle + lettersSectorAngle;
		drawSector(centerX, centerY, lettersRadius, startAngle, endAngle, i === selectedLetterSection ? "#ffa" : "#fff");
		drawText(centerX, centerY, lettersRadius, letters[i], startAngle, lettersSectorAngle);
	}
	clipInnerCircle(lettersRadius);
	ctx.restore();
}

function drawNumberSectors() {
	ctx.save();
	for (let i = 0; i < quantityOfNumbersSectors; i++) {
		const startAngle = i * (numbersSectorAngle + numbersSpacingAngle);
		const endAngle = startAngle + numbersSectorAngle;
		drawSector(centerX, centerY, numbersRadius, startAngle, endAngle, i === selectedNumberSection ? "#afa" : "#fff");
		drawText(centerX, centerY, numbersRadius, numbers[i], startAngle, numbersSectorAngle);
	}
	clipInnerCircle(numbersRadius);
	ctx.restore();
}

function drarSpetialCharsSectors() {
	ctx.save();
	for (let i = 0; i < quantityOfSpetialCharsSectors; i++) {
		const startAngle = i * (spetialCharsSectorAngle + spetialCharsSpacingAngle);
		const endAngle = startAngle + spetialCharsSectorAngle;
		drawSector(centerX, centerY, spetialCharsRadius, startAngle, endAngle, i === selectedSpetialCharsSection ? "#afa" : "#fff");
		drawText(centerX, centerY, spetialCharsRadius, spetialChars[i], startAngle, spetialCharsSectorAngle);
	}
	clipInnerCircle(spetialCharsRadius);
	ctx.restore();
}

function drawSector(x, y, radius, startAngle, endAngle, color) {
	ctx.beginPath();
	ctx.moveTo(x, y);
	ctx.arc(x, y, radius, startAngle, endAngle);
	ctx.closePath();
	ctx.fillStyle = color;
	ctx.fill();
}

function drawText(x, y, radius, text, startAngle, sectorAngle) {
	ctx.save();
	ctx.translate(x, y);
	ctx.rotate(startAngle + sectorAngle / 2);
	ctx.translate(radius * 0.9, 0);
	ctx.rotate(-startAngle - sectorAngle / 2);
	ctx.fillStyle = "black";
	ctx.font = "500 1em 'Poppins'";
	ctx.fillText(text, -ctx.measureText(text).width / 2, 4);
	ctx.restore();
}

function clipInnerCircle(radius) {
	const innerRadius = radius * 0.8;
	ctx.beginPath();
	ctx.arc(centerX, centerY, innerRadius, 0, 2 * Math.PI);
	ctx.closePath();
	ctx.clip();
	ctx.clearRect(centerX - innerRadius, centerY - innerRadius, innerRadius * 2, innerRadius * 2);
}

function updateSelection(axisX, axisY, sectionType) {
	const angle = Math.atan2(axisY * -1, axisX * -1);
	let sectorAngle = 0;
	let sector = 0;

	switch (sectionType) {
		case "letters":
			sectorAngle = (2 * Math.PI) / quantityOfLettersSectors;
			sector = Math.floor((angle + Math.PI) / sectorAngle);
			selectedLetterSection = (sector + quantityOfLettersSectors) % quantityOfLettersSectors;
			selectedNumberSection = undefined;
			selectedSpetialCharsSection = undefined;
			break;

		case "numbers":
			sectorAngle = (2 * Math.PI) / quantityOfNumbersSectors;
			sector = Math.floor((angle + Math.PI) / sectorAngle);
			selectedNumberSection = (sector + quantityOfNumbersSectors) % quantityOfNumbersSectors;
			selectedLetterSection = undefined;
			selectedSpetialCharsSection = undefined;
			break;
		case "spetialChars":
			sectorAngle = (2 * Math.PI) / quantityOfSpetialCharsSectors;
			sector = Math.floor((angle + Math.PI) / sectorAngle);
			selectedSpetialCharsSection = (sector + quantityOfSpetialCharsSectors) % quantityOfSpetialCharsSectors;
			selectedLetterSection = undefined;
			selectedNumberSection = undefined;
	}
}

function smoothInput(currentValue, previousValue, factor) {
	return previousValue + (currentValue - previousValue) * factor;
}

function applyDeadzone(value, deadzone) {
	if (Math.abs(value) < deadzone) {
		return 0;
	} else {
		return (value - Math.sign(value) * deadzone) / (1 - deadzone);
	}
}

function handleGamepadInput() {
	console.log("handleGamepadInput");
	const gamepads = navigator.getGamepads();
	if (!gamepads[0]) return;

	const gamepad = gamepads[0];

	let rawX = applyDeadzone(gamepad.axes[0], deadzone);
	let rawY = applyDeadzone(gamepad.axes[1], deadzone);

	const smoothX = smoothInput(rawX, previousX, smoothingFactor);
	const smoothY = smoothInput(rawY, previousY, smoothingFactor);

	previousX = smoothX;
	previousY = smoothY;

	const adjustedX = smoothX * sensitivity;
	const adjustedY = smoothY * sensitivity;

	if (gamepad.buttons[5].pressed && !previousDeleteButtonState) {
		inputField.value = inputField.value.slice(0, -1);
	}

	if (gamepad.buttons[4].pressed && !sectorSwicher.previousState) {
		sectorSwicher.sector = sectorTypes[(sectorTypes.indexOf(sectorSwicher.sector) + 1) % sectorTypes.length];
		console.log(sectorSwicher.sector);
	}

	if (gamepad.buttons[2].pressed && !letterChange.previousState) {
		letterChange.active = !letterChange.active;
		letters = letterChange.active
			? (letters = extraSpecialChars.split(""))
			: (letters = "abcdefghijklmnopqrstuvwxyz".split(""));
	}

	if (gamepad.buttons[3].pressed && !capslock.previousState) {
		capslock.active = !capslock.active;
		letters = capslock.active ? letters.map((letter) => letter.toUpperCase()) : letters.map((letter) => letter.toLowerCase());
	}
	if (gamepad.buttons[1].pressed) {
		closeGamepadKeyboard();
		return;
	}

	updateSelection(adjustedX, adjustedY, sectorSwicher.sector);

	drawLetterSectors();
	drawNumberSectors();
	drarSpetialCharsSectors();

	if (gamepad.buttons[0].pressed && !previousInputButtonState) {
		inputField.value +=
			sectorSwicher.sector === "letters"
				? letters[selectedLetterSection]
				: sectorSwicher.sector === "numbers"
				? numbers[selectedNumberSection]
				: spetialChars[selectedSpetialCharsSection];
		inputField.dispatchEvent(new Event("input"));
	}

	letterChange.previousState = gamepad.buttons[2].pressed;
	sectorSwicher.previousState = gamepad.buttons[4].pressed;
	capslock.previousState = gamepad.buttons[3].pressed;
	previousInputButtonState = gamepad.buttons[0].pressed;
	previousDeleteButtonState = gamepad.buttons[5].pressed;
	requestAnimationFrame(handleGamepadInput);
}

export default loadGamepadKeyboard;
