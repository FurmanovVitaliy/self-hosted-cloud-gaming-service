package game

import (
	"context"

	"github.com/FurmanovVitaliy/pixel-cloud/pkg/errors"
)

type gameStorage interface {
	Create(ctx context.Context, game Game) (string, error)
	FindAll(ctx context.Context) (games []Game, err error)
	FindOne(ctx context.Context, id string) (Game, error)
	Drop(ctx context.Context) error
}

type gameService struct {
	storage gameStorage
}

func NewGameService(storage gameStorage) *gameService {
	return &gameService{
		storage: storage,
	}
}

func (s *gameService) GetAll(ctx context.Context) ([]Game, error) {
	games, err := s.storage.FindAll(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "GS", "000001", "failed to get list of games")
	}
	return games, nil
}

func (s *gameService) Create(ctx context.Context, game Game) (string, error) {
	id, err := s.storage.Create(ctx, game)
	if err != nil {
		return "", errors.Wrap(err, "GS", "000002", "failed to create game")
	}
	return id, nil
}

func (s *gameService) GetOneById(ctx context.Context, id string) (Game, error) {
	game, err := s.storage.FindOne(ctx, id)
	if err != nil {
		return Game{}, errors.Wrap(err, "GS", "000003", "failed to find game by id")
	}
	return game, nil
}
func (s *gameService) DeleteAll(ctx context.Context) error {
	err := s.storage.Drop(ctx)
	if err != nil {
		return errors.Wrap(err, "GS", "000004", "failed to delete all games")
	}
	return nil
}
