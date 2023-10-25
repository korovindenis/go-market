package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/korovindenis/go-market/internal/domain/entity"
)

type usecase interface {
	UserRegister(ctx context.Context, user entity.User) error
	UserLogin(ctx context.Context, user entity.User) error
	GetUser(ctx context.Context, userFromReq entity.User) (entity.User, error)

	GetOrder(ctx context.Context, user entity.User) ([]entity.Order, error)
	AddOrder(ctx context.Context, order entity.Order, user entity.User) error
}

type auth interface {
	GenerateToken(user entity.User) (string, error)
}

type config interface {
	GetTokenName() string
	GetTokenLifeTime() time.Duration
}

type Handler struct {
	usecase
	auth
	config
}

func New(config config, usecase usecase, auth auth) (*Handler, error) {
	return &Handler{
		usecase,
		auth,
		config,
	}, nil
}

func (h *Handler) getUserIdFromCtx(c *gin.Context) (uint64, error) {
	userIDRaw, ok := c.Get("userId")
	if !ok {
		return 0, fmt.Errorf("%s", "getUserIdFromCtx Get UserId")
	}
	userID, ok := userIDRaw.(uint64)
	if !ok {
		return 0, fmt.Errorf("%s", "getUserIdFromCtx userIDRaw")
	}

	return userID, nil
}
