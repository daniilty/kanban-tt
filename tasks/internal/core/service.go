package core

import (
	"context"

	"github.com/daniilty/kanban-tt/tasks/internal/pg"
)

type Service interface {
	AddTask(context.Context, *Task) (error, Code)
	GetTasks(context.Context, string) ([]*Task, error)
	UpdateTask(context.Context, *Task) (error, Code)
	DeleteTask(context.Context, int) error

	AddStatus(context.Context, *Status) error
	GetStatuses(context.Context, string) ([]*Status, error)
	UpdateStatus(context.Context, *Status) (error, Code)
	DeleteStatus(context.Context, int) error
}

type service struct {
	db pg.DB
}

func NewService(db pg.DB) Service {
	return &service{
		db: db,
	}
}
