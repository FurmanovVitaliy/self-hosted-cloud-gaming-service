package hub

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WsManager struct {
	mutex sync.RWMutex
	conn  *websocket.Conn
}

func NewWsManager(conn *websocket.Conn) *WsManager {
	return &WsManager{
		conn: conn,
	}
}

func (w *WsManager) WriteJSON(v interface{}) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	err := w.conn.WriteJSON(v)
	if err != nil {
		return err
	}
	return nil
}

func (w *WsManager) ReadMessage() (*Message, error) {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	var msg Message
	err := w.conn.ReadJSON(&msg)
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			return nil, err
		}
		return nil, err
	}
	return &msg, nil
}
func (w *WsManager) Close() error {
	err := w.conn.Close()
	if err != nil {
		return err
	}
	return nil
}
