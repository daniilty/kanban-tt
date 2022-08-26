package server

import (
	"context"
	"net/http"

	"github.com/daniilty/kanban-tt/auth/claims"
)

type sub string

const (
	subContextVal sub = "sub"
)

func withParseClaimsMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sub, err := claims.ParseHTTPHeader(r.Header)
		if err != nil || sub == nil {
			resp := getUnauthorizedResponse(codeUnauthorizedNoSub)
			resp.writeJSON(w)

			return
		}

		ctx := context.WithValue(r.Context(), subContextVal, sub)
		r = r.WithContext(ctx)

		h(w, r)
	}
}
