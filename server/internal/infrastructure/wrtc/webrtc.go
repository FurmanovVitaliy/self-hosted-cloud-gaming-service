package wrtc

import (
	"fmt"
	"strings"

	"github.com/FurmanovVitaliy/pixel-cloud/pkg/logger"
	"github.com/pion/webrtc/v4"
)

type WRTC struct {
	conn      *webrtc.PeerConnection
	OnMessage func(data []byte)

	logger     logger.Logger
	IsClosed   bool
	IsAnswered bool

	aCodec string
	vCodec string

	audio *webrtc.TrackLocalStaticRTP
	video *webrtc.TrackLocalStaticRTP
	data  *webrtc.DataChannel

	candidateBuffer []webrtc.ICECandidateInit
}

func newWRTC(OnMessage func(data []byte), vCodec, aCodec string, logger logger.Logger) *WRTC {
	return &WRTC{
		aCodec:    aCodec,
		vCodec:    vCodec,
		OnMessage: OnMessage,
		logger:    logger,
	}
}

func (w *WRTC) Init(OnIceCandidate func(string)) (offer string, err error) {
	w.conn, err = webrtc.NewPeerConnection(webrtc.Configuration{
		//ICETransportPolicy: webrtc.ICETransportPolicyAll,
		//BundlePolicy:       webrtc.BundlePolicyMaxBundle,
		//RTCPMuxPolicy:      webrtc.RTCPMuxPolicyRequire,
		ICEServers: []webrtc.ICEServer{
			{URLs: []string{"stun:stun1.l.google.com:19302"}},
			{URLs: []string{"stun:stun.l.google.com:19302"}},
			{URLs: []string{"stun:stun.l.google.com:5349"}},
			{URLs: []string{"stun:stun1.l.google.com:3478"}},
			{URLs: []string{"stun:stun1.l.google.com:5349"}},
		},
	})
	if err != nil {
		return "", fmt.Errorf("error occurred while creating peer connection: %w", err)
	}
	w.conn.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		w.logger.Infof("Connection state changed: %s", connectionState.String())
		switch connectionState {
		case webrtc.ICEConnectionStateFailed, webrtc.ICEConnectionStateClosed, webrtc.ICEConnectionStateDisconnected:
			w.conn.Close()
		}
	})

	w.conn.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			w.logger.Debug("ICE gathering complete")
			OnIceCandidate("")
			return
		}
		w.logger.Debug("ice candidate found :", candidate)
		OnIceCandidate(Encode(candidate.ToJSON()))
	})

	// plug in the [video] track (out)
	video, err := newTrack("video", "pion", w.vCodec)
	if err != nil {
		return "", fmt.Errorf("error occurred while creating video track: %w", err)
	}
	_, err = w.conn.AddTrack(video)
	if err != nil {
		return "", fmt.Errorf("error occurred while adding video track: %w", err)
	}
	w.video = video
	w.logger.Debugf("added [%s] track", video.Codec().MimeType)

	// plug in the [audio] track (out)
	audio, err := newTrack("audio", "pion", w.aCodec)
	if err != nil {
		return "", fmt.Errorf("error occurred while creating audio track: %w", err)
	}
	_, err = w.conn.AddTrack(audio)
	if err != nil {
		return "", fmt.Errorf("error occurred while adding audio track: %w", err)
	}
	w.audio = audio
	w.logger.Debugf("added [%s] track", audio.Codec().MimeType)

	// plug in the [data] channel (in and out)
	err = w.addDataChannel("input")
	if err != nil {
		return "", fmt.Errorf("error occurred while adding data channel: %w", err)
	}

	w.logger.Debug("added [data] chan")

	o, err := w.conn.CreateOffer(nil)
	if err != nil {
		return "", fmt.Errorf("error occurred while creating offer: %w", err)
	}

	w.logger.Debug("server offer created")

	err = w.conn.SetLocalDescription(o)
	if err != nil {
		return "", fmt.Errorf("error occurred while setting local description: %w", err)
	}
	w.logger.Info("server local description is set")

	offer = Encode(o)
	return
}

func (w *WRTC) SetAnswer(data string) (err error) {
	var answer webrtc.SessionDescription
	Decode(data, &answer)
	err = w.conn.SetRemoteDescription(answer)
	if err != nil {
		return fmt.Errorf("error occurred while setting remote description: %w", err)
	}
	w.logger.Debug("remote description is set")
	w.IsAnswered = true

	// add buffered candidates
	for _, candidate := range w.candidateBuffer {
		err = w.conn.AddICECandidate(candidate)
		if err != nil {
			return fmt.Errorf("error occurred while adding buffered candidate: %w", err)
		}
		w.logger.Debug("buffered candidate added")
	}
	w.candidateBuffer = nil // clear the buffer

	return
}

