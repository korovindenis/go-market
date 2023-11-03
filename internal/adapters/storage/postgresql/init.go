package postgresql

import (
	"database/sql"

	"github.com/pressly/goose"
)

type config interface {
	GetStorageConnectionString() string
}

func Init(config config) (*sql.DB, error) {
	db, err := sql.Open("pgx", config.GetStorageConnectionString())
	if err != nil {
		return nil, err
	}
	if err := runMigrations(db); err != nil {
		return nil, err
	}

	return db, nil
}

func runMigrations(db *sql.DB) error {
	return goose.Run("up", db, "deployments/db/migrations")
}
