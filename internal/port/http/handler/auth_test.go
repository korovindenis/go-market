package handler

import (
	"bytes"
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

func TestHandler_Register(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		h    *Handler
		args args
	}{
		{
			name: "register",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := mocks.NewConfig(t)
			usecase := mocks.NewUsecase(t)
			auth := mocks.NewAuth(t)
			handler, _ := New(config, usecase, auth)
			router := gin.Default()

			usecase.On("UserRegister", mock.Anything, entity.User{Login: "user10", Password: "root"}).Return(int64(0), nil)
			auth.On("GenerateToken", entity.User{Login: "user10", Password: "root"}).Return("newToken", nil)
			config.On("GetTokenName").Return("gomarket_auth", nil)
			config.On("GetTokenLifeTime").Return(time.Duration(6), nil)

			// Привязываем хандлер к маршруту
			router.POST("/register", handler.Register)

			// Создаем тестовый запрос
			requestBody := []byte("{\"login\": \"user10\",\"password\": \"root\"}") // Пример данных формы
			req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(requestBody))
			if err != nil {
				t.Fatal(err)
			}

			// Создаем тестовый ответ
			w := httptest.NewRecorder()

			// Выполняем запрос к маршруту
			router.ServeHTTP(w, req)

			// Проверяем результат
			assert.Equal(t, http.StatusOK, w.Code)
			//assert.Contains(t, w.Body.String(), "User registered: John")

			//tt.h.Register(tt.args.c)
		})
	}
}
