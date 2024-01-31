package appErrors

import "errors"

var (
	ErrUserExists         = errors.New("User with that Email already exists")
	ErrUserNotExists      = errors.New("User with that do not exists")
	ErrInvalidCredentials = errors.New("Invalid Credentials")
	ErrPasswordIncorrect  = errors.New("Password Is Incorrect")
	InvalidToken          = errors.New("Token is Incorrect")
)
