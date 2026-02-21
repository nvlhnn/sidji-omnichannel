package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// RunMigrations automatically applies any pending .up.sql migrations.
// It tracks applied migrations in a `schema_migrations` table.
func RunMigrations(db *sql.DB, migrationsDir string) error {
	// 1. Create tracking table if it doesn't exist
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	// 2. Get already-applied migrations
	applied := make(map[string]bool)
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return fmt.Errorf("failed to query schema_migrations: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return err
		}
		applied[version] = true
	}

	// 3. Find all .up.sql files
	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.up.sql"))
	if err != nil {
		return fmt.Errorf("failed to glob migrations: %w", err)
	}
	sort.Strings(files) // Ensures numerical order (001, 002, ...)

	// 4. Apply pending migrations
	pending := 0
	for _, file := range files {
		name := filepath.Base(file)
		version := strings.TrimSuffix(name, ".up.sql")

		if applied[version] {
			continue
		}

		// Read migration SQL
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", name, err)
		}

		sqlStr := string(content)
		if strings.TrimSpace(sqlStr) == "" {
			continue
		}

		// Execute migration in a transaction
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin tx for %s: %w", name, err)
		}

		if _, err := tx.Exec(sqlStr); err != nil {
			tx.Rollback()
			return fmt.Errorf("migration %s failed: %w", name, err)
		}

		// Record it
		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", name, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", name, err)
		}

		log.Printf("  ✅ Applied migration: %s", name)
		pending++
	}

	if pending == 0 {
		log.Println("  ✅ Database is up to date (no pending migrations)")
	} else {
		log.Printf("  ✅ Applied %d migration(s)", pending)
	}

	return nil
}
