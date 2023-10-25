package usecases

import (
	"context"
	"fmt"

	"github.com/korovindenis/go-market/internal/domain/entity"
	"golang.org/x/crypto/bcrypt"
)

type storage interface {
	UserRegister(ctx context.Context, user entity.User) error
	UserLogin(ctx context.Context, user entity.User) error
	GetUser(ctx context.Context, userFromReq entity.User) (entity.User, error)

	GetOrder(ctx context.Context, user entity.User) ([]entity.Order, error)
	AddOrder(ctx context.Context, order entity.Order, user entity.User) error
}

type config interface {
	GetStorageSalt() string
}

type Usecases struct {
	storage
	config
}

func New(config config, storage storage) (*Usecases, error) {
	return &Usecases{
		storage,
		config,
	}, nil
}

// auth
func (u *Usecases) UserRegister(ctx context.Context, user entity.User) error {
	// add salt to password
	password, err := hashPassword(user.Password, u.config.GetStorageSalt())
	if err != nil {
		return err
	}
	user.Password = password

	return u.storage.UserRegister(ctx, user)
}
func (u *Usecases) UserLogin(ctx context.Context, user entity.User) error {
	// add salt to password
	user.Password = fmt.Sprintf("%s%s", user.Password, u.config.GetStorageSalt())

	return u.storage.UserLogin(ctx, user)
}
func (u *Usecases) GetUser(ctx context.Context, userFromReq entity.User) (entity.User, error) {
	return u.storage.GetUser(ctx, userFromReq)
}

// orders
func (u *Usecases) AddOrder(ctx context.Context, order entity.Order, user entity.User) error {
	return u.storage.AddOrder(ctx, order, user)
}

func (u *Usecases) GetOrder(ctx context.Context, user entity.User) ([]entity.Order, error) {
	return u.storage.GetOrder(ctx, user)
}

func hashPassword(password, salt string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("%s%s", password, salt)), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
