package auth

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/korovindenis/go-market/internal/adapters/auth/mocks"
	"github.com/korovindenis/go-market/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestAuth_GenerateToken(t *testing.T) {
	config := mocks.NewConfig(t)
	auth, _ := New(config)
	config.On("GetAppSecretKey").Return("xxxxxxxx", nil).Maybe()
	config.On("GetTokenName").Return("gomarket_auth", nil).Maybe()
	config.On("GetTokenLifeTime").Return(time.Duration(6), nil).Maybe()

	type args struct {
		userFromBd entity.User
	}
	tests := []struct {
		name string
		a    *Auth
		args args
		want string
		err  error
	}{
		{
			name: "positive",
			a:    auth,
			args: args{
				userFromBd: entity.User{Login: "root", Password: "root"},
			},
			want: "[a-zA-Z]+",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			token, err := tt.a.GenerateToken(tt.args.userFromBd)

			// Assert
			assert.ErrorIs(t, err, tt.err)
			if match, _ := regexp.MatchString(tt.want, token); !match {
				t.Fatal(err)
			}
		})
	}
}

func TestAuth_CheckToken(t *testing.T) {
	config := mocks.NewConfig(t)
	auth, _ := New(config)
	config.On("GetAppSecretKey").Return("xxxxxxxx", nil).Maybe()
	config.On("GetTokenName").Return("gomarket_auth", nil).Maybe()
	config.On("GetTokenLifeTime").Return(time.Duration(6), nil).Maybe()

	type args struct {
		user        entity.User
		tokenString string
	}
	tests := []struct {
		name string
		a    *Auth
		args args
		err  error
	}{
		{
			name: "positive",
			a:    auth,
			args: args{
				user:        entity.User{Login: "user7", Password: "root", IP: "127.0.0.1", UserAgent: "PostmanRuntime/7.29.2"},
				tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dpbiI6IiIsInBhc3N3b3JkIjoiIiwiSUQiOjEsIklQIjoiMTI3LjAuMC4xIiwiVXNlckFnZW50IjoiUG9zdG1hblJ1bnRpbWUvNy4yOS4yIiwiZXhwIjoxNzAzMDE2NjM2LCJuYmYiOjE3MDI5OTUwMzYsImlhdCI6MTcwMjk5NTAzNn0.9UumpeVZMCEXvbKfl1pLdCrBYRrlWs-55phxmEp1LQI",
			},
		},
		{
			name: "negative - wrong token",
			a:    auth,
			args: args{
				user:        entity.User{Login: "user7", Password: "root", IP: "127.0.0.1", UserAgent: "PostmanRuntime/7.29.2"},
				tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.xyJsb2dpbiI6IiIsInBhc3N3b3JkIjoiIiwiSUQiOjEsIklQIjoiMTI3LjAuMC4xIiwiVXNlckFnZW50IjoiUG9zdG1hblJ1bnRpbWUvNy4yOS4yIiwiZXhwIjoxNzAzMDE2NjM2LCJuYmYiOjE3MDI5OTUwMzYsImlhdCI6MTcwMjk5NTAzNn0.9UumpeVZMCEXvbKfl1pLdCrBYRrlWs-55phxmEp1LQI",
			},
			err: errors.New("auth CheckToken - that's not even a token: token is malformed: could not JSON decode claim: invalid character 'Ç' looking for beginning of value"),
		},
		{
			name: "negative - wrong IP",
			a:    auth,
			args: args{
				user:        entity.User{Login: "user7", Password: "root", IP: "0.0.0.0", UserAgent: "PostmanRuntime/7.29.2"},
				tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dpbiI6IiIsInBhc3N3b3JkIjoiIiwiSUQiOjEsIklQIjoiMTI3LjAuMC4xIiwiVXNlckFnZW50IjoiUG9zdG1hblJ1bnRpbWUvNy4yOS4yIiwiZXhwIjoxNzAzMDE2NjM2LCJuYmYiOjE3MDI5OTUwMzYsImlhdCI6MTcwMjk5NTAzNn0.9UumpeVZMCEXvbKfl1pLdCrBYRrlWs-55phxmEp1LQI",
			},
			err: errors.New("ip address in the token does not match the current one"),
		},
		{
			name: "negative - wrong UA",
			a:    auth,
			args: args{
				user:        entity.User{Login: "user7", Password: "root", IP: "127.0.0.1", UserAgent: "chrome 0.1"},
				tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dpbiI6IiIsInBhc3N3b3JkIjoiIiwiSUQiOjEsIklQIjoiMTI3LjAuMC4xIiwiVXNlckFnZW50IjoiUG9zdG1hblJ1bnRpbWUvNy4yOS4yIiwiZXhwIjoxNzAzMDE2NjM2LCJuYmYiOjE3MDI5OTUwMzYsImlhdCI6MTcwMjk5NTAzNn0.9UumpeVZMCEXvbKfl1pLdCrBYRrlWs-55phxmEp1LQI",
			},
			err: errors.New("ip address in the token does not match the current one"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := tt.a.CheckToken(tt.args.user, tt.args.tokenString)

			// Assert
			if err != nil && errors.Is(err, tt.err) {
				t.Fatal(err)
			}
		})
	}
}

func TestAuth_GetUserFromToken(t *testing.T) {
	config := mocks.NewConfig(t)
	auth, _ := New(config)
	config.On("GetAppSecretKey").Return("xxxxxxxx", nil).Maybe()
	config.On("GetTokenName").Return("gomarket_auth", nil).Maybe()
	config.On("GetTokenLifeTime").Return(time.Duration(6), nil).Maybe()

	type args struct {
		tokenString string
	}
	tests := []struct {
		name string
		a    *Auth
		args args
		err  error
		user entity.User
	}{
		{
			name: "positive",
			a:    auth,
			args: args{
				tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dpbiI6IiIsInBhc3N3b3JkIjoiIiwiSUQiOjEsIklQIjoiMTI3LjAuMC4xIiwiVXNlckFnZW50IjoiUG9zdG1hblJ1bnRpbWUvNy4yOS4yIiwiZXhwIjoxNzAzMDE2NjM2LCJuYmYiOjE3MDI5OTUwMzYsImlhdCI6MTcwMjk5NTAzNn0.9UumpeVZMCEXvbKfl1pLdCrBYRrlWs-55phxmEp1LQI",
			},
			user: entity.User{ID: 1},
		},
		{
			name: "negative - wrong token",
			a:    auth,
			args: args{
				tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.xyJsb2dpbiI6IiIsInBhc3N3b3JkIjoiIiwiSUQiOjEsIklQIjoiMTI3LjAuMC4xIiwiVXNlckFnZW50IjoiUG9zdG1hblJ1bnRpbWUvNy4yOS4yIiwiZXhwIjoxNzAzMDE2NjM2LCJuYmYiOjE3MDI5OTUwMzYsImlhdCI6MTcwMjk5NTAzNn0.9UumpeVZMCEXvbKfl1pLdCrBYRrlWs-55phxmEp1LQI",
			},
			err: errors.New("auth CheckToken - that's not even a token: token is malformed: could not JSON decode claim: invalid character 'Ç' looking for beginning of value"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			user, err := tt.a.GetUserFromToken(tt.args.tokenString)

			// Assert
			if err != nil && errors.Is(err, tt.err) {
				t.Fatal(err)
			}
			assert.Equal(t, tt.user.ID, user.ID)
		})
	}
}
