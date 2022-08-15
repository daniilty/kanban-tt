package pg

import (
	"context"
	"database/sql"
	"errors"

	"github.com/daniilty/pgxquery"
	"github.com/jmoiron/sqlx"
)

type Status struct {
	pgxquery.TableName `db:"statuses"`

	ID       int    `db:"id,primarykey"`
	Name     string `db:"name,omitempty"`
	OwnerID  string `db:"owner_id,omitempty"`
	ParentID int    `db:"parent_id,omitempty"`
	ChildID  int    `db:"child_id,omitempty"`
}

func (s *Status) isInDB() bool {
	return s.ID > 0
}

func (d *db) GetStatusWithLowestPriority(ctx context.Context, uid string) (*Status, error) {
	const q = "select * from statuses where owner_id=$1 where parent_id=0"

	status := &Status{}

	err := d.db.QueryRowContext(ctx, q, uid).Scan(&status.ID, &status.Name, &status.OwnerID, &status.ParentID, &status.ChildID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoStatuses
		}
	}

	return status, err
}

func (d *db) AddParent(ctx context.Context, s *Status) (int, error) {
	const (
		selectRootQ = "select id, name, owner_id, parent_id, child_id from statuses where parent_id=0"
		updateRootQ = "update statuses set parent_id=$1 where id=$2"
	)

	tx, err := d.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return 0, err
	}

	root := &Status{}
	err = tx.QueryRowxContext(ctx, selectRootQ).Scan(
		&root.ID, &root.Name, &root.OwnerID, &root.ParentID, &root.ChildID,
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		tx.Rollback()
		return 0, err
	}

	s.ChildID = root.ID

	id, err := addStatusTX(tx, s)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if root.isInDB() {
		_, err = tx.ExecContext(ctx, updateRootQ, id, root.ID)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	err = tx.Commit()

	return id, err
}

func (d *db) AddChild(ctx context.Context, s *Status) (int, error) {
	const (
		selectStatusQ = "select id, name, owner_id, parent_id, child_id from statuses where id=$1"
		updateParentQ = "update statuses set parent_id=$1 where id=$2"
		updateChildQ  = "update statuses set child_id=$1 where id=$2"
	)

	tx, err := d.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return 0, err
	}

	parent := &Status{}
	err = tx.QueryRowxContext(ctx, selectStatusQ, s.ParentID).Scan(
		&parent.ID, &parent.Name, &parent.OwnerID, &parent.ParentID, &parent.ChildID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrNoStatuses
		}

		tx.Rollback()
		return 0, err
	}

	s.ChildID = parent.ChildID

	id, err := addStatusTX(tx, s)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if parent.ChildID != 0 {
		_, err = tx.ExecContext(ctx, updateParentQ, id, parent.ChildID)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	_, err = tx.ExecContext(ctx, updateChildQ, id, parent.ID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit()

	return id, err
}

func addStatusTX(tx *sqlx.Tx, s *Status) (int, error) {
	const q = "insert into statuses(name, owner_id, parent_id, child_id) values(:name, :owner_id, :parent_id, :child_id) returning id"

	var id int
	rows, err := tx.NamedQuery(q, s)
	if err != nil {
		return 0, err
	}

	for rows.Next() {
		rows.Scan(&id)
	}

	return id, err
}

func (d *db) AddStatus(ctx context.Context, s *Status) (int, error) {
	const q = "insert into statuses(name, owner_id, parent_id, child_id) values(:name, :owner_id, :parent_id, :child_id) returning id"

	var id int
	rows, err := d.db.NamedQueryContext(ctx, q, s)
	if err != nil {
		return 0, err
	}

	for rows.Next() {
		rows.Scan(&id)
	}

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
	const q = "select * from statuses where owner_id=$1"

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
