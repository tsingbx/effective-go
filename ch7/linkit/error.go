package linkit

import (
	"errors"
)

var (
	ErrExists    = errors.New("already exists")
	ErrNotExists = errors.New("does not exist")
	ErrInternal  = errors.New("internal error: please try agian")
)
