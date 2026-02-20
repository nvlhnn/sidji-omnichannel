package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/sidji-omnichannel/internal/config"
)

func NewPostgres(cfg *config.DatabaseConfig) (*sql.DB, error) {
	var dsn string

	// Use DATABASE_URL if provided (e.g. Neon, Railway, Supabase)
	if cfg.URL != "" {
		dsn = cfg.URL
	} else {
		dsn = fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
		)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings (conservative for external DB / pgbouncer)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	return db, nil
}
