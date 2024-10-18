package docker

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

//! Thre was a few ty to build it another way
//! Making of  one sing container provides to freeze in webrtc srteam for some reason
//! Running conteiners via docker-compose provides to conflits with access to the same or directories
//! So, I decided to run each container separately with the same env variables and one network

type ContainerGroupConfig struct {
	User string
	//image for each container
	PulseImage   string
	VideoImage   string
	AudioImage   string
	ProtoneImage string
	//general env for all containers
	HomePath       string
	GameDirPath    string
	DevicePath     string
	NetworkMode    string
	CardPath       string
	RendererPath   string
	XauthorityPath string
	//env for each container
	PulseEnv   []string
	VideoEnv   []string
	AudioEnv   []string
	ProtoneEnv []string
}
type docker struct {
	name     string
	isStoped bool
}

type dockerGroup struct {
	client  *client.Client
	ctx     context.Context
	mu      sync.Mutex
	dockers map[string]*docker
}

func newDockerGroup(ctx context.Context, cfg *ContainerGroupConfig) (dg *dockerGroup, err error) {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return nil, err
	}
	pulseContainerConfigs := &container.Config{
		Image: cfg.PulseImage,
		Env:   cfg.PulseEnv,
		Tty:   true,
	}

	videoContainerConfigs := &container.Config{
		Image: cfg.VideoImage,
		Env:   cfg.VideoEnv,
		Tty:   true,
	}
	audioContainerConfigs := &container.Config{
		Image: cfg.AudioImage,
		Env:   cfg.AudioEnv,
		Cmd:   []string{"sleep 10"},
		Tty:   true,
	}
	protoneContainerConfigs := &container.Config{
		Image: cfg.ProtoneImage,
		Env:   cfg.ProtoneEnv,
		Cmd:   []string{"sleep 10"},
		Tty:   true,
	}

	pulseHostConfigs := &container.HostConfig{
		AutoRemove:  true,
		NetworkMode: container.NetworkMode(cfg.NetworkMode),
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: cfg.HomePath,
				Target: fmt.Sprintf("/home/%s", cfg.User),
			},
		},
	}
	videoHostConfigs := &container.HostConfig{
		Privileged:  true,
		AutoRemove:  true,
		NetworkMode: container.NetworkMode(cfg.NetworkMode),
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: cfg.HomePath,
				Target: fmt.Sprintf("/home/%s", cfg.User),
			},
		},
		Resources: container.Resources{
			Devices: []container.DeviceMapping{
				{
					PathOnHost:        cfg.RendererPath,
					PathInContainer:   "/dev/dri/renderD128",
					CgroupPermissions: "rwm",
				},
			},
		},
	}
	audioHostConfigs := &container.HostConfig{
		AutoRemove:  true,
		NetworkMode: container.NetworkMode(cfg.NetworkMode),
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: cfg.HomePath,
				Target: fmt.Sprintf("/home/%s", cfg.User),
			},
		},
	}
	protoneHostConfigs := &container.HostConfig{
		Privileged:  true,
		AutoRemove:  true,
		NetworkMode: container.NetworkMode(cfg.NetworkMode),
		Mounts: []mount.Mount{
			{
				Type:     mount.TypeBind,
				Source:   cfg.XauthorityPath,
				Target:   fmt.Sprintf("/home/%s/.Xauthority", cfg.User),
				ReadOnly: true,
			},
			{
				Type:   mount.TypeBind,
				Source: "/tmp/.X11-unix",
				Target: "/tmp/.X11-unix",
			},
			{
				Type:   mount.TypeBind,
				Source: cfg.HomePath,
				Target: fmt.Sprintf("/home/%s", cfg.User),
			},
			{
				Type:   mount.TypeBind,
				Source: cfg.GameDirPath,
				Target: fmt.Sprintf("/home/%s/game", cfg.User),
			},
		},
		Resources: container.Resources{
			Devices: []container.DeviceMapping{
				{
					PathOnHost:        cfg.DevicePath,
					PathInContainer:   "/dev/input/event99",
					CgroupPermissions: "rwm",
				},
			},
		},
	}

	dg = &dockerGroup{
		client:  cli,
		ctx:     ctx,
		dockers: make(map[string]*docker),
	}

	resp, err := cli.ContainerCreate(ctx, pulseContainerConfigs, pulseHostConfigs, nil, nil, "")
	if err != nil {
		return nil, err
	}
	dg.dockers[resp.ID] = &docker{name: cfg.User + "-pulse", isStoped: false}

	resp, err = cli.ContainerCreate(ctx, videoContainerConfigs, videoHostConfigs, nil, nil, "")
	if err != nil {
		return nil, err
	}
	dg.dockers[resp.ID] = &docker{name: cfg.User + "-video", isStoped: false}

	resp, err = cli.ContainerCreate(ctx, audioContainerConfigs, audioHostConfigs, nil, nil, "")
	if err != nil {
		return nil, err
	}
	dg.dockers[resp.ID] = &docker{name: cfg.User + "-audio", isStoped: false}

	resp, err = cli.ContainerCreate(ctx, protoneContainerConfigs, protoneHostConfigs, nil, nil, "")
	if err != nil {
		return nil, err
	}
	dg.dockers[resp.ID] = &docker{name: cfg.User + "-protone", isStoped: false}

	return dg, nil
}

