package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (h *HTTP) setRoutes(r *mux.Router) {
	api := r.PathPrefix("/api/v1/tasks").Subrouter()

	api.HandleFunc("/task", nest(h.handleAddTask, withParseClaimsMiddleware)).Methods(http.MethodPost)
	api.HandleFunc("/task", nest(h.handleUpdateTask, withParseClaimsMiddleware)).Methods(http.MethodPut)
	api.HandleFunc("/task/{id}", nest(h.handleDeleteTask, withParseClaimsMiddleware)).Methods(http.MethodDelete)
	api.HandleFunc("/tasks", nest(h.handleGetTasks, withParseClaimsMiddleware)).Methods(http.MethodGet)

	api.HandleFunc("/status", nest(h.handleAddStatus, withParseClaimsMiddleware)).Methods(http.MethodPost)
	api.HandleFunc("/status/name", nest(h.handleUpdateStatusName, withParseClaimsMiddleware)).Methods(http.MethodPut)
	api.HandleFunc("/status/parent", nest(h.handleUpdateStatusParent, withParseClaimsMiddleware)).Methods(http.MethodPut)
	api.HandleFunc("/status/{id}", nest(h.handleDeleteStatus, withParseClaimsMiddleware)).Methods(http.MethodDelete)
	api.HandleFunc("/statuses", nest(h.handleGetStatuses, withParseClaimsMiddleware)).Methods(http.MethodGet)
}
