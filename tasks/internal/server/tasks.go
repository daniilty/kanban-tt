package server

import (
	"fmt"
	"net/http"

	"github.com/daniilty/kanban-tt/auth/claims"
	"github.com/daniilty/kanban-tt/tasks/internal/core"
)

type task struct {
	ID       int    `json:"id,omitempty"`
	Content  string `json:"content"`
	Priority uint32 `json:"priority"`
	StatusID uint32 `json:"status_id"`
}

func (t *task) validate() error {
	if t.Content == "" {
		return fmt.Errorf("content cannot be empty")
	}

	return nil
}

func (h *HTTP) handleGetTasks(w http.ResponseWriter, r *http.Request) {
	resp := h.getTasksResponse(r)

	resp.writeJSON(w)
}

func (h *HTTP) getTasksResponse(r *http.Request) response {
	ctx := r.Context()
	sub := ctx.Value(subContextVal)
	if sub == nil {
		return getBadRequestWithMsgResponse("no subject")
	}

	s := sub.(*claims.Subject)

	tasks, err := h.service.GetTasks(ctx, s.UID)
	if err != nil {
		h.logger.Errorw("get tasks", "err", err)

		return getInternalServerErrorResponse()
	}

	return newOKResponse(tasks)
}

func (h *HTTP) handleAddTask(w http.ResponseWriter, r *http.Request) {
	resp := h.addTaskResponse(r)

	resp.writeJSON(w)
}

func (h *HTTP) addTaskResponse(r *http.Request) response {
	if r.Body == http.NoBody {
		return getBadRequestWithMsgResponse("no payload")
	}

	task := &task{}

	err := unmarshalReader(r.Body, task)
	if err != nil {
		return getBadRequestWithMsgResponse(fmt.Sprintf("bad body: %s", err.Error()))
	}

	err = task.validate()
	if err != nil {
		return getBadRequestWithMsgResponse(err.Error())
	}

	ctx := r.Context()
	sub := ctx.Value(subContextVal)
	if sub == nil {
		return getBadRequestWithMsgResponse("no subject")
	}

	s := sub.(*claims.Subject)

	err, ok := h.service.AddTask(ctx, &core.Task{
		Content:  task.Content,
		OwnerID:  s.UID,
		Priority: task.Priority,
		StatusID: task.StatusID,
	})
	if err != nil {
		if ok {
			return getBadRequestWithMsgResponse(err.Error())
		}

		h.logger.Errorw("add task", "err", err)

		return getInternalServerErrorResponse()
	}

	return newOKResponse(struct{}{})
}

func (h *HTTP) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	resp := h.updateTaskResponse(r)

	resp.writeJSON(w)
}

func (h *HTTP) updateTaskResponse(r *http.Request) response {
	if r.Body == http.NoBody {
		return getBadRequestWithMsgResponse("no payload")
	}

	task := &task{}

	err := unmarshalReader(r.Body, task)
	if err != nil {
		return getBadRequestWithMsgResponse(fmt.Sprintf("bad body: %s", err.Error()))
	}

	if task.ID < 0 {
		return getBadRequestWithMsgResponse("id must be positive integer")
	}

	err, ok := h.service.UpdateTask(r.Context(), &core.Task{
		ID:       task.ID,
		Content:  task.Content,
		Priority: task.Priority,
		StatusID: task.StatusID,
	})
	if err != nil {
		if ok {
			return getBadRequestWithMsgResponse(err.Error())
		}

		h.logger.Errorw("update task", "err", err)

		return getInternalServerErrorResponse()
	}

	return newOKResponse(struct{}{})
}
