package keymapping

type XboxGpadInput struct {
	A     int `json:"button0"`
	B     int `json:"button1"`
	X     int `json:"button3"`
	Y     int `json:"button4"`
	LB    int `json:"button6"`
	RB    int `json:"button7"`
	Start int `json:"button11"`
	Main  int `json:"button12"`
	LJ    int `json:"button13"`
	RJ    int `json:"button14"`

	LS struct {
		X string `json:"x"`
		Y string `json:"y"`
	} `json:"joystick0"`
	RS struct {
		X string `json:"x"`
		Y string `json:"y"`
	} `json:"joystick1"`
	Triger struct {
		RT float32 `json:"x"`
		LT float32 `json:"y"`
	} `json:"joystick2"`
	Dpad struct {
		X float32 `json:"x"`
		Y float32 `json:"y"`
	} `json:"joystick3"`
}

const (
	ButtonA     = 0x130 // A / X
	ButtonB     = 0x131 // X / Квадрат
	ButtonX     = 0x133 // Y / Треугольник
	ButtonY     = 0x134 // B / Круг
	ButtonLB    = 0x136 // Левый верхний бампер (L1)
	ButtonRB    = 0x137 // Правый верхний бампер (R1)
	ButtonStart = 0x13b // Старт
	ButtonMain  = 0x13c // Главная кнопка
	ButtonLJ    = 0x13d // Левый стик (L3)
	ButtonRJ    = 0x13e // Правый стик (R3)

	ButtonSelect       = 0x13a
	ButtonBumperLeft   = 0x136 // Левый бампер (L1)
	ButtonBumperRight  = 0x137 // Правый бампер (R1)
	ButtonTriggerLeft  = 0x138 // Левый триггер (L2)
	ButtonTriggerRight = 0x139 // Правый триггер (R2)
	ButtonThumbLeft    = 0x13d // Левый стик (L3)
	ButtonThumbRight   = 0x13e // Правый стик (R3)

	ButtonDpadUp    = 0x220 // Дпад вверх
	ButtonDpadDown  = 0x221 // Дпад вниз
	ButtonDpadLeft  = 0x222 // Дпад влево
	ButtonDpadRight = 0x223 // Дпад вправо

	ButtonMode = 0x13c // Это специальная кнопка, обычно с логотипом Xbox или Playstation
)
