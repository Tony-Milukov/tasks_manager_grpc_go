package appErrors

import "errors"

var (
	ErrUserExists         = errors.New("user with that Email already exists")
	ErrUserNotExists      = errors.New("user with that do not exists")
	ErrInvalidCredentials = errors.New("invalid Credentials")
	ErrPasswordIncorrect  = errors.New("password Is Incorrect")
	InvalidToken          = errors.New("token is Incorrect")
	NoTokenSent           = errors.New("token was not defined in metadata")
	NothingToDelete       = errors.New("nothing To delete")
	ErrStatusUndefined    = errors.New("status with that id was not defined")
	Internal              = errors.New("internal Server Error")
)
