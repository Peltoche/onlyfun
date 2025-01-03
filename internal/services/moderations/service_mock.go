// Code generated by mockery v2.46.0. DO NOT EDIT.

package moderations

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockService is an autogenerated mock type for the Service type
type MockService struct {
	mock.Mock
}

// ModeratePost provides a mock function with given fields: ctx, cmd
func (_m *MockService) ModeratePost(ctx context.Context, cmd *PostModerationCmd) (*Moderation, error) {
	ret := _m.Called(ctx, cmd)

	if len(ret) == 0 {
		panic("no return value specified for ModeratePost")
	}

	var r0 *Moderation
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *PostModerationCmd) (*Moderation, error)); ok {
		return rf(ctx, cmd)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *PostModerationCmd) *Moderation); ok {
		r0 = rf(ctx, cmd)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Moderation)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *PostModerationCmd) error); ok {
		r1 = rf(ctx, cmd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
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
