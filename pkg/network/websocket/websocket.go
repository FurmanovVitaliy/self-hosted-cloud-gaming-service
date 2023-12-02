package websocket

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Websocket struct {
	ID       string
	mutex    sync.RWMutex
	Conn     *websocket.Conn
	upgrader websocket.Upgrader
}

func New(ID string) *Websocket {
	return &Websocket{
		ID: ID + " room websocket",
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (w *Websocket) Upgrade(wr http.ResponseWriter, r *http.Request) error {
	conn, err := w.upgrader.Upgrade(wr, r, nil)
	if err != nil {
		return err
	}
	w.Conn = conn
	return nil
}

func (w *Websocket) ReadWebRtcMesInGoRoutine(WebRtcMassegeHandler func(any)) {
	go func() {
		defer func() {
			w.mutex.RLock()
			defer w.mutex.RUnlock()
			if err := recover(); err != nil {
				log.Println(err, fmt.Sprintf("%v", err))
				//TODO: log error
			}
		}()
		for {

			mesType, payload, err := w.Conn.ReadMessage()
			if err != nil {
				log.Println("Read Error :", err)
				//TODO: log error
			}
			if mesType == websocket.TextMessage {
				WebRtcMassegeHandler(payload)
			}
		}
	}()
}
func (w *Websocket) SendICE(ice any) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	err := w.Conn.WriteJSON(ice)
	if err != nil {
		return err
	}
	return nil
}
func (w *Websocket) WriteJson(v interface{}) error {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	err := w.Conn.WriteJSON(v)
	if err != nil {
		return err
	}
	return nil
}
