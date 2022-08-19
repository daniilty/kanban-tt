package server

import (
	"errors"
	"net/http"

	"github.com/daniilty/kanban-tt/auth/internal/core"
	"github.com/gorilla/mux"
)

const (
	codeNoHeader       = "NO_HEADER"
	codeInvalidSubject = "INVALID_SUB"
	codeNoKeyParam     = "NO_KEY_URL_PARAM"
	codeNoSuchKeyParam = "NO_SUCH_KEY"
)

// swagger:model
type userInfoResponse struct {
	ID             string `json:"id"`
	Email          string `json:"email"`
	EmailConfirmed bool   `json:"emailConfirmed"`
	Name           string `json:"name"`
	TaskTTL        int    `json:"taskTTL"`
}

// swagger:model
type userRequest struct {
	Email          string `json:"email"`
	EmailConfirmed bool   `json:"emailConfirmed"`
	Password       string `json:"password"`
	Name           string `json:"name"`
	TaskTTL        int    `json:"taskTTL"`
}

func (u *userInfoResponse) writeJSON(w http.ResponseWriter) error {
	return writeJSONResponse(w, http.StatusOK, u)
}

// swagger:route GET /api/v1/auth/me UserInfo GetUser
// get account info
//
// security:
//    api_key: []
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
		return getUnauthorizedErrorWithMsgResponse("no header", codeNoHeader)
	}

	userInfo, code, err := h.service.GetUserInfo(r.Context(), token)
	if err != nil {
		if code == core.CodeInternal {
			h.logger.Errorw("Get user info.", "err", err)

			return getUnauthorizedErrorWithMsgResponse(err.Error(), code)
		}

		return getInternalServerErrorResponse(code)
	}

	return convertCoreUserInfoToResponse(userInfo)
}

func (h *HTTP) confirmEmail(w http.ResponseWriter, r *http.Request) {
	resp := h.getConfirmEmailResponse(r)

	resp.writeJSON(w)
}

func (h *HTTP) getConfirmEmailResponse(r *http.Request) response {
	vars := mux.Vars(r)
	key, ok := vars["key"]
	if !ok {
		return getBadRequestWithMsgResponse("missing key parameter", codeNoHeader)
	}

	err := h.service.ConfirmUserEmail(r.Context(), key)
	if err != nil {
		if errors.Is(err, core.ErrNoSuchKey) {
			return getBadRequestWithMsgResponse(err.Error(), codeNoSuchKeyParam)
		}

		h.logger.Errorw("Confirm user email.", "err", err)

		return getInternalServerErrorResponse(core.CodeInternal)
	}

	return getOkResponse(struct{}{})
}

// swagger:route PUT /api/v1/auth/me UserInfo updateUser
// Update your account
//
// parameters:
//  + name: userRequest
//    in: body
//    required: true
//    type: userRequest
// security:
//    api_key: []
// Returns operation result
// responses:
//    200: okResp
//    400: errorResponse Bad request
//    401: errorResponse Unauthorized
//    500: errorResponse Internal server error
func (h *HTTP) updateUser(w http.ResponseWriter, r *http.Request) {
	resp := h.getUpdateUserResponse(r)

	resp.writeJSON(w)
}

func (h *HTTP) getUpdateUserResponse(r *http.Request) response {
	token := parseTokenHeader(r.Header)
	if token == "" {
		return getUnauthorizedErrorWithMsgResponse("no header", codeNoHeader)
	}

	sub, err := h.service.ParseRawToken(token)
	if err != nil {
		return getUnauthorizedErrorWithMsgResponse("invalid subject", codeInvalidSubject)
	}

	if r.Body == http.NoBody {
		return getBadRequestWithMsgResponse("no body", codeEmptyBody)
	}

	req := &userRequest{}

	err = unmarshalReader(r.Body, req)
	if err != nil {
		return getBadRequestWithMsgResponse(err.Error(), codeInvalidBody)
	}

	code, err := h.service.UpdateUser(r.Context(), &core.UserInfo{
		ID:       sub.UID,
		Name:     req.Name,
		Email:    req.Email,
		TaskTTL:  req.TaskTTL,
		Password: req.Password,
	})
	if err != nil {
		if code != core.CodeInternal {
			return getBadRequestWithMsgResponse(err.Error(), code)
		}

		h.logger.Errorw("Update user.", "err", err)
	}

	return getOkResponse(struct{}{})
}
