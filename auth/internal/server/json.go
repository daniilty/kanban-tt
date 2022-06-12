package server

import (
	"encoding/json"
	"io"
)

func unmarshalReader(r io.Reader, v interface{}) error {
	bb, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	return json.Unmarshal(bb, v)
}
