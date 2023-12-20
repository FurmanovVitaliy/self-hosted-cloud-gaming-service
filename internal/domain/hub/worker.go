package hub

import (
	"cloud/internal/domain/games"
	"cloud/internal/messages"
	"cloud/pkg/input/keymapping"
	"cloud/pkg/logger"
	"cloud/pkg/network/webrtc"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os/exec"
	"strconv"

	"github.com/bendahl/uinput"
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
	virtualDevice           uinput.Gamepad
	externalProces          [2]*exec.Cmd
	listener                *net.UDPConn
}

func (w *Worker) Run(h *Hub) {
	w.recovering()
	go w.hendleMessages(h)
	go w.hansleWsMes()
	w.webrtc.Start("h264", "opus")

}
func (w *Worker) Stop(h *Hub) {
	w.logger.Warn("closing listener of worker " + w.username)
	if w.listener != nil {
		w.listener.Close()
	}
	w.logger.Warn("stopping virtual device of worker " + w.username)
	if w.virtualDevice != nil {
		w.virtualDevice.Close()
	}
	w.logger.Warn("stopping webrtc of worker " + w.username)
	if !w.webrtc.IsClosed {
		w.webrtc.Stop()
	}
	w.logger.Warn("stopping workerof worker " + w.username)
	if h != nil {
		h.DisconnectPlayer <- w
	}
	w.logger.Warn("killingexternal processesof worker " + w.username)
	for _, cmd := range w.externalProces {
		if cmd != nil {
			w.logger.Info("kill process" + cmd.String())
			cmd.Process.Kill()
		}
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
				//w.runGame()
				go w.runCapture()
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
	protonePath := "/home/vitalii/.PP/PortProton/data/scripts/start.sh"
	cmd := exec.Command(protonePath, w.game.Path)
	w.externalProces[0] = cmd
	err := cmd.Start()
	if err != nil {
		w.ErrMes <- *messages.ErrStartGame
	}
}

func (w *Worker) runCapture() {
	defer w.recovering()
	defer func() {
		if w.listener != nil {
			w.listener.Close()
		}
	}()

	var number int
	for {
		// Генерируем случайное четырехзначное число больше 5000
		number = rand.Intn(5000) + 5001
		if number%2 == 0 {
			break // Выходим из цикла, если число четное
		}
	}

	ip := "127.0.0.1"
	port := number

	scriptPath := "/home/vitalii/dev/kms_ffmpeg/ffmpeg+ksm.sh"
	cmd := exec.Command("sudo", scriptPath, strconv.Itoa(port))
	w.externalProces[1] = cmd
	err := cmd.Start()
	if err != nil {
		w.ErrMes <- *messages.ErrStartCapture
	}

	w.listener, err = net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP(ip), Port: port})
	if err != nil {
		w.ErrMes <- *messages.NewAppError(err, fmt.Sprintf("error occurred while creating udp listener on %s:%d", ip, port), "", "")
		return
	}

	if err = w.listener.SetReadBuffer(1000000); err != nil {
		w.ErrMes <- *messages.NewAppError(err, "error occurred while setting read buffer", "", "")
		return
	}

	rtpBuf := make([]byte, 5000)
readBufData:
	for {
		n, _, err := w.listener.ReadFrom(rtpBuf)
		if err != nil {
			w.ErrMes <- *messages.NewAppError(err, "error occurred while reading data from udp listener", "", "")
			break readBufData
		}
		w.webrtc.SendVideo(rtpBuf[:n])
	}
}

func (w *Worker) configureVirtualDevice() {
	if w.playerDevice.Control.Gamepad {
		var err error
		w.virtualDevice, err = uinput.CreateGamepad("/dev/uinput", []byte("testpad"), 045, 955)
		if err != nil {
			w.ErrMes <- *messages.ErrCreateVirtualDevice
		}
	}
}

func (w *Worker) handleInput(data []byte) {

	json.Unmarshal(data, &w.virtualDeviceKeymapping)

	x, _ := strconv.ParseFloat(w.virtualDeviceKeymapping.RS.X, 32)
	y, _ := strconv.ParseFloat(w.virtualDeviceKeymapping.RS.Y, 32)

	w.virtualDevice.RightStickMove(float32(x), float32(y))

	lsX, _ := strconv.ParseFloat(w.virtualDeviceKeymapping.LS.X, 32)
	lsY, _ := strconv.ParseFloat(w.virtualDeviceKeymapping.LS.Y, 32)
	w.virtualDevice.LeftStickMove(float32(lsX), float32(lsY))

	if w.virtualDeviceKeymapping.A == 1 {
		w.virtualDevice.ButtonDown(uinput.ButtonSouth)
	} else {
		w.virtualDevice.ButtonUp(uinput.ButtonSouth)
	}
	if w.virtualDeviceKeymapping.X == 1 {
		w.virtualDevice.ButtonDown(uinput.ButtonWest)
	} else {
		w.virtualDevice.ButtonUp(uinput.ButtonWest)
	}
	if w.virtualDeviceKeymapping.Y == 1 {
		w.virtualDevice.ButtonDown(uinput.ButtonNorth)
	} else {
		w.virtualDevice.ButtonUp(uinput.ButtonNorth)
	}
	if w.virtualDeviceKeymapping.LB == 1 {
		w.virtualDevice.ButtonDown(uinput.ButtonBumperLeft)
	} else {
		w.virtualDevice.ButtonUp(uinput.ButtonBumperLeft)
	}
	if w.virtualDeviceKeymapping.RB == 1 {
		w.virtualDevice.ButtonDown(uinput.ButtonBumperRight)
	} else {
		w.virtualDevice.ButtonUp(uinput.ButtonBumperRight)
	}
	if w.virtualDeviceKeymapping.Start == 1 {
		w.virtualDevice.ButtonDown(uinput.ButtonStart)
	} else {
		w.virtualDevice.ButtonUp(uinput.ButtonStart)
	}
	if w.virtualDeviceKeymapping.Main == 1 {
		w.virtualDevice.ButtonDown(uinput.ButtonMode)
	} else {
		w.virtualDevice.ButtonUp(uinput.ButtonMode)
	}

}

func (w *Worker) recovering() {
	if r := recover(); r != nil {
		w.logger.Error(r)
	}
}
