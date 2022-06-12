package server

import (
	"net/http"
	"strings"
)

type accessTokenResponse struct {
	AccessToken string `json:"accessToken"`
}

func (a *accessTokenResponse) writeJSON(w http.ResponseWriter) error {
	return writeJSONResponse(w, http.StatusOK, &a)
}

func parseTokenHeader(h http.Header) string {
	const (
		bearer        = "Bearer "
		authorization = "Authorization"
	)

	raw := h.Get(authorization)

	return strings.TrimPrefix(raw, bearer)
}
