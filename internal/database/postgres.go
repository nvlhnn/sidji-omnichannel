package database

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
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

	// For Neon/PgBouncer compatibility: disable prepared statements
	// by adding default_query_exec_mode=simple_protocol
	if strings.Contains(dsn, "neon.tech") {
		u, err := url.Parse(dsn)
		if err == nil {
			q := u.Query()
			q.Set("default_query_exec_mode", "simple_protocol")
			u.RawQuery = q.Encode()
			dsn = u.String()
		}
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Connection pool settings (tuned for Neon serverless)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(30 * time.Second)

	return db, nil
}
