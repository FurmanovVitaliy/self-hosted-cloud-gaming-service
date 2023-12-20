package hub

import (
	"cloud/internal/messages"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type WsManager struct {
	mutex  sync.RWMutex
	conn   *websocket.Conn
	errMes chan messages.AppError
}

func NewWsManager(conn *websocket.Conn, err chan messages.AppError) WsManager {
	return WsManager{
		errMes: err,
		conn:   conn,
	}
}

func (w *WsManager) WriteJSON(v interface{}) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	err := w.conn.WriteJSON(v)
	if err != nil {
		w.errMes <- *messages.NewAppError(err, "error occurred while writing json", "", "")
	}
}

func (w *WsManager) ReadMessage() *messages.Message {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	var msg messages.Message
	err := w.conn.ReadJSON(&msg)
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			return nil
		}
		w.errMes <- *messages.NewAppError(err, "error occurred while reading json", "", "")
		return nil
	}
	return &msg
}

func (w *WsManager) Close() {
	err := w.conn.Close()
	if err != nil {
		log.Println(err)
	}
}
