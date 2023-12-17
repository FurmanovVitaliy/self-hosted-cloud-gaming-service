package main

import (
	"cloud/cmd/controller"
	"cloud/config"
	"cloud/pkg/logger"

	"github.com/rs/zerolog/log"
)

func main() {
	log.Print("Config initializing")
	config := config.GetConfig()
	log.Print("Logger initializing")
	logger := logger.Init(config.LogLevel)
	//Controller
	controller, err := controller.NewController(config, &logger)
	if err != nil {
		logger.Fatalf("Controller initializing error: %s", err) //TODO: Handle error
	}
	logger.Info("Controller starting")
	controller.Run()
}
