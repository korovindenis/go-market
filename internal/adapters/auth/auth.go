package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/korovindenis/go-market/internal/domain/entity"
)

//go:generate mockery --name config --exported
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

func (a *Auth) GenerateToken(userFromBd entity.User) (string, error) {
	claims := claims{
		userFromBd,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.GetTokenLifeTime() * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(a.GetAppSecretKey()))
}

func (a *Auth) CheckToken(user entity.User, tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.GetAppSecretKey()), nil
	})

	if !token.Valid {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return fmt.Errorf("auth CheckToken - that's not even a token: %s", err)
		} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return fmt.Errorf("auth CheckToken - invalid signature: %s", err)
		} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
			return fmt.Errorf("auth CheckToken - timing is everything: %s", err)
		}
		return fmt.Errorf("auth CheckToken - couldn't handle this token: %s", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("error when extracting claims")
	}

	// compare ip,useragent in token and current
	if user.IP != claims["IP"] {
		return fmt.Errorf("ip address in the token does not match the current one")
	}

	if user.UserAgent != claims["UserAgent"] {
		return fmt.Errorf("useragent in the token does not match the current one")
	}

	return nil
}

// get user from token
func (a *Auth) GetUserFromToken(tokenString string) (entity.User, error) {
	var user entity.User

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.GetAppSecretKey()), nil
	})

	if token.Valid && err == nil {
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return user, fmt.Errorf("error when extracting claims: %s", err)
		}
		userIDFloat, found := claims["ID"].(float64)
		if !found {
			return user, fmt.Errorf("error when extracting claims Id: %s", err)
		}

		user.ID = int64(userIDFloat)
	}

	return user, err
}
