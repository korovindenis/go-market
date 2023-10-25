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

func (h *Handler) SetOrder(c *gin.Context) {
	ctx := c.Request.Context()

	// get input data
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Error(fmt.Errorf("%s %w", "Handler SetOrder ShouldBindBodyWith", err))
		c.AbortWithError(http.StatusBadRequest, entity.ErrStatusBadRequest)
		return
	}
	num, err := strconv.ParseUint(string(body), 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("%s %w", "Handler SetOrder ParseUint", err))
		c.AbortWithError(http.StatusUnprocessableEntity, entity.ErrUnprocessableEntity)
		return
	}
	order := entity.Order{
		Number: num,
	}

	// check input data
	if err := order.IsValidNumber(); err != nil {
		c.Error(fmt.Errorf("%s %w", "Handler SetOrder IsValidNumber", err))
		c.AbortWithError(http.StatusUnprocessableEntity, entity.ErrUnprocessableEntity)
		return
	}
	userID, err := h.getUserIdFromCtx(c)
	if err != nil {
		c.Error(fmt.Errorf("%s %w", "Handler SetOrder getUserIdFromCtx", err))
		c.AbortWithError(http.StatusInternalServerError, entity.ErrInternalServerError)
		return
	}
	user := entity.User{
		ID: userID,
	}

	if err := h.usecase.AddOrder(ctx, order, user); err != nil {
		c.Error(fmt.Errorf("%s %w", "Handler SetOrder AddOrder", err))

		// order has already been uploaded by another user
		if errors.Is(err, entity.ErrOrderAlreadyUploadedAnotherUser) {
			c.Error(fmt.Errorf("%s %w", "Handler SetOrder AddOrder ErrOrderAlreadyUploadedAnotherUser", err))
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

func (h *Handler) GetOrder(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := h.getUserIdFromCtx(c)
	if err != nil {
		c.Error(fmt.Errorf("%s %w", "Handler GetOrder getUserIdFromCtx", err))
		c.AbortWithError(http.StatusInternalServerError, entity.ErrInternalServerError)
		return
	}
	user := entity.User{
		ID: userID,
	}

	h.usecase.GetOrder(ctx, user)
}
