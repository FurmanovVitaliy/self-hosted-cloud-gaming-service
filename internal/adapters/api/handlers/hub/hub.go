package hub

import (
	"cloud/internal/domain/hub"
	"cloud/pkg/input/keymapping"
	"cloud/pkg/webrtc"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"cloud/pkg/network/websocket"

	"github.com/bendahl/uinput"
	"github.com/julienschmidt/httprouter"
)

const (
	roomsURL  = "/room"
	createURL = "/room/create"
	roomURL   = "/room/join/:uuid"
)

type Handler struct {
	hub *hub.Hub
}

func NewHandler(hub *hub.Hub) *Handler {
	return &Handler{
		hub: hub,
	}
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, createURL, h.createRoom)
	router.HandlerFunc(http.MethodGet, roomsURL, h.getRoomList)
	router.GET(roomURL, h.upgradeRoomConnection)
}
func (h *Handler) getRoomList(w http.ResponseWriter, r *http.Request) {
	res := make([]hub.GetRoomRes, 0)
	for _, room := range h.hub.Rooms {
		res = append(res, hub.GetRoomRes{
			UUID: room.UUID,
			Game: room.Game,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) createRoom(w http.ResponseWriter, r *http.Request) {
	var req hub.CreateRoomReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println("create room request: ", req)
	h.hub.Rooms[req.UUID] = &hub.Room{
		UUID: req.UUID,
		Game: req.Game,
	}
	res := hub.CreateRoomRes{
		UUID: req.UUID,
		Game: req.Game,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
	fmt.Println("room created:", h.hub.Rooms)
}

func (h *Handler) upgradeRoomConnection(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	UUID := p.ByName("uuid")
	room := h.hub.Rooms[UUID]
	room.Websocket = websocket.New(UUID)
	room.Websocket.Upgrade(w, r)

	room.Peer = webrtc.NewPeer()
	offer, err := room.Peer.NewWebRTC("h264", "opus", room.Websocket.SendICE)
	if err != nil {
		fmt.Println(err)
	}
	room.Websocket.WriteJson(offer)
	fmt.Println("offer sended: ")

	room.Websocket.ReadWebRtcMesInGoRoutine(func(any interface{}) {
		room.Peer.SetCandidatesAndSDP(any)
	})
	go room.Peer.SendVideo()

	room.Peer.OnMessage = handleInput
	//vg.Close()

}

var vg, _ = uinput.CreateGamepad("/dev/uinput", []byte("testpad"), 045, 955)

var state = keymapping.XboxGpadInput{}

func handleInput(data []byte) {
	fmt.Println("handle input: ", string(data))
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

}
