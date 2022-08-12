package pg

import (
	"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB interface {
	AddTask(context.Context, *Task) (int, error)
	GetUserTasks(context.Context, string) ([]*Task, error)
	GetTasks(context.Context) ([]*Task, error)
	UpdateTask(context.Context, *Task) error
	DeleteTask(context.Context, int) error
	DeleteExpiredTasks(context.Context, string, int) error

	AddStatus(context.Context, *Status) (int, error)
	IsStatusWithIDExists(context.Context, int) (bool, error)
	IsStatusWithNameExists(context.Context, string) (bool, error)
	GetStatuses(context.Context, string) ([]*Status, error)
	GetStatusWithLowestPriority(context.Context, string) (*Status, error)
	UpdateStatus(context.Context, *Status) error
	DeleteStatus(context.Context, int) error
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
