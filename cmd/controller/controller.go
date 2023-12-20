package controller

import (
	"cloud/config"
	"time"

	"cloud/internal/api/handlers"
	"cloud/internal/domain/games"
	"cloud/internal/domain/games/library"
	"cloud/internal/domain/hub"
	"cloud/internal/domain/user"
	"cloud/pkg/client"
	"cloud/pkg/logger"
	hashsum "cloud/pkg/utils/heshsum"
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

type Controller struct {
	logger   *logger.Logger
	config   *config.Config
	hub      *hub.Hub
	storage  [2]interface{}
	handlers [3]handlers.Handler
	doneChan chan struct{}
}

func NewController(config *config.Config, logger *logger.Logger) (*Controller, error) {
	return &Controller{
		doneChan: make(chan struct{}),
		logger:   logger,
		config:   config,
	}, nil
}

func (c *Controller) Run() {
	c.servicesInition()
	go c.scanLib()
	go c.hub.Run()
	go c.startHttpServer()
	<-c.doneChan
}

func (c *Controller) Stop() {
	close(c.doneChan)
}

func (c *Controller) servicesInition() {
	cfg := c.config.MongoDb
	dbConn, err := client.NewClient(context.Background(), cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database, cfg.AuthDB)
	if err != nil {
		c.logger.Fatalf("Failed to connect to DB due to error: %s", err)
	}
	c.logger.Info("Connected to DB")

	c.storage[0] = hashsum.HashStorage(dbConn, "hashsum", c.logger)

	uStorage := user.NewStorage(dbConn, "users", c.logger)

	uSetvice := user.NewService(uStorage, c.logger)

	gStorage := games.NewStorage(dbConn, "games", c.logger)
	c.storage[1] = gStorage
	gService := games.NewService(gStorage)

	c.hub = hub.NewHub()

	c.handlers[0] = user.Handler(uSetvice, c.logger)
	c.handlers[1] = games.Handler(gService)
	c.handlers[2] = hub.Handler(c.hub, c.logger, gService)
}

func (c *Controller) startHttpServer() {
	certFile := c.config.Workdir + c.config.Certificate.Cert
	keyFile := c.config.Workdir + c.config.Certificate.Key

	cors := cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST"},
		AllowedOrigins:   []string{"https://localhost:3000", "https://192.168.1.13:3000"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		Debug:            true,
	})

	router := httprouter.New()

	for _, handler := range c.handlers {
		handler.Register(router)
	}

	httpServer := http.Server{
		Handler:      cors.Handler(router),
		Addr:         ":8000", //TODO: move to config
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  500 * time.Second,
		WriteTimeout: 500 * time.Second,
	}

	c.logger.Infof("HTTP server starting on port %s", httpServer.Addr)
	httpServer.ListenAndServeTLS(certFile, keyFile)

}
func (c *Controller) scanLib() {
	hstore := c.storage[0].(hashsum.Storage)
	gstore := c.storage[1].(games.Storage)

	serchCfg := c.config.GameSearch
	igbdCfg := c.config.IGBD

	scaner := library.NewSerchEngine(serchCfg.FileExtenstions, serchCfg.Directories, serchCfg.NamesToCompare, igbdCfg.ID, igbdCfg.Token, c.logger)
	g, hash, err := scaner.ScanLibrary()
	if err != nil {
		c.logger.Error(err)
		return
	}
	dbHash, _ := hstore.FindOne(context.Background(), "game_lib")
	if hashsum.CheckSum(dbHash, hash) {
		c.logger.Info("no changes in library")
		return
	}
	hstore.Update(context.Background(), "game_lib", hash)
	g, err = scaner.GetInfoFromIGDB(g)
	if err != nil {
		c.logger.Error(err)
		return
	}
	if err = gstore.FullyUpdate(context.Background(), g); err != nil {
		c.logger.Error(err)
		return
	}
}
