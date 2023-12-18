// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	entity "github.com/korovindenis/go-market/internal/domain/entity"

	mock "github.com/stretchr/testify/mock"
)

// Auth is an autogenerated mock type for the auth type
type Auth struct {
	mock.Mock
}

// GenerateToken provides a mock function with given fields: user
func (_m *Auth) GenerateToken(user entity.User) (string, error) {
	ret := _m.Called(user)

	if len(ret) == 0 {
		panic("no return value specified for GenerateToken")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(entity.User) (string, error)); ok {
		return rf(user)
	}
	if rf, ok := ret.Get(0).(func(entity.User) string); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(entity.User) error); ok {
		r1 = rf(user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewAuth creates a new instance of Auth. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAuth(t interface {
	mock.TestingT
	Cleanup(func())
}) *Auth {
	mock := &Auth{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}