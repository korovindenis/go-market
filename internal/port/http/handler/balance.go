package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/korovindenis/go-market/internal/domain/entity"
)

func (h *Handler) GetBalance(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := h.getUserIDFromCtx(c)
	if err != nil {
		c.Error(fmt.Errorf("%s %w", "Handler GetBalance getUserIDFromCtx", err))
		c.AbortWithError(http.StatusInternalServerError, entity.ErrInternalServerError)
		return
	}
	user := entity.User{
		ID: userID,
	}

	balance, err := h.usecase.GetBalance(ctx, user)
	if err != nil {
		c.Error(fmt.Errorf("%s %w", "Handler GetBalance usecase.GetBalance", err))
		c.AbortWithError(http.StatusInternalServerError, entity.ErrInternalServerError)
		return
	}

	c.JSON(http.StatusOK, balance)
}

func (h *Handler) WithdrawBalance(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := h.getUserIDFromCtx(c)
	if err != nil {
		c.Error(fmt.Errorf("%s %w", "Handler GetBalance getUserIDFromCtx", err))
		c.AbortWithError(http.StatusInternalServerError, entity.ErrInternalServerError)
		return
	}
	user := entity.User{
		ID: userID,
	}

	var balance entity.BalanceUpdate
	if err := c.ShouldBindJSON(&balance); err != nil {
		c.Error(fmt.Errorf("%s %w", "Handler WithdrawBalance Json", err))
		c.AbortWithError(http.StatusInternalServerError, entity.ErrInternalServerError)
		return
	}

	// check input data
	if err := balance.IsValidNumber(); err != nil {
		c.Error(fmt.Errorf("%s %w", "Handler WithdrawBalance IsValidNumber", err))
		c.AbortWithError(http.StatusUnprocessableEntity, entity.ErrUnprocessableEntity)
		return
	}

	if err := h.usecase.WithdrawBalance(ctx, balance, user); err != nil {
		if errors.Is(err, entity.ErrInsufficientBalance) {
			c.Error(fmt.Errorf("%s %w", "Handler WithdrawBalance WithdrawBalance Insufficient Balance", err))
			c.AbortWithError(http.StatusPaymentRequired, entity.ErrInsufficientBalance)
			return
		}
		c.Error(fmt.Errorf("%s %w", "Handler WithdrawBalance WithdrawBalance", err))
		c.AbortWithError(http.StatusInternalServerError, entity.ErrInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
