package server

import (
	"net/http"

	"github.com/daniilty/kanban-tt/auth/internal/core"
)

// swagger:model
type errorResponse struct {
	Status    int       `json:"-"`
	Code      core.Code `json:"code"`
	ErrorInfo string    `json:"errorInfo"`
}

func (e errorResponse) writeJSON(w http.ResponseWriter) error {
	return writeJSONResponse(w, e.Status, e)
}

func getBadRequestWithMsgResponse(msg string, code core.Code) errorResponse {
	return errorResponse{
		Status:    http.StatusBadRequest,
		Code:      code,
		ErrorInfo: msg,
	}
}

func getInternalServerErrorResponse(code core.Code) errorResponse {
	return getInternalServerErrorWithMsgResponse(http.StatusText(http.StatusInternalServerError), code)
}

func getInternalServerErrorWithMsgResponse(msg string, code core.Code) errorResponse {
	return errorResponse{
		Status:    http.StatusInternalServerError,
		Code:      code,
		ErrorInfo: msg,
	}
}

func getUnauthorizedErrorWithMsgResponse(msg string, code core.Code) errorResponse {
	return errorResponse{
		Status:    http.StatusUnauthorized,
		ErrorInfo: msg,
	}
}
