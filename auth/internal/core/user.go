package core

import (
	"context"
	"errors"
	"fmt"

	"github.com/daniilty/kanban-tt/auth/internal/pg"
	"github.com/daniilty/kanban-tt/schema"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	CodeInvalidAccessToken Code = "INVALID_ACCESS_TOKEN"
	CodeInvalidUserID      Code = "INVALID_USER_ID"
)

type UserInfo struct {
	ID             string
	Name           string
	Email          string
	EmailConfirmed bool
	// Used for registration
	Password string
}

func (s *ServiceImpl) GetUserInfo(ctx context.Context, accessToken string) (*UserInfo, Code, error) {
	sub, err := s.jwtManager.ParseRawToken(accessToken)
	if err != nil {
		return nil, CodeInvalidAccessToken, fmt.Errorf("invalid access token provided")
	}

	resp, err := s.usersClient.GetUser(ctx, &schema.GetUserRequest{
		Id: sub.UID,
	})
	if err != nil {
		if status.Code(err) == codes.InvalidArgument {
			return nil, CodeInvalidUserID, fmt.Errorf("invalid user id")
		}

		return nil, CodeInternal, err
	}

	user := resp.GetUser()

	return convertPBUserToUserInfo(user), CodeOK, nil
}

func (s *ServiceImpl) ConfirmUserEmail(ctx context.Context, key string) error {
	t, err := s.db.GetToken(ctx, key)
	if err != nil {
		if errors.Is(err, pg.ErrNoSuchToken) {
			return ErrNoSuchKey
		}

		return err
	}

	_, err = s.usersClient.UpdateUser(ctx, &schema.UpdateUserRequest{
		Id:             int64(t.UID),
		EmailConfirmed: true,
	})
	if err != nil {
		return err
	}

	return s.db.DeleteToken(ctx, t.Key)
}
