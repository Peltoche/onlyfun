// Code generated by mockery v2.46.0. DO NOT EDIT.

package response

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// Mock is an autogenerated mock type for the Writer type
type Mock struct {
	mock.Mock
}

// WriteJSON provides a mock function with given fields: w, r, statusCode, res
func (_m *Mock) WriteJSON(w http.ResponseWriter, r *http.Request, statusCode int, res interface{}) {
	_m.Called(w, r, statusCode, res)
}

// WriteJSONError provides a mock function with given fields: w, r, err
func (_m *Mock) WriteJSONError(w http.ResponseWriter, r *http.Request, err error) {
	_m.Called(w, r, err)
}

// NewMock creates a new instance of Mock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *Mock {
	mock := &Mock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
