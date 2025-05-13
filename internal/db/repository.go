package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

type Trade struct {
	ID      int
	Account string
	Symbol  string
	Volume  float64
	Open    float64
	Close   float64
	Side    string
}

type Stats struct {
	Account string  `json:"account"`
	Trades  int     `json:"trades"`
	Profit  float64 `json:"profit"`
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	if err := RunMigrations(db); err != nil {
		log.Fatalf("Migrations failed: %v", err)
	}
	return &Repository{db: db}
}

func (r *Repository) CreateTrade(trade Trade) error {
	_, err := r.db.Exec(`
		INSERT INTO trades_q 
		(account, symbol, volume, open, close, side)
		VALUES (?, ?, ?, ?, ?, ?)`,
		trade.Account, trade.Symbol, trade.Volume,
		trade.Open, trade.Close, trade.Side,
	)
	return err
}

func (r *Repository) GetUnprocessedTrades() ([]Trade, error) {
	tx, err := r.db.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	rows, err := tx.Query(`
		SELECT id, account, symbol, volume, open, close, side 
		FROM trades_q 
		WHERE processed = FALSE 
		ORDER BY created_at`)
	if err != nil {
		return nil, fmt.Errorf("failed to query trades: %w", err)
	}
	defer rows.Close()

	var trades []Trade
	for rows.Next() {
		var t Trade
		err := rows.Scan(&t.ID, &t.Account, &t.Symbol, &t.Volume, &t.Open, &t.Close, &t.Side)
		if err != nil {
			return nil, fmt.Errorf("failed to scan trade: %w", err)
		}
		trades = append(trades, t)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return trades, nil
}

func (r *Repository) GetStats(account string) (*Stats, error) {
	var stats Stats
	err := r.db.QueryRow(`
		SELECT account, trades, profit 
		FROM account_stats 
		WHERE account = ?`, account).Scan(
		&stats.Account,
		&stats.Trades,
		&stats.Profit,
	)
	return &stats, err
}

func (r *Repository) UpdateStats(account string, profit float64) error {
	_, err := r.db.Exec(`
		INSERT INTO account_stats 
		(account, trades, profit) 
		VALUES (?, 1, ?)
		ON CONFLICT(account) 
		DO UPDATE SET 
			trades = trades + 1,
			profit = profit + excluded.profit`,
		account, profit)
	return err
}

func (r *Repository) MarkAsProcessed(id int) error {
	_, err := r.db.Exec("UPDATE trades_q SET processed = TRUE WHERE id = ?", id)
	return err
}

func (r *Repository) Close() error {
	return r.db.Close()
}
func (r *Repository) Ping() error {
	return r.db.Ping()
}

func RunMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS trades_q (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			account TEXT NOT NULL,
			symbol TEXT NOT NULL,
			volume REAL NOT NULL,
			open REAL NOT NULL,
			close REAL NOT NULL,
			side TEXT NOT NULL,
			processed BOOLEAN DEFAULT FALSE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS account_stats (
			account TEXT PRIMARY KEY,
			trades INTEGER DEFAULT 0,
			profit REAL DEFAULT 0
		)`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}
	return nil
}
