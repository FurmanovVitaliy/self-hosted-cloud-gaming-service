package usecase

import (
	"context"
	"time"

	"github.com/FurmanovVitaliy/pixel-cloud/internal/adapters/igdb"
	"github.com/FurmanovVitaliy/pixel-cloud/internal/core/hub"
	"github.com/FurmanovVitaliy/pixel-cloud/internal/core/scanner"
	"github.com/FurmanovVitaliy/pixel-cloud/internal/core/srm"
	"github.com/FurmanovVitaliy/pixel-cloud/internal/domain/game"
	"github.com/FurmanovVitaliy/pixel-cloud/internal/domain/user"
	"github.com/FurmanovVitaliy/pixel-cloud/internal/infrastructure/docker"
	"github.com/FurmanovVitaliy/pixel-cloud/internal/infrastructure/wrtc"
	"github.com/FurmanovVitaliy/pixel-cloud/pkg/logger"
)

// services
type userService interface {
	CreateUser(ctx context.Context, user user.User) (string, error)
	FindByEmail(ctx context.Context, email string) (user.User, error)
	GetList(ctx context.Context) ([]user.User, error)
}

type tokenService interface {
	CreateToken(userID, username string, expiry time.Duration) (string, error)
}

type gameService interface {
	GetAll(ctx context.Context) ([]game.Game, error)
	Create(ctx context.Context, game game.Game) (string, error)
	GetOneById(ctx context.Context, id string) (game.Game, error)
	DeleteAll(ctx context.Context) error
}

type scannerSerrvice interface {
	CheckForChanges() bool
	Scan() ([]scanner.Game, error)
}

type gameInfoService interface {
	GetExtraInfoByName(name string) (gameInfo igdb.GameExtraInfo, err error)
}

type hubService interface {
	GetRooms() (rooms []hub.Room, err error)
	GetRoom(uuid string) (*hub.Room, error)
	CreateRoom(uuid, gameid string) (err error)
	CreateWorker(roomUUID, username string, msgHandler hub.MessageHandler, streamer hub.WrtcStreamer, audioR, videoR hub.UDPReader, vm hub.VM) error
}
type streamerService interface {
	CreateRTC(onMessage func(data []byte)) (wrtc.WebRTC, error)
}

type resourceService interface {
	AllocateDisplay() (srm.Display, error)
	AllocateDiscSpace(username string) (string, error)
	AllocateListeners(ip string) (readers [2]srm.UDPReader, extraPort int, err error)
	AllocateInputDevice(username, vendorID, productID, productT string) (device srm.InputDevice, err error)
}

type vmService interface {
	ConfigureAndCreate(ctx context.Context,
		username, userHomePath, inputDevicePath, gamePath, display, hostIP string,
		planeID, videoPort, audioPort, pulseServerPort int) (docker.VM, error)
}

type UseCase struct {
	userService     userService
	gameService     gameService
	tokenService    tokenService
	gameInfoService gameInfoService
	scannerSerrvice scannerSerrvice
	hubService      hubService
	streamerService streamerService
	resourceService resourceService
	vmService       vmService

	logger *logger.Logger
}

func NewUseCase(
	userService userService,
	gameService gameService,
	tokenService tokenService,
	gameInfoService gameInfoService,
	scannerSerrvice scannerSerrvice,
	hubService hubService,
	streamerService streamerService,
	resourceManagerService resourceService,
	dockerService vmService,

	logger *logger.Logger,
) *UseCase {
	return &UseCase{
		userService:     userService,
		gameService:     gameService,
		tokenService:    tokenService,
		gameInfoService: gameInfoService,
		scannerSerrvice: scannerSerrvice,
		hubService:      hubService,
		streamerService: streamerService,
		resourceService: resourceManagerService,
		vmService:       dockerService,

		logger: logger,
	}
}
