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

	AddParent(context.Context, *Status) (int, error)
	AddChild(context.Context, *Status) (int, error)
	AddStatus(context.Context, *Status) (int, error)
	IsStatusWithIDExists(context.Context, int) (bool, error)
	IsStatusWithNameExists(context.Context, string, string) (bool, error)
	GetStatuses(context.Context, string) ([]*Status, error)
	GetStatus(context.Context, int) (*Status, error)
	GetStatusWithLowestPriority(context.Context, string) (*Status, error)
	UpdateStatusName(context.Context, int, string) error
	UpdateStatusParent(context.Context, *Status, int) error
	DeleteStatus(context.Context, *Status) error
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
