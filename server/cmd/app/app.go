package app

import (
	"context"

	"github.com/FurmanovVitaliy/pixel-cloud/config"
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

	println(mongoConn)
}
