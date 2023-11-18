package hub

import (
	"cloud/pkg/network/websocket"
	"cloud/pkg/webrtc"

	"github.com/bendahl/uinput"
)

type Room struct {
	UUID      string
	Game      string
	Peer      *webrtc.Peer
	Websocket *websocket.Websocket
	Gamepad   uinput.Gamepad
}

type CreateRoomReq struct {
	UUID string `json:"uuid"`
	Game string `json:"game"`
}
type CreateRoomRes struct {
	Game string `json:"game"`
	UUID string `json:"uuid"`
}

type GetRoomRes struct {
	UUID string `json:"uuid"`
	Game string `json:"game"`
}
