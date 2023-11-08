package hub

import (
	"encoding/json"
	"fmt"
	"net/http"

	"cloud/internal/domain/hub"
	"cloud/pkg/webrtc"

	"cloud/pkg/network/websocket"

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
	room.Websocket = websocket.NewWebsocket(1024, 1024, UUID)
	room.Websocket.Upgrade(w, r)

	room.Peer = webrtc.NewPeer()
	offer, err := room.Peer.NewWebRTC("vp8", "opus")
	if err != nil {
		fmt.Println(err)
	}
	room.Websocket.Conn.WriteJSON(offer)
	for _, candidate := range room.Peer.ICECandidate {
		room.Websocket.Conn.WriteJSON(candidate)
	}
	room.Websocket.ReadWebRtcMesInGoRoutine(func(any interface{}) {
		room.Peer.SetCandidatesAndSDP(any)
	})
	go room.Peer.SendVideo()
	room.Peer.OnMessage = handleInput
}

func handleInput(data []byte) {

	fmt.Println("handle input: ", string(data))
}
