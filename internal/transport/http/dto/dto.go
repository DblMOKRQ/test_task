package dto

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type OperationWalletRequestDTO struct {
	WalletID      uuid.UUID       `json:"wallet"`
	OperationType string          `json:"operation_type"`
	Amount        decimal.Decimal `json:"amount"`
}

type CreateWalletResponseDTO struct {
	ID      uuid.UUID       `json:"id"`
	Balance decimal.Decimal `json:"balance"`
}
