package srm

import (
	"fmt"
	"net"
)

type SererResources struct {
	XServer  *xServer
	Listener listener
	VMID     string
}
type listener struct {
	Video *net.UDPConn
	Audio *net.UDPConn
	port  int
}
type xServer struct {
	ScreenNumber string `json:"screen"`
	Card         string `json:"card"`
	Port         string `json:"port_name"`
	Connector    int    `json:"connector_id"`
	Plane        int    `json:"plane_id"`
	Used         bool
}

type ServerResourceManager struct {
	Dockers   []*docker
	xServers  []xServer
	listeners []*listener
}

func NewResourceManager() *ServerResourceManager {
	initXservers()
	xservers := jsonToXserver()

	fmt.Println("Xservers:", xservers)

	srm := &ServerResourceManager{
		Dockers:  make([]*docker, 0),
		xServers: xservers,
	}
	return srm
}

func (srm *ServerResourceManager) ConfigureAndStartVm(username, gamePath, localStorage, display, portAudio, portVideo, planeId, devicePath string) (error, string) {
	d := newDocker(username, gamePath, localStorage, display, portAudio, portVideo, planeId, devicePath, "arch:simple")
	d.configureDocker()
	d.startDocker()
	srm.Dockers = append(srm.Dockers, d)
	return nil, d.ID
}

func (srm *ServerResourceManager) StopVM(dockerID string) error {
	for i, d := range srm.Dockers {
		if d.ID == dockerID {
			srm.releaseXserver(d.display)
			d.stopDocker()
			srm.Dockers = append(srm.Dockers[:i], srm.Dockers[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("docker with ID %s not found", dockerID)
}

func (srm *ServerResourceManager) AllocateResources() (*SererResources, error) {
	xServer, err := srm.allocateXserver()
	if err != nil {
		return nil, err
	}

	listener, err := srm.allocateListener()
	if err != nil {
		return nil, err
	}

	return &SererResources{
		XServer:  xServer,
		Listener: *listener,
		VMID:     "/dev/dri/card",
	}, nil
}

func (srm *ServerResourceManager) ReleaseResources(sr *SererResources) {
	srm.releaseXserver(sr.XServer.ScreenNumber)
	srm.releaseListener(sr.Listener.port)

}

/************************* not public functions ***************************/
func (srm *ServerResourceManager) allocateXserver() (*xServer, error) {
	for i := range srm.xServers {
		if !srm.xServers[i].Used {
			srm.xServers[i].Used = true
			return &srm.xServers[i], nil
		}
	}
	return nil, fmt.Errorf("no free Xservers")
}

func (srm *ServerResourceManager) releaseXserver(screenNumber string) {
	for i := range srm.xServers {
		if srm.xServers[i].ScreenNumber == screenNumber {
			srm.xServers[i].Used = false
			return
		}
	}

}

func (srm *ServerResourceManager) allocateListener() (*listener, error) {
	var err error
	ip := "127.0.0.1"
	port := randomListenerPort()

	for i := 0; i < len(srm.listeners); i++ {
		if srm.listeners[i].port == port {
			port = randomListenerPort()
			i = 0
		}
	}

	listener := listener{}
	listener.Video, err = net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP(ip), Port: port - 2})
	if err != nil {
		return nil, err
	}
	err = listener.Video.SetReadBuffer(10000000)
	if err != nil {
		return nil, err
	}
	listener.Audio, err = net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP(ip), Port: port})
	if err != nil {
		return nil, err
	}
	err = listener.Audio.SetReadBuffer(100000)
	if err != nil {
		return nil, err
	}
	listener.port = port
	srm.listeners = append(srm.listeners, &listener)
	return &listener, nil

}
func (srm *ServerResourceManager) releaseListener(port int) {
	for i := range srm.listeners {
		if srm.listeners[i].port == port {
			err := srm.listeners[i].Video.Close()
			if err != nil {
				fmt.Println("Error while closing video listener:", err)
			}
			fmt.Println("Video listener closed")
			srm.listeners[i].Audio.Close()
			if err != nil {
				fmt.Println("Error while closing audio listener:", err)
			}
			fmt.Println("Audio listener closed")
			srm.listeners = append(srm.listeners[:i], srm.listeners[i+1:]...)
			return
		}
	}
}
