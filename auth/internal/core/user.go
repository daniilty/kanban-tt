package core

import (
	"context"
	"errors"
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
	CodeInvalidAccessToken Code = "INVALID_ACCESS_TOKEN"
	CodeInvalidUserID      Code = "INVALID_USER_ID"
	CodeInvalidData        Code = "INVALID_DATA"
)

type UserInfo struct {
	ID             string
	Name           string
	Email          string
	EmailConfirmed bool
	// Used for registration
	Password string
	TaskTTL  int
}

func (u *UserInfo) hasOneChangedField() bool {
	switch {
	case u.Name != "":
		return true
	case u.Email != "":
		return true
	case u.Password != "":
		return true
	case u.TaskTTL != 0:
		return true
	default:
		return false
	}
}

func (u *UserInfo) toUpdateUser() (*schema.UpdateUserRequest, error) {
	id, err := strconv.Atoi(u.ID)
	if err != nil {
		return nil, err
	}

	passwordHash := ""
	if u.Password != "" {
		passwordHash = getMD5Sum(u.Password)
	}

	return &schema.UpdateUserRequest{
		Id:           int64(id),
		Name:         u.Name,
		Email:        u.Email,
		PasswordHash: passwordHash,
		TaskTtl:      int64(u.TaskTTL),
	}, nil
}

func (s *ServiceImpl) ParseRawToken(token string) (*claims.Subject, error) {
	return s.jwtManager.ParseRawToken(token)
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

func (s *ServiceImpl) UpdateUser(ctx context.Context, user *UserInfo) (Code, error) {
	// if we don't need to change anything just fuck it
	if !user.hasOneChangedField() {
		return CodeOK, nil
	}

	req, err := user.toUpdateUser()
	if err != nil {
		return CodeInvalidUserID, errors.New("invalid user id")
	}

	if user.Email != "" {
		resp, err := s.usersClient.GetUser(ctx, &schema.GetUserRequest{
			Id: user.ID,
		})
		if err != nil {
			if status.Code(err) == codes.InvalidArgument {
				return CodeInvalidUserID, fmt.Errorf("invalid user id")
			}

			return CodeInternal, err
		}

		if user.Email != resp.User.Email {
			_, err = s.usersClient.UnconfirmUserEmail(ctx, &schema.UnconfirmUserEmailRequest{
				Id: req.Id,
			})
			if err != nil {
				return CodeInternal, err
			}
		}

		key, err := generate.SecureToken(tokenLen)
		if err != nil {
			return CodeInternal, err
		}

		now := time.Now()

		err = s.db.AddToken(ctx, &pg.Token{
			Key:       key,
			UID:       int(req.GetId()),
			CreatedAt: &now,
		})
		if err != nil {
			return CodeInternal, err
		}

		confirmURL := generateConfirmLink(s.confirmURL, key)

		err = s.kafkaProducer.SendMessage(ctx, &schema.Email{
			To:  user.Email,
			Msg: "Please confirm your new email with this link: " + confirmURL.String() + " or your account will be deleted in a week",
		})
		if err != nil {
			return CodeInternal, err
		}
	}

	_, err = s.usersClient.UpdateUser(ctx, req)
	if err != nil {
		if status.Code(err) == codes.InvalidArgument {
			return CodeInvalidData, err
		}

		return CodeInternal, err
	}

	return CodeOK, nil
}
