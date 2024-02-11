package appErrors

import "errors"

var (
	ErrUserExists         = errors.New("user with that Email already exists")
	ErrUserNotExists      = errors.New("user with that id do not exists")
	ErrTaskNotExists      = errors.New("task with that id do not exists")
	ErrInvalidCredentials = errors.New("invalid Credentials")
	ErrPasswordIncorrect  = errors.New("password Is Incorrect")
	InvalidToken          = errors.New("token is Incorrect")
	NoTokenSent           = errors.New("token was not defined in metadata")
	NothingToDelete       = errors.New("nothing to delete")
	ErrStatusUndefined    = errors.New("status with that id was not defined")
	NoArguments           = errors.New("there are not enough arguments to continue")
	TaskAlreadyAssigned   = errors.New("task was already assigned to the user before")
	ErrNoPermission       = errors.New("you have no permission to do that")
	Internal              = errors.New("internal Server Error")
	TaskNotAssigned       = errors.New("this task was not assigned to this user")
)
