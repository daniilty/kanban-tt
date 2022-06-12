package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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
