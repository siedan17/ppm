package database

import (
	"database/sql"
	"fmt"
	"io/fs"

	"github.com/pressly/goose/v3"
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

	// Bootstrap: if old schema_migrations exists but goose table doesn't,
	// seed goose_db_version so goose won't re-run migration 1 on existing DBs
	if err := bootstrapGoose(sqlDB); err != nil {
		return nil, fmt.Errorf("bootstrap goose: %w", err)
	}

	// Run goose migrations
	goose.SetBaseFS(migrationsFS)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return nil, fmt.Errorf("set goose dialect: %w", err)
	}
	if err := goose.Up(sqlDB, "."); err != nil {
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	return &DB{sqlDB}, nil
}

// bootstrapGoose checks if the old custom migration runner's schema_migrations
// table exists. If so and goose_db_version doesn't exist yet, it creates and
// seeds the goose version table so goose treats migration 1 as already applied.
func bootstrapGoose(db *sql.DB) error {
	var oldExists int
	err := db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='schema_migrations'`).Scan(&oldExists)
	if err != nil || oldExists == 0 {
		return nil // fresh database or no old migrations table
	}

	var gooseExists int
	err = db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='goose_db_version'`).Scan(&gooseExists)
	if err != nil {
		return err
	}
	if gooseExists > 0 {
		return nil // already bootstrapped
	}

	// Old migrations table exists but goose table doesn't — bootstrap
	_, err = db.Exec(`CREATE TABLE goose_db_version (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		version_id INTEGER NOT NULL,
		is_applied INTEGER NOT NULL,
		tstamp TEXT NOT NULL DEFAULT (datetime('now'))
	)`)
	if err != nil {
		return fmt.Errorf("create goose_db_version: %w", err)
	}

	// Mark version 0 (base) and version 1 (initial migration) as applied
	_, err = db.Exec(`INSERT INTO goose_db_version (version_id, is_applied) VALUES (0, 1)`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO goose_db_version (version_id, is_applied) VALUES (1, 1)`)
	if err != nil {
		return err
	}

	fmt.Println("Bootstrapped goose_db_version from existing schema_migrations")
	return nil
}
