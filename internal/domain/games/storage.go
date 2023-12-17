package games

import (
	"context"
)

type Storage interface {
	Create(ctx context.Context, game Game) (string, error)
	FindAll(ctx context.Context) (games []Game, err error)
	FindOne(ctx context.Context, id string) (Game, error)
	FindOneByName(ctx context.Context, name string) (Game, error)
	Update(ctx context.Context, game Game) error
	FullyUpdate(ctx context.Context, games []Game) error
	Delete(ctx context.Context, id string) error
	Drop(ctx context.Context) error
}
