package server

import (
	"context"
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/daniilty/kanban-tt/tasks/internal/core"
	"github.com/gorilla/mux"
)

// HTTP - http server.
type HTTP struct {
	innerServer *http.Server

	logger  *zap.SugaredLogger
	service core.Service
}

func (h *HTTP) Run(ctx context.Context) {
	h.logger.Infow("HTTP server starting.", "addr", h.innerServer.Addr)

	go func() {
		err := h.innerServer.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}

			h.logger.Errorw("Listen and serve HTTP", "addr", h.innerServer.Addr, "err", err)
		}
	}()

	<-ctx.Done()

	h.logger.Info("Graceful server shutdown.")
	h.innerServer.Shutdown(context.Background())
}

// NewHTTP - constructor.
func NewHTTP(addr string, logger *zap.SugaredLogger, service core.Service) *HTTP {
	h := &HTTP{
		logger:  logger,
		service: service,
	}

	r := mux.NewRouter()
	h.setRoutes(r)

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	h.innerServer = srv

	return h
}
