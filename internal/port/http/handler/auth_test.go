package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/korovindenis/go-market/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUsecase struct {
	mock.Mock
}

func (m *mockUsecase) UserRegister(ctx context.Context, user entity.User) (int64, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockUsecase) UserLogin(ctx context.Context, user entity.User) error {
	return m.Called(ctx, user).Error(0)
}
func (m *mockUsecase) AddOrder(ctx context.Context, order entity.Order, user entity.User) error {
	return m.Called(ctx, user).Error(0)
}

func (m *mockUsecase) GetAllOrders(ctx context.Context, user entity.User) ([]entity.Order, error) {
	return []entity.Order{}, nil
}
func (m *mockUsecase) GetBalance(ctx context.Context, user entity.User) (entity.Balance, error) {
	return entity.Balance{}, nil
}
func (m *mockUsecase) Withdrawals(ctx context.Context, user entity.User) ([]entity.BalanceUpdate, error) {
	return []entity.BalanceUpdate{}, nil
}
func (m *mockUsecase) GetUser(ctx context.Context, user entity.User) (entity.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *mockUsecase) GetUserBalance(ctx context.Context, user entity.User) (entity.Balance, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(entity.Balance), args.Error(1)
}

func (m *mockUsecase) WithdrawBalance(ctx context.Context, balance entity.BalanceUpdate, user entity.User) error {
	return m.Called(ctx, balance, user).Error(0)
}

func TestRegisterHandler(t *testing.T) {
	r := gin.Default()
	handler := &Handler{}
	mockUsecase := new(mockUsecase)
	handler.usecase = mockUsecase

	user := entity.User{}

	jsonUser, _ := json.Marshal(user)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonUser))
	req.Header.Set("Content-Type", "application/json")

	mockUsecase.On("UserRegister", mock.Anything, mock.Anything).Return(int64(1), nil)

	r.POST("/register", handler.Register)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLoginHandler(t *testing.T) {
	r := gin.Default()
	handler := &Handler{}
	mockUsecase := new(mockUsecase)
	handler.usecase = mockUsecase

	user := entity.User{}

	jsonUser, _ := json.Marshal(user)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonUser))
	req.Header.Set("Content-Type", "application/json")

	mockUsecase.On("GetUser", mock.Anything, mock.Anything).Return(user, nil)
	mockUsecase.On("UserLogin", mock.Anything, mock.Anything).Return(nil)

	r.POST("/login", handler.Login)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
