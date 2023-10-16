package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/korovindenis/go-market/internal/domain/entity"
)

func (h *Handler) Order(c *gin.Context) {
	// TO DO
	// 200 — номер заказа уже был загружен этим пользователем;
	// 409 — номер заказа уже был загружен другим пользователем;

	ctx := c.Request.Context()

	// get input data
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Error(fmt.Errorf("%s %w", "Handler Order ShouldBindBodyWith", err))
		c.AbortWithError(http.StatusBadRequest, entity.ErrStatusBadRequest)
		return
	}
	num, err := strconv.ParseUint(string(body), 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("%s %w", "Handler Order ParseUint", err))
		c.AbortWithError(http.StatusUnprocessableEntity, entity.ErrUnprocessableEntity)
		return
	}
	order := entity.Order{
		Number: num,
	}

	// check input data
	if err := order.IsValidNumber(); err != nil {
		c.Error(fmt.Errorf("%s %w", "Handler Order IsValidNumber", err))
		c.AbortWithError(http.StatusUnprocessableEntity, entity.ErrUnprocessableEntity)
		return
	}

	userIdRaw, ok := c.Get("userId")
	if !ok {
		c.Error(fmt.Errorf("%s %w", "Handler Order Get userId", err))
		c.AbortWithError(http.StatusInternalServerError, entity.ErrInternalServerError)
		return
	}
	userId, ok := userIdRaw.(uint64)
	if !ok {
		c.Error(fmt.Errorf("%s %w", "Handler Order Get userId to uint64", err))
		c.AbortWithError(http.StatusInternalServerError, entity.ErrInternalServerError)
		return
	}
	user := entity.User{
		Id: userId,
	}

	if err := h.usecase.AddOrder(ctx, order, user); err != nil {
		c.Error(fmt.Errorf("%s %w", "Handler Order AddOrder", err))

		// order has already been uploaded by another user
		if errors.Is(err, entity.ErrOrderAlreadyUploadedAnotherUser) {
			c.Error(fmt.Errorf("%s %w", "Handler Order AddOrder ErrOrderAlreadyUploadedAnotherUser", err))
			c.AbortWithError(http.StatusConflict, entity.ErrOrderAlreadyUploadedAnotherUser)
			return
		}
		// order has already been uploaded by this user
		if errors.Is(err, entity.ErrOrderAlreadyUploaded) {
			c.Status(http.StatusOK)
			return
		}
		c.AbortWithError(http.StatusInternalServerError, entity.ErrInternalServerError)
		return
	}

	c.Status(http.StatusAccepted)
}
