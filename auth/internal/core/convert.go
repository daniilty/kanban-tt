package core

import "github.com/daniilty/kanban-tt/schema"

func convertUserInfoToAddUser(u *UserInfo) *schema.AddUserRequest {
	passwordHash := getMD5Sum(u.Password)

	return &schema.AddUserRequest{
		Email:        u.Email,
		Name:         u.Name,
		PasswordHash: passwordHash,
	}
}

func convertPBUserToUserInfo(u *schema.User) *UserInfo {
	return &UserInfo{
		ID:             u.Id,
		Email:          u.Email,
		Name:           u.Name,
		EmailConfirmed: u.EmailConfirmed,
		TaskTTL:        int(u.TaskTtl),
	}
}
