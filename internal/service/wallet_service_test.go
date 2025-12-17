package service

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"testing"

	"testtask/internal/domain"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWalletRepository struct {
	mock.Mock
}

func (m *MockWalletRepository) GetBalanceForUpdate(ctx context.Context, walletID uuid.UUID) (decimal.Decimal, error) {
	args := m.Called(ctx, walletID)
	return args.Get(0).(decimal.Decimal), args.Error(1)
}

func (m *MockWalletRepository) UpdateBalance(ctx context.Context, walletID uuid.UUID, newBalance decimal.Decimal) error {
	args := m.Called(ctx, walletID, newBalance)
	return args.Error(0)
}

func (m *MockWalletRepository) GetBalance(ctx context.Context, walletID uuid.UUID) (decimal.Decimal, error) {
	// ...
	return decimal.Zero, nil
}
func (m *MockWalletRepository) Create(ctx context.Context, wallet *domain.Wallet) error {
	// ...
	return nil
}

type MockUoW struct {
	mock.Mock
	Repo *MockWalletRepository
}

func (m *MockUoW) Wallets() domain.WalletRepository {
	return m.Repo
}

func (m *MockUoW) Do(ctx context.Context, fn func(uow domain.UnitOfWork) error) error {
	return fn(m)
}

func TestWalletService_PerformOperation(t *testing.T) {
	logger := zap.NewNop()
	walletID := uuid.New()

	tests := []struct {
		name          string
		req           domain.OperationRequest
		setupMock     func(repo *MockWalletRepository)
		expectedError error
	}{
		{
			name: "Успешное пополнение",
			req:  domain.OperationRequest{ID: walletID, OperationType: domain.Deposit, Amount: decimal.NewFromInt(100)},
			setupMock: func(repo *MockWalletRepository) {
				repo.On("GetBalanceForUpdate", mock.Anything, walletID).Return(decimal.NewFromInt(50), nil)
				repo.On("UpdateBalance", mock.Anything, walletID, decimal.NewFromInt(150)).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Успешное списание",
			req:  domain.OperationRequest{ID: walletID, OperationType: domain.Withdraw, Amount: decimal.NewFromInt(50)},
			setupMock: func(repo *MockWalletRepository) {
				repo.On("GetBalanceForUpdate", mock.Anything, walletID).Return(decimal.NewFromInt(100), nil)
				repo.On("UpdateBalance", mock.Anything, walletID, decimal.NewFromInt(50)).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Ошибка: кошелек не найден",
			req:  domain.OperationRequest{ID: walletID, OperationType: domain.Deposit, Amount: decimal.NewFromInt(100)},
			setupMock: func(repo *MockWalletRepository) {
				repo.On("GetBalanceForUpdate", mock.Anything, walletID).Return(decimal.Zero, domain.ErrWalletNotFound)
			},
			expectedError: domain.ErrWalletNotFound,
		},
		{
			name: "Ошибка: недостаточно средств",
			req:  domain.OperationRequest{ID: walletID, OperationType: domain.Withdraw, Amount: decimal.NewFromInt(200)},
			setupMock: func(repo *MockWalletRepository) {
				repo.On("GetBalanceForUpdate", mock.Anything, walletID).Return(decimal.NewFromInt(100), nil)
			},
			expectedError: domain.ErrInsufficientFunds,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockWalletRepository)
			mockUOW := &MockUoW{Repo: mockRepo}

			tt.setupMock(mockRepo)

			service := NewWalletService(mockUOW, logger)
			err := service.PerformOperation(context.Background(), tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectedError))
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
