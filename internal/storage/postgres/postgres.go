package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
	"sso_3.0/cmd/migrations"
	configParser "sso_3.0/internal/config"
	"sso_3.0/internal/storage/postgres/task"
	"sso_3.0/internal/storage/postgres/user"
)

type Storage struct {
	TaskStorage *task.Storage
	UserStorage *user.Storage
}

func New(cfg *configParser.Config, log *slog.Logger) (*Storage, error) {
	db, err := sql.Open("postgres", cfg.DbUrl)
	fmt.Println(err)
	if err != nil {
		panic("Error on connecting to database")
	}

	log.Info("Successfully connected to db")

	err = Migrate(cfg.DbUrl)

	if err != nil {
		return nil, err
	}

	taskStorage := task.New(db, log)
	userStorage := user.New(db, log)

	return &Storage{taskStorage, userStorage}, nil
}

func Migrate(dbUrl string) error {
	err := migrations.MigrateDb(dbUrl, "/app/migrations", "up")

	if err != nil {
		return err
	}
	return nil
}
