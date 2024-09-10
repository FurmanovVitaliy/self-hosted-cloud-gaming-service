package srm

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

type docker struct {
	username     string
	display      string
	gamePath     string
	portAudio    string
	portVideo    string
	planeID      string
	localStorage string
	image        string
	devicePath   string
	ID           string
	client       *client.Client
}

func newDocker(username, gamePath, localStorage, display, portAudio, portVideo, planeId, devicePath, image string) *docker {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}
	docker := &docker{
		username:     username,
		display:      display,
		gamePath:     gamePath,
		localStorage: localStorage,
		portAudio:    portAudio,
		portVideo:    portVideo,
		planeID:      planeId,
		image:        image,
		devicePath:   devicePath,
		client:       cli,
	}
	return docker
}

func (d *docker) configureDocker() *container.CreateResponse {
	ctx := context.Background()

	cfg := &container.Config{
		Image: "arch:simple",
		User:  "user",
		Env: []string{
			fmt.Sprintf("EXECUTIONFILE=%s", filepath.Base(d.gamePath)),
			fmt.Sprintf("DISPLAY=:%s", d.display),
			"HOST_IP=127.0.0.1",
			fmt.Sprintf("AUDIO_PORT=%s", d.portAudio),
			fmt.Sprintf("VIDEO_PORT=%s", d.portVideo),
			fmt.Sprintf("PLANE_ID=%s", d.planeID),
			//"DISPLAY=:0",           //!!for test
			"THREAD_QUEUE_SIZE=64", //  128
			"FORMAT=nv12",
			"PRESET=ultrafast", // ultrafast, superfast, veryfast, faster, fast, medium, slow, slower, veryslow  (default: main)
			"PROFILE=100",
			"BUF_SIZE=0",
			"QP=25", // <30
			"VAAPI_DEVICE=/dev/dri/renderD128",
			"FPS=60",
			"CODEC=h264_vaapi",
			"BITRAIT=10000k",
			"G=120",
			"RESOLUTION=w=1920:h=1080",
			"PKT_SIZE=1200", //1200 max size for RTP and WebRTC
		},
		Tty: true,
	}

	cfgHost := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:     mount.TypeBind,
				Source:   os.Getenv("XAUTHORITY"),
				Target:   "/root/.Xauthority",
				ReadOnly: true,
			},
			{
				Type:   mount.TypeBind,
				Source: "/tmp/.X11-unix",
				Target: "/tmp/.X11-unix",
			},
			{
				Type:   mount.TypeBind,
				Source: createUserHome(d.username, d.localStorage),
				Target: "/home/user",
			},
			{
				Type:   mount.TypeBind,
				Source: filepath.Dir(d.gamePath),
				Target: "/home/user/game",
			},
		},
		Resources: container.Resources{
			Devices: []container.DeviceMapping{
				{
					PathOnHost:        d.devicePath,
					PathInContainer:   "/dev/input/event99",
					CgroupPermissions: "rwm",
				},
				{
					PathOnHost:        "/dev/dri/card0",
					PathInContainer:   "/dev/dri/card0",
					CgroupPermissions: "rwm",
				},
				{
					PathOnHost:        "/dev/dri/renderD128",
					PathInContainer:   "/dev/dri/renderD128",
					CgroupPermissions: "rwm",
				},
			},
		},
		Privileged:  true,
		AutoRemove:  true,
		NetworkMode: "host",
	}
	resp, err := d.client.ContainerCreate(ctx, cfg, cfgHost, nil, nil, "")
	if err != nil {
		panic(err)
	}
	d.ID = resp.ID
	return &resp
}
func (d *docker) startDocker() {
	ctx := context.Background()
	err := d.client.ContainerStart(ctx, d.ID, container.StartOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Контейнер успешно запущен:", d.ID)
}

func (d *docker) stopDocker() error {
	ctx := context.Background()
	if err := d.client.ContainerStop(ctx, d.ID, container.StopOptions{}); err != nil {
		return fmt.Errorf("failed to stop the container: %v", err)
	}
	fmt.Printf("Контейнер успешно остановлен: %s\n", d.ID)
	return nil
}
