package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/korovindenis/go-market/internal/domain/entity"
)

func (h *Handler) Withdrawals(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := h.GetUserIDFromCtx(c)
	if err != nil {
		c.Error(fmt.Errorf("%s %w", "Handler GetBalance GetUserIDFromCtx", err))
		c.AbortWithError(http.StatusInternalServerError, entity.ErrInternalServerError)
		return
	}
	user := entity.User{
		ID: userID,
	}

	withdrawals, err := h.usecase.Withdrawals(ctx, user)
	if err != nil {
		if errors.Is(err, entity.ErrNoContent) {
			c.Error(fmt.Errorf("%s %w", "Handler Withdrawals usecase.Withdrawals", err))
			c.AbortWithError(http.StatusNoContent, entity.ErrNoContent)
			return
		}
		c.Error(fmt.Errorf("%s %w", "Handler Withdrawals usecase.Withdrawals", err))
		c.AbortWithError(http.StatusInternalServerError, entity.ErrInternalServerError)
		return
	}

	c.JSON(http.StatusOK, withdrawals)
}
