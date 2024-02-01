package user

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"log/slog"
	"sso_3.0/internal/domain/user"
	appErrors "sso_3.0/internal/errors"
)

type StorageInterFace interface {
	Register(username string) *user.Model
	GetUserByEmail(ctx context.Context, email string) (*user.Model, error)
}

type Storage struct {
	db  *sql.DB
	log *slog.Logger
	StorageInterFace
}

func New(db *sql.DB, logger *slog.Logger) *Storage {
	return &Storage{db: db, log: logger}
}

func (s *Storage) Register(ctx context.Context, email, hash string) (*user.Model, error) {
	op := "storage.auth.Register"
	log := s.log.With("op", op)

	var userId int64
	err := s.db.QueryRowContext(ctx, "INSERT INTO users (email,password) VALUES ($1, $2) RETURNING id", email, hash).Scan(&userId)
	log.Info("Created new User")

	pqErr, _ := err.(*pq.Error)

	if err != nil {
		if pqErr.Code == "23505" {
			return nil, appErrors.ErrUserExists
		}

		log.Error("Error on creating User", "errors", err)
		return nil, err
	}

	return &user.Model{
		Id:    userId,
		Email: email,
	}, nil
}

func (s *Storage) GetUserByEmail(ctx context.Context, email string) (*user.Model, error) {
	op := "storage.auth.Register"
	log := s.log.With("op", op)

	var userId int64
	var hash string
	err := s.db.QueryRowContext(ctx, "SELECT id,password FROM users WHERE email=$1", email).Scan(&userId, &hash)

	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, appErrors.ErrUserNotExists
		}

		log.Error("Error on getting auth", "errors", err)
		return nil, err
	}

	return &user.Model{
		Id:    userId,
		Email: email,
		Hash:  hash,
	}, nil
}
