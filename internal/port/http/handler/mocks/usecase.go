// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	context "context"

	entity "github.com/korovindenis/go-market/internal/domain/entity"

	mock "github.com/stretchr/testify/mock"
)

// Usecase is an autogenerated mock type for the usecase type
type Usecase struct {
	mock.Mock
}

// AddOrder provides a mock function with given fields: ctx, order, user
func (_m *Usecase) AddOrder(ctx context.Context, order entity.Order, user entity.User) error {
	ret := _m.Called(ctx, order, user)

	if len(ret) == 0 {
		panic("no return value specified for AddOrder")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.Order, entity.User) error); ok {
		r0 = rf(ctx, order, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllOrders provides a mock function with given fields: ctx, user
func (_m *Usecase) GetAllOrders(ctx context.Context, user entity.User) ([]entity.Order, error) {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for GetAllOrders")
	}

	var r0 []entity.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.User) ([]entity.Order, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, entity.User) []entity.Order); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, entity.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBalance provides a mock function with given fields: ctx, user
func (_m *Usecase) GetBalance(ctx context.Context, user entity.User) (entity.Balance, error) {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for GetBalance")
	}

	var r0 entity.Balance
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.User) (entity.Balance, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, entity.User) entity.Balance); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Get(0).(entity.Balance)
	}

	if rf, ok := ret.Get(1).(func(context.Context, entity.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUser provides a mock function with given fields: ctx, userFromReq
func (_m *Usecase) GetUser(ctx context.Context, userFromReq entity.User) (entity.User, error) {
	ret := _m.Called(ctx, userFromReq)

	if len(ret) == 0 {
		panic("no return value specified for GetUser")
	}

	var r0 entity.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.User) (entity.User, error)); ok {
		return rf(ctx, userFromReq)
	}
	if rf, ok := ret.Get(0).(func(context.Context, entity.User) entity.User); ok {
		r0 = rf(ctx, userFromReq)
	} else {
		r0 = ret.Get(0).(entity.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, entity.User) error); ok {
		r1 = rf(ctx, userFromReq)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserLogin provides a mock function with given fields: ctx, user
func (_m *Usecase) UserLogin(ctx context.Context, user entity.User) error {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for UserLogin")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.User) error); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserRegister provides a mock function with given fields: ctx, user
func (_m *Usecase) UserRegister(ctx context.Context, user entity.User) (int64, error) {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for UserRegister")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.User) (int64, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, entity.User) int64); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, entity.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WithdrawBalance provides a mock function with given fields: ctx, balance, user
func (_m *Usecase) WithdrawBalance(ctx context.Context, balance entity.BalanceUpdate, user entity.User) error {
	ret := _m.Called(ctx, balance, user)

	if len(ret) == 0 {
		panic("no return value specified for WithdrawBalance")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.BalanceUpdate, entity.User) error); ok {
		r0 = rf(ctx, balance, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Withdrawals provides a mock function with given fields: ctx, user
func (_m *Usecase) Withdrawals(ctx context.Context, user entity.User) ([]entity.BalanceUpdate, error) {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for Withdrawals")
	}

	var r0 []entity.BalanceUpdate
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.User) ([]entity.BalanceUpdate, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, entity.User) []entity.BalanceUpdate); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.BalanceUpdate)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, entity.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUsecase creates a new instance of Usecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUsecase(t interface {
	mock.TestingT
	Cleanup(func())
}) *Usecase {
	mock := &Usecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
