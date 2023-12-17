package hub

import (
	"fmt"
)

type Message struct {
	RoomUUID    string `json:"room_uuid"`
	Username    string `json:"username"`
	ContentType string `json:"content_type"`
	Content     string `json:"content"`
}

type Hub struct {
	Rooms            map[string]*Room
	Broadcast        chan *Message
	ConnectPlayer    chan *Worker
	DisconnectPlayer chan *Worker
}

func NewHub() *Hub {
	return &Hub{
		Rooms:            make(map[string]*Room),
		Broadcast:        make(chan *Message),
		ConnectPlayer:    make(chan *Worker),
		DisconnectPlayer: make(chan *Worker),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case worker := <-h.ConnectPlayer:
			if _, ok := h.Rooms[worker.roomUUID]; ok {
				room := h.Rooms[worker.roomUUID]
				if _, ok := room.Workers[worker.username]; !ok {
					room.Workers[worker.username] = worker
				}
			}
		case worker := <-h.DisconnectPlayer:
			if _, ok := h.Rooms[worker.roomUUID]; ok {
				if _, ok := h.Rooms[worker.roomUUID].Workers[worker.username]; ok {
					if len(h.Rooms[worker.roomUUID].Workers) != 0 {
						h.Broadcast <- &Message{
							RoomUUID:    worker.roomUUID,
							Username:    worker.username,
							ContentType: "broadcast",
							Content:     fmt.Sprintf("%s left the room", worker.username),
						}
						delete(h.Rooms[worker.roomUUID].Workers, worker.username)
						close(worker.RegMes)
					}
				}
			}
		case message := <-h.Broadcast:
			if _, ok := h.Rooms[message.RoomUUID]; ok {
				for _, worker := range h.Rooms[message.RoomUUID].Workers {
					worker.RegMes <- *message

				}
			}

		}
	}
}
