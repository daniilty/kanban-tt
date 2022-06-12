package pg

import (
	"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB interface {
	AddUser(context.Context, *User) error
	GetUser(context.Context, string) (*User, error)
	GetUserByEmail(context.Context, string) (*User, error)
	IsUserWithIDExists(context.Context, string) (bool, error)
	IsUserWithEmailExists(context.Context, string) (bool, error)
	IsUserWithEmailPasswordExists(context.Context, string, string) (bool, error)
	UpdateUser(context.Context, *User) error
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
