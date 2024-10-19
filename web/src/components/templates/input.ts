import log from "@/common/log";
import constants from "@/common/constants";
import  loadGamepadKeyboard  from "@comp/templates/gamepad-keyboard";

let interval: ReturnType<typeof setInterval> | null = null;

const updateGamepadState = (): string => {
    const gamepadState: { [key: string]: { x: number; y: number } | number } = {};
    const gamepad = navigator.getGamepads()[0];

    if (gamepad) {
        try {
			// Buttons (triggers, buttons)
			gamepad.buttons.forEach((button, index) => {
				gamepadState[`btn${index}`] = button.pressed ? button.value : 0;
			});
			// Axes (joysticks)
			gamepad.axes.forEach((axis, index) => {
				if (index % 2 === 0) {
					const joystickIndex = index / 2;
					gamepadState[`axes${joystickIndex}`] = {
						x: parseFloat(gamepad.axes[index].toFixed(3)),
						y: parseFloat(gamepad.axes[index + 1].toFixed(3)),
					};
				}
			});
        } catch (error) {
            console.error("Error processing gamepad data:", error);
            throw error;
        }
    } else {
        console.log("No gamepad connected.");
    }
	//console.log(JSON.stringify(gamepadState));
    return JSON.stringify(gamepadState);
};

export const startUpdatingGamepadState = (inputHandler: (state: string) => void): void => {
    if (typeof inputHandler !== 'function') {
        throw new Error("inputHandler must be a function");
    }

    if (interval) {
        stopUpdatingGamepadState();
    }

    interval = setInterval(() => {
        try {
            const state = updateGamepadState();
            inputHandler(state);
        } catch (error) {
            console.error(`Failed to update gamepad due to error: ${error}`);
        }
    }, constants.input.UPDATE_INTERVAL);
};

export const stopUpdatingGamepadState = (): void => {
    if (interval) {
        clearInterval(interval);
        interval = null;
    }
};


let focusedGroups: NodeListOf<HTMLElement> | null = null;
let currentFocusedGroup = 0;
let focusedElements: HTMLElement[][] = [];
let currentFocusedElement = 0;

export const enableGamepadNavigation = () => {
	clearFocus();
	focusedGroups = document.querySelectorAll("[gamepad-focus-group]");
	log.debug(focusedGroups);
	if (!focusedGroups) {
		log.debug("No focus groups found.");
		return;
	}
	focusedGroups.forEach((group) => {
		const elements = Array.from(group.querySelectorAll("[gamepad-focus-element]")) as HTMLElement[];
		if (elements.length > 0) {
			focusedElements.push(elements);
		}
	});
	log.debug(focusedElements);
};
export const clearFocus = () => {
	focusedGroups = null;
	focusedElements=[];
	currentFocusedGroup = 0;
	currentFocusedElement = 0;
	interval = null;
}
const focusElement = () => {
	if (focusedGroups && focusedGroups.length > 0 && focusedElements && focusedElements.length > 0) {
		focusedElements.forEach((elements) => {
			elements.forEach((element) => {
				element.classList.remove("focused");
			});
		});
		let elementInFocus = focusedElements[currentFocusedGroup][currentFocusedElement];
		if (!elementInFocus) {
			currentFocusedElement = 0;
			elementInFocus = focusedElements[currentFocusedGroup][currentFocusedElement];
		}
		elementInFocus.classList.add("focused");
		elementInFocus.focus();
	}
};
let buttonStates = Array(16).fill(false);

export const simpleHandling = () => {
	if (interval) clearInterval(interval);

	interval = setInterval(() => {
		const gamepad = navigator.getGamepads()[0];
		if (!gamepad || !focusedElements || !focusedGroups) {
			return;
		}

		const buttonPressed = (buttonIndex: number) => gamepad.buttons[buttonIndex].pressed;

		for (let i = 0; i < buttonStates.length; i++) {
			const pressed = buttonPressed(i);
			if (pressed && !buttonStates[i]) {
				buttonStates[i] = true;

				switch (i) {
					case 12:
						currentFocusedGroup = currentFocusedGroup > 0 ? currentFocusedGroup - 1 : focusedGroups.length - 1;
						focusElement();
						break;
					case 13:
						currentFocusedGroup = currentFocusedGroup < focusedGroups.length - 1 ? currentFocusedGroup + 1 : 0;
						focusElement();
						break;
					case 0:
						let el = focusedElements[currentFocusedGroup][currentFocusedElement];
						el.click();
						if (el instanceof HTMLInputElement) {
							el.focus();
							loadGamepadKeyboard(el);
						}
						break;
					case 5:
						emulateKeyboard(39);
						break;
					case 4:
						emulateKeyboard(37);
						break;
					case 14:
						if (currentFocusedElement >0 ){
							currentFocusedElement = currentFocusedElement - 1;
						}else if (currentFocusedGroup > 0){
							currentFocusedGroup = currentFocusedGroup - 1;
							currentFocusedElement = focusedElements[currentFocusedGroup].length - 1;
						}else{
							currentFocusedGroup = focusedGroups.length - 1;
							currentFocusedElement = focusedElements[currentFocusedGroup].length - 1;
						}
						focusElement();
						break;
					case 15:
						if (currentFocusedElement < focusedElements[currentFocusedGroup].length - 1){
							currentFocusedElement = currentFocusedElement + 1;
						}else if (currentFocusedGroup < focusedGroups.length - 1){
							currentFocusedGroup = currentFocusedGroup + 1;
							currentFocusedElement = 0;
						}else{
							currentFocusedGroup = 0;
							currentFocusedElement = 0;
						}
						focusElement();
						break;
				}
			} else if (!pressed) {
				buttonStates[i] = false;
			}
		}
	}, 100);
};

const emulateKeyboard = (keyCode: number) => {
	const event = new KeyboardEvent("keydown", { keyCode });
	document.dispatchEvent(event);
	setTimeout(() => {
		const event = new KeyboardEvent("keyup", { keyCode });
		document.dispatchEvent(event);
	}, 100);
};
