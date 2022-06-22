package server

import "net/http"

// swagger:model
type userInfoResponse struct {
	ID             string `json:"id"`
	Email          string `json:"email"`
	EmailConfirmed bool   `json:"emailConfirmed"`
	Name           string `json:"name"`
}

func (u *userInfoResponse) writeJSON(w http.ResponseWriter) error {
	return writeJSONResponse(w, http.StatusOK, u)
}

// swagger:route GET /api/v1/auth/me UserInfo user
// get account info
//
// security:
//    api-key: Bearer
// Returns operation result
// responses:
//    200: userInfoResponse
//    400: errorResponse Bad request
//    401: errorResponse Unauthorized
//    500: errorResponse Internal server error
func (h *HTTP) me(w http.ResponseWriter, r *http.Request) {
	resp := h.getMeResponse(r)

	resp.writeJSON(w)
}

func (h *HTTP) getMeResponse(r *http.Request) response {
	token := parseTokenHeader(r.Header)
	if token == "" {
		return getUnauthorizedErrorWithMsgResponse("no header")
	}

	userInfo, ok, err := h.service.GetUserInfo(r.Context(), token)
	if err != nil {
		if ok {
			return getUnauthorizedErrorWithMsgResponse(err.Error())
		}

		h.logger.Errorw("Get user info.", "err", err)

		return getInternalServerErrorResponse()
	}

	return convertCoreUserInfoToResponse(userInfo)
}
