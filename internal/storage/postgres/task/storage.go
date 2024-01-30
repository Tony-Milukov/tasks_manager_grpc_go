package task

import (
	"database/sql"
	"log/slog"
)

type Storage struct {
	db  *sql.DB
	log *slog.Logger
}

func New(db *sql.DB, log *slog.Logger) *Storage {

	return &Storage{db: db}
}
