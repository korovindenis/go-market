package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/korovindenis/go-market/internal/domain/entity"
	"github.com/korovindenis/go-market/internal/port/http/handler/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type callTimes struct {
	userRegister  int
	generateToken int
	userLogin     int
	getUser       int
}

func Example() {
	t := &testing.T{}

	// run test with testify
	TestHandler_AuthRegister(t)
	TestHandler_AuthLogin(t)
}

func TestHandler_AuthRegister(t *testing.T) {
	config := mocks.NewConfig(t)
	usecase := mocks.NewUsecase(t)
	auth := mocks.NewAuth(t)
	ctxInf := mocks.NewCtxinfo(t)
	handler, _ := New(config, usecase, auth, ctxInf)
	router := gin.Default()

	config.On("GetTokenName").Return("gomarket_auth", nil)
	config.On("GetTokenLifeTime").Return(time.Duration(6), nil)
	router.POST("/register", handler.Register)

	tests := []struct {
		name       string
		args       entity.User
		statusCode int
		err        error
		callTimes  callTimes
	}{
		{
			name:       "register positive",
			args:       entity.User{Login: "user10", Password: "root"},
			statusCode: http.StatusOK,
			callTimes: callTimes{
				userRegister:  1,
				generateToken: 1,
			},
		},
		{
			name:       "register wrong param",
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			args, _ := json.Marshal(tt.args)
			userRegister := usecase.On("UserRegister", mock.Anything, tt.args).Return(int64(0), nil).Times(tt.callTimes.userRegister)
			generateToken := auth.On("GenerateToken", tt.args).Return("newToken", nil).Times(tt.callTimes.generateToken)

			// Act
			req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer([]byte(args)))
			if err != tt.err {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)

			// Unset
			userRegister.Unset()
			generateToken.Unset()
		})
	}
}

func TestHandler_AuthLogin(t *testing.T) {
	config := mocks.NewConfig(t)
	usecase := mocks.NewUsecase(t)
	auth := mocks.NewAuth(t)
	ctxInf := mocks.NewCtxinfo(t)
	handler, _ := New(config, usecase, auth, ctxInf)
	router := gin.Default()

	config.On("GetTokenName").Return("gomarket_auth", nil).Maybe()
	config.On("GetTokenLifeTime").Return(time.Duration(6), nil).Maybe()
	router.POST("/login", handler.Login)

	tests := []struct {
		name       string
		args       entity.User
		statusCode int
		err        error
		callTimes  callTimes
	}{
		{
			name:       "login positive",
			args:       entity.User{Login: "user10", Password: "root"},
			statusCode: http.StatusOK,
			callTimes: callTimes{
				userLogin:     1,
				generateToken: 1,
				getUser:       1,
			},
		},
		{
			name:       "login wrong param",
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			args, _ := json.Marshal(tt.args)
			getUser := usecase.On("GetUser", mock.Anything, tt.args).Return(tt.args, nil).Times(tt.callTimes.getUser)
			generateToken := auth.On("GenerateToken", mock.Anything).Return("newToken", nil).Times(tt.callTimes.generateToken)
			userLogin := usecase.On("UserLogin", mock.Anything, tt.args).Return(nil).Maybe().Times(tt.callTimes.userLogin)

			// Act
			req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(args))
			if err != tt.err {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)

			// Unset
			userLogin.Unset()
			generateToken.Unset()
			getUser.Unset()
		})
	}
}
