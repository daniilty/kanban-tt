package server

import (
	"fmt"
	"net/http"

	"github.com/daniilty/kanban-tt/auth/internal/validate"
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

func (r *registerRequest) validate() error {
	if r.Email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	err := validate.Email(r.Email)
	if err != nil {
		return err
	}

	if r.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if r.Password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	err = validate.Password(r.Password, 8, true)

	return err
}

// swagger:route POST /api/v1/auth/register Register user
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
		return getBadRequestWithMsgResponse("no body")
	}

	req := &registerRequest{}

	err := unmarshalReader(r.Body, req)
	if err != nil {
		return getBadRequestWithMsgResponse(err.Error())
	}

	err = req.validate()
	if err != nil {
		return getBadRequestWithMsgResponse(err.Error())
	}

	accessToken, ok, err := h.service.Register(r.Context(), req.toService())
	if err != nil {
		if ok {
			return getBadRequestWithMsgResponse(err.Error())
		}

		h.logger.Errorw("Register user.", "err", err)

		return getInternalServerErrorResponse()
	}

	return &accessTokenResponse{
		AccessToken: accessToken,
	}
}
