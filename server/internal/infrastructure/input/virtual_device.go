package input

import (
	"encoding/json"
	"fmt"

	"github.com/FurmanovVitaliy/uinput"
)

// TODO: own virtual device realization to get out this abstraction hell

type DeviceType string

const (
	Keyboard DeviceType = "keyboard"
	Mouse    DeviceType = "mouse"
	Gamepad  DeviceType = "gamepad"
)

type Device interface {
	Close() error
}
type Keymapping interface {
	State() gpadState
}

type InputDevice struct {
	evdevPath  string
	username   string
	vendor     string
	product    string
	device     Device
	Keymapping Keymapping
}

func NewInputDevice(username, vendor, product string, deviceType DeviceType) (device *InputDevice, err error) {
	switch deviceType {
	case Keyboard:
		return nil, fmt.Errorf("keyboard device is not supported yet")
	case Mouse:
		return
	case Gamepad:
		gamepad, err := uinput.CreateGamepad("/dev/uinput", []byte(username+"-gamepad"), 045, 955)
		if err != nil {
			return nil, err
		}
		keymap, err := selectController(vendor, product)
		if err != nil {
			return nil, err
		}
		path, err := findEvdevPath(username + "-gamepad")
		if err != nil {
			return nil, err
		}
		device = &InputDevice{
			evdevPath:  path,
			vendor:     vendor,
			device:     gamepad,
			product:    product,
			username:   username,
			Keymapping: keymap.(Keymapping),
		}
	default:
		err = fmt.Errorf("unsupported device type: %s", deviceType)
	}
	return
}

func (id *InputDevice) HandleInput(data []byte) {
	if _, ok := id.device.(uinput.Gamepad); ok {
		err := json.Unmarshal(data, &id.Keymapping)
		if err != nil {
			fmt.Println("error unmarshalling data: ", err)
		}
		state := id.Keymapping.State()
		id.device.(uinput.Gamepad).LeftStickMove(state.LS.X, state.LS.Y)
		id.device.(uinput.Gamepad).RightStickMove(state.RS.X, state.RS.Y)
		id.device.(uinput.Gamepad).LeftTriggerPress(state.LT)
		id.device.(uinput.Gamepad).RightTriggerPress(state.RT)
		handleButton(id.device.(uinput.Gamepad), uinput.ButtonSouth, state.A)
		handleButton(id.device.(uinput.Gamepad), uinput.ButtonEast, state.B)
		handleButton(id.device.(uinput.Gamepad), uinput.ButtonWest, state.X)
		handleButton(id.device.(uinput.Gamepad), uinput.ButtonNorth, state.Y)
		handleButton(id.device.(uinput.Gamepad), uinput.ButtonThumbLeft, state.LB)
		handleButton(id.device.(uinput.Gamepad), uinput.ButtonThumbRight, state.RB)
		handleButton(id.device.(uinput.Gamepad), uinput.ButtonStart, state.Start)
		handleButton(id.device.(uinput.Gamepad), uinput.ButtonSelect, state.Select)
		handleButton(id.device.(uinput.Gamepad), uinput.ButtonMode, state.Mode)
		handleButton(id.device.(uinput.Gamepad), uinput.ButtonDpadUp, state.Dup)
		handleButton(id.device.(uinput.Gamepad), uinput.ButtonDpadDown, state.Ddown)
		handleButton(id.device.(uinput.Gamepad), uinput.ButtonDpadLeft, state.Dleft)
		handleButton(id.device.(uinput.Gamepad), uinput.ButtonDpadRight, state.Dright)
	}

}
func (id *InputDevice) GetEvdevPath() (string, error) {
	return findEvdevPath(id.username + "-gamepad")
}

func (id *InputDevice) Close() error {
	return id.device.Close()

}
