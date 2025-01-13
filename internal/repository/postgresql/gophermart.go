package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/Makovey/gophermart/internal/config"
	"github.com/Makovey/gophermart/internal/logger"
)

const (
	migrationPath = "internal/db/migrations"
)

type Repo struct {
	log  logger.Logger
	pool *pgxpool.Pool
}

func NewPostgresRepo(log logger.Logger, cfg config.Config) (*Repo, error) {
	fn := "postgresql.NewPostgresRepo"

	path, err := filepath.Abs(migrationPath)
	if err != nil {
		return nil, fmt.Errorf("[%s]: could not determine absolute path for migrations: %w", fn, err)
	}

	err = upMigrations(cfg.DatabaseURI(), path)
	if err != nil {
		return nil, fmt.Errorf("[%s]: could not up migrations: %w", fn, err)
	}

	conn, err := pgxpool.New(context.Background(), cfg.DatabaseURI())
	if err != nil {
		return nil, fmt.Errorf("[%s]: could not connect to database: %w", fn, err)
	}

	return &Repo{
		log:  log,
		pool: conn,
	}, nil
}

func upMigrations(databaseURI, migrationsDir string) error {
	db, err := sql.Open("pgx", databaseURI)
	if err != nil {
		return err
	}
	defer db.Close()

	if err = goose.Up(db, migrationsDir); err != nil {
		return err
	}

	return nil
}

func (r *Repo) Close() error {
	r.pool.Close()
	return nil
}
