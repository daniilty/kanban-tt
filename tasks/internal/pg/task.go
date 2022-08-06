package pg

import (
	"context"
	"errors"
	"time"

	"github.com/daniilty/pgxquery"
)

type Task struct {
	pgxquery.TableName `db:"tasks"`

	ID        int        `db:"id,primarykey"`
	Content   string     `db:"content,omitempty"`
	Priority  int        `db:"priority,omitempty"`
	OwnerID   string     `db:"owner_id,omitempty"`
	StatusID  int        `db:"status_id,omitempty"`
	CreatedAt *time.Time `db:"created_at,omitempty"`
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
	q, err := pgxquery.GenerateNamedUpdate(t)
	if err != nil {
		if errors.Is(err, pgxquery.ErrEmptyModel) {
			return ErrEmptyModel
		}

		return err
	}

	_, err = d.db.NamedExecContext(ctx, q, t)

	return err
}

func (d *db) DeleteTask(ctx context.Context, id int) error {
	const q = "delete from tasks where id=$1"

	_, err := d.db.ExecContext(ctx, q, id)

	return err
}

func (d *db) DeleteExpiredTasks(ctx context.Context) error {
	const q = "delete from tasks where CURRENT_DATE - created_at > 0"

	_, err := d.db.ExecContext(ctx, q)

	return err
}
