const device = {};

export const collectDeviceInfo = () => {
    device['display'] = {
        height: window.screen.height,
        width: window.screen.width
    }
    device['control'] = {
        gamepad: ('getGamepads' in navigator) ? true : false,
        keyboard: ('keyboard' in navigator) ? true : false,
        mouse: ('mouse' in navigator) ? true : false,
        touch: ('touch' in navigator) ? true : false,
    }

    return device;
}


export default collectDeviceInfo;

