package model

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrInvalidInput  = errors.New("invalid input parameters")
	ErrAlreadyExists = errors.New("already exists")
	ErrInvalidSort   = errors.New("invalid sort parameters")
)
