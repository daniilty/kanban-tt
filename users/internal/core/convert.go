package core

import "github.com/daniilty/kanban-tt/users/internal/pg"

func (n *User) toDB() *pg.User {
	return &pg.User{
		ID:             n.ID,
		Name:           n.Name,
		Email:          n.Email,
		EmailConfirmed: n.EmailConfirmed,
		PasswordHash:   n.PasswordHash,
		TaskTTL:        n.TaskTTL,
	}
}

func convertDBUserToService(user *pg.User) *User {
	return &User{
		ID:             user.ID,
		Name:           user.Name,
		Email:          user.Email,
		EmailConfirmed: user.EmailConfirmed,
		PasswordHash:   user.PasswordHash,
	}
}
