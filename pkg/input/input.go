package input

import "github.com/bendahl/uinput"

type VirtualGamepad struct {
	Path     string
	Vendor   uint16
	Produkt  uint16
	Username []byte
}

func (v *VirtualGamepad) Create() {
	var vg, _ = uinput.CreateGamepad("/dev/uinput", []byte("Virtual Gamepad"), 1, 1)
	_ = vg
}

func (v *VirtualGamepad) New(path string, vendor uint16, produkt uint16, username []byte) (*VirtualGamepad, error) {
	return &VirtualGamepad{
		Path:     path,
		Vendor:   vendor,
		Produkt:  produkt,
		Username: username,
	}, nil
}
