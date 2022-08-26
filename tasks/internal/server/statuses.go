package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/daniilty/kanban-tt/auth/claims"
	"github.com/daniilty/kanban-tt/tasks/internal/core"
	"github.com/gorilla/mux"
)

const (
	// пустое имя
	codeNameEmpty core.Code = "EMPTY_NAME"
	// статус с таким именем уже существует
	codeStatusWithNameExists core.Code = "STATUS_WITH_NAME_EXISTS"
	// родителя с таким айди не существует
	codeParentDoesNotExists core.Code = "PARENT_WITH_ID_DOES_NOT_EXISTS"
)

// swagger:model
type status struct {
	// required: true
	ID int `json:"id"`
	// required: true
	Name string `json:"name"`
	// required: true
	ParentID uint32 `json:"parentId"`
}

func (s *status) validate() (error, core.Code) {
	if s.Name == "" {
		return fmt.Errorf("name cannot be empty"), codeNameEmpty
	}

	return nil, core.CodeOK
}

// swagger:route GET /api/v1/tasks/statuses Status statusesGet
// get user created statuses
//
// security:
//    api-key: Bearer
// Returns operation result
// responses:
//    200: status
//    400: errorResponse Bad request
//    401: errorResponse Unauthorized
//    500: errorResponse Internal server error
func (h *HTTP) handleGetStatuses(w http.ResponseWriter, r *http.Request) {
	resp := h.getStatusesResponse(r)

	resp.writeJSON(w)
}

func (h *HTTP) getStatusesResponse(r *http.Request) response {
	ctx := r.Context()
	sub := ctx.Value(subContextVal)
	if sub == nil {
		return getUnauthorizedResponse(codeUnauthorizedNoSub)
	}

	s := sub.(*claims.Subject)

	tasks, err := h.service.GetStatuses(ctx, s.UID)
	if err != nil {
		return getInternalServerErrorResponse(core.CodeDBFail)
	}

	return newOKResponse(tasks)
}

// swagger:route POST /api/v1/tasks/status Status statusAdd
// Add status
//
// security:
//    api-key: []
//
// parameters:
//  + name: status
//    in: body
//    required: true
//    type: status
//
// Returns operation result
// responses:
//    200: addResponse
//    400: errorResponse Bad request
//    401: errorResponse Unauthorized
//    500: errorResponse Internal server error
func (h *HTTP) handleAddStatus(w http.ResponseWriter, r *http.Request) {
	resp := h.addStatusResponse(r)

	resp.writeJSON(w)
}

func (h *HTTP) addStatusResponse(r *http.Request) response {
	if r.Body == http.NoBody {
		return getBadRequestWithMsgResponse("no payload", codeEmptyBody)
	}

	status := &status{}

	err := unmarshalReader(r.Body, status)
	if err != nil {
		return getBadRequestWithMsgResponse(fmt.Sprintf("bad body: %s", err.Error()), codeInvalidBodyStructure)
	}

	err, code := status.validate()
	if err != nil {
		return getBadRequestWithMsgResponse(err.Error(), code)
	}

	ctx := r.Context()
	sub := ctx.Value(subContextVal)
	if sub == nil {
		return getUnauthorizedResponse(codeUnauthorizedNoSub)
	}

	s := sub.(*claims.Subject)

	id, err := h.service.AddStatus(ctx, &core.Status{
		Name:     status.Name,
		ParentID: int(status.ParentID),
		OwnerID:  s.UID,
	})
	if err != nil {
		if errors.Is(err, core.ErrStatusWithNameExists) {
			return getBadRequestWithMsgResponse(err.Error(), codeStatusWithNameExists)
		}

		if errors.Is(err, core.ErrNoSuchParent) {
			return getBadRequestWithMsgResponse(err.Error(), codeParentDoesNotExists)
		}

		h.logger.Errorw("Add status.", "err", err)

		return getInternalServerErrorResponse(core.CodeDBFail)
	}

	return &addResponse{
		ID: id,
	}
}

func (h *HTTP) handleUpdateStatusName(w http.ResponseWriter, r *http.Request) {
	resp := h.updateStatusNameResponse(r)

	resp.writeJSON(w)
}

