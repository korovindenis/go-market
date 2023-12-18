package handler

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/korovindenis/go-market/internal/domain/entity"
)

//go:generate mockery --name usecase --exported
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

//go:generate mockery --name auth --exported
type auth interface {
	GenerateToken(user entity.User) (string, error)
}

//go:generate mockery --name config --exported
type config interface {
	GetTokenName() string
	GetTokenLifeTime() time.Duration
}

//go:generate mockery --name ctxinfo --exported
type ctxinfo interface {
	GetUserIDFromCtx(c *gin.Context) (int64, error)
}

type Handler struct {
	usecase
	auth
	config
	ctxinfo
}

func New(config config, usecase usecase, auth auth, ctxinfo ctxinfo) (*Handler, error) {
	return &Handler{
		usecase,
		auth,
		config,
		ctxinfo,
	}, nil
}
