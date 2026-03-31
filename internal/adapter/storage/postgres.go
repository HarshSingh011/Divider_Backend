package storage

import (
	"database/sql"
	"errors"
	"time"

	"stocktrack-backend/internal/domain"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) SaveUser(user *domain.User) error {
	if user == nil || user.Email == "" {
		return errors.New("invalid user")
	}

	query := `
		INSERT INTO users (id, email, username, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
			password_hash = EXCLUDED.password_hash,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.Exec(query,
		user.ID,
		user.Email,
		user.Username,
		user.Password,
		time.Now(),
		time.Now(),
	)

	return err
}

func (r *PostgresUserRepository) FindByEmail(email string) (*domain.User, error) {
	query := `SELECT id, email, username, password_hash FROM users WHERE email = $1`

	var user domain.User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) FindByUsername(username string) (*domain.User, error) {
	query := `SELECT id, email, username, password_hash FROM users WHERE username = $1`

	var user domain.User
	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}
func (r *PostgresUserRepository) FindByID(id string) (*domain.User, error) {
	query := `SELECT id, email, username, password FROM users WHERE id = $1`

	var user domain.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) Exists(email string) bool {
	query := `SELECT 1 FROM users WHERE email = $1 LIMIT 1`
	var exists int
	err := r.db.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return false
	}
	return exists == 1
}


type PostgresCandleRepository struct {
	db *sql.DB
}

func NewPostgresCandleRepository(db *sql.DB) *PostgresCandleRepository {
	return &PostgresCandleRepository{db: db}
}

func (r *PostgresCandleRepository) SaveCandle(candle *domain.Candle) error {
	if candle == nil || candle.Symbol == "" {
		return errors.New("invalid candle")
	}

	query := `
		INSERT INTO candles (symbol, open, high, low, close, volume, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(query,
		candle.Symbol,
		candle.Open,
		candle.High,
		candle.Low,
		candle.Close,
		candle.Volume,
		candle.Timestamp,
	)

	return err
}

func (r *PostgresCandleRepository) GetCandles(symbol string, limit int) ([]domain.Candle, error) {
	query := `
		SELECT symbol, open, high, low, close, volume, timestamp
		FROM candles
		WHERE symbol = $1
		ORDER BY timestamp DESC
		LIMIT $2
	`

	rows, err := r.db.Query(query, symbol, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var candles []domain.Candle
	for rows.Next() {
		var candle domain.Candle
		if err := rows.Scan(
			&candle.Symbol,
			&candle.Open,
			&candle.High,
			&candle.Low,
			&candle.Close,
			&candle.Volume,
			&candle.Timestamp,
		); err != nil {
			return nil, err
		}
		candles = append(candles, candle)
	}

	return candles, rows.Err()
}

func (r *PostgresCandleRepository) GetCandlesByTimeRange(symbol string, from, to time.Time) ([]domain.Candle, error) {
	query := `
		SELECT symbol, open, high, low, close, volume, timestamp
		FROM candles
		WHERE symbol = $1 AND timestamp BETWEEN $2 AND $3
		ORDER BY timestamp ASC
	`

	rows, err := r.db.Query(query, symbol, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var candles []domain.Candle
	for rows.Next() {
		var candle domain.Candle
		if err := rows.Scan(
			&candle.Symbol,
			&candle.Open,
			&candle.High,
			&candle.Low,
			&candle.Close,
			&candle.Volume,
			&candle.Timestamp,
		); err != nil {
			return nil, err
		}
		candles = append(candles, candle)
	}

	return candles, rows.Err()
}

type PostgresAlertRepository struct {
	db *sql.DB
}

func NewPostgresAlertRepository(db *sql.DB) *PostgresAlertRepository {
	return &PostgresAlertRepository{db: db}
}

func (r *PostgresAlertRepository) SaveAlert(alert *domain.Alert) error {
	if alert == nil || alert.ID == "" {
		return errors.New("invalid alert")
	}

	query := `
		INSERT INTO alerts (id, user_id, symbol, price, condition, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
			price = EXCLUDED.price,
			condition = EXCLUDED.condition,
			is_active = EXCLUDED.is_active
	`

	_, err := r.db.Exec(query,
		alert.ID,
		alert.UserID,
		alert.Symbol,
		alert.ThresholdPrice,
		alert.Condition,
		alert.IsActive,
	)

	return err
}

func (r *PostgresAlertRepository) FindAlertsByUser(userID string) ([]domain.Alert, error) {
	query := `SELECT id, user_id, symbol, price, condition, is_active FROM alerts WHERE user_id = $1`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []domain.Alert
	for rows.Next() {
		var alert domain.Alert
		if err := rows.Scan(
			&alert.ID,
			&alert.UserID,
			&alert.Symbol,
			&alert.ThresholdPrice,
			&alert.Condition,
			&alert.IsActive,
		); err != nil {
			return nil, err
		}
		alerts = append(alerts, alert)
	}

	return alerts, rows.Err()
}

func (r *PostgresAlertRepository) FindActiveAlerts() ([]domain.Alert, error) {
	query := `SELECT id, user_id, symbol, price, condition, is_active FROM alerts WHERE is_active = true`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []domain.Alert
	for rows.Next() {
		var alert domain.Alert
		if err := rows.Scan(
			&alert.ID,
			&alert.UserID,
			&alert.Symbol,
			&alert.ThresholdPrice,
			&alert.Condition,
			&alert.IsActive,
		); err != nil {
			return nil, err
		}
		alerts = append(alerts, alert)
	}

	return alerts, rows.Err()
}

func (r *PostgresAlertRepository) UpdateAlert(alert *domain.Alert) error {
	query := `UPDATE alerts SET price = $1, condition = $2, is_active = $3 WHERE id = $4`

	result, err := r.db.Exec(query, alert.ThresholdPrice, alert.Condition, alert.IsActive, alert.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("alert not found")
	}

	return nil
}

func (r *PostgresAlertRepository) DeleteAlert(alertID string) error {
	query := `DELETE FROM alerts WHERE id = $1`

	_, err := r.db.Exec(query, alertID)
	return err
}

type PostgresTransactionRepository struct {
	db *sql.DB
}

func NewPostgresTransactionRepository(db *sql.DB) *PostgresTransactionRepository {
	return &PostgresTransactionRepository{db: db}
}

func (r *PostgresTransactionRepository) SaveTransaction(transaction *domain.Transaction) error {
	if transaction == nil || transaction.UserID == "" {
		return errors.New("invalid transaction")
	}

	query := `
		INSERT INTO transactions (user_id, symbol, type, quantity, price, amount, fee, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(query,
		transaction.UserID,
		transaction.Symbol,
		transaction.Type,
		transaction.Quantity,
		transaction.Price,
		transaction.Amount,
		transaction.Fee,
		transaction.CreatedAt,
	)

	return err
}

func (r *PostgresTransactionRepository) FindTransactionsByUser(userID string) ([]domain.Transaction, error) {
	query := `
		SELECT user_id, symbol, type, quantity, price, amount, fee, timestamp
		FROM transactions
		WHERE user_id = $1
		ORDER BY timestamp DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []domain.Transaction
	for rows.Next() {
		var txn domain.Transaction
		if err := rows.Scan(
			&txn.UserID,
			&txn.Symbol,
			&txn.Type,
			&txn.Quantity,
			&txn.Price,
			&txn.Amount,
			&txn.Fee,
			&txn.CreatedAt,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, txn)
	}

	return transactions, rows.Err()
}

func (r *PostgresTransactionRepository) FindTransactionsBySymbol(userID, symbol string) ([]domain.Transaction, error) {
	query := `
		SELECT user_id, symbol, type, quantity, price, amount, fee, timestamp
		FROM transactions
		WHERE user_id = $1 AND symbol = $2
		ORDER BY timestamp DESC
	`

	rows, err := r.db.Query(query, userID, symbol)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []domain.Transaction
	for rows.Next() {
		var txn domain.Transaction
		if err := rows.Scan(
			&txn.UserID,
			&txn.Symbol,
			&txn.Type,
			&txn.Quantity,
			&txn.Price,
			&txn.Amount,
			&txn.Fee,
			&txn.CreatedAt,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, txn)
	}

	return transactions, rows.Err()
}

func (r *PostgresTransactionRepository) GetUserBalance(userID string) (float64, error) {
	query := `
		SELECT COALESCE(SUM(
			CASE
				WHEN type = 'DEPOSIT' THEN amount
				WHEN type = 'WITHDRAWAL' OR type = 'BUY' THEN -(amount + fee)
				WHEN type = 'SELL' THEN (amount - fee)
				WHEN type = 'BROKERAGE_FEE' THEN -fee
				ELSE 0
			END
		), 0) as balance
		FROM transactions
		WHERE user_id = $1
	`

	balance := 100000.0
	var txnBalance sql.NullFloat64

	err := r.db.QueryRow(query, userID).Scan(&txnBalance)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	if txnBalance.Valid {
		return balance + txnBalance.Float64, nil
	}

	return balance, nil
}
