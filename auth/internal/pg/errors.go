package pg

import "errors"

var (
	ErrNoSuchToken = errors.New("such token does not exist")
)
