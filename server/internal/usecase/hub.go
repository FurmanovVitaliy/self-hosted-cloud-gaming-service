package usecase

import (
	"context"
	"net/http"
	"time"

	"github.com/FurmanovVitaliy/pixel-cloud/internal/infrastructure/broker"
	"github.com/FurmanovVitaliy/pixel-cloud/pkg/errors"
	"github.com/gorilla/websocket"
)

// dto
type CreateRoomReq struct {
	UUID   string `json:"uuid"`
	GameID string `json:"game_id"`
}
type CreateRoomRes struct {
	UUID     string `json:"uuid"`
	GameName string `json:"game_name"`
}

type GetRoomRes struct {
	UUID    string   `json:"uuid"`
	Game    string   `json:"game"`
	Players []string `json:"players"`
}

type RoomStatusRes struct {
	Exist          bool   `json:"exist"`
	Bysy           bool   `json:"bysy"`
	Game           string `json:"game"`
	PlayerQuantity int    `json:"player_quantity"`
}

type JoinRoomRes struct {
	Display struct {
		Height int `json:"height,omitempty"`
		Width  int `json:"width,omitempty"`
	} `json:"display,omitempty"`
	Control struct {
		Type      string `json:"type,omitempty"`
		VendorID  string `json:"vendorID,omitempty"`
		ProductID string `json:"productID,omitempty"`
	} `json:"control,omitempty"`
}

// custom error
var (
	RoomNotExist      = errors.New(http.StatusNotFound, "HS", "00000", "Room not exist")
	ErrInvalidRequest = errors.New(http.StatusBadRequest, "HS", "00001", "invalid request")
	ErrRoomJoin       = errors.New(http.StatusInternalServerError, "HS", "00002", "Failed to join room")
)

func (uc *UseCase) GetRooms() (rooms []GetRoomRes, err error) {
	rms, err := uc.hubService.GetRooms()
	if err != nil {
		return []GetRoomRes{}, err
	}
	rooms = make([]GetRoomRes, 0, len(rms))
	for _, r := range rms {
		var players []string
		for _, p := range r.Workers {
			players = append(players, p.Username)
		}
		game, err := uc.gameService.GetOneById(context.Background(), r.GameID)
		if err != nil {
			return []GetRoomRes{}, err
		}
		rooms = append(rooms, GetRoomRes{
			UUID:    r.UUID,
			Game:    game.Name,
			Players: players,
		})
	}
	return rooms, nil
}

func (uc *UseCase) CreateRoom(req CreateRoomReq) (res CreateRoomRes, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = uc.hubService.CreateRoom(req.UUID, req.GameID)
	if err != nil {
		return CreateRoomRes{}, errors.Wrap(err, "HS", "00003", "Failed to create room: "+err.Error())
	}
	game, err := uc.gameService.GetOneById(ctx, req.GameID)
	if err != nil {
		return CreateRoomRes{}, errors.Wrap(err, "HS", "00003", "Failed to create room: "+err.Error())
	}
	return CreateRoomRes{UUID: req.UUID, GameName: game.Name}, nil
}

// ! very strict logic
func (uc *UseCase) JoinRoom(uuid, username string, ws *websocket.Conn, deviceInfo JoinRoomRes) (err error) {
	_, err = uc.hubService.GetRoom(uuid)
	if err != nil {
		return RoomNotExist
	}
	if username == "" {
		return ErrInvalidRequest
	}
	//TODO: get ip from config or from srm service if will be connect  more than one server-machine
	var ip = "127.0.0.1"
	//!!check resources first
	inputDevice, virtualDisplay, aReader, vReader, pulsePort, userDiskSpace, err := uc.allocateAllResources(ip, username, deviceInfo)
	if err != nil {
		uc.logger.Error(err)
		return ErrRoomJoin
	}
	//! create massage bus/habdler (general for worker, room, player)
	br := broker.New(ws)
	//!!create streamer
	streamer, err := uc.streamerService.CreateRTC(inputDevice.HandleInput)
	if err != nil {
		uc.releaseAllResources(inputDevice, virtualDisplay, aReader, vReader)
		uc.logger.Error(err)
		return ErrRoomJoin
	}
	//!!create vm
	vm, err := uc.createVmWithSrmresources(ip, uuid, username, userDiskSpace, inputDevice, virtualDisplay, vReader, aReader, pulsePort)
	if err != nil {
		uc.releaseAllResources(inputDevice, virtualDisplay, aReader, vReader)
		uc.logger.Error(err)
		return ErrRoomJoin
	}
	//! add player in room
	if err = uc.hubService.CreateWorker(uuid, username, br, streamer, aReader, vReader, vm); err != nil {
		br.Stop()
		uc.releaseAllResources(inputDevice, virtualDisplay, aReader, vReader)
		return ErrRoomJoin
	}
	//! start listen client messages and send to worker (worker will send to streamer and vm)
	go func() {
		defer func() {
			uc.logger.Warn("Stop listen client messages")
			uc.releaseAllResources(inputDevice, virtualDisplay, aReader, vReader)
		}()
		br.Read()
	}()

	return nil

}
