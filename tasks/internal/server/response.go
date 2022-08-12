package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// swagger:model
type addResponse struct {
	// required: true
	ID int `json:"id"`
}

func (a *addResponse) writeJSON(w http.ResponseWriter) error {
	return writeJSONResponse(w, http.StatusOK, a)
}

// swagger:model
type okResponse struct {
	// required: true
	jsonData interface{}
}

func newOKResponse(data interface{}) response {
	return &okResponse{
		jsonData: data,
	}
}

func (o *okResponse) writeJSON(w http.ResponseWriter) error {
	return writeJSONResponse(w, http.StatusOK, o.jsonData)
}

type response interface {
	writeJSON(http.ResponseWriter) error
}

func writeJSONResponse(w http.ResponseWriter, status int, v interface{}) error {
	bb, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}

	w.WriteHeader(status)

	_, err = w.Write(bb)

	return err
}
