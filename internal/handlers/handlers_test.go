package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	repository "gitlab.com/digineat/go-broker-test/internal/db"
)

type mockRepo struct{}

func (m *mockRepo) CreateTrade(t repository.Trade) error           { return nil }
func (m *mockRepo) GetStats(acc string) (*repository.Stats, error) { return &repository.Stats{}, nil }
func (m *mockRepo) UpdateStats(acc string, profit float64) error   { return nil }
func (m *mockRepo) Ping() error                                    { return nil }

const validTradeJSON = `{
	"account": "test",
	"symbol": "BTCUSD",
	"volume": 1.0,
	"open": 50000,
	"close": 51000,
	"side": "buy"
}`

func TestTradeValidation(t *testing.T) {
	h := NewHandlers(&mockRepo{})

	tests := []struct {
		body         string
		expectedCode int
	}{
		{`{}`, 400},
		{`{"account":"","symbol":"BTCUSD","volume":1}`, 400},
		{validTradeJSON, 200},
	}

	for _, tt := range tests {
		req := httptest.NewRequest("POST", "/trades", bytes.NewBufferString(tt.body))
		w := httptest.NewRecorder()
		h.CreateTrade(w, req)

		if w.Code != tt.expectedCode {
			t.Errorf("Expected %d, got %d for %s", tt.expectedCode, w.Code, tt.body)
		}
	}
}

func TestGetStatsHandler(t *testing.T) {
	h := NewHandlers(&mockRepo{})

	req := httptest.NewRequest("GET", "/stats/acc123", nil)
	w := httptest.NewRecorder()
	h.GetStats(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestHealthCheck(t *testing.T) {
	h := NewHandlers(&mockRepo{})
	req := httptest.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()

	h.HealthCheck(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
