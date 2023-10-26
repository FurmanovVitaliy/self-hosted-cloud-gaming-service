package controller

import (
	"cloud/config"
	gamehandler "cloud/internal/adapters/api/handlers/games"
	hubhandler "cloud/internal/adapters/api/handlers/hub"
	userhandler "cloud/internal/adapters/api/handlers/user"
	"cloud/internal/domain/games"
	"cloud/internal/domain/hub"
	"cloud/internal/domain/user"
	"cloud/pkg/client"
	"cloud/pkg/logger"
	"context"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

type Controller struct {
	logger *logger.Logger
	config *config.Config
}

func NewController(config *config.Config, logger *logger.Logger) (*Controller, error) {

	http.Handle("/", fileServer("/home/vitalii/dev/go-code/fileSearch/web/build"))

	return &Controller{
		logger: logger,
		config: config,
	}, nil
}
func (c *Controller) Run() {
	go c.startHttpServer()
	go c.connectToDB()
	//c.startServices()
	//c.scanLib()
	//c.starHandlingUsers()
}

func (c *Controller) Stop() {
	//c.stopServices()
}

func (c *Controller) connectToDB() {

	cfg := c.config.MongoDb
	dbConn, err := client.NewClient(context.Background(), cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database, cfg.AuthDB)
	if err != nil {
		c.logger.Fatalf("Failed to connect to DB due to error: %s", err)
	}
	c.logger.Info("Connected to DB")
	userStorage := user.NewStorage(dbConn, "users", c.logger)
	userService := user.NewService(userStorage, c.logger)
	userHandler := userhandler.NewHandler(userService, c.logger)

	gamesStorage := games.NewStorage(dbConn, "games", c.logger)
	gamesService := games.NewService(gamesStorage)
	gamesHandler := gamehandler.NewHandler(gamesService)

	//!!hub must be in controller
	hub := hub.NewHub()
	hubHandler := hubhandler.NewHandler(hub)
	cors := cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST"},
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		Debug:            true,
	})

	router := httprouter.New()
	userHandler.Register(router)
	hubHandler.Register(router)
	gamesHandler.Register(router)
	go hub.Run()
	http.ListenAndServe(":8000", cors.Handler(router))

}

func (c *Controller) startHttpServer() {
	cors := cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST"},
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Debug:            true,
	})

	handler := cors.Handler(http.DefaultServeMux)

	httpServer := http.Server{
		Handler:      handler,
		Addr:         ":8080", //TODO: from config
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  500 * time.Second,
		WriteTimeout: 500 * time.Second,
	}

	c.logger.Info("HTTP server starting on port 8080")
	httpServer.ListenAndServe()

}

func fileServer(dir string) http.Handler { return http.FileServer(http.Dir(dir)) }
