package domain

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	ErrWalletNotFound = errors.New("wallet not found")
)

type WalletRepository interface {
	GetBalanceForUpdate(ctx context.Context, walletID uuid.UUID) (decimal.Decimal, error)
	UpdateBalance(ctx context.Context, walletID uuid.UUID, newBalance decimal.Decimal) error

	GetBalance(ctx context.Context, walletID uuid.UUID) (decimal.Decimal, error)
	Create(ctx context.Context, wallet *Wallet) error
}

type UnitOfWork interface {
	// Wallets возвращает репозиторий для кошельков
	Wallets() WalletRepository

	// Do выполняет функцию в транзакции
	// Если функция возвращает ошибку, транзакция откатывается
	// Если нет - коммитится
	Do(ctx context.Context, fn func(uow UnitOfWork) error) error
}
