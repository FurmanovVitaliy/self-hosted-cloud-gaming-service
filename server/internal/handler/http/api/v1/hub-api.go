package v1

import (
	"net/http"

	"github.com/FurmanovVitaliy/pixel-cloud/internal/usecase"
	"github.com/gorilla/websocket"
)

func (h *handler) GetRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.uc.GetRooms()
	if err != nil {
		errorResponse(err, w, r)
		return
	}
	if err := encode(w, r, http.StatusOK, rooms); err != nil {
		errorResponse(err, w, r)
		return
	}
}

func (h *handler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	req, err := decode[usecase.CreateRoomReq](r)
	if err != nil {
		errorResponse(err, w, r)
	}
	res, err := h.uc.CreateRoom(req)
	if err != nil {
		errorResponse(err, w, r)
		return
	}
	if err := encode(w, r, http.StatusCreated, res); err != nil {
		errorResponse(err, w, r)
		return
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) (b bool) { return true },
}

func (h *handler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Path[len("/room/join/"):]
	username := r.URL.Query().Get("username")

	websocket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		errorResponseWebsocket(err, websocket)
		return
	}
	var deviceInfo usecase.JoinRoomRes

	err = websocket.ReadJSON(&deviceInfo)
	if err != nil {
		errorResponseWebsocket(err, websocket)
		return
	}

	err = h.uc.JoinRoom(uuid, username, websocket, deviceInfo)
	if err != nil {
		errorResponseWebsocket(err, websocket)
		return
	}
}
