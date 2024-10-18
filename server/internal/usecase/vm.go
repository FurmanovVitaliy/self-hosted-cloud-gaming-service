package usecase

import (
	"context"

	"github.com/FurmanovVitaliy/pixel-cloud/internal/core/srm"
	"github.com/FurmanovVitaliy/pixel-cloud/internal/infrastructure/docker"
)

func (uc *UseCase) createVmWithSrmresources(ip, uuid, username, userHomeDir string, inputDevice srm.InputDevice, display srm.Display, vReade, aReader srm.UDPReader, pulsePort int) (vm docker.VM, err error) {

	room, err := uc.hubService.GetRoom(uuid)
	if err != nil {
		return nil, RoomNotExist
	}

	game, err := uc.gameService.GetOneById(context.Background(), room.GameID)
	if err != nil {
		uc.logger.Error("failed to get game: %w", err)
		return nil, ErrGameUnavailable
	}
	inputPath, err := inputDevice.GetEvdevPath()
	if err != nil {
		uc.logger.Error("failed to get input device path: %w", err)
		return nil, err
	}

	vm, err = uc.vmService.ConfigureAndCreate(context.Background(), username, userHomeDir, inputPath, game.Path, display.GetDisplayNumper(), ip, display.GetPlaneID(), vReade.GetPort(), aReader.GetPort(), pulsePort)
	return vm, err
}
