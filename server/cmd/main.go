package main

import (
	"log"

	"github.com/FurmanovVitaliy/pixel-cloud/cmd/app"
	"github.com/FurmanovVitaliy/pixel-cloud/config"
	"github.com/FurmanovVitaliy/pixel-cloud/pkg/logger"
)

func main() {
	log.Print("config initializing")
	config := config.GetConfig()
	log.Print("logger initializing")
	logger := logger.Init(config.LogLevel)
	//Controller
	app, err := app.NewApp(config, &logger)
	if err != nil {
		logger.Fatalf("app initializing error: %s", err)
	}
	logger.Info("app starting")
	app.Run()

}
