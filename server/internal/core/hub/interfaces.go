package hub

type MessageHandler interface {
	RegisterChannel(tag string, channel interface{}) error
	Write(tag string, msg interface{}) error
	Stop()
}

type WrtcStreamer interface {
	Init(OnIceCandidate func(string)) (offer string, err error)
	SetAnswer(data string) error
	AddCandidate(data string) error
	SendData(data []byte) error
	SendVideo(data []byte) error
	SendAudio(data []byte) error
	Stop()
}

type UDPReader interface {
	Read() ([]byte, error)
	Close() error
}

type VM interface {
	Start() error
	Stop()
}
