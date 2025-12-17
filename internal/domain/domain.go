package domain

import (
	"errors"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	ErrOperationNotSupported = errors.New("operation not supported")
	ErrIDIsNil               = errors.New("id is nil")
	ErrAmountZeroOrNegative  = errors.New("amount is zero or is negative")
	ErrInsufficientFunds     = errors.New("insufficient funds")
	ErrUnknownOperationType  = errors.New("unknown operation type")
)

type OperationType string

const (
	Deposit  OperationType = "DEPOSIT"
	Withdraw OperationType = "WITHDRAW"
)

type OperationRequest struct {
	ID            uuid.UUID
	OperationType OperationType
	Amount        decimal.Decimal
}

type Wallet struct {
	ID      uuid.UUID
	Balance decimal.Decimal
}
