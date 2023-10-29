package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/korovindenis/go-market/internal/domain/entity"
)

type usecase interface {
	UserRegister(ctx context.Context, user entity.User) (int64, error)
	UserLogin(ctx context.Context, user entity.User) error
	GetUser(ctx context.Context, userFromReq entity.User) (entity.User, error)

	GetAllOrders(ctx context.Context, user entity.User) ([]entity.Order, error)
	AddOrder(ctx context.Context, order entity.Order, user entity.User) error

	GetBalance(ctx context.Context, user entity.User) (entity.Balance, error)
	WithdrawBalance(ctx context.Context, balance entity.BalanceUpdate, user entity.User) error

	Withdrawals(ctx context.Context, user entity.User) ([]entity.BalanceUpdate, error)
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

func (h *Handler) getUserIDFromCtx(c *gin.Context) (int64, error) {
	userIDRaw, ok := c.Get("userId")
	if !ok {
		return 0, fmt.Errorf("%s", "getUserIDFromCtx Get UserId")
	}
	userID, ok := userIDRaw.(int64)
	if !ok {
		return 0, fmt.Errorf("%s", "getUserIDFromCtx userIDRaw")
	}

	return userID, nil
}
