package custom_errors

import "errors"

var (
	ErrNotFound          = errors.New("entity not found")
	ErrIncorrectPassword = errors.New("incorrect password")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type ApiError struct {
	Message string `json:"message"`
}
