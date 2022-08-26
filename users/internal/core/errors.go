package core

import "errors"

var (
	ErrNoSuchUser = errors.New("no such user")
	ErrNoSuchTTL  = errors.New("no such ttl")
)
