package input

import (
	"cloud/pkg/input/keymapping"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/FurmanovVitaliy/uinput"

	evdev "github.com/gvalkov/golang-evdev"
)

type VirtualGamepad struct {
	Path     string
	Username string
	Device   uinput.Gamepad
}

func CreateGamepad(username string) (*VirtualGamepad, error) {
	vg := &VirtualGamepad{
		Username: username,
	}
	var err error
	fmt.Println(vg.Username)
	vg.Device, err = uinput.CreateGamepad("/dev/uinput", []byte(vg.Username+"-gamepad"), 045, 955)
	if err != nil {
		return nil, err
	}

	devices, err := evdev.ListInputDevices()
	if err != nil {
		fmt.Print("error listing input devices: ", err)
		return nil, err
	}

	for _, dev := range devices {
		if dev.Name == vg.Username+"-gamepad" {
			vg.Path = dev.Fn
		}
	}
	fmt.Println("Виртуальный геймпад создан:", vg.Username+"-gamepad")
	fmt.Println(vg.Path)
	fmt.Print(vg.Device)
	return vg, nil
}

func (vg *VirtualGamepad) Close() {
	vg.Device.Close()
}

func (vg *VirtualGamepad) HandleInput(data []byte, keymap keymapping.XboxGpadInput) {
	err := json.Unmarshal(data, &keymap)
	if err != nil {
		fmt.Println("error unmarshalling data: ", err)
	}
	//fmt.Println(keymap)

	lsX, _ := strconv.ParseFloat(keymap.LS.X, 32)
	lsY, _ := strconv.ParseFloat(keymap.LS.Y, 32)
	vg.Device.LeftStickMove(float32(lsX), float32(lsY))

	rsX, _ := strconv.ParseFloat(keymap.RS.X, 32)
	rsY, _ := strconv.ParseFloat(keymap.RS.Y, 32)
	vg.Device.RightStickMove(float32(rsX), float32(rsY))

	lt, err := strconv.ParseFloat(keymap.LT, 32)
	if err == nil {
		vg.Device.LeftTriggerPress(float32(lt))
	} else {
		vg.Device.LeftTriggerPress(0)
	}

	rt, err := strconv.ParseFloat(keymap.RT, 32)
	if err == nil {
		vg.Device.RightTriggerPress(float32(rt))
	} else {
		vg.Device.RightTriggerPress(0)
	}

	if keymap.A == 1 {
		vg.Device.ButtonDown(uinput.ButtonSouth)
	} else {
		vg.Device.ButtonUp(uinput.ButtonSouth)
	}
	if keymap.B == 1 {
		vg.Device.ButtonDown(uinput.ButtonEast)
	} else {
		vg.Device.ButtonUp(uinput.ButtonEast)
	}
	if keymap.X == 1 {
		vg.Device.ButtonDown(uinput.ButtonNorth)
	} else {
		vg.Device.ButtonUp(uinput.ButtonNorth)
	}
	if keymap.Y == 1 {
		vg.Device.ButtonDown(uinput.ButtonWest)
	} else {
		vg.Device.ButtonUp(uinput.ButtonWest)
	}
	if keymap.LB == 1 {
		vg.Device.ButtonDown(uinput.ButtonBumperLeft)
	} else {
		vg.Device.ButtonUp(uinput.ButtonBumperLeft)
	}
	if keymap.RB == 1 {
		vg.Device.ButtonDown(uinput.ButtonBumperRight)
	} else {
		vg.Device.ButtonUp(uinput.ButtonBumperRight)
	}
	if keymap.Start == 1 {
		vg.Device.ButtonDown(uinput.ButtonStart)
	} else {
		vg.Device.ButtonUp(uinput.ButtonStart)
	}
	if keymap.LJ == 1 {
		vg.Device.ButtonDown(uinput.ButtonThumbLeft)
	} else {
		vg.Device.ButtonUp(uinput.ButtonThumbLeft)
	}
	if keymap.RJ == 1 {
		vg.Device.ButtonDown(uinput.ButtonThumbRight)
	} else {
		vg.Device.ButtonUp(uinput.ButtonThumbRight)
	}
	if keymap.Dup == 1 {
		vg.Device.ButtonDown(uinput.ButtonDpadUp)
	} else {
		vg.Device.ButtonUp(uinput.ButtonDpadUp)
	}
	if keymap.Ddown == 1 {
		vg.Device.ButtonDown(uinput.ButtonDpadDown)
	} else {
		vg.Device.ButtonUp(uinput.ButtonDpadDown)
	}
	if keymap.Dleft == 1 {
		vg.Device.ButtonDown(uinput.ButtonDpadLeft)
	} else {
		vg.Device.ButtonUp(uinput.ButtonDpadLeft)
	}
	if keymap.Dright == 1 {
		vg.Device.ButtonDown(uinput.ButtonDpadRight)
	} else {
		vg.Device.ButtonUp(uinput.ButtonDpadRight)
	}
	if keymap.Mode == 1 {
		vg.Device.ButtonDown(uinput.ButtonMode)
	} else {
		vg.Device.ButtonUp(uinput.ButtonMode)
	}
	if keymap.Select == 1 {
		vg.Device.ButtonDown(uinput.ButtonSelect)
	} else {
		vg.Device.ButtonUp(uinput.ButtonSelect)
	}

}
