package usecase

import (
	"context"
	"time"

	"github.com/FurmanovVitaliy/pixel-cloud/internal/domain/game"
	"github.com/FurmanovVitaliy/pixel-cloud/internal/domain/user"
	"github.com/FurmanovVitaliy/pixel-cloud/pkg/logger"
)

// services
type userService interface {
	CreateUser(ctx context.Context, user user.User) (string, error)
	FindByEmail(ctx context.Context, email string) (user.User, error)
	GetList(ctx context.Context) ([]user.User, error)
}

type tokenService interface {
	CreateToken(userID, username string, expiry time.Duration) (string, error)
}

type gameService interface {
	GetAll(ctx context.Context) ([]game.Game, error)
	Create(ctx context.Context, game game.Game) (string, error)
	GetOneById(ctx context.Context, id string) (game.Game, error)
	DeleteAll(ctx context.Context) error
}

type UseCase struct {
	userService  userService
	gameService  gameService
	tokenService tokenService
	logger       *logger.Logger
}

func NewUseCase(
	userService userService,
	gameService gameService,
	tokenService tokenService,
	logger *logger.Logger,
) *UseCase {
	return &UseCase{
		userService:  userService,
		gameService:  gameService,
		tokenService: tokenService,
		logger:       logger,
	}
}
