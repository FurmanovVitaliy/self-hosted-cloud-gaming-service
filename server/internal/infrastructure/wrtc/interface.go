package wrtc

type WebRTC interface {
	Init(OnIceCandidate func(string)) (offer string, err error)
	SetAnswer(data string) (err error)
	AddCandidate(data string) (err error)
	SendData(data []byte) error
	SendVideo(data []byte) error
	SendAudio(data []byte) error
	Stop()
}
