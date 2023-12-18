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

func TestHandler_Auth(t *testing.T) {
	config := mocks.NewConfig(t)
	usecase := mocks.NewUsecase(t)
	auth := mocks.NewAuth(t)
	ctxInf := mocks.NewCtxinfo(t)
	handler, _ := New(config, usecase, auth, ctxInf)
	router := gin.Default()

	config.On("GetTokenName").Return("gomarket_auth", nil)
	config.On("GetTokenLifeTime").Return(time.Duration(6), nil)

	tests := []struct {
		name       string
		route      string
		args       entity.User
		statusCode int
		err        error
	}{
		{
			name:       "register",
			route:      "/register",
			args:       entity.User{Login: "user10", Password: "root"},
			statusCode: http.StatusOK,
		},
		{
			name:       "login",
			route:      "/login",
			args:       entity.User{Login: "user10", Password: "root"},
			statusCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			args, _ := json.Marshal(tt.args)
			usecase.On("UserRegister", mock.Anything, tt.args).Return(int64(0), nil)
			auth.On("GenerateToken", tt.args).Return("newToken", nil)
			router.POST(tt.route, handler.Register)

			// Act
			req, err := http.NewRequest(http.MethodPost, tt.route, bytes.NewBuffer([]byte(args)))
			if err != tt.err {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}
