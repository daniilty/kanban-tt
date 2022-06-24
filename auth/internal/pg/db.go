package pg

import (
	"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB interface {
	AddToken(context.Context, *Token) error
	GetToken(context.Context, string) (*Token, error)
	DeleteToken(context.Context, string) error
}

func Connect(ctx context.Context, addr string) (DB, error) {
	d, err := sqlx.ConnectContext(ctx, "postgres", addr)
	if err != nil {
		return nil, err
	}

	return &db{
		db: d,
	}, nil
}

type db struct {
	db *sqlx.DB
}
