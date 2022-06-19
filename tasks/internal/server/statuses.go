package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/daniilty/kanban-tt/auth/claims"
	"github.com/daniilty/kanban-tt/tasks/internal/core"
)

type status struct {
	Name     string `json:"name"`
	Priority uint32 `json:"priority"`
}

func (s *status) validate() error {
	if s.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	return nil
}

func (h *HTTP) handleGetStatuses(w http.ResponseWriter, r *http.Request) {
	resp := h.getStatusesResponse(r)

	resp.writeJSON(w)
}

func (h *HTTP) getStatusesResponse(r *http.Request) response {
	ctx := r.Context()
	sub := ctx.Value(subContextVal)
	if sub == nil {
		return getBadRequestWithMsgResponse("no subject")
	}

	s := sub.(*claims.Subject)

	tasks, err := h.service.GetStatuses(ctx, s.UID)
	if err != nil {
		return getInternalServerErrorResponse()
	}

	return newOKResponse(tasks)
}

func (h *HTTP) handleAddStatus(w http.ResponseWriter, r *http.Request) {
	resp := h.addStatusResponse(r)

	resp.writeJSON(w)
}

func (h *HTTP) addStatusResponse(r *http.Request) response {
	if r.Body == http.NoBody {
		return getBadRequestWithMsgResponse("no payload")
	}

	status := &status{}

	err := unmarshalReader(r.Body, status)
	if err != nil {
		return getBadRequestWithMsgResponse(fmt.Sprintf("bad body: %s", err.Error()))
	}

	err = status.validate()
	if err != nil {
		return getBadRequestWithMsgResponse(err.Error())
	}

	ctx := r.Context()
	sub := ctx.Value(subContextVal)
	if sub == nil {
		return getBadRequestWithMsgResponse("no subject")
	}

	s := sub.(*claims.Subject)

	err = h.service.AddStatus(ctx, &core.Status{
		Name:     status.Name,
		Priority: status.Priority,
		OwnerID:  s.UID,
	})
	if err != nil {
		if errors.Is(err, core.ErrStatusWithNameExists) {
			return getBadRequestWithMsgResponse(err.Error())
		}

		h.logger.Errorw("add task", "err", err)

		return getInternalServerErrorResponse()
	}

	return newOKResponse(struct{}{})
}
