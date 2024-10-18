package input

import "fmt"

type gpadState struct {
	A     int
	B     int
	Y     int
	X     int
	LB    int
	RB    int
	Start int
	LJ    int
	RJ    int

	Dup    int
	Ddown  int
	Dleft  int
	Dright int
	Mode   int
	Select int

	LS struct {
		X float32
		Y float32
	}
	RS struct {
		X float32
		Y float32
	}

	LT float32
	RT float32

	Extra1 int
	Extra2 int
	Extra3 int
	Extra4 int
	Extra5 int
	Extra6 int
}

func selectController(vendor, product string) (interface{}, error) {

	if products, exists := ControllerMap[vendor]; exists {
		if controller, exists := products[product]; exists {
			return controller, nil
		}
		if controller, exists := products["default"]; exists {
			return controller, nil
		}
		return nil, fmt.Errorf("unsupported product: %s", product)
	}
	return ControllerMap["default"]["default"], nil
}

var ControllerMap = map[string]map[string]interface{}{
	"045e": {
		"02fd":    &xbox{},
		"default": &xbox{},
	},
	"2dc8": {
		"3106":    &bitDo{},
		"default": &bitDo{},
	},
	"default": {
		"default": &xbox{},
	},
}

type xbox struct {
	A     int `json:"btn0"`
	B     int `json:"btn1"`
	Y     int `json:"btn3"`
	X     int `json:"btn2"`
	LB    int `json:"btn4"`
	RB    int `json:"btn5"`
	Start int `json:"btn9"`
	LJ    int `json:"btn10"`
	RJ    int `json:"btn11"`

	Dup    int `json:"btn12"`
	Ddown  int `json:"btn13"`
	Dleft  int `json:"btn14"`
	Dright int `json:"btn15"`
	Mode   int `json:"btn16"`
	Select int `json:"btn8"`

	LS struct {
		X float32 `json:"x"`
		Y float32 `json:"y"`
	} `json:"axes0"`
	RS struct {
		X float32 `json:"x"`
		Y float32 `json:"y"`
	} `json:"axes1"`

	LT float32 `json:"btn6"`
	RT float32 `json:"btn7"`
}

func (x *xbox) State() gpadState {
	return gpadState{
		A:      x.A,
		B:      x.B,
		Y:      x.Y,
		X:      x.X,
		LB:     x.LB,
		RB:     x.RB,
		Start:  x.Start,
		LJ:     x.LJ,
		RJ:     x.RJ,
		Dup:    x.Dup,
		Ddown:  x.Ddown,
		Dleft:  x.Dleft,
		Dright: x.Dright,
		Mode:   x.Mode,
		Select: x.Select,
		LS: struct {
			X float32
			Y float32
		}{X: x.LS.X, Y: x.LS.Y},
		RS: struct {
			X float32
			Y float32
		}{X: x.RS.X, Y: x.RS.Y},
		LT: x.LT,
		RT: x.RT,
	}
}

type bitDo struct {
	A     int `json:"btn0"`
	B     int `json:"btn1"`
	Y     int `json:"btn3"`
	X     int `json:"btn2"`
	LB    int `json:"btn4"`
	RB    int `json:"btn5"`
	Start int `json:"btn9"`
	LJ    int `json:"btn10"`
	RJ    int `json:"btn11"`

	Dup    int `json:"btn12"`
	Ddown  int `json:"btn13"`
	Dleft  int `json:"btn14"`
	Dright int `json:"btn15"`
	Mode   int `json:"btn16"`
	Select int `json:"button8"`

	LS struct {
		X float32 `json:"x"`
		Y float32 `json:"y"`
	} `json:"axes0"`
	RS struct {
		X float32 `json:"x"`
		Y float32 `json:"y"`
	} `json:"axes1"`

	LT float32 `json:"trigger1"`
	RT float32 `json:"trigger2"`
}

func (b *bitDo) State() gpadState {
	return gpadState{
		A:      b.A,
		B:      b.B,
		Y:      b.Y,
		X:      b.X,
		LB:     b.LB,
		RB:     b.RB,
		Start:  b.Start,
		LJ:     b.LJ,
		RJ:     b.RJ,
		Dup:    b.Dup,
		Ddown:  b.Ddown,
		Dleft:  b.Dleft,
		Dright: b.Dright,
		Mode:   b.Mode,
		Select: b.Select,
		LS: struct {
			X float32
			Y float32
		}{X: b.LS.X, Y: b.LS.Y},
		RS: struct {
			X float32
			Y float32
		}{X: b.RS.X, Y: b.RS.Y},
		LT: b.LT,
		RT: b.RT,
	}
}
