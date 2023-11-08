package websocket

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Websocket struct {
	ID       string
	Conn     *websocket.Conn
	upgrader websocket.Upgrader
}

func NewWebsocket(readBuf, writeBuf int, ID string) *Websocket {
	return &Websocket{
		ID: ID + " room websocket",
		upgrader: websocket.Upgrader{
			ReadBufferSize:  readBuf,
			WriteBufferSize: writeBuf,
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
func (w *Websocket) WriteICE(ice any) error {
	err := w.Conn.WriteJSON(ice)
	log.Println("Write ICE Error :", err)
	return err
}