// swagger:route PUT /api/v1/tasks/status/name Status statusNameUpdate
// Update status name
//
// security:
//    api-key: []
//
// parameters:
//  + name: status
//    in: body
//    required: true
//    type: status
//
// Returns operation result
// responses:
//    200: okResponse
//    400: errorResponse Bad request
//    401: errorResponse Unauthorized
//    403: errorResponse Forbidden
//    500: errorResponse Internal server error
func (h *HTTP) updateStatusNameResponse(r *http.Request) response {
	if r.Body == http.NoBody {
		return getBadRequestWithMsgResponse("no payload", codeEmptyBody)
	}

	status := &status{}

	err := unmarshalReader(r.Body, status)
	if err != nil {
		return getBadRequestWithMsgResponse(fmt.Sprintf("bad body: %s", err.Error()), codeInvalidBodyStructure)
	}

	if status.ID < 0 {
		return getBadRequestWithMsgResponse("id must be positive integer", codeIDPositive)
	}

	ctx := r.Context()
	sub := ctx.Value(subContextVal)
	if sub == nil {
		return getUnauthorizedResponse(codeUnauthorizedNoSub)
	}

	s := sub.(*claims.Subject)

	code, err := h.service.UpdateStatusName(ctx, &core.Status{
		ID:      status.ID,
		Name:    status.Name,
		OwnerID: s.UID,
	})
	if err != nil {
		if code == core.CodeDBFail {
			h.logger.Errorw("Update status.", "err", err)

			return getInternalServerErrorResponse(code)
		}

		if code == core.CodeNotPermitted {
			return getForbiddenResponse(code)
		}

		return getBadRequestWithMsgResponse(err.Error(), code)
	}

	return newOKResponse(struct{}{})
}

func (h *HTTP) handleUpdateStatusParent(w http.ResponseWriter, r *http.Request) {
	resp := h.updateStatusParentResponse(r)

	resp.writeJSON(w)
}

// swagger:route PUT /api/v1/tasks/status/parent Status statusParentUpdate
// Update status parent(provide 0 parentId if you want to add status to the head)
//
// security:
//    api-key: []
//
// parameters:
//  + name: status
//    in: body
//    required: true
//    type: status
//
// Returns operation result
// responses:
//    200: okResponse
//    400: errorResponse Bad request
//    401: errorResponse Unauthorized
//    403: errorResponse Forbidden
//    500: errorResponse Internal server error
func (h *HTTP) updateStatusParentResponse(r *http.Request) response {
	if r.Body == http.NoBody {
		return getBadRequestWithMsgResponse("no payload", codeEmptyBody)
	}

	status := &status{}

	err := unmarshalReader(r.Body, status)
	if err != nil {
		return getBadRequestWithMsgResponse(fmt.Sprintf("bad body: %s", err.Error()), codeInvalidBodyStructure)
	}

	if status.ID < 0 {
		return getBadRequestWithMsgResponse("id must be positive integer", codeIDPositive)
	}

	ctx := r.Context()
	sub := ctx.Value(subContextVal)
	if sub == nil {
		return getUnauthorizedResponse(codeUnauthorizedNoSub)
	}

	s := sub.(*claims.Subject)

	code, err := h.service.UpdateStatusParent(ctx, &core.Status{
		ID:       status.ID,
		ParentID: int(status.ParentID),
		OwnerID:  s.UID,
	})
	if err != nil {
		if code == core.CodeDBFail {
			h.logger.Errorw("Update status.", "err", err)

			return getInternalServerErrorResponse(code)
		}

		if code == core.CodeNotPermitted {
			return getForbiddenResponse(code)
		}

		return getBadRequestWithMsgResponse(err.Error(), code)
	}

	return newOKResponse(struct{}{})
}

// swagger:route Delete /api/v1/tasks/status/{id} Status statusesDelete
// Delete status
//
// security:
//    api-key: []
//
// parameters:
//  + name: id
//    in: path
//    required: true
//    type: integer
//
// Returns operation result
// responses:
//    200: okResponse
//    401: errorResponse Unauthorized
//    500: errorResponse Internal server error
func (h *HTTP) handleDeleteStatus(w http.ResponseWriter, r *http.Request) {
	resp := h.deleteStatusResponse(r)

	resp.writeJSON(w)
}

func (h *HTTP) deleteStatusResponse(r *http.Request) response {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		return getBadRequestWithMsgResponse("no id", codeNoID)
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		msg := fmt.Sprintf("invalid id: %s", err)

		return getBadRequestWithMsgResponse(msg, codeInvalidIDType)
	}

	ctx := r.Context()
	sub := ctx.Value(subContextVal)
	if sub == nil {
		return getUnauthorizedResponse(codeUnauthorizedNoSub)
	}

	s := sub.(*claims.Subject)
	uid, err := strconv.Atoi(s.UID)
	if err != nil {
		return getUnauthorizedResponse(codeUnauthorizedNoSub)
	}

	code, err := h.service.DeleteStatus(ctx, uid, id)
	if err != nil {
		if code == core.CodeDBFail {
			h.logger.Errorw("Delete status.", "err", err)

			return getInternalServerErrorResponse(core.CodeDBFail)
		}

		if code == core.CodeNotPermitted {
			return getForbiddenResponse(code)
		}

		return getBadRequestWithMsgResponse(err.Error(), code)
	}

	return newOKResponse(struct{}{})
}
