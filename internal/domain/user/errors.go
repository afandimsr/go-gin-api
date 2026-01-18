package user

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrPasswordWeak = errors.New("password weak")
)
