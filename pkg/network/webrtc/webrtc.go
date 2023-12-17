package webrtc

import (
	"cloud/pkg/logger"
	"net"
	"time"

	"fmt"

	"strings"

	"github.com/pion/webrtc/v4"
)

type RTCMessage struct {
	ContentType string `json:"content_type"`
	Content     string `json:"content"`
}

type WRTC struct {
	conn      *webrtc.PeerConnection
	logger    logger.Logger
	OnMessage func(data []byte)

	isAnswer bool
	stop     chan (bool)

	audio *webrtc.TrackLocalStaticRTP
	video *webrtc.TrackLocalStaticRTP
	data  *webrtc.DataChannel
}

func New() *WRTC {
	return &WRTC{
		stop:     make(chan bool),
		isAnswer: false,
		logger:   logger.Init("7"),
	}
}
func (w *WRTC) Start(vCodec, aCodec string, sendToClient func(any) error) (err error) {
	w.conn, err = webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{{URLs: []string{"stun:stun1.l.google.com:19302"}}},
	})
	if err != nil {
		return
	}
	w.conn.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		w.logger.Infof("connection state has changed:  %s \n", connectionState.String())
		switch connectionState {
		case webrtc.ICEConnectionStateClosed:
			w.Stop()
		case webrtc.ICEConnectionStateDisconnected:
			w.Stop()
		}
	})

	w.conn.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			sendToClient(RTCMessage{
				ContentType: SIGNAL,
				Content:     SERVER_ICE_GATHERING_COMPLETE,
			})
			w.logger.Warn("ice gathering is complete")
		}
		if candidate != nil {
			sendToClient(RTCMessage{
				ContentType: SERVER_CANDIDATE,
				Content:     Encode(candidate.ToJSON()),
			})
			w.logger.Debug("ice candidate found :", candidate)
		}
	})

	// plug in the [video] track (out)
	video, err := newTrack("video", "pion", vCodec)
	if err != nil {
		w.logger.Error("error while creating video track", err)
	}
	_, err = w.conn.AddTrack(video)
	if err != nil {
		w.logger.Error("error while adding video track to Peer", err)
		return
	}
	w.video = video
	w.logger.Debugf("added [%s] track", video.Codec().MimeType)

	// plug in the [audio] track (out)
	audio, err := newTrack("audio", "pion", aCodec)
	if err != nil {
		w.logger.Error("error while creating audio track", err)
	}
	_, err = w.conn.AddTrack(audio)
	if err != nil {
		w.logger.Error("error while adding audio track to Peer", err)
	}
	w.audio = audio
	w.logger.Debugf("added [%s] track", audio.Codec().MimeType)
	// plug in the [data] channel (in and out)
	if err = w.addDataChannel("input"); err != nil {
		return
	}
	w.logger.Debug("added [data] chan")

	offer, err := w.conn.CreateOffer(nil)
	if err != nil {
		w.logger.Error(err)
		return
	}
	w.logger.Info("server created offer")

	err = w.conn.SetLocalDescription(offer)
	if err != nil {
		w.logger.Error(err)
		return
	}
	w.logger.Info("server local description is set")

	sendToClient(RTCMessage{
		ContentType: OFFER,
		Content:     Encode(offer),
	})

	return
}
func (w *WRTC) SetAnswer(data string) error {
	var answear webrtc.SessionDescription
	Decode(data, &answear)
	err := w.conn.SetRemoteDescription(answear)
	if err != nil {
		w.logger.Errorf("error occurred while setting remote description: %s", err)
		return err
	}
	w.logger.Info("remote description is set")
	w.isAnswer = true
	return nil
}

func (w *WRTC) AddCandidate(data string) error {
	w.waitForAction(w.isAnswer)
	var candidate webrtc.ICECandidateInit
	Decode(data, &candidate)
	w.logger.Debug("received candidate: ", candidate)
	err := w.conn.AddICECandidate(candidate)
	if err != nil {
		w.logger.Errorf("error occurred while adding candidate: %s", err)
		return err
	}
	w.logger.Info("candidate added")
	return nil
}

func (w *WRTC) SendData(data []byte) {
	err := w.data.Send(data)
	if err != nil {
		w.logger.Errorf("error occurred while data sending: %s", err)
	}
}

func (w *WRTC) SendVideo(ip string, port int) {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP(ip), Port: port})
	if err != nil {
		w.logger.Errorf("error occurred while socket creating: %s", err)
		return
	}
	defer listener.Close()

	if err := listener.SetReadBuffer(1000000); err != nil {
		w.logger.Errorf("error occurred while setting read buffer: %s", err)
		return
	}

	rtpBuf := make([]byte, 5000)
readBufData:
	for {
		select {
		case <-w.stop:
			break readBufData // выход из цикла при получении сигнала остановки
		default:
			n, _, err := listener.ReadFrom(rtpBuf)
			if err != nil {
				w.logger.Errorf("error occurred while data reading: %s", err)
				break readBufData // выход из цикла при ошибке чтения данных
			}
			if _, err := w.video.Write(rtpBuf[:n]); err != nil {
				w.logger.Errorf("error occurred while video sending: %s", err)
				break readBufData // выход из цикла при ошибке отправки видео
			}
		}
	}
}

func (w *WRTC) SendAudio(data []byte) {}

func (w *WRTC) Stop() {
	w.stop <- true
	w.data.Close()
	w.conn.Close()

}

func newTrack(id string, label string, codec string) (*webrtc.TrackLocalStaticRTP, error) {
	codec = strings.ToLower(codec)
	var mime string
	switch id {
	case "audio":
		switch codec {
		case "opus":
			mime = webrtc.MimeTypeOpus
		}
	case "video":
		switch codec {
		case "h264":
			mime = webrtc.MimeTypeH264
		case "vpx", "vp8":
			mime = webrtc.MimeTypeVP8
		}
	}
	if mime == "" {
		return nil, fmt.Errorf("unsupported codec %s:%s", id, codec)
	}

	return webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: mime}, id, label)
}

func (w *WRTC) addDataChannel(label string) error {
	ch, err := w.conn.CreateDataChannel(label, nil)
	if err != nil {
		return err
	}
	ch.OnOpen(func() {
		w.logger.Info("data channel [input] opened")
		w.logger.Debugf("label: %s, id: %d", ch.Label(), ch.ID())
	})
	ch.OnError(func(err error) { w.logger.Error(err) })
	ch.OnMessage(func(m webrtc.DataChannelMessage) {
		if len(m.Data) == 0 {
			return
		}
		if w.OnMessage != nil {
			w.OnMessage(m.Data)
		}
	})
	w.data = ch
	ch.OnClose(func() { w.logger.Info("data channel [input] closed") })
	return nil
}

// hack
func (w *WRTC) waitForAction(action bool) {
	for {
		if !action {
			time.Sleep(3 * time.Second)
			w.logger.Warn("waiting for action")
			w.waitForAction(action)
		} else {
			time.Sleep(2 * time.Second)
			break
		}
	}
}
