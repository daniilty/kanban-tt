package core

import (
	"context"

	"github.com/daniilty/kanban-tt/schema"
)

func (s *ServiceImpl) GetTTLs(ctx context.Context) []int64 {
	resp, err := s.usersClient.GetTTLs(ctx, &schema.GetTTLsRequest{})
	if err != nil {
		// yeah, whatever...
		return []int64{}
	}

	return resp.GetTtls()
}
