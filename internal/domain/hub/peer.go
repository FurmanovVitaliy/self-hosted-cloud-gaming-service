package hub

import (
	"log"

	"github.com/gorilla/websocket"
)

type Message struct {
	Content  string `json:"content"`
	RoomID   string `json:"roomID"`
	Username string `json:"username"`
}

type Peer struct {
	ID       string
	Host     bool
	RoomID   string
	Username string
	Message  chan *Message
	Conn     *websocket.Conn
}
type GetPeersRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (p *Peer) WriteMessage() {
	defer func() {
		p.Conn.Close()
	}()
	for {
		message, ok := <-p.Message
		if !ok {
			return
		}
		p.Conn.WriteJSON(message)
	}
}
func (p *Peer) ReadMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- p
		p.Conn.Close()
	}()
	for {
		_, payload, err := p.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		meg := &Message{
			Content:  string(payload),
			RoomID:   p.RoomID,
			Username: p.Username,
		}
		hub.Broadcast <- meg
	}
}
