package client

import (
	"context"
	"fmt"
	"time"

	"github.com/FurmanovVitaliy/pixel-cloud/util"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgresClient interface {
	Exec(ctx context.Context, sql string, arguments ...any)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func NewPostgresClient(ctx context.Context, maxAttempts int, host, port, username, password, database string) (pool *pgxpool.Pool, err error) {

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, database)
	util.DoWithRetry(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		pool, err = pgxpool.Connect(ctx, dsn)
		if err != nil {
			return err
		}
		return nil

	}, maxAttempts, 5*time.Second)

	return pool, err
}
