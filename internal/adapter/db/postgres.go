package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type Database struct {
	conn *sql.DB
}

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewDatabase(cfg Config) (*Database, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(5 * time.Minute)

	if err := createSchema(conn); err != nil {
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return &Database{conn: conn}, nil
}

func (d *Database) GetConn() *sql.DB {
	return d.conn
}

func (d *Database) Close() error {
	return d.conn.Close()
}

func createSchema(conn *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(255) PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		username VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

	CREATE TABLE IF NOT EXISTS candles (
		id SERIAL PRIMARY KEY,
		symbol VARCHAR(50) NOT NULL,
		open FLOAT NOT NULL,
		high FLOAT NOT NULL,
		low FLOAT NOT NULL,
		close FLOAT NOT NULL,
		volume FLOAT DEFAULT 0,
		timestamp TIMESTAMP NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_candles_symbol_timestamp ON candles(symbol, timestamp DESC);

	CREATE TABLE IF NOT EXISTS alerts (
		id VARCHAR(255) PRIMARY KEY,
		user_id VARCHAR(255) NOT NULL REFERENCES users(id),
		symbol VARCHAR(50) NOT NULL,
		price FLOAT NOT NULL,
		condition VARCHAR(20) NOT NULL,
		is_active BOOLEAN DEFAULT true,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_alerts_user_id ON alerts(user_id);
	CREATE INDEX IF NOT EXISTS idx_alerts_active ON alerts(is_active);

	CREATE TABLE IF NOT EXISTS transactions (
		id SERIAL PRIMARY KEY,
		user_id VARCHAR(255) NOT NULL REFERENCES users(id),
		symbol VARCHAR(50),
		type VARCHAR(50) NOT NULL,
		quantity FLOAT,
		price FLOAT,
		amount FLOAT NOT NULL,
		fee FLOAT DEFAULT 0,
		timestamp TIMESTAMP NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
	CREATE INDEX IF NOT EXISTS idx_transactions_symbol ON transactions(symbol);
	`

	_, err := conn.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	return nil
}
