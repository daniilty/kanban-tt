package pg

import (
	"context"
	"time"

	"github.com/daniilty/pgxquery"
)

type User struct {
	pgxquery.TableName `db:"users"`

	ID             string     `db:"id,primarykey"`
	Name           string     `db:"name,omitempty"`
	Email          string     `db:"email,omitempty"`
	PasswordHash   string     `db:"password_hash,omitempty"`
	EmailConfirmed bool       `db:"email_confirmed,omitempty"`
	TaskTTL        int        `db:"task_ttl,omitempty"`
	CreatedAt      *time.Time `db:"created_at,omitempty"`
}

func (d *db) AddUser(ctx context.Context, u *User) (int, error) {
	const q = "insert into users(name, email, email_confirmed, password_hash, task_ttl, created_at) values($1, $2, $3, $4, $5, $6) returning id;"

	var id int

	err := d.db.QueryRowContext(ctx, q, u.Name, u.Email, u.EmailConfirmed, u.PasswordHash, u.TaskTTL, u.CreatedAt).Scan(&id)

	return id, err
}

func (d *db) GetUser(ctx context.Context, id string) (*User, error) {
	const q = "select * from users where id=$1"

	u := &User{}
	err := d.db.GetContext(ctx, u, q, id)

	return u, err
}

func (d *db) GetUserTaskTTL(ctx context.Context, id string) (int, error) {
	const q = "select task_ttl from users where id=$1"

	ttl := 0
	err := d.db.GetContext(ctx, &ttl, q, id)

	return ttl, err
}

func (d *db) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	const q = "select * from users where email=$1"

	u := &User{}
	err := d.db.GetContext(ctx, u, q, email)
	return u, err
}

func (d *db) IsUserWithIDExists(ctx context.Context, id string) (bool, error) {
	const q = "select exists(select from users where id=$1)"

	exists := false
	err := d.db.GetContext(ctx, &exists, q, id)

	return exists, err
}

func (d *db) IsUserWithEmailExists(ctx context.Context, email string) (bool, error) {
	const q = "select exists(select from users where email=$1)"

	exists := false
	err := d.db.GetContext(ctx, &exists, q, email)

	return exists, err
}

func (d *db) IsUserWithEmailPasswordExists(ctx context.Context, email string, passwordHash string) (bool, error) {
	const q = "select exists(select from users where email=$1 and password_hash=$2)"

	exists := false
	err := d.db.GetContext(ctx, &exists, q, email, passwordHash)

	return exists, err
}

func (d *db) UpdateUser(ctx context.Context, u *User) error {
	q, err := pgxquery.GenerateNamedUpdate(u)
	if err != nil {
		return err
	}

	_, err = d.db.NamedExecContext(ctx, q, u)

	return err
}
