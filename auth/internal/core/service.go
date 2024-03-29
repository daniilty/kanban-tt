package core

import (
	"context"
	"net/url"

	"github.com/daniilty/kanban-tt/auth/claims"
	"github.com/daniilty/kanban-tt/auth/internal/jwt"
	"github.com/daniilty/kanban-tt/auth/internal/kafka"
	"github.com/daniilty/kanban-tt/auth/internal/pg"
	schema "github.com/daniilty/kanban-tt/schema"
)

var _ Service = (*ServiceImpl)(nil)

type Service interface {
	Login(context.Context, *LoginData) (string, Code, error)
	Register(context.Context, *UserInfo) (string, Code, error)
	ParseRawToken(string) (*claims.Subject, error)
	RefreshSession(string) (string, error)
	GetUserInfo(context.Context, string) (*UserInfo, Code, error)
	GetTTLs(context.Context) []int64
	UpdateUser(context.Context, *UserInfo) (Code, error)
	ConfirmUserEmail(context.Context, string) error
	JWKS() []byte
}

type ServiceImpl struct {
	usersClient   schema.UsersClient
	jwtManager    jwt.Manager
	kafkaProducer kafka.Producer
	confirmURL    *url.URL
	db            pg.DB
}

func NewServiceImpl(usersClient schema.UsersClient, jwtManager jwt.Manager, kafkaProducer kafka.Producer, confirmURL *url.URL, db pg.DB) *ServiceImpl {
	return &ServiceImpl{
		usersClient:   usersClient,
		jwtManager:    jwtManager,
		kafkaProducer: kafkaProducer,
		confirmURL:    confirmURL,
		db:            db,
	}
}
