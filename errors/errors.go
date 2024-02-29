package errors

import "errors"

var (
	ErrUserNotFound     = errors.New("user credentials not valid")
	ErrUserAlreadyExist = errors.New("user already exist")
)
