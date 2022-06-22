package pg

import "errors"

var (
	ErrNoStatuses = errors.New("no statuses")
	ErrEmptyModel = errors.New("empty model")
)
