package authService

import (
	"context"
	"errors"
	"log/slog"
	appErrors "sso_3.0/internal/errors"
	"sso_3.0/internal/pkg/bcrypt"
	"sso_3.0/internal/pkg/jwt"
	"sso_3.0/internal/storage/postgres"
	"sso_3.0/internal/storage/postgres/user"
)

type Service struct {
	log         *slog.Logger
	userStorage *user.Storage
}

func New(log *slog.Logger, storage *postgres.Storage) *Service {
	return &Service{log: log, userStorage: storage.UserStorage}
}

func (s *Service) Register(ctx context.Context, email, password string) (string, error) {
	op := "service.auth.Register"
	log := s.log.With("op", op)

	hash, err := bcrypt.HashPassword(password)
	if err != nil {
		log.Error("Error on Hashing Password")
		return "", err
	}
	user, err := s.userStorage.Register(ctx, email, hash)
	if err != nil {
		return "", err
	}

	token, err := jwt.NewToken(user)

	return token, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	op := "service.auth.Login"
	log := s.log.With("op", op)

	user, err := s.userStorage.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(appErrors.ErrUserNotExists, err) {
			return "", appErrors.ErrInvalidCredentials
		}
		log.Error("Error: ", err)
		return "", err
	}

	err = bcrypt.CheckPasswordHash(password, user.Hash)
	if err != nil {
		if errors.Is(appErrors.ErrPasswordIncorrect, err) {
			return "", appErrors.ErrInvalidCredentials
		}
		log.Error("Error: ", err)
		return "", err
	}

	token, err := jwt.NewToken(user)

	return token, nil
}
