package hub

import (
	"cloud/internal/domain/games"
	"cloud/pkg/input/keymapping"
	"cloud/pkg/logger"
	"cloud/pkg/network/webrtc"
	"encoding/json"
	"os"
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
	RegMes                  chan Message
	ErrMes                  chan error
	logger                  *logger.Logger
	websocket               WsManager
	webrtc                  *webrtc.WRTC
	game                    games.Game
	PlayerDevice            PlayerDevice
	virtualDeviceKeymapping keymapping.XboxGpadInput
	virtualDevice           uinput.Gamepad
	appPid                  int
	capturePid              int
}

func (w *Worker) Run(h *Hub) {
	go w.handleErrorMes(h)
	go w.handleIncomeMes(h)
	w.webrtc.Start("h264", "opus", w.websocket.WriteJSON)
	w.runGame()
	w.runCapture()
	w.configureVirtualDevice()
	go w.webrtc.SendVideo("127.0.0.1", 5004)
	w.webrtc.OnMessage = w.handleInput
}

func (w *Worker) handleIncomeMes(*Hub) {
	defer w.websocket.Close()
	defer w.logger.Info("ws closed")
readingWsMes:
	for {
		msg, err := w.websocket.ReadMessage()
		if err != nil {
			w.ErrMes <- err
			return
		}
		switch msg.ContentType {
		case webrtc.ANSWER:
			err = w.webrtc.SetAnswer(msg.Content)
			if err != nil {
				w.ErrMes <- err
				return
			}
		case webrtc.CLIENT_CANDIDATE:
			err = w.webrtc.AddCandidate(msg.Content)
			if err != nil {
				w.ErrMes <- err
				return
			}
		case webrtc.SIGNAL:
			if msg.Content == webrtc.CONNECTION_READY {
				break readingWsMes
			}
		case "deviceInfo":
			webrtc.Decode(msg.Content, &w.PlayerDevice)
		}
	}
}

func (w *Worker) handleErrorMes(h *Hub) {
	defer w.logger.Infof("error handler closed for %s worker", w.username)
	for err := range w.ErrMes {
		w.logger.Error(err)
		w.logger.Infof("%s worker disconnected", w.username)
		w.websocket.Close()
		w.webrtc.Stop()
		app, err := os.FindProcess(w.appPid)
		if err == nil {
			app.Kill()
		}
		capture, err := os.FindProcess(w.capturePid)
		if err == nil {
			capture.Kill()
		}
		h.DisconnectPlayer <- w
		h.Broadcast <- &Message{
			RoomUUID:    w.roomUUID,
			Username:    w.username,
			ContentType: "broadcast",
			Content:     w.username + " left the room",
		}
		return
	}
}

func (w *Worker) runGame() {
	protonePath := "/home/vitalii/.PP/PortProton/data/scripts/start.sh"
	cmd := exec.Command(protonePath, w.game.Path)
	err := cmd.Start()
	if err != nil {
		w.ErrMes <- err
	}
	w.appPid = cmd.Process.Pid
}
func (w *Worker) runCapture() {
	scriptPath := "/home/vitalii/dev/kms_ffmpeg/ffmpeg+ksm.sh"
	cmd := exec.Command("sudo", scriptPath)
	err := cmd.Start()
	if err != nil {
		w.ErrMes <- err
	}
}

func (w *Worker) configureVirtualDevice() {
	if w.PlayerDevice.Control.Gamepad {
		var err error
		w.virtualDevice, err = uinput.CreateGamepad("/dev/uinput", []byte("testpad"), 045, 955)
		if err != nil {
			w.ErrMes <- err
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
