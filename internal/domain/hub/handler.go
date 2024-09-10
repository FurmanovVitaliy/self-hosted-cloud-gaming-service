package hub

import (
	"cloud/internal/domain/games"
	"cloud/internal/messages"
	"cloud/pkg/logger"
	"cloud/pkg/network/webrtc"
	"context"
	"fmt"

	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) (b bool) { return true },
}

const (
	roomsURL  = "/room"
	createURL = "/room/create"
	roomURL   = "/room/join/:uuid"
	roomState = "/room/state/:uuid"
)

type handler struct {
	hub         *Hub
	logger      *logger.Logger
	gameService games.Service
}

func Handler(hub *Hub, logger *logger.Logger, gameService games.Service) *handler {
	return &handler{
		hub:         hub,
		logger:      logger,
		gameService: gameService,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, createURL, h.createRoom)
	router.HandlerFunc(http.MethodGet, roomsURL, h.getRoomList)
	router.GET(roomState, h.getRoomState)
	router.GET(roomURL, h.joinRoom)
}

func (h *handler) getRoomList(w http.ResponseWriter, r *http.Request) {
	res := make([]GetRoomRes, 0)
	for _, room := range h.hub.Rooms {
		res = append(res, GetRoomRes{
			UUID: room.UUID,
			Game: room.Game.Name,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *handler) createRoom(w http.ResponseWriter, r *http.Request) {
	h.hub.Mutex.Lock()
	var req CreateRoomReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	game, err := h.gameService.GetOne(req.GameID)
	if err != nil {
		h.logger.Errorf("Game with id '%s' not found while room (uuid %s) creating : %s", req.GameID, req.UUID, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.hub.Rooms[req.UUID] = &Room{
		UUID:    req.UUID,
		Game:    game,
		Workers: make(map[string]*Worker),
	}

	res := CreateRoomRes{
		UUID:     req.UUID,
		GameName: game.Name,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
	h.logger.Info("room created:", h.hub.Rooms)
	h.hub.Mutex.Unlock()
}

func (h *handler) getRoomState(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var res RoomStatusRes
	roomUUID := p.ByName("uuid")
	username := r.URL.Query().Get("username")
	if roomUUID == "" || username == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, ok := h.hub.Rooms[roomUUID]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	res.Exist = true
	res.Game = h.hub.Rooms[roomUUID].Game.Name
	res.PlayerQuantity = len(h.hub.Rooms[roomUUID].Workers)
	if res.PlayerQuantity > 0 {
		res.Bysy = true
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&res)
}

func (h *handler) joinRoom(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	roomUUID := p.ByName("uuid")
	username := r.URL.Query().Get("username")

	websocket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Errorf("error while upgrading websocket connection: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	Mes := make(chan messages.Message)
	ErrMes := make(chan messages.AppError)

	fmt.Println("hub SRM : ", h.hub.SRM)
	ctx, cancel := context.WithCancel(context.Background())
	worker := &Worker{
		username:        username,
		roomUUID:        roomUUID,
		logger:          h.logger,
		Message:         Mes,
		ErrMes:          ErrMes,
		websocket:       NewWsManager(websocket, ErrMes),
		webrtc:          webrtc.New(Mes, ErrMes),
		game:            h.hub.Rooms[roomUUID].Game,
		resourceManager: h.hub.SRM,
		ctx:             ctx,
		cancelFunc:      cancel,
	}
	h.hub.ConnectPlayer <- worker

	go worker.Run(h.hub)
}
