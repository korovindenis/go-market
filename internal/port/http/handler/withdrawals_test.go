package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/korovindenis/go-market/internal/domain/entity"
	"github.com/korovindenis/go-market/internal/port/http/handler/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_Withdrawals(t *testing.T) {
	config := mocks.NewConfig(t)
	usecase := mocks.NewUsecase(t)
	auth := mocks.NewAuth(t)
	ctxInf := mocks.NewCtxinfo(t)
	handler, _ := New(config, usecase, auth, ctxInf)
	router := gin.Default()
	router.GET("/withdrawals", handler.Withdrawals)

	type CustErr struct {
		withdrawals    error
		getUserFromCtx error
		request        error
	}
	tests := []struct {
		name       string
		orders     []entity.BalanceUpdate
		statusCode int
		err        CustErr
	}{
		{
			name:       "withdrawals positive",
			statusCode: http.StatusOK,
		},
		{
			name:       "withdrawals err Get withdrawals",
			statusCode: http.StatusInternalServerError,
			err: CustErr{
				withdrawals: errors.New("get orwithdrawals was failed"),
			},
		},
		{
			name:       "withdrawals err GetUserIDFromCtx",
			statusCode: http.StatusInternalServerError,
			err: CustErr{
				getUserFromCtx: errors.New("get withdrawals was failed"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			Withdrawals := usecase.On("Withdrawals", mock.Anything, mock.Anything).Return(tt.orders, tt.err.withdrawals).Maybe()
			GetUserIDFromCtx := ctxInf.On("GetUserIDFromCtx", mock.Anything).Return(int64(0), tt.err.getUserFromCtx)

			// Act
			req, err := http.NewRequest(http.MethodGet, "/withdrawals", nil)
			if err != tt.err.request {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)

			// Unset
			GetUserIDFromCtx.Unset()
			Withdrawals.Unset()
		})
	}
}
