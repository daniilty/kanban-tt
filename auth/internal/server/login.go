package server

import (
	"fmt"
	"net/http"
)

// swagger:model
type loginRequest struct {
	// required: true
	Email string `json:"email"`
	// required: true
	Password string `json:"password"`
}

func (l *loginRequest) validate() error {
	if l.Email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	if l.Password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	return nil
}

// swagger:route POST /api/v1/auth/login Authorize user
// Login to your account
//
// parameters:
//  + name: loginRequest
//    in: body
//    required: true
//    type: loginRequest
//
// Returns operation result
// responses:
//    200: accessTokenResponse
//    400: errorResponse Bad request
//    500: errorResponse Internal server error
func (h *HTTP) login(w http.ResponseWriter, r *http.Request) {
	resp := h.getLoginResponse(r)

	resp.writeJSON(w)
}

func (h *HTTP) getLoginResponse(r *http.Request) response {
	if r.Body == http.NoBody {
		return getBadRequestWithMsgResponse("no body")
	}

	l := &loginRequest{}

	err := unmarshalReader(r.Body, l)
	if err != nil {
		return getBadRequestWithMsgResponse(err.Error())
	}

	err = l.validate()
	if err != nil {
		return getBadRequestWithMsgResponse(err.Error())
	}

	accessToken, ok, err := h.service.Login(r.Context(), l.toService())
	if err != nil {
		if ok {
			return getBadRequestWithMsgResponse(err.Error())
		}

		h.logger.Errorw("Login user.", "err", err)

		return getInternalServerErrorResponse()
	}

	return &accessTokenResponse{
		AccessToken: accessToken,
	}
}
