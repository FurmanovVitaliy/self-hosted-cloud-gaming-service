package hhtpx

import (
	"cloud/pkg/logger"
	"time"
)

type (
	Options struct {
		Https              bool
		HttpsRedirect      bool
		HttpsRedirectAdres string
		HttpsCert          string
		HttpsKey           string
		HttpsDomain        string
		PortRoll           bool
		IdleTimeout        time.Duration
		ReadTimeout        time.Duration
		WriteTimeout       time.Duration
		Logger             *logger.Logger
		Zone               string
	}
	Option func(*Options)
)

func (o *Options) override(options ...Option) {
	for _, opt := range options {
		opt(o)
	}
}
