package http

import (
	"context"
	"net/http"
	"time"

	"github.com/FurmanovVitaliy/pixel-cloud/pkg/logger"

	"github.com/rs/cors"
)

type Router interface {
	RegisterRoutes()
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

type ServerConfig struct {
	ServerPort             string
	CorsAlloedMethods      []string
	CorsAllowedHeaders     []string
	CorsAllowedOrigins     []string
	CorsExposedHeaders     []string
	CorsMaxAge             int
	IsDebug                bool
	CorsAllowedCredentials bool
	CertFilePath           string
	KeyFilePath            string
}

type Server struct {
	logger     *logger.Logger
	httpServer *http.Server

	certFile string
	keyFile  string
}

func NewServer(logger *logger.Logger, config *ServerConfig, router Router) *Server {
	var handler http.Handler = router
	router.RegisterRoutes()

	c := cors.New(cors.Options{
		AllowedMethods:   config.CorsAlloedMethods,
		AllowedOrigins:   config.CorsAllowedOrigins,
		AllowedHeaders:   config.CorsAllowedHeaders,
		ExposedHeaders:   config.CorsExposedHeaders,
		MaxAge:           config.CorsMaxAge,
		Debug:            config.IsDebug,
		AllowCredentials: true,
	})
	//Bind Middleware
	handler = c.Handler(handler)

	// Set up the http server
	httpServer := &http.Server{
		Addr:           ":" + config.ServerPort,
		Handler:        handler,
		IdleTimeout:    120 * time.Second,
		ReadTimeout:    50 * time.Second,
		WriteTimeout:   50 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	return &Server{
		logger:     logger,
		httpServer: httpServer,
		certFile:   config.CertFilePath,
		keyFile:    config.KeyFilePath,
	}
}

func (s *Server) Run() error {
	if s.certFile == "" || s.keyFile == "" {
		s.logger.Infof("HTTP server starting on port %s", s.httpServer.Addr)
		err := s.httpServer.ListenAndServe()
		if err != nil {
			s.logger.Errorf("failed to start HTTP server: %v", err)
			return err
		}
		return nil
	}
	s.logger.Infof("HTTPs server starting on port %s", s.httpServer.Addr)
	err := s.httpServer.ListenAndServeTLS(s.certFile, s.keyFile)
	if err != nil {
		s.logger.Errorf("failed to start HTTPs server: %v", err)
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
