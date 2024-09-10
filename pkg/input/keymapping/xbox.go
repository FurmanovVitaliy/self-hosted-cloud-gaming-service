package keymapping

type XboxGpadInput struct {
	A     int `json:"button0"`
	B     int `json:"button1"`
	Y     int `json:"button3"`
	X     int `json:"button2"`
	LB    int `json:"button4"`
	RB    int `json:"button5"`
	Start int `json:"button9"`
	LJ    int `json:"button10"`
	RJ    int `json:"button11"`

	Dup    int `json:"button12"`
	Ddown  int `json:"button13"`
	Dleft  int `json:"button14"`
	Dright int `json:"button15"`
	Mode   int `json:"button16"`
	Select int `json:"button8"`

	LS struct {
		X string `json:"x"`
		Y string `json:"y"`
	} `json:"joystick0"`
	RS struct {
		X string `json:"x"`
		Y string `json:"y"`
	} `json:"joystick1"`

	LT string `json:"trigger6"`
	RT string `json:"trigger7"`
}
