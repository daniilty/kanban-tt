package pg

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type parentSetFunc func(context.Context, *sqlx.Tx, *Status, int) (*Status, error)

type Status struct {
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
	const q = "select id, name, owner_id, parent_id, child_id from statuses where owner_id=$1 and parent_id=0"

	status := &Status{}

	err := d.db.QueryRowContext(ctx, q, uid).Scan(
		&status.ID, &status.Name, &status.OwnerID, &status.ParentID, &status.ChildID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoStatuses
		}
	}

	return status, err
}

func (d *db) AddParent(ctx context.Context, s *Status) (int, error) {
	const (
		selectRootQ = "select id, name, owner_id, parent_id, child_id from statuses where parent_id=0 and owner_id=$1"
		updateRootQ = "update statuses set parent_id=$1 where id=$2"
	)

	tx, err := d.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return 0, err
	}

	root := &Status{}
	err = tx.QueryRowxContext(ctx, selectRootQ, s.OwnerID).Scan(
		&root.ID, &root.Name, &root.OwnerID, &root.ParentID, &root.ChildID,
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		tx.Rollback()
		return 0, err
	}

	s.ChildID = root.ID

	id, err := addStatusTx(tx, s)
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

	id, err := addStatusTx(tx, s)
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

func addStatusTx(tx *sqlx.Tx, s *Status) (int, error) {
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

func (d *db) IsStatusWithNameExists(ctx context.Context, name string, uid string) (bool, error) {
	const q = "select exists(select from statuses where name=$1 and owner_id=$2)"

	exists := false
	err := d.db.GetContext(ctx, &exists, q, name, uid)

	return exists, err
}

func (d *db) GetStatuses(ctx context.Context, uid string) ([]*Status, error) {
	const q = "select * from statuses where owner_id=$1"

	statuses := []*Status{}
	err := d.db.SelectContext(ctx, &statuses, q, uid)

	return statuses, err
}

func (d *db) UpdateStatusName(ctx context.Context, id int, name string) error {
	const q = "update statuses set name=$1 where id=$2"

	_, err := d.db.ExecContext(ctx, q, name, id)

	return err
}

func (d *db) UpdateStatusParent(ctx context.Context, s *Status, parentID int) error {
	const (
		updateParentQ = "update statuses set child_id=$1 where id=$2"
		updateChildQ  = "update statuses set parent_id=$1 where id=$2"
		updateStatusQ = "update statuses set parent_id=$1, child_id=$2 where id=$3"
	)

	opts := &sql.TxOptions{}
	tx, err := d.db.BeginTxx(ctx, opts)
	if err != nil {
		return err
	}

	err = removeStatusFromListTx(ctx, tx, s)
	if err != nil {
		tx.Rollback()
		return err
	}

	set := getParentSetFunc(parentID)

	parent, err := set(ctx, tx, s, parentID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, updateStatusQ, parent.ID, parent.ChildID, s.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()

	return err
}

func (d *db) DeleteStatus(ctx context.Context, s *Status) error {
	const (
		deleteQ       = "delete from statuses where id=$1"
		updateParentQ = "update statuses set child_id=$1 where id=$2"
		updateChildQ  = "update statuses set parent_id=$1 where id=$2"
	)

	// short path
	if s.ParentID == 0 && s.ChildID == 0 {
		_, err := d.db.ExecContext(ctx, deleteQ, s.ID)
		return err
	}

	opts := &sql.TxOptions{}
	tx, err := d.db.BeginTxx(ctx, opts)
	if err != nil {
		return err
	}

	if s.ParentID != 0 {
		_, err = tx.ExecContext(ctx, updateParentQ, s.ChildID, s.ParentID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if s.ChildID != 0 {
		_, err = tx.ExecContext(ctx, updateChildQ, s.ParentID, s.ChildID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	_, err = tx.ExecContext(ctx, deleteQ, s.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()

	return err
}

func (d *db) GetStatus(ctx context.Context, id int) (*Status, error) {
	const selectQ = "select id, name, owner_id, parent_id, child_id from statuses where id=$1"

	s := &Status{}
	err := d.db.QueryRowContext(ctx, selectQ, id).Scan(
		&s.ID, &s.Name, &s.OwnerID, &s.ParentID, &s.ChildID,
	)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoStatuses
	}

	return s, err
}

func removeStatusFromListTx(ctx context.Context, tx *sqlx.Tx, s *Status) error {
	const (
		updateParentQ = "update statuses set child_id=$1 where id=$2"
		updateChildQ  = "update statuses set parent_id=$1 where id=$2"
	)

	// set a new child for a task's parent
	if s.ParentID != 0 {
		_, err := tx.ExecContext(ctx, updateParentQ, s.ChildID, s.ParentID)
		if err != nil {
			return err
		}
	}

	// set a new parent for a task's child
	if s.ChildID != 0 {
		_, err := tx.ExecContext(ctx, updateChildQ, s.ParentID, s.ChildID)
		if err != nil {
			return err
		}
	}

	return nil
}

func getStatusTx(tx *sqlx.Tx, id int) (*Status, error) {
	const selectQ = "select id, name, owner_id, parent_id, child_id from statuses where id=$1"

	s := &Status{}
	err := tx.QueryRow(selectQ, id).Scan(
		&s.ID, &s.Name, &s.OwnerID, &s.ParentID, &s.ChildID,
	)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoStatuses
	}

	return s, err
}

func getHeadStatusTx(tx *sqlx.Tx, uid string) (*Status, error) {
	const q = "select id, name, owner_id, parent_id, child_id from statuses where owner_id=$1 and parent_id=0"

	s := &Status{}
	err := tx.QueryRow(q, uid).Scan(
		&s.ID, &s.Name, &s.OwnerID, &s.ParentID, &s.ChildID,
	)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoStatuses
	}

	return s, err
}

func getParentSetFunc(parentID int) parentSetFunc {
	if parentID == 0 {
		return setNewHeadStatusTx
	}

	return setNewParentStatusTx
}

func setNewHeadStatusTx(ctx context.Context, tx *sqlx.Tx, s *Status, _ int) (*Status, error) {
	const updateChildQ = "update statuses set parent_id=$1 where id=$2"

	head, err := getHeadStatusTx(tx, s.OwnerID)
	if err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx, updateChildQ, s.ID, head.ID)
	if err != nil {
		return nil, err
	}

	// set fake parent
	head.ChildID = head.ID
	head.ID = 0

	return head, nil
}

func setNewParentStatusTx(ctx context.Context, tx *sqlx.Tx, s *Status, parentID int) (*Status, error) {
	const updateParentQ = "update statuses set child_id=$1 where id=$2"

	parent, err := getStatusTx(tx, parentID)
	if err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx, updateParentQ, s.ID, parent.ID)
	if err != nil {
		return nil, err
	}

	return parent, nil
}
