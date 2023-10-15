package handler

import (
	"encoding/binary"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/korovindenis/go-market/internal/domain/entity"
)

func (h *Handler) Order(c *gin.Context) {
	// 200 — номер заказа уже был загружен этим пользователем;
	// 202 — новый номер заказа принят в обработку;
	// 400 — неверный формат запроса;
	// 401 — пользователь не аутентифицирован;
	// 409 — номер заказа уже был загружен другим пользователем;
	// 500 — внутренняя ошибка сервера.

	// get input data
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Error(fmt.Errorf("%s %w", "Handler Order ShouldBindBodyWith", err))
		c.AbortWithError(http.StatusBadRequest, entity.ErrStatusBadRequest)
		return
	}
	order := entity.Order{
		Number: binary.BigEndian.Uint64(body),
	}

	// check input data
	if err := order.IsValidNumber(); err != nil {
		c.Error(fmt.Errorf("%s %w", "Handler Order IsValidNumber", err))
		c.AbortWithError(http.StatusUnprocessableEntity, entity.ErrUnprocessableEntity)
		return
	}

	c.Status(http.StatusOK)
}
