package server

import (
	"net/http"
)

type errorResponse struct {
	Status    int    `json:"status"`
	ErrorInfo string `json:"errorInfo"`
}

func (e errorResponse) writeJSON(w http.ResponseWriter) error {
	return writeJSONResponse(w, e.Status, e)
}

func getBadRequestWithMsgResponse(msg string) errorResponse {
	return errorResponse{
		Status:    http.StatusBadRequest,
		ErrorInfo: msg,
	}
}

func getInternalServerErrorResponse() errorResponse {
	return getInternalServerErrorWithMsgResponse(http.StatusText(http.StatusInternalServerError))
}

func getInternalServerErrorWithMsgResponse(msg string) errorResponse {
	return errorResponse{
		Status:    http.StatusInternalServerError,
		ErrorInfo: msg,
	}
}

func getUnauthorizedErrorWithMsgResponse(msg string) errorResponse {
	return errorResponse{
		Status:    http.StatusUnauthorized,
		ErrorInfo: msg,
	}
}
