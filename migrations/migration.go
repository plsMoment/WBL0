package migrations

import (
	"WBL0/config"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func MigrateUp(cfg *config.Config) error {
	connStr := fmt.Sprintf(
		"pgx5://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.SSLMode,
	)
	m, err := migrate.New(cfg.MigrationPath, connStr)
	if err != nil {
		return err
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}
