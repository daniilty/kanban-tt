package server

import "net/http"

func (h *HTTP) jwks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	_, err := w.Write(h.service.JWKS())
	if err != nil {
		h.logger.Debugw("Write jwks.", "err", err)
	}
}
