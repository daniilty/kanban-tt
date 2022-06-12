package core

import (
	"context"
	"fmt"

	"github.com/daniilty/kanban-tt/auth/claims"
	"github.com/daniilty/kanban-tt/schema"
)

type LoginData struct {
	Email    string
	Password string
}

func (s *ServiceImpl) Login(ctx context.Context, data *LoginData) (string, bool, error) {
	passwordHash := getMD5Sum(data.Password)

	resp, err := s.usersClient.IsValidUserCredentials(ctx, &schema.IsValidUserCredentialsRequest{
		Email:        data.Email,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return "", false, err
	}

	if !resp.IsValid {
		return "", true, fmt.Errorf("invalid credentials")
	}

	userResp, err := s.usersClient.GetUserByEmail(ctx, &schema.GetUserByEmailRequest{
		Email: data.Email,
	})
	if err != nil {
		return "", false, err
	}

	accessToken, err := s.jwtManager.Generate(&claims.Subject{
		UID: userResp.User.GetId(),
	})
	if err != nil {
		return "", false, err
	}

	return accessToken, true, nil
}
