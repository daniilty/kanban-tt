package pg

import (
	"context"
	"database/sql"
	"errors"

	"github.com/daniilty/pgxquery"
)

type Status struct {
	pgxquery.TableName `db:"statuses"`

	ID       int    `db:"id,primarykey"`
	Name     string `db:"name,omitempty"`
	Priority int    `db:"priority,omitempty"`
	OwnerID  string `db:"owner_id,omitempty"`
}

func (d *db) GetStatusWithLowestPriority(ctx context.Context, uid string) (*Status, error) {
	const q = "select * from statuses where owner_id=$1 order by priority asc limit 1"

	status := &Status{}

	err := d.db.QueryRowContext(ctx, q, uid).Scan(&status.ID, &status.Name, &status.Priority, &status.OwnerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoStatuses
		}
	}

	return status, err
}

func (d *db) AddStatus(ctx context.Context, s *Status) (int, error) {
	const q = "insert into statuses(name, priority, owner_id) values(:name, :priority, :owner_id) returning id"

	var id int
	err := d.db.QueryRowContext(ctx, q, s).Scan(&id)

	return id, err
}

func (d *db) IsStatusWithIDExists(ctx context.Context, id int) (bool, error) {
	const q = "select exists(select from statuses where id=$1)"

	exists := false
	err := d.db.GetContext(ctx, &exists, q, id)

	return exists, err
}

func (d *db) IsStatusWithNameExists(ctx context.Context, name string) (bool, error) {
	const q = "select exists(select from statuses where name=$1)"

	exists := false
	err := d.db.GetContext(ctx, &exists, q, name)

	return exists, err
}

func (d *db) GetStatuses(ctx context.Context, uid string) ([]*Status, error) {
	const q = "select * from statuses where owner_id=$1 order by priority asc"

	statuses := []*Status{}
	err := d.db.SelectContext(ctx, &statuses, q, uid)

	return statuses, err
}

func (d *db) UpdateStatus(ctx context.Context, s *Status) error {
	q, err := pgxquery.GenerateNamedUpdate(s)
	if err != nil {
		if errors.Is(err, pgxquery.ErrEmptyModel) {
			return ErrEmptyModel
		}

		return err
	}

	_, err = d.db.NamedExecContext(ctx, q, s)

	return err
}

func (d *db) DeleteStatus(ctx context.Context, id int) error {
	const q = "delete from statuses where id=$1"

	_, err := d.db.ExecContext(ctx, q, id)

	return err
}
