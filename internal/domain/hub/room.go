package hub

import (
	"cloud/internal/domain/games"

	"github.com/bendahl/uinput"
)

type Room struct {
	UUID    string
	Game    games.Game
	Workers map[string]*Worker
	Gamepad uinput.Gamepad
}

type CreateRoomReq struct {
	UUID   string `json:"uuid"`
	GameID string `json:"game_id"`
}
type CreateRoomRes struct {
	UUID     string `json:"uuid"`
	GameName string `json:"game_name"`
}

type GetRoomRes struct {
	UUID string `json:"uuid"`
	Game string `json:"game"`
}
