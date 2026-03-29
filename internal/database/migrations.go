package database

import (
	"fmt"
	"io/fs"
	"sort"
	"strings"
)

func runMigrations(db *DB, migrationsFS fs.FS) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version TEXT PRIMARY KEY,
		applied_at TEXT NOT NULL DEFAULT (datetime('now'))
	)`)
	if err != nil {
		return fmt.Errorf("create migrations table: %w", err)
	}

	entries, err := fs.ReadDir(migrationsFS, ".")
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	var names []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)

	for _, name := range names {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = ?", name).Scan(&count)
		if err != nil {
			return fmt.Errorf("check migration %s: %w", name, err)
		}
		if count > 0 {
			continue
		}

		content, err := fs.ReadFile(migrationsFS, name)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}

		statements := splitStatements(string(content))
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			if _, err := db.Exec(stmt); err != nil {
				return fmt.Errorf("run migration %s: %w\nstatement: %s", name, err, stmt)
			}
		}

		if _, err := db.Exec("INSERT INTO schema_migrations (version) VALUES (?)", name); err != nil {
			return fmt.Errorf("record migration %s: %w", name, err)
		}
		fmt.Printf("Applied migration: %s\n", name)
	}

	return nil
}

func splitStatements(sql string) []string {
	var statements []string
	var current strings.Builder
	inTrigger := false

	for _, line := range strings.Split(sql, "\n") {
		trimmed := strings.TrimSpace(strings.ToUpper(line))

		if strings.HasPrefix(trimmed, "CREATE TRIGGER") {
			inTrigger = true
		}

		current.WriteString(line)
		current.WriteString("\n")

		if inTrigger {
			if strings.HasSuffix(trimmed, "END;") {
				statements = append(statements, current.String())
				current.Reset()
				inTrigger = false
			}
		} else if strings.HasSuffix(strings.TrimSpace(line), ";") {
			statements = append(statements, current.String())
			current.Reset()
		}
	}

	if s := strings.TrimSpace(current.String()); s != "" {
		statements = append(statements, s)
	}
	return statements
}
