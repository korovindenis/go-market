package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/korovindenis/go-market/internal/domain/entity"
)

type config interface {
	GetTokenName() string
}

type auth interface {
	CheckToken(user entity.User, tokenString string) error
	GetUserFromToken(tokenString string) (entity.User, error)
}

type Middleware struct {
	config
	auth
}

func New(config config, auth auth) (*Middleware, error) {
	return &Middleware{
		config,
		auth,
	}, nil
}

func (m *Middleware) CheckMethod() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodPost && c.Request.Method != http.MethodGet {
			c.AbortWithError(http.StatusMethodNotAllowed, entity.ErrMethodNotAllowed)
			return
		}
		c.Next()
	}
}

func (m *Middleware) CheckContentTypeJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet {
			c.Next()
			return
		}
		if c.GetHeader("Content-Type") != "application/json" {
			c.Error(fmt.Errorf("%s", "Middleware CheckContentTypeJSON"))
			c.AbortWithError(http.StatusBadRequest, entity.ErrStatusBadRequest)
			return
		}
		c.Next()
	}
}

func (m *Middleware) CheckContentTypeText() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet {
			c.Next()
			return
		}
		if c.GetHeader("Content-Type") != "text/plain" {
			c.Error(fmt.Errorf("%s", "Middleware CheckContentTypeText"))
			c.AbortWithError(http.StatusBadRequest, entity.ErrStatusBadRequest)
			return
		}
		c.Next()
	}
}

func (m *Middleware) CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(m.config.GetTokenName())
		if err != nil {
			c.Error(fmt.Errorf("%s %w", "Middleware CheckAuth get Cookie", err))
			c.AbortWithError(http.StatusUnauthorized, entity.ErrUserLoginUnauthorized)
			return
		}

		userDevice := entity.User{
			IP:        c.ClientIP(),
			UserAgent: c.GetHeader("User-Agent"),
		}

		if err := m.auth.CheckToken(userDevice, token); err != nil {
			c.Error(fmt.Errorf("error: %s %w", "Middleware CheckToken", err))
			c.AbortWithError(http.StatusUnauthorized, entity.ErrUserLoginUnauthorized)
			return
		}

		c.Next()
	}
}

func (m *Middleware) AddUserInfoToCtx() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(m.config.GetTokenName())
		if err != nil {
			c.Error(fmt.Errorf("%s %w", "Middleware AddUserInfoToCtx get Cookie", err))
			c.AbortWithError(http.StatusInternalServerError, entity.ErrInternalServerError)
			return
		}

		user, err := m.auth.GetUserFromToken(token)
		if err != nil {
			c.Error(fmt.Errorf("error: %s %w", "Middleware AddUserInfoToCtx GetUserFromToken", err))
			c.AbortWithError(http.StatusInternalServerError, entity.ErrInternalServerError)
			return
		}

		c.Set("userId", user.ID)

		c.Next()
	}
}
