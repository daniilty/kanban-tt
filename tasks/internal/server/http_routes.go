package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (h *HTTP) setRoutes(r *mux.Router) {
	api := r.PathPrefix("/api/v1/tasks").Subrouter()

	api.HandleFunc("/task", nest(h.handleAddTask, parseClaimsMiddleware)).Methods(http.MethodPost)
	api.HandleFunc("/tasks", nest(h.handleGetTasks, parseClaimsMiddleware)).Methods(http.MethodGet)
	api.HandleFunc("/statuses", nest(h.handleGetStatuses, parseClaimsMiddleware)).Methods(http.MethodGet)
	api.HandleFunc("/status", nest(h.handleAddStatus, parseClaimsMiddleware)).Methods(http.MethodPost)
}
