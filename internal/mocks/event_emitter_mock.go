// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	context "context"

	event "github.com/douglasdgoulart/video-editor-api/pkg/event"

	mock "github.com/stretchr/testify/mock"
)

// EventEmitterMock is an autogenerated mock type for the EventEmitter type
type EventEmitterMock struct {
	mock.Mock
}

type EventEmitterMock_Expecter struct {
	mock *mock.Mock
}

func (_m *EventEmitterMock) EXPECT() *EventEmitterMock_Expecter {
	return &EventEmitterMock_Expecter{mock: &_m.Mock}
}

// Send provides a mock function with given fields: ctx, _a1
func (_m *EventEmitterMock) Send(ctx context.Context, _a1 event.Event) error {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Send")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, event.Event) error); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EventEmitterMock_Send_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Send'
type EventEmitterMock_Send_Call struct {
	*mock.Call
}

// Send is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 event.Event
func (_e *EventEmitterMock_Expecter) Send(ctx interface{}, _a1 interface{}) *EventEmitterMock_Send_Call {
	return &EventEmitterMock_Send_Call{Call: _e.mock.On("Send", ctx, _a1)}
}

func (_c *EventEmitterMock_Send_Call) Run(run func(ctx context.Context, _a1 event.Event)) *EventEmitterMock_Send_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(event.Event))
	})
	return _c
}

func (_c *EventEmitterMock_Send_Call) Return(_a0 error) *EventEmitterMock_Send_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *EventEmitterMock_Send_Call) RunAndReturn(run func(context.Context, event.Event) error) *EventEmitterMock_Send_Call {
	_c.Call.Return(run)
	return _c
}

// NewEventEmitterMock creates a new instance of EventEmitterMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewEventEmitterMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *EventEmitterMock {
	mock := &EventEmitterMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
