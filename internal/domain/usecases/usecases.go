package usecases

import (
	"context"

	"github.com/korovindenis/go-market/internal/domain/entity"
)

type storage interface {
	UserRegister(ctx context.Context, user entity.User) error
	UserLogin(ctx context.Context, user entity.User) error
	GetUser(ctx context.Context, userFromReq entity.User) (entity.User, error)

	AddOrder(ctx context.Context, order entity.Order, user entity.User) error
}

type Usecases struct {
	storage storage
}

func New(storage storage) (*Usecases, error) {
	return &Usecases{
		storage,
	}, nil
}

// auth
func (u *Usecases) UserRegister(ctx context.Context, user entity.User) error {
	return u.storage.UserRegister(ctx, user)
}
func (u *Usecases) UserLogin(ctx context.Context, user entity.User) error {
	return u.storage.UserLogin(ctx, user)
}
func (u *Usecases) GetUser(ctx context.Context, userFromReq entity.User) (entity.User, error) {
	return u.storage.GetUser(ctx, userFromReq)
}

// orders
func (u *Usecases) AddOrder(ctx context.Context, order entity.Order, user entity.User) error {
	return u.storage.AddOrder(ctx, order, user)
}
