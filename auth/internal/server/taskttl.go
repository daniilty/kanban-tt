package server

import (
	"context"
	"net/http"
)

func (h *HTTP) getTaskTTLs(w http.ResponseWriter, r *http.Request) {
	resp := h.getTaskTTLsResponse(r.Context())

	resp.writeJSON(w)
}

func (h *HTTP) getTaskTTLsResponse(ctx context.Context) response {
	return getOkResponse(h.service.GetTTLs(ctx))
}