// Start starts all containers in the group. If any container stops or fails to start, it stops all other containers.
// This function is blocking and will return when all containers have stopped.
// Go func if best way to run it in parallel
func (dg *dockerGroup) Start() error {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(dg.ctx)
	defer cancel()

	var startErrs []string
	var cenceled []string

	for id := range dg.dockers {
		wg.Add(1)
		go func(id, name string) {
			defer wg.Done()
			log.Printf("Starting container %s\n", name)
			err := dg.client.ContainerStart(ctx, id, container.StartOptions{})
			if err != nil {
				log.Printf("Error starting container %s: %v\n", name, err)
				startErrs = append(startErrs, fmt.Sprintf("Error starting container %s: %v", name, err))
				cancel() // Cancel the context to stop other containers
				return
			}

			statusCh, waitErrCh := dg.client.ContainerWait(ctx, id, container.WaitConditionNotRunning)
			select {
			case err := <-waitErrCh:
				if err != nil {
					cenceled = append(cenceled, fmt.Sprintf("Error waiting for container %s to stop: %v", name, err))
					log.Printf("Error waiting for container %s to stop: %v\n", name, err)
					cancel() // Cancel the context to stop other containers
					return
				}
			case status := <-statusCh:
				log.Printf("Container %s stopped with status code: %d\n", name, status.StatusCode)
				cancel() // Cancel the context to stop other containers
				return

			}

			dg.mu.Lock()
			dg.dockers[id].isStoped = true
			dg.mu.Unlock()

		}(id, dg.dockers[id].name)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	if len(cenceled) > 0 {
		dg.Stop()
		return nil
	}

	if len(startErrs) > 0 {
		dg.Stop()
		return fmt.Errorf("failed to start containers: %v", startErrs)
	}
	return nil
}

// Stop stops all running containers in the group.
func (dg *dockerGroup) Stop() {
	dg.mu.Lock()
	defer dg.mu.Unlock()
	defer dg.client.Close()

	for id := range dg.dockers {
		// Skip if the container has already been stopped
		if dg.dockers[id].isStoped {
			continue
		}

		// Check if the container exists before trying to stop or kill
		_, err := dg.client.ContainerInspect(dg.ctx, id)
		if err != nil {
			if client.IsErrNotFound(err) {
				continue
			}
			log.Printf("Error inspecting container %s: %v\n", dg.dockers[id].name, err)
			continue
		}

		fmt.Printf("Stopping container %s\n", dg.dockers[id].name)
		err = dg.client.ContainerStop(dg.ctx, id, container.StopOptions{Timeout: nil})
		if err != nil {
			fmt.Printf("Error stopping container %s, trying to kill: %v\n", dg.dockers[id].name, err)
			killErr := dg.client.ContainerKill(dg.ctx, id, "SIGKILL")
			if killErr != nil {
				fmt.Printf("Error killing container %s: %v\n", dg.dockers[id].name, killErr)
			}
		} else {
			dg.dockers[id].isStoped = true // Mark as stopped
		}
	}
}
