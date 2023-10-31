package hub

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"cloud/internal/domain/hub"
	"cloud/pkg/network/socket"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/pion/webrtc/v3"
)

const (
	roomsURL  = "/room"
	createURL = "/room/create"
	roomURL   = "/room/join/:uuid"
)

type Handler struct {
	hub *hub.Hub
}

func NewHandler(hub *hub.Hub) *Handler {
	return &Handler{
		hub: hub,
	}
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, createURL, h.createRoom)
	router.HandlerFunc(http.MethodGet, roomsURL, h.getRoomList)
	router.GET(roomURL, h.upgradeRoomConnection)
	router.GET("/dv", h.getPeers)
}
func (h *Handler) getRoomList(w http.ResponseWriter, r *http.Request) {
	res := make([]hub.GetRoomRes, 0)
	for _, room := range h.hub.Rooms {
		res = append(res, hub.GetRoomRes{
			UUID: room.UUID,
			Game: room.Game,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) createRoom(w http.ResponseWriter, r *http.Request) {
	var req hub.CreateRoomReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println("create room request: ", req)
	h.hub.Rooms[req.UUID] = &hub.Room{
		UUID: req.UUID,
		Game: req.Game,
	}
	res := hub.CreateRoomRes{
		UUID: req.UUID,
		Game: req.Game,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
	fmt.Println("room created:", h.hub.Rooms)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		//origin := r.Header.Get("Origin")
		//if origin == "http://localhost:3000" {
		return true
	},
}

func (h *Handler) upgradeRoomConnection(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// roon/joinRooom/:roomUUID
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer conn.Close()
	UUID := p.ByName("uuid")
	room := h.hub.Rooms[UUID]
	room.Conn = conn
	room.WebRTC = createPeerConnection()
	peerConnection := createPeerConnection()
	defer peerConnection.Close()
	// Создание UDP-сокета для прослушивания порта 5004
	listener, err := socket.NewVideoUDPListener()
	if err != nil {
		fmt.Println("Ошибка при создании сокета:", err)
	}
	// Создание видео-трека
	videoTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeVP8}, "video", "pion")
	if err != nil {
		panic(err)
	}
	// Добавление видео-трека в отправитель (RTP Sender)
	rtpSender, err := peerConnection.AddTrack(videoTrack)
	if err != nil {
		panic(err)
	}
	// Горутина для чтения входящих пакетов RTCP
	go func() {
		rtcpBuf := make([]byte, 1500)
		for {
			if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
				return
			}
		}
	}()

	// Обработка изменений состояния ICE-соединения
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Connection State has changed:  %s \n", connectionState.String())
	})
	peerConnection.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate != nil {
			fmt.Println("server get candidate")
			// Отправляем ICE кандидата клиенту через WebSocket
			jsonCandidate, _ := json.Marshal(candidate)
			err = conn.WriteMessage(websocket.TextMessage, jsonCandidate)
			log.Println("server send candidate")
			if err != nil {
				log.Println(err)
				return
			}
		}
	})

	// Создание оффера
	offer, err := peerConnection.CreateOffer(nil)
	log.Println("servet creationg offer")
	if err != nil {
		log.Println(err)
		return
	}
	// Установка локального описания сессии и запуск UDP-слушателей
	err = peerConnection.SetLocalDescription(offer)
	log.Println("servet set local description")
	if err != nil {
		log.Println(err)
		return
	}

	// Отправляем SDP оффер клиенту через WebSocket
	err = conn.WriteJSON(offer)
	log.Println("server send offer")
	if err != nil {
		log.Println(err)
		return
	}
	// Ожидаем SDP ответ от клиента через WebSocket
	go ListenForWs(conn, peerConnection)

	inboundRTPPacket := make([]byte, 1600) // UDP MTU
	for {
		n, _, err := listener.ReadFrom(inboundRTPPacket)
		if err != nil {
			panic(fmt.Sprintf("Ошибка при чтении: %s", err))
		}

		if _, err = videoTrack.Write(inboundRTPPacket[:n]); err != nil {
			if errors.Is(err, io.ErrClosedPipe) {

				return
			}

			panic(err)
		}
	}
}

func (h *Handler) getPeers(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var res []hub.GetPeersRes
	roomId := p.ByName("uuid")
	if _, ok := h.hub.Rooms[roomId]; !ok {
		res = make([]hub.GetPeersRes, 0)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func createPeerConnection() *webrtc.PeerConnection {
	// Создание и настройка объекта PeerConnection
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	return peerConnection
}
func ListenForWs(conn *websocket.Conn, peerConnection *webrtc.PeerConnection) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, fmt.Sprintf("%v", err))
		}
	}()

	for {
		messageType, payload, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
		}

		if messageType == websocket.TextMessage {
			var candidate webrtc.ICECandidateInit
			var answer webrtc.SessionDescription
			_ = json.Unmarshal([]byte(payload), &candidate)
			_ = json.Unmarshal([]byte(payload), &answer)

			if answer.SDP != "" {
				fmt.Println("server get answer:")
				fmt.Println(answer)
				_ = peerConnection.SetRemoteDescription(answer)

			}
			if candidate.UsernameFragment != nil {
				fmt.Println("server get candidate:")
				fmt.Println(candidate)
				_ = peerConnection.AddICECandidate(candidate)
			}
		}
	}
}
