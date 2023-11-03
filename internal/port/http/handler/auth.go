package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/korovindenis/go-market/internal/domain/entity"
)

func (h *Handler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	var user entity.User

	// check input data
	if err := c.ShouldBindJSON(&user); err != nil {
		c.Error(fmt.Errorf("%s %w", "Handler Register ShouldBindJSON", err))
		c.AbortWithError(http.StatusBadRequest, entity.ErrStatusBadRequest)
		return
	}

	// get user device info
	user.IP = c.ClientIP()
	user.UserAgent = c.GetHeader("User-Agent")

	// attempt registration user
	// with check unique login
	userID, err := h.usecase.UserRegister(ctx, user)
	if err != nil {
		c.Error(fmt.Errorf("%s %w", "Handler Register UserRegister", err))

		if errors.Is(err, entity.ErrUserLoginNotUnique) {
			c.AbortWithError(http.StatusConflict, entity.ErrUserLoginNotUnique)
			return
		}
		c.AbortWithError(http.StatusInternalServerError, entity.ErrInternalServerError)
		return
	}
	user.ID = userID

	// generation token
	token, err := h.auth.GenerateToken(user)
	if err != nil {
		c.Error(fmt.Errorf("%s %w", "Handler Register GenerateToken", err))
		c.AbortWithError(http.StatusInternalServerError, entity.ErrInternalServerError)
		return
	}

	// create and set cookie
	cookie, err := h.createCookie(token)
	if err != nil {
		c.Error(fmt.Errorf("%s %w", "Handler Register createCookie", err))
		c.AbortWithError(http.StatusInternalServerError, entity.ErrInternalServerError)
		return
	}
	http.SetCookie(c.Writer, cookie)

	c.Status(http.StatusOK)
}

func (h *Handler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var userFromReq entity.User

	// check input data
	if err := c.ShouldBindJSON(&userFromReq); err != nil {
		c.Error(fmt.Errorf("%s %w", "Handler Login ShouldBindJSON", err))
		c.AbortWithError(http.StatusBadRequest, entity.ErrStatusBadRequest)
		return
	}

	// get user device info
	user := entity.User{
		IP:        c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
	}

	// get user from storage
	userFromStorage, err := h.usecase.GetUser(ctx, userFromReq)
	if err != nil {
		c.Error(fmt.Errorf("%s %w", "Handler Register GetUser", err))
		c.AbortWithError(http.StatusInternalServerError, entity.ErrInternalServerError)
		return
	}
	user.ID = userFromStorage.ID

	// generation token
	token, err := h.auth.GenerateToken(user)
	if err != nil {
		c.Error(fmt.Errorf("%s %w", "Handler Login GenerateToken", err))
		c.AbortWithError(http.StatusInternalServerError, entity.ErrInternalServerError)
		return
	}

	// attempt auth user
	if err := h.usecase.UserLogin(ctx, userFromReq); err != nil {
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

func (h *Handler) createCookie(token string) (*http.Cookie, error) {
	return &http.Cookie{
		Name:     h.GetTokenName(),
		Value:    token,
		Expires:  time.Now().Add(h.GetTokenLifeTime() * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   false,
		Path:     "/",
	}, nil
}
