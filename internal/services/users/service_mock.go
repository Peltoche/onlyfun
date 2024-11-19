// Code generated by mockery v2.46.0. DO NOT EDIT.

package users

import (
	context "context"

	secret "github.com/Peltoche/onlyfun/internal/tools/secret"
	mock "github.com/stretchr/testify/mock"

	sqlstorage "github.com/Peltoche/onlyfun/internal/tools/sqlstorage"

	uuid "github.com/Peltoche/onlyfun/internal/tools/uuid"
)

// MockService is an autogenerated mock type for the Service type
type MockService struct {
	mock.Mock
}

// AddToDeletion provides a mock function with given fields: ctx, userID
func (_m *MockService) AddToDeletion(ctx context.Context, userID uuid.UUID) error {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for AddToDeletion")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Authenticate provides a mock function with given fields: ctx, username, password
func (_m *MockService) Authenticate(ctx context.Context, username string, password secret.Text) (*User, error) {
	ret := _m.Called(ctx, username, password)

	if len(ret) == 0 {
		panic("no return value specified for Authenticate")
	}

	var r0 *User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, secret.Text) (*User, error)); ok {
		return rf(ctx, username, password)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, secret.Text) *User); ok {
		r0 = rf(ctx, username, password)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, secret.Text) error); ok {
		r1 = rf(ctx, username, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Bootstrap provides a mock function with given fields: ctx, cmd
func (_m *MockService) Bootstrap(ctx context.Context, cmd *BootstrapCmd) (*User, error) {
	ret := _m.Called(ctx, cmd)

	if len(ret) == 0 {
		panic("no return value specified for Bootstrap")
	}

	var r0 *User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *BootstrapCmd) (*User, error)); ok {
		return rf(ctx, cmd)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *BootstrapCmd) *User); ok {
		r0 = rf(ctx, cmd)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *BootstrapCmd) error); ok {
		r1 = rf(ctx, cmd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Create provides a mock function with given fields: ctx, user
func (_m *MockService) Create(ctx context.Context, user *CreateCmd) (*User, error) {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *CreateCmd) (*User, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *CreateCmd) *User); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *CreateCmd) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields: ctx, paginateCmd
func (_m *MockService) GetAll(ctx context.Context, paginateCmd *sqlstorage.PaginateCmd) ([]User, error) {
	ret := _m.Called(ctx, paginateCmd)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 []User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *sqlstorage.PaginateCmd) ([]User, error)); ok {
		return rf(ctx, paginateCmd)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *sqlstorage.PaginateCmd) []User); ok {
		r0 = rf(ctx, paginateCmd)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *sqlstorage.PaginateCmd) error); ok {
		r1 = rf(ctx, paginateCmd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllWithStatus provides a mock function with given fields: ctx, status, cmd
func (_m *MockService) GetAllWithStatus(ctx context.Context, status Status, cmd *sqlstorage.PaginateCmd) ([]User, error) {
	ret := _m.Called(ctx, status, cmd)

	if len(ret) == 0 {
		panic("no return value specified for GetAllWithStatus")
	}

	var r0 []User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, Status, *sqlstorage.PaginateCmd) ([]User, error)); ok {
		return rf(ctx, status, cmd)
	}
	if rf, ok := ret.Get(0).(func(context.Context, Status, *sqlstorage.PaginateCmd) []User); ok {
		r0 = rf(ctx, status, cmd)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, Status, *sqlstorage.PaginateCmd) error); ok {
		r1 = rf(ctx, status, cmd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: ctx, userID
func (_m *MockService) GetByID(ctx context.Context, userID uuid.UUID) (*User, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 *User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*User, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *User); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HardDelete provides a mock function with given fields: ctx, userID
func (_m *MockService) HardDelete(ctx context.Context, userID uuid.UUID) error {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for HardDelete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateUserPassword provides a mock function with given fields: ctx, cmd
func (_m *MockService) UpdateUserPassword(ctx context.Context, cmd *UpdatePasswordCmd) error {
	ret := _m.Called(ctx, cmd)

	if len(ret) == 0 {
		panic("no return value specified for UpdateUserPassword")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *UpdatePasswordCmd) error); ok {
		r0 = rf(ctx, cmd)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockService creates a new instance of MockService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockService {
	mock := &MockService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
