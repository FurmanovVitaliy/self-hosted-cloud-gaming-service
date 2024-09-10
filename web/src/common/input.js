let interval;
export const startUpdatingGamepadState = (inputHandler) => {
  // Make sure no previous interval is running
  if (interval) clearInterval(interval);

  interval = setInterval(() => {
    let state;
    try {
      state = updateGamepadState();
    } catch (error) {
      console.error('Failed to update gamepad state:', error);
      return; // Exit the interval callback if there's an error
    }
    return inputHandler(state);
  }, 25); // Very short interval, consider increasing if performance issues occur
};

export const stopUpdatingGamepadState = () => {
  clearInterval(interval);
  interval = null; // Clear the reference to avoid potential leaks or errors
};

export const updateGamepadState = () => {
  let gamepadState = {};
  const gamepad = navigator.getGamepads()[0];
  if (gamepad) {
    try {
      for (let j = 0; j < gamepad.buttons.length; j++) {
        gamepadState[`button${j}`] = gamepad.buttons[j].pressed ? 1 : 0;
      }

      for (let j = 0; j < gamepad.axes.length; j += 2) { // Ensuring step increments by 2
        gamepadState[`joystick${j / 2}`] = {
          x: gamepad.axes[j].toFixed(3),
          y: gamepad.axes[j + 1].toFixed(3),
        };
      }

      for (let j = 0; j < gamepad.buttons.length; j++) {
        if (gamepad.buttons[j].value) {
          gamepadState[`trigger${j}`] = gamepad.buttons[j].value.toFixed(3);
        }
      }
    } catch (error) {
      console.error('Error processing gamepad data:', error);
      throw error; // Re-throw to be handled in setInterval
    }
  } else {
    console.log('No gamepad connected.');
  }
  console.log(gamepadState);
  return JSON.stringify(gamepadState);

};

export default startUpdatingGamepadState;