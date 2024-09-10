package webrtc

import (
	"cloud/internal/messages"
	"cloud/pkg/logger"
	"fmt"
	"strings"
	"time"

	"github.com/pion/webrtc/v4"
)

type WRTC struct {
	conn             *webrtc.PeerConnection
	logger           logger.Logger
	remoteCandidates []webrtc.ICECandidateInit
	OnMessage        func(data []byte)
	IsClosed         bool

	sysMes chan messages.Message
	errMes chan messages.AppError

	audio *webrtc.TrackLocalStaticRTP
	video *webrtc.TrackLocalStaticRTP
	data  *webrtc.DataChannel
}

func New(sys chan messages.Message, err chan messages.AppError) WRTC {
	return WRTC{
		remoteCandidates: make([]webrtc.ICECandidateInit, 0, 6),
		IsClosed:         false,
		sysMes:           sys,
		errMes:           err,
		logger:           logger.Init("7"),
	}
}
func (w *WRTC) Start(vCodec, aCodec string) {
	defer w.recovering()
	var err error
	w.conn, err = webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{{URLs: []string{"stun:stun1.l.google.com:19302"}}},
	})
	if err != nil {
		w.errMes <- *messages.NewAppError(err, "error occurred while creating peer connection", "", "")
		return
	}
	w.conn.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		w.recovering()
		w.logger.Infof("connection state has changed:  %s \n", connectionState.String())
		switch connectionState {
		case webrtc.ICEConnectionStateFailed:
			w.sysMes <- *messages.RtcConnectionClosed
			w.Stop()
		case webrtc.ICEConnectionStateClosed:
			w.sysMes <- *messages.RtcConnectionClosed
			w.Stop()
		case webrtc.ICEConnectionStateDisconnected:
			w.sysMes <- *messages.RtcConnectionClosed
			w.Stop()
		}
	})

	w.conn.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			w.sysMes <- *messages.NewMessage("", "", messages.RTC_SIGNAL, messages.RtcIceGatheringComplete)
		}
		if candidate != nil {
			w.sysMes <- *messages.NewMessage("", "", messages.RTC_SERVER_CANDIDATE, Encode(candidate.ToJSON()))
			w.logger.Debug("ice candidate found :", candidate)
		}
	})

	// plug in the [video] track (out)
	video, err := newTrack("video", "pion", vCodec)
	if err != nil {
		w.errMes <- *messages.NewAppError(err, "error occurred while creating video track", "", "")
		return
	}
	_, err = w.conn.AddTrack(video)
	if err != nil {
		w.errMes <- *messages.NewAppError(err, "error occurred while adding video track", "", "")
		return
	}
	w.video = video
	w.logger.Debugf("added [%s] track", video.Codec().MimeType)

	// plug in the [audio] track (out)
	audio, err := newTrack("audio", "pion", aCodec)
	if err != nil {
		w.errMes <- *messages.NewAppError(err, "error occurred while creating audio track", "", "")
		return
	}
	_, err = w.conn.AddTrack(audio)
	if err != nil {
		w.errMes <- *messages.NewAppError(err, "error occurred while adding audio track", "", "")
		return
	}
	w.audio = audio
	w.logger.Debugf("added [%s] track", audio.Codec().MimeType)

	// plug in the [data] channel (in and out)
	w.addDataChannel("input")

	w.logger.Debug("added [data] chan")

	offer, err := w.conn.CreateOffer(nil)
	if err != nil {
		w.errMes <- *messages.NewAppError(err, "error occurred while creating offer", "", "")
		return
	}
	w.logger.Info("server created offer")

	err = w.conn.SetLocalDescription(offer)
	if err != nil {
		w.errMes <- *messages.NewAppError(err, "error occurred while setting local description", "", "")
		return
	}
	w.logger.Info("server local description is set")

	w.sysMes <- *messages.NewMessage("", "", messages.RTC_OFFER, Encode(offer))
}

func (w *WRTC) SetAnswer(data string) {
	defer w.recovering()
	var answear webrtc.SessionDescription
	Decode(data, &answear)
	err := w.conn.SetRemoteDescription(answear)
	if err != nil {
		w.errMes <- *messages.NewAppError(err, "error occurred while setting remote description", "", "")
		return
	}
	w.logger.Info("remote description is set")
	time.Sleep(3 * time.Second)
	/*if (len(w.remoteCandidates)) > 0 {
		w.setCandidate()
	} else {
		w.logger.Warn("no candidates")
	}*/
}

/*
	func (w *WRTC) setCandidate() {
		for _, candidate := range w.remoteCandidates {
			w.logger.Debug("adding candidate: ", candidate)
			err := w.conn.AddICECandidate(candidate)
			if err != nil {
				w.errMes <- *messages.NewAppError(err, "error occurred while adding candidate", "", "")
			}
			w.logger.Info("candidate added")
		}
	}
*/
func (w *WRTC) AddCandidate(data string) {
	var candidate webrtc.ICECandidateInit
	Decode(data, &candidate)
	w.logger.Debug("received candidate: ", candidate)
	err := w.conn.AddICECandidate(candidate)
	if err != nil {
		w.errMes <- *messages.NewAppError(err, "error occurred while adding candidate", "", "")
	}
}

func (w *WRTC) SendData(data []byte) {

	defer w.recovering()
	err := w.data.Send(data)
	if err != nil {
		w.errMes <- *messages.NewAppError(err, "error occurred while sending data", "", "")
	}
}

func (w *WRTC) SendVideo(data []byte) {
	w.WRTCrecovering()
	_, err := w.video.Write(data)
	if err != nil {
		w.errMes <- *messages.NewAppError(err, "error occurred while sending video", "", "")
	}

}
func (w *WRTC) SendAudio(data []byte) {
	w.WRTCrecovering()
	_, err := w.audio.Write(data)
	if err != nil {
		w.errMes <- *messages.NewAppError(err, "error occurred while sending audio", "", "")
	}
}

func (w *WRTC) Stop() {
	if w.conn == nil {
		return
	}
	if w.conn.ConnectionState() < webrtc.PeerConnectionStateConnecting {
		_ = w.conn.Close()
	}
	w.IsClosed = true
	w.logger.Warn("webrtc stoped")
}

func newTrack(id string, label string, codec string) (*webrtc.TrackLocalStaticRTP, error) {
	fmt.Println("enter newTrack: ")
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
		case "h265", "hevc":
			mime = webrtc.MimeTypeH265
		case "av1":
			mime = webrtc.MimeTypeAV1
		case "vpx", "vp8":
			mime = webrtc.MimeTypeVP8
		}
		fmt.Println(mime)
	}
	if mime == "" {
		fmt.Println("unsupported codec")
		return nil, fmt.Errorf("unsupported codec %s:%s", id, codec)
	}

	return webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: mime}, id, label)
}

func (w *WRTC) addDataChannel(label string) {
	defer w.recovering()
	ch, err := w.conn.CreateDataChannel(label, nil)
	if err != nil {
		w.errMes <- *messages.NewAppError(err, fmt.Sprintf("error occurred while creating data channel %s", label), "", "")
		return
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
}

func (w *WRTC) recovering() {
	if r := recover(); r != nil {
		w.logger.Error(r)
		w.Stop()
	}
}
func (w *WRTC) WRTCrecovering() {
	if r := recover(); r != nil {
		w.logger.Error(r)
	}
}
