package hhtpx

import (
	"cloud/pkg/logger"
	"net/http"
	"time"
)

type Server struct {
	http.Server

	opts Options
	//listener net.Listener
	//redirect *Server
	log *logger.Logger
}

type (
	Mux struct {
		*http.ServeMux
		//prefix string
	}
	Handler        = http.Handler
	HandlerFunc    = http.HandlerFunc
	ResponseWriter = http.ResponseWriter
	Request        = http.Request
)

func NewHTTPServer(adres string, handler func(*Server) Handler, options ...Option) (*Server, error) {
	opts := &Options{
		Https:         false,
		HttpsRedirect: true,
		IdleTimeout:   120 * time.Second,
		ReadTimeout:   500 * time.Second,
		WriteTimeout:  500 * time.Second,
	}
	opts.override(options...)

	server := &Server{
		Server: http.Server{
			Addr:         adres,
			IdleTimeout:  opts.IdleTimeout,
			ReadTimeout:  opts.ReadTimeout,
			WriteTimeout: opts.WriteTimeout,
		},
		opts: *opts,
		log:  opts.Logger,
	}

	server.Handler = handler(server)
	return server, nil
}

func FileServer(dir string) http.Handler { return http.FileServer(http.Dir(dir)) }
