package hub

import (
	"fmt"
	"sync"

	"github.com/FurmanovVitaliy/pixel-cloud/pkg/logger"
)

type Room struct {
	UUID    string
	GameID  string
	Workers map[string]*Worker
}

type Hub struct {
	Mutex      sync.Mutex
	Rooms      map[string]*Room
	Connect    chan *Worker
	Disconnect chan *Worker
	Broadcast  chan *ChatMsg
}

func New() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Connect:    make(chan *Worker),
		Disconnect: make(chan *Worker),
		Broadcast:  make(chan *ChatMsg, 5),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case worker := <-h.Connect:
			h.Mutex.Lock()
			if room, exists := h.Rooms[worker.roomUUID]; exists {
				if len(room.Workers) == 0 {
					worker.host = true
				}
				room.Workers[worker.Username] = worker
				h.Broadcast <- &ChatMsg{
					RoomUUID: worker.roomUUID,
					Username: worker.Username,
					Content:  fmt.Sprintf("%s joined the room", worker.Username),
				}
			}
			h.Mutex.Unlock()

		case worker := <-h.Disconnect:
			fmt.Println("disconnect worker in HUB")
			h.Mutex.Lock()
			if room, exist := h.Rooms[worker.roomUUID]; exist {
				if _, exists := room.Workers[worker.Username]; exists {
					worker.vm.Stop()
					worker.streamer.Stop()
					worker.msgHandler.Stop()
					delete(room.Workers, worker.Username)
					if len(room.Workers) > 0 {
						h.Broadcast <- &ChatMsg{
							RoomUUID: worker.roomUUID,
							Username: worker.Username,
							Content:  fmt.Sprintf("%s left the room", worker.Username),
						}
					}
					if len(room.Workers) == 0 {
						delete(h.Rooms, room.UUID)
					}
				}
			}
			h.Mutex.Unlock()

		case chatMsg := <-h.Broadcast:
			if room, exists := h.Rooms[chatMsg.RoomUUID]; exists {
				for _, player := range room.Workers {
					player.WriteChatMsg(chatMsg)
				}
			}
		}
	}
}

func (h *Hub) GetRooms() ([]Room, error) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	var rooms []Room
	for _, room := range h.Rooms {
		rooms = append(rooms, *room)
	}
	return rooms, nil
}

func (h *Hub) GetRoom(uuid string) (*Room, error) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	room, ok := h.Rooms[uuid]
	if !ok {
		return nil, fmt.Errorf("room not found")
	}
	return room, nil
}

func (h *Hub) CreateRoom(uuid, gameID string) error {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	if _, ok := h.Rooms[uuid]; ok {
		return fmt.Errorf("room already exists")
	}
	h.Rooms[uuid] = &Room{
		UUID:    uuid,
		GameID:  gameID,
		Workers: make(map[string]*Worker),
	}
	return nil
}

func (h *Hub) CreateWorker(roomUUID, username string, msgHandler MessageHandler, streamer WrtcStreamer, audioR, videoR UDPReader, vm VM) error {
	worker := &Worker{
		msgHandler:  msgHandler,
		streamer:    streamer,
		roomUUID:    roomUUID,
		Username:    username,
		chatMsg:     make(chan *ChatMsg, 10),
		wrtcMsg:     make(chan *WrtcMsg, 50),
		audioReader: audioR,
		videoReader: videoR,
		vm:          vm,
		logger:      logger.Init("debug"),
	}

	err := msgHandler.RegisterChannel("chat", worker.chatMsg)
	if err != nil {
		return err
	}
	err = msgHandler.RegisterChannel("wrtc", worker.wrtcMsg)
	if err != nil {
		return err
	}

	h.Connect <- worker

	go worker.ReadChatMsg(h)

	go worker.ReadStreamerMsg(h)

	go worker.SrartVM(h)

	return nil
}
