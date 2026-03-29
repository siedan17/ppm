package database

import (
	"database/sql"
	"fmt"
	"io/fs"

	_ "modernc.org/sqlite"
)

type DB struct {
	*sql.DB
}

func Open(path string, migrationsFS fs.FS) (*DB, error) {
	sqlDB, err := sql.Open("sqlite", path+"?_pragma=busy_timeout(5000)")
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	// Enable WAL and foreign keys
	for _, pragma := range []string{
		"PRAGMA journal_mode=WAL",
		"PRAGMA foreign_keys=ON",
	} {
		if _, err := sqlDB.Exec(pragma); err != nil {
			return nil, fmt.Errorf("exec %s: %w", pragma, err)
		}
	}

	db := &DB{sqlDB}

	if migrationsFS != nil {
		if err := runMigrations(db, migrationsFS); err != nil {
			return nil, fmt.Errorf("run migrations: %w", err)
		}
	}

	return db, nil
}