func (w *WRTC) AddCandidate(data string) (err error) {
	var candidate webrtc.ICECandidateInit
	Decode(data, &candidate)
	w.logger.Debugf("Received candidate: %+v", candidate)

	// if answer is not set, buffer the candidate
	if !w.IsAnswered {
		w.logger.Debug("Buffering candidate as answer is not yet set")
		w.candidateBuffer = append(w.candidateBuffer, candidate)
		return
	}

	// if answer is set, add the candidate
	err = w.conn.AddICECandidate(candidate)
	if err != nil {
		return fmt.Errorf("error occurred while adding candidate: %w", err)
	}
	w.logger.Debug("Candidate added")
	return
}

func (w *WRTC) SendData(data []byte) error {
	err := w.data.Send(data)
	if err != nil {
		return fmt.Errorf("error occurred while sending data: %w", err)
	}
	return nil
}

func (w *WRTC) SendVideo(data []byte) error {
	//w.logger.Debugf("Sending video of size: %d bytes", len(data))
	_, err := w.video.Write(data)
	if err != nil {
		w.logger.Errorf("Error while sending video: %v", err)
		return fmt.Errorf("error occurred while sending video: %w", err)
	}
	return nil
}

func (w *WRTC) SendAudio(data []byte) error {
	//w.logger.Debugf("Sending audio of size: %d bytes", len(data))
	_, err := w.audio.Write(data)
	if err != nil {
		w.logger.Errorf("Error while sending audio: %v", err)
		return fmt.Errorf("error occurred while sending audio: %w", err)
	}
	return nil
}

func (w *WRTC) Stop() {
	if w.conn == nil {
		return
	}
	if w.conn.ConnectionState() < webrtc.PeerConnectionStateConnecting {
		_ = w.conn.Close()
	}
	w.IsClosed = true
	w.logger.Info("webrtc stoped")
}

func newTrack(id string, label string, codec string) (*webrtc.TrackLocalStaticRTP, error) {
	var mime string
	var clockRate uint32
	var channels uint16
	var sdpFmtpLine string
	var rtcpFeedback []webrtc.RTCPFeedback

	codec = strings.ToLower(codec)
	switch id {
	case "audio":
		switch codec {
		case "opus":
			mime = webrtc.MimeTypeOpus
			clockRate = 48000
			channels = 2
			sdpFmtpLine = "minptime=10; useinbandfec=1"
		default:
			return nil, fmt.Errorf("unsupported audio codec %s", codec)
		}
	case "video":
		switch codec {
		case "h264":
			mime = webrtc.MimeTypeH264
			clockRate = 90000
			sdpFmtpLine = "profile-level-id=42e01f;packetization-mode=1"
			rtcpFeedback = []webrtc.RTCPFeedback{
				{Type: "nack", Parameter: ""},
				{Type: "nack", Parameter: "pli"},
			}
		case "h265", "hevc":
			mime = webrtc.MimeTypeH265
			clockRate = 90000
			sdpFmtpLine = "profile-id=1;tier=main;level-id=1"
			rtcpFeedback = []webrtc.RTCPFeedback{
				{Type: "nack", Parameter: ""},
				{Type: "nack", Parameter: "pli"},
			}
		case "av1":
			mime = webrtc.MimeTypeAV1
			clockRate = 90000
			sdpFmtpLine = "profile-id=0;level-id=1"
		case "vpx", "vp8":
			mime = webrtc.MimeTypeVP8
			clockRate = 90000
			sdpFmtpLine = "profile-id=0"
			rtcpFeedback = []webrtc.RTCPFeedback{
				{Type: "nack", Parameter: ""},
				{Type: "nack", Parameter: "pli"},
			}
		default:
			return nil, fmt.Errorf("unsupported video codec %s", codec)
		}
	default:
		return nil, fmt.Errorf("unsupported track type %s", id)
	}

	// Создание RTPCodecCapability с настройками
	codecCapability := webrtc.RTPCodecCapability{
		MimeType:     mime,
		ClockRate:    clockRate,
		Channels:     channels,
		SDPFmtpLine:  sdpFmtpLine,
		RTCPFeedback: rtcpFeedback,
	}

	return webrtc.NewTrackLocalStaticRTP(codecCapability, id, label)
}

func (w *WRTC) addDataChannel(label string) error {
	ch, err := w.conn.CreateDataChannel(label, nil)
	if err != nil {
		return fmt.Errorf("error occurred while creating data channel: %w", err)
	}
	ch.OnOpen(func() {
		w.logger.Info("data channel [input] opened")
		w.logger.Debugf("label: %s, id: %d", ch.Label(), ch.ID())
	})
	ch.OnError(func(err error) { w.logger.Warn(err) })
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
