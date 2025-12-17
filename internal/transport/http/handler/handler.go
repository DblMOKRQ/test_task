package handler

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"net/http"

	"testtask/internal/domain"
	"testtask/internal/transport/http/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type WalletService interface {
	CreateWallet(ctx context.Context) (*domain.Wallet, error)
	PerformOperation(ctx context.Context, wallet domain.OperationRequest) error
	GetBalance(ctx context.Context, id uuid.UUID) (decimal.Decimal, error)
}

type Handler struct {
	walletService WalletService
}

func NewHandler(walletService WalletService) *Handler {
	return &Handler{
		walletService: walletService,
	}
}

func (h *Handler) Operation(c *gin.Context) {
	log := c.MustGet("logger").(*zap.Logger)
	var req dto.OperationWalletRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("Failed to decode request body", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	wallet := domain.OperationRequest{
		ID:            req.WalletID,
		OperationType: domain.OperationType(req.OperationType),
		Amount:        req.Amount,
	}

	err := h.walletService.PerformOperation(c.Request.Context(), wallet)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrIDIsNil), errors.Is(err, domain.ErrAmountZeroOrNegative):
			log.Warn("Invalid request data", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, domain.ErrWalletNotFound):
			log.Warn("Wallet not found", zap.String("wallet_id", wallet.ID.String()))
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, domain.ErrInsufficientFunds):
			log.Warn("Insufficient funds", zap.String("wallet_id", wallet.ID.String()))
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		case errors.Is(err, domain.ErrUnknownOperationType):
			log.Warn("Unknown operation type", zap.String("operation_type", req.OperationType))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			log.Error("Failed to perform operation", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) GetBalance(c *gin.Context) {
	log := c.MustGet("logger").(*zap.Logger)
	walletIDStr := c.Param("id")
	walletID, err := uuid.Parse(walletIDStr)
	if err != nil {
		log.Warn("Failed to parse walletID", zap.String("walletIDStr", walletIDStr), zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "walletID is not a valid UUID"})
		return
	}

	balance, err := h.walletService.GetBalance(c.Request.Context(), walletID)
	if err != nil {
		if errors.Is(err, domain.ErrWalletNotFound) {
			log.Warn("Wallet not found for get balance", zap.String("walletID", walletID.String()))
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Error("Failed to get balance", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"walletID": walletID,
		"balance":  balance.String()},
	)
}

func (h *Handler) CreateWallet(c *gin.Context) {
	log := c.MustGet("logger").(*zap.Logger)

	log.Info("Handling create wallet request")

	newWallet, err := h.walletService.CreateWallet(c.Request.Context())
	if err != nil {
		log.Error("Failed to create wallet", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	log.Info("Successfully created new wallet", zap.String("wallet_id", newWallet.ID.String()))

	responseDTO := dto.CreateWalletResponseDTO{
		ID:      newWallet.ID,
		Balance: newWallet.Balance,
	}
	c.JSON(http.StatusCreated, responseDTO)
}
