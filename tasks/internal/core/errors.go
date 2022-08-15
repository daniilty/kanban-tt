package core

import "errors"

var (
	ErrStatusWithNameExists = errors.New("status with such name already exist")
	ErrNoSuchParent         = errors.New("parent with such id does not exist")
)
