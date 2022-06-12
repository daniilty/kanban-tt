package pg

import (
	"context"
)

type User struct {
	ID             string `db:"id"`
	Name           string `db:"name"`
	Email          string `db:"email"`
	PasswordHash   string `db:"password_hash"`
	EmailConfirmed bool   `db:"email_confirmed"`
}

func (d *db) AddUser(ctx context.Context, u *User) error {
	const q = "insert into users(name, email, email_confirmed, password_hash) values(:name, :email, :email_confirmed, :password_hash)"

	_, err := d.db.NamedExecContext(ctx, q, u)

	return err
}

func (d *db) GetUser(ctx context.Context, id string) (*User, error) {
	const q = "select * from users where id=$1"

	u := &User{}
	err := d.db.GetContext(ctx, u, q, id)

	return u, err
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
	const q = "update users set email=coalesce(:email,email), email_confirmed=coalesce(:email_confirmed,email_confirmed), password_hash=coalesce(:password_hash,password_hash), name=coalesce(:name, name) where id=:id"

	_, err := d.db.NamedExecContext(ctx, q, u)

	return err
}
