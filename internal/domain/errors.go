package domain

import (
	"errors"
)

// Domain errors
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidInput      = errors.New("invalid input")
	ErrInternalServer    = errors.New("internal server error")
)
