package pg

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Token struct {
	Key       string     `db:"key"`
	UID       int        `db:"uid"`
	CreatedAt *time.Time `db:"created_at"`
}

func (d *db) AddToken(ctx context.Context, t *Token) error {
	const q = "insert into tokens(key, uid, created_at) values(:key, :uid, :created_at)"

	_, err := d.db.NamedExecContext(ctx, q, t)

	return err
}

func (d *db) GetToken(ctx context.Context, key string) (*Token, error) {
	const q = "select key, uid, created_at from tokens where key=$1"

	t := &Token{}
	err := d.db.QueryRowContext(ctx, q, key).Scan(&t.Key, &t.UID, &t.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoSuchToken
		}
	}

	return t, err
}

func (d *db) DeleteToken(ctx context.Context, key string) error {
	const q = "delete from tokens where key=$1"

	_, err := d.db.ExecContext(ctx, q, key)

	return err
}
