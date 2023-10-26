package hub

import (
	"encoding/json"
	"fmt"
	"net/http"

	"cloud/internal/domain/hub"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

const (
	roomURL   = "/room"
	createURL = "/room/create"
	joinURL   = "/room/join/:roomId"
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
	router.HandlerFunc(http.MethodGet, roomURL, h.getRoomList)
	router.GET(joinURL, h.joinRoom)
	router.GET("/dv", h.getPeers)
}
func (h *Handler) getRoomList(w http.ResponseWriter, r *http.Request) {
	res := make([]hub.GetRoomRes, 0)
	for _, room := range h.hub.Rooms {
		res = append(res, hub.GetRoomRes{
			ID:   room.ID,
			Name: room.Name,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) createRoom(w http.ResponseWriter, r *http.Request) {
	var req hub.CreateRoomRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	h.hub.Rooms[req.ID] = &hub.Room{
		ID:    req.ID,
		Name:  req.Name,
		Peers: make(map[string]*hub.Peer),
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)

}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		//origin := r.Header.Get("Origin")
		//if origin == "http://localhost:3000" {
		return true
	},
}

func (h *Handler) joinRoom(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// roon/joinRooom/:roomID?peerID=123
	fmt.Println("router params: ", p)
	roomID := p.ByName("roomID")
	peerID := r.URL.Query().Get("peerID")
	username := r.URL.Query().Get("username")

	peer := &hub.Peer{
		ID:       peerID,
		Username: username,
		Host:     false,
		Conn:     conn,
		Message:  make(chan *hub.Message),
		RoomID:   roomID,
	}
	m := &hub.Message{
		Content:  "Hello",
		RoomID:   roomID,
		Username: username,
	}
	//register peer through channel
	h.hub.Register <- peer
	//broadcast message to all peers in room
	h.hub.Broadcast <- m
	//send message to all peers in room
	//read message from peer
	go peer.WriteMessage()
	peer.ReadMessage(h.hub)

}

func (h *Handler) getPeers(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var res []hub.GetPeersRes
	roomId := p.ByName("uuid")
	if _, ok := h.hub.Rooms[roomId]; !ok {
		res = make([]hub.GetPeersRes, 0)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
		return
	}
	for _, peer := range h.hub.Rooms[roomId].Peers {
		res = append(res, hub.GetPeersRes{
			ID:   peer.ID,
			Name: peer.Username,
		})
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
