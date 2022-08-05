package core

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/daniilty/kanban-tt/auth/claims"
	"github.com/daniilty/kanban-tt/auth/internal/generate"
	"github.com/daniilty/kanban-tt/auth/internal/pg"
	"github.com/daniilty/kanban-tt/schema"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	tokenLen = 12

	CodeUserWithEmailExists Code = "USER_WITH_SUCH_EMAIL_EXISTS"
)

func (s *ServiceImpl) Register(ctx context.Context, user *UserInfo) (string, Code, error) {
	_, err := s.usersClient.GetUserByEmail(ctx, &schema.GetUserByEmailRequest{Email: user.Email})
	if err == nil {
		return "", CodeUserWithEmailExists, fmt.Errorf("user with such email already exists: %s", user.Email)
	}

	if status.Code(err) != codes.InvalidArgument {
		return "", CodeInternal, err
	}

	resp, err := s.usersClient.AddUser(ctx, convertUserInfoToAddUser(user))
	if err != nil {
		return "", CodeInternal, err
	}

	key, err := generate.SecureToken(tokenLen)
	if err != nil {
		return "", CodeInternal, err
	}

	now := time.Now()

	err = s.db.AddToken(ctx, &pg.Token{
		Key:       key,
		UID:       int(resp.GetId()),
		CreatedAt: &now,
	})
	if err != nil {
		return "", CodeInternal, err
	}

	confirmURL := generateConfirmLink(s.confirmURL, key)

	err = s.kafkaProducer.SendMessage(ctx, &schema.Email{
		To:  user.Email,
		Msg: "Please confirm your email with this link: " + confirmURL.String() + " or your account will be deleted in a week",
	})
	if err != nil {
		return "", CodeInternal, err
	}

	uid := strconv.Itoa(int(resp.GetId()))

	accessToken, err := s.jwtManager.Generate(&claims.Subject{
		UID: uid,
	})
	if err != nil {
		return "", CodeInternal, err
	}

	return accessToken, CodeOK, nil
}
