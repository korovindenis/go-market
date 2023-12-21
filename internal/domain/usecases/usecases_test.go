package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/korovindenis/go-market/internal/domain/entity"
	"github.com/korovindenis/go-market/internal/domain/usecases/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUsecases_UserRegister(t *testing.T) {
	config := mocks.NewConfig(t)
	storage := mocks.NewStorage(t)
	usecases, _ := New(config, storage)
	config.On("GetStorageSalt").Return("xxxxxxxx", nil).Maybe()

	tests := []struct {
		ctx  context.Context
		name string
		want int64
		err  error
		u    *Usecases
		user entity.User
	}{
		{
			name: "positive",
			want: 0,
			u:    usecases,
			ctx:  context.Background(),
		},
		{
			name: "negative",
			want: 0,
			u:    usecases,
			ctx:  context.Background(),
			err:  errors.New(""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			userRegister := storage.On("UserRegister", mock.Anything, mock.Anything).Return(int64(0), tt.err)

			// Act
			user, err := tt.u.UserRegister(tt.ctx, tt.user)

			// Assert
			if err != nil && !errors.Is(err, tt.err) {
				t.Fatal(err)
			}
			assert.Equal(t, tt.user.ID, user)

			// Unset
			userRegister.Unset()
		})
	}
}
func TestUsecases_UserLogin(t *testing.T) {
	config := mocks.NewConfig(t)
	storage := mocks.NewStorage(t)
	usecases, _ := New(config, storage)
	config.On("GetStorageSalt").Return("xxxxxxxx", nil).Maybe()

	tests := []struct {
		ctx  context.Context
		name string
		want int64
		err  error
		u    *Usecases
		user entity.User
	}{
		{
			name: "positive",
			want: 0,
			u:    usecases,
			ctx:  context.Background(),
		},
		{
			name: "negative",
			want: 0,
			u:    usecases,
			ctx:  context.Background(),
			err:  errors.New(""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			userLogin := storage.On("UserLogin", mock.Anything, mock.Anything).Return(tt.err)

			// Act
			err := tt.u.UserLogin(tt.ctx, tt.user)

			// Assert
			if err != nil && !errors.Is(err, tt.err) {
				t.Fatal(err)
			}

			// Unset
			userLogin.Unset()
		})
	}
}
func TestUsecases_GetUser(t *testing.T) {
	config := mocks.NewConfig(t)
	storage := mocks.NewStorage(t)
	usecases, _ := New(config, storage)
	config.On("GetStorageSalt").Return("xxxxxxxx", nil).Maybe()

	tests := []struct {
		ctx  context.Context
		name string
		want int64
		err  error
		u    *Usecases
		user entity.User
	}{
		{
			name: "positive",
			want: 0,
			u:    usecases,
			ctx:  context.Background(),
		},
		{
			name: "negative",
			want: 0,
			u:    usecases,
			ctx:  context.Background(),
			err:  errors.New(""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			getUser := storage.On("GetUser", mock.Anything, mock.Anything).Return(entity.User{}, tt.err)

			// Act
			user, err := tt.u.GetUser(tt.ctx, tt.user)

			// Assert
			if err != nil && !errors.Is(err, tt.err) {
				t.Fatal(err)
			}
			assert.Equal(t, tt.user, user)

			// Unset
			getUser.Unset()
		})
	}
}
func TestUsecases_AddOrder(t *testing.T) {
	config := mocks.NewConfig(t)
	storage := mocks.NewStorage(t)
	usecases, _ := New(config, storage)

	tests := []struct {
		ctx   context.Context
		name  string
		want  int64
		err   error
		u     *Usecases
		user  entity.User
		order entity.Order
	}{
		{
			name: "positive",
			want: 0,
			u:    usecases,
			ctx:  context.Background(),
		},
		{
			name: "negative",
			want: 0,
			u:    usecases,
			ctx:  context.Background(),
			err:  errors.New(""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			addOrder := storage.On("AddOrder", mock.Anything, mock.Anything, mock.Anything).Return(tt.err)

			// Act
			err := tt.u.AddOrder(tt.ctx, tt.order, tt.user)

			// Assert
			if err != nil && !errors.Is(err, tt.err) {
				t.Fatal(err)
			}

			// Unset
			addOrder.Unset()
		})
	}
}
func TestUsecases_GetAllOrders(t *testing.T) {
	config := mocks.NewConfig(t)
	storage := mocks.NewStorage(t)
	usecases, _ := New(config, storage)

	tests := []struct {
		ctx   context.Context
		name  string
		want  int64
		err   error
		u     *Usecases
		user  entity.User
		order []entity.Order
	}{
		{
			name:  "positive",
			want:  0,
			u:     usecases,
			ctx:   context.Background(),
			order: []entity.Order{},
		},
		{
			name:  "negative",
			want:  0,
			u:     usecases,
			ctx:   context.Background(),
			err:   errors.New(""),
			order: []entity.Order{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			getAllOrders := storage.On("GetAllOrders", mock.Anything, mock.Anything).Return([]entity.Order{}, tt.err)

			// Act
			order, err := tt.u.GetAllOrders(tt.ctx, tt.user)

			// Assert
			if err != nil && !errors.Is(err, tt.err) {
				t.Fatal(err)
			}
			assert.Equal(t, tt.order, order)

			// Unset
			getAllOrders.Unset()
		})
	}
}

