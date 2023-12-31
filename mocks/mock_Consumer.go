// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/erupshis/effective_mobile/internal/msgbroker (interfaces: Consumer)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	msgbroker "github.com/erupshis/effective_mobile/internal/msgbroker"
	gomock "github.com/golang/mock/gomock"
)

// MockConsumer is a mock of Consumer interface.
type MockConsumer struct {
	ctrl     *gomock.Controller
	recorder *MockConsumerMockRecorder
}

// MockConsumerMockRecorder is the mock recorder for MockConsumer.
type MockConsumerMockRecorder struct {
	mock *MockConsumer
}

// NewMockConsumer creates a new mock instance.
func NewMockConsumer(ctrl *gomock.Controller) *MockConsumer {
	mock := &MockConsumer{ctrl: ctrl}
	mock.recorder = &MockConsumerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConsumer) EXPECT() *MockConsumerMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockConsumer) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockConsumerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockConsumer)(nil).Close))
}

// Listen mocks base method.
func (m *MockConsumer) Listen(arg0 context.Context, arg1 chan<- msgbroker.Message) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Listen", arg0, arg1)
}

// Listen indicates an expected call of Listen.
func (mr *MockConsumerMockRecorder) Listen(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Listen", reflect.TypeOf((*MockConsumer)(nil).Listen), arg0, arg1)
}

// ReadMessage mocks base method.
func (m *MockConsumer) ReadMessage(arg0 context.Context) (msgbroker.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadMessage", arg0)
	ret0, _ := ret[0].(msgbroker.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadMessage indicates an expected call of ReadMessage.
func (mr *MockConsumerMockRecorder) ReadMessage(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadMessage", reflect.TypeOf((*MockConsumer)(nil).ReadMessage), arg0)
}
