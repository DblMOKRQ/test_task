package service

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"

	"testtask/internal/domain"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var initialBalance = decimal.NewFromInt(0)

type WalletService struct {
	uowFactory domain.UnitOfWork
	log        *zap.Logger
}

func NewWalletService(uowFactory domain.UnitOfWork, log *zap.Logger) *WalletService {
	return &WalletService{
		uowFactory: uowFactory,
		log:        log.Named("WalletService"),
	}
}

func (s *WalletService) PerformOperation(ctx context.Context, req domain.OperationRequest) error {
	s.log.Debug("Perform operation", zap.Any("req", req))
	if req.ID == uuid.Nil {
		s.log.Warn("Wallet ID is nil", zap.Any("req", req))
		return domain.ErrIDIsNil
	}

	if req.Amount.IsZero() || req.Amount.IsNegative() {
		s.log.Warn("Wallet amount is zero or is negative", zap.Any("req", req))
		return domain.ErrAmountZeroOrNegative
	}

	return s.uowFactory.Do(ctx, func(uow domain.UnitOfWork) error {
		walletRepo := uow.Wallets()

		balance, err := walletRepo.GetBalanceForUpdate(ctx, req.ID)
		if err != nil {
			if errors.Is(err, domain.ErrWalletNotFound) {
				s.log.Error("Wallet not found", zap.Any("req", req))
				return domain.ErrWalletNotFound
			}
			s.log.Error("Error getting balance for req", zap.Any("req", req), zap.Error(err))
			return fmt.Errorf("failed to get balance: %w", err)
		}

		var newBalance decimal.Decimal
		switch req.OperationType {
		case domain.Deposit:
			s.log.Debug("Deposit operation", zap.String("balance", balance.String()), zap.Any("req", req))
			newBalance = balance.Add(req.Amount)
		case domain.Withdraw:
			s.log.Debug("Withdraw operation", zap.String("balance", balance.String()), zap.Any("req", req))
			if balance.LessThan(req.Amount) {
				s.log.Warn("Insufficient funds", zap.String("balance", balance.String()), zap.Any("req", req))
				return domain.ErrInsufficientFunds
			}
			newBalance = balance.Sub(req.Amount)
		default:
			s.log.Warn("Unknown operation type", zap.Any("req", req))
			return domain.ErrUnknownOperationType
		}
		return walletRepo.UpdateBalance(ctx, req.ID, newBalance)

	})

}

func (s *WalletService) GetBalance(ctx context.Context, id uuid.UUID) (decimal.Decimal, error) {
	s.log.Debug("Get balance", zap.Any("id", id))
	balance, err := s.uowFactory.Wallets().GetBalance(ctx, id)
	if err != nil {
		s.log.Error("Failed to get balance", zap.Error(err), zap.Any("id", id))
		return decimal.Zero, err
	}
	return balance, nil
}
func (s *WalletService) CreateWallet(ctx context.Context) (*domain.Wallet, error) {
	newID := uuid.New()
	s.log.Debug("Generated new wallet ID", zap.String("id", newID.String()))

	newWallet := &domain.Wallet{
		ID:      newID,
		Balance: initialBalance,
	}

	err := s.uowFactory.Wallets().Create(ctx, newWallet)
	if err != nil {
		s.log.Error("Failed to save new wallet to repository", zap.Error(err))
		return nil, fmt.Errorf("failed to save new wallet to repository: %w", err)
	}

	return newWallet, nil
}
