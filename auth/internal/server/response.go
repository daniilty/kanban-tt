package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type response interface {
	writeJSON(http.ResponseWriter) error
}

// swagger:model
type okResp struct {
	data interface{}
}

func (o *okResp) writeJSON(w http.ResponseWriter) error {
	return writeJSONResponse(w, http.StatusOK, o.data)
}

func writeJSONResponse(w http.ResponseWriter, status int, v interface{}) error {
	bb, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(bb)

	return err
}

func getOkResponse(data interface{}) response {
	return &okResp{
		data: data,
	}
}
