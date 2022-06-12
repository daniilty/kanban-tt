package core

import (
	"context"
	"fmt"

	"github.com/daniilty/kanban-tt/schema"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserInfo struct {
	ID             string
	Name           string
	Email          string
	EmailConfirmed bool
	// Used for registration
	Password string
}

func (s *ServiceImpl) GetUserInfo(ctx context.Context, accessToken string) (*UserInfo, bool, error) {
	sub, err := s.jwtManager.ParseRawToken(accessToken)
	if err != nil {
		return nil, true, fmt.Errorf("invalid access token provided")
	}

	resp, err := s.usersClient.GetUser(ctx, &schema.GetUserRequest{
		Id: sub.UID,
	})
	if err != nil {
		if status.Code(err) == codes.InvalidArgument {
			return nil, true, fmt.Errorf("invalid user id")
		}

		return nil, false, err
	}

	user := resp.GetUser()

	return convertPBUserToUserInfo(user), true, nil
}
