package server

import (
	"net/http"

	"github.com/daniilty/kanban-tt/auth/internal/core"
)

const (
	codeEmptyBody             core.Code = "EMPTY_BODY"
	codeInvalidBody           core.Code = "INVALID_BODY_STRUCTURE"
	CodeEmailCannotBeEmpty    core.Code = "EMAIL_CANNOT_BE_EMPTY"
	CodePasswordCannotBeEmpty core.Code = "PASSWORD_CANNOT_BE_EMPTY"
)

// swagger:model
type loginRequest struct {
	// required: true
	Email string `json:"email"`
	// required: true
	Password string `json:"password"`
}

func (l *loginRequest) validate() core.Code {
	if l.Email == "" {
		return CodeEmailCannotBeEmpty
	}

	if l.Password == "" {
		return CodeEmailCannotBeEmpty
	}

	return core.CodeOK
}

// swagger:route POST /api/v1/auth/login Authorize user
// Login to your account
//
// parameters:
//   - name: loginRequest
//     in: body
//     required: true
//     type: loginRequest
//
// Returns operation result
// responses:
//
//	200: accessTokenResponse
//	400: errorResponse Bad request
//	500: errorResponse Internal server error
func (h *HTTP) login(w http.ResponseWriter, r *http.Request) {
	resp := h.getLoginResponse(r)

	resp.writeJSON(w)
}

func (h *HTTP) getLoginResponse(r *http.Request) response {
	if r.Body == http.NoBody {
		return getBadRequestWithMsgResponse("no body", codeEmptyBody)
	}

	l := &loginRequest{}

	err := unmarshalReader(r.Body, l)
	if err != nil {
		return getBadRequestWithMsgResponse(err.Error(), codeInvalidBody)
	}

	code := l.validate()
	if code != core.CodeOK {
		return getBadRequestWithMsgResponse(err.Error(), code)
	}

	accessToken, code, err := h.service.Login(r.Context(), l.toService())
	if err != nil {
		if code == core.CodeInternal {
			h.log.Errorw("Login user.", "err", err)

			return getInternalServerErrorResponse(code)
		}

		return getBadRequestWithMsgResponse(err.Error(), code)
	}

	return &accessTokenResponse{
		AccessToken: accessToken,
	}
}
