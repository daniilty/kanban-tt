package core

import (
	"context"

	"github.com/daniilty/kanban-tt/auth/internal/jwt"
	schema "github.com/daniilty/kanban-tt/schema"
)

var _ Service = (*ServiceImpl)(nil)

type Service interface {
	Login(context.Context, *LoginData) (string, bool, error)
	Register(context.Context, *UserInfo) (string, bool, error)
	RefreshSession(string) (string, error)
	GetUserInfo(context.Context, string) (*UserInfo, bool, error)
	JWKS() []byte
}

type ServiceImpl struct {
	usersClient schema.UsersClient
	jwtManager  jwt.Manager
}

func NewServiceImpl(usersClient schema.UsersClient, jwtManager jwt.Manager) *ServiceImpl {
	return &ServiceImpl{
		usersClient: usersClient,
		jwtManager:  jwtManager,
	}
}
