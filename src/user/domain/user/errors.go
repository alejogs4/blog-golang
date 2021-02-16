package user

import "errors"

// Error for user behaviors
var (
	ErrBadUserData          = errors.New("User: All fields must be present")
	ErrTooShortUserPassword = errors.New("User: Password should be at least 6 characters long")
	ErrInvalidUserLogin     = errors.New("User: Either email or password are wrong")
)
