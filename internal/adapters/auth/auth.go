package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/korovindenis/go-market/internal/domain/entity"
)

type config interface {
	GetAppSecretKey() string
	GetTokenLifeTime() time.Duration
}

type Auth struct {
	config
}

type claims struct {
	entity.User
	jwt.RegisteredClaims
}

func New(config config) (*Auth, error) {
	return &Auth{
		config,
	}, nil
}

func (a *Auth) GetNewToken(user entity.User) (string, error) {
	claims := claims{
		user,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.GetTokenLifeTime() * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(a.GetAppSecretKey()))
}
