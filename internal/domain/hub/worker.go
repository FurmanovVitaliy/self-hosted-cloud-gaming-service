package hub

import (
	"cloud/internal/domain/games"
	"cloud/internal/domain/srm"
	"cloud/internal/messages"
	"cloud/pkg/input"
	"cloud/pkg/input/keymapping"
	"cloud/pkg/logger"
	"cloud/pkg/network/webrtc"
	"context"
	"fmt"
	"net"
	"runtime"

	"strconv"
)

type PlayerDevice struct {
	Display struct {
		Height float64 `json:"height,omitempty"`
		Width  float64 `json:"width,omitempty"`
	} `json:"display,omitempty"`
	Control struct {
		Gamepad  bool `json:"gamepad,omitempty"`
		Keyboard bool `json:"keyboard,omitempty"`
		Mouse    bool `json:"mouse,omitempty"`
		Touch    bool `json:"touch,omitempty"`
	} `json:"control,omitempty"`
}

type Worker struct {
	username                string
	roomUUID                string
	Message                 chan messages.Message
	ErrMes                  chan messages.AppError
	logger                  *logger.Logger
	websocket               WsManager
	webrtc                  webrtc.WRTC
	game                    games.Game
	playerDevice            PlayerDevice
	virtualDeviceKeymapping keymapping.XboxGpadInput
	virtualDevice           *input.VirtualGamepad
	ctx                     context.Context
	cancelFunc              context.CancelFunc

	resourceManager *srm.ServerResourceManager
	serverResources *srm.SererResources
}

func (w *Worker) Run(h *Hub) {
	w.recovering()
	go w.hendleMessages(h)
	go w.hansleWsMes()
	w.webrtc.Start("h264", "opus")
	w.serverResources, _ = w.resourceManager.AllocateResources()

}
func (w *Worker) Stop(h *Hub) {
	w.logger.Warn("stopping capture in goroutine of worker " + w.username)
	w.cancelFunc()

	w.logger.Warn("stopping virtual device of worker " + w.username)
	if w.virtualDevice.Device != nil {
		w.virtualDevice.Device.Close()
	}
	w.logger.Warn("stopping webrtc of worker " + w.username)
	if !w.webrtc.IsClosed {
		w.webrtc.Stop()
	}
	w.logger.Warn("stopping workerof worker " + w.username)
	if h != nil {
		h.DisconnectPlayer <- w
	}
	w.logger.Infof("worker %s stopped", w.username)
}

func (w *Worker) hendleMessages(h *Hub) {
	defer w.recovering()
	defer w.logger.Infof("message handler closed for %s worker", w.username)
	defer w.logger.Infof("error handler closed for %s worker", w.username)
	defer w.Stop(h)
hendleMessages:
	for {
		select {
		case msg := <-w.Message:
			switch msg.ContentType {
			case messages.RTC_OFFER:
				w.websocket.WriteJSON(msg)
			case messages.RTC_SERVER_CANDIDATE:
				w.websocket.WriteJSON(msg)
			case messages.RTC_SIGNAL:
				if msg.Content == messages.RtcIceGatheringComplete {
					w.websocket.WriteJSON(msg)
				}
				if msg == *messages.RtcConnectionClosed {
					w.logger.Warn("webrtc connection closed")
					break hendleMessages
				}
			}
		case err := <-w.ErrMes:
			w.logger.Error(err)
			break hendleMessages
		}
	}
}

func (w *Worker) hansleWsMes() {
	defer w.recovering()
	defer w.websocket.Close()
	defer w.logger.Info("ws closed")
readingWsMes:
	for {
		msg := w.websocket.ReadMessage()
		switch msg.ContentType {
		case messages.RTC_ANSWER:
			w.webrtc.SetAnswer(msg.Content)
		case messages.RTC_CLIENT_CANDIDATE:
			w.webrtc.AddCandidate(msg.Content)
		case messages.RTC_SIGNAL:

			if msg.Content == "interrupt" {
				w.logger.Warn("interrupt signal received")
				w.webrtc.Stop()
				break readingWsMes
			}
			if msg.Content == messages.RtcConnectionReady {
				w.runGame()
				go w.runVideoCapture()
				go w.runAudioCapture()
				break readingWsMes
			}
		case "deviceInfo":
			webrtc.Decode(msg.Content, &w.playerDevice)
			w.configureVirtualDevice()
			w.webrtc.OnMessage = w.handleInput
		}
	}
}

func (w *Worker) runGame() {
	_, id := w.resourceManager.ConfigureAndStartVm(
		w.username,
		w.game.Path,
		"/home/vitalii/Dev/arch-containrer/simple-arch-container/LOCAL_STARAGE",
		w.serverResources.XServer.ScreenNumber,
		strconv.Itoa(w.serverResources.Listener.Audio.LocalAddr().(*net.UDPAddr).Port),
		strconv.Itoa(w.serverResources.Listener.Video.LocalAddr().(*net.UDPAddr).Port),
		strconv.Itoa(w.serverResources.XServer.Plane),
		w.virtualDevice.Path,
	)

	w.serverResources.VMID = id
}

func (w *Worker) runVideoCapture() {
	defer w.recovering()

	rtpBuf := make([]byte, 10000000)
	for {
		select {
		case <-w.ctx.Done():
			fmt.Println("stop video capture")
			return // Остановка горутины при отмене контекста
		default:
			n, _, err := w.serverResources.Listener.Video.ReadFrom(rtpBuf)
			if err != nil {
				w.ErrMes <- *messages.NewAppError(err, "error occurred while reading data from udp listener", "", "")
				return
			}
			w.webrtc.SendVideo(rtpBuf[:n])
		}
	}
}

func (w *Worker) runAudioCapture() {
	defer w.recovering()

	rtpBuf := make([]byte, 5000)
	for {
		select {
		case <-w.ctx.Done():
			fmt.Println("stop audio capture")
			return
		default:
			n, _, err := w.serverResources.Listener.Audio.ReadFrom(rtpBuf)
			if err != nil {
				w.ErrMes <- *messages.NewAppError(err, "error occurred while reading data from udp listener", "", "")
				return
			}
			w.webrtc.SendAudio(rtpBuf[:n])
		}
	}
}

func (w *Worker) configureVirtualDevice() {
	if w.playerDevice.Control.Gamepad {
		var err error
		w.virtualDevice, err = input.CreateGamepad(w.username)
		if err != nil {
			w.ErrMes <- *messages.ErrCreateVirtualDevice
		}
	}
}

func (w *Worker) handleInput(data []byte) {
	w.virtualDevice.HandleInput(data, w.virtualDeviceKeymapping)
}

func (w *Worker) recovering() {
	if r := recover(); r != nil {
		// Создаем срез байт для хранения стека вызовов
		buf := make([]byte, 1024)
		// Читаем стек вызовов
		n := runtime.Stack(buf, false)
		// Логируем ошибку и стек вызовов
		w.logger.Printf("Recovered from panic: %v\nStack trace: %s", r, buf[:n])
	}
}
