package server

import (
	"context"
	"net/http"

	"go.uber.org/zap"
)

type loggerMiddleware struct {
	log *zap.SugaredLogger

	handler http.Handler
}

func (l *loggerMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), "log", l.log)
	r = r.WithContext(ctx)

	l.handler.ServeHTTP(w, r)
}

func (s *HTTP) loggerMiddleware(h http.Handler) http.Handler {
	return &loggerMiddleware{
		log:     s.log,
		handler: h,
	}
}
