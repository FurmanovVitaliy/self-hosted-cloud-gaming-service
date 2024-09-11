package app

import (
	"context"
	"time"

	"github.com/FurmanovVitaliy/pixel-cloud/config"
	"github.com/FurmanovVitaliy/pixel-cloud/internal/domain/game"
	"github.com/FurmanovVitaliy/pixel-cloud/internal/domain/user"
	"github.com/FurmanovVitaliy/pixel-cloud/internal/usecase"
	"github.com/FurmanovVitaliy/pixel-cloud/pkg/client"
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

	//init services
	gameService := game.NewGameService(gStorage)
	userService := user.NewUserService(uStorage, time.Second*5)

	//init usecase
	uc := usecase.NewUseCase(userService, gameService, a.logger)

	println(uc)
}
