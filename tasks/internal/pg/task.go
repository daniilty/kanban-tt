package pg

import (
	"context"
	"time"
)

type Task struct {
	ID        int        `db:"id"`
	Content   string     `db:"content"`
	Priority  int        `db:"priority"`
	OwnerID   string     `db:"owner_id"`
	StatusID  int        `db:"status_id"`
	CreatedAt *time.Time `db:"created_at"`
}

func (d *db) AddTask(ctx context.Context, t *Task) error {
	const q = "insert into tasks(content, priority, owner_id, status_id, created_at) values(:content, :priority, :owner_id, :status_id, :created_at)"

	_, err := d.db.NamedExecContext(ctx, q, t)

	return err
}

func (d *db) GetTasks(ctx context.Context, uid string) ([]*Task, error) {
	const q = "select * from tasks where owner_id=$1"

	tasks := []*Task{}
	err := d.db.SelectContext(ctx, &tasks, q, uid)

	return tasks, err
}

func (d *db) UpdateTask(ctx context.Context, t *Task) error {
	const q = "update tasks set content=coalesce(:content, content), priority=coalesce(:priority, priority), status_id=coalesce(:status_id, status_id) where id=:id"

	_, err := d.db.NamedExecContext(ctx, q, t)

	return err
}

func (d *db) DeleteTask(ctx context.Context, id string) error {
	const q = "delete from tasks where id=$1"

	_, err := d.db.ExecContext(ctx, q, id)

	return err
}
