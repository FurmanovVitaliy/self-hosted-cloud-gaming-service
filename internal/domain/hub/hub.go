package hub

import (
	"cloud/internal/domain/srm"
	"cloud/internal/messages"
	"fmt"
	"log"
	"sync"
)

type Hub struct {
	Mutex            sync.Mutex
	Rooms            map[string]*Room
	Broadcast        chan *messages.Message
	ConnectPlayer    chan *Worker
	DisconnectPlayer chan *Worker
	SRM              *srm.ServerResourceManager
}

func NewHub() *Hub {
	return &Hub{
		Rooms:            make(map[string]*Room),
		Broadcast:        make(chan *messages.Message),
		ConnectPlayer:    make(chan *Worker),
		DisconnectPlayer: make(chan *Worker),
		SRM:              srm.NewResourceManager(),
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
					log.Println("disconnecting player:", worker.username)
					if len(h.Rooms[worker.roomUUID].Workers) > 1 {
						h.Broadcast <- &messages.Message{
							RoomUUID:    worker.roomUUID,
							Username:    worker.username,
							ContentType: "broadcast",
							Content:     fmt.Sprintf("%s left the room", worker.username),
						}

					}
					h.SRM.StopVM(worker.serverResources.VMID)
					fmt.Println("killed docker: ", worker.serverResources.VMID)
					h.SRM.ReleaseResources(worker.serverResources)

					fmt.Println("resources released:", worker.serverResources.XServer.ScreenNumber, worker.serverResources.Listener)

					delete(h.Rooms[worker.roomUUID].Workers, worker.username)
					close(worker.Message)
					close(worker.ErrMes)
					if len(h.Rooms[worker.roomUUID].Workers) == 0 {
						delete(h.Rooms, worker.roomUUID)
					}
				}
			}
		case message := <-h.Broadcast:
			if _, ok := h.Rooms[message.RoomUUID]; ok {
				for _, worker := range h.Rooms[message.RoomUUID].Workers {
					worker.Message <- *message

				}
			}

		}
	}
}