func TestUsecases_GetBalance(t *testing.T) {
	config := mocks.NewConfig(t)
	storage := mocks.NewStorage(t)
	usecases, _ := New(config, storage)

	tests := []struct {
		ctx     context.Context
		name    string
		want    int64
		err     error
		u       *Usecases
		user    entity.User
		balance entity.Balance
	}{
		{
			name:    "positive",
			want:    0,
			u:       usecases,
			ctx:     context.Background(),
			balance: entity.Balance{},
		},
		{
			name:    "negative",
			want:    0,
			u:       usecases,
			ctx:     context.Background(),
			err:     errors.New(""),
			balance: entity.Balance{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			getBalance := storage.On("GetBalance", mock.Anything, mock.Anything).Return(entity.Balance{}, tt.err)

			// Act
			balance, err := tt.u.GetBalance(tt.ctx, tt.user)

			// Assert
			if err != nil && !errors.Is(err, tt.err) {
				t.Fatal(err)
			}
			assert.Equal(t, tt.balance, balance)

			// Unset
			getBalance.Unset()
		})
	}
}
func TestUsecases_WithdrawBalance(t *testing.T) {
	config := mocks.NewConfig(t)
	storage := mocks.NewStorage(t)
	usecases, _ := New(config, storage)

	tests := []struct {
		ctx     context.Context
		name    string
		want    int64
		err     error
		u       *Usecases
		user    entity.User
		balance entity.BalanceUpdate
	}{
		{
			name:    "positive",
			want:    0,
			u:       usecases,
			ctx:     context.Background(),
			balance: entity.BalanceUpdate{},
		},
		{
			name:    "negative",
			want:    0,
			u:       usecases,
			ctx:     context.Background(),
			err:     errors.New(""),
			balance: entity.BalanceUpdate{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			withdrawBalance := storage.On("WithdrawBalance", mock.Anything, mock.Anything, mock.Anything).Return(tt.err)

			// Act
			err := tt.u.WithdrawBalance(tt.ctx, tt.balance, tt.user)

			// Assert
			if err != nil && !errors.Is(err, tt.err) {
				t.Fatal(err)
			}

			// Unset
			withdrawBalance.Unset()
		})
	}
}
func TestUsecases_Withdrawals(t *testing.T) {
	config := mocks.NewConfig(t)
	storage := mocks.NewStorage(t)
	usecases, _ := New(config, storage)

	tests := []struct {
		ctx     context.Context
		name    string
		want    int64
		err     error
		u       *Usecases
		user    entity.User
		balance []entity.BalanceUpdate
	}{
		{
			name:    "positive",
			want:    0,
			u:       usecases,
			ctx:     context.Background(),
			balance: []entity.BalanceUpdate{},
		},
		{
			name:    "negative",
			want:    0,
			u:       usecases,
			ctx:     context.Background(),
			err:     errors.New(""),
			balance: []entity.BalanceUpdate{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			withdrawBalance := storage.On("Withdrawals", mock.Anything, mock.Anything).Return([]entity.BalanceUpdate{}, tt.err)

			// Act
			balance, err := tt.u.Withdrawals(tt.ctx, tt.user)

			// Assert
			if err != nil && !errors.Is(err, tt.err) {
				t.Fatal(err)
			}
			assert.Equal(t, balance, tt.balance)

			// Unset
			withdrawBalance.Unset()
		})
	}
}
