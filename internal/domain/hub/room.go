package hub

import (
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

type Room struct {
	UUID   string
	Game   string
	Conn   *websocket.Conn
	WebRTC *webrtc.PeerConnection
}

type CreateRoomReq struct {
	UUID string `json:"uuid"`
	Game string `json:"game"`
	Peer *Peer  `json:"peer"`
}
type CreateRoomRes struct {
	Game string `json:"game"`
	UUID string `json:"uuid"`
}

type GetRoomRes struct {
	UUID string `json:"uuid"`
	Game string `json:"game"`
}
