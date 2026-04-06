package db

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	_ "github.com/lib/pq"
)

// Config holds database configuration
type Config struct {
	URL      string
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Database wraps the SQL database connection
type Database struct {
	conn *sql.DB
}

// NewDatabase creates a new database connection
func NewDatabase(cfg Config) (*Database, error) {
	if strings.TrimSpace(cfg.URL) != "" {
		conn, err := sql.Open("postgres", cfg.URL)
		if err == nil {
			if pingErr := conn.Ping(); pingErr == nil {
				return &Database{conn: conn}, nil
			}
			_ = conn.Close()
		}

		// Some hosted Postgres URLs include channel_binding which lib/pq may reject.
		sanitizedURL := removeURLParams(cfg.URL, "channel_binding")
		if sanitizedURL != cfg.URL {
			fallbackConn, fallbackErr := sql.Open("postgres", sanitizedURL)
			if fallbackErr == nil {
				if pingErr := fallbackConn.Ping(); pingErr == nil {
					return &Database{conn: fallbackConn}, nil
				}
				_ = fallbackConn.Close()
			}
		}

		return nil, fmt.Errorf("failed to connect using DATABASE_URL")
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := conn.Ping(); err != nil {
		return nil, err
	}

	return &Database{conn: conn}, nil
}

func removeURLParams(raw string, params ...string) string {
	parsed, err := url.Parse(raw)
	if err != nil {
		return raw
	}

	query := parsed.Query()
	for _, param := range params {
		query.Del(param)
	}
	parsed.RawQuery = query.Encode()
	return parsed.String()
}

// Migrate creates required tables and compatibility columns when missing.
func (d *Database) Migrate() error {
	if d.conn == nil {
		return fmt.Errorf("database connection is nil")
	}

	statements := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			email TEXT NOT NULL UNIQUE,
			username TEXT NOT NULL,
			password_hash TEXT,
			password TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS password_hash TEXT`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS password TEXT`,
		`UPDATE users SET password_hash = COALESCE(password_hash, password) WHERE password_hash IS NULL`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users (email)`,

		`CREATE TABLE IF NOT EXISTS candles (
			id BIGSERIAL PRIMARY KEY,
			symbol TEXT NOT NULL,
			open DOUBLE PRECISION NOT NULL,
			high DOUBLE PRECISION NOT NULL,
			low DOUBLE PRECISION NOT NULL,
			close DOUBLE PRECISION NOT NULL,
			volume INTEGER NOT NULL DEFAULT 0,
			timestamp TIMESTAMPTZ NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_candles_symbol_timestamp ON candles (symbol, timestamp DESC)`,

		`CREATE TABLE IF NOT EXISTS alerts (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			symbol TEXT NOT NULL,
			price DOUBLE PRECISION NOT NULL,
			condition TEXT NOT NULL,
			is_active BOOLEAN NOT NULL DEFAULT TRUE,
			triggered_at TIMESTAMPTZ,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_alerts_user_id ON alerts (user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_alerts_active_symbol ON alerts (is_active, symbol)`,

		`CREATE TABLE IF NOT EXISTS transactions (
			id BIGSERIAL PRIMARY KEY,
			user_id TEXT NOT NULL,
			symbol TEXT NOT NULL,
			type TEXT NOT NULL,
			quantity DOUBLE PRECISION NOT NULL,
			price DOUBLE PRECISION NOT NULL,
			amount DOUBLE PRECISION NOT NULL,
			fee DOUBLE PRECISION NOT NULL DEFAULT 0,
			timestamp TIMESTAMPTZ NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_user_timestamp ON transactions (user_id, timestamp DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_symbol_timestamp ON transactions (symbol, timestamp DESC)`,
	}

	for _, stmt := range statements {
		if _, err := d.conn.Exec(stmt); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	return nil
}

// GetConn returns the underlying SQL connection
func (d *Database) GetConn() *sql.DB {
	return d.conn
}

// Close closes the database connection
func (d *Database) Close() error {
	if d.conn != nil {
		return d.conn.Close()
	}
	return nil
}
