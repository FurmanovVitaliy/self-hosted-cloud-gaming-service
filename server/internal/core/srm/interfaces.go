package srm

type InputDevice interface {
	GetEvdevPath() (string, error)
	HandleInput(data []byte)
	Close() error
}

type UDPReader interface {
	Read() ([]byte, error)
	GetPort() int
	Close() error
}

type Display interface {
	GetDisplayNumper() string
	GetPlaneID() int
	Close()
}
