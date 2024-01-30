package user

import (
	"database/sql"
	"log/slog"
	"sso_3.0/internal/domain/user"
)

type StorageInterFace interface {
	CreateUser(username string) *user.Model
}

type Storage struct {
	db  *sql.DB
	log *slog.Logger
	StorageInterFace
}

func New(db *sql.DB, log *slog.Logger) *Storage {
	return &Storage{db: db}
}

//func (s *Storage) CreateUser(username string) *user.Model {
//
//}
