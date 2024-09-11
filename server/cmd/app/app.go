package app

import (
	"context"
	"fmt"
	"time"

	"github.com/FurmanovVitaliy/pixel-cloud/config"
	"github.com/FurmanovVitaliy/pixel-cloud/internal/adapters/igdb"
	"github.com/FurmanovVitaliy/pixel-cloud/internal/adapters/jwt"
	"github.com/FurmanovVitaliy/pixel-cloud/internal/adapters/scanner"
	"github.com/FurmanovVitaliy/pixel-cloud/internal/domain/game"
	"github.com/FurmanovVitaliy/pixel-cloud/internal/domain/user"
	v1 "github.com/FurmanovVitaliy/pixel-cloud/internal/handler/http/api/v1"
	"github.com/FurmanovVitaliy/pixel-cloud/internal/usecase"
	"github.com/FurmanovVitaliy/pixel-cloud/pkg/client"
	"github.com/FurmanovVitaliy/pixel-cloud/pkg/http"
	"github.com/FurmanovVitaliy/pixel-cloud/pkg/logger"
)

type App struct {
	logger *logger.Logger
	config *config.Config
}

func NewApp(config *config.Config, logger *logger.Logger) (*App, error) {
	return &App{
		logger: logger,
		config: config,
	}, nil
}

func (a *App) Run() {

	// Connect to DB
	ctx := context.Background()
	mongoConn, err := client.NewClient(ctx, a.config.MongoDb.Host, a.config.MongoDb.Port, a.config.MongoDb.Username, a.config.MongoDb.Password, a.config.MongoDb.Database, a.config.MongoDb.AuthDB)
	if err != nil {
		a.logger.Fatalf("Failed to connect to DB due to error: %s", err)
	}
	a.logger.Info("Connected to DB")

	//init repositories
	gStorage := game.NewStorage(mongoConn, "games", a.logger)
	uStorage := user.NewStorage(mongoConn, "users", a.logger)
	sStorage := scanner.NewStorage(mongoConn, "hashsum", a.logger)
	//init services
	gameService := game.NewGameService(gStorage)
	userService := user.NewUserService(uStorage, time.Second*5)
	tokenService := jwt.NewJwtService(a.config.JWT.SecretKey)
	//init game-scanner
	scannerParams := scanner.CreateParams(a.config.GameSearch.SystemDirectories, a.config.GameSearch.Directories, a.config.GameSearch.NamesToCompare, a.config.GameSearch.FileExtenstions)
	scannerService := scanner.New(ctx, a.logger, scannerParams, sStorage)
	igdbClient, err := igdb.NewClient(a.config.IGBD.ID, a.config.IGBD.Token)
	if err != nil {
		a.logger.Fatalf("failed to connect to IGDB due to error: %s", err)
	}
	igdbService := igdb.New(igdbClient, a.logger)
	//init usecase
	uc := usecase.NewUseCase(
		userService,
		gameService,
		tokenService,
		igdbService,
		scannerService,
		a.logger)

	uc.ScanLibrary()

	//init http server
	var serverConfig = &http.ServerConfig{
		ServerPort:             a.config.Server.Port,
		CorsAlloedMethods:      a.config.Cors.AllowedMethods,
		CorsAllowedHeaders:     a.config.Cors.AllowedHeaders,
		CorsAllowedOrigins:     a.config.Cors.AllowedOrigins,
		CorsExposedHeaders:     a.config.Cors.ExposedHeaders,
		CorsMaxAge:             a.config.Cors.MaxAge,
		IsDebug:                a.config.IsDebug,
		CorsAllowedCredentials: a.config.Cors.AllowCredentials,
		CertFilePath:           fmt.Sprintf("%s/%s", a.config.Workdir, a.config.Certificates.Cert),
		KeyFilePath:            fmt.Sprintf("%s/%s", a.config.Workdir, a.config.Certificates.Key),
	}

	handler := v1.NewHandler(uc, a.logger)
	router := http.NewHttpRouter(handler)
	server := http.NewServer(a.logger, serverConfig, router)
	server.Run()
}
