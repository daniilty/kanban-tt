package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/daniilty/kanban-tt/auth/claims"
	"github.com/daniilty/kanban-tt/tasks/internal/core"
	"github.com/gorilla/mux"
)

const (
	// пустое поле контента
	codeContentEmpty core.Code = "EMPTY_CONTENT"
)

// swagger:model
type task struct {
	// required: true
	ID int `json:"id,omitempty"`
	// required: true
	Content string `json:"content"`
	// required: true
	Priority uint32 `json:"priority"`
	// required: true
	StatusID uint32 `json:"status_id"`
}

func (t *task) validate() (error, core.Code) {
	if t.Content == "" {
		return fmt.Errorf("content cannot be empty"), codeContentEmpty
	}

	return nil, core.CodeOK
}

// swagger:route GET /api/v1/tasks/tasks Task tasksGet
// get user created tasks
//
// security:
//    api-key: []
// Returns operation result
// responses:
//    200: task
//    400: errorResponse Bad request
//    401: errorResponse Unauthorized
//    500: errorResponse Internal server error
func (h *HTTP) handleGetTasks(w http.ResponseWriter, r *http.Request) {
	resp := h.getTasksResponse(r)

	resp.writeJSON(w)
}

func (h *HTTP) getTasksResponse(r *http.Request) response {
	ctx := r.Context()
	sub := ctx.Value(subContextVal)
	if sub == nil {
		return getUnauthorizedWithResponse(codeUnauthorizedNoSub)
	}

	s := sub.(*claims.Subject)

	tasks, err := h.service.GetUserTasks(ctx, s.UID)
	if err != nil {
		h.logger.Errorw("get tasks", "err", err)

		return getInternalServerErrorResponse(core.CodeDBFail)
	}

	return newOKResponse(tasks)
}

// swagger:route POST /api/v1/tasks/task Task taskAdd
// Add task
//
// security:
//    api-key: []
//
// parameters:
//  + name: status
//    in: body
//    required: true
//    type: task
//
// Returns operation result
// responses:
//    200: okResponse
//    400: errorResponse Bad request
//    401: errorResponse Unauthorized
//    500: errorResponse Internal server error
func (h *HTTP) handleAddTask(w http.ResponseWriter, r *http.Request) {
	resp := h.addTaskResponse(r)

	resp.writeJSON(w)
}

func (h *HTTP) addTaskResponse(r *http.Request) response {
	if r.Body == http.NoBody {
		return getBadRequestWithMsgResponse("no payload", codeEmptyBody)
	}

	task := &task{}

	err := unmarshalReader(r.Body, task)
	if err != nil {
		return getBadRequestWithMsgResponse(fmt.Sprintf("bad body: %s", err.Error()), codeInvalidBodyStructure)
	}

	err, code := task.validate()
	if err != nil {
		return getBadRequestWithMsgResponse(err.Error(), code)
	}

	ctx := r.Context()
	sub := ctx.Value(subContextVal)
	if sub == nil {
		return getUnauthorizedWithResponse(codeUnauthorizedNoSub)
	}

	s := sub.(*claims.Subject)

	err, code = h.service.AddTask(ctx, &core.Task{
		Content:  task.Content,
		OwnerID:  s.UID,
		Priority: task.Priority,
		StatusID: task.StatusID,
	})
	if err != nil {
		if code != core.CodeDBFail {
			return getBadRequestWithMsgResponse(err.Error(), code)
		}

		h.logger.Errorw("add task", "err", err)

		return getInternalServerErrorResponse(code)
	}

	return newOKResponse(struct{}{})
}

func (h *HTTP) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	resp := h.updateTaskResponse(r)

	resp.writeJSON(w)
}

// swagger:route PUT /api/v1/tasks/task Task taskUpdate
// Update task
//
// security:
//    api-key: []
//
// parameters:
//  + name: task
//    in: body
//    required: true
//    type: task
//
// Returns operation result
// responses:
//    200: okResponse
//    400: errorResponse Bad request
//    401: errorResponse Unauthorized
//    500: errorResponse Internal server error
func (h *HTTP) updateTaskResponse(r *http.Request) response {
	if r.Body == http.NoBody {
		return getBadRequestWithMsgResponse("no payload", codeEmptyBody)
	}

	task := &task{}

	err := unmarshalReader(r.Body, task)
	if err != nil {
		return getBadRequestWithMsgResponse(fmt.Sprintf("bad body: %s", err.Error()), codeInvalidBodyStructure)
	}

	if task.ID < 0 {
		return getBadRequestWithMsgResponse("id must be positive integer", codeIDPositive)
	}

	err, code := h.service.UpdateTask(r.Context(), &core.Task{
		ID:       task.ID,
		Content:  task.Content,
		Priority: task.Priority,
		StatusID: task.StatusID,
	})
	if err != nil {
		if code != core.CodeDBFail {
			return getBadRequestWithMsgResponse(err.Error(), code)
		}

		h.logger.Errorw("update task", "err", err)

		return getInternalServerErrorResponse(code)
	}

	return newOKResponse(struct{}{})
}

// swagger:route Delete /api/v1/tasks/task/{id} Task taskDelete
// Delete task
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
func (h *HTTP) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	resp := h.deleteTaskResponse(r)

	resp.writeJSON(w)
}

func (h *HTTP) deleteTaskResponse(r *http.Request) response {
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

	err = h.service.DeleteTask(r.Context(), id)
	if err != nil {
		return getInternalServerErrorResponse(core.CodeDBFail)
	}

	return newOKResponse(struct{}{})
}
