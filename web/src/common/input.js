let interval;
export const startUpdatingGamepadState = (inputHandler) => {
  interval = setInterval(() => {
  return inputHandler(updateGamepadState());
  }, 2); // Интервал в миллисекундах 
};

export const stopUpdatingGamepadState = () => {
  clearInterval(interval); 
};
export const updateGamepadState = () => {
  let gamepadState = {};
  const gamepad = navigator.getGamepads()[0];
  if (gamepad) {
    for (let j = 0; j < gamepad.buttons.length; j++) {
      gamepadState[`button${j}`] = gamepad.buttons[j].pressed ? 1 : 0;
    }

    for (let j = 0; j < gamepad.axes.length; j++) {
      if (j % 2 === 0) {
        gamepadState[`joystick${j / 2}`] = {
          x: gamepad.axes[j].toFixed(3),
          y: gamepad.axes[j + 1].toFixed(3),
        };
      }
    }

    for (let j = 0; j < gamepad.buttons.length; j++) {
      if (gamepad.buttons[j].value) {
        gamepadState[`trigger${j}`] = gamepad.buttons[j].value.toFixed(3);
      }
    }
  }
  return JSON.stringify(gamepadState);};

export default startUpdatingGamepadState;
