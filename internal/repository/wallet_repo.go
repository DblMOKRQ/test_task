package postgres

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"

	"testtask/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

const (
	getBalanceForUpdateQuery = `SELECT balance FROM wallets WHERE id = $1 FOR UPDATE;`
	updateBalanceQuery       = `UPDATE wallets SET balance = $1 WHERE id = $2;`
	getBalanceQuery          = `SELECT balance FROM wallets WHERE id = $1;`
	createWalletQuery        = `INSERT INTO wallets (id, balance) VALUES ($1, $2);`
)

// GetBalanceForUpdate получает баланс кошелька, используя пессимистическую блокировку
func (r *WalletRepo) GetBalanceForUpdate(ctx context.Context, id uuid.UUID) (decimal.Decimal, error) {

	var balance decimal.Decimal
	err := r.exec.QueryRow(ctx, getBalanceForUpdateQuery, id).Scan(&balance)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return decimal.Zero, domain.ErrWalletNotFound
		}
		return decimal.Zero, err
	}

	return balance, nil
}

// UpdateBalance обновляет баланс кошелька
func (r *WalletRepo) UpdateBalance(ctx context.Context, id uuid.UUID, newBalance decimal.Decimal) error {
	cmdTag, err := r.exec.Exec(ctx, updateBalanceQuery, newBalance, id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() != 1 {
		return errors.New("wallet not found or not updated")
	}

	return nil
}

// GetBalance получает текущий баланс без блокировки
func (r *WalletRepo) GetBalance(ctx context.Context, id uuid.UUID) (decimal.Decimal, error) {
	var balance decimal.Decimal
	err := r.exec.QueryRow(ctx, getBalanceQuery, id).Scan(&balance)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return decimal.Zero, domain.ErrWalletNotFound
		}
		return decimal.Zero, err
	}

	return balance, nil
}

// Create создает новый кошелек в базе данных.
func (r *WalletRepo) Create(ctx context.Context, wallet *domain.Wallet) error {
	r.log.Debug("Executing create wallet query", zap.String("id", wallet.ID.String()))

	cmdTag, err := r.exec.Exec(ctx, createWalletQuery, wallet.ID, wallet.Balance)
	if err != nil {
		r.log.Error("Failed to execute insert query for new wallet", zap.Error(err))
		return fmt.Errorf("failed to execute insert query for new wallet: %w", err)
	}
	if cmdTag.RowsAffected() != 1 {
		return errors.New("failed to create wallet: no rows affected")
	}

	return nil
}
