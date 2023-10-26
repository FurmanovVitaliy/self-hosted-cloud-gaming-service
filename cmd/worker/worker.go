package worker

import (
	"cloud/config"
	"cloud/internal/adapters/api/handlers/metric"
	"context"

	"fmt"

	"cloud/pkg/logger"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

type Worker struct {
	config     *config.Config
	logger     *logger.Logger
	router     *httprouter.Router
	httpServer *http.Server
}

func NewWorker(config *config.Config, logger *logger.Logger) (Worker, error) {
	logger.Info("Router initializing")
	router := httprouter.New()
	//TODO: logger.Info("Swagger docs initializing")

	logger.Info("Routes registration on server:")

	logger.Info("---Metric rout registration")
	metricHandler := metric.Handler{}
	metricHandler.Register(router)

	return Worker{
		config: config,
		logger: logger,
		router: router,
	}, nil
}
func (w *Worker) Run() {
	w.startHtttp()
}

func (w *Worker) startHtttp() {

	cors := cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedOrigins:   []string{"http://localhost:9090"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
		Debug:            true,
	})

	handler := cors.Handler(w.router)

	w.httpServer = &http.Server{
		Handler:      handler,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	var listener net.Listener
	var listenErr error

	if w.config.Listen.Type == "socket" {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		w.logger.Infof("ditect app path :%s", appDir)
		if err != nil {
			w.logger.Fatal(err)
		}
		w.logger.Info("create socket")
		socketPath := path.Join(appDir, "app.socket")
		w.logger.Debugf("socket path:%s", socketPath)

		w.logger.Info("listen unix socket")
		listener, listenErr = net.Listen("unix", socketPath)
		w.logger.Infof("server is listening unix socket:%s", socketPath)

	} else {
		w.logger.Info("listen tcp")
		listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", w.config.Listen.BindIP, w.config.Listen.Port))
		w.logger.Infof("server is listening port %s:%s", w.config.Listen.BindIP, w.config.Listen.Port)
	}
	if listenErr != nil {
		w.logger.Fatal(listenErr)

	}
	w.logger.Info("App initialazing complet")
	if err := w.httpServer.Serve(listener); err != nil {
		switch {
		//TODO: case errors.Is(err, http.ErrServerClosed):
		//TODO	w.logger.Info("server closed")
		default:
			w.logger.Fatal(err)
		}
	}
	err := w.httpServer.Shutdown(context.Background())
	if err != nil {
		w.logger.Fatal(err)
	}

}
