package server

import (
	"fmt"
	"net/http"

	"github.com/daniilty/kanban-tt/auth/internal/core"
	"github.com/daniilty/kanban-tt/auth/internal/validate"
)

const (
	codeEmailCannotBeEmpty    = "EMAIL_CANNOT_BE_EMPTY"
	codeNameCannotBeEmpty     = "NAME_CANNOT_BE_EMPTY"
	codePasswordCannotBeEmpty = "PASSWORD_CANNOT_BE_EMPTY"
	codeInvalidPassword       = "PASSWORD_IS_INVALID"
	codeInvalidEmail          = "EMAIL_IS_INVALID"
)

// swagger:model
type registerRequest struct {
	// required: true
	Email string `json:"email"`
	// required: true
	Name string `json:"name"`
	// required: true
	Password string `json:"password"`
}

func (r *registerRequest) validate() (core.Code, error) {
	if r.Email == "" {
		return CodeEmailCannotBeEmpty, fmt.Errorf("email cannot be empty")
	}

	err := validate.Email(r.Email)
	if err != nil {
		return codeInvalidEmail, err
	}

	if r.Name == "" {
		return codeNameCannotBeEmpty, fmt.Errorf("name cannot be empty")
	}

	if r.Password == "" {
		return codePasswordCannotBeEmpty, fmt.Errorf("password cannot be empty")
	}

	err = validate.Password(r.Password, 8, false)
	if err != nil {
		return codeInvalidPassword, err
	}

	return core.CodeOK, err
}

// swagger:route POST /api/v1/auth/register Register registerUser
// Register user
//
// parameters:
//  + name: registerRequest
//    in: body
//    required: true
//    type: registerRequest
//
// Returns operation result
// responses:
//    200: accessTokenResponse
//    400: errorResponse Bad request
//    500: errorResponse Internal server error
func (h *HTTP) register(w http.ResponseWriter, r *http.Request) {
	resp := h.getRegisterResponse(r)

	resp.writeJSON(w)
}

func (h *HTTP) getRegisterResponse(r *http.Request) response {
	if r.Body == http.NoBody {
		return getBadRequestWithMsgResponse("no body", codeEmptyBody)
	}

	req := &registerRequest{}

	err := unmarshalReader(r.Body, req)
	if err != nil {
		return getBadRequestWithMsgResponse(err.Error(), codeInvalidBody)
	}

	code, err := req.validate()
	if err != nil {
		return getBadRequestWithMsgResponse(err.Error(), code)
	}

	accessToken, code, err := h.service.Register(r.Context(), req.toService())
	if err != nil {
		if code != core.CodeInternal {
			return getBadRequestWithMsgResponse(err.Error(), code)
		}

		h.logger.Errorw("Register user.", "err", err)

		return getInternalServerErrorResponse(code)
	}

	return &accessTokenResponse{
		AccessToken: accessToken,
	}
}
