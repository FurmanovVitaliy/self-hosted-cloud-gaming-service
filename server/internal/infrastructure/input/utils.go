package input

import (
	"fmt"

	"github.com/FurmanovVitaliy/uinput"
	evdev "github.com/gvalkov/golang-evdev"
)

func findEvdevPath(name string) (string, error) {
	var path string
	devices, err := evdev.ListInputDevices()
	if err != nil {
		return "", err
	}
	for _, device := range devices {
		if device.Name == name {
			path = device.Fn
			return path, nil
		}
	}
	return "", fmt.Errorf("device not found: %s", name)
}

func handleButton(device uinput.Gamepad, button int, pressed int) {
	if pressed == 1 {
		device.ButtonDown(button)
	} else {
		device.ButtonUp(button)
	}
}
