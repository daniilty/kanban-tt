package server

import (
	"context"
	"net/http"

	"github.com/daniilty/kanban-tt/auth/claims"
)

const (
	subContextVal = "sub"
)

func parseClaimsMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sub, err := claims.ParseHTTPHeader(r.Header)
		if err != nil || sub == nil {
			resp := getUnauthorizedWithResponse()
			resp.writeJSON(w)

			return
		}

		ctx := context.WithValue(r.Context(), subContextVal, sub)
		r = r.WithContext(ctx)

		h(w, r)
	}
}
