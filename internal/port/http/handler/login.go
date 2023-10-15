package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/korovindenis/go-market/internal/domain/entity"
)

func (h *Handler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	user := entity.User{
		Ip:        c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
	}

	// check input data
	if err := c.ShouldBindJSON(&user); err != nil {
		c.Error(fmt.Errorf("%s %w", "Handler Login ShouldBindJSON", err))
		c.AbortWithError(http.StatusBadRequest, entity.ErrStatusBadRequest)
		return
	}

	// generation token
	token, err := h.auth.GetNewToken(user)
	if err != nil {
		c.Error(fmt.Errorf("%s %w", "Handler Login GetNewToken", err))
		c.AbortWithError(http.StatusInternalServerError, entity.ErrInternalServerError)
		return
	}

	// attempt auth user
	if err := h.usecase.UserLogin(ctx, user); err != nil {
		c.Error(fmt.Errorf("%s %w", "Handler Login UserLogin", err))

		if errors.Is(err, entity.ErrUserLoginUnauthorized) {
			c.AbortWithError(http.StatusUnauthorized, entity.ErrUserLoginUnauthorized)
			return
		}
		c.AbortWithError(http.StatusInternalServerError, entity.ErrInternalServerError)
		return
	}

	// create and set cookie
	cookie, err := h.createCookie(token)
	if err != nil {
		c.Error(fmt.Errorf("%s %w", "Handler Login createCookie", err))
		c.AbortWithError(http.StatusInternalServerError, entity.ErrInternalServerError)
		return
	}
	http.SetCookie(c.Writer, cookie)

	c.Status(http.StatusOK)
}
