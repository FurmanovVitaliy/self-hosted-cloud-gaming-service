package broker

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sync"

	"github.com/gorilla/websocket"
)

type Broker struct {
	conn   *websocket.Conn
	routes map[string]reflect.Value
	mu     sync.Mutex
}

func New(ws *websocket.Conn) *Broker {
	return &Broker{
		conn:   ws,
		routes: make(map[string]reflect.Value),
	}
}

func (b *Broker) RegisterChannel(tag string, channel interface{}) error {
	if _, ok := b.routes[tag]; ok {
		return fmt.Errorf("Channel already registered for tag: %s", tag)
	}

	channelType := reflect.TypeOf(channel)
	if channelType.Kind() != reflect.Chan {
		return fmt.Errorf("Expected a channel, got %s", channelType.Kind())
	}

	b.routes[tag] = reflect.ValueOf(channel)
	log.Println("Registering channel for tag:", tag)
	return nil
}

func (b *Broker) Read() {
	defer func() {
		for tag, channel := range b.routes {
			if channel.Kind() == reflect.Chan && channel.IsValid() {
				channel.Close()
				delete(b.routes, tag)
				log.Printf("Closed channel for tag: '%s'", tag)
			}
		}
		b.conn.Close()
	}()
	log.Println("Starting message broker")
	for {
		if b.conn == nil {
			log.Println("Connection closed")
			break
		}
		var rawMsg map[string]interface{}
		err := b.conn.ReadJSON(&rawMsg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			log.Println("WebSocket read error:", err)
			break
		}

		tag, ok := rawMsg["tag"].(string)
		if !ok {
			log.Println("Invalid message format: Tag missing")
			continue
		}

		channel, ok := b.routes[tag]
		if !ok {
			log.Println("No channel registered for content type:", tag)
			continue
		}

		msgType := channel.Type().Elem()
		msgBytes, err := json.Marshal(rawMsg)
		if err != nil {
			log.Printf("Error marshalling message for tag %s: %v", tag, err)
			log.Println("Message:", rawMsg)
			continue
		}

		msgValue := reflect.New(msgType).Interface()
		err = json.Unmarshal(msgBytes, msgValue)
		if err != nil {
			log.Printf("Error unmarshalling message for tag %s: %v", tag, err)
			log.Println("Message:", rawMsg)
			continue
		}

		channel.Send(reflect.ValueOf(msgValue).Elem())
	}
}

func (b *Broker) Write(tag string, msg interface{}) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	msgMap := make(map[string]interface{})

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Println("Error marshalling message:", err)
		return err
	}

	err = json.Unmarshal(msgBytes, &msgMap)
	if err != nil {
		log.Println("Error unmarshalling message to map:", err)
		return err
	}

	msgMap["tag"] = tag

	err = b.conn.WriteJSON(msgMap)
	if err != nil {
		log.Println("Error writing message:", err)
	}
	return err
}

func (b *Broker) Stop() {
	log.Println("Stopping message broker")
	err := b.conn.Close()
	if err != nil {
		log.Println("Error closing WebSocket connection:", err.Error())
		log.Println("Error closing WebSocket connection:", err)
	}
}
