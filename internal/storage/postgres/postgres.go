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
	"time"
)

type Storage struct {
	TaskStorage *task.Storage
	UserStorage *user.Storage
}

func New(cfg *configParser.Config, log *slog.Logger) (*Storage, error) {
	db, err := sql.Open("postgres", cfg.DbUrl)

	if err != nil {
		log.Error("Successfully connected to db")
		panic("Error on connecting to database")
	}

	log.Info("Successfully connected to db")

	err = Migrate(cfg.DbUrl, 5)

	taskStorage := task.New(db, log)
	userStorage := user.New(db, log)

	return &Storage{taskStorage, userStorage}, nil
}

func Migrate(dbUrl string, triesCount int) error {
	for try := 1; try < triesCount; try++ {
		err := migrations.MigrateDb(dbUrl, "/app/migrations", "up")

		if err == nil {
			fmt.Println("\nSuccessfully Migrated DB\n")
			break
		}
		fmt.Printf("ERROR on migrate: %e", err)
		for count := 15; count != 0; count-- {
			fmt.Printf("\nReTrying to Migrate db in %d..\n", count)
			time.Sleep(time.Second * 1)
		}
	}

	return nil
}
