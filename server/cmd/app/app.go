package app

import (
	"context"
	"fmt"
	"time"

	"github.com/FurmanovVitaliy/pixel-cloud/config"
	"github.com/FurmanovVitaliy/pixel-cloud/internal/adapters/igdb"
	"github.com/FurmanovVitaliy/pixel-cloud/internal/adapters/jwt"
	"github.com/FurmanovVitaliy/pixel-cloud/internal/core/hub"
	"github.com/FurmanovVitaliy/pixel-cloud/internal/core/scanner"
	"github.com/FurmanovVitaliy/pixel-cloud/internal/core/srm"
	"github.com/FurmanovVitaliy/pixel-cloud/internal/domain/game"
	"github.com/FurmanovVitaliy/pixel-cloud/internal/domain/user"
	v1 "github.com/FurmanovVitaliy/pixel-cloud/internal/handler/http/api/v1"
	"github.com/FurmanovVitaliy/pixel-cloud/internal/infrastructure/display"
	"github.com/FurmanovVitaliy/pixel-cloud/internal/infrastructure/docker"
	"github.com/FurmanovVitaliy/pixel-cloud/internal/infrastructure/wrtc"
	"github.com/FurmanovVitaliy/pixel-cloud/internal/usecase"
	"github.com/FurmanovVitaliy/pixel-cloud/pkg/client"
	"github.com/FurmanovVitaliy/pixel-cloud/pkg/http"
	"github.com/FurmanovVitaliy/pixel-cloud/pkg/logger"
)

type App struct {
	logger *logger.Logger
	c      *config.Config
}

func NewApp(config *config.Config, logger *logger.Logger) (*App, error) {
	return &App{
		logger: logger,
		c:      config,
	}, nil
}

func (a *App) Run() {

	// Connect to DB
	ctx := context.Background()
	mongoConn, err := client.NewMongoClient(ctx, a.c.MongoDb.Host, a.c.MongoDb.Port, a.c.MongoDb.Username, a.c.MongoDb.Password, a.c.MongoDb.Database, a.c.MongoDb.AuthDB)
	if err != nil {
		a.logger.Fatalf("failed to connect to DB due to error: %s", err)
	}
	a.logger.Info("connected to MongoDB")
	/*
		postgresConn, err := client.NewPostgresClient(ctx, 5, a.config.Postgres.Host, a.config.Postgres.Port, a.config.Postgres.Username, a.config.Postgres.Password, a.config.Postgres.Database)
		if err != nil {
			a.logger.Fatalf("failed to connect to DB due to error: %s", err)
		}
	*/
	igdbClient, err := igdb.NewClient(a.c.IGBD.ID, a.c.IGBD.Token)
	if err != nil {
		a.logger.Fatalf("failed to connect to IGDB due to error: %s", err)
	}

	//init repositories
	gStorage := game.NewStorage(mongoConn, "games", a.logger)
	uStorage := user.NewStorage(mongoConn, "users", a.logger)
	sStorage := scanner.NewStorage(mongoConn, "hashsum", a.logger)
	dStorage := display.NewArrayXServerRepository()
	//init services
	gameService := game.NewGameService(gStorage)
	userService := user.NewUserService(uStorage, time.Second*5)
	igdbService := igdb.New(igdbClient, a.logger)
	tokenService := jwt.NewJwtService(a.c.JWT.SecretKey)
	stremService := wrtc.NewStreamerService(a.c.Streamer.VideoCodec, a.c.Streamer.AudioCodec)
	vmService := docker.NewGameVmService(a.c.Docker.PulseImage, a.c.Docker.VideoImage, a.c.Docker.AudioImage, a.c.Docker.ProtoneImage, a.c.Docker.NetworkMode, a.c.Docker.RendererPath, nil, a.c.VideoCapture.Env, a.c.AudioCapture.Env, nil)

	//spawn core modules-services (hub, resorses-manager, game-scanner etc)
	//these services will integrate and used in usecase
	scannerParams := scanner.CreateParams(a.c.GameSearch.SystemDirectories, a.c.GameSearch.Directories, a.c.GameSearch.NamesToCompare, a.c.GameSearch.FileExtenstions)
	scannerService := scanner.New(ctx, a.logger, scannerParams, sStorage)

	displayService := display.NewXServerService(dStorage)
	displayService.PopulateViaLocalScript(a.c.VirtualDisplayInitializer.EnableVirtualDisplaysScriptPath, a.c.VirtualDisplayInitializer.DisplayInfoJsonPath)
	srmDisplays := displayService.GetAll()
	srmParams := srm.CreateParams(a.c.UsersFileStorage.FsInitFilesPath, a.c.UsersFileStorage.Path, a.c.UDPReader.MinPort, a.c.UDPReader.MaxPort, a.c.UDPReader.UdpBuffer, a.c.UDPReader.ReadBuffer)
	srmService := srm.New(srmParams, srmDisplays, a.logger)

	hubService := hub.New()
	go hubService.Run()

	//init usecase
	uc := usecase.NewUseCase(
		userService,
		gameService,
		tokenService,
		igdbService,
		scannerService,
		hubService,
		stremService,
		srmService,
		vmService,
		a.logger)

	uc.ScanLibrary()

	//init http server
	var serverConfig = &http.ServerConfig{
		ServerPort:             a.c.Server.Port,
		CorsAlloedMethods:      a.c.Cors.AllowedMethods,
		CorsAllowedHeaders:     a.c.Cors.AllowedHeaders,
		CorsAllowedOrigins:     a.c.Cors.AllowedOrigins,
		CorsExposedHeaders:     a.c.Cors.ExposedHeaders,
		CorsMaxAge:             a.c.Cors.MaxAge,
		IsDebug:                a.c.IsDebug,
		CorsAllowedCredentials: a.c.Cors.AllowCredentials,
		CertFilePath:           fmt.Sprintf("%s/%s", a.c.Workdir, a.c.Certificates.Cert),
		KeyFilePath:            fmt.Sprintf("%s/%s", a.c.Workdir, a.c.Certificates.Key),
	}

	handler := v1.NewHandler(uc, a.logger)
	router := http.NewHttpRouter(handler)
	server := http.NewServer(a.logger, serverConfig, router)

	server.Run()
}
