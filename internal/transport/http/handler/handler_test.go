package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"

	"testtask/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWalletService struct {
	mock.Mock
}

func (m *MockWalletService) CreateWallet(ctx context.Context) (*domain.Wallet, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Wallet), args.Error(1)
}

func (m *MockWalletService) PerformOperation(ctx context.Context, req domain.OperationRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockWalletService) GetBalance(ctx context.Context, id uuid.UUID) (decimal.Decimal, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(decimal.Decimal), args.Error(1)
}

func setupTest() (*gin.Engine, *MockWalletService) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("logger", zap.NewNop())
		c.Next()
	})

	mockService := new(MockWalletService)
	handler := NewHandler(mockService)

	v1 := router.Group("/api/v1")
	{
		v1.POST("/wallets", handler.CreateWallet)
		v1.GET("/wallets/:walletID", handler.GetBalance)
		v1.POST("/wallet", handler.Operation)
	}

	return router, mockService
}

func TestHandler_CreateWallet(t *testing.T) {
	router, mockService := setupTest()

	t.Run("Success", func(t *testing.T) {
		walletID := uuid.New()
		expectedWallet := &domain.Wallet{ID: walletID, Balance: decimal.Zero}

		mockService.On("CreateWallet", mock.Anything).Return(expectedWallet, nil).Once()

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/wallets", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var respBody map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &respBody)
		assert.Equal(t, walletID.String(), respBody["id"])
		assert.Equal(t, "0", respBody["balance"])

		mockService.AssertExpectations(t)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		mockService.On("CreateWallet", mock.Anything).Return(nil, errors.New("db is down")).Once()

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/wallets", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestHandler_GetBalance(t *testing.T) {
	router, mockService := setupTest()

	t.Run("Success", func(t *testing.T) {
		walletID := uuid.New()
		expectedBalance, _ := decimal.NewFromString("123.45")

		mockService.On("GetBalance", mock.Anything, walletID).Return(expectedBalance, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/wallets/"+walletID.String(), nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var respBody map[string]string
		_ = json.Unmarshal(w.Body.Bytes(), &respBody)
		assert.Equal(t, "123.45", respBody["balance"])
		mockService.AssertExpectations(t)
	})

	t.Run("Wallet Not Found", func(t *testing.T) {
		walletID := uuid.New()
		mockService.On("GetBalance", mock.Anything, walletID).Return(decimal.Zero, domain.ErrWalletNotFound).Once()

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/wallets/"+walletID.String(), nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestHandler_Operation(t *testing.T) {
	router, mockService := setupTest()

	t.Run("Success - Deposit", func(t *testing.T) {
		walletID := uuid.New()
		amount, _ := decimal.NewFromString("100")
		reqBody := map[string]interface{}{
			"walletId":      walletID,
			"operationType": "DEPOSIT",
			"amount":        amount,
		}
		jsonBody, _ := json.Marshal(reqBody)

		expectedReq := domain.OperationRequest{
			ID:            walletID,
			OperationType: "DEPOSIT",
			Amount:        amount,
		}
		mockService.On("PerformOperation", mock.Anything, expectedReq).Return(nil).Once()

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Insufficient Funds", func(t *testing.T) {
		walletID := uuid.New()
		amount, _ := decimal.NewFromString("500")
		reqBody := map[string]interface{}{
			"walletId":      walletID,
			"operationType": "WITHDRAW",
			"amount":        amount,
		}
		jsonBody, _ := json.Marshal(reqBody)

		expectedReq := domain.OperationRequest{
			ID:            walletID,
			OperationType: "WITHDRAW",
			Amount:        amount,
		}
		mockService.On("PerformOperation", mock.Anything, expectedReq).Return(domain.ErrInsufficientFunds).Once()

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Bad Request - Invalid JSON", func(t *testing.T) {
		invalidJson := []byte(`{"walletId": "not-a-uuid"`)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewBuffer(invalidJson))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
