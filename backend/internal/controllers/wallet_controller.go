package controllers

import (
	"net/http"

	"github.com/dfanso/parkme-backend/internal/services"
	"github.com/dfanso/parkme-backend/pkg/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WalletController struct {
	service *services.WalletService
}

func NewWalletController(service *services.WalletService) *WalletController {
	return &WalletController{
		service: service,
	}
}

type TopUpRequest struct {
	Amount float64 `json:"amount"`
}

func (c *WalletController) TopUp(ctx echo.Context) error {
	var req TopUpRequest
	if err := ctx.Bind(&req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body", err)
	}

	userID := ctx.Get("userID").(primitive.ObjectID)

	transaction, err := c.service.TopUp(ctx.Request().Context(), userID, req.Amount)
	if err != nil {
		switch err {
		case services.ErrInvalidAmount:
			return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid amount", err)
		default:
			return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to top up wallet", err)
		}
	}

	return utils.SuccessResponse(ctx, http.StatusOK, "Wallet topped up successfully", transaction)
}

func (c *WalletController) GetBalance(ctx echo.Context) error {
	userID := ctx.Get("userID").(primitive.ObjectID)

	wallet, err := c.service.GetOrCreateWallet(ctx.Request().Context(), userID)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get balance", err)
	}

	return utils.SuccessResponse(ctx, http.StatusOK, "Balance retrieved successfully", map[string]float64{
		"balance": wallet.Balance,
	})
}

func (c *WalletController) GetTransactions(ctx echo.Context) error {
	userID := ctx.Get("userID").(primitive.ObjectID)

	// Just create wallet if it doesn't exist
	_, err := c.service.GetOrCreateWallet(ctx.Request().Context(), userID)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get transactions", err)
	}

	transactions, err := c.service.GetTransactions(ctx.Request().Context(), userID)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get transactions", err)
	}

	return utils.SuccessResponse(ctx, http.StatusOK, "Transactions retrieved successfully", transactions)
}
