package core

import (
	"context"

	"github.com/daniilty/kanban-tt/auth/internal/jwt"
	"github.com/daniilty/kanban-tt/auth/internal/kafka"
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
	usersClient   schema.UsersClient
	jwtManager    jwt.Manager
	kafkaProducer kafka.Producer
}

func NewServiceImpl(usersClient schema.UsersClient, jwtManager jwt.Manager, kafkaProducer kafka.Producer) *ServiceImpl {
	return &ServiceImpl{
		usersClient:   usersClient,
		jwtManager:    jwtManager,
		kafkaProducer: kafkaProducer,
	}
}
