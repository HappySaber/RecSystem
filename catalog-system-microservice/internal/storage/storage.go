package storage

import "errors"

var (
	ErrUserExists   = errors.New("user already exists")
	ErrShowNotFound = errors.New("show not found")
	ErrAppNotFound  = errors.New("app not found")
)
