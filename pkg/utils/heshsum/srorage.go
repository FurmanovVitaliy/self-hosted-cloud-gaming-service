package hashsum

import (
	"context"
)

type Storage interface {
	Create(ctx context.Context, id string, hash string) error
	FindOne(ctx context.Context, id string) (string, error)
	Update(ctx context.Context, id string, hash string) error
	Delete(ctx context.Context, id string) error
}
