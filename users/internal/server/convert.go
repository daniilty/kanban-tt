package server

import (
	"strconv"

	"github.com/daniilty/kanban-tt/schema"
	"github.com/daniilty/kanban-tt/users/internal/core"
)

func convertPBAddUserToCore(user *schema.AddUserRequest) *core.User {
	return &core.User{
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}
}

func convertPBUpdateUserToCore(user *schema.UpdateUserRequest) *core.User {
	return &core.User{
		ID:             strconv.Itoa(int(user.Id)),
		Name:           user.Name,
		Email:          user.Email,
		PasswordHash:   user.PasswordHash,
		EmailConfirmed: user.EmailConfirmed,
		TaskTTL:        int(user.TaskTtl),
	}
}

func convertCoreUserToPB(user *core.User) *schema.User {
	return &schema.User{
		Id:             user.ID,
		Name:           user.Name,
		Email:          user.Email,
		EmailConfirmed: user.EmailConfirmed,
		TaskTtl:        int64(user.TaskTTL),
	}
}

func convertCoreUserToResponse(u *core.User) user {
	return user{
		ID:   u.ID,
		Name: u.Name,
	}
}

func convertPBUserToCore(user *schema.User) *core.User {
	return &core.User{
		ID:             user.Id,
		Name:           user.Name,
		Email:          user.Email,
		EmailConfirmed: user.EmailConfirmed,
	}
}

func convertCoreUsersToPB(users []*core.User) []*schema.User {
	converted := make([]*schema.User, 0, len(users))

	for i := range users {
		converted = append(converted, convertCoreUserToPB(users[i]))
	}

	return converted
}

func convertCoreUsersToResponse(users []*core.User) []user {
	converted := make([]user, 0, len(users))

	for i := range users {
		converted = append(converted, convertCoreUserToResponse(users[i]))
	}

	return converted
}
