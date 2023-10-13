package postgresql

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"
)

type Storage struct {
	db *sql.DB
}

type config interface {
	GetStorageConnectionString() string
}

func New(cfg config) (*Storage, error) {
	db, err := sql.Open("pgx", cfg.GetStorageConnectionString())
	if err != nil {
		return nil, err
	}

	storage := &Storage{
		db: db,
	}

	if err := storage.runMigrations(); err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *Storage) runMigrations() error {
	return goose.Run("up", s.db, "deployments/db/migrations")
}
