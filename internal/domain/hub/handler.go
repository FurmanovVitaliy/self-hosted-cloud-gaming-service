package hub

import (
	"cloud/internal/domain/games"
	"cloud/pkg/logger"
	"cloud/pkg/network/webrtc"

	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

const (
	roomsURL  = "/room"
	createURL = "/room/create"
	roomURL   = "/room/join/:uuid"
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
	var req CreateRoomReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	game, err := h.gameService.GetOne(req.GameID)
	if err != nil {
		h.logger.Errorf("Game with id %s not found while room (uuid %s) creating : %s", req.GameID, req.UUID, err)
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
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) (b bool) { return true },
}

func (h *handler) joinRoom(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	websocket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Errorf("error while upgrading websocket connection: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	roomUUID := p.ByName("uuid")
	username := r.URL.Query().Get("username")

	worker := &Worker{
		username:  username,
		roomUUID:  roomUUID,
		logger:    h.logger,
		RegMes:    make(chan Message),
		ErrMes:    make(chan error),
		websocket: *NewWsManager(websocket),
		webrtc:    webrtc.New(),
		game:      h.hub.Rooms[roomUUID].Game,
	}
	h.hub.ConnectPlayer <- worker

	worker.Run(h.hub)
}

/*
var vg, _ = uinput.CreateGamepad("/dev/uinput", []byte("testpad"), 045, 955)

var state = keymapping.XboxGpadInput{}

func handleInput(data []byte) {
	//fmt.Println("handle input: ", string(data))
	json.Unmarshal(data, &state)

	x, _ := strconv.ParseFloat(state.RS.X, 32)
	y, _ := strconv.ParseFloat(state.RS.Y, 32)

	vg.RightStickMove(float32(x), float32(y))

	lsX, _ := strconv.ParseFloat(state.LS.X, 32)
	lsY, _ := strconv.ParseFloat(state.LS.Y, 32)
	vg.LeftStickMove(float32(lsX), float32(lsY))

	if state.A == 1 {
		vg.ButtonDown(uinput.ButtonSouth)
	} else {
		vg.ButtonUp(uinput.ButtonSouth)
	}
	if state.X == 1 {
		vg.ButtonPress(uinput.ButtonWest)
	}
	if state.Y == 1 {
		vg.ButtonPress(uinput.ButtonNorth)
	}
	if state.LB == 1 {
		vg.ButtonPress(uinput.ButtonBumperLeft)
	}
	if state.RB == 1 {
		vg.ButtonPress(uinput.ButtonBumperRight)
	}
	if state.Start == 1 {
		vg.ButtonPress(uinput.ButtonStart)
	}
	if state.Main == 1 {
		vg.ButtonPress(uinput.ButtonMode)
	}
*/
