package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	repository "gitlab.com/digineat/go-broker-test/internal/db"
)

type TradeRequest struct {
	Account string  `json:"account"`
	Symbol  string  `json:"symbol"`
	Volume  float64 `json:"volume"`
	Open    float64 `json:"open"`
	Close   float64 `json:"close"`
	Side    string  `json:"side"`
}

func (req *TradeRequest) ToDomain() repository.Trade {
	return repository.Trade{
		Account: req.Account,
		Symbol:  req.Symbol,
		Volume:  req.Volume,
		Open:    req.Open,
		Close:   req.Close,
		Side:    req.Side,
	}
}

type TradeRepo interface {
	CreateTrade(trade repository.Trade) error
	GetStats(account string) (*repository.Stats, error)
	UpdateStats(account string, profit float64) error
	Ping() error
}

type Handlers struct {
	repo TradeRepo
}

func NewHandlers(repo TradeRepo) *Handlers {
	return &Handlers{repo: repo}
}

func (h *Handlers) CreateTrade(w http.ResponseWriter, r *http.Request) {
	var trade TradeRequest
	if err := json.NewDecoder(r.Body).Decode(&trade); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if err := trade.validateTrade(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.repo.CreateTrade(trade.ToDomain()); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	log.Printf("Trade created for account %s", trade.Account)
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) GetStats(w http.ResponseWriter, r *http.Request) {
	account := r.PathValue("acc")
	stats, err := h.repo.GetStats(account)
	if err != nil {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}
	log.Printf("Stats retrieved for account %s", account)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (h *Handlers) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if err := h.repo.Ping(); err != nil {
		http.Error(w, "ERROR", http.StatusInternalServerError)
		return
	}
	log.Println("Health check OK")
	w.Write([]byte("OK"))
	w.WriteHeader(http.StatusOK)
}

func (req *TradeRequest) validateTrade() error {
	if req.Account == "" {
		return errors.New("Account is required")
	}
	if req.Symbol == "" {
		return errors.New("Symbol is required")
	}
	if req.Volume == 0 {
		return errors.New("Volume is required")
	}
	if req.Open <= 0 {
		return errors.New("Open is not valid")
	}
	if req.Close <= 0 {
		return errors.New("Close is not valid")
	}
	if req.Side != "buy" && req.Side != "sell" {
		return errors.New("Side is not valid")
	}
	return nil
}
