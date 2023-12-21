package handler

import (
	"bytes"
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

func TestHandler_SetOrder(t *testing.T) {
	config := mocks.NewConfig(t)
	usecase := mocks.NewUsecase(t)
	auth := mocks.NewAuth(t)
	ctxInf := mocks.NewCtxinfo(t)
	handler, _ := New(config, usecase, auth, ctxInf)
	router := gin.Default()
	router.POST("/orders", handler.SetOrder)

	type CustErr struct {
		addOrder       error
		getUserFromCtx error
		request        error
	}
	tests := []struct {
		name       string
		args       string
		statusCode int
		err        CustErr
	}{
		{
			name:       "orders positive",
			args:       "9278923470",
			statusCode: http.StatusAccepted,
		},
		{
			name:       "orders check Luhn",
			args:       "1",
			statusCode: http.StatusUnprocessableEntity,
		},
		{
			name:       "orders err AddOrder",
			args:       "9278923470",
			statusCode: http.StatusInternalServerError,
			err: CustErr{
				addOrder: errors.New("add order was failed"),
			},
		},
		{
			name:       "orders err GetUserIDFromCtx",
			args:       "9278923470",
			statusCode: http.StatusInternalServerError,
			err: CustErr{
				getUserFromCtx: errors.New("get id was failed"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			AddOrder := usecase.On("AddOrder", mock.Anything, mock.Anything, mock.Anything).Return(tt.err.addOrder).Maybe()
			GetUserIDFromCtx := ctxInf.On("GetUserIDFromCtx", mock.Anything).Return(int64(0), tt.err.getUserFromCtx)

			// Act
			req, err := http.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer([]byte(tt.args)))
			if err != tt.err.request {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)

			// Unset
			GetUserIDFromCtx.Unset()
			AddOrder.Unset()
		})
	}
}

func TestHandler_GetAllOrders(t *testing.T) {
	config := mocks.NewConfig(t)
	usecase := mocks.NewUsecase(t)
	auth := mocks.NewAuth(t)
	ctxInf := mocks.NewCtxinfo(t)
	handler, _ := New(config, usecase, auth, ctxInf)
	router := gin.Default()
	router.GET("/orders", handler.GetAllOrders)

	type CustErr struct {
		getOrder       error
		getUserFromCtx error
		request        error
	}
	tests := []struct {
		name       string
		orders     []entity.Order
		statusCode int
		err        CustErr
	}{
		{
			name:       "orders positive",
			statusCode: http.StatusOK,
		},
		{
			name:       "orders err GetAllOrders",
			statusCode: http.StatusInternalServerError,
			err: CustErr{
				getOrder: errors.New("get order was failed"),
			},
		},
		{
			name:       "orders err GetUserIDFromCtx",
			statusCode: http.StatusInternalServerError,
			err: CustErr{
				getUserFromCtx: errors.New("get id was failed"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			GetAllOrders := usecase.On("GetAllOrders", mock.Anything, mock.Anything, mock.Anything).Return(tt.orders, tt.err.getOrder).Maybe()
			GetUserIDFromCtx := ctxInf.On("GetUserIDFromCtx", mock.Anything).Return(int64(0), tt.err.getUserFromCtx)

			// Act
			req, err := http.NewRequest(http.MethodGet, "/orders", nil)
			if err != tt.err.request {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)

			// Unset
			GetUserIDFromCtx.Unset()
			GetAllOrders.Unset()
		})
	}
}
