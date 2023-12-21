package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/korovindenis/go-market/internal/domain/entity"
	"github.com/korovindenis/go-market/internal/port/http/handler/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_GetBalance(t *testing.T) {
	config := mocks.NewConfig(t)
	usecase := mocks.NewUsecase(t)
	auth := mocks.NewAuth(t)
	ctxInf := mocks.NewCtxinfo(t)
	handler, _ := New(config, usecase, auth, ctxInf)
	router := gin.Default()
	router.GET("/balance", handler.GetBalance)

	tests := []struct {
		name       string
		balance    entity.Balance
		statusCode int
		err        error
		user       entity.User
	}{
		{
			name:       "get balance",
			balance:    entity.Balance{Current: 100, Withdrawn: 100},
			statusCode: http.StatusOK,
		},
		{
			name:       "get balance wrong param",
			statusCode: http.StatusInternalServerError,
			err:        errors.New("api was failed"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			assertData, _ := json.Marshal(tt.balance)
			w := httptest.NewRecorder()

			getUserIDFromCtx := ctxInf.On("GetUserIDFromCtx", mock.Anything).Return(int64(0), tt.err).Maybe()
			getBalance := usecase.On("GetBalance", mock.Anything, mock.Anything).Return(tt.balance, tt.err).Maybe()

			// Act
			req, err := http.NewRequest(http.MethodGet, "/balance", nil)
			if err != nil {
				t.Fatal(err)
			}

			router.ServeHTTP(w, req)

			// Assert
			if tt.statusCode == http.StatusOK {
				responseData, _ := io.ReadAll(w.Body)
				assert.Equal(t, string(assertData), string(responseData))
			}
			assert.Equal(t, tt.statusCode, w.Code)

			// Unset
			getUserIDFromCtx.Unset()
			getBalance.Unset()
		})
	}
}

func TestHandler_WithdrawBalance(t *testing.T) {
	config := mocks.NewConfig(t)
	usecase := mocks.NewUsecase(t)
	auth := mocks.NewAuth(t)
	ctxInf := mocks.NewCtxinfo(t)
	handler, _ := New(config, usecase, auth, ctxInf)
	router := gin.Default()
	router.POST("/balance/withdraw", handler.WithdrawBalance)

	tests := []struct {
		name          string
		statusCode    int
		err           error
		balanceUpdate entity.BalanceUpdate
		user          entity.User
	}{
		{
			name:          "withdraw balance - positive",
			statusCode:    http.StatusOK,
			balanceUpdate: entity.BalanceUpdate{Order: "2377225624", Sum: float64(100)},
		},
		{
			name:          "check Luhn",
			statusCode:    http.StatusUnprocessableEntity,
			balanceUpdate: entity.BalanceUpdate{Order: "1", Sum: float64(100)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			args, _ := json.Marshal(tt.balanceUpdate)
			w := httptest.NewRecorder()

			ctxInf.On("GetUserIDFromCtx", mock.Anything).Return(int64(0), tt.err)

			if tt.statusCode == http.StatusOK {
				usecase.On("WithdrawBalance", mock.Anything, tt.balanceUpdate, tt.user).Return(tt.err)
			}

			// Act
			req, err := http.NewRequest(http.MethodPost, "/balance/withdraw", bytes.NewBuffer([]byte(args)))
			if err != tt.err {
				t.Fatal(err)
			}

			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}
