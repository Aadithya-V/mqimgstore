// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// IMsgB is an autogenerated mock type for the IMsgB type
type IMsgB struct {
	mock.Mock
}

// Publish provides a mock function with given fields: ctx, queueName, body
func (_m *IMsgB) Publish(ctx context.Context, queueName string, body []byte) error {
	ret := _m.Called(ctx, queueName, body)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, []byte) error); ok {
		r0 = rf(ctx, queueName, body)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewIMsgB creates a new instance of IMsgB. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIMsgB(t interface {
	mock.TestingT
	Cleanup(func())
}) *IMsgB {
	mock := &IMsgB{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
