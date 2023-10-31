package hub

type Hub struct {
	Rooms map[string]*Room
	///Register   chan *Peer
	//Unregister chan *Peer
	//Broadcast  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms: make(map[string]*Room),
		//Register:   make(chan *Peer),
		//Unregister: make(chan *Peer),
		//Broadcast:  make(chan *Message),
	}
}

/*func (h *Hub) Run() {
	for {
		select {
		case peer := <-h.Register:
			if _, ok := h.Rooms[peer.RoomID]; ok {
				r := h.Rooms[peer.RoomID]
				if _, ok := r.Peers[peer.ID]; !ok {
					r.Peers[peer.ID] = peer
				}
			}
		case peer := <-h.Unregister:
			if _, ok := h.Rooms[peer.RoomID]; ok {
				if _, ok := h.Rooms[peer.RoomID].Peers[peer.ID]; ok {
					//send info that client left he room
					delete(h.Rooms[peer.RoomID].Peers, peer.ID)
					close(peer.Message)
				}
			}

		case message := <-h.Broadcast:
			if _, ok := h.Rooms[message.RoomID]; ok {
				r := h.Rooms[message.RoomID]
				for _, peer := range r.Peers {
					select {
					case peer.Message <- message:
					default:
						close(peer.Message)
						delete(r.Peers, peer.ID)
					}
				}
			}
		}
	}
}
*/
