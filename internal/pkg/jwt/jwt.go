package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"os"
	"sso_3.0/internal/domain/user"
	appErrors "sso_3.0/internal/errors"
	"time"
)

func NewToken(user *user.Model) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["uid"] = user.Id
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 48).Unix()

	tokenString, err := token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
func CheckToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	var uid = claims["uid"]
	if err != nil {
		return "", err
	}

	if !token.Valid || !ok {
		return "", appErrors.InvalidToken
	}

	return uid.(string), nil
}
