package handler

import (
	"context"
	"time"

	"github.com/korovindenis/go-market/internal/domain/entity"
)

type usecase interface {
	UserRegister(ctx context.Context, user entity.User) error
	UserLogin(ctx context.Context, user entity.User) error
}

type auth interface {
	GetNewToken(user entity.User) (string, error)
}

type config interface {
	GetAppName() string
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
