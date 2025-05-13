package main

import (
	"database/sql"
	"flag"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	repository "gitlab.com/digineat/go-broker-test/internal/db"
)

func main() {
	// Command line flags
	// dbPath := flag.String("db", "../server/data.db", "path to SQLite database")
	dbPath := flag.String("db", "data.db", "path to SQLite database")
	pollInterval := flag.Duration("poll", 100*time.Millisecond, "polling interval")
	flag.Parse()

	// Initialize database connection
	db, err := sql.Open("sqlite3", *dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	repo := repository.NewRepository(db)

	log.Printf("Worker started with polling interval: %v", *pollInterval)

	// Main worker loop
	for {
		trades, err := repo.GetUnprocessedTrades()
		if err != nil {
			log.Printf("Failed to get unprocessed trades: %v", err)
			continue
		}
		for _, trade := range trades {
			log.Printf("Processing trade: %v", trade)
			profit := calculateProfit(trade)
			if err := repo.UpdateStats(trade.Account, profit); err == nil {
				repo.MarkAsProcessed(trade.ID)
			}

		}
		// Sleep for the specified interval
		time.Sleep(*pollInterval)
	}
}

func calculateProfit(trade repository.Trade) float64 {
	lot := 100000.0
	profit := (trade.Close - trade.Open) * trade.Volume * lot
	if trade.Side == "sell" {
		profit = -profit
	}
	return profit

}
