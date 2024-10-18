package hub

import (
	"github.com/FurmanovVitaliy/pixel-cloud/pkg/logger"
)

type Worker struct {
	Username    string
	roomUUID    string
	host        bool
	msgHandler  MessageHandler
	streamer    WrtcStreamer
	chatMsg     chan *ChatMsg
	wrtcMsg     chan *WrtcMsg
	audioReader UDPReader
	videoReader UDPReader
	vm          VM
	logger      logger.Logger
}

func (w *Worker) WriteChatMsg(msg *ChatMsg) {
	if err := w.msgHandler.Write("chat", msg); err != nil {
		w.logger.Error("error occurred while writing chat message: %w", err)
	}
}

func (w *Worker) ReadChatMsg(h *Hub) {
	defer func() {
		w.logger.Warn("chat message reader stopped")
		h.Disconnect <- w
	}()
	for {
		msg, ok := <-w.chatMsg
		if !ok {
			return
		}
		h.Broadcast <- msg
	}
}

func (w *Worker) ReadStreamerMsg(h *Hub) {
	defer func() {
		w.logger.Warn("streamer message reader stopped")
		w.streamer.Stop()
		h.Disconnect <- w
	}()

	offer, err := w.streamer.Init(w.streamerOnIceCandidate)
	if err != nil {
		w.logger.Error("error occurred while initializing streamer: %w", err)
		return
	}
	if err := w.msgHandler.Write("wrtc", &WrtcMsg{ContentType: "offer", Content: offer}); err != nil {
		w.logger.Error("error occurred while writing offer: %w", err)
		return
	}
	for {
		msg, ok := <-w.wrtcMsg
		if !ok {
			return
		}
		switch msg.ContentType {
		case "answer":
			if err := w.streamer.SetAnswer(msg.Content); err != nil {
				w.logger.Error("error occurred while setting answer: %w", err)
				return
			}
		case "candidate":
			if err := w.streamer.AddCandidate(msg.Content); err != nil {
				w.logger.Error("error occurred while adding candidate: %w", err)
				return
			}
		case "connection_ready":
			w.logger.Info("connection is ready")
			go w.streamerProcesAudio()
			go w.streamerProcesVideo()
		}
	}
}

func (w *Worker) streamerOnIceCandidate(candidate string) {
	if candidate == "" {
		if err := w.msgHandler.Write("wrtc", &WrtcMsg{ContentType: "server_ice_ready"}); err != nil {
			w.logger.Error("error occurred while writing server_ice_ready: %w", err)
		}
		return
	}
	if err := w.msgHandler.Write("wrtc", &WrtcMsg{ContentType: "candidate", Content: candidate}); err != nil {
		w.logger.Error("error occurred while writing candidate: %w", err)
	}
}

func (w *Worker) streamerProcesAudio() {
	defer func() {
		w.logger.Info("streamer audio reader stopped")
		w.audioReader.Close()
	}()
	w.logger.Info("sending audio data started")
	for {
		data, err := w.audioReader.Read()
		if err != nil {
			w.logger.Error("error occurred while reading audio data: %w", err)
			return
		}
		if err := w.streamer.SendAudio(data); err != nil {
			w.logger.Error("error occurred while sending audio data: %w", err)
			return
		}
	}

}
func (w *Worker) streamerProcesVideo() {
	defer func() {
		w.logger.Info("streamer video reader stopped")
		w.videoReader.Close()
	}()
	w.logger.Info("sending video data started")
	for {
		data, err := w.videoReader.Read()
		if err != nil {
			w.logger.Error("error occurred while reading video data: %w", err)
			return
		}
		if err := w.streamer.SendVideo(data); err != nil {
			w.logger.Error("error occurred while sending video data: %w", err)
			return
		}
	}
}

func (w *Worker) SrartVM(h *Hub) {
	defer func() {
		w.logger.Warn("VM stopped")
		h.Disconnect <- w
	}()
	w.logger.Info("starting VM")
	if err := w.vm.Start(); err != nil {
		w.logger.Error("error occurred while starting VM: %w", err)
		return
	}
}
