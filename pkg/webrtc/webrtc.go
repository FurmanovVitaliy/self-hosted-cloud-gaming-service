package webrtc

import (
	"cloud/pkg/logger"
	"net"
	"time"

	"encoding/json"
	"fmt"

	"strings"

	"github.com/pion/webrtc/v4"
)

//TODO: add madia engine for codecs controle & all data must be in base 64

type Peer struct {
	conn       *webrtc.PeerConnection
	logger     logger.Logger
	OnMessage  func(data []byte)
	Candidates chan webrtc.ICECandidateInit

	signalStop chan struct{}

	audio *webrtc.TrackLocalStaticRTP
	video *webrtc.TrackLocalStaticRTP
	data  *webrtc.DataChannel
}

func NewPeer() *Peer {
	return &Peer{
		signalStop: make(chan struct{}),
		Candidates: make(chan webrtc.ICECandidateInit),
		logger:     logger.Init("7"),
	}
}
func (p *Peer) NewWebRTC(vCodec, aCodec string, sendICE func(any) error) (sdp any, err error) {
	p.conn, err = webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{{URLs: []string{"stun:stun1.l.google.com:19302"}}},
	})
	if err != nil {
		return "", err
	}
	p.conn.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		p.logger.Infof("Connection State has changed:  %s \n", connectionState.String())
		switch connectionState {
		case webrtc.ICEConnectionStateDisconnected:
		case webrtc.ICEConnectionStateFailed:
		case webrtc.ICEConnectionStateClosed:
			p.Disconnect()
		}
	})

	p.conn.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			singnal := make(map[string]string)
			singnal["signal"] = "Server ICE gathering is complete"
			sendICE(singnal)
			p.logger.Debug("ICE gathering is complete")
		}
		if candidate != nil {
			ice := candidate.ToJSON()
			if err := sendICE(ice); err != nil {
				p.logger.Error("Error while sending ICE candidate", err)
			}
			p.logger.Debug("ICE candidate found :", ice)
		}
	})
	// plug in the [video] track (out)
	video, err := newTrack("video", "pion", vCodec)
	if err != nil {
		p.logger.Error("Error while creating video track", err)
	}
	_, err = p.conn.AddTrack(video)
	if err != nil {
		p.logger.Error("Error while adding video track to Peer", err)
		return "", err
	}
	p.video = video
	p.logger.Debugf("Added [%s] track", video.Codec().MimeType)

	// plug in the [audio] track (out)
	audio, err := newTrack("audio", "pion", aCodec)
	if err != nil {
		p.logger.Error("Error while creating audio track", err)
	}
	_, err = p.conn.AddTrack(audio)
	if err != nil {
		p.logger.Error("Error while adding audio track to Peer", err)
	}
	p.audio = audio
	p.logger.Debugf("Added [%s] track", audio.Codec().MimeType)
	// plug in the [data] channel (in and out)
	if err = p.addDataChannel("input"); err != nil {
		return "", err
	}
	p.logger.Debug("Added [data] chan")

	offer, err := p.conn.CreateOffer(nil)
	if err != nil {
		p.logger.Error(err)
		return "", err
	}
	p.logger.Info("Server created offer")

	err = p.conn.SetLocalDescription(offer)
	if err != nil {
		p.logger.Error(err)
		return "", err
	}
	p.logger.Info("Server local description is set")

	return offer, nil
}

func (p *Peer) SetCandidatesAndSDP(data interface{}) {
	var candidate webrtc.ICECandidateInit
	var answer webrtc.SessionDescription
	switch d := data.(type) {
	case []byte:
		_ = json.Unmarshal(d, &candidate)
		_ = json.Unmarshal(d, &answer)
	default:
		p.logger.Errorf("invalid data type for SetCandidatesAndSDP: %T", data)
		return
	}

	if answer.SDP != "" {
		p.logger.Debug("Received answer:")
		p.logger.Debug(answer)
		_ = p.conn.SetRemoteDescription(answer)
	}
	if candidate.Candidate != "" {
		p.logger.Debug("Received candidate:")
		p.logger.Debug(candidate)
		_ = p.conn.AddICECandidate(candidate)
	}
}

func (p *Peer) SendData(data []byte) {
	err := p.data.Send(data)
	if err != nil {
		p.logger.Errorf("error occurred while data sending: %s", err)
	}
}

func (p *Peer) SendVideo() {

	//listener, err := socket.NewVideoUDPListener()
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 5004})
	if err != nil {
		p.logger.Errorf("error occurred while socket creating: %s", err)
	}
	defer listener.Close()
	listener.SetReadBuffer(1000000)
	rtpBuf := make([]byte, 5000)
	for {
		select {
		case <-p.signalStop:
			p.logger.Info("Stop video sending")
			listener.Close()
			listener = nil

			return
		default:
			n, _, err := listener.ReadFrom(rtpBuf)
			if err != nil {
				p.logger.Errorf("error occurred while data reading: %s", err)
				return
			}
			if _, err = p.video.Write(rtpBuf[:n]); err != nil {
				p.logger.Errorf("error occurred while video sending: %s", err)
				return
			}
		}
	}
}

func (p *Peer) SendAudio(data []byte) {}

func (p *Peer) Disconnect() {
	close(p.signalStop)
	time.Sleep(1 * time.Second)
	p.logger.Info("Disconnecting from peer")
	p.video = nil
	p.audio = nil
	p.data.Close()
	p.data = nil
	p.conn.Close()

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

func (p *Peer) addDataChannel(label string) error {
	ch, err := p.conn.CreateDataChannel(label, nil)
	if err != nil {
		return err
	}
	ch.OnOpen(func() {
		p.logger.Info("Data channel [input] opened")
		p.logger.Debugf("label: %s, id: %d", ch.Label(), ch.ID())
	})
	ch.OnError(func(err error) { p.logger.Error(err) })
	ch.OnMessage(func(m webrtc.DataChannelMessage) {
		if len(m.Data) == 0 {
			return
		}
		if p.OnMessage != nil {
			p.OnMessage(m.Data)
		}
	})
	p.data = ch
	ch.OnClose(func() { p.logger.Info("Data channel [input] closed") })
	return nil
}
