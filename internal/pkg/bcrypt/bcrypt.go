package bcrypt

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	appErrors "sso_3.0/internal/errors"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}
func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	if err != nil {
		if errors.Is(bcrypt.ErrMismatchedHashAndPassword, err) {
			return appErrors.ErrPasswordIncorrect
		}
		return err
	}

	return nil
}
