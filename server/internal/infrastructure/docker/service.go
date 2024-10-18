package docker

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

//! There were a few tries to build it another way;
//! Making one single container provides to freeze in webrtc stream for some reason;
//! Running containers via docker-compose provides conflicts with access to the same directories;
//! So, I decided to run each container separately with the same env variables and one network;

type vmService struct {
	pulseImage   string
	videoImage   string
	audioImage   string
	protoneImage string
	networkMode  string
	RendererPath string

	pulseEnv   []string
	videoEnv   []string
	audioEnv   []string
	protoneEnv []string
}

type VM interface {
	Start() error
	Stop()
}

func NewGameVmService(
	pulseImage, videoImage, audioImage, protoneImage, networkMode, rendererPath string,
	pulseEnv, videoEnv, audioEnv, protoneEnv []string) *vmService {

	if pulseEnv == nil {
		pulseEnv = []string{}
	}
	if videoEnv == nil {
		videoEnv = []string{}
	}
	if audioEnv == nil {
		audioEnv = []string{}
	}
	if protoneEnv == nil {
		protoneEnv = []string{}
	}
	return &vmService{
		pulseImage:   pulseImage,
		videoImage:   videoImage,
		audioImage:   audioImage,
		protoneImage: protoneImage,
		networkMode:  networkMode,
		RendererPath: rendererPath,
		pulseEnv:     pulseEnv,
		videoEnv:     videoEnv,
		audioEnv:     audioEnv,
		protoneEnv:   protoneEnv,
	}
}

// TODO: add xauthority path and remove hard code variables
func (f *vmService) ConfigureAndCreate(ctx context.Context,
	username, userHomePath, inputDevicePath, gamePath, display, hostIP string,
	planeID, videoPort, audioPort, pulseServerPort int) (dgi VM, err error) {
	//add temp pulse env
	pulseEnv := f.pulseEnv
	pulseEnv = append(pulseEnv, fmt.Sprintf("UNAME=%s", username))
	pulseEnv = append(pulseEnv, fmt.Sprintf("PULSE_SERVER_TCP_PORT=%d", pulseServerPort))
	//add temp video env
	videoEnv := f.videoEnv
	videoEnv = append(videoEnv, fmt.Sprintf("UNAME=%s", username))
	videoEnv = append(videoEnv, fmt.Sprintf("PLANE_ID=%d", planeID))
	videoEnv = append(videoEnv, fmt.Sprintf("HOST_IP=%s", hostIP))
	videoEnv = append(videoEnv, fmt.Sprintf("VIDEO_PORT=%d", videoPort))
	videoEnv = append(videoEnv, fmt.Sprintf("VAAPI_DEVICE=%s", f.RendererPath))
	//add temp audio env
	audioEnv := f.audioEnv
	audioEnv = append(audioEnv, fmt.Sprintf("UNAME=%s", username))
	audioEnv = append(audioEnv, fmt.Sprintf("HOST_IP=%s", hostIP))
	audioEnv = append(audioEnv, fmt.Sprintf("AUDIO_PORT=%d", audioPort))
	audioEnv = append(audioEnv, fmt.Sprintf("PULSE_SERVER_TCP_PORT=%d", pulseServerPort))
	//add temp protone env
	protoneEnv := f.protoneEnv
	protoneEnv = append(protoneEnv, fmt.Sprintf("UNAME=%s", username))
	protoneEnv = append(protoneEnv, fmt.Sprintf("DISPLAY=:%s", display))
	protoneEnv = append(protoneEnv, fmt.Sprintf("EXECUTIONFILE=%s", filepath.Base(gamePath)))
	protoneEnv = append(protoneEnv, fmt.Sprintf("PULSE_SERVER_TCP_PORT=%d", pulseServerPort))

	cfg := &ContainerGroupConfig{
		User:           username,
		PulseImage:     f.pulseImage,
		VideoImage:     f.videoImage,
		AudioImage:     f.audioImage,
		ProtoneImage:   f.protoneImage,
		NetworkMode:    f.networkMode,
		HomePath:       userHomePath,
		XauthorityPath: os.Getenv("XAUTHORITY"),
		RendererPath:   "/dev/dri/renderD128",
		CardPath:       "/dev/dri/card0",
		DevicePath:     inputDevicePath,
		GameDirPath:    filepath.Dir(gamePath),
		PulseEnv:       pulseEnv,
		VideoEnv:       videoEnv,
		AudioEnv:       audioEnv,
		ProtoneEnv:     protoneEnv,
	}

	fmt.Println(cfg)

	dg, err := newDockerGroup(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return dg, nil
}
