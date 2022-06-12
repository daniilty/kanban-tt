package server

import (
	"github.com/daniilty/kanban-tt/schema"
	"github.com/daniilty/kanban-tt/users/internal/core"
)

type GRPC struct {
	schema.UnimplementedUsersServer

	service core.Service
}

func NewGRPC(usersService core.Service) *GRPC {
	return &GRPC{
		service: usersService,
	}
}
