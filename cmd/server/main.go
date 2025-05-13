package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	repository "gitlab.com/digineat/go-broker-test/internal/db"
	"gitlab.com/digineat/go-broker-test/internal/handlers"
)

func main() {
	// Command line flags
	dbPath := flag.String("db", "data.db", "path to SQLite database")
	listenAddr := flag.String("listen", "8080", "HTTP server listen address")
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
	handler := handlers.NewHandlers(repo)
	// Initialize HTTP server
	mux := http.NewServeMux()

	// POST /trades endpoint
	mux.HandleFunc("POST /trades", handler.CreateTrade)

	// GET /stats/{acc} endpoint
	mux.HandleFunc("GET /stats/{acc}", handler.GetStats)

	// GET /healthz endpoint
	mux.HandleFunc("GET /healthz", handler.HealthCheck)

	// Start server
	serverAddr := fmt.Sprintf(":%s", *listenAddr)
	log.Printf("Starting server on %s", serverAddr)
	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
