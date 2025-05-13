package repository

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestRepository(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:?_foreign_keys=1")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := RunMigrations(db); err != nil {
		t.Fatal(err)
	}

	repo := NewRepository(db)

	t.Run("CreateTrade", func(t *testing.T) {
		err := repo.CreateTrade(Trade{
			Account: "acc1",
			Symbol:  "EURUSD",
			Volume:  1.0,
			Open:    1.08,
			Close:   1.09,
			Side:    "buy",
		})
		if err != nil {
			t.Errorf("CreateTrade failed: %v", err)
		}
	})

	t.Run("GetUnprocessedTrades", func(t *testing.T) {
		_, err := db.Exec(`INSERT INTO trades_q 
			(account, symbol, volume, open, close, side, processed) 
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			"acc2", "USDJPY", 2.0, 110.5, 111.0, "sell", false)
		if err != nil {
			t.Fatal(err)
		}

		trades, err := repo.GetUnprocessedTrades()
		if err != nil {
			t.Errorf("GetUnprocessedTrades failed: %v", err)
		}

		if len(trades) == 0 {
			t.Error("Expected unprocessed trades, got none")
		}
	})

	t.Run("UpdateStats", func(t *testing.T) {
		err := repo.UpdateStats("acc3", 1500.0)
		if err != nil {
			t.Errorf("UpdateStats failed: %v", err)
		}

		var trades int
		err = db.QueryRow("SELECT trades FROM account_stats WHERE account = ?", "acc3").Scan(&trades)
		if err != nil {
			t.Errorf("Failed to verify UpdateStats: %v", err)
		}

		if trades != 1 {
			t.Errorf("Expected 1 trade, got %d", trades)
		}
	})

	t.Run("MarkAsProcessed", func(t *testing.T) {
		res, err := db.Exec(`INSERT INTO trades_q 
			(account, symbol, volume, open, close, side) 
			VALUES (?, ?, ?, ?, ?, ?)`,
			"acc4", "GBPUSD", 1.5, 1.25, 1.26, "buy")
		if err != nil {
			t.Fatal(err)
		}

		id, _ := res.LastInsertId()
		err = repo.MarkAsProcessed(int(id))
		if err != nil {
			t.Errorf("MarkAsProcessed failed: %v", err)
		}

		var processed bool
		err = db.QueryRow("SELECT processed FROM trades_q WHERE id = ?", id).Scan(&processed)
		if err != nil || !processed {
			t.Error("Failed to mark trade as processed")
		}
	})

	t.Run("GetStats", func(t *testing.T) {
		_, err := db.Exec(`INSERT INTO account_stats 
			(account, trades, profit) 
			VALUES (?, ?, ?)`,
			"acc5", 3, 4500.0)
		if err != nil {
			t.Fatal(err)
		}

		stats, err := repo.GetStats("acc5")
		if err != nil {
			t.Errorf("GetStats failed: %v", err)
		}

		if stats.Trades != 3 || stats.Profit != 4500.0 {
			t.Errorf("Unexpected stats: %+v", stats)
		}
	})

	t.Run("Ping", func(t *testing.T) {
		err := repo.Ping()
		if err != nil {
			t.Errorf("Ping failed: %v", err)
		}
	})
}
