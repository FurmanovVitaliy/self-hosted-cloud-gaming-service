package v1

import (
	"github.com/FurmanovVitaliy/pixel-cloud/internal/usecase"
	"github.com/FurmanovVitaliy/pixel-cloud/pkg/logger"
)

type handler struct {
	uc     *usecase.UseCase
	logger *logger.Logger
}

func NewHandler(uc *usecase.UseCase, logger *logger.Logger) *handler {
	return &handler{uc: uc, logger: logger}
}
