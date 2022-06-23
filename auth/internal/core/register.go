package core

import (
	"context"
	"fmt"
	"strconv"

	"github.com/daniilty/kanban-tt/auth/claims"
	"github.com/daniilty/kanban-tt/schema"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServiceImpl) Register(ctx context.Context, user *UserInfo) (string, bool, error) {
	_, err := s.usersClient.GetUserByEmail(ctx, &schema.GetUserByEmailRequest{Email: user.Email})
	if err == nil {
		return "", true, fmt.Errorf("user with such email already exists: %s", user.Email)
	}

	if status.Code(err) != codes.InvalidArgument {
		return "", false, err
	}

	resp, err := s.usersClient.AddUser(ctx, convertUserInfoToAddUser(user))
	if err != nil {
		return "", false, err
	}

	uid := strconv.Itoa(int(resp.GetId()))

	accessToken, err := s.jwtManager.Generate(&claims.Subject{
		UID: uid,
	})
	if err != nil {
		return "", false, err
	}

	return accessToken, true, nil
}
