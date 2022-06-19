package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type okResponse struct {
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
