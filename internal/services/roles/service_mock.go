// Code generated by mockery v2.46.0. DO NOT EDIT.

package roles

import mock "github.com/stretchr/testify/mock"

// MockService is an autogenerated mock type for the Service type
type MockService struct {
	mock.Mock
}

// IsRoleAuthorized provides a mock function with given fields: roleName, askedPerm
func (_m *MockService) IsRoleAuthorized(roleName string, askedPerm Permission) bool {
	ret := _m.Called(roleName, askedPerm)

	if len(ret) == 0 {
		panic("no return value specified for IsRoleAuthorized")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, Permission) bool); ok {
		r0 = rf(roleName, askedPerm)
	} else {
		r0 = ret.Get(0).(bool)
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
