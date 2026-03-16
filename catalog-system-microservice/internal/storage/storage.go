package storage

import "errors"

var (
	ErrUserExists  = errors.New("user already exists")
	ErrNotFound    = errors.New("not found")
	ErrAppNotFound = errors.New("app not found")
)
