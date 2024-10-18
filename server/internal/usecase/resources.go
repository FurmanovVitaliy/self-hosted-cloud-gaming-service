package usecase

import (
	"net/http"

	"github.com/FurmanovVitaliy/pixel-cloud/internal/core/srm"
	"github.com/FurmanovVitaliy/pixel-cloud/pkg/errors"
)

var FailedToAllocateResources = errors.New(http.StatusInternalServerError, "RS", "00000", "failed to allocate resources")

func (uc *UseCase) releaseAllResources(
	inputDevice srm.InputDevice,
	virtualDisplay srm.Display,
	aReader, vReader srm.UDPReader) {
	uc.logger.Info("releasing resources")
	// Попытка закрыть input device
	if err := inputDevice.Close(); err != nil {
		uc.logger.Warn("failed to close input device")
	}
	// Закрытие виртуального дисплея
	virtualDisplay.Close()
	// Попытка закрыть audio reader
	if err := aReader.Close(); err != nil {
		uc.logger.Warn("failed to close audio reader")
	}
	// Попытка закрыть video reader
	if err := vReader.Close(); err != nil {
		uc.logger.Warn("failed to close video reader")
	}
}

func (uc *UseCase) allocateAllResources(
	ip, username string,
	deviceInfo JoinRoomRes) (srm.InputDevice, srm.Display, srm.UDPReader, srm.UDPReader, int, string, error) {

	inputDevice, err := uc.resourceService.AllocateInputDevice(username, deviceInfo.Control.VendorID, deviceInfo.Control.ProductID, deviceInfo.Control.Type)
	if err != nil {
		return nil, nil, nil, nil, 0, "", errors.New(http.StatusInternalServerError, "RS", "00001", "input device not available!")
	}

	virtualDisplay, err := uc.resourceService.AllocateDisplay()
	if err != nil || virtualDisplay == nil {
		inputDevice.Close()
		return nil, nil, nil, nil, 0, "", errors.New(http.StatusInternalServerError, "RS", "00002", "xserver not available")
	}

	listeners, emptyPort, err := uc.resourceService.AllocateListeners(ip)
	if err != nil {
		inputDevice.Close()
		virtualDisplay.Close()
		return nil, nil, nil, nil, 0, "", errors.New(http.StatusInternalServerError, "RS", "00003", "listeners not available")
	}

	userDiskSpace, err := uc.resourceService.AllocateDiscSpace(username)
	if err != nil {
		inputDevice.Close()
		virtualDisplay.Close()
		if err = listeners[0].Close(); err != nil {
			uc.logger.Warn("failed to close listener 0")
		}
		if err = listeners[1].Close(); err != nil {
			uc.logger.Warn("failed to close listener 1")
		}
		return nil, nil, nil, nil, 0, "", errors.New(http.StatusInternalServerError, "RS", "00004", "disk space not available")
	}

	return inputDevice, virtualDisplay, listeners[0], listeners[1], emptyPort, userDiskSpace, nil
}
