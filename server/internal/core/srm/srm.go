package srm

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/FurmanovVitaliy/pixel-cloud/internal/infrastructure/display"
	"github.com/FurmanovVitaliy/pixel-cloud/internal/infrastructure/input"
	"github.com/FurmanovVitaliy/pixel-cloud/internal/infrastructure/listener"
	"github.com/FurmanovVitaliy/pixel-cloud/pkg/logger"
)

var _ InputDevice = &input.InputDevice{}
var _ UDPReader = &listener.UDPReader{}
var _ Display = &display.XServer{}

type RMconfig struct {
	fsInitialDirectoryPath string
	userDirectoriesPath    string
	minUdpPortRange        int
	maxUdpPortRange        int
	udpBufferSize          int
	udpReadBuffer          int
}
type ResourceManager struct {
	displays []display.XServer
	config   *RMconfig
	logger   *logger.Logger
}

type Resources struct {
	Display      *display.XServer
	VideoReader  *listener.UDPReader
	AudioReader  *listener.UDPReader
	InputDevice  *input.InputDevice
	UserHomePath string
}

func CreateParams(fsInitialDirectoryPath, userDirectoriesPath string, minUdpPortRange, maxUdpPortRange, udpBufferSize, udpReadBuffer int) *RMconfig {
	return &RMconfig{
		fsInitialDirectoryPath: fsInitialDirectoryPath,
		userDirectoriesPath:    userDirectoriesPath,
		minUdpPortRange:        minUdpPortRange,
		maxUdpPortRange:        maxUdpPortRange,
		udpBufferSize:          udpBufferSize,
		udpReadBuffer:          udpReadBuffer,
	}
}

func New(config *RMconfig, displays []display.XServer, logger *logger.Logger) *ResourceManager {
	return &ResourceManager{
		displays: displays,
		config:   config,
		logger:   logger,
	}
}

func (rm *ResourceManager) AllocateListeners(ip string) ([2]UDPReader, int, error) {

	var aReader, vReader UDPReader
	var al, vl *net.UDPConn
	var port int
	var err error

	for port = rm.config.minUdpPortRange; port <= rm.config.maxUdpPortRange-2; port++ {
		al, err = net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP(ip), Port: port + 2})
		if err != nil {
			continue
		}

		vl, err = net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP(ip), Port: port})
		if err == nil {
			break
		}

		al.Close()
	}

	if err != nil {
		return [2]UDPReader{}, 0, fmt.Errorf("No free ports in range %d-%d", rm.config.minUdpPortRange, rm.config.maxUdpPortRange)
	}

	if err = al.SetReadBuffer(rm.config.udpBufferSize); err != nil {
		return [2]UDPReader{}, 0, err
	}

	if err = vl.SetReadBuffer(rm.config.udpBufferSize); err != nil {
		return [2]UDPReader{}, 0, err
	}

	aReader = listener.NewUdpReader(port+2, uint(rm.config.udpReadBuffer), al)
	vReader = listener.NewUdpReader(port, uint(rm.config.udpReadBuffer), vl)

	return [2]UDPReader{aReader, vReader}, port + 1, nil
}

func (rm *ResourceManager) AllocateDisplay() (Display, error) {
	for i := range rm.displays {
		if !rm.displays[i].Used {
			rm.displays[i].Used = true
			return &rm.displays[i], nil
		}
	}
	return nil, fmt.Errorf("no available xservers")
}

func (rm *ResourceManager) AllocateDiscSpace(username string) (string, error) {

	userHome := filepath.Join(rm.config.userDirectoriesPath, username+"-home")

	if _, err := os.Stat(userHome); !os.IsNotExist(err) {
		rm.logger.Debugf("User folder found: %v", userHome)
		return userHome, nil
	}

	if err := os.MkdirAll(userHome, 0777); err != nil {
		return "", fmt.Errorf("error user folder creation: %v", err)
	}
	rm.logger.Debugf("User folder created: %v", userHome)

	source := rm.config.fsInitialDirectoryPath + "/"
	dest := userHome + "/"
	fmt.Println("source: ", source)
	fmt.Println("dest: ", dest)

	cmd := exec.Command("rsync", "-a", source, dest)

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed wfile file base rsync: %v", err)
	}

	return userHome, nil
}

func (rm *ResourceManager) AllocateInputDevice(username, vendorID, productID, productT string) (device InputDevice, err error) {
	var T input.DeviceType //TODO: refactor
	switch productT {
	case "keyboard":
		T = input.Keyboard
	case "mouse":
		T = input.Mouse
	case "gamepad":
		T = input.Gamepad
	default:
		return nil, fmt.Errorf("unknown device type: %v", productT)
	}
	device, err = input.NewInputDevice(username, vendorID, productID, T)
	if err != nil {
		return nil, err
	}
	rm.logger.Debugf("Input device allocated: %v", device)
	return device, nil
}
