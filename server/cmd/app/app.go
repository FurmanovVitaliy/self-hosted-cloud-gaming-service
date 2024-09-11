package app

import (
	"github.com/FurmanovVitaliy/pixel-cloud/config"
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
	a.logger.Info("app is running")
}
