package migrations

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {

	dbUrl, migrationsPath, op := getFlags()

	err := MigrateDb(dbUrl, migrationsPath, op)

	if err != nil {
		panic(err)
	}

	fmt.Println("Database Successfully migrated")
}

func MigrateDb(url, path, op string) error {
	sourcePath := fmt.Sprintf("file://%s", path)
	m, err := migrate.New(sourcePath, url)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	switch op {
	case "down":
		if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("migration down errors: %w", err)
		}
	default: // Treat any non-"down" operation as "up".
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("migration up errors: %w", err)
		}
	}

	return nil
}

func getFlags() (dbUrl, migrationsPath, op string) {
	flag.StringVar(&dbUrl, "db-url", "", "Your connection url to the postgres  db")
	flag.StringVar(&migrationsPath, "migrations-path", "", "Your connection url to the postgres  db")
	flag.StringVar(&op, "op", "", "Migrations operator: up | down")
	flag.Parse()

	if dbUrl == "" || migrationsPath == "" {
		panic("Flags were not defined: db_url, migrations_path")
	}

	return dbUrl, migrationsPath, op
}
