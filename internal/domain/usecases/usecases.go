package usecases

import (
	"context"

	"github.com/korovindenis/go-market/internal/domain/entity"
)

type storage interface {
	UserRegister(ctx context.Context, user entity.User) error
	UserLogin(ctx context.Context, user entity.User) error
}

type Usecases struct {
	storage storage
}

func New(storage storage) (*Usecases, error) {
	return &Usecases{
		storage,
	}, nil
}

func (u *Usecases) UserRegister(ctx context.Context, user entity.User) error {
	return u.storage.UserRegister(ctx, user)
}

func (u *Usecases) UserLogin(ctx context.Context, user entity.User) error {
	return u.storage.UserLogin(ctx, user)
}
